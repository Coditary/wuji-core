package modelformat_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/modelformat"
)

func TestAllFormatsAreKnown(t *testing.T) {
	all := modelformat.All()
	if len(all) < 50 {
		t.Fatalf("expected comprehensive format list, got %d entries", len(all))
	}
	for _, f := range all {
		if !modelformat.IsKnown(f) {
			t.Fatalf("format %q has no description", f)
		}
		if f.Description() == "" {
			t.Fatalf("format %q has empty description", f)
		}
	}
}

func TestIsKnownRejectsUnknown(t *testing.T) {
	if modelformat.IsKnown(modelformat.Type("not-a-format")) {
		t.Fatal("expected unknown format to be rejected")
	}
}
