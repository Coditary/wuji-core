package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
)

func TestDefaultVLLMWithoutFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cwd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	vllm := cfg.ResolvedVLLM()
	if vllm.Host != "127.0.0.1" || vllm.Port != 8000 {
		t.Fatalf("unexpected defaults: %+v", vllm)
	}
	if !vllm.ManageServerEnabled() {
		t.Fatal("expected manage_server true by default")
	}
	if vllm.ResolvedAPIBase() != "http://127.0.0.1:8000/v1" {
		t.Fatalf("unexpected api base: %s", vllm.ResolvedAPIBase())
	}
}

func TestSetVLLMFieldPersists(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cwd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)

	if err := config.SetVLLMField(dir, "default_model", "demo-model"); err != nil {
		t.Fatalf("SetVLLMField: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.ResolvedVLLM().DefaultModel != "demo-model" {
		t.Fatalf("expected persisted model, got %q", cfg.ResolvedVLLM().DefaultModel)
	}
}

func TestLoadUsesWujiRootWithoutGoMod(t *testing.T) {
	wujiHome := t.TempDir()
	novelDir := t.TempDir()
	t.Setenv("WUJI_ROOT", wujiHome)
	cwd, _ := os.Getwd()
	if err := os.Chdir(novelDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Root != wujiHome {
		t.Fatalf("root=%q want %q", cfg.Root, wujiHome)
	}
}
