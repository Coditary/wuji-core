package media_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/media"
)

func TestIsVideoContainer(t *testing.T) {
	if !media.IsVideoContainer("clip.MP4") {
		t.Fatal("expected mp4 to be a video container")
	}
	if media.IsVideoContainer("speech.wav") {
		t.Fatal("expected wav not to be a video container")
	}
}
