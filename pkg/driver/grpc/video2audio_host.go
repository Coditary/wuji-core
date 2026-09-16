package grpc

import (
	"context"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (s *Host) Video2Audio(ctx context.Context, req *wujiv1.Video2AudioRequest) (*wujiv1.Video2AudioResponse, error) {
	gen, err := driver.As[driver.VideoToAudioConverter](s.drv, capability.Video2Audio)
	if err != nil {
		return nil, err
	}
	resp, err := gen.VideoToAudio(ctx, driver.Video2AudioRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.Video2AudioResponseToProto(resp), nil
}
