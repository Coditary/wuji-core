package media_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/media"
)

func TestParseAudioFormat(t *testing.T) {
	got, err := media.ParseAudioFormat("mp3")
	if err != nil || got != "mp3" {
		t.Fatalf("mp3: got %q err %v", got, err)
	}
	got, err = media.ParseAudioFormat(".flac")
	if err != nil || got != "flac" {
		t.Fatalf("flac: got %q err %v", got, err)
	}
	got, err = media.ParseAudioFormat("aac")
	if err != nil || got != "m4a" {
		t.Fatalf("aac: got %q err %v", got, err)
	}
	if _, err := media.ParseAudioFormat("wma"); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
