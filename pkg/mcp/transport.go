package mcp

import (
	"fmt"
	"strings"

	wujicfg "github.com/coditary/wuji-core/pkg/config"
)

// Transport selects how an MCP server is exposed to clients.
type Transport string

const (
	TransportDaemon Transport = "daemon"
	TransportStdio  Transport = "stdio"
	TransportHTTP   Transport = "http"
	TransportSSE    Transport = "sse"
)

const defaultGatewayHost = "127.0.0.1"
const defaultGatewayPort = 9333

// GatewayOptions configures HTTP/SSE gateway exposure.
type GatewayOptions struct {
	Host     string
	Port     int
	Path     string
	PortBase int
}

// ResolvedGateway holds the effective bind settings for a gateway.
type ResolvedGateway struct {
	Host string
	Port int
	Path string
	URL  string
}

// ResolveGateway fills defaults for a named server.
func ResolveGateway(name string, opts GatewayOptions, portOffset int) ResolvedGateway {
	host := opts.Host
	if host == "" {
		host = defaultGatewayHost
	}

	port := opts.Port
	if port <= 0 {
		base := opts.PortBase
		if base <= 0 {
			base = defaultGatewayPort
		}
		port = base + portOffset
	}

	path := opts.Path
	if path == "" {
		path = "/mcp/" + name
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return ResolvedGateway{
		Host: host,
		Port: port,
		Path: path,
		URL:  fmt.Sprintf("http://%s:%d%s", host, port, path),
	}
}

func shellQuote(s string) string {
	if s == "" {
		return "''"
	}
	if !strings.ContainsAny(s, " \t\n'\"\\$") {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func formatShellCommand(command string, args []string) string {
	parts := []string{shellQuote(command)}
	for _, arg := range args {
		parts = append(parts, shellQuote(arg))
	}
	return strings.Join(parts, " ")
}

func buildStdioShellCommand(def wujicfg.MCPServerDef) string {
	cmd := formatShellCommand(def.Command, def.Args)
	if def.Cwd != "" {
		cmd = fmt.Sprintf("cd %s && %s", shellQuote(def.Cwd), cmd)
	}
	return cmd
}

func gatewayTransportName(t Transport) string {
	switch t {
	case TransportHTTP:
		return "streamableHttp"
	case TransportSSE:
		return "sse"
	default:
		return "streamableHttp"
	}
}
