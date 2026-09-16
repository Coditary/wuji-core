package config

import "testing"

func TestResolveAPICapabilitySpecMergesTask(t *testing.T) {
	base := &APICapabilitySpec{
		Protocol:  "http",
		URL:       "https://api.example.com/v1",
		APIKeyEnv: "KEY",
		Body:      map[string]string{"model": "${model}"},
		Tasks: map[string]*APICapabilitySpec{
			"generate": {
				URL:  "https://api.example.com/v1/generate",
				Body: map[string]string{"prompt": "${prompt}"},
			},
		},
	}
	spec := ResolveAPICapabilitySpec(base, "generate")
	if spec.URL != "https://api.example.com/v1/generate" {
		t.Fatalf("url=%s", spec.URL)
	}
	if spec.Body["prompt"] != "${prompt}" {
		t.Fatalf("body=%v", spec.Body)
	}
	if spec.APIKeyEnv != "KEY" {
		t.Fatalf("api_key_env=%s", spec.APIKeyEnv)
	}
}
