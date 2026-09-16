package config_test

import (
	"strings"
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
)

func TestDefaultDriverEndpointUsesUnixSocket(t *testing.T) {
	ep := config.DefaultDriverEndpoint("/tmp/wuji", "llama")
	if !strings.HasPrefix(ep, "unix://") {
		t.Fatalf("expected unix endpoint, got %q", ep)
	}
	if !strings.HasSuffix(ep, "/.wuji/run/drivers/llama.sock") {
		t.Fatalf("unexpected socket path: %q", ep)
	}
}

func TestDriverGRPCAddrPrefersExplicitOverride(t *testing.T) {
	cfg := &config.Config{
		Root: "/tmp/wuji",
		Drivers: config.DriversConfig{
			A1111: &config.A1111Config{GRPCAddr: "127.0.0.1:59999"},
		},
	}
	if cfg.DriverGRPCAddr("a1111") != "127.0.0.1:59999" {
		t.Fatalf("expected explicit override")
	}
}
