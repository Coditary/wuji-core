package api

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
)

func TestHTTPAPIErrorWithPath(t *testing.T) {
	spec := &config.APICapabilitySpec{ErrorPath: "error.message"}
	raw := []byte(`{"error":{"message":"nope"}}`)
	err := httpAPIError(spec, 400, raw)
	if err == nil || err.Error() == "" {
		t.Fatal("expected error")
	}
}

func TestExpandBodyRaw(t *testing.T) {
	spec := &config.APICapabilitySpec{
		BodyRaw: `{"prompt":"${prompt}","model":"${model}"}`,
	}
	vars := map[string]string{"prompt": "hi", "model": "xl"}
	out := expandString(spec.BodyRaw, vars)
	if out != `{"prompt":"hi","model":"xl"}` {
		t.Fatalf("got %s", out)
	}
}
