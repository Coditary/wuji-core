package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	wujicfg "github.com/coditary/wuji-core/pkg/config"
)

// OpenCodeConfig is the MCP section of opencode.json.
type OpenCodeConfig struct {
	Schema string                       `json:"$schema,omitempty"`
	MCP    map[string]OpenCodeMCPServer `json:"mcp"`
}

// OpenCodeMCPServer is a single OpenCode MCP server entry.
type OpenCodeMCPServer struct {
	Type        string            `json:"type"`
	Command     []string          `json:"command,omitempty"`
	URL         string            `json:"url,omitempty"`
	Enabled     bool              `json:"enabled,omitempty"`
	OAuth       *bool             `json:"oauth,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

// BuildOpenCodeConfig generates OpenCode MCP entries from Wuji definitions.
func BuildOpenCodeConfig(root, wujiBin string, transport Transport, gw GatewayOptions) (*OpenCodeConfig, error) {
	cfg, err := wujicfg.LoadMCP(root)
	if err != nil {
		return nil, err
	}
	if len(cfg.Servers) == 0 {
		return &OpenCodeConfig{MCP: map[string]OpenCodeMCPServer{}}, nil
	}

	out := &OpenCodeConfig{
		Schema: "https://opencode.ai/config.json",
		MCP:    make(map[string]OpenCodeMCPServer, len(cfg.Servers)),
	}

	offset := 0
	for name, def := range cfg.Servers {
		if def.Disabled {
			continue
		}
		if def.IsRemote() {
			oauth := false
			out.MCP[name] = OpenCodeMCPServer{
				Type:    "remote",
				URL:     def.URL,
				Enabled: true,
				OAuth:   &oauth,
			}
			continue
		}

		switch transport {
		case TransportHTTP, TransportSSE:
			resolved := ResolveGateway(name, gw, offset)
			oauth := false
			out.MCP[name] = OpenCodeMCPServer{
				Type:    "remote",
				URL:     resolved.URL,
				Enabled: true,
				OAuth:   &oauth,
			}
			offset++
		default:
			entry := OpenCodeMCPServer{
				Type:    "local",
				Command: []string{wujiBin, "mcp", "start", name, "--stdio"},
				Enabled: true,
			}
			if len(def.Env) > 0 {
				entry.Environment = def.Env
			}
			out.MCP[name] = entry
		}
	}

	return out, nil
}

// WriteOpenCodeConfig merges MCP entries into opencode.json at project root.
func WriteOpenCodeConfig(root string, cfg *OpenCodeConfig) (string, error) {
	path := filepath.Join(root, "opencode.json")
	merged := map[string]json.RawMessage{}
	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, &merged); err != nil {
			return "", fmt.Errorf("parse %s: %w", path, err)
		}
	} else if !os.IsNotExist(err) {
		return "", err
	}

	mcpData, err := json.Marshal(cfg.MCP)
	if err != nil {
		return "", err
	}
	merged["$schema"] = json.RawMessage(`"https://opencode.ai/config.json"`)
	merged["mcp"] = mcpData

	data, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return "", err
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// WujiBinary returns the path to the current wuji executable for OpenCode command arrays.
func WujiBinary() string {
	exe, err := os.Executable()
	if err != nil {
		return "wuji"
	}
	return exe
}
