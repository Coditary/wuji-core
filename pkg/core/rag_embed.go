package core

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/rag"
)

const embedBatchSize = 32

func (c *Core) embedTexts(ctx context.Context, ragDriverID, model string, texts []string) ([][]float32, int, error) {
	if len(texts) == 0 {
		return nil, 0, nil
	}
	dataD, err := c.resolveDataDriver(ragDriverID)
	if err != nil {
		return rag.PseudoEmbedder(ctx, model, texts)
	}
	prod, err := driver.As[driver.DataProducer](dataD, capability.Data)
	if err != nil {
		return rag.PseudoEmbedder(ctx, model, texts)
	}

	var rows [][]float32
	dims := 0
	for i := 0; i < len(texts); i += embedBatchSize {
		end := i + embedBatchSize
		if end > len(texts) {
			end = len(texts)
		}
		batch := texts[i:end]
		shape, err := prod.ProduceData(ctx, driver.DataRequest{
			Task: driver.DataTaskEmbed, Model: model, Texts: batch,
		})
		if err != nil {
			return nil, 0, fmt.Errorf("embed batch: %w", err)
		}
		vecs, d, err := vectorRowsFromShape(shape)
		if err != nil {
			return nil, 0, err
		}
		if dims == 0 {
			dims = d
		} else if d != dims {
			return nil, 0, fmt.Errorf("embedding dims changed: got %d want %d", d, dims)
		}
		rows = append(rows, vecs...)
	}
	return rows, dims, nil
}

func (c *Core) resolveDataDriver(ragDriverID string) (driver.Driver, error) {
	if c.appConfig != nil {
		if id := c.appConfig.DriverForCapability(capability.Data); id != "" {
			return c.resolve(id)
		}
	}
	if ragDriverID != "" {
		if d, err := c.resolve(ragDriverID); err == nil {
			if _, err := driver.As[driver.DataProducer](d, capability.Data); err == nil {
				return d, nil
			}
		}
	}
	return c.resolve(c.ResolveDriverID("", capability.Data))
}

func vectorRowsFromShape(shape data.DataShape) ([][]float32, int, error) {
	batch, ok := shape.(*data.VectorBatch)
	if !ok || batch == nil {
		return nil, 0, fmt.Errorf("embed returned %T, want vector batch", shape)
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

func (c *Core) ragContext(ctx context.Context, preferredDriverID string) context.Context {
	return rag.WithEmbedder(ctx, func(ctx context.Context, model string, texts []string) ([][]float32, int, error) {
		return c.embedTexts(ctx, preferredDriverID, model, texts)
	})
}
