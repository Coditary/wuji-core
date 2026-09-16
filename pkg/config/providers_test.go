package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadProvidersFromConfig(t *testing.T) {
	dir := t.TempDir()
	wujiDir := filepath.Join(dir, ".wuji")
	if err := os.MkdirAll(wujiDir, 0o755); err != nil {
		t.Fatal(err)
	}
	raw := `default_provider: echo
providers:
  echo:
    type: driver
    driver: echo
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
	if cfg.DefaultProvider != "echo" {
		t.Fatalf("default_provider: %q", cfg.DefaultProvider)
	}
	if cfg.Providers["echo"].Driver != "echo" {
		t.Fatalf("providers: %+v", cfg.Providers)
	}
}
