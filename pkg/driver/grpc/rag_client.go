package grpc

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *RemoteDriver) ProduceRAG(ctx context.Context, req driver.RAGRequest) (data.DataShape, error) {
	resp, err := d.client.RunRAG(ctx, driver.RAGRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote run rag: %w", err)
	}
	return driver.RAGShapeFromProto(resp)
}
