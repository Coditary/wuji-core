package chatformat

import (
	"regexp"
	"strings"
)

var anthropicToolNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,128}$`)

const anthropicToolNameMaxLen = 128

// NameMaps tracks tool-name rewrites for providers with strict name rules.
type NameMaps struct {
	Forward map[string]string // original -> wire
	Reverse map[string]string // wire -> original
}

// ToolsOpenAI encodes tool definitions for OpenAI / Ollama / Cursor-compatible APIs.
func ToolsOpenAI(tools []ToolDef) []map[string]any {
	out := make([]map[string]any, 0, len(tools))
	for _, t := range tools {
		params := t.Parameters
		if params == nil {
			params = map[string]any{"type": "object", "properties": map[string]any{}}
		}
		out = append(out, map[string]any{
			"type": "function",
			"function": map[string]any{
				"name":        t.Name,
				"description": t.Description,
				"parameters":  params,
			},
		})
	}
	return out
}

// ToolsAnthropic encodes tool definitions for Anthropic Messages API.
// Invalid names are sanitized; use Reverse to restore original names in responses.
func ToolsAnthropic(tools []ToolDef) ([]map[string]any, NameMaps) {
	maps := buildAnthropicNameMaps(toolNames(tools))
	out := make([]map[string]any, 0, len(tools))
	for _, t := range tools {
		name := t.Name
		if wire, ok := maps.Forward[name]; ok {
			name = wire
		}
		params := t.Parameters
		if params == nil {
			params = map[string]any{"type": "object", "properties": map[string]any{}}
		}
		out = append(out, map[string]any{
			"name":         name,
			"description":  t.Description,
			"input_schema": params,
		})
	}
	return out, maps
}

func toolNames(tools []ToolDef) []string {
	out := make([]string, 0, len(tools))
	for _, t := range tools {
		if strings.TrimSpace(t.Name) != "" {
			out = append(out, t.Name)
		}
	}
	return out
}

func buildAnthropicNameMaps(originals []string) NameMaps {
	forward := map[string]string{}
	used := map[string]bool{}
	for _, name := range originals {
		candidate := basicSanitizeAnthropicToolName(name)
		if candidate == name {
			used[candidate] = true
		}
	}
	for _, name := range originals {
		candidate := basicSanitizeAnthropicToolName(name)
		if candidate == name {
			continue
		}
		if forward[name] != "" {
			continue
		}
		unique := candidate
		n := 1
		for used[unique] {
			n++
			suffix := "_" + itoa(n)
			head := candidate
			if len(head)+len(suffix) > anthropicToolNameMaxLen {
				head = head[:anthropicToolNameMaxLen-len(suffix)]
			}
			unique = head + suffix
		}
		forward[name] = unique
		used[unique] = true
	}
	reverse := map[string]string{}
	for orig, wire := range forward {
		reverse[wire] = orig
	}
	return NameMaps{Forward: forward, Reverse: reverse}
}

func basicSanitizeAnthropicToolName(name string) string {
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '-':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	out := strings.Trim(b.String(), "_")
	if out == "" {
		out = "tool"
	}
	if len(out) > anthropicToolNameMaxLen {
		out = out[:anthropicToolNameMaxLen]
	}
	return out
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

// RestoreToolName maps a wire name back to the original when sanitization ran.
func RestoreToolName(name string, maps NameMaps) string {
	if maps.Reverse == nil {
		return name
	}
	if orig, ok := maps.Reverse[name]; ok {
		return orig
	}
	return name
}

// IsValidAnthropicToolName reports whether a name can be sent as-is.
func IsValidAnthropicToolName(name string) bool {
	return anthropicToolNamePattern.MatchString(name)
}
