package catalog

import (
	"context"
	"fmt"
	"time"

	wujicfg "github.com/coditary/wuji-core/pkg/config"
)

// MaybeAutoSync refreshes the catalog when auto_sync_hours is configured and the cache is stale.
func MaybeAutoSync(ctx context.Context, cfg *wujicfg.Config) error {
	if cfg == nil {
		return nil
	}
	catalogCfg := cfg.ResolvedCatalog()
	if catalogCfg.AutoSyncHours <= 0 {
		return nil
	}
	age, err := IndexAge(cfg.Root)
	if err == nil && age < time.Duration(catalogCfg.AutoSyncHours)*time.Hour {
		return nil
	}
	_, err = Sync(ctx, cfg.Root, catalogCfg.Source)
	return err
}

// LoadOrSync returns the cached index, optionally auto-syncing when configured.
func LoadOrSync(ctx context.Context, cfg *wujicfg.Config) (*Index, error) {
	if cfg != nil {
		if err := MaybeAutoSync(ctx, cfg); err != nil {
			return nil, fmt.Errorf("catalog auto-sync: %w", err)
		}
	}
	root := "."
	if cfg != nil && cfg.Root != "" {
		root = cfg.Root
	}
	return LoadIndex(root)
}

// NPMProtocol maps models.dev npm package names to Wuji API driver protocols.
func NPMProtocol(npm string) string {
	switch npm {
	case "@ai-sdk/anthropic":
		return "anthropic"
	case "@ai-sdk/openai-compatible", "@ai-sdk/openai", "@ai-sdk/azure", "@ai-sdk/deepinfra", "@ai-sdk/cerebras", "@ai-sdk/cohere", "@ai-sdk/amazon-bedrock":
		return "openai"
	default:
		if npm != "" {
			return "openai"
		}
		return ""
	}
}
