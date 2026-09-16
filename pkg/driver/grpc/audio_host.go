package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (s *Host) GenerateMusic(ctx context.Context, req *wujiv1.GenerateMusicRequest) (*wujiv1.GenerateAudioResponse, error) {
	gen, ok := s.drv.(driver.MusicGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "generate music not supported")
	}
	resp, err := gen.GenerateMusic(ctx, driver.AudioMusicRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.AudioResponseToProto(resp), nil
}

func (s *Host) MelodyToMusic(ctx context.Context, req *wujiv1.MelodyToMusicRequest) (*wujiv1.GenerateAudioResponse, error) {
	gen, ok := s.drv.(driver.MelodyToMusicGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "melody-to-music not supported")
	}
	resp, err := gen.MelodyToMusic(ctx, driver.AudioMelodyToMusicRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.AudioResponseToProto(resp), nil
}

func (s *Host) GenerateSFX(ctx context.Context, req *wujiv1.GenerateSFXRequest) (*wujiv1.GenerateAudioResponse, error) {
	gen, ok := s.drv.(driver.SFXGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "generate sfx not supported")
	}
	resp, err := gen.GenerateSFX(ctx, driver.AudioSFXRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.AudioResponseToProto(resp), nil
}

func (s *Host) GenerateAudioSpeech(ctx context.Context, req *wujiv1.GenerateAudioSpeechRequest) (*wujiv1.GenerateAudioResponse, error) {
	gen, ok := s.drv.(driver.AudioSpeechGenerator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "generate audio speech not supported")
	}
	resp, err := gen.GenerateAudioSpeech(ctx, driver.AudioSpeechRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.AudioResponseToProto(resp), nil
}
