package netx

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const unixScheme = "unix://"

// IsUnix reports whether endpoint uses a unix domain socket.
func IsUnix(endpoint string) bool {
	ep := strings.TrimSpace(endpoint)
	return strings.HasPrefix(ep, unixScheme) || strings.HasPrefix(ep, "unix:")
}

// UnixSocketPath returns the filesystem path for a unix endpoint.
func UnixSocketPath(endpoint string) string {
	ep := strings.TrimSpace(endpoint)
	if strings.HasPrefix(ep, unixScheme) {
		return ep[len(unixScheme):]
	}
	if strings.HasPrefix(ep, "unix:") {
		return strings.TrimPrefix(ep, "unix:")
	}
	return ""
}

// DialTarget normalizes an endpoint for grpc.NewClient.
func DialTarget(endpoint string) string {
	ep := strings.TrimSpace(endpoint)
	if ep == "" {
		return ep
	}
	if strings.HasPrefix(ep, unixScheme) {
		return ep
	}
	if strings.HasPrefix(ep, "unix:") {
		path := strings.TrimPrefix(ep, "unix:")
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		}
		return unixScheme + path
	}
	return ep
}

// Listen binds a gRPC server endpoint (unix://… or host:port).
// The returned dial target is what clients should pass to grpc.NewClient.
func Listen(endpoint string) (net.Listener, string, error) {
	ep := strings.TrimSpace(endpoint)
	if ep == "" {
		return nil, "", fmt.Errorf("empty endpoint")
	}
	if IsUnix(ep) {
		path := UnixSocketPath(ep)
		if !filepath.IsAbs(path) {
			abs, err := filepath.Abs(path)
			if err != nil {
				return nil, "", fmt.Errorf("resolve socket path: %w", err)
			}
			path = abs
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, "", fmt.Errorf("create socket dir: %w", err)
		}
		_ = os.Remove(path)
		lis, err := net.Listen("unix", path)
		if err != nil {
			return nil, "", fmt.Errorf("listen unix %s: %w", path, err)
		}
		return lis, unixScheme + path, nil
	}

	lis, err := net.Listen("tcp", ep)
	if err != nil {
		return nil, "", fmt.Errorf("listen tcp %s: %w", ep, err)
	}
	return lis, lis.Addr().String(), nil
}

// EndpointListening reports whether the endpoint accepts connections.
func EndpointListening(endpoint string) bool {
	ep := strings.TrimSpace(endpoint)
	if ep == "" {
		return false
	}
	if path := UnixSocketPath(ep); path != "" {
		conn, err := net.DialTimeout("unix", path, 150*time.Millisecond)
		if err != nil {
			return false
		}
		_ = conn.Close()
		return true
	}
	conn, err := net.DialTimeout("tcp", ep, 150*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// StopListenerOnEndpoint terminates a listener on the given endpoint.
func StopListenerOnEndpoint(endpoint string) error {
	path := UnixSocketPath(endpoint)
	if path != "" {
		if fuser, err := exec.LookPath("fuser"); err == nil {
			out, err := exec.Command(fuser, "-k", path).CombinedOutput()
			if err != nil {
				msg := strings.TrimSpace(string(out))
				if msg != "" && !strings.Contains(msg, "No such file") {
					return fmt.Errorf("fuser -k %s: %v (%s)", path, err, msg)
				}
			}
		}
		_ = os.Remove(path)
		return nil
	}

	_, port, err := ParseHostPort(endpoint)
	if err != nil {
		return err
	}
	return stopTCPPort(port)
}

// BinaryNewerThanEndpoint reports whether binPath is newer than the listener.
func BinaryNewerThanEndpoint(binPath, endpoint string) bool {
	fi, err := os.Stat(binPath)
	if err != nil {
		return false
	}
	path := UnixSocketPath(endpoint)
	if path != "" {
		sock, err := os.Stat(path)
		if err != nil {
			return false
		}
		return fi.ModTime().After(sock.ModTime())
	}
	_, port, err := ParseHostPort(endpoint)
	if err != nil || port <= 0 {
		return false
	}
	return binaryNewerThanTCPPort(binPath, port)
}

// ParseHostPort splits a TCP host:port address.
func ParseHostPort(addr string) (host string, port int, err error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return "", 0, fmt.Errorf("empty address")
	}
	if IsUnix(addr) {
		return "", 0, fmt.Errorf("address %q is not tcp", addr)
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

func stopTCPPort(port int) error {
	if port <= 0 {
		return nil
	}
	portSpec := fmt.Sprintf("%d/tcp", port)
	if path, err := exec.LookPath("fuser"); err == nil {
		out, err := exec.Command(path, "-k", portSpec).CombinedOutput()
		if err != nil {
			msg := strings.TrimSpace(string(out))
			if msg != "" && !strings.Contains(msg, "No such file") {
				return fmt.Errorf("fuser -k %s: %v (%s)", portSpec, err, msg)
			}
		}
	}
	return nil
}
