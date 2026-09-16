package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

func mediaResponseFromHTTP(spec *config.APICapabilitySpec, raw []byte, kind string) (*driver.ImageResponse, error) {
	path, err := mediaPathFromHTTP(spec, raw, kind)
	if err != nil {
		return nil, err
	}
	return &driver.ImageResponse{Path: path, Format: filepath.Ext(path)}, nil
}

func mediaPathFromHTTP(spec *config.APICapabilitySpec, raw []byte, kind string) (string, error) {
	if spec.Response == nil {
		return "", fmt.Errorf("%s http response mapping is required", kind)
	}
	urlKey := kind + "_url"
	b64Key := kind + "_base64"
	if p := spec.Response[urlKey]; p != "" {
		url, err := extractJSONPath(raw, p)
		if err != nil {
			return "", err
		}
		return downloadToTemp(url, kind)
	}
	if p := spec.Response[b64Key]; p != "" {
		b64, err := extractJSONPath(raw, p)
		if err != nil {
			return "", err
		}
		return writeBase64File(b64, kind)
	}
	if p := spec.Response["url"]; p != "" {
		url, err := extractJSONPath(raw, p)
		if err != nil {
			return "", err
		}
		return downloadToTemp(url, kind)
	}
	if p := spec.Response["path"]; p != "" {
		return extractJSONPath(raw, p)
	}
	return "", fmt.Errorf("%s response needs %s_url, %s_base64, url, or path", kind, kind, kind)
}

func generateImageOpenAI(ctx context.Context, spec *config.APICapabilitySpec, req driver.ImageRequest) (*driver.ImageResponse, error) {
	apiKey, err := resolveAPIKey(ctx, spec.APIKey, spec.APIKeyEnv)
	if err != nil {
		return nil, err
	}
	baseURL := strings.TrimRight(strings.TrimSpace(spec.BaseURL), "/")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	model := firstNonEmpty(req.Model, specModel(spec))
	body := map[string]any{
		"model":  model,
		"prompt": req.Prompt,
		"size":   fmt.Sprintf("%dx%d", defaultInt(req.Width, 1024), defaultInt(req.Height, 1024)),
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/images/generations", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	applyAuth(httpReq, spec, apiKey, nil)
	resp, err := (&http.Client{Timeout: httpTimeout(spec)}).Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, httpAPIError(spec, resp.StatusCode, respBody)
	}
	specCopy := *spec
	if specCopy.Response == nil {
		specCopy.Response = map[string]string{
			"image_url":    "data[0].url",
			"image_base64": "data[0].b64_json",
		}
	}
	return mediaResponseFromHTTP(&specCopy, respBody, "image")
}

func downloadToTemp(url, prefix string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(url)
	if ext == "" || len(ext) > 8 {
		ext = ".bin"
	}
	f, err := os.CreateTemp("", "wuji-"+prefix+"-*"+ext)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return "", err
	}
	return f.Name(), nil
}

func writeBase64File(b64, prefix string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	f, err := os.CreateTemp("", "wuji-"+prefix+"-*.bin")
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return "", err
	}
	return f.Name(), nil
}

func defaultInt(v, def int) int {
	if v > 0 {
		return v
	}
	return def
}
