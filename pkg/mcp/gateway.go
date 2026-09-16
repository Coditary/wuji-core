package mcp

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	wujicfg "github.com/coditary/wuji-core/pkg/config"
)

// StartGateway exposes a stdio MCP server over HTTP or SSE using supergateway.
func (m *Manager) StartGateway(ctx context.Context, name string, transport Transport, opts GatewayOptions, portOffset int) (ResolvedGateway, error) {
	def, err := wujicfg.GetMCPServer(m.root, name)
	if err != nil {
		return ResolvedGateway{}, err
	}
	if def.Disabled {
		return ResolvedGateway{}, fmt.Errorf("MCP server %q is disabled", name)
	}

	gw := ResolveGateway(name, opts, portOffset)

	if def.IsRemote() {
		if !m.checkRemote(ctx, def.URL) {
			return gw, fmt.Errorf("remote MCP server %q at %s is not reachable", name, def.URL)
		}
		gw.URL = def.URL
		return gw, nil
	}

	st, err := m.Status(ctx, name)
	if err != nil {
		return ResolvedGateway{}, err
	}
	if st.Running && st.GatewayURL != "" {
		gw.URL = st.GatewayURL
		return gw, nil
	}
	if st.Running {
		_ = m.Stop(name)
	}

	if err := os.MkdirAll(m.runtimeDir(), 0o755); err != nil {
		return ResolvedGateway{}, fmt.Errorf("create runtime dir: %w", err)
	}

	logFile := m.logPath(name)
	log, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return ResolvedGateway{}, fmt.Errorf("open log file: %w", err)
	}

	stdioCmd := buildStdioShellCommand(def)
	gatewayArgs := []string{
		"-y", "supergateway",
		"--stdio", stdioCmd,
		"--port", fmt.Sprintf("%d", gw.Port),
		"--baseUrl", fmt.Sprintf("http://%s:%d", gw.Host, gw.Port),
		"--outputTransport", gatewayTransportName(transport),
	}
	switch transport {
	case TransportHTTP:
		gatewayArgs = append(gatewayArgs, "--streamableHttpPath", gw.Path)
	case TransportSSE:
		ssePath := gw.Path + "/sse"
		if gw.Path == "/"+name {
			ssePath = "/sse"
		}
		gatewayArgs = append(gatewayArgs, "--ssePath", ssePath, "--messagePath", "/message")
		gw.URL = fmt.Sprintf("http://%s:%d%s", gw.Host, gw.Port, ssePath)
	}

	cmd := exec.CommandContext(ctx, "npx", gatewayArgs...)
	cmd.Stdout = log
	cmd.Stderr = log
	cmd.Env = buildEnv(def.Env)

	if err := cmd.Start(); err != nil {
		_ = log.Close()
		return ResolvedGateway{}, fmt.Errorf("start gateway for %s: %w", name, err)
	}

	state := RuntimeState{
		PID:         cmd.Process.Pid,
		StartedAt:   time.Now().UTC(),
		Command:     "npx " + strings.Join(gatewayArgs, " "),
		LogFile:     logFile,
		Transport:   string(transport),
		GatewayURL:  gw.URL,
		GatewayHost: gw.Host,
		GatewayPort: gw.Port,
		GatewayPath: gw.Path,
	}
	if err := m.saveState(name, state); err != nil {
		_ = cmd.Process.Kill()
		_ = log.Close()
		return ResolvedGateway{}, err
	}

	go func() {
		_ = cmd.Wait()
		_ = log.Close()
		if s, _ := m.loadState(name); s != nil && s.PID == state.PID {
			_ = m.removeState(name)
		}
	}()

	time.Sleep(300 * time.Millisecond)
	if !processAlive(state.PID) {
		_ = m.removeState(name)
		return ResolvedGateway{}, fmt.Errorf("MCP gateway for %q exited immediately — check %s", name, logFile)
	}

	return gw, nil
}

// StartAllGateways starts gateways for all enabled stdio servers.
func (m *Manager) StartAllGateways(ctx context.Context, transport Transport, opts GatewayOptions) ([]string, []ResolvedGateway, []error) {
	cfg, err := wujicfg.LoadMCP(m.root)
	if err != nil {
		return nil, nil, []error{err}
	}

	var started []string
	var gateways []ResolvedGateway
	var errs []error
	offset := 0
	for name, def := range cfg.Servers {
		if def.Disabled {
			continue
		}
		perServer := opts
		if opts.Port > 0 {
			perServer.Port = opts.Port + offset
		}
		gw, err := m.StartGateway(ctx, name, transport, perServer, offset)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
		} else {
			started = append(started, name)
			gateways = append(gateways, gw)
		}
		if opts.Port <= 0 {
			offset++
		} else {
			offset++
		}
	}
	return started, gateways, errs
}

