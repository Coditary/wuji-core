package dummy

import (
	"context"
	"fmt"
	"strings"

	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/rag"
)

func (d *Driver) Index(ctx context.Context, req driver.RAGRequest) (*driver.RAGIndexResponse, error) {
	return rag.Index(ctx, req)
}

func (d *Driver) Query(ctx context.Context, req driver.RAGRequest) (*driver.RAGQueryResponse, error) {
	return rag.Query(ctx, req)
}

func (d *Driver) Answer(ctx context.Context, req driver.RAGRequest) (*driver.RAGAnswerResponse, error) {
	q, err := d.Query(ctx, req)
	if err != nil {
		return nil, err
	}
	answer := fmt.Sprintf("Based on %d retrieved chunks: %s", len(q.Hits), strings.TrimSpace(req.Query))
	return &driver.RAGAnswerResponse{
		Answer: answer, Collection: req.Collection, Query: req.Query, Hits: q.Hits,
	}, nil
}

func (d *Driver) ListCollections(ctx context.Context, req driver.RAGRequest) ([]driver.RAGCollectionInfo, error) {
	return rag.ListCollections(ctx, req)
}

func (d *Driver) CollectionInfo(ctx context.Context, req driver.RAGRequest) (driver.RAGCollectionInfo, error) {
	return rag.CollectionInfo(ctx, req)
}

func (d *Driver) DeleteCollection(ctx context.Context, req driver.RAGRequest) error {
	return rag.DeleteCollection(ctx, req)
}

func (d *Driver) PurgeSource(ctx context.Context, req driver.RAGRequest) error {
	return rag.PurgeSource(ctx, req)
}

func (d *Driver) ExportCollection(ctx context.Context, req driver.RAGRequest) error {
	return rag.ExportCollection(ctx, req)
}

func (d *Driver) ImportCollection(ctx context.Context, req driver.RAGRequest) error {
	return rag.ImportCollection(ctx, req)
}

func (d *Driver) RenameCollection(ctx context.Context, req driver.RAGRequest) error {
	return rag.RenameCollection(ctx, req)
}
