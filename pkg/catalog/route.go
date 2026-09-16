package catalog

import (
	"sort"
	"strings"

	"github.com/coditary/wuji-core/pkg/config"
)

// APIEntryFromProvider builds a Wuji apis: entry from a catalog provider.
func APIEntryFromProvider(p ProviderEntry) config.APIEntry {
	protocol := NPMProtocol(p.NPM)
	spec := &config.APICapabilitySpec{
		Protocol:  protocol,
		BaseURL:   strings.TrimSpace(p.API),
		Streaming: true,
	}
	if len(p.Env) > 0 {
		spec.APIKeyEnv = p.Env[0]
	}
	if model := defaultProviderModelID(p); model != "" {
		spec.Model = model
		spec.Defaults = &config.APIDefaults{Model: model}
	}
	if protocol == "anthropic" {
		spec.AnthropicVersion = "2023-06-01"
	}
	return config.APIEntry{Text: spec}
}

func defaultProviderModelID(p ProviderEntry) string {
	if len(p.Models) == 0 {
		return ""
	}
	ids := make([]string, 0, len(p.Models))
	for id := range p.Models {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids[0]
}

// EnsureAPIRoute merges a catalog provider into cfg.APIs.
func EnsureAPIRoute(cfg *config.Config, providerID string, p ProviderEntry) {
	if cfg == nil {
		return
	}
	if cfg.APIs == nil {
		cfg.APIs = map[string]config.APIEntry{}
	}
	entry := APIEntryFromProvider(p)
	if existing, ok := cfg.APIs[providerID]; ok {
		entry = mergeAPIEntry(existing, entry)
	}
	cfg.APIs[providerID] = entry
}

func mergeAPIEntry(existing, next config.APIEntry) config.APIEntry {
	if existing.Text != nil && next.Text != nil {
		merged := *existing.Text
		if strings.TrimSpace(merged.Protocol) == "" {
			merged.Protocol = next.Text.Protocol
		}
		if strings.TrimSpace(merged.BaseURL) == "" {
			merged.BaseURL = next.Text.BaseURL
		}
		if strings.TrimSpace(merged.APIKeyEnv) == "" {
			merged.APIKeyEnv = next.Text.APIKeyEnv
		}
		if merged.Defaults == nil {
			merged.Defaults = next.Text.Defaults
		} else if merged.Defaults.Model == "" && next.Text.Defaults != nil {
			merged.Defaults.Model = next.Text.Defaults.Model
		}
		if strings.TrimSpace(merged.Model) == "" {
			merged.Model = next.Text.Model
		}
		existing.Text = &merged
	} else if existing.Text == nil {
		existing.Text = next.Text
	}
	return existing
}
