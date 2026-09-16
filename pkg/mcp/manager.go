package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	wujicfg "github.com/coditary/wuji-core/pkg/config"
)

const runtimeDirName = "mcp"

// RuntimeState tracks a running stdio MCP server process.
type RuntimeState struct {
	PID         int       `json:"pid"`
	StartedAt   time.Time `json:"started_at"`
	Command     string    `json:"command"`
	LogFile     string    `json:"log_file,omitempty"`
	Transport   string    `json:"transport,omitempty"`
	GatewayURL  string    `json:"gateway_url,omitempty"`
	GatewayHost string    `json:"gateway_host,omitempty"`
	GatewayPort int       `json:"gateway_port,omitempty"`
	GatewayPath string    `json:"gateway_path,omitempty"`
}

// ServerStatus is the combined config + runtime view of an MCP server.
type ServerStatus struct {
	Name       string
	Def        wujicfg.MCPServerDef
	Running    bool
	PID        int
	Started    time.Time
	LogFile    string
	Transport  string
	GatewayURL string
	RemoteOK   bool // for URL-based servers: last health check result
}

// Manager handles MCP server process lifecycle.
type Manager struct {
	root string
}

// NewManager creates an MCP manager for the given project root.
func NewManager(root string) *Manager {
	return &Manager{root: root}
}

func (m *Manager) runtimeDir() string {
	return filepath.Join(m.root, ".wuji", "runtime", runtimeDirName)
}

func (m *Manager) statePath(name string) string {
	return filepath.Join(m.runtimeDir(), name+".json")
}

func (m *Manager) logPath(name string) string {
	return filepath.Join(m.runtimeDir(), name+".log")
}

// ListStatus returns status for all configured servers.
func (m *Manager) ListStatus(ctx context.Context) ([]ServerStatus, error) {
	cfg, err := wujicfg.LoadMCP(m.root)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(cfg.Servers))
	for name := range cfg.Servers {
		names = append(names, name)
	}
	sort.Strings(names)

	out := make([]ServerStatus, 0, len(names))
	for _, name := range names {
		st, err := m.Status(ctx, name)
		if err != nil {
			return nil, err
		}
		out = append(out, st)
	}
	return out, nil
}

// Status returns the current status of a single server.
func (m *Manager) Status(ctx context.Context, name string) (ServerStatus, error) {
	def, err := wujicfg.GetMCPServer(m.root, name)
	if err != nil {
		return ServerStatus{}, err
	}

	st := ServerStatus{Name: name, Def: def}

	if def.IsRemote() {
		st.RemoteOK = m.checkRemote(ctx, def.URL)
		st.Running = st.RemoteOK
		return st, nil
	}

	state, err := m.loadState(name)
	if err != nil {
		return st, nil
	}
	if state != nil && processAlive(state.PID) {
		st.Running = true
		st.PID = state.PID
		st.Started = state.StartedAt
		st.LogFile = state.LogFile
		st.Transport = state.Transport
		st.GatewayURL = state.GatewayURL
	} else if state != nil {
		_ = m.removeState(name)
	}
	return st, nil
}

