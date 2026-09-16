package config

import (
	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/sysmem"
)

// DefaultIdleShutdownSeconds is how long an unused remote driver may stay loaded in daemon mode.
const DefaultIdleShutdownSeconds = 1800

// ResourcesConfig is the on-disk drivers.resources block.
type ResourcesConfig struct {
	TotalVRAMMB           int                         `yaml:"total_vram_mb,omitempty"`
	TotalRAMMB            int                         `yaml:"total_ram_mb,omitempty"`
	EvictAndRestore       *bool                       `yaml:"evict_and_restore,omitempty"`
	UnloadTransientModels *bool                       `yaml:"unload_transient_models,omitempty"`
	CapabilityBudgets     map[string]CapabilityBudget `yaml:"capability_budgets,omitempty"`
	ModelVRAMMB           map[string]int              `yaml:"model_vram_mb,omitempty"`
	ModelRAMMB            map[string]int              `yaml:"model_ram_mb,omitempty"`
	LazyRestore           *bool                       `yaml:"lazy_restore,omitempty"`
	IdleShutdownSeconds   *int                        `yaml:"idle_shutdown_seconds,omitempty"`
}

// EffectiveResources is the resolved runtime resource policy.
type EffectiveResources struct {
	TotalVRAMMB           int
	TotalRAMMB            int
	EvictAndRestore       bool
	UnloadTransientModels bool
	LazyRestore           bool
	CapabilityBudgets     map[string]CapabilityBudget
	ModelVRAMMB           map[string]int
	ModelRAMMB            map[string]int
	// RAMBudgetFromAuto / VRAMBudgetFromAuto are true when the budget was detected, not read from config.
	RAMBudgetFromAuto  bool
	VRAMBudgetFromAuto bool
	// IdleShutdownSeconds: 0 disables idle shutdown; daemon default when unset is DefaultIdleShutdownSeconds.
	IdleShutdownSeconds int
}

// CapabilityBudget limits memory and concurrency for one capability.
type CapabilityBudget struct {
	VRAMMB        int `yaml:"vram_mb,omitempty"`
	MaxConcurrent int `yaml:"max_concurrent,omitempty"`
}

// ResolvedResources returns resource settings with defaults applied.
func (c *Config) ResolvedResources() EffectiveResources {
	out := EffectiveResources{
		EvictAndRestore:       true,
		UnloadTransientModels: true,
		LazyRestore:           true,
		CapabilityBudgets:     defaultCapabilityBudgets(),
		IdleShutdownSeconds:   DefaultIdleShutdownSeconds,
	}
	if c == nil || c.Drivers.Resources == nil {
		applyAutoBudget(&out)
		return out
	}
	src := c.Drivers.Resources
	if src.TotalVRAMMB > 0 {
		out.TotalVRAMMB = src.TotalVRAMMB
		out.VRAMBudgetFromAuto = false
	}
	if src.TotalRAMMB > 0 {
		out.TotalRAMMB = src.TotalRAMMB
		out.RAMBudgetFromAuto = false
	}
	if src.EvictAndRestore != nil {
		out.EvictAndRestore = *src.EvictAndRestore
	}
	if src.UnloadTransientModels != nil {
		out.UnloadTransientModels = *src.UnloadTransientModels
	}
	for k, v := range src.CapabilityBudgets {
		base := out.CapabilityBudgets[k]
		if v.VRAMMB > 0 {
			base.VRAMMB = v.VRAMMB
		}
		if v.MaxConcurrent > 0 {
			base.MaxConcurrent = v.MaxConcurrent
		}
		out.CapabilityBudgets[k] = base
	}
	if len(src.ModelVRAMMB) > 0 {
		out.ModelVRAMMB = src.ModelVRAMMB
	}
	if len(src.ModelRAMMB) > 0 {
		out.ModelRAMMB = src.ModelRAMMB
	}
	if src.LazyRestore != nil {
		out.LazyRestore = *src.LazyRestore
	}
	if src.IdleShutdownSeconds != nil {
		out.IdleShutdownSeconds = *src.IdleShutdownSeconds
	}
	if src.TotalVRAMMB == 0 && src.TotalRAMMB == 0 {
		applyAutoBudget(&out)
	}
	return out
}

