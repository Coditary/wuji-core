package textapi

import (
	"context"
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

type stubGen struct {
	lastDriver string
	lastReq    driver.TextRequest
}

func (s *stubGen) GenerateText(_ context.Context, driverID string, req driver.TextRequest) (*driver.TextResponse, error) {
	s.lastDriver = driverID
	s.lastReq = req
	return &driver.TextResponse{Text: "ok", FinishReason: "stop"}, nil
}

func TestRouterDriverProvider(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: "echo",
		Providers: map[string]config.ProviderConfig{
			"echo": {Type: "driver", Driver: "echo", Model: "m1"},
		},
	}
	gen := &stubGen{}
	r := NewRouter(cfg, gen)
	resp, err := r.ChatComplete(context.Background(), "echo", ChatRequest{
		Messages: []Message{{Role: RoleUser, Content: "hi"}},
		Model:    "override",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content != "ok" || gen.lastDriver != "echo" {
		t.Fatalf("unexpected: %+v driver=%s", resp, gen.lastDriver)
	}
	if gen.lastReq.Model != "override" {
		t.Fatalf("model: %q", gen.lastReq.Model)
	}
}
