package driver_test

import (
	"context"
	"os/exec"
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/driver/dummy"
)

func TestRunAudioTaskTranscodeToMP3(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not available")
	}

	d := dummy.New()
	resp, err := driver.RunAudioTask(context.Background(), d, driver.AudioRequest{
		Task: driver.AudioTaskTextToMusic, Prompt: "beat", Duration: 1, Format: "mp3",
	})
	if err != nil {
		t.Fatalf("RunAudioTask: %v", err)
	}
	if resp.Format != "mp3" || !hasSuffix(resp.Path, ".mp3") {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestRunAudioTaskOverlapStitch(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not available")
	}

	d := dummy.New()
	resp, err := driver.RunAudioTask(context.Background(), d, driver.AudioRequest{
		Task: driver.AudioTaskTextToMusic, Prompt: "long beat", Duration: 18, Overlap: 2,
	})
	if err != nil {
		t.Fatalf("RunAudioTask: %v", err)
	}
	if resp.Duration != 18 {
		t.Fatalf("duration = %g, want 18", resp.Duration)
	}
}

func hasSuffix(path, suffix string) bool {
	return len(path) >= len(suffix) && path[len(path)-len(suffix):] == suffix
}
