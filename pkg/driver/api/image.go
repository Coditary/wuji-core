package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *Driver) runImage(ctx context.Context, req driver.ImageRequest) (*driver.ImageResponse, error) {
	if d.entry.Image == nil {
		return nil, fmt.Errorf("driver %q does not support image", d.id)
	}
	task := req.TaskOrDefault()
	spec := config.ResolveAPICapabilitySpec(d.entry.Image, string(task))
	proto := strings.ToLower(strings.TrimSpace(spec.Protocol))
	if proto == "openai-images" || proto == "openai-image" {
		if task == driver.ImageTaskGenerate {
			return generateImageOpenAI(ctx, spec, req)
		}
	}
	raw, err := doHTTP(d.scoped(ctx), spec, imageHTTPVars(req, spec))
	if err != nil {
		return nil, err
	}
	return mediaResponseFromHTTP(spec, raw, "image")
}

func (d *Driver) GenerateImage(ctx context.Context, req driver.ImageGenerateRequest) (*driver.ImageResponse, error) {
	return d.runImage(ctx, unionImageFromGenerate(req))
}

func (d *Driver) ImageToImage(ctx context.Context, req driver.ImageToImageRequest) (*driver.ImageResponse, error) {
	return d.runImage(ctx, unionImageFromImg2Img(req))
}

func (d *Driver) InpaintImage(ctx context.Context, req driver.ImageInpaintRequest) (*driver.ImageResponse, error) {
	return d.runImage(ctx, unionImageFromInpaint(req))
}

func (d *Driver) UpscaleImage(ctx context.Context, req driver.ImageUpscaleRequest) (*driver.ImageResponse, error) {
	return d.runImage(ctx, unionImageFromUpscale(req, driver.ImageTaskUpscale))
}

func (d *Driver) DownscaleImage(ctx context.Context, req driver.ImageDownscaleRequest) (*driver.ImageResponse, error) {
	return d.runImage(ctx, unionImageFromDownscale(req))
}

func (d *Driver) ScaleImage(ctx context.Context, req driver.ImageScaleRequest) (*driver.ImageResponse, error) {
	return d.runImage(ctx, unionImageFromScale(req))
}

func (d *Driver) EditImage(ctx context.Context, req driver.ImageEditRequest) (*driver.ImageResponse, error) {
	return d.runImage(ctx, unionImageFromEdit(req))
}

func (d *Driver) ControlNetImage(ctx context.Context, req driver.ImageControlNetRequest) (*driver.ImageResponse, error) {
	return d.runImage(ctx, unionImageFromControlNet(req))
}

func (d *Driver) DepthToImage(ctx context.Context, req driver.ImageDepthToImageRequest) (*driver.ImageResponse, error) {
	return d.runImage(ctx, unionImageFromDepth(req))
}

func (d *Driver) ImageVariation(ctx context.Context, req driver.ImageVariationRequest) (*driver.ImageResponse, error) {
	return d.runImage(ctx, unionImageFromVariation(req))
}

func (d *Driver) StyleTransferImage(ctx context.Context, req driver.ImageStyleTransferRequest) (*driver.ImageResponse, error) {
	return d.runImage(ctx, unionImageFromStyleTransfer(req))
}

func (d *Driver) SpriteImage(ctx context.Context, req driver.ImageSpriteRequest) (*driver.ImageResponse, error) {
	return d.runImage(ctx, unionImageFromSprite(req))
}

func unionImageFromGenerate(req driver.ImageGenerateRequest) driver.ImageRequest {
	return driver.ImageRequest{
		Task: driver.ImageTaskGenerate, Prompt: req.Prompt,
		NegativePrompt: req.NegativePrompt, Model: req.Model, Width: req.Width, Height: req.Height,
		Steps: req.Steps, Sampler: req.Sampler, CFGScale: req.CFGScale, BatchSize: req.BatchSize,
		BatchCount: req.BatchCount, Seed: req.Seed, LoRAs: req.LoRAs, ControlUnits: req.ControlUnits,
		ControlMode: req.ControlMode, Mode: req.Mode,
	}
}

func unionImageFromImg2Img(req driver.ImageToImageRequest) driver.ImageRequest {
	r := unionImageFromGenerate(driver.ImageGenerateRequest{Prompt: req.Prompt, ImageCommonParams: req.ImageCommonParams})
	r.Task = driver.ImageTaskImg2Img
	r.InitImagePath = req.InitImagePath
	r.DenoisingStrength = req.DenoisingStrength
	return r
}

