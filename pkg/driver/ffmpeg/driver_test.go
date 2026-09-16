package ffmpeg_test

import (
	"context"
	"testing"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/driver/ffmpeg"
)

func TestFFmpegDriverVideoToAudioRequiresVideo(t *testing.T) {
	d := ffmpeg.New()
	_, err := d.VideoToAudio(context.Background(), driver.Video2AudioRequest{VideoPath: "notes.txt"})
	if err == nil {
		t.Fatal("expected error for non-video path")
	}
}

func TestFFmpegDriverCapability(t *testing.T) {
	d := ffmpeg.New()
	if !driver.HasCapability(d, capability.Video2Audio) {
		t.Fatal("ffmpeg driver should advertise video2audio")
	}
}
