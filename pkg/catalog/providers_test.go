package catalog

import (
	"testing"
	"time"

	wujicfg "github.com/coditary/wuji-core/pkg/config"
)

func TestMergeIntoProviderCatalog(t *testing.T) {
	root := t.TempDir()
	raw := &rawCatalog{
		Providers: map[string]rawProvider{
			"github-copilot": {
				ID: "github-copilot", Name: "GitHub Copilot",
				API: "https://api.githubcopilot.com",
				Models: map[string]rawModel{
					"claude-opus-5": {Name: "Claude Opus 5"},
				},
			},
		},
	}
	idx, meta, err := BuildIndexFromRaw(raw, DefaultSource, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if err := Save(root, idx, meta); err != nil {
		t.Fatal(err)
	}

	summaries := map[string]wujicfg.ProviderSummary{
		"llama": {Driver: "llama", Kind: "driver"},
	}
	MergeIntoProviderCatalog(root, summaries)

	if _, ok := summaries["github-copilot"]; !ok {
		t.Fatal("expected github-copilot in summaries")
	}
	if summaries["github-copilot"].Kind != "catalog" {
		t.Fatalf("kind = %q", summaries["github-copilot"].Kind)
	}
	if len(summaries["github-copilot"].Models) != 1 {
		t.Fatalf("models = %v", summaries["github-copilot"].Models)
	}
}
