package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

var knownRemoteDriverIDs = []string{"a1111", "vllm", "llama", "echo", "raggo", "gorag"}

// DriverIDForEndpoint maps a gRPC endpoint back to a driver ID.
func (c *Config) DriverIDForEndpoint(endpoint string) string {
	if id := driverIDFromSocketEndpoint(endpoint); id != "" {
		return id
	}
	for _, id := range knownRemoteDriverIDs {
		if c.DriverGRPCAddr(id) == endpoint {
			return id
		}
	}
	return ""
}

func driverIDFromSocketEndpoint(endpoint string) string {
	const prefix = "unix://"
	if !strings.HasPrefix(endpoint, prefix) {
		return ""
	}
	path := endpoint[len(prefix):]
	base := filepath.Base(path)
	if !strings.HasSuffix(base, ".sock") {
		return ""
	}
	return strings.TrimSuffix(base, ".sock")
}
func (c *Config) AutoStartDriverEnabled(driverID string) bool {
	if c == nil {
		return false
	}
	switch driverID {
	case "dummy", "":
		return false
	case "a1111":
		if c.Drivers.A1111 != nil && c.Drivers.A1111.AutoStartDriver != nil {
			return *c.Drivers.A1111.AutoStartDriver
		}
		return true
	case "vllm":
		if c.Drivers.VLLM != nil && c.Drivers.VLLM.AutoStartDriver != nil {
			return *c.Drivers.VLLM.AutoStartDriver
		}
		return true
	case "llama":
		if c.Drivers.Llama != nil && c.Drivers.Llama.AutoStartDriver != nil {
			return *c.Drivers.Llama.AutoStartDriver
		}
		return true
	case "echo":
		return true
	case "raggo":
		if c.Drivers.Raggo != nil && c.Drivers.Raggo.AutoStartDriver != nil {
			return *c.Drivers.Raggo.AutoStartDriver
		}
		return false
	case "gorag":
		if c.Drivers.Gorag != nil && c.Drivers.Gorag.AutoStartDriver != nil {
			return *c.Drivers.Gorag.AutoStartDriver
		}
		return false
	case "local-rag":
		if c.Drivers.LocalRAG != nil && c.Drivers.LocalRAG.AutoStartDriver != nil {
			return *c.Drivers.LocalRAG.AutoStartDriver
		}
		return false
	default:
		return false
	}
}

// DriverGRPCAddr returns the gRPC endpoint for a driver.
// Uses an explicit grpc_addr from config when set, otherwise a local unix socket.
func (c *Config) DriverGRPCAddr(driverID string) string {
	if c == nil || driverID == "" {
		return ""
	}
	if explicit := c.explicitGRPCAddr(driverID); explicit != "" {
		return explicit
	}
	return DefaultDriverEndpoint(c.Root, driverID)
}

// DriverBinaryPath returns the path to a driver executable.
func (c *Config) DriverBinaryPath(driverID string) string {
	root := "."
	if c != nil && c.Root != "" {
		root = c.Root
	}
	switch driverID {
	case "a1111":
		if c != nil && c.Drivers.A1111 != nil && c.Drivers.A1111.DriverBin != "" {
			return c.Drivers.A1111.DriverBin
		}
	case "vllm":
		if c != nil && c.Drivers.VLLM != nil && c.Drivers.VLLM.DriverBin != "" {
			return c.Drivers.VLLM.DriverBin
		}
	case "llama":
		if c != nil && c.Drivers.Llama != nil && c.Drivers.Llama.DriverBin != "" {
			return c.Drivers.Llama.DriverBin
		}
	case "echo":
		return filepath.Join(root, "bin", "wuji-driver-echo")
	}
	return filepath.Join(root, "bin", fmt.Sprintf("wuji-driver-%s", driverID))
}
func ParseHostPort(addr string) (host string, port int, err error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return "", 0, fmt.Errorf("empty address")
	}
	if !strings.Contains(addr, ":") {
		return "127.0.0.1", 0, fmt.Errorf("address %q missing port", addr)
	}
	i := strings.LastIndex(addr, ":")
	host = addr[:i]
	if host == "" {
		host = "127.0.0.1"
	}
	if _, err := fmt.Sscanf(addr[i+1:], "%d", &port); err != nil {
		return "", 0, fmt.Errorf("parse port from %q: %w", addr, err)
	}
	return host, port, nil
}
