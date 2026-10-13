package grpc

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *RemoteDriver) Decide(ctx context.Context, req driver.DecisionRequest) (*driver.DecisionResponse, error) {
	pbReq, err := driver.ProduceDecisionRequestToProto(req)
	if err != nil {
		return nil, err
	}
	resp, err := d.client.ProduceDecision(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("remote produce decision: %w", err)
	}
	return driver.ProduceDecisionResponseFromProto(resp)
}

var _ driver.DecisionProducer = (*RemoteDriver)(nil)
