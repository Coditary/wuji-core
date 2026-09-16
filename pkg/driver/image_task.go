package driver

import (
	"fmt"
	"strings"
)

// ImageTask identifies the image generation or transformation mode.
type ImageTask string

const (
	ImageTaskGenerate      ImageTask = "generate"
	ImageTaskImg2Img       ImageTask = "img2img"
	ImageTaskInpaint       ImageTask = "inpaint"
	ImageTaskUpscale       ImageTask = "upscale"
	ImageTaskEdit          ImageTask = "edit"
	ImageTaskControlNet    ImageTask = "controlnet"
	ImageTaskDepth2Img     ImageTask = "depth2img"
	ImageTaskVariation     ImageTask = "variation"
	ImageTaskStyleTransfer ImageTask = "style-transfer"
	ImageTaskSprite        ImageTask = "sprite"
)

// ImageControlType identifies a ControlNet conditioning mode.
type ImageControlType string

const (
	ImageControlCanny        ImageControlType = "canny"
	ImageControlDepth        ImageControlType = "depth"
	ImageControlPose         ImageControlType = "pose"
	ImageControlNormal       ImageControlType = "normal"
	ImageControlSegmentation ImageControlType = "segmentation"
	ImageControlScribble     ImageControlType = "scribble"
	ImageControlLineart      ImageControlType = "lineart"
	ImageControlLineartAnime ImageControlType = "lineart-anime"
	ImageControlMLSD         ImageControlType = "mlsd"
	ImageControlTile         ImageControlType = "tile"
)

// ImageTaskInfo describes a task and its required inputs for CLI help and validation.
type ImageTaskInfo struct {
	Task                 ImageTask
	Description          string
	RequiresPrompt       bool
	RequiresInitImage    bool
	RequiresMask         bool
	RequiresControlImage bool
	RequiresControlType  bool
	RequiresStyleImage   bool
}

// AllImageTasks returns every supported image task in display order.
func AllImageTasks() []ImageTask {
	return []ImageTask{
		ImageTaskGenerate,
		ImageTaskImg2Img,
		ImageTaskInpaint,
		ImageTaskUpscale,
		ImageTaskEdit,
		ImageTaskControlNet,
		ImageTaskDepth2Img,
		ImageTaskVariation,
		ImageTaskStyleTransfer,
		ImageTaskSprite,
	}
}

// AllImageControlTypes returns every supported ControlNet control type.
func AllImageControlTypes() []ImageControlType {
	return []ImageControlType{
		ImageControlCanny,
		ImageControlDepth,
		ImageControlPose,
		ImageControlNormal,
		ImageControlSegmentation,
		ImageControlScribble,
		ImageControlLineart,
		ImageControlLineartAnime,
		ImageControlMLSD,
		ImageControlTile,
	}
}

// ImageTaskCatalog returns metadata for all image tasks.
func ImageTaskCatalog() []ImageTaskInfo {
	return []ImageTaskInfo{
		{Task: ImageTaskGenerate, Description: "Text-to-image generation from a prompt", RequiresPrompt: true},
		{Task: ImageTaskImg2Img, Description: "Transform an image guided by a prompt", RequiresPrompt: true, RequiresInitImage: true},
		{Task: ImageTaskInpaint, Description: "Fill masked regions in an image", RequiresPrompt: true, RequiresInitImage: true, RequiresMask: true},
		{Task: ImageTaskUpscale, Description: "Scale image up or down (coordinates upscale/downscale capabilities)", RequiresInitImage: true},
		{Task: ImageTaskEdit, Description: "Edit an image from an instruction prompt (InstructPix2Pix-style)", RequiresPrompt: true, RequiresInitImage: true},
		{Task: ImageTaskControlNet, Description: "Generate with structural control (edges, pose, depth, …)", RequiresPrompt: true, RequiresControlImage: true, RequiresControlType: true},
		{Task: ImageTaskDepth2Img, Description: "Generate from a depth map", RequiresPrompt: true, RequiresControlImage: true},
		{Task: ImageTaskVariation, Description: "Create variations of an input image", RequiresInitImage: true},
		{Task: ImageTaskStyleTransfer, Description: "Apply a reference style (IP-Adapter / style image)", RequiresPrompt: true, RequiresStyleImage: true},
		{Task: ImageTaskSprite, Description: "Generate a sprite sheet atlas from a prompt and references", RequiresPrompt: true},
	}
}

