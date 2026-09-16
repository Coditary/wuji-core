package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
)

func TestDefaultLlamaWithoutFile(t *testing.T) {
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

	llama := cfg.ResolvedLlama()
	if llama.InferenceHost != "127.0.0.1" || llama.InferencePort != 8080 {
		t.Fatalf("unexpected defaults: %+v", llama)
	}
	if llama.ServerBin != "../../plugins/wuji/llama/vendor/llama/llama-server" {
		t.Fatalf("unexpected server_bin: %s", llama.ServerBin)
	}
}

func TestSetLlamaFieldPersists(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cwd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)

	if err := config.SetLlamaField(dir, "default_model", "demo.gguf"); err != nil {
		t.Fatalf("SetLlamaField: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.ResolvedLlama().DefaultModel != "demo.gguf" {
		t.Fatalf("expected persisted model, got %q", cfg.ResolvedLlama().DefaultModel)
	}
}

func TestOllamaThinkDefaultsWithAPI(t *testing.T) {
	useAPI := true
	llama := config.LlamaConfig{UseOllamaAPI: &useAPI}
	if !llama.OllamaThinkEnabled() {
		t.Fatal("expected think enabled when use_ollama_api is true and ollama_think unset")
	}

	disabled := false
	llama.OllamaThink = &disabled
	if llama.OllamaThinkEnabled() {
		t.Fatal("expected explicit ollama_think: false to disable thinking")
	}
}