func unionImageFromInpaint(req driver.ImageInpaintRequest) driver.ImageRequest {
	r := unionImageFromGenerate(driver.ImageGenerateRequest{Prompt: req.Prompt, ImageCommonParams: req.ImageCommonParams})
	r.Task = driver.ImageTaskInpaint
	r.InitImagePath = req.InitImagePath
	r.MaskImagePath = req.MaskImagePath
	r.DenoisingStrength = req.DenoisingStrength
	return r
}

func unionImageFromUpscale(req driver.ImageUpscaleRequest, task driver.ImageTask) driver.ImageRequest {
	return driver.ImageRequest{
		Task: task, InitImagePath: req.InitImagePath, Model: req.Model, Scale: req.Scale, LoRAs: req.LoRAs,
	}
}

func unionImageFromDownscale(req driver.ImageDownscaleRequest) driver.ImageRequest {
	return unionImageFromUpscale(driver.ImageUpscaleRequest{
		InitImagePath: req.InitImagePath, Model: req.Model, Scale: req.Scale, LoRAs: req.LoRAs,
	}, driver.ImageTaskUpscale)
}

func unionImageFromScale(req driver.ImageScaleRequest) driver.ImageRequest {
	return unionImageFromUpscale(driver.ImageUpscaleRequest{
		InitImagePath: req.InitImagePath, Model: req.Model, Scale: req.Scale, LoRAs: req.LoRAs,
	}, driver.ImageTaskUpscale)
}

func unionImageFromEdit(req driver.ImageEditRequest) driver.ImageRequest {
	r := unionImageFromGenerate(driver.ImageGenerateRequest{Prompt: req.Prompt, ImageCommonParams: req.ImageCommonParams})
	r.Task = driver.ImageTaskEdit
	r.InitImagePath = req.InitImagePath
	r.DenoisingStrength = req.DenoisingStrength
	return r
}

func unionImageFromControlNet(req driver.ImageControlNetRequest) driver.ImageRequest {
	r := unionImageFromGenerate(driver.ImageGenerateRequest{Prompt: req.Prompt, ImageCommonParams: req.ImageCommonParams})
	r.Task = driver.ImageTaskControlNet
	r.ControlUnits = req.ControlUnits
	r.ControlMode = req.ControlMode
	return r
}

func unionImageFromDepth(req driver.ImageDepthToImageRequest) driver.ImageRequest {
	r := unionImageFromGenerate(driver.ImageGenerateRequest{Prompt: req.Prompt, ImageCommonParams: req.ImageCommonParams})
	r.Task = driver.ImageTaskDepth2Img
	r.ControlImagePath = req.DepthImagePath
	return r
}

func unionImageFromVariation(req driver.ImageVariationRequest) driver.ImageRequest {
	r := driver.ImageRequest{
		Task: driver.ImageTaskVariation, InitImagePath: req.InitImagePath,
		DenoisingStrength: req.DenoisingStrength, Model: req.Model, LoRAs: req.LoRAs,
	}
	return r
}

func unionImageFromStyleTransfer(req driver.ImageStyleTransferRequest) driver.ImageRequest {
	r := unionImageFromGenerate(driver.ImageGenerateRequest{Prompt: req.Prompt, ImageCommonParams: req.ImageCommonParams})
	r.Task = driver.ImageTaskStyleTransfer
	r.StyleImagePath = req.StyleImagePath
	r.StyleWeight = req.StyleWeight
	return r
}

func unionImageFromSprite(req driver.ImageSpriteRequest) driver.ImageRequest {
	return driver.ImageRequest{
		Task: driver.ImageTaskSprite, Prompt: req.Prompt, Model: req.Model, Mode: req.Mode,
		InitImagePath: req.InitImagePath, ReferenceImagePath: req.ReferenceImagePath,
		StyleImagePath: req.StyleImagePath, StyleWeight: req.StyleWeight,
		FrameWidth: req.FrameWidth, FrameHeight: req.FrameHeight,
		Columns: req.Columns, Rows: req.Rows, FrameCount: req.FrameCount,
		DenoisingStrength: req.DenoisingStrength,
		SpriteAction:      req.Action, SpriteView: req.View, SpriteDirections: req.Directions,
		SpriteLoop: req.Loop, SpritePadding: req.Padding, SpriteTransparent: req.Transparent,
		LoRAs: req.LoRAs,
	}
}
