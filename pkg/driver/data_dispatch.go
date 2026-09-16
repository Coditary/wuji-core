package driver

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/data"
)

// RunDataTask dispatches a data request to a driver or returns input-only transforms.
func RunDataTask(ctx context.Context, d Driver, req DataRequest) (data.DataShape, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	target := InferOutputShape(req)
	if req.Task == "" {
		if req.Input == nil {
			return nil, fmt.Errorf("structured input is required when no task flag is set")
		}
		return coerceOutputShape(req.Input, target)
	}
	if d == nil {
		return nil, fmt.Errorf("no driver available for data task %q", req.Task)
	}
	prod, err := As[DataProducer](d, capability.Data)
	if err != nil {
		return nil, err
	}
	out, err := prod.ProduceData(ctx, req)
	if err != nil {
		return nil, err
	}
	return coerceOutputShape(out, target)
}

func coerceOutputShape(shape data.DataShape, target data.ShapeKind) (data.DataShape, error) {
	if target == "" || target == "auto" {
		return shape, nil
	}
	return data.ConvertShape(shape, target)
}
