package grpc

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *RemoteDriver) ProduceData(ctx context.Context, req driver.DataRequest) (data.DataShape, error) {
	pbReq, err := driver.ProduceDataRequestToProto(req)
	if err != nil {
		return nil, err
	}
	resp, err := d.client.ProduceData(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("remote produce data: %w", err)
	}
	return driver.ProduceDataResponseFromProto(resp)
}
