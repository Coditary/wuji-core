package catalog

import (
	"testing"
	"time"
)

func TestBuildIndex(t *testing.T) {
	raw := &rawCatalog{
		Models: map[string]rawModel{
			"google/gemini-3.8-flash": {
				ID: "google/gemini-3.8-flash", Name: "Gemini 3.8 Flash",
				Reasoning: true, ToolCall: true, OpenWeights: false,
				Limit: rawModelLimit{Context: 1048576, Output: 65536},
			},
		},
		Providers: map[string]rawProvider{
			"github-copilot": {
				ID: "github-copilot", Name: "GitHub Copilot",
				API: "https://api.githubcopilot.com", NPM: "@ai-sdk/openai-compatible",
				Env: []string{"GITHUB_TOKEN"},
				Models: map[string]rawModel{
					"claude-opus-5": {ID: "claude-opus-5", Name: "Claude Opus 5"},
					"gpt-5.6-sol":   {ID: "gpt-5.6-sol", Name: "GPT-5.6 Sol"},
				},
			},
		},
	}

	idx, meta, err := BuildIndexFromRaw(raw, DefaultSource, time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("BuildIndexFromRaw: %v", err)
	}
	if meta.ProviderCount != 1 || meta.ModelCount != 1 || meta.LabCount != 1 {
		t.Fatalf("unexpected meta: %+v", meta)
	}
	if meta.ProviderModels != 2 {
		t.Fatalf("provider models = %d, want 2", meta.ProviderModels)
	}

	p := idx.Providers["github-copilot"]
	if p.API != "https://api.githubcopilot.com" {
		t.Fatalf("provider API: %q", p.API)
	}
	if len(p.Models) != 2 {
		t.Fatalf("provider model count = %d", len(p.Models))
	}

	m := idx.Models["google/gemini-3.8-flash"]
	if m.Lab != "google" || !m.Reasoning || m.Context != 1048576 {
		t.Fatalf("model entry: %+v", m)
	}

	lab := idx.Labs["google"]
	if lab.Name != "Google" || len(lab.Models) != 1 {
		t.Fatalf("lab entry: %+v", lab)
	}
}

func TestSaveAndLoadIndex(t *testing.T) {
	root := t.TempDir()
	raw := &rawCatalog{
		Models: map[string]rawModel{
			"anthropic/claude-opus-5": {ID: "anthropic/claude-opus-5", Name: "Claude Opus 5"},
		},
		Providers: map[string]rawProvider{
			"anthropic": {ID: "anthropic", Name: "Anthropic", Models: map[string]rawModel{}},
		},
	}
	idx, meta, err := BuildIndexFromRaw(raw, DefaultSource, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if err := Save(root, idx, meta); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := LoadIndex(root)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	if len(loaded.Providers) != 1 || len(loaded.Models) != 1 {
		t.Fatalf("loaded index: providers=%d models=%d", len(loaded.Providers), len(loaded.Models))
	}
	gotMeta, err := LoadMeta(root)
	if err != nil || gotMeta == nil {
		t.Fatalf("LoadMeta: %v meta=%v", err, gotMeta)
	}
}

func TestFormatLabName(t *testing.T) {
	if got := formatLabName("zhipu-ai"); got != "Zhipu Ai" {
		t.Fatalf("formatLabName = %q", got)
	}
}
