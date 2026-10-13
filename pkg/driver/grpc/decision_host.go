package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (s *Host) ProduceDecision(ctx context.Context, req *wujiv1.ProduceDecisionRequest) (*wujiv1.ProduceDecisionResponse, error) {
	decisionReq, err := driver.ProduceDecisionRequestFromProto(req)
	if err != nil {
		return nil, err
	}
	prod, err := driver.As[driver.DecisionProducer](s.drv, capability.Decision)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, "produce decision not supported")
	}
	out, err := prod.Decide(ctx, decisionReq)
	if err != nil {
		return nil, err
	}
	return driver.ProduceDecisionResponseToProto(out)
}
