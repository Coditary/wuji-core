package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const mcpFileName = "mcp.json"
const legacyMCPFileName = "mcp.yaml"

// MCPServersConfig holds MCP server definitions in Cursor-compatible format.
type MCPServersConfig struct {
	Servers map[string]MCPServerDef `json:"mcpServers" yaml:"servers"`
}

// MCPServerDef describes a single MCP server (stdio or remote URL).
type MCPServerDef struct {
	Command  string            `json:"command,omitempty" yaml:"command,omitempty"`
	Args     []string          `json:"args,omitempty" yaml:"args,omitempty"`
	Env      map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	Cwd      string            `json:"cwd,omitempty" yaml:"cwd,omitempty"`
	URL      string            `json:"url,omitempty" yaml:"url,omitempty"`
	Disabled bool              `json:"disabled,omitempty" yaml:"disabled,omitempty"`
}

// IsRemote reports whether the server is a remote HTTP/SSE endpoint.
func (d MCPServerDef) IsRemote() bool {
	return d.URL != ""
}

// IsStdio reports whether the server is started as a local subprocess.
func (d MCPServerDef) IsStdio() bool {
	return d.Command != ""
}

// Validate checks that the definition is usable.
func (d MCPServerDef) Validate() error {
	if d.IsRemote() && d.IsStdio() {
		return fmt.Errorf("cannot set both url and command")
	}
	if !d.IsRemote() && !d.IsStdio() {
		return fmt.Errorf("either command or url is required")
	}
	if d.IsStdio() && d.Command == "" {
		return fmt.Errorf("command is required for stdio servers")
	}
	return nil
}

// MCPConfigPath returns the path to the MCP config file.
func MCPConfigPath(root string) string {
	return filepath.Join(root, dirName, mcpFileName)
}

// ResolveMCPImportPath returns the file to import from.
// An explicit path is used as-is; otherwise common Cursor locations are tried.
func ResolveMCPImportPath(explicit string) (string, error) {
	if explicit != "" {
		return expandHome(explicit), nil
	}

	candidates := []string{".cursor/mcp.json"}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".cursor", "mcp.json"))
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, nil
		}
	}
	return "", fmt.Errorf("no mcp.json found — pass --import <path> or create .cursor/mcp.json")
}

// LoadMCP reads MCP server definitions from .wuji/mcp.json.
func LoadMCP(root string) (*MCPServersConfig, error) {
	path := MCPConfigPath(root)
	cfg := &MCPServersConfig{Servers: map[string]MCPServerDef{}}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return migrateLegacyMCP(root)
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if cfg.Servers == nil {
		cfg.Servers = map[string]MCPServerDef{}
	}
	return cfg, nil
}

func migrateLegacyMCP(root string) (*MCPServersConfig, error) {
	legacyPath := filepath.Join(root, dirName, legacyMCPFileName)
	data, err := os.ReadFile(legacyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &MCPServersConfig{Servers: map[string]MCPServerDef{}}, nil
		}
		return nil, fmt.Errorf("read %s: %w", legacyPath, err)
	}

	var legacy struct {
		Servers map[string]MCPServerDef `yaml:"servers"`
	}
	if err := yaml.Unmarshal(data, &legacy); err != nil {
		return nil, fmt.Errorf("parse %s: %w", legacyPath, err)
	}

	cfg := &MCPServersConfig{Servers: legacy.Servers}
	if cfg.Servers == nil {
		cfg.Servers = map[string]MCPServerDef{}
	}
	if err := SaveMCP(root, cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// SaveMCP writes MCP server definitions to .wuji/mcp.json.
func SaveMCP(root string, cfg *MCPServersConfig) error {
	dir := filepath.Join(root, dirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if cfg.Servers == nil {
		cfg.Servers = map[string]MCPServerDef{}
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(MCPConfigPath(root), data, 0o644)
}

// AddMCPServer adds or updates a server definition.
func AddMCPServer(root, name string, def MCPServerDef) error {
	if name == "" {
		return fmt.Errorf("server name is required")
	}
	if err := def.Validate(); err != nil {
		return err
	}

	cfg, err := LoadMCP(root)
	if err != nil {
		return err
	}
	cfg.Servers[name] = def
	return SaveMCP(root, cfg)
}

// RemoveMCPServer removes a server definition.
func RemoveMCPServer(root, name string) error {
	cfg, err := LoadMCP(root)
	if err != nil {
		return err
	}
	if _, ok := cfg.Servers[name]; !ok {
		return fmt.Errorf("MCP server %q not found", name)
	}
	delete(cfg.Servers, name)
	return SaveMCP(root, cfg)
}

// GetMCPServer returns a single server definition.
func GetMCPServer(root, name string) (MCPServerDef, error) {
	cfg, err := LoadMCP(root)
	if err != nil {
		return MCPServerDef{}, err
	}
	def, ok := cfg.Servers[name]
	if !ok {
		return MCPServerDef{}, fmt.Errorf("MCP server %q not found", name)
	}
	return def, nil
}

// ImportMCP merges servers from another Cursor-compatible mcp.json file.
func ImportMCP(root, jsonPath string, overwrite bool) (int, error) {
	path := expandHome(jsonPath)
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read %s: %w", path, err)
	}

	var src MCPServersConfig
	if err := json.Unmarshal(data, &src); err != nil {
		return 0, fmt.Errorf("parse %s: %w", path, err)
	}
	if len(src.Servers) == 0 {
		return 0, fmt.Errorf("no mcpServers found in %s", path)
	}

	cfg, err := LoadMCP(root)
	if err != nil {
		return 0, err
	}

	imported := 0
	for name, def := range src.Servers {
		if !overwrite {
			if _, exists := cfg.Servers[name]; exists {
				continue
			}
		}
		if err := def.Validate(); err != nil {
			return imported, fmt.Errorf("server %q: %w", name, err)
		}
		cfg.Servers[name] = def
		imported++
	}

	if imported == 0 {
		return 0, nil
	}
	if err := SaveMCP(root, cfg); err != nil {
		return 0, err
	}
	return imported, nil
}

func expandHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if path == "~" {
		return home
	}
	return filepath.Join(home, strings.TrimPrefix(path, "~/"))
}
