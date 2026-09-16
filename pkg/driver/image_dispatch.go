package driver

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/capability"
)

// HasImageTask reports whether a driver advertises support for an image task.
func HasImageTask(d Driver, task ImageTask) bool {
	for _, t := range d.Info().ImageTasks {
		if t == task {
			return true
		}
	}
	return false
}

// SupportsImage reports whether a driver supports any image task.
func SupportsImage(d Driver) bool {
	return len(d.Info().ImageTasks) > 0
}

// RunImageTask dispatches a CLI union request to the task-specific driver interface.
func RunImageTask(ctx context.Context, d Driver, req ImageRequest) (*ImageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if !SupportsImage(d) {
		return nil, fmt.Errorf("driver %q does not support image tasks", d.Info().ID)
	}

	task := req.TaskOrDefault()

	if task == ImageTaskUpscale {
		return RunImageScale(ctx, ScaleDrivers{Image: d, Upscale: d, Downscale: d, Scale: d}, req)
	}

	task = ResolveImageTaskForDriver(d, task)

	if !HasImageTask(d, task) {
		return nil, fmt.Errorf("driver %q does not support image task %q", d.Info().ID, req.TaskOrDefault())
	}

	switch task {
	case ImageTaskGenerate:
		gen, ok := d.(ImageGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement ImageGenerator", d.Info().ID)
		}
		return gen.GenerateImage(ctx, req.ToGenerateRequest())
	case ImageTaskImg2Img:
		gen, ok := d.(ImageToImageGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement ImageToImageGenerator", d.Info().ID)
		}
		return gen.ImageToImage(ctx, req.ToImageToImageRequest())
	case ImageTaskInpaint:
		gen, ok := d.(ImageInpaintGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement ImageInpaintGenerator", d.Info().ID)
		}
		return gen.InpaintImage(ctx, req.ToInpaintRequest())
	case ImageTaskEdit:
		gen, ok := d.(ImageEditGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement ImageEditGenerator", d.Info().ID)
		}
		return gen.EditImage(ctx, req.ToEditRequest())
	case ImageTaskControlNet:
		gen, ok := d.(ImageControlNetGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement ImageControlNetGenerator", d.Info().ID)
		}
		return gen.ControlNetImage(ctx, req.ToControlNetRequest())
	case ImageTaskDepth2Img:
		gen, ok := d.(ImageDepthToImageGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement ImageDepthToImageGenerator", d.Info().ID)
		}
		return gen.DepthToImage(ctx, req.ToDepthToImageRequest())
	case ImageTaskVariation:
		gen, ok := d.(ImageVariationGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement ImageVariationGenerator", d.Info().ID)
		}
		return gen.ImageVariation(ctx, req.ToVariationRequest())
	case ImageTaskStyleTransfer:
		gen, ok := d.(ImageStyleTransferGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement ImageStyleTransferGenerator", d.Info().ID)
		}
		return gen.StyleTransferImage(ctx, req.ToStyleTransferRequest())
	case ImageTaskSprite:
		gen, ok := d.(ImageSpriteGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement ImageSpriteGenerator", d.Info().ID)
		}
		return gen.SpriteImage(ctx, req.ToSpriteRequest())
	default:
		return nil, fmt.Errorf("unsupported image task %q", task)
	}
}

// ResolveImageTaskForDriver maps edit to img2img when the driver lacks a dedicated edit backend.
func ResolveImageTaskForDriver(d Driver, task ImageTask) ImageTask {
	if task != ImageTaskEdit {
		return task
	}
	if HasImageTask(d, ImageTaskEdit) {
		if _, ok := d.(ImageEditGenerator); ok {
			return ImageTaskEdit
		}
	}
	if HasImageTask(d, ImageTaskImg2Img) {
		if _, ok := d.(ImageToImageGenerator); ok {
			return ImageTaskImg2Img
		}
	}
	if _, ok := d.(ImageToImageGenerator); ok {
		return ImageTaskImg2Img
	}
	return task
}

// ImageCapabilitiesFromTasks returns capability.ImageGeneration when any image task is supported.
func ImageCapabilitiesFromTasks(tasks []ImageTask) []capability.Type {
	if len(tasks) == 0 {
		return nil
	}
	return []capability.Type{capability.ImageGeneration}
}
