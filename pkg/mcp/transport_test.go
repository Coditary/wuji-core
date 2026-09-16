package mcp_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/mcp"
)

func TestResolveGatewayDefaults(t *testing.T) {
	gw := mcp.ResolveGateway("filesystem", mcp.GatewayOptions{}, 1)
	if gw.Port != 9334 {
		t.Fatalf("expected port 9334, got %d", gw.Port)
	}
	if gw.Path != "/mcp/filesystem" {
		t.Fatalf("unexpected path %q", gw.Path)
	}
	if gw.URL != "http://127.0.0.1:9334/mcp/filesystem" {
		t.Fatalf("unexpected url %q", gw.URL)
	}
}

func TestBuildOpenCodeConfigStdio(t *testing.T) {
	dir := t.TempDir()
	def := config.MCPServerDef{Command: "npx", Args: []string{"-y", "demo"}}
	if err := config.AddMCPServer(dir, "demo", def); err != nil {
		t.Fatal(err)
	}

	cfg, err := mcp.BuildOpenCodeConfig(dir, "/usr/bin/wuji", mcp.TransportStdio, mcp.GatewayOptions{})
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := cfg.MCP["demo"]
	if !ok {
		t.Fatal("missing demo entry")
	}
	if entry.Type != "local" {
		t.Fatalf("expected local, got %q", entry.Type)
	}
	want := []string{"/usr/bin/wuji", "mcp", "start", "demo", "--stdio"}
	if len(entry.Command) != len(want) {
		t.Fatalf("command %v", entry.Command)
	}
	for i := range want {
		if entry.Command[i] != want[i] {
			t.Fatalf("command[%d]=%q want %q", i, entry.Command[i], want[i])
		}
	}
}

func TestBuildOpenCodeConfigHTTP(t *testing.T) {
	dir := t.TempDir()
	def := config.MCPServerDef{Command: "npx", Args: []string{"-y", "demo"}}
	if err := config.AddMCPServer(dir, "demo", def); err != nil {
		t.Fatal(err)
	}

	cfg, err := mcp.BuildOpenCodeConfig(dir, "/usr/bin/wuji", mcp.TransportHTTP, mcp.GatewayOptions{})
	if err != nil {
		t.Fatal(err)
	}
	entry := cfg.MCP["demo"]
	if entry.Type != "remote" {
		t.Fatalf("expected remote, got %q", entry.Type)
	}
	if entry.URL != "http://127.0.0.1:9333/mcp/demo" {
		t.Fatalf("unexpected url %q", entry.URL)
	}
	if entry.OAuth == nil || *entry.OAuth {
		t.Fatal("expected oauth false")
	}
}
