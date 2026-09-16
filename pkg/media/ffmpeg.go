package media

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var videoExtensions = map[string]bool{
	".mp4": true, ".mkv": true, ".webm": true, ".mov": true,
	".avi": true, ".flv": true, ".m4v": true, ".wmv": true, ".mpeg": true, ".mpg": true,
}

// IsVideoContainer reports whether the path looks like a video file by extension.
func IsVideoContainer(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return videoExtensions[ext]
}

// ResolveFFmpegBin returns the ffmpeg binary path from bin override, PATH, or an error.
func ResolveFFmpegBin(override string) (string, error) {
	if strings.TrimSpace(override) != "" {
		return strings.TrimSpace(override), nil
	}
	if env := strings.TrimSpace(os.Getenv("WUJI_FFMPEG_PATH")); env != "" {
		return env, nil
	}
	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", fmt.Errorf("ffmpeg is required for video2audio (install ffmpeg, set WUJI_FFMPEG_PATH, or configure ffmpeg.bin)")
	}
	return path, nil
}

// ExtractAudioFromVideo writes a mono 16 kHz WAV file for STT backends.
func ExtractAudioFromVideo(ctx context.Context, ffmpegBin, videoPath string) (audioPath string, err error) {
	videoPath = strings.TrimSpace(videoPath)
	if videoPath == "" {
		return "", fmt.Errorf("video path is required")
	}
	if _, err := os.Stat(videoPath); err != nil {
		return "", fmt.Errorf("video file %q: %w", videoPath, err)
	}

	bin, err := ResolveFFmpegBin(ffmpegBin)
	if err != nil {
		return "", err
	}

	tmp, err := os.CreateTemp("", "wuji-audio-*.wav")
	if err != nil {
		return "", fmt.Errorf("create temp audio file: %w", err)
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()

	cmd := exec.CommandContext(ctx, bin,
		"-y", "-i", videoPath,
		"-vn", "-acodec", "pcm_s16le", "-ar", "16000", "-ac", "1",
		tmpPath,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		_ = os.Remove(tmpPath)
		return "", fmt.Errorf("extract audio from %q: %w\n%s", videoPath, err, strings.TrimSpace(string(out)))
	}

	return tmpPath, nil
}
