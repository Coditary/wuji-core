package driver

import (
	"context"
	"fmt"
	"strings"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/media"
)

func HasAudioTask(d Driver, task AudioTask) bool {
	for _, t := range d.Info().AudioTasks {
		if t == task {
			return true
		}
	}
	return false
}

func SupportsAudio(d Driver) bool {
	return len(d.Info().AudioTasks) > 0
}

func RunAudioTask(ctx context.Context, d Driver, req AudioRequest) (*AudioResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if !SupportsAudio(d) {
		return nil, fmt.Errorf("driver %q does not support audio tasks", d.Info().ID)
	}

	task := req.TaskOrDefault()
	if !HasAudioTask(d, task) {
		return nil, fmt.Errorf("driver %q does not support audio task %q", d.Info().ID, task)
	}

	var resp *AudioResponse
	var err error
	if audioTaskSupportsOverlap(task) && req.Overlap > 0 {
		resp, err = runAudioWithOverlap(ctx, d, req, task, "")
	} else {
		resp, err = dispatchAudioTask(ctx, d, req, task)
	}
	if err != nil {
		return nil, err
	}
	return finalizeAudioOutput(ctx, req, resp, "")
}

func (r AudioRequest) validateOverlap() error {
	if r.Overlap <= 0 {
		return nil
	}
	task := r.TaskOrDefault()
	if !audioTaskSupportsOverlap(task) {
		return fmt.Errorf("--overlap is only supported for text-to-music and melody-to-music, not %q", task)
	}
	if r.Duration > 0 && r.Overlap >= r.Duration {
		return fmt.Errorf("--overlap must be less than --duration")
	}
	segment := DefaultAudioSegmentDuration
	if r.Duration > 0 && r.Duration < segment {
		segment = r.Duration
	}
	if r.Overlap >= segment {
		return fmt.Errorf("--overlap must be less than %g seconds (segment length)", segment)
	}
	return nil
}

func (r AudioRequest) validateFormat() error {
	if strings.TrimSpace(r.Format) == "" {
		return nil
	}
	_, err := media.ParseAudioFormat(r.Format)
	return err
}

func AudioCapabilitiesFromTasks(tasks []AudioTask) []capability.Type {
	if len(tasks) == 0 {
		return nil
	}
	return []capability.Type{capability.AudioGeneration}
}
