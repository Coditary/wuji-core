package config

import (
	"os"
	"path/filepath"
	"strings"
)

const driverSocketsDir = "run/drivers"
const coreSocketName = "core.sock"

// CoreSocketPath returns the default unix socket path for the Wuji core daemon.
func CoreSocketPath(root string) string {
	if root == "" {
		root = "."
	}
	return filepath.Join(root, dirName, driverSocketsDir, coreSocketName)
}

// CoreBinaryPath returns the wuji-core daemon executable.
func CoreBinaryPath(root string) string {
	if v := strings.TrimSpace(os.Getenv("WUJI_CORE_BIN")); v != "" {
		return v
	}
	if root == "" {
		root = "."
	}
	return filepath.Join(root, "bin", "wuji-core")
}

// DefaultCoreEndpoint returns the local core daemon endpoint.
func DefaultCoreEndpoint(root string) string {
	if v := strings.TrimSpace(os.Getenv("WUJI_CORE_ADDR")); v != "" {
		return v
	}
	return "unix://" + CoreSocketPath(root)
}

// DriverSocketPath returns the default unix socket path for a driver.
func DriverSocketPath(root, driverID string) string {
	if root == "" {
		root = "."
	}
	return filepath.Join(root, dirName, driverSocketsDir, sanitizeDriverID(driverID)+".sock")
}

// DefaultDriverEndpoint returns the default local gRPC endpoint (unix socket).
func DefaultDriverEndpoint(root, driverID string) string {
	return "unix://" + DriverSocketPath(root, driverID)
}

func sanitizeDriverID(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return "driver"
	}
	var b strings.Builder
	for _, r := range id {
		switch r {
		case '/', '\\', ':', ' ':
			b.WriteByte('_')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ResolveListenAddr returns the CLI --addr flag or the default driver endpoint.
func ResolveListenAddr(driverID, flagOverride string) (string, error) {
	if strings.TrimSpace(flagOverride) != "" {
		return strings.TrimSpace(flagOverride), nil
	}
	cfg, err := Load()
	if err != nil {
		return "", err
	}
	return cfg.DriverGRPCAddr(driverID), nil
}

func (c *Config) explicitGRPCAddr(driverID string) string {
	if c == nil {
		return ""
	}
	switch driverID {
	case "a1111":
		if c.Drivers.A1111 != nil {
			return c.Drivers.A1111.GRPCAddr
		}
	case "vllm":
		if c.Drivers.VLLM != nil {
			return c.Drivers.VLLM.GRPCAddr
		}
	case "llama":
		if c.Drivers.Llama != nil {
			return c.Drivers.Llama.GRPCAddr
		}
	case "raggo":
		if c.Drivers.Raggo != nil {
			return c.Drivers.Raggo.GRPCAddr
		}
	case "gorag":
		if c.Drivers.Gorag != nil {
			return c.Drivers.Gorag.GRPCAddr
		}
	}
	return ""
}
