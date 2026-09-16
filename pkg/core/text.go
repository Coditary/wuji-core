package core

import (
	"context"
	"os"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/media"
)

func (c *Core) generateTextFromAudio(ctx context.Context, driverID string, req driver.TextRequest) (*driver.TextResponse, error) {
	audioPath := req.MediaPath
	cleanup := func() {}
	if media.IsVideoContainer(audioPath) {
		extracted, err := c.Video2Audio(ctx, driverID, driver.Video2AudioRequest{VideoPath: audioPath})
		if err != nil {
			return nil, err
		}
		audioPath = extracted.AudioPath
		if extracted.Temporary {
			cleanup = func() { _ = os.Remove(audioPath) }
		}
	}
	defer cleanup()

	req.MediaPath = audioPath
	textDriver, err := c.resolveForCapability(driverID, capability.TextGeneration)
	if err != nil {
		return nil, err
	}
	return driver.RunTextInput(ctx, textDriver, req)
}
