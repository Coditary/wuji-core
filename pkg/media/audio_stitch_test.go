package media_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/coditary/wuji-core/pkg/media"
)

func TestConcatAudioWithCrossfade(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not available")
	}

	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.wav")
	p2 := filepath.Join(dir, "b.wav")
	out := filepath.Join(dir, "out.wav")
	if err := media.WritePlaceholderWAV(p1, 10, 44100); err != nil {
		t.Fatal(err)
	}
	if err := media.WritePlaceholderWAV(p2, 10, 44100); err != nil {
		t.Fatal(err)
	}
	if err := media.ConcatAudioWithCrossfade(context.Background(), "", []string{p1, p2}, 2, out); err != nil {
		t.Fatalf("ConcatAudioWithCrossfade: %v", err)
	}
	info, err := os.Stat(out)
	if err != nil || info.Size() < 44 {
		t.Fatalf("output file: %v size=%d", err, info.Size())
	}
}
