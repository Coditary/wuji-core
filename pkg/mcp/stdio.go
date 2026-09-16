package mcp

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	wujicfg "github.com/coditary/wuji-core/pkg/config"
)

// RunStdio runs an MCP server in the foreground, proxying stdin/stdout for MCP clients
// such as OpenCode (type: local, command: ["wuji", "mcp", "start", NAME, "--stdio"]).
func (m *Manager) RunStdio(ctx context.Context, name string) error {
	def, err := wujicfg.GetMCPServer(m.root, name)
	if err != nil {
		return err
	}
	if def.Disabled {
		return fmt.Errorf("MCP server %q is disabled", name)
	}

	var cmd *exec.Cmd
	if def.IsRemote() {
		cmd = exec.CommandContext(ctx, "npx", "-y", "mcp-remote", def.URL)
		cmd.Env = buildEnv(def.Env)
	} else {
		cmd = exec.CommandContext(ctx, def.Command, def.Args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if def.Cwd != "" {
			cmd.Dir = def.Cwd
		}
		cmd.Env = buildEnv(def.Env)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("stdio proxy for %q: %w", name, err)
	}
	return nil
}
