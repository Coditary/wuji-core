package driver

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/capability"
)

func HasVoiceTask(d Driver, task VoiceTask) bool {
	for _, t := range d.Info().VoiceTasks {
		if t == task {
			return true
		}
	}
	return false
}

func SupportsVoice(d Driver) bool {
	return len(d.Info().VoiceTasks) > 0
}

func RunVoiceTask(ctx context.Context, d Driver, req VoiceRequest) (*VoiceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if !SupportsVoice(d) {
		return nil, fmt.Errorf("driver %q does not support voice tasks", d.Info().ID)
	}

	task := req.TaskOrDefault()
	if !HasVoiceTask(d, task) {
		return nil, fmt.Errorf("driver %q does not support voice task %q", d.Info().ID, task)
	}

	switch task {
	case VoiceTaskZeroShot:
		gen, ok := d.(VoiceCreator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement VoiceCreator", d.Info().ID)
		}
		return gen.CreateVoice(ctx, req.ToCreateRequest())
	case VoiceTaskConversion:
		gen, ok := d.(VoiceConverter)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement VoiceConverter", d.Info().ID)
		}
		return gen.ConvertVoice(ctx, req.ToConvertRequest())
	default:
		return nil, fmt.Errorf("unsupported voice task %q", task)
	}
}

func VoiceCapabilitiesFromTasks(tasks []VoiceTask) []capability.Type {
	if len(tasks) == 0 {
		return nil
	}
	return []capability.Type{capability.VoiceCloning}
}
