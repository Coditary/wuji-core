package core_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/core"
	"github.com/coditary/wuji-core/pkg/driver"
)

func TestResolveAudioDriverIDSpeechFallsBackToCapability(t *testing.T) {
	cfg := &config.Config{
		DefaultDriver: "dummy",
		CapabilityDrivers: config.CapabilityDrivers{
			"audio": "dummy",
		},
	}
	c, err := core.New(core.Config{AppConfig: cfg, DefaultDriverID: "dummy"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer c.Close()

	req := driver.AudioRequest{Task: driver.AudioTaskTextToSpeech, Prompt: "hello"}
	if got := c.ResolveAudioDriverID("", &req); got != "dummy" {
		t.Fatalf("got %q want dummy", got)
	}
}

func TestResolveAudioDriverIDRespectsExplicitBrandAlias(t *testing.T) {
	cfg := &config.Config{DefaultDriver: "dummy"}
	c, err := core.New(core.Config{AppConfig: cfg, DefaultDriverID: "dummy"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer c.Close()

	req := driver.AudioRequest{Task: driver.AudioTaskTextToSpeech, Prompt: "hello"}
	// onnx-tts not registered → piper brand driver would be used if registered; else stays piper id
	if got := c.ResolveAudioDriverID("piper", &req); got != "piper" {
		t.Fatalf("got %q want piper", got)
	}
}

func TestApplyAudioBrandAliasSetsPipeline(t *testing.T) {
	req := driver.AudioRequest{Prompt: "hi"}
	resolved, ok := driver.ApplyAudioBrandAlias("piper", &req, func(id string) bool {
		return id == driver.AudioRuntimeONNX
	})
	if !ok {
		t.Fatal("expected alias applied")
	}
	if resolved != driver.AudioRuntimeONNX {
		t.Fatalf("got runtime %q", resolved)
	}
	if req.Pipeline != "piper" {
		t.Fatalf("pipeline %q want piper", req.Pipeline)
	}
}
