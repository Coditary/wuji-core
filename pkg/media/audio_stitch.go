package media

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ConcatAudioWithCrossfade stitches audio files with a triangular crossfade overlap.
func ConcatAudioWithCrossfade(ctx context.Context, ffmpegBin string, paths []string, overlapSec float32, outputPath string) error {
	if len(paths) == 0 {
		return fmt.Errorf("no audio segments to stitch")
	}
	if len(paths) == 1 {
		return copyFile(paths[0], outputPath)
	}
	if overlapSec <= 0 {
		return concatAudio(ctx, ffmpegBin, paths, outputPath)
	}

	bin, err := ResolveFFmpegBin(ffmpegBin)
	if err != nil {
		return err
	}

	dir := filepath.Dir(outputPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}
	}

	args := []string{"-y"}
	for _, path := range paths {
		args = append(args, "-i", path)
	}

	filter := buildAcrossfadeFilter(len(paths), overlapSec)
	args = append(args, "-filter_complex", filter, "-map", "[out]", outputPath)

	cmd := exec.CommandContext(ctx, bin, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("stitch audio segments: %w\n%s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// TrimAudio truncates an audio file to the given duration in seconds.
func TrimAudio(ctx context.Context, ffmpegBin, inputPath, outputPath string, durationSec float32) error {
	if durationSec <= 0 {
		return copyFile(inputPath, outputPath)
	}
	bin, err := ResolveFFmpegBin(ffmpegBin)
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, bin,
		"-y", "-i", inputPath,
		"-t", fmt.Sprintf("%.3f", durationSec),
		"-c", "copy",
		outputPath,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		// Stream copy may fail for some codecs; re-encode on fallback.
		cmd = exec.CommandContext(ctx, bin,
			"-y", "-i", inputPath,
			"-t", fmt.Sprintf("%.3f", durationSec),
			outputPath,
		)
		if out, err = cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("trim audio: %w\n%s", err, strings.TrimSpace(string(out)))
		}
	}
	return nil
}

func concatAudio(ctx context.Context, ffmpegBin string, paths []string, outputPath string) error {
	bin, err := ResolveFFmpegBin(ffmpegBin)
	if err != nil {
		return err
	}

	listFile, err := os.CreateTemp("", "wuji-concat-*.txt")
	if err != nil {
		return fmt.Errorf("create concat list: %w", err)
	}
	listPath := listFile.Name()
	defer os.Remove(listPath)

	for _, path := range paths {
		escaped := strings.ReplaceAll(path, "'", "'\\''")
		if _, err := fmt.Fprintf(listFile, "file '%s'\n", escaped); err != nil {
			_ = listFile.Close()
			return err
		}
	}
	if err := listFile.Close(); err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, bin, "-y", "-f", "concat", "-safe", "0", "-i", listPath, "-c", "copy", outputPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("concat audio: %w\n%s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func buildAcrossfadeFilter(count int, overlapSec float32) string {
	if count < 2 {
		return "[0]anull[out]"
	}
	d := fmt.Sprintf("%.3f", overlapSec)
	if count == 2 {
		return "[0][1]acrossfade=d=" + d + ":c1=tri:c2=tri[out]"
	}

	var filter string
	for i := 1; i < count; i++ {
		outLabel := "out"
		if i < count-1 {
			outLabel = fmt.Sprintf("a%02d", i)
		}
		if i == 1 {
			filter = fmt.Sprintf("[0][1]acrossfade=d=%s:c1=tri:c2=tri[%s]", d, outLabel)
			continue
		}
		prev := fmt.Sprintf("a%02d", i-1)
		filter += fmt.Sprintf(";[%s][%d]acrossfade=d=%s:c1=tri:c2=tri[%s]", prev, i, d, outLabel)
	}
	return filter
}

func copyFile(src, dst string) error {
	in, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	dir := filepath.Dir(dst)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(dst, in, 0o644)
}
