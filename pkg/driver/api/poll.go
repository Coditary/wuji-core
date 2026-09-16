package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/coditary/wuji-core/pkg/config"
)

func pollUntilDone(ctx context.Context, spec *config.APICapabilitySpec, client *http.Client, initial []byte, vars map[string]string) ([]byte, error) {
	p := spec.Poll
	if p == nil {
		return initial, nil
	}
	interval := time.Duration(defaultInt(p.IntervalSeconds, 2)) * time.Second
	timeout := httpTimeout(spec)
	if p.TimeoutSeconds > 0 {
		timeout = time.Duration(p.TimeoutSeconds) * time.Second
	}
	deadline := time.Now().Add(timeout)

	if strings.TrimSpace(p.JobIDPath) != "" {
		jobID, err := extractJSONPath(initial, p.JobIDPath)
		if err != nil {
			return nil, fmt.Errorf("poll job_id: %w", err)
		}
		vars = copyVars(vars)
		vars["job_id"] = jobID
	}

	pollURL := expandString(strings.TrimSpace(p.PollURL), vars)
	if pollURL == "" {
		return initial, nil
	}

	for {
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("poll timeout after %s", timeout)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, pollURL, nil)
		if err != nil {
			return nil, err
		}
		apiKey, _ := resolveAPIKey(ctx, spec.APIKey, spec.APIKeyEnv)
		applyAuth(req, spec, apiKey, vars)
		for k, v := range expandHeaders(spec.Headers, vars, apiKey) {
			req.Header.Set(k, v)
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		raw, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 400 {
			return nil, httpAPIError(spec, resp.StatusCode, raw)
		}
		if strings.TrimSpace(p.StatusPath) == "" {
			return raw, nil
		}
		status, err := extractJSONPath(raw, p.StatusPath)
		if err != nil {
			return nil, err
		}
		status = strings.ToLower(strings.TrimSpace(status))
		for _, fail := range p.FailValues {
			if status == strings.ToLower(strings.TrimSpace(fail)) {
				return nil, fmt.Errorf("poll failed with status %q", status)
			}
		}
		done := false
		for _, dv := range p.DoneValues {
			if status == strings.ToLower(strings.TrimSpace(dv)) {
				done = true
				break
			}
		}
		if len(p.DoneValues) == 0 && (status == "succeeded" || status == "completed" || status == "done") {
			done = true
		}
		if done {
			return raw, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(interval):
		}
	}
}
