package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const reqpackDriversFileName = "drivers-reqpack.yaml"

// ReqpackDriverBinaryPath returns the canonical binary installed by rqp-plugin-wuji.
func ReqpackDriverBinaryPath(driverID string) string {
	home, err := os.UserHomeDir()
	if err != nil || driverID == "" {
		return ""
	}
	path := filepath.Join(home, ".local", "share", "wuji", "bin", "wuji-driver-"+driverID)
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return path
}

func mergeReqpackDriverOverlay(cfg *Config) error {
	if cfg == nil {
		return nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	path := filepath.Join(home, dirName, reqpackDriversFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read %s: %w", path, err)
	}

	var overlay fileConfig
	if err := yaml.Unmarshal(data, &overlay); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	if len(overlay.CapabilityDrivers) > 0 {
		if cfg.CapabilityDrivers == nil {
			cfg.CapabilityDrivers = CapabilityDrivers{}
		}
		for cap, driverID := range overlay.CapabilityDrivers {
			if driverID != "" {
				cfg.CapabilityDrivers[cap] = driverID
			}
		}
	}
	if overlay.DefaultDriver != "" {
		cfg.DefaultDriver = overlay.DefaultDriver
	}
	mergeDriversOverlay(cfg, overlay.Drivers)
	return nil
}

func mergeDriversOverlay(cfg *Config, overlay DriversConfig) {
	if overlay.VLLM != nil {
		if cfg.Drivers.VLLM == nil {
			cfg.Drivers.VLLM = &VLLMConfig{}
		}
		mergeVLLM(cfg.Drivers.VLLM, *overlay.VLLM)
	}
	if overlay.Llama != nil {
		if cfg.Drivers.Llama == nil {
			cfg.Drivers.Llama = &LlamaConfig{}
		}
		mergeLlama(cfg.Drivers.Llama, *overlay.Llama)
	}
	if overlay.A1111 != nil {
		if cfg.Drivers.A1111 == nil {
			cfg.Drivers.A1111 = &A1111Config{}
		}
		mergeA1111(cfg.Drivers.A1111, *overlay.A1111)
	}
	if overlay.Raggo != nil {
		if cfg.Drivers.Raggo == nil {
			cfg.Drivers.Raggo = &RaggoConfig{}
		}
		mergeRaggo(cfg.Drivers.Raggo, *overlay.Raggo)
	}
	if overlay.Gorag != nil {
		if cfg.Drivers.Gorag == nil {
			cfg.Drivers.Gorag = &GoragConfig{}
		}
		mergeGorag(cfg.Drivers.Gorag, *overlay.Gorag)
	}
	cfg.Drivers.Registry = mergeDriverRegistry(cfg.Drivers.Registry, overlay.Registry)
}
