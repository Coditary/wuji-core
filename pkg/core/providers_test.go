package core_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/core"
	"github.com/coditary/wuji-core/pkg/driver/dummy"
)

func TestCoreListProvidersIncludesDrivers(t *testing.T) {
	c := testCore(t)
	defer c.Close()

	catalog := c.ListProviders()
	if len(catalog.Providers) == 0 {
		t.Fatal("expected providers in catalog")
	}
	if _, ok := catalog.Providers[dummy.DriverID]; !ok {
		t.Fatalf("expected dummy driver in catalog, got: %v", catalog.Providers)
	}
	if catalog.Providers[dummy.DriverID].Kind != "driver" {
		t.Fatalf("unexpected kind: %q", catalog.Providers[dummy.DriverID].Kind)
	}
}

func TestCoreListProvidersMergesAPIsAndLegacy(t *testing.T) {
	c, err := core.New(core.Config{
		AppConfig: &config.Config{
			DefaultProvider: "echo",
			APIs: map[string]config.APIEntry{
				"anthropic": {Text: &config.APICapabilitySpec{Model: "claude-sonnet"}},
			},
			Providers: map[string]config.ProviderConfig{
				"echo": {Type: "driver", Driver: "echo"},
			},
		},
		DefaultDriverID: dummy.DriverID,
		Lazy:            true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	catalog := c.ListProviders()
	if catalog.Providers["anthropic"].Kind != "api" {
		t.Fatalf("anthropic: %+v", catalog.Providers["anthropic"])
	}
	if catalog.Providers["echo"].Kind != "provider" {
		t.Fatalf("echo: %+v", catalog.Providers["echo"])
	}
	if _, ok := catalog.Providers[dummy.DriverID]; !ok {
		t.Fatal("expected registered dummy driver")
	}
}
