package core_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/core"
	"github.com/coditary/wuji-core/pkg/driver/dummy"
	ffmpegdrv "github.com/coditary/wuji-core/pkg/driver/ffmpeg"
	localdrv "github.com/coditary/wuji-core/pkg/driver/local"
)

func TestResolveDriverIDPrecedence(t *testing.T) {
	cfg := &config.Config{
		DefaultDriver: "dummy",
		CapabilityDrivers: config.CapabilityDrivers{
			"text":  "llama",
			"image": "a1111",
		},
	}
	c, err := core.New(core.Config{AppConfig: cfg, DefaultDriverID: "dummy"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer c.Close()

	tests := []struct {
		name      string
		preferred string
		cap       capability.Type
		want      string
	}{
		{"explicit override", "dummy", capability.TextGeneration, "dummy"},
		{"capability driver", "", capability.TextGeneration, "llama"},
		{"capability driver image", "", capability.ImageGeneration, "a1111"},
		{"default driver", "", capability.Mesh, "dummy"},
		{"video2audio default", "", capability.Video2Audio, ffmpegdrv.DriverID},
		{"rag default", "", capability.RAG, "raggo"},
		{"rag legacy alias", "local", capability.RAG, localdrv.DriverID},
	}
	for _, tc := range tests {
		got := c.ResolveDriverID(tc.preferred, tc.cap)
		if got != tc.want {
			t.Errorf("%s: got %q want %q", tc.name, got, tc.want)
		}
	}
}

func TestGenerateTextRespectsPreferredOverCapabilityDriver(t *testing.T) {
	cfg := &config.Config{
		DefaultDriver: dummy.DriverID,
		CapabilityDrivers: config.CapabilityDrivers{
			"text": "llama",
		},
	}
	c, err := core.New(core.Config{AppConfig: cfg, DefaultDriverID: dummy.DriverID})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer c.Close()

	got := c.ResolveDriverID(dummy.DriverID, capability.TextGeneration)
	if got != dummy.DriverID {
		t.Fatalf("expected preferred driver dummy, got %q", got)
	}
}