// Start launches a stdio MCP server or verifies a remote endpoint.
func (m *Manager) Start(ctx context.Context, name string) error {
	def, err := wujicfg.GetMCPServer(m.root, name)
	if err != nil {
		return err
	}
	if def.Disabled {
		return fmt.Errorf("MCP server %q is disabled", name)
	}

	if def.IsRemote() {
		if !m.checkRemote(ctx, def.URL) {
			return fmt.Errorf("remote MCP server %q at %s is not reachable", name, def.URL)
		}
		return nil
	}

	st, err := m.Status(ctx, name)
	if err != nil {
		return err
	}
	if st.Running {
		return nil
	}

	if err := os.MkdirAll(m.runtimeDir(), 0o755); err != nil {
		return fmt.Errorf("create runtime dir: %w", err)
	}

	logFile := m.logPath(name)
	log, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}

	stdinR, stdinW, err := os.Pipe()
	if err != nil {
		_ = log.Close()
		return fmt.Errorf("stdin pipe: %w", err)
	}

	cmd := exec.CommandContext(ctx, def.Command, def.Args...)
	cmd.Stdin = stdinR
	cmd.Stdout = log
	cmd.Stderr = log
	if def.Cwd != "" {
		cmd.Dir = def.Cwd
	}
	cmd.Env = buildEnv(def.Env)

	if err := cmd.Start(); err != nil {
		_ = stdinR.Close()
		_ = stdinW.Close()
		_ = log.Close()
		return fmt.Errorf("start %s: %w", name, err)
	}
	_ = stdinR.Close()

	state := RuntimeState{
		PID:       cmd.Process.Pid,
		StartedAt: time.Now().UTC(),
		Command:   formatCommand(def),
		LogFile:   logFile,
		Transport: string(TransportDaemon),
	}
	if err := m.saveState(name, state); err != nil {
		_ = cmd.Process.Kill()
		_ = log.Close()
		return err
	}

	// Reap child in background; clear state when process exits.
	go func() {
		_ = cmd.Wait()
		_ = stdinW.Close()
		_ = log.Close()
		if s, _ := m.loadState(name); s != nil && s.PID == state.PID {
			_ = m.removeState(name)
		}
	}()

	// Brief pause to catch immediate crashes.
	time.Sleep(200 * time.Millisecond)
	if !processAlive(state.PID) {
		_ = m.removeState(name)
		return fmt.Errorf("MCP server %q exited immediately — check %s", name, logFile)
	}

	return nil
}

// Stop terminates a running stdio MCP server.
func (m *Manager) Stop(name string) error {
	def, err := wujicfg.GetMCPServer(m.root, name)
	if err != nil {
		return err
	}
	if def.IsRemote() {
		return fmt.Errorf("cannot stop remote MCP server %q (url: %s)", name, def.URL)
	}

	state, err := m.loadState(name)
	if err != nil {
		return err
	}
	if state == nil || !processAlive(state.PID) {
		_ = m.removeState(name)
		return nil
	}

	proc, err := os.FindProcess(state.PID)
	if err != nil {
		_ = m.removeState(name)
		return nil
	}

	_ = proc.Signal(syscall.SIGTERM)
	deadline := time.After(10 * time.Second)
	for processAlive(state.PID) {
		select {
		case <-deadline:
			_ = proc.Kill()
			_ = m.removeState(name)
			return nil
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}

	_ = m.removeState(name)
	return nil
}

// Restart stops and starts a stdio MCP server.
func (m *Manager) Restart(ctx context.Context, name string) error {
	if err := m.Stop(name); err != nil {
		return err
	}
	return m.Start(ctx, name)
}

// StartAll starts all enabled stdio servers and checks remote endpoints.
func (m *Manager) StartAll(ctx context.Context) ([]string, []error) {
	cfg, err := wujicfg.LoadMCP(m.root)
	if err != nil {
		return nil, []error{err}
	}

	var started []string
	var errs []error
	for name, def := range cfg.Servers {
		if def.Disabled {
			continue
		}
		if err := m.Start(ctx, name); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
		} else {
			started = append(started, name)
		}
	}
	return started, errs
}

// StopAll stops all running stdio MCP servers.
func (m *Manager) StopAll() ([]string, []error) {
	cfg, err := wujicfg.LoadMCP(m.root)
	if err != nil {
		return nil, []error{err}
	}

	var stopped []string
	var errs []error
	for name, def := range cfg.Servers {
		if def.IsRemote() {
			continue
		}
		if err := m.Stop(name); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
		} else {
			stopped = append(stopped, name)
		}
	}
	return stopped, errs
}

func (m *Manager) loadState(name string) (*RuntimeState, error) {
	path := m.statePath(name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var state RuntimeState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (m *Manager) saveState(name string, state RuntimeState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.statePath(name), data, 0o644)
}

func (m *Manager) removeState(name string) error {
	err := os.Remove(m.statePath(name))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func processAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// Signal 0 checks existence without sending a signal.
	return proc.Signal(syscall.Signal(0)) == nil
}

func buildEnv(extra map[string]string) []string {
	env := os.Environ()
	for k, v := range extra {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env
}

func formatCommand(def wujicfg.MCPServerDef) string {
	parts := []string{def.Command}
	parts = append(parts, def.Args...)
	return strings.Join(parts, " ")
}

func (m *Manager) checkRemote(ctx context.Context, url string) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	// Any HTTP response means something is listening.
	return true
}
