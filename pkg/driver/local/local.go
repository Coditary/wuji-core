package local

import (
	"context"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/rag"
	"github.com/coditary/wuji-core/pkg/ragstore"
)

const (
	DriverID   = "local-rag"
	DriverName = "Local RAG"
)

// Driver indexes and searches documents under .wuji/rag using Ollama embeddings.
type Driver struct {
	store ragstore.Backend
}

func New() *Driver {
	return &Driver{store: ragstore.DefaultBackend}
}

func (d *Driver) ragCtx(ctx context.Context) context.Context {
	return ragstore.WithBackend(ctx, d.store)
}

func (d *Driver) Info() driver.Info {
	return driver.Info{
		ID:           DriverID,
		Name:         DriverName,
		Version:      "0.2.0",
		Description:  "Filesystem RAG store (.wuji/rag/) with Ollama embeddings.",
		Capabilities: []capability.Type{capability.RAG},
		RAGTasks:     driver.AllRAGTasks(),
		Remote:       false,
	}
}

func (d *Driver) Capabilities() []capability.Type { return d.Info().Capabilities }

func (d *Driver) Index(ctx context.Context, req driver.RAGRequest) (*driver.RAGIndexResponse, error) {
	return rag.Index(d.ragCtx(ctx), req)
}

func (d *Driver) Query(ctx context.Context, req driver.RAGRequest) (*driver.RAGQueryResponse, error) {
	return rag.Query(d.ragCtx(ctx), req)
}

func (d *Driver) ListCollections(ctx context.Context, req driver.RAGRequest) ([]driver.RAGCollectionInfo, error) {
	return rag.ListCollections(d.ragCtx(ctx), req)
}

func (d *Driver) CollectionInfo(ctx context.Context, req driver.RAGRequest) (driver.RAGCollectionInfo, error) {
	return rag.CollectionInfo(d.ragCtx(ctx), req)
}

func (d *Driver) DeleteCollection(ctx context.Context, req driver.RAGRequest) error {
	return rag.DeleteCollection(d.ragCtx(ctx), req)
}

func (d *Driver) PurgeSource(ctx context.Context, req driver.RAGRequest) error {
	return rag.PurgeSource(d.ragCtx(ctx), req)
}

func (d *Driver) ExportCollection(ctx context.Context, req driver.RAGRequest) error {
	return rag.ExportCollection(d.ragCtx(ctx), req)
}

func (d *Driver) ImportCollection(ctx context.Context, req driver.RAGRequest) error {
	return rag.ImportCollection(d.ragCtx(ctx), req)
}

func (d *Driver) RenameCollection(ctx context.Context, req driver.RAGRequest) error {
	return rag.RenameCollection(d.ragCtx(ctx), req)
}

func (d *Driver) Close() error { return nil }
