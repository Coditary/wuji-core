package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
	"gopkg.in/yaml.v3"
)

func TestInitResourcesPersists(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	oldWD, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWD)

	res, err := config.InitResources(dir)
	if err != nil {
		t.Fatalf("InitResources: %v", err)
	}
	if !res.Enabled() {
		t.Fatal("expected enabled after init")
	}

	data, err := os.ReadFile(filepath.Join(dir, ".wuji", "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	var parsed struct {
		Drivers struct {
			Resources struct {
				TotalRAMMB  int `yaml:"total_ram_mb"`
				TotalVRAMMB int `yaml:"total_vram_mb"`
			} `yaml:"resources"`
		} `yaml:"drivers"`
	}
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed.Drivers.Resources.TotalRAMMB != res.TotalRAMMB {
		t.Fatalf("persisted RAM %d != resolved %d", parsed.Drivers.Resources.TotalRAMMB, res.TotalRAMMB)
	}
}

func TestSetResourcesField(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	oldWD, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWD)

	if err := config.SetResourcesField(dir, "total_ram_mb", "48000"); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	res := cfg.ResolvedResources()
	if res.TotalRAMMB != 48000 {
		t.Fatalf("expected 48000, got %d", res.TotalRAMMB)
	}
}

func TestSetResourcesIdleShutdown(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	oldWD, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWD)

	if err := config.SetResourcesField(dir, "idle_shutdown_seconds", "120"); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	res := cfg.ResolvedResources()
	if res.IdleShutdownSeconds != 120 {
		t.Fatalf("expected 120, got %d", res.IdleShutdownSeconds)
	}
}
