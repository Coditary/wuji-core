package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// extractJSONPath reads a simple JSON path like $.data[0].url or data.0.text
func extractJSONPath(raw []byte, path string) (string, error) {
	path = strings.TrimSpace(path)
	path = strings.TrimPrefix(path, "$.")
	path = strings.TrimPrefix(path, "$")
	path = strings.Trim(path, ".")
	if path == "" {
		return "", fmt.Errorf("empty json path")
	}
	var doc any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return "", err
	}
	cur := doc
	for _, part := range strings.Split(path, ".") {
		if part == "" {
			continue
		}
		if idx, err := strconv.Atoi(part); err == nil {
			arr, ok := cur.([]any)
			if !ok || idx < 0 || idx >= len(arr) {
				return "", fmt.Errorf("path segment %q: not an array index", part)
			}
			cur = arr[idx]
			continue
		}
		// support key[index] like data[0]
		if i := strings.Index(part, "["); i > 0 && strings.HasSuffix(part, "]") {
			key := part[:i]
			idxStr := strings.TrimSuffix(part[i+1:], "]")
			idx, err := strconv.Atoi(idxStr)
			if err != nil {
				return "", err
			}
			m, ok := cur.(map[string]any)
			if !ok {
				return "", fmt.Errorf("path segment %q: not an object", key)
			}
			arr, ok := m[key].([]any)
			if !ok || idx < 0 || idx >= len(arr) {
				return "", fmt.Errorf("path segment %q: not an array", part)
			}
			cur = arr[idx]
			continue
		}
		m, ok := cur.(map[string]any)
		if !ok {
			return "", fmt.Errorf("path segment %q: not an object", part)
		}
		cur, ok = m[part]
		if !ok {
			return "", fmt.Errorf("path segment %q: not found", part)
		}
	}
	switch v := cur.(type) {
	case string:
		return v, nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
}
