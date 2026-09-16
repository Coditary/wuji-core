package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (s *Host) GenerateVideo(ctx context.Context, req *wujiv1.GenerateVideoRequest) (*wujiv1.GenerateVideoResponse, error) {
	gen, ok := s.drv.(driver.VideoGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "generate video not supported")
	}
	resp, err := gen.GenerateVideo(ctx, driver.VideoGenerateRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.VideoResponseToProto(resp), nil
}

func (s *Host) ImageToVideo(ctx context.Context, req *wujiv1.ImageToVideoRequest) (*wujiv1.GenerateVideoResponse, error) {
	gen, ok := s.drv.(driver.VideoImageToVideoGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "image-to-video not supported")
	}
	resp, err := gen.ImageToVideo(ctx, driver.VideoImageToVideoRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.VideoResponseToProto(resp), nil
}

func (s *Host) InterpolateVideo(ctx context.Context, req *wujiv1.InterpolateVideoRequest) (*wujiv1.GenerateVideoResponse, error) {
	gen, ok := s.drv.(driver.VideoInterpolateGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "interpolate video not supported")
	}
	resp, err := gen.InterpolateVideo(ctx, driver.VideoInterpolateRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.VideoResponseToProto(resp), nil
}

func (s *Host) UpscaleVideo(ctx context.Context, req *wujiv1.UpscaleVideoRequest) (*wujiv1.GenerateVideoResponse, error) {
	gen, ok := s.drv.(driver.VideoUpscaleGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "upscale video not supported")
	}
	resp, err := gen.UpscaleVideo(ctx, driver.VideoUpscaleRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.VideoResponseToProto(resp), nil
}
