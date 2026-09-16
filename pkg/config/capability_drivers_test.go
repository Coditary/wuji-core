package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
)

func TestSetCapabilityDriverPersists(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cwd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)

	if err := config.SetCapabilityDriver(dir, "text", "vllm"); err != nil {
		t.Fatalf("SetCapabilityDriver: %v", err)
	}
	if err := config.SetCapabilityDriver(dir, "image", "dummy"); err != nil {
		t.Fatalf("SetCapabilityDriver image: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.DriverForCapability(capability.TextGeneration) != "vllm" {
		t.Fatalf("expected text driver vllm, got %q", cfg.DriverForCapability(capability.TextGeneration))
	}
	if cfg.DriverForCapability(capability.ImageGeneration) != "dummy" {
		t.Fatalf("expected image driver dummy, got %q", cfg.DriverForCapability(capability.ImageGeneration))
	}
}

func TestSetCapabilityDriverRejectsUnknownCapability(t *testing.T) {
	dir := t.TempDir()
	if err := config.SetCapabilityDriver(dir, "unknown", "dummy"); err == nil {
		t.Fatal("expected error for unknown capability")
	}
}
