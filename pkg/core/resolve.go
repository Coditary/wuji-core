package core

import (
	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/driver/dummy"
	ffmpegdrv "github.com/coditary/wuji-core/pkg/driver/ffmpeg"
	localdrv "github.com/coditary/wuji-core/pkg/driver/local"
)

// ResolveDriverID picks the driver for a capability.
// Precedence: explicit preferredID (--driver or caller override) → capability_drivers →
// video2audio default (ffmpeg) → default_driver.
func (c *Core) ResolveDriverID(preferredID string, cap capability.Type) string {
	return ResolveDriverID(c.appConfig, c.defaultDriver, preferredID, cap)
}

// ResolveDriverID picks the driver for a capability from config.
func ResolveDriverID(cfg *config.Config, defaultDriver, preferredID string, cap capability.Type) string {
	if preferredID != "" {
		return localdrv.NormalizeDriverID(preferredID)
	}
	if cfg != nil {
		if d := cfg.DriverForCapability(cap); d != "" {
			return localdrv.NormalizeDriverID(d)
		}
	}
	if cap == capability.Video2Audio {
		return ffmpegdrv.DriverID
	}
	if cap == capability.RAG {
		return "raggo"
	}
	if defaultDriver != "" {
		return localdrv.NormalizeDriverID(defaultDriver)
	}
	return dummy.DriverID
}

func (c *Core) resolveForCapability(preferredID string, cap capability.Type) (driver.Driver, error) {
	return c.resolve(c.ResolveDriverID(preferredID, cap))
}
