package grpc

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *RemoteDriver) GenerateImage(ctx context.Context, req driver.ImageGenerateRequest) (*driver.ImageResponse, error) {
	resp, err := d.client.GenerateImage(ctx, driver.ImageGenerateRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote generate image: %w", err)
	}
	return driver.ImageResponseFromProto(resp), nil
}

func (d *RemoteDriver) ImageToImage(ctx context.Context, req driver.ImageToImageRequest) (*driver.ImageResponse, error) {
	resp, err := d.client.ImageToImage(ctx, driver.ImageToImageRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote image-to-image: %w", err)
	}
	return driver.ImageResponseFromProto(resp), nil
}

func (d *RemoteDriver) InpaintImage(ctx context.Context, req driver.ImageInpaintRequest) (*driver.ImageResponse, error) {
	resp, err := d.client.InpaintImage(ctx, driver.ImageInpaintRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote inpaint image: %w", err)
	}
	return driver.ImageResponseFromProto(resp), nil
}

func (d *RemoteDriver) UpscaleImage(ctx context.Context, req driver.ImageUpscaleRequest) (*driver.ImageResponse, error) {
	resp, err := d.client.UpscaleImage(ctx, driver.ImageUpscaleRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote upscale image: %w", err)
	}
	return driver.ImageResponseFromProto(resp), nil
}

func (d *RemoteDriver) EditImage(ctx context.Context, req driver.ImageEditRequest) (*driver.ImageResponse, error) {
	resp, err := d.client.EditImage(ctx, driver.ImageEditRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote edit image: %w", err)
	}
	return driver.ImageResponseFromProto(resp), nil
}

func (d *RemoteDriver) ControlNetImage(ctx context.Context, req driver.ImageControlNetRequest) (*driver.ImageResponse, error) {
	resp, err := d.client.ControlNetImage(ctx, driver.ImageControlNetRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote controlnet image: %w", err)
	}
	return driver.ImageResponseFromProto(resp), nil
}

func (d *RemoteDriver) DepthToImage(ctx context.Context, req driver.ImageDepthToImageRequest) (*driver.ImageResponse, error) {
	resp, err := d.client.DepthToImage(ctx, driver.ImageDepthToImageRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote depth-to-image: %w", err)
	}
	return driver.ImageResponseFromProto(resp), nil
}

func (d *RemoteDriver) ImageVariation(ctx context.Context, req driver.ImageVariationRequest) (*driver.ImageResponse, error) {
	resp, err := d.client.ImageVariation(ctx, driver.ImageVariationRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote image variation: %w", err)
	}
	return driver.ImageResponseFromProto(resp), nil
}

func (d *RemoteDriver) StyleTransferImage(ctx context.Context, req driver.ImageStyleTransferRequest) (*driver.ImageResponse, error) {
	resp, err := d.client.StyleTransferImage(ctx, driver.ImageStyleTransferRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote style transfer image: %w", err)
	}
	return driver.ImageResponseFromProto(resp), nil
}

func (d *RemoteDriver) SpriteImage(ctx context.Context, req driver.ImageSpriteRequest) (*driver.ImageResponse, error) {
	resp, err := d.client.SpriteImage(ctx, driver.ImageSpriteRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote sprite sheet image: %w", err)
	}
	return driver.ImageResponseFromProto(resp), nil
}
