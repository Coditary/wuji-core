package ffmpeg

import (
	"context"
	"fmt"
	"strings"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/media"
)

const (
	DriverID   = "ffmpeg"
	DriverName = "FFmpeg Driver"
)

// Driver provides video2audio extraction and ffmpeg-based video transforms.
type Driver struct {
	FFmpegBin string
}

func New() *Driver {
	return &Driver{}
}

func (d *Driver) Info() driver.Info {
	return driver.Info{
		ID:          DriverID,
		Name:        DriverName,
		Version:     "0.2.0",
		Description: "Video transforms via ffmpeg: audio extraction (video2audio), frame interpolation, and scaling.",
		Capabilities: []capability.Type{
			capability.Video2Audio,
			capability.VideoGeneration,
		},
		VideoTasks: []driver.VideoTask{
			driver.VideoTaskInterpolate,
			driver.VideoTaskScale,
		},
	}
}

func (d *Driver) Capabilities() []capability.Type {
	return d.Info().Capabilities
}

func (d *Driver) VideoToAudio(ctx context.Context, req driver.Video2AudioRequest) (*driver.Video2AudioResponse, error) {
	videoPath := strings.TrimSpace(req.VideoPath)
	if videoPath == "" {
		return nil, fmt.Errorf("video path is required for video2audio")
	}
	if !media.IsVideoContainer(videoPath) {
		return nil, fmt.Errorf("%q is not a supported video container for video2audio", videoPath)
	}

	audioPath, err := media.ExtractAudioFromVideo(ctx, d.FFmpegBin, videoPath)
	if err != nil {
		return nil, err
	}
	return &driver.Video2AudioResponse{AudioPath: audioPath, Temporary: true}, nil
}

func (d *Driver) InterpolateVideo(ctx context.Context, req driver.VideoInterpolateRequest) (*driver.VideoResponse, error) {
	videoPath := strings.TrimSpace(req.InputVideoPath)
	if videoPath == "" {
		return nil, fmt.Errorf("input video path is required for interpolation")
	}
	path, meta, err := media.InterpolateVideo(ctx, d.FFmpegBin, videoPath, req.TargetFPS)
	if err != nil {
		return nil, err
	}
	return &driver.VideoResponse{
		Path:     path,
		Duration: meta.Duration,
		Frames:   meta.Frames,
		FPS:      meta.FPS,
	}, nil
}

func (d *Driver) UpscaleVideo(ctx context.Context, req driver.VideoUpscaleRequest) (*driver.VideoResponse, error) {
	videoPath := strings.TrimSpace(req.InputVideoPath)
	if videoPath == "" {
		return nil, fmt.Errorf("input video path is required for scale")
	}
	path, meta, err := media.ScaleVideo(ctx, d.FFmpegBin, videoPath, req.Scale)
	if err != nil {
		return nil, err
	}
	return &driver.VideoResponse{
		Path:     path,
		Duration: meta.Duration,
		Frames:   meta.Frames,
		FPS:      meta.FPS,
	}, nil
}

func (d *Driver) Close() error {
	return nil
}
