package rag

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/ragstore"
)

type embedderKey struct{}

// Embedder produces vector embeddings for chunk texts.
type Embedder func(ctx context.Context, model string, texts []string) ([][]float32, int, error)

// WithEmbedder attaches an embedder to the context for RAG drivers.
func WithEmbedder(ctx context.Context, emb Embedder) context.Context {
	if emb == nil {
		return ctx
	}
	return context.WithValue(ctx, embedderKey{}, emb)
}

// EmbedderFrom returns the embedder from context, if any.
func EmbedderFrom(ctx context.Context) (Embedder, bool) {
	emb, ok := ctx.Value(embedderKey{}).(Embedder)
	return emb, ok && emb != nil
}

// ResolveEmbedder returns the context embedder or a pseudo fallback.
func ResolveEmbedder(ctx context.Context) Embedder {
	if emb, ok := EmbedderFrom(ctx); ok {
		return emb
	}
	return PseudoEmbedder
}

// PseudoEmbedder uses deterministic dev vectors when no data driver is available.
func PseudoEmbedder(ctx context.Context, model string, texts []string) ([][]float32, int, error) {
	_ = ctx
	_ = model
	dims := ragstore.DefaultEmbedDims()
	return ragstore.EmbedTexts(texts, dims), dims, nil
}

// DriverEmbedder embeds via a driver's data capability.
func DriverEmbedder(d driver.Driver) Embedder {
	return func(ctx context.Context, model string, texts []string) ([][]float32, int, error) {
		if len(texts) == 0 {
			return nil, 0, nil
		}
		prod, err := driver.As[driver.DataProducer](d, capability.Data)
		if err != nil {
			return PseudoEmbedder(ctx, model, texts)
		}
		const batchSize = 32
		var rows [][]float32
		dims := 0
		for i := 0; i < len(texts); i += batchSize {
			end := i + batchSize
			if end > len(texts) {
				end = len(texts)
			}
			batch := texts[i:end]
			shape, err := prod.ProduceData(ctx, driver.DataRequest{
				Task: driver.DataTaskEmbed, Model: model, Texts: batch,
			})
			if err != nil {
				return nil, 0, err
			}
			vecs, d, err := vectorsFromShape(shape)
			if err != nil {
				return nil, 0, err
			}
			if dims == 0 {
				dims = d
			} else if d != dims {
				return nil, 0, fmt.Errorf("embedding dims changed mid-batch: got %d want %d", d, dims)
			}
			rows = append(rows, vecs...)
		}
		return rows, dims, nil
	}
}

func vectorsFromShape(shape data.DataShape) ([][]float32, int, error) {
	batch, ok := shape.(*data.VectorBatch)
	if !ok || batch == nil {
		return nil, 0, fmt.Errorf("embed task returned %T, want vector batch", shape)
	}
	if batch.Dims <= 0 || batch.Count() == 0 {
		return nil, 0, fmt.Errorf("empty embedding batch")
	}
	out := make([][]float32, batch.Count())
	for i := 0; i < batch.Count(); i++ {
		row := batch.At(i)
		if len(row) != batch.Dims {
			return nil, 0, fmt.Errorf("vector row %d has %d dims, want %d", i, len(row), batch.Dims)
		}
		out[i] = append([]float32(nil), row...)
	}
	return out, batch.Dims, nil
}
