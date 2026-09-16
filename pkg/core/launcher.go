package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/coditary/wuji-core/pkg/netx"
)

// EnsureDriver makes sure the configured remote driver process is reachable.
func (c *Core) EnsureDriver(ctx context.Context, driverID string) error {
	if driverID == "" || driverID == "dummy" {
		return nil
	}
	if d, err := c.registry.Get(driverID); err == nil && !d.Info().Remote {
		return nil
	}
	endpoint := c.appConfig.DriverGRPCAddr(driverID)
	if endpoint == "" {
		return fmt.Errorf("unknown remote driver %q", driverID)
	}

	bin := c.appConfig.DriverBinaryPath(driverID)

	if netx.EndpointListening(endpoint) {
		if netx.BinaryNewerThanEndpoint(bin, endpoint) {
			_ = c.registry.Unregister(driverID)
			if err := netx.StopListenerOnEndpoint(endpoint); err != nil {
				return fmt.Errorf("restart stale driver %q: %w", driverID, err)
			}
			for i := 0; i < 20 && netx.EndpointListening(endpoint); i++ {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(250 * time.Millisecond):
				}
			}
		} else if _, err := c.registry.Get(driverID); err == nil {
			c.touchDriver(driverID)
			return nil
		} else {
			if err := c.ConnectRemote(ctx, endpoint, false); err != nil {
				return err
			}
			c.touchDriver(driverID)
			return nil
		}
	} else if _, err := c.registry.Get(driverID); err == nil {
		_ = c.registry.Unregister(driverID)
	}

	if !c.appConfig.AutoStartDriverEnabled(driverID) {
		return fmt.Errorf("driver %q not running at %s (start it manually or enable auto_start_driver)", driverID, endpoint)
	}

	if _, err := os.Stat(bin); err != nil {
		return fmt.Errorf("driver binary not found at %s: %w", bin, err)
	}

	logPath := filepath.Join(c.appConfig.Root, ".wuji", "logs", fmt.Sprintf("driver-%s.log", driverID))
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return err
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	// Detach from the RPC context: auto-started drivers must survive individual requests
	// (e.g. taiji auto-continue issues several sequential GenerateText calls).
	cmd := exec.Command(bin, "--addr", endpoint)
	cmd.Dir = c.appConfig.Root
	cmd.Env = driverProcessEnv()
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if err := cmd.Start(); err != nil {
		_ = logFile.Close()
		return fmt.Errorf("start driver %q: %w", driverID, err)
	}
	_ = logFile.Close()

	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		if netx.EndpointListening(endpoint) {
			if err := c.ConnectRemote(ctx, endpoint, false); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "started driver %q at %s\n", driverID, endpoint)
			c.touchDriver(driverID)
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(300 * time.Millisecond):
		}
	}
	return fmt.Errorf("driver %q did not start listening on %s within 60s (see %s)", driverID, endpoint, logPath)
}

func driverProcessEnv() []string {
	env := os.Environ()
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return env
	}
	extra := filepath.Join(home, ".local", "bin")
	for i, e := range env {
		if !strings.HasPrefix(e, "PATH=") {
			continue
		}
		if strings.Contains(e, extra) {
			return env
		}
		env[i] = e + ":" + extra
		return env
	}
	return append(env, "PATH="+extra)
}
