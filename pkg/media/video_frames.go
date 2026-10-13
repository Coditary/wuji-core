package media

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const FrameFilePattern = "frame_%06d.png"

// ExtractVideoFrames writes PNG frames from a video into outDir.
func ExtractVideoFrames(ctx context.Context, ffmpegBinOverride, videoPath, outDir string) (VideoMeta, error) {
	videoPath = strings.TrimSpace(videoPath)
	if videoPath == "" {
		return VideoMeta{}, fmt.Errorf("video path is required")
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return VideoMeta{}, fmt.Errorf("create frame dir: %w", err)
	}

	meta, err := ProbeVideo(ctx, ffmpegBinOverride, videoPath)
	if err != nil {
		return VideoMeta{}, err
	}

	bin, err := ResolveFFmpegBin(ffmpegBinOverride)
	if err != nil {
		return VideoMeta{}, err
	}

	pattern := filepath.Join(outDir, FrameFilePattern)
	cmd := exec.CommandContext(ctx, bin, "-y", "-i", videoPath, "-vsync", "0", pattern)
	if out, runErr := cmd.CombinedOutput(); runErr != nil {
		return VideoMeta{}, fmt.Errorf("extract frames from %q: %w\n%s", videoPath, runErr, strings.TrimSpace(string(out)))
	}
	return meta, nil
}

// CountPNGFrames counts PNG files in a directory.
func CountPNGFrames(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
			count++
		}
	}
	return count, nil
}

// AssembleVideoFromFrames encodes numbered PNG frames to MP4 and optionally muxes audio from sourceVideo.
func AssembleVideoFromFrames(ctx context.Context, ffmpegBinOverride, framesDir string, fps int, sourceVideo string) (string, VideoMeta, error) {
	if fps <= 0 {
		return "", VideoMeta{}, fmt.Errorf("fps must be > 0")
	}

	bin, err := ResolveFFmpegBin(ffmpegBinOverride)
	if err != nil {
		return "", VideoMeta{}, err
	}

	outPath, err := tempVideoPath()
	if err != nil {
		return "", VideoMeta{}, err
	}

	pattern := filepath.Join(framesDir, FrameFilePattern)
	args := []string{"-y", "-framerate", fmt.Sprintf("%d", fps), "-i", pattern}

	audioPath := ""
	if strings.TrimSpace(sourceVideo) != "" {
		audioPath, err = extractAudioCopy(ctx, bin, sourceVideo)
		if err == nil && audioPath != "" {
			args = append(args, "-i", audioPath)
		}
	}

	args = append(args,
		"-c:v", "libx264", "-crf", "18", "-preset", "medium", "-pix_fmt", "yuv420p",
		"-movflags", "+faststart",
	)
	if audioPath != "" {
		args = append(args, "-c:a", "aac", "-shortest")
	} else {
		args = append(args, "-an")
	}
	args = append(args, outPath)

	cmd := exec.CommandContext(ctx, bin, args...)
	if out, runErr := cmd.CombinedOutput(); runErr != nil {
		_ = os.Remove(outPath)
		if audioPath != "" {
			_ = os.Remove(audioPath)
		}
		return "", VideoMeta{}, fmt.Errorf("assemble video: %w\n%s", runErr, strings.TrimSpace(string(out)))
	}
	if audioPath != "" {
		_ = os.Remove(audioPath)
	}

	meta, err := ProbeVideo(ctx, ffmpegBinOverride, outPath)
	if err != nil {
		return outPath, VideoMeta{FPS: fps}, nil
	}
	return outPath, meta, nil
}

func extractAudioCopy(ctx context.Context, ffmpegBin, videoPath string) (string, error) {
	tmp, err := os.CreateTemp("", "wuji-audio-mux-*")
	if err != nil {
		return "", err
	}
	audioPath := tmp.Name() + ".m4a"
	_ = tmp.Close()

	cmd := exec.CommandContext(ctx, ffmpegBin, "-y", "-i", videoPath, "-vn", "-acodec", "copy", audioPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		_ = os.Remove(audioPath)
		return "", fmt.Errorf("extract audio track: %w\n%s", err, strings.TrimSpace(string(out)))
	}
	return audioPath, nil
}

// RIFEPassesNeeded returns how many 2x RIFE passes are required to approach targetMultiplier.
func RIFEPassesNeeded(targetMultiplier float64) int {
	if targetMultiplier <= 1 {
		return 0
	}
	passes := 0
	m := 1.0
	for m < targetMultiplier {
		m *= 2
		passes++
	}
	return passes
}
