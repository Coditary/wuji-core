package driver

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/capability"
)

// HasDecisionTask reports whether a driver advertises support for a decision task type.
func HasDecisionTask(d Driver, task DecisionTask) bool {
	for _, t := range d.Info().DecisionTasks {
		if t == task {
			return true
		}
	}
	return false
}

// SupportsDecision reports whether the driver exposes any decision tasks.
func SupportsDecision(d Driver) bool {
	return len(d.Info().DecisionTasks) > 0
}

// RunDecision dispatches a decision request to a driver.
func RunDecision(ctx context.Context, d Driver, req DecisionRequest) (*DecisionResponse, error) {
	if err := req.ResolveQuestions(); err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if d == nil {
		return nil, fmt.Errorf("no driver available for decision")
	}
	prod, err := As[DecisionProducer](d, capability.Decision)
	if err != nil {
		return nil, err
	}
	for _, q := range req.Questions {
		if !HasDecisionTask(d, q.Type) {
			return nil, fmt.Errorf("driver %q does not support decision task %q", d.Info().ID, q.Type)
		}
	}
	return prod.Decide(ctx, req)
}
