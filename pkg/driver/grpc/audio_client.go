package grpc

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *RemoteDriver) GenerateMusic(ctx context.Context, req driver.AudioMusicRequest) (*driver.AudioResponse, error) {
	resp, err := d.client.GenerateMusic(ctx, driver.AudioMusicRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote generate music: %w", err)
	}
	return driver.AudioResponseFromProto(resp), nil
}

func (d *RemoteDriver) MelodyToMusic(ctx context.Context, req driver.AudioMelodyToMusicRequest) (*driver.AudioResponse, error) {
	resp, err := d.client.MelodyToMusic(ctx, driver.AudioMelodyToMusicRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote melody-to-music: %w", err)
	}
	return driver.AudioResponseFromProto(resp), nil
}

func (d *RemoteDriver) GenerateSFX(ctx context.Context, req driver.AudioSFXRequest) (*driver.AudioResponse, error) {
	resp, err := d.client.GenerateSFX(ctx, driver.AudioSFXRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote generate sfx: %w", err)
	}
	return driver.AudioResponseFromProto(resp), nil
}

func (d *RemoteDriver) GenerateAudioSpeech(ctx context.Context, req driver.AudioSpeechRequest) (*driver.AudioResponse, error) {
	resp, err := d.client.GenerateAudioSpeech(ctx, driver.AudioSpeechRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote generate audio speech: %w", err)
	}
	return driver.AudioResponseFromProto(resp), nil
}
