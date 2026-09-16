package grpc

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *RemoteDriver) VideoToAudio(ctx context.Context, req driver.Video2AudioRequest) (*driver.Video2AudioResponse, error) {
	resp, err := d.client.Video2Audio(ctx, driver.Video2AudioRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote video2audio: %w", err)
	}
	return driver.Video2AudioResponseFromProto(resp), nil
}
