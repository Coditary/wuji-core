package media_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/coditary/wuji-core/pkg/media"
)

func TestWritePlaceholderWAV(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.wav")
	if err := media.WritePlaceholderWAV(path, 1, 44100); err != nil {
		t.Fatalf("WritePlaceholderWAV: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	// 44 byte header + 1s * 44100 samples * 2 bytes
	want := 44 + 44100*2
	if info.Size() != int64(want) {
		t.Fatalf("size = %d, want %d", info.Size(), want)
	}
}