func (t ImageTask) String() string {
	return string(t)
}

// IsValid reports whether the task is a known image task.
func (t ImageTask) IsValid() bool {
	for _, known := range AllImageTasks() {
		if t == known {
			return true
		}
	}
	return false
}

// Info returns catalog metadata for the task, or nil if unknown.
func (t ImageTask) Info() *ImageTaskInfo {
	for _, info := range ImageTaskCatalog() {
		if info.Task == t {
			copy := info
			return &copy
		}
	}
	return nil
}

func (t ImageControlType) String() string {
	return string(t)
}

// IsValid reports whether the control type is known.
func (t ImageControlType) IsValid() bool {
	for _, known := range AllImageControlTypes() {
		if t == known {
			return true
		}
	}
	return false
}

// ImageTaskInputs holds the CLI fields used to infer an image task implicitly.
type ImageTaskInputs struct {
	Prompt           string
	InitImagePath    string
	MaskImagePath    string
	ControlImagePath string
	ControlType      string
	ControlUnits     []ControlNetUnit
	StyleImagePath   string
	Scale            float32
	UpscaleRequested bool
	EditRequested    bool
	SpriteRequested  bool
}

// InferImageTask selects the image task from provided inputs when --task is omitted.
// More specific signals win (--mask → inpaint before --init-image → img2img).
func InferImageTask(in ImageTaskInputs) ImageTask {
	if in.MaskImagePath != "" {
		return ImageTaskInpaint
	}
	if in.SpriteRequested {
		return ImageTaskSprite
	}
	if in.EditRequested && in.InitImagePath != "" {
		return ImageTaskEdit
	}
	if in.StyleImagePath != "" {
		return ImageTaskStyleTransfer
	}
	if len(in.ControlUnits) > 0 && in.InitImagePath == "" {
		if len(in.ControlUnits) == 1 && in.ControlUnits[0].Type == ImageControlDepth {
			return ImageTaskDepth2Img
		}
		return ImageTaskControlNet
	}
	if in.ControlImagePath != "" {
		if in.ControlType != "" {
			return ImageTaskControlNet
		}
		return ImageTaskDepth2Img
	}
	if in.InitImagePath != "" {
		if in.UpscaleRequested || (in.Scale != 0 && strings.TrimSpace(in.Prompt) == "") {
			return ImageTaskUpscale
		}
		if strings.TrimSpace(in.Prompt) == "" {
			return ImageTaskVariation
		}
		return ImageTaskImg2Img
	}
	return ImageTaskGenerate
}

// ParseImageTask normalizes and validates a task name from CLI or config.
func ParseImageTask(raw string) (ImageTask, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	if normalized == "" || normalized == "txt2img" || normalized == "text-to-image" {
		return ImageTaskGenerate, nil
	}
	if normalized == "image-to-image" {
		return ImageTaskImg2Img, nil
	}
	if normalized == "style" || normalized == "ip-adapter" {
		return ImageTaskStyleTransfer, nil
	}
	if normalized == "spritesheet" || normalized == "sprite-sheet" {
		return ImageTaskSprite, nil
	}

	task := ImageTask(normalized)
	if !task.IsValid() {
		return "", fmt.Errorf("unknown image task %q (valid: %s)", raw, joinImageTasks())
	}
	return task, nil
}

// TaskOrDefault returns generate when Task is empty.
func (r ImageRequest) TaskOrDefault() ImageTask {
	if r.Task == "" {
		return ImageTaskGenerate
	}
	return r.Task
}

