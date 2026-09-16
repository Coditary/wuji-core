package core

import (
	"context"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
)

// ListVoiceProfiles returns speakable voice metadata from the configured voice driver.
func (c *Core) ListVoiceProfiles(ctx context.Context, driverID string) ([]driver.VoiceProfileInfo, error) {
	d, err := c.resolveForCapability(driverID, capability.VoiceCloning)
	if err != nil {
		return nil, err
	}
	return driver.ListVoiceProfiles(ctx, d)
}

// FindVoiceProfile looks up one voice profile by name or id.
func (c *Core) FindVoiceProfile(ctx context.Context, driverID, nameOrID string) (*driver.VoiceProfileInfo, error) {
	d, err := c.resolveForCapability(driverID, capability.VoiceCloning)
	if err != nil {
		return nil, err
	}
	return driver.FindVoiceProfile(ctx, d, nameOrID)
}
