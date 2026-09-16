package metrics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coditary/wuji-core/pkg/config"
)

const (
	dirName     = "metrics"
	logFileName = "requests.jsonl"
)

// TextRequestRecord is one logged text generation call.
type TextRequestRecord struct {
	Timestamp        time.Time `json:"timestamp"`
	Driver           string    `json:"driver"`
	Model            string    `json:"model"`
	InputTokens      int       `json:"input_tokens"`
	OutputTokens     int       `json:"output_tokens"`
	LatencyMs        int64     `json:"latency_ms"`
	TokensPerSec     float64   `json:"tokens_per_sec"`
	EstimatedCostUSD float64   `json:"estimated_cost_usd"`
	Streamed         bool      `json:"streamed"`
}

// Summary aggregates metrics log entries.
type Summary struct {
	Requests         int     `json:"requests"`
	InputTokens      int     `json:"input_tokens"`
	OutputTokens     int     `json:"output_tokens"`
	TotalLatencyMs   int64   `json:"total_latency_ms"`
	AvgTokensPerSec  float64 `json:"avg_tokens_per_sec"`
	EstimatedCostUSD float64 `json:"estimated_cost_usd"`
	StreamedRequests int     `json:"streamed_requests"`
}

// Log appends one record to .wuji/metrics/requests.jsonl.
func Log(root string, rec TextRequestRecord) error {
	if root == "" {
		return nil
	}
	dir := filepath.Join(root, ".wuji", dirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	raw, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(dir, logFileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(raw, '\n'))
	return err
}

// EstimateTokens approximates token count from text length (~4 chars per token).
func EstimateTokens(text string) int {
	n := len([]rune(text))
	if n == 0 {
		return 0
	}
	return (n + 3) / 4
}

// ComputeCostUSD estimates cost from pricing config and token counts.
func ComputeCostUSD(pricing map[string]config.ModelPricing, model string, inputTokens, outputTokens int) float64 {
	cfg := &config.Config{Pricing: pricing}
	p := cfg.PricingFor(model)
	inCost := float64(inputTokens) / 1_000_000 * p.InputPerMillion
	outCost := float64(outputTokens) / 1_000_000 * p.OutputPerMillion
	return inCost + outCost
}

// ReadSince loads records newer than cutoff (zero = all).
func ReadSince(root string, since time.Time) ([]TextRequestRecord, error) {
	path := filepath.Join(root, ".wuji", dirName, logFileName)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var out []TextRequestRecord
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var rec TextRequestRecord
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			continue
		}
		if !since.IsZero() && rec.Timestamp.Before(since) {
			continue
		}
		out = append(out, rec)
	}
	return out, scanner.Err()
}

// Summarize aggregates records into a summary.
func Summarize(records []TextRequestRecord) Summary {
	var s Summary
	var tpsSum float64
	var tpsCount int
	for _, r := range records {
		s.Requests++
		s.InputTokens += r.InputTokens
		s.OutputTokens += r.OutputTokens
		s.TotalLatencyMs += r.LatencyMs
		s.EstimatedCostUSD += r.EstimatedCostUSD
		if r.Streamed {
			s.StreamedRequests++
		}
		if r.TokensPerSec > 0 {
			tpsSum += r.TokensPerSec
			tpsCount++
		}
	}
	if tpsCount > 0 {
		s.AvgTokensPerSec = tpsSum / float64(tpsCount)
	}
	return s
}

// FormatSummary prints a human-readable stats summary.
func FormatSummary(s Summary, window string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Wuji text generation stats (%s)\n", window)
	fmt.Fprintf(&b, "  requests:          %d\n", s.Requests)
	fmt.Fprintf(&b, "  input tokens:      %d\n", s.InputTokens)
	fmt.Fprintf(&b, "  output tokens:     %d\n", s.OutputTokens)
	fmt.Fprintf(&b, "  total latency:     %s\n", formatDuration(time.Duration(s.TotalLatencyMs)*time.Millisecond))
	if s.AvgTokensPerSec > 0 {
		fmt.Fprintf(&b, "  avg output tok/s:  %.1f\n", s.AvgTokensPerSec)
	}
	fmt.Fprintf(&b, "  estimated cost:    $%.4f\n", s.EstimatedCostUSD)
	if s.StreamedRequests > 0 {
		fmt.Fprintf(&b, "  streamed requests: %d\n", s.StreamedRequests)
	}
	return b.String()
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return d.Round(time.Millisecond).String()
}
