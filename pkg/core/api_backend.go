package core

import (
	"context"
	"strings"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/scheduler"
)

func (c *Core) ensureAPIInferenceBackend(ctx context.Context, driverID string, cap capability.Type, model string) error {
	if c.appConfig == nil {
		return nil
	}
	entry, ok := c.appConfig.APIs[driverID]
	if !ok {
		return nil
	}
	ensureID := strings.TrimSpace(apiEnsureDriver(entry, cap))
	if ensureID == "" {
		return nil
	}
	if err := c.EnsureDriver(ctx, ensureID); err != nil {
		return err
	}
	warmModel := strings.TrimSpace(model)
	if warmModel == "" {
		warmModel = apiDefaultModel(entry, cap)
	}
	if warmModel == "" {
		warmModel = c.defaultModelFor(ensureID)
	}
	return c.warmSlot(ctx, scheduler.SessionSnapshot{
		DriverID:   ensureID,
		Model:      warmModel,
		Capability: cap,
	})
}

func apiEnsureDriver(entry config.APIEntry, cap capability.Type) string {
	spec := apiSpecForCapability(entry, cap)
	if spec == nil {
		return ""
	}
	return spec.EnsureDriver
}

func apiDefaultModel(entry config.APIEntry, cap capability.Type) string {
	spec := apiSpecForCapability(entry, cap)
	if spec == nil {
		return ""
	}
	if spec.Model != "" {
		return spec.Model
	}
	if spec.Defaults != nil && spec.Defaults.Model != "" {
		return spec.Defaults.Model
	}
	return ""
}

func apiSpecForCapability(entry config.APIEntry, cap capability.Type) *config.APICapabilitySpec {
	switch cap {
	case capability.TextGeneration:
		return entry.Text
	default:
		return nil
	}
}
