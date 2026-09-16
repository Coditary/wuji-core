package grpc

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *RemoteDriver) GenerateVideo(ctx context.Context, req driver.VideoGenerateRequest) (*driver.VideoResponse, error) {
	resp, err := d.client.GenerateVideo(ctx, driver.VideoGenerateRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote generate video: %w", err)
	}
	return driver.VideoResponseFromProto(resp), nil
}

func (d *RemoteDriver) ImageToVideo(ctx context.Context, req driver.VideoImageToVideoRequest) (*driver.VideoResponse, error) {
	resp, err := d.client.ImageToVideo(ctx, driver.VideoImageToVideoRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote image-to-video: %w", err)
	}
	return driver.VideoResponseFromProto(resp), nil
}

func (d *RemoteDriver) InterpolateVideo(ctx context.Context, req driver.VideoInterpolateRequest) (*driver.VideoResponse, error) {
	resp, err := d.client.InterpolateVideo(ctx, driver.VideoInterpolateRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote interpolate video: %w", err)
	}
	return driver.VideoResponseFromProto(resp), nil
}

func (d *RemoteDriver) UpscaleVideo(ctx context.Context, req driver.VideoUpscaleRequest) (*driver.VideoResponse, error) {
	resp, err := d.client.UpscaleVideo(ctx, driver.VideoUpscaleRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote upscale video: %w", err)
	}
	return driver.VideoResponseFromProto(resp), nil
}
