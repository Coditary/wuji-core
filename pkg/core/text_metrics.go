package core

import (
	"strings"
	"time"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/metrics"
)

func (c *Core) enrichAndLogTextMetrics(driverID, model string, req driver.TextRequest, resp *driver.TextResponse, start time.Time, streamed bool) {
	if resp == nil {
		return
	}
	latency := time.Since(start)
	resp.LatencyMs = latency.Milliseconds()
	if resp.InputTokens == 0 {
		resp.InputTokens = estimateInputTokens(req)
	}
	if resp.TokensUsed == 0 && strings.TrimSpace(resp.Text) != "" {
		resp.TokensUsed = metrics.EstimateTokens(resp.Text)
	}
	if latency > 0 && resp.TokensUsed > 0 {
		resp.TokensPerSec = float64(resp.TokensUsed) / latency.Seconds()
	}
	if model == "" {
		model = c.defaultModelFor(driverID)
	}
	var pricing map[string]config.ModelPricing
	if c.appConfig != nil {
		pricing = c.appConfig.Pricing
	}
	resp.EstimatedCostUSD = metrics.ComputeCostUSD(pricing, model, resp.InputTokens, resp.TokensUsed)
	_ = metrics.Log(c.metricsRoot(), metrics.TextRequestRecord{
		Timestamp:        time.Now().UTC(),
		Driver:           driverID,
		Model:            model,
		InputTokens:      resp.InputTokens,
		OutputTokens:     resp.TokensUsed,
		LatencyMs:        resp.LatencyMs,
		TokensPerSec:     resp.TokensPerSec,
		EstimatedCostUSD: resp.EstimatedCostUSD,
		Streamed:         streamed,
	})
}

func (c *Core) metricsRoot() string {
	if c.appConfig != nil && c.appConfig.Root != "" {
		return c.appConfig.Root
	}
	return ""
}

func estimateInputTokens(req driver.TextRequest) int {
	if strings.TrimSpace(req.Prompt) != "" {
		return metrics.EstimateTokens(req.Prompt)
	}
	var n int
	for _, m := range req.Messages {
		n += metrics.EstimateTokens(m.Content)
	}
	if strings.TrimSpace(req.SystemPrompt) != "" {
		n += metrics.EstimateTokens(req.SystemPrompt)
	}
	return n
}
