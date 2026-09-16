package grpc

import (
	"context"
	"fmt"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *RemoteDriver) CreateVoice(ctx context.Context, req driver.VoiceCreateRequest) (*driver.VoiceResponse, error) {
	resp, err := d.client.CreateVoice(ctx, driver.VoiceCreateRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote create voice: %w", err)
	}
	return driver.VoiceResponseFromProto(resp), nil
}

func (d *RemoteDriver) ConvertVoice(ctx context.Context, req driver.VoiceConvertRequest) (*driver.VoiceResponse, error) {
	resp, err := d.client.ConvertVoice(ctx, driver.VoiceConvertRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote convert voice: %w", err)
	}
	return driver.VoiceResponseFromProto(resp), nil
}

func (d *RemoteDriver) ListVoiceProfiles(ctx context.Context) ([]driver.VoiceProfileInfo, error) {
	resp, err := d.client.ListVoiceProfiles(ctx, &wujiv1.ListVoiceProfilesRequest{})
	if err != nil {
		return nil, fmt.Errorf("remote list voice profiles: %w", err)
	}
	return driver.ListVoiceProfilesResponseFromProto(resp), nil
}
