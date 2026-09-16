package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coditary/wuji-core/pkg/config"
)

func httpTimeout(spec *config.APICapabilitySpec) time.Duration {
	if spec != nil && spec.TimeoutSeconds > 0 {
		return time.Duration(spec.TimeoutSeconds) * time.Second
	}
	return 15 * time.Minute
}

func doHTTP(ctx context.Context, spec *config.APICapabilitySpec, vars map[string]string) ([]byte, error) {
	if spec == nil {
		return nil, fmt.Errorf("api spec is nil")
	}
	apiKey, err := resolveAPIKey(ctx, spec.APIKey, spec.APIKeyEnv)
	if err != nil {
		return nil, err
	}
	if apiKey != "" {
		vars = copyVars(vars)
		vars["API_KEY"] = apiKey
	}
	method := strings.TrimSpace(spec.Method)
	if method == "" {
		method = http.MethodPost
	}
	urlStr := expandString(strings.TrimSpace(spec.URL), vars)
	if urlStr == "" {
		urlStr = strings.TrimRight(expandString(strings.TrimSpace(spec.BaseURL), vars), "/")
	}
	if urlStr == "" {
		return nil, fmt.Errorf("url or base_url is required for http protocol")
	}
	urlStr = appendQuery(urlStr, spec.Query, vars)

	var bodyReader io.Reader
	contentType := strings.TrimSpace(spec.ContentType)

	if len(spec.Multipart) > 0 {
		body, ct, err := buildMultipart(spec.Multipart, vars)
		if err != nil {
			return nil, err
		}
		bodyReader = body
		contentType = ct
	} else if strings.TrimSpace(spec.BodyRaw) != "" {
		raw := expandString(spec.BodyRaw, vars)
		bodyReader = strings.NewReader(raw)
		if contentType == "" {
			contentType = "application/json"
		}
	} else if len(spec.Body) > 0 {
		raw, err := marshalBody(expandMap(spec.Body, vars))
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(raw)
		if contentType == "" {
			contentType = "application/json"
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr, bodyReader)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for k, v := range expandHeaders(spec.Headers, vars, apiKey) {
		req.Header.Set(k, v)
	}
	applyAuth(req, spec, apiKey, vars)

	client := &http.Client{Timeout: httpTimeout(spec)}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if spec.Streaming {
		return readSSE(resp.Body)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, httpAPIError(spec, resp.StatusCode, raw)
	}

	if spec.Poll != nil {
		return pollUntilDone(ctx, spec, client, raw, vars)
	}
	return raw, nil
}

func httpAPIError(spec *config.APICapabilitySpec, status int, raw []byte) error {
	msg := strings.TrimSpace(string(raw))
	if spec != nil && strings.TrimSpace(spec.ErrorPath) != "" {
		if extracted, err := extractJSONPath(raw, spec.ErrorPath); err == nil && strings.TrimSpace(extracted) != "" {
			msg = extracted
		}
	}
	return fmt.Errorf("http api status %d: %s", status, msg)
}

func buildMultipart(fields []config.APIMultipartField, vars map[string]string) (io.Reader, string, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	for _, f := range fields {
		name := strings.TrimSpace(f.Name)
		if name == "" {
			continue
		}
		if strings.TrimSpace(f.File) != "" {
			path := expandString(f.File, vars)
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, "", fmt.Errorf("multipart file %q: %w", path, err)
			}
			part, err := w.CreateFormFile(name, filepathBase(path))
			if err != nil {
				return nil, "", err
			}
			if _, err := part.Write(data); err != nil {
				return nil, "", err
			}
			continue
		}
		if err := w.WriteField(name, expandString(f.Value, vars)); err != nil {
			return nil, "", err
		}
	}
	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return buf, w.FormDataContentType(), nil
}

func filepathBase(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[i+1:]
	}
	return path
}
