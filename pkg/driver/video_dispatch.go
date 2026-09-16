package driver

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/capability"
)

func HasVideoTask(d Driver, task VideoTask) bool {
	for _, t := range d.Info().VideoTasks {
		if t == task {
			return true
		}
	}
	return false
}

func SupportsVideo(d Driver) bool {
	return len(d.Info().VideoTasks) > 0
}

func RunVideoTask(ctx context.Context, d Driver, req VideoRequest) (*VideoResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if !SupportsVideo(d) {
		return nil, fmt.Errorf("driver %q does not support video tasks", d.Info().ID)
	}

	task := req.TaskOrDefault()
	if !HasVideoTask(d, task) {
		return nil, fmt.Errorf("driver %q does not support video task %q", d.Info().ID, task)
	}

	switch task {
	case VideoTaskGenerate:
		gen, ok := d.(VideoGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement VideoGenerator", d.Info().ID)
		}
		return gen.GenerateVideo(ctx, req.ToGenerateRequest())
	case VideoTaskImageToVideo:
		gen, ok := d.(VideoImageToVideoGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement VideoImageToVideoGenerator", d.Info().ID)
		}
		return gen.ImageToVideo(ctx, req.ToImageToVideoRequest())
	case VideoTaskInterpolate:
		gen, ok := d.(VideoInterpolateGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement VideoInterpolateGenerator", d.Info().ID)
		}
		return gen.InterpolateVideo(ctx, req.ToInterpolateRequest())
	case VideoTaskScale:
		gen, ok := d.(VideoUpscaleGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement VideoUpscaleGenerator", d.Info().ID)
		}
		return gen.UpscaleVideo(ctx, req.ToUpscaleRequest())
	default:
		return nil, fmt.Errorf("unsupported video task %q", task)
	}
}

func VideoCapabilitiesFromTasks(tasks []VideoTask) []capability.Type {
	if len(tasks) == 0 {
		return nil
	}
	return []capability.Type{capability.VideoGeneration}
}
