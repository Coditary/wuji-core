package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (s *Host) CreateVoice(ctx context.Context, req *wujiv1.CreateVoiceRequest) (*wujiv1.CloneVoiceResponse, error) {
	gen, ok := s.drv.(driver.VoiceCreator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "create voice not supported")
	}
	resp, err := gen.CreateVoice(ctx, driver.VoiceCreateRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.VoiceResponseToProto(resp), nil
}

func (s *Host) ConvertVoice(ctx context.Context, req *wujiv1.ConvertVoiceRequest) (*wujiv1.CloneVoiceResponse, error) {
	gen, ok := s.drv.(driver.VoiceConverter)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "convert voice not supported")
	}
	resp, err := gen.ConvertVoice(ctx, driver.VoiceConvertRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.VoiceResponseToProto(resp), nil
}

func (s *Host) ListVoiceProfiles(ctx context.Context, _ *wujiv1.ListVoiceProfilesRequest) (*wujiv1.ListVoiceProfilesResponse, error) {
	provider, ok := s.drv.(driver.VoiceCatalogProvider)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "voice catalog not supported")
	}
	profiles, err := provider.ListVoiceProfiles(ctx)
	if err != nil {
		return nil, err
	}
	return driver.ListVoiceProfilesResponseToProto(profiles), nil
}
