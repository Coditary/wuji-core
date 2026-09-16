package config

import (
	"fmt"
	"strings"

	"github.com/coditary/wuji-core/pkg/capability"
)

// CapabilityDrivers maps capability names (text, image, …) to driver IDs.
type CapabilityDrivers map[string]string

var configurableCapabilities = []capability.Type{
	capability.TextGeneration,
	capability.ImageGeneration,
	capability.ImageUpscale,
	capability.ImageDownscale,
	capability.ImageScale,
	capability.VideoGeneration,
	capability.AudioGeneration,
	capability.Mesh,
	capability.Video2Audio,
	capability.VoiceCloning,
	capability.DatasetMgmt,
	capability.Data,
	capability.RAG,
}

// DriverForCapability returns the configured driver for a capability, or "".
func (c *Config) DriverForCapability(cap capability.Type) string {
	if c == nil || c.CapabilityDrivers == nil {
		return ""
	}
	if id := c.CapabilityDrivers[string(cap)]; id != "" {
		return id
	}
	if cap == capability.Mesh {
		return c.CapabilityDrivers["3d"]
	}
	return ""
}

// SetCapabilityDriver sets the driver for a capability and persists config.yaml.
func SetCapabilityDriver(root, capName, driverID string) error {
	cap, err := parseCapabilityName(capName)
	if err != nil {
		return err
	}
	if strings.TrimSpace(driverID) == "" {
		return fmt.Errorf("driver id must not be empty")
	}

	cfg, err := Load()
	if err != nil {
		return err
	}
	if cfg.Root == "" {
		cfg.Root = root
	}
	if cfg.CapabilityDrivers == nil {
		cfg.CapabilityDrivers = CapabilityDrivers{}
	}
	cfg.CapabilityDrivers[string(cap)] = driverID
	return cfg.Save()
}

func parseCapabilityName(name string) (capability.Type, error) {
	cap, err := capability.Parse(name)
	if err != nil {
		return "", fmt.Errorf("unknown capability %q (valid: %s)", name, capabilityDriverKeys())
	}
	for _, c := range configurableCapabilities {
		if c == cap {
			return cap, nil
		}
	}
	return "", fmt.Errorf("unknown capability %q (valid: %s)", name, capabilityDriverKeys())
}

func capabilityDriverKeys() string {
	parts := make([]string, 0, len(configurableCapabilities))
	for _, c := range configurableCapabilities {
		parts = append(parts, string(c))
	}
	return strings.Join(parts, ", ")
}
