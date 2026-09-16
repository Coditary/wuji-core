package core

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/scheduler"
	"github.com/coditary/wuji-core/pkg/sysmem"
)

var (
	warnResourcesMu         sync.Mutex
	warnedResourcesDisabled bool
)

func (c *Core) initScheduler() {
	res := c.appConfig.ResolvedResources()
	c.resources = scheduler.New(c.appConfig.Root, res, c.unloadSlot, c.warmSlot, c.probeLoaded)
}

// ResourceStatus returns the persisted scheduler state.
func (c *Core) ResourceStatus() (scheduler.PersistedState, error) {
	if c.resources == nil {
		return scheduler.PersistedState{}, nil
	}
	return c.resources.Status()
}

// ResourcesConfig returns effective resource settings.
func (c *Core) ResourcesConfig() config.EffectiveResources {
	return c.appConfig.ResolvedResources()
}

// PurgeInference stops llama and vllm inference servers on configured ports.
func (c *Core) PurgeInference(ctx context.Context) error {
	return c.UnloadAllInference(ctx)
}

// LoadInference preloads a model for the given driver.
func (c *Core) LoadInference(ctx context.Context, driverID, model string, cap capability.Type) error {
	driverID = c.ResolveDriverID(driverID, cap)
	c.beginDriverUse(driverID)
	defer c.endDriverUse(driverID)
	if err := c.EnsureDriver(ctx, driverID); err != nil {
		return err
	}
	if model == "" {
		model = c.defaultModelFor(driverID)
	}
	job := jobSpec(driverID, model, cap)
	job.Transient = false

	if c.resources != nil {
		return c.resources.LoadDriver(ctx, job)
	}
	return c.warmSlot(ctx, scheduler.SessionSnapshot{
		DriverID: driverID, Model: model, Capability: cap,
	})
}

// UnloadInferenceDriver unloads one driver or all when driverID is "all".
func (c *Core) UnloadInferenceDriver(ctx context.Context, driverID string) error {
	if driverID == "all" {
		return c.UnloadAllInference(ctx)
	}
	if c.resources != nil {
		return c.resources.UnloadDriver(ctx, driverID)
	}
	return c.unloadSlot(ctx, driverID)
}

// UnloadAllInference stops all known inference backends and clears scheduler state.
func (c *Core) UnloadAllInference(ctx context.Context) error {
	if c.resources != nil {
		var firstErr error
		for _, id := range []string{"llama", "vllm", "a1111"} {
			if err := c.resources.UnloadDriver(ctx, id); err != nil && firstErr == nil {
				if isDriverNotFound(err) {
					continue
				}
				firstErr = err
			}
		}
		return firstErr
	}
	for _, id := range []string{"llama", "vllm"} {
		if err := c.unloadSlot(ctx, id); err != nil && !isDriverNotFound(err) {
			return err
		}
	}
	return nil
}

func (c *Core) defaultModelFor(driverID string) string {
	switch driverID {
	case "llama":
		return c.appConfig.ResolvedLlama().DefaultModel
	case "vllm":
		return c.appConfig.ResolvedVLLM().DefaultModel
	case "a1111":
		return c.appConfig.ResolvedA1111().DefaultModel
	default:
		return ""
	}
}

func isDriverNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}

func (c *Core) maybeWarnResourcesDisabled(driverID string) {
	if c.resources != nil && c.resources.Enabled() {
		return
	}
	if driverID == "dummy" {
		return
	}
	warnResourcesMu.Lock()
	defer warnResourcesMu.Unlock()
	if warnedResourcesDisabled {
		return
	}
	warnedResourcesDisabled = true

	ram := sysmem.TotalRAMMB()
	msg := "warning: resource management is disabled — inference can exhaust all RAM and crash the system."
	if ram > 0 {
		suggested := ram * 80 / 100
		msg += fmt.Sprintf("\n  Add to .wuji/config.yaml: drivers.resources.total_ram_mb: %d  (detected %d MB system RAM)", suggested, ram)
	} else {
		msg += "\n  Set drivers.resources.total_ram_mb in .wuji/config.yaml to enable budgeting."
	}
	fmt.Fprintln(os.Stderr, msg)
}

