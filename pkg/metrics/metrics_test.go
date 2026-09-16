package metrics

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/coditary/wuji-core/pkg/config"
)

func TestLogAndSummarize(t *testing.T) {
	root := t.TempDir()
	rec := TextRequestRecord{
		Timestamp:        time.Now().UTC(),
		Driver:           "gemma4",
		Model:            "gemma4:26b",
		InputTokens:      100,
		OutputTokens:     50,
		LatencyMs:        2000,
		TokensPerSec:     25,
		EstimatedCostUSD: 0,
		Streamed:         true,
	}
	if err := Log(root, rec); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, ".wuji", dirName, logFileName)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("log file missing: %v", err)
	}
	records, err := ReadSince(root, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("records = %d, want 1", len(records))
	}
	s := Summarize(records)
	if s.Requests != 1 || s.InputTokens != 100 || s.OutputTokens != 50 {
		t.Fatalf("unexpected summary: %+v", s)
	}
}

func TestComputeCostUSD(t *testing.T) {
	pricing := map[string]config.ModelPricing{
		"gpt-4o": {InputPerMillion: 2.5, OutputPerMillion: 10},
	}
	got := ComputeCostUSD(pricing, "gpt-4o", 1_000_000, 500_000)
	want := 2.5 + 5.0
	if got != want {
		t.Fatalf("cost = %f, want %f", got, want)
	}
}
