package local_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/driver/local"
)

func TestNormalizeDriverID(t *testing.T) {
	cases := map[string]string{
		"local-rag": "local-rag",
		"local":     "local-rag",
		"rag":       "local-rag",
		"vllm":      "vllm",
		"":          "",
	}
	for in, want := range cases {
		if got := local.NormalizeDriverID(in); got != want {
			t.Fatalf("NormalizeDriverID(%q)=%q want %q", in, got, want)
		}
	}
}
