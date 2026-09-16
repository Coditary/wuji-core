package driver

import (
	"fmt"
	"strings"
)

// VideoTask identifies the video generation or transformation mode.
type VideoTask string

const (
	VideoTaskGenerate    VideoTask = "generate"
	VideoTaskImageToVideo VideoTask = "i2v"
	VideoTaskInterpolate VideoTask = "interpolate"
	VideoTaskScale       VideoTask = "scale"
)

// VideoTaskInfo describes a task and its required inputs.
type VideoTaskInfo struct {
	Task              VideoTask
	Description       string
	RequiresPrompt    bool
	RequiresInitImage bool
	RequiresInitVideo bool
}

func AllVideoTasks() []VideoTask {
	return []VideoTask{
		VideoTaskGenerate,
		VideoTaskImageToVideo,
		VideoTaskInterpolate,
		VideoTaskScale,
	}
}

func VideoTaskCatalog() []VideoTaskInfo {
	return []VideoTaskInfo{
		{Task: VideoTaskGenerate, Description: "Text-to-video generation from a prompt", RequiresPrompt: true},
		{Task: VideoTaskImageToVideo, Description: "Animate an image into video (image-to-video)", RequiresPrompt: true, RequiresInitImage: true},
		{Task: VideoTaskInterpolate, Description: "Increase frame rate via interpolation (RIFE, FILM, …)", RequiresInitVideo: true},
		{Task: VideoTaskScale, Description: "Scale video up or down (super-resolution / resize)", RequiresInitVideo: true},
	}
}

func (t VideoTask) String() string { return string(t) }

func (t VideoTask) IsValid() bool {
	for _, known := range AllVideoTasks() {
		if t == known {
			return true
		}
	}
	return false
}

func (t VideoTask) Info() *VideoTaskInfo {
	for _, info := range VideoTaskCatalog() {
		if info.Task == t {
			copy := info
			return &copy
		}
	}
	return nil
}

func ParseVideoTask(raw string) (VideoTask, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	switch normalized {
	case "", "generate", "t2v", "text-to-video", "txt2vid":
		return VideoTaskGenerate, nil
	case "i2v", "image-to-video", "img2vid":
		return VideoTaskImageToVideo, nil
	case "interpolate", "interpolation", "frame-interpolation":
		return VideoTaskInterpolate, nil
	case "scale":
		return VideoTaskScale, nil
	case "upscale", "super-resolution", "sr":
		return VideoTaskScale, nil
	}
	task := VideoTask(normalized)
	if !task.IsValid() {
		return "", fmt.Errorf("unknown video task %q (valid: %s)", raw, joinVideoTasks())
	}
	return task, nil
}

func (r VideoRequest) TaskOrDefault() VideoTask {
	if r.Task == "" {
		return VideoTaskGenerate
	}
	return r.Task
}

// VideoTaskInputs collects CLI fields used to infer a video task from flags.
type VideoTaskInputs struct {
	ImagePath string
	VideoPath string
	Scale     float32
	ScaleSet  bool
	FPSSet    bool
}

// InferVideoTask selects a task from input flags when --task is omitted.
func InferVideoTask(in VideoTaskInputs) (VideoTask, error) {
	hasImage := strings.TrimSpace(in.ImagePath) != ""
	hasVideo := strings.TrimSpace(in.VideoPath) != ""

	if hasImage && hasVideo {
		return "", fmt.Errorf("use either --image or --video, not both")
	}
	if hasImage {
		return VideoTaskImageToVideo, nil
	}
	if hasVideo {
		if in.ScaleSet && in.Scale > 0 {
			if in.FPSSet {
				return "", fmt.Errorf("with --video, use either --scale for scaling or --fps for interpolation, not both")
			}
			return VideoTaskScale, nil
		}
		if in.FPSSet {
			return VideoTaskInterpolate, nil
		}
		return "", fmt.Errorf("with --video, set --fps for frame interpolation or --scale (e.g. 2, 2.5, 50%%, 300%%)")
	}
	return VideoTaskGenerate, nil
}

func (r VideoRequest) Validate() error {
	task := r.TaskOrDefault()
	if !task.IsValid() {
		return fmt.Errorf("unknown video task %q (valid: %s)", task, joinVideoTasks())
	}
	info := task.Info()
	if info == nil {
		return fmt.Errorf("unknown video task %q", task)
	}
	if info.RequiresPrompt && strings.TrimSpace(r.Prompt) == "" {
		return fmt.Errorf("prompt is required for video task %q", task)
	}
	if info.RequiresInitImage && r.InitImagePath == "" {
		return fmt.Errorf("--image is required for video task %q", task)
	}
	if info.RequiresInitVideo && r.InitVideoPath == "" {
		return fmt.Errorf("--video is required for video task %q", task)
	}
	if task == VideoTaskScale {
		if r.Scale < 0 {
			return fmt.Errorf("--scale must be > 0 for scale (got %g)", r.Scale)
		}
		if r.Scale == 0 {
			return fmt.Errorf("--scale is required for video task %q", task)
		}
	}
	return nil
}

func ParseCameraControl(raw string) (CameraControl, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	if normalized == "" {
		return CameraControlNone, nil
	}
	switch CameraControl(normalized) {
	case CameraControlZoomIn, CameraControlZoomOut,
		CameraControlPanLeft, CameraControlPanRight,
		CameraControlPanUp, CameraControlPanDown:
		return CameraControl(normalized), nil
	default:
		return "", fmt.Errorf("unknown camera preset %q (valid: zoom_in, zoom_out, pan_left, pan_right, pan_up, pan_down)", raw)
	}
}

func joinVideoTasks() string {
	parts := make([]string, len(AllVideoTasks()))
	for i, t := range AllVideoTasks() {
		parts[i] = string(t)
	}
	return strings.Join(parts, ", ")
}
