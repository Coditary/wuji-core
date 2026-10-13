package modelformat_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/modelformat"
)

func TestMergeTypesDedupes(t *testing.T) {
	got := modelformat.MergeTypes(
		[]modelformat.Type{modelformat.GGUF, modelformat.LoRA},
		[]modelformat.Type{modelformat.LoRA, modelformat.Ollama},
	)
	want := []modelformat.Type{modelformat.GGUF, modelformat.LoRA, modelformat.Ollama}
	if len(got) != len(want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%q want %q", i, got[i], want[i])
		}
	}
}

func TestDiffusionImageFormatsIncludeAux(t *testing.T) {
	formats := modelformat.DiffusionImageFormats()
	if len(formats) < 10 {
		t.Fatalf("expected diffusion image preset, got %d formats", len(formats))
	}
	has := func(t modelformat.Type) bool {
		for _, f := range formats {
			if f == t {
				return true
			}
		}
		return false
	}
	for _, want := range []modelformat.Type{modelformat.ControlNet, modelformat.VAE, modelformat.LoRA} {
		if !has(want) {
			t.Fatalf("missing %q in DiffusionImageFormats", want)
		}
	}
}
