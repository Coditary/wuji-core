package config

import (
	"fmt"
	"strings"

	"github.com/coditary/wuji-core/pkg/sysmem"
)

// InitResources writes detected RAM/VRAM budgets to config.yaml.
func InitResources(root string) (EffectiveResources, error) {
	cfg, err := Load()
	if err != nil {
		return EffectiveResources{}, err
	}
	if cfg.Root == "" {
		cfg.Root = root
	}
	if cfg.Drivers.Resources == nil {
		cfg.Drivers.Resources = &ResourcesConfig{}
	}

	ram, vram := sysmem.SuggestedBudget()
	if ram > 0 {
		cfg.Drivers.Resources.TotalRAMMB = ram
	}
	if vram > 0 {
		cfg.Drivers.Resources.TotalVRAMMB = vram
	}
	if ram == 0 && vram == 0 {
		return EffectiveResources{}, fmt.Errorf("could not detect system RAM or GPU VRAM")
	}
	if err := cfg.Save(); err != nil {
		return EffectiveResources{}, err
	}
	return cfg.ResolvedResources(), nil
}

// SetResourcesField updates a single resources config key and persists config.yaml.
func SetResourcesField(root, key, value string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	if cfg.Root == "" {
		cfg.Root = root
	}
	if cfg.Drivers.Resources == nil {
		cfg.Drivers.Resources = &ResourcesConfig{}
	}
	res := cfg.Drivers.Resources

	switch key {
	case "total_ram_mb", "ram_budget":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil || n < 0 {
			return fmt.Errorf("total_ram_mb must be a non-negative integer")
		}
		res.TotalRAMMB = n
	case "total_vram_mb", "vram_budget":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil || n < 0 {
			return fmt.Errorf("total_vram_mb must be a non-negative integer")
		}
		res.TotalVRAMMB = n
	case "lazy_restore", "evict_and_restore", "unload_transient_models":
		v, err := parseBool(value)
		if err != nil {
			return err
		}
		switch key {
		case "lazy_restore":
			res.LazyRestore = &v
		case "evict_and_restore":
			res.EvictAndRestore = &v
		case "unload_transient_models":
			res.UnloadTransientModels = &v
		}
	case "idle_shutdown_seconds":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil || n < 0 {
			return fmt.Errorf("idle_shutdown_seconds must be a non-negative integer")
		}
		res.IdleShutdownSeconds = &n
	default:
		return fmt.Errorf("unknown resources config key %q", key)
	}

	return cfg.Save()
}

func parseBool(value string) (bool, error) {
	switch strings.ToLower(value) {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return false, fmt.Errorf("expected true or false, got %q", value)
	}
}
