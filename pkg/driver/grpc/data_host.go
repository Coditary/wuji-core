package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (s *Host) ProduceData(ctx context.Context, req *wujiv1.ProduceDataRequest) (*wujiv1.ProduceDataResponse, error) {
	dataReq, err := driver.ProduceDataRequestFromProto(req)
	if err != nil {
		return nil, err
	}
	if dataReq.Task == "" {
		shape, err := driver.RunDataTask(ctx, s.drv, dataReq)
		if err != nil {
			return nil, err
		}
		return driver.ProduceDataResponseToProto(shape)
	}
	prod, err := driver.As[driver.DataProducer](s.drv, capability.Data)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, "produce data not supported")
	}
	out, err := prod.ProduceData(ctx, dataReq)
	if err != nil {
		return nil, err
	}
	return driver.ProduceDataResponseToProto(out)
}
