package core

import (
	"context"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (c *Core) trackDriverUse(ctx context.Context, driverID string, cap capability.Type) (func(), error) {
	resolvedID := c.ResolveDriverID(driverID, cap)
	c.beginDriverUse(resolvedID)
	if err := c.EnsureDriver(ctx, resolvedID); err != nil {
		c.endDriverUse(resolvedID)
		return nil, err
	}
	return func() { c.endDriverUse(resolvedID) }, nil
}

func (c *Core) TrainText(ctx context.Context, driverID string, req driver.TextTrainRequest) (*driver.TrainResponse, error) {
	end, err := c.trackDriverUse(ctx, driverID, capability.TextGeneration)
	if err != nil {
		return nil, err
	}
	defer end()
	d, err := c.resolveForCapability(driverID, capability.TextGeneration)
	if err != nil {
		return nil, err
	}
	tr, err := driver.As[driver.TextTrainer](d, capability.TextGeneration)
	if err != nil {
		return nil, err
	}
	return tr.TrainText(ctx, req)
}

func (c *Core) TrainImage(ctx context.Context, driverID string, req driver.ImageTrainRequest) (*driver.TrainResponse, error) {
	end, err := c.trackDriverUse(ctx, driverID, capability.ImageGeneration)
	if err != nil {
		return nil, err
	}
	defer end()
	d, err := c.resolveForCapability(driverID, capability.ImageGeneration)
	if err != nil {
		return nil, err
	}
	tr, err := driver.As[driver.ImageTrainer](d, capability.ImageGeneration)
	if err != nil {
		return nil, err
	}
	return tr.TrainImage(ctx, req)
}

func (c *Core) TrainVideo(ctx context.Context, driverID string, req driver.VideoTrainRequest) (*driver.TrainResponse, error) {
	end, err := c.trackDriverUse(ctx, driverID, capability.VideoGeneration)
	if err != nil {
		return nil, err
	}
	defer end()
	d, err := c.resolveForCapability(driverID, capability.VideoGeneration)
	if err != nil {
		return nil, err
	}
	tr, err := driver.As[driver.VideoTrainer](d, capability.VideoGeneration)
	if err != nil {
		return nil, err
	}
	return tr.TrainVideo(ctx, req)
}

func (c *Core) TrainAudio(ctx context.Context, driverID string, req driver.AudioTrainRequest) (*driver.TrainResponse, error) {
	end, err := c.trackDriverUse(ctx, driverID, capability.AudioGeneration)
	if err != nil {
		return nil, err
	}
	defer end()
	d, err := c.resolveForCapability(driverID, capability.AudioGeneration)
	if err != nil {
		return nil, err
	}
	tr, err := driver.As[driver.AudioTrainer](d, capability.AudioGeneration)
	if err != nil {
		return nil, err
	}
	return tr.TrainAudio(ctx, req)
}

func (c *Core) TrainMesh(ctx context.Context, driverID string, req driver.MeshTrainRequest) (*driver.TrainResponse, error) {
	end, err := c.trackDriverUse(ctx, driverID, capability.Mesh)
	if err != nil {
		return nil, err
	}
	defer end()
	d, err := c.resolveForCapability(driverID, capability.Mesh)
	if err != nil {
		return nil, err
	}
	tr, err := driver.As[driver.MeshTrainer](d, capability.Mesh)
	if err != nil {
		return nil, err
	}
	return tr.TrainMesh(ctx, req)
}

func (c *Core) TrainVoice(ctx context.Context, driverID string, req driver.VoiceTrainRequest) (*driver.TrainResponse, error) {
	end, err := c.trackDriverUse(ctx, driverID, capability.VoiceCloning)
	if err != nil {
		return nil, err
	}
	defer end()
	d, err := c.resolveForCapability(driverID, capability.VoiceCloning)
	if err != nil {
		return nil, err
	}
	tr, err := driver.As[driver.VoiceTrainer](d, capability.VoiceCloning)
	if err != nil {
		return nil, err
	}
	return tr.TrainVoice(ctx, req)
}
