package media

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// VideoMeta describes a probed video stream.
type VideoMeta struct {
	Duration float32
	FPS      int
	Frames   int
	Width    int
	Height   int
}

// ResolveFFprobeBin returns the ffprobe binary path.
func ResolveFFprobeBin(ffmpegBinOverride string) (string, error) {
	if env := strings.TrimSpace(os.Getenv("WUJI_FFPROBE_PATH")); env != "" {
		return env, nil
	}
	if bin := strings.TrimSpace(ffmpegBinOverride); bin != "" {
		dir := filepath.Dir(bin)
		candidate := filepath.Join(dir, "ffprobe")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	path, err := exec.LookPath("ffprobe")
	if err != nil {
		return "", fmt.Errorf("ffprobe is required for video processing (install ffmpeg, set WUJI_FFPROBE_PATH, or configure ffmpeg.bin)")
	}
	return path, nil
}

// ProbeVideo reads duration, frame rate, and dimensions from a video file.
func ProbeVideo(ctx context.Context, ffmpegBinOverride, videoPath string) (VideoMeta, error) {
	videoPath = strings.TrimSpace(videoPath)
	if videoPath == "" {
		return VideoMeta{}, fmt.Errorf("video path is required")
	}
	if _, err := os.Stat(videoPath); err != nil {
		return VideoMeta{}, fmt.Errorf("video file %q: %w", videoPath, err)
	}

	probeBin, err := ResolveFFprobeBin(ffmpegBinOverride)
	if err != nil {
		return VideoMeta{}, err
	}

	cmd := exec.CommandContext(ctx, probeBin,
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=avg_frame_rate,width,height,nb_frames",
		"-show_entries", "format=duration",
		"-of", "json",
		videoPath,
	)
	out, err := cmd.Output()
	if err != nil {
		return VideoMeta{}, fmt.Errorf("probe video %q: %w", videoPath, err)
	}

	var payload struct {
		Streams []struct {
			Width        int    `json:"width"`
			Height       int    `json:"height"`
			NbFrames     string `json:"nb_frames"`
			AvgFrameRate string `json:"avg_frame_rate"`
		} `json:"streams"`
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}
	if err := json.Unmarshal(out, &payload); err != nil {
		return VideoMeta{}, fmt.Errorf("parse ffprobe output: %w", err)
	}
	if len(payload.Streams) == 0 {
		return VideoMeta{}, fmt.Errorf("no video stream in %q", videoPath)
	}
	stream := payload.Streams[0]

	fps, err := parseFrameRate(stream.AvgFrameRate)
	if err != nil {
		return VideoMeta{}, fmt.Errorf("parse frame rate for %q: %w", videoPath, err)
	}

	duration, _ := strconv.ParseFloat(strings.TrimSpace(payload.Format.Duration), 64)
	frames := 0
	if stream.NbFrames != "" && stream.NbFrames != "N/A" {
		frames, _ = strconv.Atoi(stream.NbFrames)
	}
	if frames <= 0 && duration > 0 && fps > 0 {
		frames = int(math.Round(duration * float64(fps)))
	}

	return VideoMeta{
		Duration: float32(duration),
		FPS:      fps,
		Frames:   frames,
		Width:    stream.Width,
		Height:   stream.Height,
	}, nil
}

func parseFrameRate(raw string) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "0/0" {
		return 0, fmt.Errorf("missing frame rate")
	}
	if strings.Contains(raw, "/") {
		parts := strings.SplitN(raw, "/", 2)
		num, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return 0, err
		}
		den, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil || den == 0 {
			return 0, fmt.Errorf("invalid frame rate %q", raw)
		}
		return int(math.Round(num / den)), nil
	}
	f, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, err
	}
	return int(math.Round(f)), nil
}

