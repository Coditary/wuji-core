package media

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

var audioFormatExtensions = map[string]string{
	"wav":  "wav",
	"mp3":  "mp3",
	"flac": "flac",
	"ogg":  "ogg",
	"opus": "opus",
	"m4a":  "m4a",
	"aac":  "m4a",
}

// ParseAudioFormat normalizes a user-facing output format name.
func ParseAudioFormat(raw string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.TrimPrefix(normalized, ".")
	if normalized == "" {
		return "", fmt.Errorf("audio format must not be empty")
	}
	if _, ok := audioFormatExtensions[normalized]; !ok {
		return "", fmt.Errorf("unsupported audio format %q (valid: wav, mp3, flac, ogg, opus, m4a)", raw)
	}
	if normalized == "aac" {
		return "m4a", nil
	}
	return normalized, nil
}

// AudioFormatFromPath guesses the format from a file extension.
func AudioFormatFromPath(path string) string {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	if ext == "aac" {
		return "m4a"
	}
	if _, ok := audioFormatExtensions[ext]; ok {
		return ext
	}
	return ""
}

// AudioFormatExtension returns the file extension for a normalized format.
func AudioFormatExtension(format string) string {
	if ext, ok := audioFormatExtensions[format]; ok {
		return ext
	}
	return format
}

// TranscodeAudio converts input audio to the requested format with ffmpeg.
func TranscodeAudio(ctx context.Context, ffmpegBin, inputPath, outputPath, format string, sampleRate int) error {
	format, err := ParseAudioFormat(format)
	if err != nil {
		return err
	}
	bin, err := ResolveFFmpegBin(ffmpegBin)
	if err != nil {
		return err
	}

	args := []string{"-y", "-i", inputPath, "-vn"}
	switch format {
	case "wav":
		args = append(args, "-acodec", "pcm_s16le")
	case "mp3":
		args = append(args, "-acodec", "libmp3lame", "-q:a", "2")
	case "flac":
		args = append(args, "-acodec", "flac")
	case "ogg":
		args = append(args, "-acodec", "libvorbis", "-q:a", "4")
	case "opus":
		args = append(args, "-acodec", "libopus", "-b:a", "128k")
	case "m4a":
		args = append(args, "-acodec", "aac", "-b:a", "192k")
	}
	if sampleRate > 0 {
		args = append(args, "-ar", fmt.Sprintf("%d", sampleRate))
	}
	args = append(args, outputPath)

	cmd := exec.CommandContext(ctx, bin, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("transcode audio to %s: %w\n%s", format, err, strings.TrimSpace(string(out)))
	}
	return nil
}
