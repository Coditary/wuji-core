package core

import (
	"context"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
)

// ProduceDecision runs typed decision questions through the selected driver.
func (c *Core) ProduceDecision(ctx context.Context, driverID string, req driver.DecisionRequest) (*driver.DecisionResponse, error) {
	if err := req.ResolveQuestions(); err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return scheduleGenerate(c, ctx, driverID, req.Model, capability.Decision, func(ctx context.Context) (*driver.DecisionResponse, error) {
		d, err := c.resolveForCapability(driverID, capability.Decision)
		if err != nil {
			return nil, err
		}
		return driver.RunDecision(ctx, d, req)
	})
}
