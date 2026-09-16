package driver_test

import (
	"testing"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/modelformat"
)

func TestCapabilityFormatsRoundTrip(t *testing.T) {
	original := driver.CapabilityFormats{
		capability.TextGeneration: {modelformat.GGUF, modelformat.Ollama},
	}

	proto := driver.CapabilityFormatsToProto(original)
	restored := driver.CapabilityFormatsFromProto(proto)

	if len(restored) != len(original) {
		t.Fatalf("expected %d entries, got %d", len(original), len(restored))
	}
	got := restored.FormatsFor(capability.TextGeneration)
	if len(got) != 2 || got[0] != modelformat.GGUF || got[1] != modelformat.Ollama {
		t.Fatalf("unexpected formats: %#v", got)
	}
}

func TestCapabilityFormatsFromProtoSkipsEmpty(t *testing.T) {
	restored := driver.CapabilityFormatsFromProto([]*wujiv1.CapabilityFormatSupport{
		{Capability: "text", Formats: []string{"gguf"}},
		{Capability: "image"},
	})
	if len(restored) != 1 {
		t.Fatalf("expected one entry, got %#v", restored)
	}
}
