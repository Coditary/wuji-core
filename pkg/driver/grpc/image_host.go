package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (s *Host) GenerateImage(ctx context.Context, req *wujiv1.GenerateImageRequest) (*wujiv1.GenerateImageResponse, error) {
	gen, ok := s.drv.(driver.ImageGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "generate image not supported")
	}
	resp, err := gen.GenerateImage(ctx, driver.ImageGenerateRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.ImageResponseToProto(resp), nil
}

func (s *Host) ImageToImage(ctx context.Context, req *wujiv1.ImageToImageRequest) (*wujiv1.GenerateImageResponse, error) {
	gen, ok := s.drv.(driver.ImageToImageGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "image-to-image not supported")
	}
	resp, err := gen.ImageToImage(ctx, driver.ImageToImageRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.ImageResponseToProto(resp), nil
}

func (s *Host) InpaintImage(ctx context.Context, req *wujiv1.InpaintImageRequest) (*wujiv1.GenerateImageResponse, error) {
	gen, ok := s.drv.(driver.ImageInpaintGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "inpaint not supported")
	}
	resp, err := gen.InpaintImage(ctx, driver.ImageInpaintRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.ImageResponseToProto(resp), nil
}

func (s *Host) UpscaleImage(ctx context.Context, req *wujiv1.UpscaleImageRequest) (*wujiv1.GenerateImageResponse, error) {
	gen, ok := s.drv.(driver.ImageUpscaleGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "upscale not supported")
	}
	resp, err := gen.UpscaleImage(ctx, driver.ImageUpscaleRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.ImageResponseToProto(resp), nil
}

func (s *Host) EditImage(ctx context.Context, req *wujiv1.EditImageRequest) (*wujiv1.GenerateImageResponse, error) {
	gen, ok := s.drv.(driver.ImageEditGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "edit image not supported")
	}
	resp, err := gen.EditImage(ctx, driver.ImageEditRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.ImageResponseToProto(resp), nil
}

func (s *Host) ControlNetImage(ctx context.Context, req *wujiv1.ControlNetImageRequest) (*wujiv1.GenerateImageResponse, error) {
	gen, ok := s.drv.(driver.ImageControlNetGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "controlnet not supported")
	}
	resp, err := gen.ControlNetImage(ctx, driver.ImageControlNetRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.ImageResponseToProto(resp), nil
}

func (s *Host) DepthToImage(ctx context.Context, req *wujiv1.DepthToImageRequest) (*wujiv1.GenerateImageResponse, error) {
	gen, ok := s.drv.(driver.ImageDepthToImageGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "depth-to-image not supported")
	}
	resp, err := gen.DepthToImage(ctx, driver.ImageDepthToImageRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.ImageResponseToProto(resp), nil
}

func (s *Host) ImageVariation(ctx context.Context, req *wujiv1.ImageVariationRequest) (*wujiv1.GenerateImageResponse, error) {
	gen, ok := s.drv.(driver.ImageVariationGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "image variation not supported")
	}
	resp, err := gen.ImageVariation(ctx, driver.ImageVariationRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.ImageResponseToProto(resp), nil
}

func (s *Host) StyleTransferImage(ctx context.Context, req *wujiv1.StyleTransferImageRequest) (*wujiv1.GenerateImageResponse, error) {
	gen, ok := s.drv.(driver.ImageStyleTransferGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "style transfer not supported")
	}
	resp, err := gen.StyleTransferImage(ctx, driver.ImageStyleTransferRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.ImageResponseToProto(resp), nil
}

func (s *Host) SpriteImage(ctx context.Context, req *wujiv1.SpriteImageRequest) (*wujiv1.GenerateImageResponse, error) {
	gen, ok := s.drv.(driver.ImageSpriteGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "sprite not supported")
	}
	resp, err := gen.SpriteImage(ctx, driver.ImageSpriteRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.ImageResponseToProto(resp), nil
}
