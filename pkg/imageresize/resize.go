package imageresize

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const defaultOutputDir = "/tmp/wuji/scale"

// ScaleImage resizes an image by factor using ImageMagick (magick or convert).
// factor 0.5 halves the image; 2.5 enlarges to 250%.
func ScaleImage(ctx context.Context, inputPath, outputDir string, factor float32) (string, error) {
	if strings.TrimSpace(inputPath) == "" {
		return "", fmt.Errorf("input image path is required")
	}
	if factor <= 0 {
		return "", fmt.Errorf("scale factor must be > 0 (got %g)", factor)
	}
	if _, err := os.Stat(inputPath); err != nil {
		return "", fmt.Errorf("input image: %w", err)
	}

	bin, err := lookupMagick()
	if err != nil {
		return "", err
	}

	if outputDir == "" {
		outputDir = defaultOutputDir
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", fmt.Errorf("create output dir: %w", err)
	}

	ext := filepath.Ext(inputPath)
	if ext == "" {
		ext = ".png"
	}
	outPath := filepath.Join(outputDir, fmt.Sprintf("wuji-scale-%d%s", time.Now().UnixNano(), ext))
	resizeArg := formatResizeArg(factor)

	cmd := exec.CommandContext(ctx, bin, inputPath, "-filter", "Lanczos", "-resize", resizeArg, outPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("magick resize: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return outPath, nil
}

func lookupMagick() (string, error) {
	if path, err := exec.LookPath("magick"); err == nil {
		return path, nil
	}
	if path, err := exec.LookPath("convert"); err == nil {
		return path, nil
	}
	return "", fmt.Errorf("ImageMagick not found (install magick or convert for arbitrary scale factors)")
}

func formatResizeArg(factor float32) string {
	pct := float64(factor) * 100
	pctStr := strconv.FormatFloat(pct, 'f', -1, 64)
	return pctStr + "%"
}
