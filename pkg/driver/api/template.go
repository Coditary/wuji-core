package api

import (
	"context"
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/coditary/wuji-core/pkg/auth"
)

var templateVar = regexp.MustCompile(`\$\{([^}]+)\}`)

// expandString replaces ${key} and ${env:VAR} in s using vars.
func expandString(s string, vars map[string]string) string {
	return templateVar.ReplaceAllStringFunc(s, func(m string) string {
		key := strings.TrimSpace(m[2 : len(m)-1])
		if strings.HasPrefix(key, "env:") {
			return os.Getenv(strings.TrimPrefix(key, "env:"))
		}
		if v, ok := vars[key]; ok {
			return v
		}
		return m
	})
}

// expandMap expands all string values in a map.
func expandMap(m map[string]string, vars map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = expandString(v, vars)
	}
	return out
}

// expandHeaders expands header values and injects API_KEY when set.
func expandHeaders(headers map[string]string, vars map[string]string, apiKey string) map[string]string {
	if apiKey != "" {
		vars = copyVars(vars)
		vars["API_KEY"] = apiKey
	}
	return expandMap(headers, vars)
}

func copyVars(in map[string]string) map[string]string {
	out := make(map[string]string, len(in)+1)
	for k, v := range in {
		out[k] = v
	}
	return out
}

// resolveAPIKey returns inline api key, auth store, or environment variable.
func resolveAPIKey(ctx context.Context, inline, envName string) (string, error) {
	return auth.ResolveKey(providerIDFromContext(ctx), inline, envName)
}

// messagesJSON encodes chat messages for HTTP body templates.
func messagesJSON(messages []byte) string {
	if len(messages) == 0 {
		return "[]"
	}
	return string(messages)
}

func marshalBody(body map[string]string) ([]byte, error) {
	typed := map[string]any{}
	for k, v := range body {
		if strings.HasPrefix(strings.TrimSpace(v), "[") || strings.HasPrefix(strings.TrimSpace(v), "{") {
			var parsed any
			if err := json.Unmarshal([]byte(v), &parsed); err == nil {
				typed[k] = parsed
				continue
			}
		}
		typed[k] = v
	}
	return json.Marshal(typed)
}
