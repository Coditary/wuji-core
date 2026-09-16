package api

import (
	"net/http"
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
)

func TestApplyAuthQuery(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://api.example.com/v1", nil)
	spec := &config.APICapabilitySpec{
		Auth: &config.APIAuthConfig{Type: "query", QueryParam: "key"},
	}
	applyAuth(req, spec, "secret", nil)
	if req.URL.Query().Get("key") != "secret" {
		t.Fatalf("query auth: %s", req.URL.RawQuery)
	}
}

func TestAppendQuery(t *testing.T) {
	url := appendQuery("https://api.example.com/gen", map[string]string{"model": "${model}"}, map[string]string{"model": "xl"})
	if !contains(url, "model=xl") {
		t.Fatalf("url: %s", url)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
