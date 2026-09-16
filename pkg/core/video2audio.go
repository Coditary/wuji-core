package core

import (
	"context"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (c *Core) Video2Audio(ctx context.Context, driverID string, req driver.Video2AudioRequest) (*driver.Video2AudioResponse, error) {
	return scheduleGenerate(c, ctx, driverID, "", capability.Video2Audio, func(ctx context.Context) (*driver.Video2AudioResponse, error) {
		d, err := c.resolveForCapability(driverID, capability.Video2Audio)
		if err != nil {
			return nil, err
		}
		gen, err := driver.As[driver.VideoToAudioConverter](d, capability.Video2Audio)
		if err != nil {
			return nil, err
		}
		return gen.VideoToAudio(ctx, req)
	})
}