// Validate checks task-specific required fields.
func (r ImageRequest) Validate() error {
	task := r.TaskOrDefault()
	if !task.IsValid() {
		return fmt.Errorf("unknown image task %q (valid: %s)", task, joinImageTasks())
	}
	info := task.Info()
	if info == nil {
		return fmt.Errorf("unknown image task %q", task)
	}

	if info.RequiresPrompt && strings.TrimSpace(r.Prompt) == "" {
		return fmt.Errorf("prompt is required for image task %q", task)
	}
	if info.RequiresInitImage && r.InitImagePath == "" {
		return fmt.Errorf("--init-image is required for image task %q", task)
	}
	if info.RequiresMask && r.MaskImagePath == "" {
		return fmt.Errorf("--mask is required for image task %q", task)
	}
	if info.RequiresControlImage && r.ControlImagePath == "" && len(r.ControlUnits) == 0 {
		return fmt.Errorf("a control flag (e.g. --pose, --depth) is required for image task %q", task)
	}
	if info.RequiresControlType {
		if r.ControlType == "" && !hasControlUnits(r.ControlUnits) {
			return fmt.Errorf("--control-type or a typed control flag is required for image task %q", task)
		}
		if r.ControlType != "" && !r.ControlType.IsValid() {
			return fmt.Errorf("unknown control type %q", r.ControlType)
		}
	}
	if info.RequiresStyleImage && r.StyleImagePath == "" {
		return fmt.Errorf("--style is required for image task %q", task)
	}
	if err := validateControlUnits(r.ControlUnits); err != nil {
		return err
	}
	if r.ControlMode != "" {
		if _, err := ParseControlNetMode(string(r.ControlMode)); err != nil {
			return err
		}
	}

	if r.Mode != "" {
		if _, err := ParseAssetMode(string(r.Mode)); err != nil {
			return err
		}
	}

	switch task {
	case ImageTaskControlNet:
		if len(r.ControlUnits) == 0 && r.ControlImagePath == "" {
			return fmt.Errorf("at least one control flag (e.g. --pose, --depth) is required for image task %q", task)
		}
	case ImageTaskDepth2Img:
		if len(r.ControlUnits) == 0 && r.ControlImagePath == "" {
			return fmt.Errorf("--depth is required for image task %q", task)
		}
	case ImageTaskImg2Img, ImageTaskEdit, ImageTaskVariation, ImageTaskInpaint, ImageTaskGenerate, ImageTaskStyleTransfer:
		// control units are optional attachments
	case ImageTaskSprite:
		if r.FrameWidth <= 0 || r.FrameHeight <= 0 {
			return fmt.Errorf("--frame-width and --frame-height (or -W/-H) are required for image task %q", task)
		}
		if spriteSheetFrameTotal(r) <= 0 {
			return fmt.Errorf("set --columns/--rows, --frames, or --batch for image task %q", task)
		}
		if r.SpriteAction != "" {
			if _, err := ParseSpriteAction(string(r.SpriteAction)); err != nil {
				return err
			}
		}
		if r.SpriteView != "" {
			if _, err := ParseSpriteView(string(r.SpriteView)); err != nil {
				return err
			}
		}
		if _, err := ParseSpriteDirections(r.SpriteDirections); err != nil {
			return err
		}
		if r.SpritePadding < 0 {
			return fmt.Errorf("--padding must be >= 0 for image task %q", task)
		}
	case ImageTaskUpscale:
		if r.Scale < 0 {
			return fmt.Errorf("--scale must be > 0 for upscale (got %g)", r.Scale)
		}
		if r.Scale == 0 {
			return fmt.Errorf("--scale is required (or use --upscale for 4× default)")
		}
	}

	return nil
}

func spriteSheetFrameTotal(r ImageRequest) int {
	if r.FrameCount > 0 {
		return r.FrameCount
	}
	if r.Columns > 0 && r.Rows > 0 {
		return r.Columns * r.Rows
	}
	if r.BatchSize > 0 && r.BatchCount > 0 {
		return r.BatchSize * r.BatchCount
	}
	if r.BatchSize > 0 {
		return r.BatchSize
	}
	return 0
}

func joinImageTasks() string {
	parts := make([]string, len(AllImageTasks()))
	for i, t := range AllImageTasks() {
		parts[i] = string(t)
	}
	return strings.Join(parts, ", ")
}

func joinImageControlTypes() string {
	parts := make([]string, len(AllImageControlTypes()))
	for i, t := range AllImageControlTypes() {
		parts[i] = string(t)
	}
	return strings.Join(parts, ", ")
}

func hasControlUnits(units []ControlNetUnit) bool {
	for _, u := range units {
		if u.ImagePath != "" {
			return true
		}
	}
	return false
}
