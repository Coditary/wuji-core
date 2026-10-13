package driver

import "testing"

func TestApplyAudioBrandAliasCoquiEngine(t *testing.T) {
	req := AudioRequest{}
	id, ok := ApplyAudioBrandAlias("coqui", &req, func(runtime string) bool {
		return runtime == AudioRuntimePyTorch
	})
	if !ok || id != AudioRuntimePyTorch {
		t.Fatalf("alias: ok=%v id=%q", ok, id)
	}
	if req.Engine != "coqui" {
		t.Fatalf("engine %q", req.Engine)
	}
}

func TestApplyAudioBrandAliasFallbackWhenRuntimeMissing(t *testing.T) {
	req := AudioRequest{}
	id, ok := ApplyAudioBrandAlias("piper", &req, func(string) bool { return false })
	if ok || id != "piper" {
		t.Fatalf("expected no alias, got ok=%v id=%q", ok, id)
	}
}
