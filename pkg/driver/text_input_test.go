package driver_test

import (
	"strings"
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
)

func TestTextRequestValidatePrompt(t *testing.T) {
	req := driver.TextRequest{InputMode: driver.TextInputPrompt, Prompt: "hi"}
	if err := req.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestTextRequestValidatePromptEmpty(t *testing.T) {
	req := driver.TextRequest{InputMode: driver.TextInputPrompt}
	if err := req.Validate(); err == nil {
		t.Fatal("expected error for empty prompt")
	}
}

func TestTextRequestValidateTranslateWithMedia(t *testing.T) {
	req := driver.TextRequest{
		InputMode: driver.TextInputAudio,
		MediaPath: "x.wav",
		Translate: true,
	}
	if err := req.Validate(); err == nil {
		t.Fatal("expected error for --translate with media input")
	}
}

func TestTextRequestForGenerationTranslate(t *testing.T) {
	req := driver.TextRequest{
		Prompt:    "Hallo",
		Translate: true,
		TargetLang: "French",
	}
	out := req.ForGeneration()
	if out.SystemPrompt == "" {
		t.Fatal("expected translation system prompt")
	}
	if out.Prompt != "Hallo" {
		t.Fatalf("prompt changed: %q", out.Prompt)
	}
}

func TestTextRequestForGenerationTranslateWithSourceLang(t *testing.T) {
	req := driver.TextRequest{
		Prompt:     "Hallo",
		Translate:  true,
		TargetLang: "French",
		Language:   "de",
	}
	out := req.ForGeneration()
	if out.SystemPrompt == "" || !containsAll(out.SystemPrompt, "de", "French") {
		t.Fatalf("unexpected system prompt: %q", out.SystemPrompt)
	}
}

func containsAll(s string, parts ...string) bool {
	for _, p := range parts {
		if !strings.Contains(s, p) {
			return false
		}
	}
	return true
}