func (c *Core) scheduleRun(ctx context.Context, job scheduler.JobSpec, fn func(context.Context) error) error {
	c.maybeWarnResourcesDisabled(job.DriverID)

	// RAG embeds via Ollama HTTP inside the driver; no local inference RAM slot.
	if job.Capability == capability.RAG {
		return fn(ctx)
	}

	if c.resources == nil || !c.resources.Enabled() {
		return fn(ctx)
	}
	res := c.appConfig.ResolvedResources()
	if job.VRAMMB <= 0 {
		job.VRAMMB = res.EstimateModelVRAM(job.Capability, job.Model)
	}
	if job.RAMMB <= 0 {
		job.RAMMB = res.EstimateModelRAM(job.Capability, job.Model)
	}
	return c.resources.Run(ctx, job, fn)
}

func (c *Core) unloadSlot(ctx context.Context, driverID string) error {
	d, err := c.registry.Get(driverID)
	if err != nil {
		// Stale scheduler entry (e.g. provider id "gemma4" instead of driver "llama").
		return nil
	}
	if err := driver.UnloadInference(ctx, d); err != nil {
		return err
	}
	switch driverID {
	case "llama":
		cfg := c.appConfig.ResolvedLlama()
		return scheduler.StopListenerOnPort(cfg.InferencePort)
	case "vllm":
		cfg := c.appConfig.ResolvedVLLM()
		return scheduler.StopListenerOnPort(cfg.Port)
	case "a1111":
		cfg := c.appConfig.ResolvedA1111()
		_, port, err := config.ParseHostPort(cfg.ResolvedAPIBase())
		if err != nil {
			return err
		}
		return scheduler.StopListenerOnPort(port)
	}
	return nil
}

func (c *Core) probeLoaded(driverID string) bool {
	switch driverID {
	case "llama":
		if _, err := c.registry.Get("llama"); err != nil {
			return false
		}
		return scheduler.PortListening(c.appConfig.ResolvedLlama().InferencePort)
	case "vllm":
		if _, err := c.registry.Get("vllm"); err != nil {
			return false
		}
		return scheduler.PortListening(c.appConfig.ResolvedVLLM().Port)
	case "a1111":
		if _, err := c.registry.Get("a1111"); err != nil {
			return false
		}
		_, port, err := config.ParseHostPort(c.appConfig.ResolvedA1111().ResolvedAPIBase())
		if err != nil {
			return false
		}
		return scheduler.PortListening(port)
	case "raggo", "gorag", "local-rag", "local", "rag":
		return false
	default:
		return false
	}
}

func (c *Core) warmSlot(ctx context.Context, snap scheduler.SessionSnapshot) error {
	c.beginDriverUse(snap.DriverID)
	defer c.endDriverUse(snap.DriverID)
	d, err := c.registry.Get(snap.DriverID)
	if err != nil {
		return err
	}
	if w, ok := d.(driver.ModelWarmer); ok {
		return w.WarmModel(ctx, snap.Model)
	}
	return nil
}

func capabilityTransient(cap capability.Type) bool {
	switch cap {
	case capability.ImageGeneration, capability.VideoGeneration,
		capability.AudioGeneration, capability.Mesh:
		return true
	default:
		return false
	}
}

func jobSpec(driverID, model string, cap capability.Type) scheduler.JobSpec {
	return scheduler.JobSpec{
		DriverID:   driverID,
		Model:      model,
		Capability: cap,
		Transient:  capabilityTransient(cap),
	}
}

func scheduleGenerate[T any](
	c *Core,
	ctx context.Context,
	driverID string,
	model string,
	cap capability.Type,
	run func(context.Context) (T, error),
) (T, error) {
	var zero T
	resolvedID := c.ResolveDriverID(driverID, cap)
	c.beginDriverUse(resolvedID)
	defer c.endDriverUse(resolvedID)
	if err := c.ensureAPIInferenceBackend(ctx, resolvedID, cap, model); err != nil {
		return zero, err
	}
	if err := c.EnsureDriver(ctx, resolvedID); err != nil {
		return zero, err
	}
	var result T
	err := c.scheduleRun(ctx, jobSpec(resolvedID, model, cap), func(ctx context.Context) error {
		out, err := run(ctx)
		if err != nil {
			return err
		}
		result = out
		return nil
	})
	if err != nil {
		return zero, err
	}
	return result, nil
}
