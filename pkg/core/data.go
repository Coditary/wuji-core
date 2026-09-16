package core

import (
	"context"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
)

// ProduceData runs a structured data task through the selected driver.
func (c *Core) ProduceData(ctx context.Context, driverID string, req driver.DataRequest) (data.DataShape, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if req.Task == "" {
		return driver.RunDataTask(ctx, nil, req)
	}
	return scheduleGenerate(c, ctx, driverID, req.Model, capability.Data, func(ctx context.Context) (data.DataShape, error) {
		d, err := c.resolveForCapability(driverID, capability.Data)
		if err != nil {
			return nil, err
		}
		return driver.RunDataTask(ctx, d, req)
	})
}
