package textapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

// DriverTextGen generates text through a Wuji driver.
type DriverTextGen interface {
	GenerateText(ctx context.Context, driverID string, req driver.TextRequest) (*driver.TextResponse, error)
}

// Router resolves provider ids from config and completes chat requests.
type Router struct {
	cfg     *config.Config
	drivers DriverTextGen
}

// NewRouter creates a provider router.
func NewRouter(cfg *config.Config, drivers DriverTextGen) *Router {
	return &Router{cfg: cfg, drivers: drivers}
}

// ChatComplete routes a chat request through the configured provider id.
func (r *Router) ChatComplete(ctx context.Context, providerID string, req ChatRequest) (*ChatResponse, error) {
	if r.cfg == nil {
		return nil, fmt.Errorf("config not loaded")
	}
	id := strings.TrimSpace(providerID)
	if id == "" {
		id = strings.TrimSpace(r.cfg.DefaultProvider)
	}
	if id == "" {
		return nil, fmt.Errorf("provider_id is required")
	}
	pc, ok := r.cfg.Providers[id]
	if !ok {
		return nil, fmt.Errorf("provider %q not found in .wuji/config.yaml", id)
	}
	typ := strings.ToLower(strings.TrimSpace(pc.Type))
	switch typ {
	case "driver":
		return r.completeDriver(ctx, pc, req)
	case "openai":
		return completeOpenAI(ctx, pc, req)
	case "anthropic":
		return completeAnthropic(ctx, pc, req)
	default:
		return nil, fmt.Errorf("provider %q: unknown type %q", id, pc.Type)
	}
}

func (r *Router) completeDriver(ctx context.Context, pc config.ProviderConfig, req ChatRequest) (*ChatResponse, error) {
	driverID := strings.TrimSpace(pc.Driver)
	if driverID == "" {
		return nil, fmt.Errorf("driver provider: driver is required")
	}
	if r.drivers == nil {
		return nil, fmt.Errorf("driver provider: text generator not configured")
	}
	textReq := toTextRequest(req, pc.Model)
	resp, err := r.drivers.GenerateText(ctx, driverID, textReq)
	if err != nil {
		return nil, err
	}
	return &ChatResponse{
		Content:      resp.Text,
		FinishReason: resp.FinishReason,
		Usage:        Usage{OutputTokens: resp.TokensUsed},
	}, nil
}
