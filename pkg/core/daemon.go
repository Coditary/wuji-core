package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/netx"
)

// ConnectRuntime connects to the core daemon or starts it when needed.
func ConnectRuntime(ctx context.Context, cfg *config.Config) (Runtime, func(), error) {
	if os.Getenv("WUJI_EMBEDDED") == "1" {
		c, err := New(Config{AppConfig: cfg, Lazy: true})
		if err != nil {
			return nil, nil, err
		}
		return c, func() { _ = c.Close() }, nil
	}

	root := cfg.Root
	if root == "" {
		root = "."
	}
	endpoint := config.DefaultCoreEndpoint(root)
	if client, err := DialClient(ctx, endpoint, root); err == nil {
		return client, func() { _ = client.Close() }, nil
	}

	if !netx.IsUnix(endpoint) {
		return nil, nil, fmt.Errorf("core daemon not reachable at %s", endpoint)
	}

	if err := startDaemon(ctx, root, endpoint); err != nil {
		return nil, nil, err
	}

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if client, err := DialClient(ctx, endpoint, root); err == nil {
			return client, func() { _ = client.Close() }, nil
		}
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
	return nil, nil, fmt.Errorf("core daemon did not become ready at %s", endpoint)
}

func startDaemon(ctx context.Context, root, endpoint string) error {
	if netx.EndpointListening(endpoint) {
		return nil
	}
	bin := config.CoreBinaryPath(root)
	if bin == "" {
		return fmt.Errorf("wuji-core binary not found; build with: make build-core")
	}
	logPath := filepath.Join(root, ".wuji", "logs", "core.log")
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return err
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, bin, "--addr", endpoint)
	cmd.Dir = root
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if err := cmd.Start(); err != nil {
		_ = logFile.Close()
		return fmt.Errorf("start core daemon: %w", err)
	}
	_ = logFile.Close()
	return nil
}
