package dummy_test

import (
	"context"
	"testing"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/driver/dummy"
)

func TestDummyDriverCapabilities(t *testing.T) {
	d := dummy.New()
	caps := d.Capabilities()

	expected := []capability.Type{
		capability.TextGeneration, capability.ImageGeneration,
		capability.ImageUpscale, capability.ImageDownscale, capability.ImageScale,
		capability.VideoGeneration,
		capability.AudioGeneration, capability.Mesh,
		capability.Video2Audio, capability.VoiceCloning, capability.Training, capability.DatasetMgmt,
		capability.Data,
		capability.RAG,
	}

	if len(caps) != len(expected) {
		t.Fatalf("expected %d capabilities, got %d", len(expected), len(caps))
	}

	for _, c := range expected {
		if !driver.HasCapability(d, c) {
			t.Fatalf("dummy driver should support %s", c)
		}
	}
}

func TestDummyDriverGenerateText(t *testing.T) {
	d := dummy.New()
	resp, err := d.GenerateText(context.Background(), driver.TextRequest{Prompt: "hello"})
	if err != nil {
		t.Fatalf("GenerateText: %v", err)
	}
	if resp.Text == "" {
		t.Fatal("expected non-empty response")
	}
}

func TestDummyDriverInfo(t *testing.T) {
	d := dummy.New()
	if d.Info().ID != dummy.DriverID {
		t.Fatalf("expected ID %q, got %q", dummy.DriverID, d.Info().ID)
	}
}
