package config_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/sysmem"
)

func TestResolvedResourcesDefaults(t *testing.T) {
	cfg := &config.Config{}
	res := cfg.ResolvedResources()
	if res.BudgetVRAM(capability.TextGeneration) != 6144 {
		t.Fatalf("unexpected text budget: %d", res.BudgetVRAM(capability.TextGeneration))
	}
	if !res.LazyRestore {
		t.Fatal("expected lazy restore by default")
	}
	if res.IdleShutdownSeconds != config.DefaultIdleShutdownSeconds {
		t.Fatalf("expected idle shutdown default %d, got %d", config.DefaultIdleShutdownSeconds, res.IdleShutdownSeconds)
	}
}

func TestResolvedResourcesIdleShutdownDisabled(t *testing.T) {
	cfg := &config.Config{
		Drivers: config.DriversConfig{
			Resources: &config.ResourcesConfig{
				IdleShutdownSeconds: intPtr(0),
			},
		},
	}
	res := cfg.ResolvedResources()
	if res.IdleShutdownSeconds != 0 {
		t.Fatalf("expected 0, got %d", res.IdleShutdownSeconds)
	}
}

func TestResolvedResourcesIdleShutdownCustom(t *testing.T) {
	cfg := &config.Config{
		Drivers: config.DriversConfig{
			Resources: &config.ResourcesConfig{
				IdleShutdownSeconds: intPtr(600),
			},
		},
	}
	res := cfg.ResolvedResources()
	if res.IdleShutdownSeconds != 600 {
		t.Fatalf("expected 600, got %d", res.IdleShutdownSeconds)
	}
}

func intPtr(n int) *int { return &n }

func TestResolvedResourcesEnabled(t *testing.T) {
	cfg := &config.Config{
		Drivers: config.DriversConfig{
			Resources: &config.ResourcesConfig{TotalVRAMMB: 12000},
		},
	}
	res := cfg.ResolvedResources()
	if !res.Enabled() || res.TotalVRAMMB != 12000 {
		t.Fatalf("unexpected resources: %+v", res)
	}
}

func TestResolvedResourcesRAMOnly(t *testing.T) {
	cfg := &config.Config{
		Drivers: config.DriversConfig{
			Resources: &config.ResourcesConfig{TotalRAMMB: 49152},
		},
	}
	res := cfg.ResolvedResources()
	if !res.Enabled() || res.TotalRAMMB != 49152 {
		t.Fatalf("unexpected resources: %+v", res)
	}
	if res.RAMBudgetFromAuto {
		t.Fatal("expected manual RAM budget")
	}
	if res.BudgetSource() != config.BudgetSourceConfig {
		t.Fatalf("expected config source, got %s", res.BudgetSource())
	}
	if !res.LazyRestore {
		t.Fatal("expected lazy restore by default")
	}
}

func TestResolvedResourcesPartialBudgetKeepsMaxConcurrent(t *testing.T) {
	cfg := &config.Config{
		Drivers: config.DriversConfig{
			Resources: &config.ResourcesConfig{
				CapabilityBudgets: map[string]config.CapabilityBudget{
					string(capability.TextGeneration): {VRAMMB: 8192},
				},
			},
		},
	}
	res := cfg.ResolvedResources()
	if res.MaxConcurrent(capability.TextGeneration) != 2 {
		t.Fatalf("expected default max concurrent 2, got %d", res.MaxConcurrent(capability.TextGeneration))
	}
	if res.BudgetVRAM(capability.TextGeneration) != 8192 {
		t.Fatalf("expected overridden vram 8192, got %d", res.BudgetVRAM(capability.TextGeneration))
	}
}

func TestResolvedResourcesAutoBudget(t *testing.T) {
	cfg := &config.Config{}
	res := cfg.ResolvedResources()
	if sysmem.TotalRAMMB() == 0 && sysmem.TotalVRAMMB() == 0 {
		t.Skip("no detectable memory on this host")
	}
	if !res.Enabled() {
		t.Fatal("expected auto-enabled resources")
	}
	if res.BudgetSource() != config.BudgetSourceAuto && res.BudgetSource() != config.BudgetSourceMixed {
		t.Fatalf("unexpected budget source: %s", res.BudgetSource())
	}
}
