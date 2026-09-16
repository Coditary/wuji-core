package driver

import (
	"context"
	"math"
	"os"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/imageresize"
)

const defaultScaleOutputDir = "/tmp/wuji/scale"

type scalePlan struct {
	aiStep       int
	magickFactor float32
}

// planScale splits a requested factor into an AI upscaler step (2, 4, or 8) and a
// follow-up ImageMagick resize. For fractional targets (e.g. 2.5×) it overshoots with
// the next higher AI tier (4×) then downscales (62.5%) for better detail than upscaling
// after a smaller AI step.
func planScale(factor float32) scalePlan {
	if factor < 1 {
		return scalePlan{magickFactor: factor}
	}
	if nearOne(factor) {
		return scalePlan{magickFactor: 1}
	}

	for _, step := range []int{2, 4, 8} {
		if float32(step) >= factor-1e-4 {
			return scalePlan{aiStep: step, magickFactor: factor / float32(step)}
		}
	}

	return scalePlan{aiStep: 8, magickFactor: factor / 8}
}

func nearOne(f float32) bool {
	return math.Abs(float64(f)-1) < 1e-4
}

// RunImageScale coordinates downscale, AI upscale, and ImageMagick resize backends.
func RunImageScale(ctx context.Context, drivers ScaleDrivers, req ImageRequest) (*ImageResponse, error) {
	factor := req.Scale
	if factor <= 0 {
		factor = 4
	}

	switch ClassifyScale(factor) {
	case ScaleKindDownscale:
		out, err := runDownscale(ctx, drivers.Downscale, req, factor)
		if err != nil {
			return nil, err
		}
		return &ImageResponse{Path: out, Format: "png"}, nil
	case ScaleKindUpscaleExact:
		if nearOne(factor) {
			return &ImageResponse{Path: req.InitImagePath, Format: "png"}, nil
		}
		out, err := runUpscaleStep(ctx, drivers.Upscale, req, int(factor))
		if err != nil {
			return nil, err
		}
		return &ImageResponse{Path: out, Format: "png"}, nil
	default:
		return runHybridScale(ctx, drivers, req, factor)
	}
}

func runHybridScale(ctx context.Context, drivers ScaleDrivers, req ImageRequest, factor float32) (*ImageResponse, error) {
	if drivers.Scale != nil {
		if gen, err := As[ImageScaleGenerator](drivers.Scale, capability.ImageScale); err == nil {
			resp, err := gen.ScaleImage(ctx, ImageScaleRequest{
				InitImagePath: req.InitImagePath,
				Model:         req.Model,
				Scale:         factor,
				LoRAs:         req.LoRAs,
			})
			if err == nil && fileExists(resp.Path) {
				return resp, nil
			}
		}
	}
	return runBuiltInHybridScale(ctx, drivers, req, factor)
}

// RunBuiltInHybridScale applies the default wuji orchestration for non-integer scale factors.
func RunBuiltInHybridScale(ctx context.Context, drivers ScaleDrivers, req ImageRequest, factor float32) (*ImageResponse, error) {
	return runBuiltInHybridScale(ctx, drivers, req, factor)
}

func runBuiltInHybridScale(ctx context.Context, drivers ScaleDrivers, req ImageRequest, factor float32) (*ImageResponse, error) {
	plan := planScale(factor)
	current := req.InitImagePath
	if plan.aiStep > 0 {
		out, err := runUpscaleStep(ctx, drivers.Upscale, req, plan.aiStep)
		if err != nil {
			out, err = imageresize.ScaleImage(ctx, req.InitImagePath, defaultScaleOutputDir, factor)
			if err != nil {
				return nil, err
			}
			return &ImageResponse{Path: out, Format: "png"}, nil
		}
		current = out
	}
	if nearOne(plan.magickFactor) {
		return &ImageResponse{Path: current, Format: "png"}, nil
	}
	var out string
	var err error
	if plan.magickFactor < 1 {
		out, err = runDownscale(ctx, drivers.Downscale, imageReqAt(current, req), plan.magickFactor)
	} else {
		out, err = imageresize.ScaleImage(ctx, current, defaultScaleOutputDir, plan.magickFactor)
	}
	if err != nil {
		return nil, err
	}
	return &ImageResponse{Path: out, Format: "png"}, nil
}

func imageReqAt(path string, req ImageRequest) ImageRequest {
	req.InitImagePath = path
	return req
}

func runDownscale(ctx context.Context, d Driver, req ImageRequest, factor float32) (string, error) {
	if d != nil {
		if gen, err := As[ImageDownscaleGenerator](d, capability.ImageDownscale); err == nil {
			resp, err := gen.DownscaleImage(ctx, ImageDownscaleRequest{
				InitImagePath: req.InitImagePath,
				Model:         req.Model,
				Scale:         factor,
				LoRAs:         req.LoRAs,
			})
			if err == nil && fileExists(resp.Path) {
				return resp.Path, nil
			}
		}
	}
	return imageresize.ScaleImage(ctx, req.InitImagePath, defaultScaleOutputDir, factor)
}

func runUpscaleStep(ctx context.Context, d Driver, req ImageRequest, step int) (string, error) {
	if d != nil {
		if gen, err := As[ImageUpscaleGenerator](d, capability.ImageUpscale); err == nil {
			resp, err := gen.UpscaleImage(ctx, ImageUpscaleRequest{
				InitImagePath: req.InitImagePath,
				Model:         req.Model,
				Scale:         float32(step),
				LoRAs:         req.LoRAs,
			})
			if err == nil && fileExists(resp.Path) {
				return resp.Path, nil
			}
		}
	}
	return imageresize.ScaleImage(ctx, req.InitImagePath, defaultScaleOutputDir, float32(step))
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}
