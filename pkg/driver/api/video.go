package api

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *Driver) runVideo(ctx context.Context, req driver.VideoRequest) (*driver.VideoResponse, error) {
	if d.entry.Video == nil {
		return nil, fmt.Errorf("driver %q does not support video", d.id)
	}
	task := req.TaskOrDefault()
	spec := config.ResolveAPICapabilitySpec(d.entry.Video, string(task))
	raw, err := doHTTP(d.scoped(ctx), spec, videoHTTPVars(req, spec))
	if err != nil {
		return nil, err
	}
	path, err := mediaPathFromHTTP(spec, raw, "video")
	if err != nil {
		return nil, err
	}
	return &driver.VideoResponse{Path: path}, nil
}

func (d *Driver) GenerateVideo(ctx context.Context, req driver.VideoGenerateRequest) (*driver.VideoResponse, error) {
	union := driver.VideoRequest{
		Task: driver.VideoTaskGenerate, Prompt: req.Prompt, CameraControl: req.CameraControl,
		Duration: req.Duration, FPS: req.FPS, Frames: req.Frames, MotionStrength: req.MotionStrength,
		ContextLength: req.ContextLength, Sampler: req.Sampler, Scheduler: req.Scheduler,
		NegativePrompt: req.NegativePrompt, Model: req.Model, Seed: req.Seed,
	}
	return d.runVideo(ctx, union)
}

func (d *Driver) ImageToVideo(ctx context.Context, req driver.VideoImageToVideoRequest) (*driver.VideoResponse, error) {
	union := driver.VideoRequest{
		Task: driver.VideoTaskImageToVideo, Prompt: req.Prompt, InitImagePath: req.InitImagePath,
		Duration: req.Duration, FPS: req.FPS, Frames: req.Frames, MotionStrength: req.MotionStrength,
		ContextLength: req.ContextLength, Sampler: req.Sampler, Scheduler: req.Scheduler,
		NegativePrompt: req.NegativePrompt, Model: req.Model, Seed: req.Seed,
	}
	return d.runVideo(ctx, union)
}

func (d *Driver) InterpolateVideo(ctx context.Context, req driver.VideoInterpolateRequest) (*driver.VideoResponse, error) {
	union := driver.VideoRequest{
		Task: driver.VideoTaskInterpolate, InitVideoPath: req.InputVideoPath,
		FPS: req.TargetFPS, Model: req.Model,
	}
	return d.runVideo(ctx, union)
}

func (d *Driver) UpscaleVideo(ctx context.Context, req driver.VideoUpscaleRequest) (*driver.VideoResponse, error) {
	union := driver.VideoRequest{
		Task: driver.VideoTaskScale, InitVideoPath: req.InputVideoPath, Scale: req.Scale, Model: req.Model,
	}
	return d.runVideo(ctx, union)
}
