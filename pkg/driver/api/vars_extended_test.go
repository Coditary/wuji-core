package api

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

func TestRAGHTTPVarsIncludesFlags(t *testing.T) {
	vars := ragHTTPVars(driver.RAGRequest{
		Task: driver.RAGTaskQuery, Collection: "docs", Query: "hello",
		TopK: 5, ChunkOverlap: 50, DryRun: true, Temperature: 0.2,
		Filter: map[string]string{"lang": "de"},
	}, &config.APICapabilitySpec{})
	if vars["query"] != "hello" || vars["dry_run"] != "true" || vars["filter_json"] == "" {
		t.Fatalf("vars=%v", vars)
	}
}

func TestTextTrainHTTPVars(t *testing.T) {
	vars := textTrainHTTPVars(driver.TextTrainRequest{
		Name: "lora-ft", DatasetID: "ds-1", BaseModel: "llama", Epochs: 3, LoRARank: 16,
	}, &config.APICapabilitySpec{})
	if vars["dataset_id"] != "ds-1" || vars["lora_rank"] != "16" {
		t.Fatalf("vars=%v", vars)
	}
}

func TestAugmentImageScaleCapabilities(t *testing.T) {
	caps := augmentCapabilities(config.APIEntry{
		Image: &config.APICapabilitySpec{Protocol: "http"},
	}, []capability.Type{capability.ImageGeneration})
	found := false
	for _, c := range caps {
		if c == capability.ImageUpscale {
			found = true
		}
	}
	if !found {
		t.Fatalf("caps=%v", caps)
	}
}