func applyAutoBudget(out *EffectiveResources) {
	ram, vram := sysmem.SuggestedBudget()
	if ram > 0 && out.TotalRAMMB == 0 {
		out.TotalRAMMB = ram
		out.RAMBudgetFromAuto = true
	}
	if vram > 0 && out.TotalVRAMMB == 0 {
		out.TotalVRAMMB = vram
		out.VRAMBudgetFromAuto = true
	}
}

// BudgetSource describes where memory budgets come from.
type BudgetSource string

const (
	BudgetSourceDisabled BudgetSource = "disabled"
	BudgetSourceAuto     BudgetSource = "auto-detected"
	BudgetSourceConfig   BudgetSource = "config"
	BudgetSourceMixed    BudgetSource = "mixed"
)

// BudgetSource reports whether budgets are auto-detected, from config, or disabled.
func (r EffectiveResources) BudgetSource() BudgetSource {
	if !r.Enabled() {
		return BudgetSourceDisabled
	}
	auto := r.RAMBudgetFromAuto || r.VRAMBudgetFromAuto
	manual := (r.TotalRAMMB > 0 && !r.RAMBudgetFromAuto) || (r.TotalVRAMMB > 0 && !r.VRAMBudgetFromAuto)
	switch {
	case auto && manual:
		return BudgetSourceMixed
	case auto:
		return BudgetSourceAuto
	default:
		return BudgetSourceConfig
	}
}

func defaultCapabilityBudgets() map[string]CapabilityBudget {
	return map[string]CapabilityBudget{
		string(capability.TextGeneration):  {VRAMMB: 6144, MaxConcurrent: 2},
		string(capability.ImageGeneration): {VRAMMB: 8192, MaxConcurrent: 1},
		string(capability.VideoGeneration): {VRAMMB: 10240, MaxConcurrent: 1},
		string(capability.AudioGeneration): {VRAMMB: 4096, MaxConcurrent: 1},
	}
}

// Enabled reports whether the scheduler should manage resources.
func (r EffectiveResources) Enabled() bool {
	return r.TotalVRAMMB > 0 || r.TotalRAMMB > 0
}

// BudgetVRAM returns the configured VRAM budget for a capability.
func (r EffectiveResources) BudgetVRAM(cap capability.Type) int {
	if b, ok := r.CapabilityBudgets[string(cap)]; ok && b.VRAMMB > 0 {
		return b.VRAMMB
	}
	return 4096
}

// MaxConcurrent returns how many jobs of this capability may run in parallel.
func (r EffectiveResources) MaxConcurrent(cap capability.Type) int {
	key := string(cap)
	if b, ok := r.CapabilityBudgets[key]; ok && b.MaxConcurrent > 0 {
		return b.MaxConcurrent
	}
	if b, ok := defaultCapabilityBudgets()[key]; ok && b.MaxConcurrent > 0 {
		return b.MaxConcurrent
	}
	if cap == capability.TextGeneration {
		return 2
	}
	return 1
}

// EstimateModelVRAM returns configured or default VRAM for a model.
func (r EffectiveResources) EstimateModelVRAM(cap capability.Type, model string) int {
	if model != "" {
		if mb, ok := r.ModelVRAMMB[model]; ok && mb > 0 {
			return mb
		}
	}
	return r.BudgetVRAM(cap)
}

// EstimateModelRAM returns configured or estimated system RAM for a model.
func (r EffectiveResources) EstimateModelRAM(cap capability.Type, model string) int {
	if model != "" {
		if mb, ok := r.ModelRAMMB[model]; ok && mb > 0 {
			return mb
		}
	}
	base := r.EstimateModelVRAM(cap, model)
	switch cap {
	case capability.TextGeneration, capability.VideoGeneration:
		return base * 2
	default:
		return base
	}
}
