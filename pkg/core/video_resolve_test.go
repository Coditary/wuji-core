package core_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/core"
	"github.com/coditary/wuji-core/pkg/driver"
	ffmpegdrv "github.com/coditary/wuji-core/pkg/driver/ffmpeg"
)

func TestResolveVideoDriverIDRoutesTransformsToFFmpeg(t *testing.T) {
	cfg := &config.Config{
		DefaultDriver: "dummy",
		CapabilityDrivers: config.CapabilityDrivers{
			"video": "a1111",
		},
	}
	c, err := core.New(core.Config{AppConfig: cfg, DefaultDriverID: "dummy"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer c.Close()

	interpolate := driver.VideoRequest{Task: driver.VideoTaskInterpolate, InitVideoPath: "clip.mp4", FPS: 60}
	if got := c.ResolveVideoDriverID("", interpolate); got != ffmpegdrv.DriverID {
		t.Fatalf("interpolate: got %q want %q", got, ffmpegdrv.DriverID)
	}

	scale := driver.VideoRequest{Task: driver.VideoTaskScale, InitVideoPath: "clip.mp4", Scale: 2}
	if got := c.ResolveVideoDriverID("", scale); got != ffmpegdrv.DriverID {
		t.Fatalf("scale: got %q want %q", got, ffmpegdrv.DriverID)
	}

	generate := driver.VideoRequest{Task: driver.VideoTaskGenerate, Prompt: "waves"}
	if got := c.ResolveVideoDriverID("", generate); got != "a1111" {
		t.Fatalf("generate: got %q want %q", got, "a1111")
	}
}

func TestResolveVideoDriverIDRespectsExplicitOverride(t *testing.T) {
	cfg := &config.Config{
		CapabilityDrivers: config.CapabilityDrivers{
			"video": "a1111",
		},
	}
	c, err := core.New(core.Config{AppConfig: cfg, DefaultDriverID: "dummy"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer c.Close()

	req := driver.VideoRequest{Task: driver.VideoTaskScale, InitVideoPath: "clip.mp4", Scale: 2}
	if got := c.ResolveVideoDriverID("dummy", req); got != "dummy" {
		t.Fatalf("expected explicit override dummy, got %q", got)
	}
}

func TestResolveVideoDriverIDGenerateUsesCapabilityDriver(t *testing.T) {
	cfg := &config.Config{
		CapabilityDrivers: config.CapabilityDrivers{
			string(capability.VideoGeneration): "forge",
		},
	}
	c, err := core.New(core.Config{AppConfig: cfg, DefaultDriverID: "dummy"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer c.Close()

	req := driver.VideoRequest{Task: driver.VideoTaskGenerate, Prompt: "cat"}
	if got := c.ResolveVideoDriverID("", req); got != "forge" {
		t.Fatalf("got %q want forge", got)
	}
}
