package driver_test

import (
	"context"
	"strings"
	"testing"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/driver/dummy"
)

func TestRunTextInputImage(t *testing.T) {
	d := dummy.New()
	resp, err := driver.RunTextInput(context.Background(), d, driver.TextRequest{
		InputMode: driver.TextInputImage,
		MediaPath: "photo.jpg",
		Prompt:    "what is this?",
	})
	if err != nil {
		t.Fatalf("RunTextInput: %v", err)
	}
	if !strings.Contains(resp.Text, "image2text") {
		t.Fatalf("unexpected response: %q", resp.Text)
	}
}

func TestRunTextInputAudio(t *testing.T) {
	d := dummy.New()
	resp, err := driver.RunTextInput(context.Background(), d, driver.TextRequest{
		InputMode: driver.TextInputAudio,
		MediaPath: "speech.wav",
	})
	if err != nil {
		t.Fatalf("RunTextInput: %v", err)
	}
	if !strings.Contains(resp.Text, "audio2text") {
		t.Fatalf("unexpected response: %q", resp.Text)
	}
}

func TestRunTextInputUnsupportedVideo(t *testing.T) {
	d := textOnlyDriver{}
	resp, err := driver.RunTextInput(context.Background(), d, driver.TextRequest{
		InputMode: driver.TextInputVideo,
		MediaPath: "clip.mp4",
	})
	if err == nil {
		t.Fatalf("expected error, got %q", resp.Text)
	}
}

type textOnlyDriver struct{}

func (textOnlyDriver) Info() driver.Info {
	return driver.Info{ID: "text-only", Capabilities: []capability.Type{capability.TextGeneration}}
}
func (textOnlyDriver) Capabilities() []capability.Type { return []capability.Type{capability.TextGeneration} }
func (textOnlyDriver) Close() error { return nil }
func (textOnlyDriver) GenerateText(_ context.Context, req driver.TextRequest) (*driver.TextResponse, error) {
	return &driver.TextResponse{Text: req.Prompt}, nil
}
