package driver

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/capability"
)

func HasMeshTask(d Driver, task MeshTask) bool {
	for _, t := range d.Info().MeshTasks {
		if t == task {
			return true
		}
	}
	return false
}

func SupportsMesh(d Driver) bool {
	return len(d.Info().MeshTasks) > 0
}

func RunMeshTask(ctx context.Context, d Driver, req MeshRequest) (*MeshResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if !SupportsMesh(d) {
		return nil, fmt.Errorf("driver %q does not support mesh tasks", d.Info().ID)
	}

	task := req.TaskOrDefault()
	if !HasMeshTask(d, task) {
		return nil, fmt.Errorf("driver %q does not support mesh task %q", d.Info().ID, task)
	}

	gen, ok := d.(MeshGenerator)
	if !ok {
		return nil, fmt.Errorf("driver %q does not implement MeshGenerator", d.Info().ID)
	}
	resp, err := gen.GenerateMesh(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp != nil && resp.Task == "" {
		resp.Task = task
	}
	if resp != nil && resp.Mode == "" {
		resp.Mode = req.ModeOrDefault()
	}
	if resp != nil && resp.Representation == "" {
		resp.Representation = req.RepresentationOrDefault()
	}
	return resp, nil
}

func MeshCapabilitiesFromTasks(tasks []MeshTask) []capability.Type {
	if len(tasks) == 0 {
		return nil
	}
	return []capability.Type{capability.Mesh}
}
