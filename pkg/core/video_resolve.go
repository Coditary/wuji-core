package core

import (
	"math"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
	ffmpegdrv "github.com/coditary/wuji-core/pkg/driver/ffmpeg"
	localdrv "github.com/coditary/wuji-core/pkg/driver/local"
)

const (
	rifeDriverID       = "rife"
	realesrganDriverID = "realesrgan"
	waifu2xDriverID    = "waifu2x"
	chainnerDriverID   = "chainner"
)

// ResolveVideoDriverID picks the driver for a video request.
// Transform tasks prefer specialized backends when registered:
//   interpolate → rife, then ffmpeg
//   scale (AI upscale) → realesrgan / waifu2x / chainner (by scale), then ffmpeg
func (c *Core) ResolveVideoDriverID(preferredID string, req driver.VideoRequest) string {
	if preferredID != "" {
		return localdrv.NormalizeDriverID(preferredID)
	}

	task := req.TaskOrDefault()
	switch task {
	case driver.VideoTaskInterpolate:
		if d, err := c.registry.Get(rifeDriverID); err == nil && driver.HasVideoTask(d, task) {
			return rifeDriverID
		}
		if d, err := c.registry.Get(ffmpegdrv.DriverID); err == nil && driver.HasVideoTask(d, task) {
			return ffmpegdrv.DriverID
		}
	case driver.VideoTaskScale:
		for _, id := range aiScaleDriverOrder(req.Scale) {
			if d, err := c.registry.Get(id); err == nil && driver.HasVideoTask(d, task) {
				return id
			}
		}
		if d, err := c.registry.Get(ffmpegdrv.DriverID); err == nil && driver.HasVideoTask(d, task) {
			return ffmpegdrv.DriverID
		}
	}

	return c.ResolveDriverID("", capability.VideoGeneration)
}

// aiScaleDriverOrder returns upscale driver IDs to try, highest priority first.
func aiScaleDriverOrder(scale float32) []string {
	// waifu2x: scale 1 (denoise) or 2 (anime upscale)
	if scale >= 1 && scale < 2 {
		return []string{waifu2xDriverID, chainnerDriverID}
	}
	if !shouldUseAIUpscale(scale) {
		return nil
	}
	s := int(scale)
	switch s {
	case 2:
		return []string{realesrganDriverID, waifu2xDriverID, chainnerDriverID}
	default:
		return []string{realesrganDriverID, chainnerDriverID}
	}
}

func shouldUseAIUpscale(scale float32) bool {
	if scale < 2 {
		return false
	}
	f := float64(scale)
	return math.Abs(f-math.Round(f)) < 0.001
}
