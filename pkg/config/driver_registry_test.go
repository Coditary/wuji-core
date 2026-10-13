package config

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDriversConfigUnmarshalRegistry(t *testing.T) {
	const raw = `
resources:
  total_vram_mb: 8192
ollama:
  driver_bin: /tmp/wuji-driver-ollama
  server_bin: ollama
  manage_server: true
my-backend:
  driver_bin: /tmp/wuji-driver-my-backend
  host: 127.0.0.1
  port: 9999
`
	var drivers DriversConfig
	if err := yaml.Unmarshal([]byte(raw), &drivers); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if drivers.Resources == nil || drivers.Resources.TotalVRAMMB != 8192 {
		t.Fatalf("resources not parsed: %#v", drivers.Resources)
	}
	ollama := drivers.Section("ollama")
	if ollama == nil || ollama.DriverBin != "/tmp/wuji-driver-ollama" {
		t.Fatalf("ollama section: %#v", ollama)
	}
	backend := drivers.Section("my-backend")
	if backend == nil || backend.Port != 9999 {
		t.Fatalf("my-backend section: %#v", backend)
	}
}

func TestMergeDriverRegistryOverlay(t *testing.T) {
	cfg := &Config{
		Drivers: DriversConfig{
			Registry: map[string]DriverSection{
				"ollama": {DriverBin: "/old/bin"},
			},
		},
	}
	mergeDriversOverlay(cfg, DriversConfig{
		Registry: map[string]DriverSection{
			"ollama": {ServerBin: "ollama", ManageServer: boolPtr(true)},
			"tgi":    {DriverBin: "/tmp/wuji-driver-tgi"},
		},
	})
	ollama := cfg.Drivers.Section("ollama")
	if ollama == nil || ollama.DriverBin != "/old/bin" || ollama.ServerBin != "ollama" {
		t.Fatalf("merged ollama: %#v", ollama)
	}
	tgi := cfg.Drivers.Section("tgi")
	if tgi == nil || tgi.DriverBin != "/tmp/wuji-driver-tgi" {
		t.Fatalf("merged tgi: %#v", tgi)
	}
}

func TestConfiguredDriverIDsIncludesRegistry(t *testing.T) {
	cfg := &Config{
		CapabilityDrivers: CapabilityDrivers{"text": "ollama"},
		Drivers: DriversConfig{
			Registry: map[string]DriverSection{
				"ollama": {DriverBin: "/tmp/bin"},
			},
		},
	}
	ids := cfg.configuredDriverIDs()
	if len(ids) != 1 || ids[0] != "ollama" {
		t.Fatalf("unexpected ids: %#v", ids)
	}
}

func TestDriverLifecycleUsesRegistry(t *testing.T) {
	auto := true
	cfg := &Config{
		Root: "/tmp/wuji",
		Drivers: DriversConfig{
			Registry: map[string]DriverSection{
				"ollama": {
					DriverBin:       "/tmp/wuji-driver-ollama",
					AutoStartDriver: &auto,
					GRPCAddr:        "unix:///tmp/ollama.sock",
				},
			},
		},
	}
	if !cfg.AutoStartDriverEnabled("ollama") {
		t.Fatal("expected auto start for registry driver")
	}
	if cfg.DriverBinaryPath("ollama") != "/tmp/wuji-driver-ollama" {
		t.Fatalf("binary path: %q", cfg.DriverBinaryPath("ollama"))
	}
	if cfg.DriverGRPCAddr("ollama") != "unix:///tmp/ollama.sock" {
		t.Fatalf("grpc addr: %q", cfg.DriverGRPCAddr("ollama"))
	}
}

func TestReqpackOverlayUnmarshalIntoRegistry(t *testing.T) {
	const raw = `
capability_drivers:
  text: ollama
drivers:
  ollama:
    driver_bin: /home/user/.local/share/wuji/bin/wuji-driver-ollama
    auto_start_driver: true
    server_bin: ollama
    manage_server: true
`
	var overlay fileConfig
	if err := yaml.Unmarshal([]byte(raw), &overlay); err != nil {
		t.Fatalf("unmarshal overlay: %v", err)
	}
	cfg := &Config{}
	mergeDriversOverlay(cfg, overlay.Drivers)
	if cfg.CapabilityDrivers != nil {
		t.Fatalf("capability drivers should be merged separately, got %#v", cfg.CapabilityDrivers)
	}
	section := cfg.Drivers.Section("ollama")
	if section == nil || !strings.Contains(section.DriverBin, "wuji-driver-ollama") {
		t.Fatalf("overlay section: %#v", section)
	}
}

func boolPtr(v bool) *bool { return &v }
