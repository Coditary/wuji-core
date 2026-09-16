package capability_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/capability"
)

func TestParseKnown(t *testing.T) {
	got, err := capability.Parse("text")
	if err != nil || got != capability.TextGeneration {
		t.Fatalf("Parse(text) = %q err=%v", got, err)
	}
}

func TestParseMeshAlias(t *testing.T) {
	got, err := capability.Parse("3d")
	if err != nil || got != capability.Mesh {
		t.Fatalf("Parse(3d) = %q err=%v", got, err)
	}
}

func TestParseUnknown(t *testing.T) {
	_, err := capability.Parse("nope")
	if err == nil {
		t.Fatal("expected error for unknown capability")
	}
}

func TestAllIncludesText(t *testing.T) {
	for _, c := range capability.All() {
		if c == capability.TextGeneration {
			return
		}
	}
	t.Fatal("All() missing text capability")
}
