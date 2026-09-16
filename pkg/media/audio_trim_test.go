package media_test

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/coditary/wuji-core/pkg/media"
)

func TestTrimAudio(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not available")
	}

	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.wav")
	p2 := filepath.Join(dir, "b.wav")
	stitched := filepath.Join(dir, "stitched.wav")
	trimmed := filepath.Join(dir, "trimmed.wav")
	if err := media.WritePlaceholderWAV(p1, 10, 44100); err != nil {
		t.Fatal(err)
	}
	if err := media.WritePlaceholderWAV(p2, 10, 44100); err != nil {
		t.Fatal(err)
	}
	if err := media.ConcatAudioWithCrossfade(context.Background(), "", []string{p1, p2}, 2, stitched); err != nil {
		t.Fatal(err)
	}
	if err := media.TrimAudio(context.Background(), "", stitched, trimmed, 18); err != nil {
		t.Fatalf("TrimAudio: %v", err)
	}
}
