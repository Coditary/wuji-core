package catalog

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
)

func TestAPIEntryFromProvider(t *testing.T) {
	entry := APIEntryFromProvider(ProviderEntry{
		ID:   "github-copilot",
		Name: "GitHub Copilot",
		API:  "https://api.githubcopilot.com",
		NPM:  "@ai-sdk/openai-compatible",
		Env:  []string{"GITHUB_TOKEN"},
		Models: map[string]string{
			"claude-opus-5": "Claude Opus 5",
		},
	})
	if entry.Text == nil {
		t.Fatal("expected text spec")
	}
	if entry.Text.Protocol != "openai" {
		t.Fatalf("protocol = %q", entry.Text.Protocol)
	}
	if entry.Text.APIKeyEnv != "GITHUB_TOKEN" {
		t.Fatalf("api_key_env = %q", entry.Text.APIKeyEnv)
	}
	if entry.Text.Defaults == nil || entry.Text.Defaults.Model != "claude-opus-5" {
		t.Fatalf("defaults = %+v", entry.Text.Defaults)
	}
}

func TestEnsureAPIRoutePreservesExisting(t *testing.T) {
	cfg := &config.Config{
		APIs: map[string]config.APIEntry{
			"anthropic": {
				Text: &config.APICapabilitySpec{
					Protocol:  "anthropic",
					APIKeyEnv: "ANTHROPIC_API_KEY",
					Defaults:  &config.APIDefaults{Model: "claude-opus-5"},
				},
			},
		},
	}
	EnsureAPIRoute(cfg, "anthropic", ProviderEntry{
		ID: "anthropic", API: "https://api.anthropic.com", NPM: "@ai-sdk/anthropic",
		Env:    []string{"ANTHROPIC_API_KEY"},
		Models: map[string]string{"claude-opus-5": "Claude Opus 5"},
	})
	if cfg.APIs["anthropic"].Text.Defaults.Model != "claude-opus-5" {
		t.Fatalf("model overwritten: %+v", cfg.APIs["anthropic"].Text)
	}
}
