package core_test

import (
	"context"
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/core"
)

func TestConnectRuntimeEmbeddedByDefault(t *testing.T) {
	t.Setenv("WUJI_DAEMON", "")
	t.Setenv("WUJI_EMBEDDED", "")

	cfg := &config.Config{Root: t.TempDir()}
	rt, cleanup, err := core.ConnectRuntime(context.Background(), cfg)
	if err != nil {
		t.Fatalf("ConnectRuntime: %v", err)
	}
	if cleanup == nil {
		t.Fatal("expected cleanup func")
	}
	if rt == nil {
		t.Fatal("expected runtime")
	}
	cleanup()
}

func TestConnectRuntimeDaemonMode(t *testing.T) {
	t.Setenv("WUJI_DAEMON", "1")
	t.Setenv("WUJI_CORE_BIN", "/nonexistent/wuji-core")

	cfg := &config.Config{Root: t.TempDir()}
	_, _, err := core.ConnectRuntime(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected daemon connect error")
	}
}
