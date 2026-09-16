package grpc

import (
	"context"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (s *Host) RunRAG(ctx context.Context, req *wujiv1.RAGRequestMsg) (*wujiv1.RAGResponseMsg, error) {
	ragReq, err := driver.RAGRequestFromProto(req)
	if err != nil {
		return nil, err
	}
	shape, err := runRAGOnHost(ctx, s.drv, ragReq)
	if err != nil {
		return nil, err
	}
	return driver.RAGShapeToProto(shape)
}
