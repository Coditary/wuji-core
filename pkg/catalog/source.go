package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type rawCatalog struct {
	Models    map[string]rawModel    `json:"models"`
	Providers map[string]rawProvider `json:"providers"`
}

type rawProvider struct {
	ID     string               `json:"id"`
	Name   string               `json:"name"`
	API    string               `json:"api"`
	NPM    string               `json:"npm"`
	Doc    string               `json:"doc"`
	Env    []string             `json:"env"`
	Models map[string]rawModel  `json:"models"`
}

type rawModel struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Family      string         `json:"family"`
	Reasoning   bool           `json:"reasoning"`
	ToolCall    bool           `json:"tool_call"`
	OpenWeights bool           `json:"open_weights"`
	Limit       rawModelLimit  `json:"limit"`
}

type rawModelLimit struct {
	Context int `json:"context"`
	Output  int `json:"output"`
}

// Fetch downloads catalog.json from source (default https://models.dev).
func Fetch(ctx context.Context, source string) ([]byte, string, error) {
	source = strings.TrimRight(strings.TrimSpace(source), "/")
	if source == "" {
		source = DefaultSource
	}
	url := source + "/catalog.json"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", "wuji-catalog/1.0")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 2 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, "", fmt.Errorf("fetch %s: HTTP %d: %s", url, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read %s: %w", url, err)
	}
	return data, resp.Header.Get("ETag"), nil
}

func decodeRawCatalog(data []byte) (*rawCatalog, error) {
	var raw rawCatalog
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode catalog.json: %w", err)
	}
	if raw.Providers == nil {
		raw.Providers = map[string]rawProvider{}
	}
	if raw.Models == nil {
		raw.Models = map[string]rawModel{}
	}
	return &raw, nil
}
