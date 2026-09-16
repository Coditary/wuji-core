package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
)

func TestDriverEndpointsPersistInConfigYAML(t *testing.T) {
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
	if err := cfg.AddDriverEndpoint("127.0.0.1:50054"); err != nil {
		t.Fatalf("AddDriverEndpoint: %v", err)
	}

	reloaded, err := config.Load()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(reloaded.DriverEndpoints) != 1 || reloaded.DriverEndpoints[0] != "127.0.0.1:50054" {
		t.Fatalf("unexpected endpoints: %#v", reloaded.DriverEndpoints)
	}

	data, err := os.ReadFile(filepath.Join(dir, ".wuji", "config.yaml"))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(data), "driver_endpoints") {
		t.Fatalf("expected driver_endpoints in config.yaml, got:\n%s", data)
	}
}

func TestLegacyDriversYAMLMigratesIntoConfig(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	wujiDir := filepath.Join(dir, ".wuji")
	if err := os.MkdirAll(wujiDir, 0o755); err != nil {
		t.Fatal(err)
	}
	legacy := "drivers:\n  - endpoint: 127.0.0.1:50053\n"
	if err := os.WriteFile(filepath.Join(wujiDir, "drivers.yaml"), []byte(legacy), 0o644); err != nil {
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
	if len(cfg.DriverEndpoints) != 1 || cfg.DriverEndpoints[0] != "127.0.0.1:50053" {
		t.Fatalf("unexpected endpoints: %#v", cfg.DriverEndpoints)
	}
	if _, err := os.Stat(filepath.Join(wujiDir, "drivers.yaml")); !os.IsNotExist(err) {
		t.Fatal("expected legacy drivers.yaml to be removed after migration")
	}
}

func TestResolvedDriverEndpointsIncludesConfiguredGRPCAddr(t *testing.T) {
	cfg := &config.Config{
		Root:              "/tmp/wuji",
		CapabilityDrivers: config.CapabilityDrivers{"image": "a1111"},
		Drivers: config.DriversConfig{
			A1111: &config.A1111Config{},
		},
		DriverEndpoints: []string{"127.0.0.1:50054"},
	}
	eps := cfg.ResolvedDriverEndpoints()
	if len(eps) != 2 {
		t.Fatalf("expected 2 endpoints, got %#v", eps)
	}
}
