package mcp_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/mcp"
)

func TestStartStopSleepServer(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	def := config.MCPServerDef{
		Command: "sleep",
		Args:    []string{"60"},
	}
	if err := config.AddMCPServer(dir, "sleeper", def); err != nil {
		t.Fatal(err)
	}

	mgr := mcp.NewManager(dir)
	ctx := context.Background()

	if err := mgr.Start(ctx, "sleeper"); err != nil {
		t.Fatalf("Start: %v", err)
	}

	st, err := mgr.Status(ctx, "sleeper")
	if err != nil {
		t.Fatal(err)
	}
	if !st.Running || st.PID <= 0 {
		t.Fatalf("expected running server, got %+v", st)
	}

	if err := mgr.Stop("sleeper"); err != nil {
		t.Fatalf("Stop: %v", err)
	}

	time.Sleep(200 * time.Millisecond)
	st, err = mgr.Status(ctx, "sleeper")
	if err != nil {
		t.Fatal(err)
	}
	if st.Running {
		t.Fatalf("expected stopped server, got %+v", st)
	}
}

func TestStartAlreadyRunning(t *testing.T) {
	dir := t.TempDir()
	def := config.MCPServerDef{Command: "sleep", Args: []string{"60"}}
	if err := config.AddMCPServer(dir, "sleeper", def); err != nil {
		t.Fatal(err)
	}

	mgr := mcp.NewManager(dir)
	ctx := context.Background()
	if err := mgr.Start(ctx, "sleeper"); err != nil {
		t.Fatal(err)
	}
	if err := mgr.Start(ctx, "sleeper"); err != nil {
		t.Fatalf("second Start should be no-op: %v", err)
	}
	_ = mgr.Stop("sleeper")
}
