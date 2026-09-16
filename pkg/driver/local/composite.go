package local

import (
	"context"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/embed"
)

// CompositeHost combines local-rag storage with a data-capable embedder over gRPC.
type CompositeHost struct {
	RAG   *Driver
	embed driver.DataProducer
}

// NewCompositeHost wraps local-rag and an embedding data producer.
func NewCompositeHost(data driver.DataProducer) *CompositeHost {
	return &CompositeHost{
		RAG:   New(),
		embed: data,
	}
}

// NewCompositeFromConfig wraps local-rag with Ollama embeddings from Wuji config.
func NewCompositeFromConfig(cfg *config.Config) *CompositeHost {
	ragCfg := config.LocalRAGConfig{}
	if cfg != nil {
		ragCfg = cfg.ResolvedLocalRAG()
	} else {
		ragCfg = config.DefaultLocalRAGConfig()
	}
	emb := embed.NewOllama(ragCfg.OllamaAPI, ragCfg.DefaultEmbedModel, DriverID)
	return NewCompositeHost(emb)
}

func (c *CompositeHost) ragCtx(ctx context.Context) context.Context {
	return c.RAG.ragCtx(ctx)
}

func (c *CompositeHost) Info() driver.Info {
	info := c.RAG.Info()
	info.Description = "Filesystem RAG (.wuji/rag/) with Ollama embeddings."
	info.Capabilities = append([]capability.Type{capability.Data}, info.Capabilities...)
	info.DataTasks = driver.AllDataTasks()
	return info
}

func (c *CompositeHost) Capabilities() []capability.Type { return c.Info().Capabilities }

func (c *CompositeHost) Close() error { return nil }

func (c *CompositeHost) ProduceData(ctx context.Context, req driver.DataRequest) (data.DataShape, error) {
	return c.embed.ProduceData(ctx, req)
}

func (c *CompositeHost) Index(ctx context.Context, req driver.RAGRequest) (*driver.RAGIndexResponse, error) {
	return c.RAG.Index(c.ragCtx(ctx), req)
}

func (c *CompositeHost) Query(ctx context.Context, req driver.RAGRequest) (*driver.RAGQueryResponse, error) {
	return c.RAG.Query(c.ragCtx(ctx), req)
}

func (c *CompositeHost) ListCollections(ctx context.Context, req driver.RAGRequest) ([]driver.RAGCollectionInfo, error) {
	return c.RAG.ListCollections(c.ragCtx(ctx), req)
}

func (c *CompositeHost) CollectionInfo(ctx context.Context, req driver.RAGRequest) (driver.RAGCollectionInfo, error) {
	return c.RAG.CollectionInfo(c.ragCtx(ctx), req)
}

func (c *CompositeHost) DeleteCollection(ctx context.Context, req driver.RAGRequest) error {
	return c.RAG.DeleteCollection(c.ragCtx(ctx), req)
}

func (c *CompositeHost) PurgeSource(ctx context.Context, req driver.RAGRequest) error {
	return c.RAG.PurgeSource(c.ragCtx(ctx), req)
}

func (c *CompositeHost) ExportCollection(ctx context.Context, req driver.RAGRequest) error {
	return c.RAG.ExportCollection(c.ragCtx(ctx), req)
}

func (c *CompositeHost) ImportCollection(ctx context.Context, req driver.RAGRequest) error {
	return c.RAG.ImportCollection(c.ragCtx(ctx), req)
}

func (c *CompositeHost) RenameCollection(ctx context.Context, req driver.RAGRequest) error {
	return c.RAG.RenameCollection(c.ragCtx(ctx), req)
}

// Ensure CompositeHost implements DataProducer at compile time.
var _ driver.DataProducer = (*CompositeHost)(nil)