// ScaleVideo resizes a video by the given factor using ffmpeg (Lanczos).
func ScaleVideo(ctx context.Context, ffmpegBinOverride, videoPath string, scale float32) (string, VideoMeta, error) {
	if scale <= 0 {
		return "", VideoMeta{}, fmt.Errorf("scale factor must be > 0")
	}
	if !IsVideoContainer(videoPath) {
		return "", VideoMeta{}, fmt.Errorf("%q is not a supported video container", videoPath)
	}

	bin, err := ResolveFFmpegBin(ffmpegBinOverride)
	if err != nil {
		return "", VideoMeta{}, err
	}

	outPath, err := tempVideoPath()
	if err != nil {
		return "", VideoMeta{}, err
	}

	filter := fmt.Sprintf("scale=trunc(iw*%.6g/2)*2:trunc(ih*%.6g/2)*2", scale, scale)
	cmd := exec.CommandContext(ctx, bin,
		"-y", "-i", videoPath,
		"-vf", filter,
		"-c:v", "libx264", "-crf", "18", "-preset", "medium",
		"-c:a", "copy",
		outPath,
	)
	if out, runErr := cmd.CombinedOutput(); runErr != nil {
		_ = os.Remove(outPath)
		return "", VideoMeta{}, fmt.Errorf("scale video %q: %w\n%s", videoPath, runErr, strings.TrimSpace(string(out)))
	}

	meta, err := ProbeVideo(ctx, ffmpegBinOverride, outPath)
	if err != nil {
		return outPath, VideoMeta{}, nil
	}
	return outPath, meta, nil
}

// InterpolateVideo changes the frame rate using ffmpeg (minterpolate when upsampling, fps when downsampling).
func InterpolateVideo(ctx context.Context, ffmpegBinOverride, videoPath string, targetFPS int) (string, VideoMeta, error) {
	if targetFPS <= 0 {
		return "", VideoMeta{}, fmt.Errorf("target fps must be > 0")
	}
	if !IsVideoContainer(videoPath) {
		return "", VideoMeta{}, fmt.Errorf("%q is not a supported video container", videoPath)
	}

	bin, err := ResolveFFmpegBin(ffmpegBinOverride)
	if err != nil {
		return "", VideoMeta{}, err
	}

	src, err := ProbeVideo(ctx, ffmpegBinOverride, videoPath)
	if err != nil {
		return "", VideoMeta{}, err
	}
	if src.FPS > 0 && targetFPS <= src.FPS {
		return "", VideoMeta{}, fmt.Errorf("target fps %d must be greater than source fps %d for interpolation", targetFPS, src.FPS)
	}

	outPath, err := tempVideoPath()
	if err != nil {
		return "", VideoMeta{}, err
	}

	filter := fmt.Sprintf("minterpolate=fps=%d:mi_mode=mci:mc_mode=aobmc:me_mode=bidir:vsbmc=1", targetFPS)
	cmd := exec.CommandContext(ctx, bin,
		"-y", "-i", videoPath,
		"-vf", filter,
		"-c:v", "libx264", "-crf", "18", "-preset", "medium",
		"-c:a", "copy",
		outPath,
	)
	if out, runErr := cmd.CombinedOutput(); runErr != nil {
		_ = os.Remove(outPath)
		return "", VideoMeta{}, fmt.Errorf("interpolate video %q: %w\n%s", videoPath, runErr, strings.TrimSpace(string(out)))
	}

	meta, err := ProbeVideo(ctx, ffmpegBinOverride, outPath)
	if err != nil {
		return outPath, VideoMeta{FPS: targetFPS}, nil
	}
	return outPath, meta, nil
}

func tempVideoPath() (string, error) {
	tmp, err := os.CreateTemp("", "wuji-video-*.mp4")
	if err != nil {
		return "", fmt.Errorf("create temp video file: %w", err)
	}
	path := tmp.Name()
	_ = tmp.Close()
	return path, nil
}

// ConvertToMP4 transcodes webp/gif/mkv inputs to mp4 for uniform Wuji video output.
func ConvertToMP4(ctx context.Context, ffmpegBinOverride, inputPath string) (string, VideoMeta, error) {
	bin, err := ResolveFFmpegBin(ffmpegBinOverride)
	if err != nil {
		return "", VideoMeta{}, err
	}

	outPath, err := tempVideoPath()
	if err != nil {
		return "", VideoMeta{}, err
	}

	cmd := exec.CommandContext(ctx, bin,
		"-y", "-i", inputPath,
		"-c:v", "libx264", "-crf", "18", "-preset", "medium", "-pix_fmt", "yuv420p",
		"-movflags", "+faststart",
		"-an",
		outPath,
	)
	if out, runErr := cmd.CombinedOutput(); runErr != nil {
		_ = os.Remove(outPath)
		return "", VideoMeta{}, fmt.Errorf("convert to mp4 %q: %w\n%s", inputPath, runErr, strings.TrimSpace(string(out)))
	}

	meta, err := ProbeVideo(ctx, ffmpegBinOverride, outPath)
	if err != nil {
		return outPath, VideoMeta{}, nil
	}
	return outPath, meta, nil
}
