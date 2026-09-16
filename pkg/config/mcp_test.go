package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
)

func TestAddAndLoadMCPServer(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	def := config.MCPServerDef{
		Command: "npx",
		Args:    []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"},
	}
	if err := config.AddMCPServer(dir, "filesystem", def); err != nil {
		t.Fatalf("AddMCPServer: %v", err)
	}

	cfg, err := config.LoadMCP(dir)
	if err != nil {
		t.Fatalf("LoadMCP: %v", err)
	}
	got, ok := cfg.Servers["filesystem"]
	if !ok {
		t.Fatal("expected filesystem server")
	}
	if got.Command != "npx" || len(got.Args) != 3 {
		t.Fatalf("unexpected def: %+v", got)
	}

	data, err := os.ReadFile(config.MCPConfigPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]map[string]config.MCPServerDef
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	if raw["mcpServers"]["filesystem"].Command != "npx" {
		t.Fatalf("expected cursor-compatible json format, got %s", string(data))
	}
}

func TestRemoveMCPServer(t *testing.T) {
	dir := t.TempDir()
	def := config.MCPServerDef{Command: "echo", Args: []string{"hi"}}
	if err := config.AddMCPServer(dir, "demo", def); err != nil {
		t.Fatal(err)
	}
	if err := config.RemoveMCPServer(dir, "demo"); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.LoadMCP(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Servers) != 0 {
		t.Fatalf("expected empty servers, got %+v", cfg.Servers)
	}
}

func TestImportMCP(t *testing.T) {
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "mcp.json")
	content := `{
  "mcpServers": {
    "echo": {
      "command": "echo",
      "args": ["hello"]
    },
    "remote": {
      "url": "http://127.0.0.1:9999"
    }
  }
}`
	if err := os.WriteFile(jsonPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	n, err := config.ImportMCP(dir, jsonPath, false)
	if err != nil {
		t.Fatalf("ImportMCP: %v", err)
	}
	if n != 2 {
		t.Fatalf("expected 2 imported, got %d", n)
	}

	cfg, err := config.LoadMCP(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Servers["echo"].IsStdio() || !cfg.Servers["remote"].IsRemote() {
		t.Fatalf("unexpected servers: %+v", cfg.Servers)
	}
}

func TestResolveMCPImportPathAutoDetect(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)

	cursorDir := filepath.Join(dir, ".cursor")
	if err := os.MkdirAll(cursorDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cursorDir, "mcp.json"), []byte(`{"mcpServers":{}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	path, err := config.ResolveMCPImportPath("")
	if err != nil {
		t.Fatal(err)
	}
	if path != ".cursor/mcp.json" {
		t.Fatalf("expected .cursor/mcp.json, got %q", path)
	}
}

func TestMCPServerDefValidate(t *testing.T) {
	if err := (config.MCPServerDef{Command: "npx"}).Validate(); err != nil {
		t.Fatalf("expected valid stdio def: %v", err)
	}
	if err := (config.MCPServerDef{URL: "http://localhost/mcp"}).Validate(); err != nil {
		t.Fatalf("expected valid remote def: %v", err)
	}
	if err := (config.MCPServerDef{}).Validate(); err == nil {
		t.Fatal("expected error for empty def")
	}
	if err := (config.MCPServerDef{Command: "npx", URL: "http://x"}).Validate(); err == nil {
		t.Fatal("expected error for mixed def")
	}
}

func TestMigrateLegacyMCPYAML(t *testing.T) {
	dir := t.TempDir()
	wujiDir := filepath.Join(dir, ".wuji")
	if err := os.MkdirAll(wujiDir, 0o755); err != nil {
		t.Fatal(err)
	}
	legacy := `servers:
  demo:
    command: sleep
    args: ["1"]
`
	if err := os.WriteFile(filepath.Join(wujiDir, "mcp.yaml"), []byte(legacy), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadMCP(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Servers["demo"].Command != "sleep" {
		t.Fatalf("unexpected migrated config: %+v", cfg.Servers)
	}
	if _, err := os.Stat(config.MCPConfigPath(dir)); err != nil {
		t.Fatal("expected mcp.json to be created from legacy yaml")
	}
}
