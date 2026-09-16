package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateProvidersToAPIs(t *testing.T) {
	cfg := &Config{
		Providers: map[string]ProviderConfig{
			"claude": {Type: "anthropic", Model: "claude-sonnet", APIKeyEnv: "ANTHROPIC_API_KEY"},
			"echo":   {Type: "driver", Driver: "echo"},
		},
	}
	cfg.migrateProvidersToAPIs()
	if _, ok := cfg.APIs["claude"]; !ok {
		t.Fatal("expected claude api")
	}
	if _, ok := cfg.APIs["echo"]; ok {
		t.Fatal("driver providers should not migrate")
	}
}

func TestLoadAPIsFromConfig(t *testing.T) {
	dir := t.TempDir()
	wujiDir := filepath.Join(dir, ".wuji")
	if err := os.MkdirAll(wujiDir, 0o755); err != nil {
		t.Fatal(err)
	}
	raw := `apis:
  demo:
    text:
      protocol: openai
      model: gpt-4o
`
	if err := os.WriteFile(filepath.Join(wujiDir, "config.yaml"), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIs["demo"].Text.Protocol != "openai" {
		t.Fatalf("apis: %+v", cfg.APIs)
	}
}
