package api

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

func TestMergeTextDefaults(t *testing.T) {
	spec := &config.APICapabilitySpec{
		Model: "gpt-4o",
		Defaults: &config.APIDefaults{
			MaxTokens:   2048,
			Temperature: 0.5,
			TopP:        0.9,
		},
	}
	out := mergeTextDefaults(spec, driver.TextRequest{})
	if out.Model != "gpt-4o" || out.MaxTokens != 2048 || out.Temperature != 0.5 || out.TopP != 0.9 {
		t.Fatalf("merged %+v", out)
	}
}
