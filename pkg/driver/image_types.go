package driver

// ImageCommonParams are shared diffusion parameters for most image tasks.
type ImageCommonParams struct {
	NegativePrompt string
	Model          string
	Width          int
	Height         int
	Steps          int
	Sampler        string
	CFGScale       float32
	BatchSize      int
	BatchCount     int
	Seed           *int
	LoRAs          []LoRARef
	ControlUnits   []ControlNetUnit
	ControlMode    ControlNetMode
	Mode           AssetMode
}

// ImageGenerateRequest is input for text-to-image (GenerateImage RPC).
type ImageGenerateRequest struct {
	Prompt string
	ImageCommonParams
}

// ImageToImageRequest is input for image-to-image (ImageToImage RPC).
type ImageToImageRequest struct {
	Prompt            string
	InitImagePath     string
	DenoisingStrength float32
	ImageCommonParams
}

// ImageInpaintRequest is input for inpainting (InpaintImage RPC).
type ImageInpaintRequest struct {
	Prompt            string
	InitImagePath     string
	MaskImagePath     string
	DenoisingStrength float32
	ImageCommonParams
}

// ImageUpscaleRequest is input for super-resolution (UpscaleImage RPC).
type ImageUpscaleRequest struct {
	InitImagePath string
	Model         string
	Scale         float32
	LoRAs         []LoRARef
}

// ImageDownscaleRequest is input for resolution reduction (DownscaleImage RPC).
type ImageDownscaleRequest struct {
	InitImagePath string
	Model         string
	Scale         float32
	LoRAs         []LoRARef
}

// ImageScaleRequest is input for arbitrary-factor scaling (ScaleImage RPC, e.g. 3×).
type ImageScaleRequest struct {
	InitImagePath string
	Model         string
	Scale         float32
	LoRAs         []LoRARef
}

// ImageEditRequest is input for instruction-based editing (EditImage RPC).
type ImageEditRequest struct {
	Prompt            string
	InitImagePath     string
	DenoisingStrength float32
	ImageCommonParams
}

// ImageControlNetRequest is input for ControlNet generation (ControlNetImage RPC).
type ImageControlNetRequest struct {
	Prompt       string
	ControlUnits []ControlNetUnit
	ControlMode  ControlNetMode
	ImageCommonParams
}

// ImageDepthToImageRequest is input for depth-conditioned generation (DepthToImage RPC).
type ImageDepthToImageRequest struct {
	Prompt         string
	DepthImagePath string
	ImageCommonParams
}

// ImageVariationRequest is input for image variations (ImageVariation RPC).
type ImageVariationRequest struct {
	InitImagePath     string
	DenoisingStrength float32
	ImageCommonParams
}

// ImageStyleTransferRequest is input for style / IP-Adapter transfer (StyleTransferImage RPC).
type ImageStyleTransferRequest struct {
	Prompt         string
	StyleImagePath string
	StyleWeight    float32
	ImageCommonParams
}

// ImageSpriteRequest is input for sprite sheet generation (SpriteImage RPC).
type ImageSpriteRequest struct {
	Prompt             string
	InitImagePath      string
	ReferenceImagePath string
	StyleImagePath     string
	StyleWeight        float32
	Action             SpriteAction
	View               SpriteView
	Directions         int
	Loop               bool
	Padding            int
	Transparent        bool
	FrameWidth         int
	FrameHeight        int
	Columns            int
	Rows               int
	FrameCount         int
	DenoisingStrength  float32
	ImageCommonParams
}

func commonFromImageRequest(r ImageRequest) ImageCommonParams {
	return ImageCommonParams{
		NegativePrompt: r.NegativePrompt,
		Model:          r.Model,
		Width:          r.Width,
		Height:         r.Height,
		Steps:          r.Steps,
		Sampler:        r.Sampler,
		CFGScale:       r.CFGScale,
		BatchSize:      r.BatchSize,
		BatchCount:     r.BatchCount,
		Seed:           r.Seed,
		LoRAs:          r.LoRAs,
		ControlUnits:   append([]ControlNetUnit(nil), r.ControlUnits...),
		ControlMode:    r.ControlMode,
		Mode:           r.Mode,
	}
}

// ToGenerateRequest converts a CLI union request to a generate request.
func (r ImageRequest) ToGenerateRequest() ImageGenerateRequest {
	return ImageGenerateRequest{Prompt: r.Prompt, ImageCommonParams: commonFromImageRequest(r)}
}

func (r ImageRequest) ToImageToImageRequest() ImageToImageRequest {
	return ImageToImageRequest{
		Prompt: r.Prompt, InitImagePath: r.InitImagePath,
		DenoisingStrength: r.DenoisingStrength, ImageCommonParams: commonFromImageRequest(r),
	}
}

func (r ImageRequest) ToInpaintRequest() ImageInpaintRequest {
	return ImageInpaintRequest{
		Prompt: r.Prompt, InitImagePath: r.InitImagePath, MaskImagePath: r.MaskImagePath,
		DenoisingStrength: r.DenoisingStrength, ImageCommonParams: commonFromImageRequest(r),
	}
}

func (r ImageRequest) ToUpscaleRequest() ImageUpscaleRequest {
	return ImageUpscaleRequest{
		InitImagePath: r.InitImagePath, Model: r.Model, Scale: r.Scale, LoRAs: r.LoRAs,
	}
}

func (r ImageRequest) ToEditRequest() ImageEditRequest {
	return ImageEditRequest{
		Prompt: r.Prompt, InitImagePath: r.InitImagePath,
		DenoisingStrength: r.DenoisingStrength, ImageCommonParams: commonFromImageRequest(r),
	}
}

func (r ImageRequest) ToControlNetRequest() ImageControlNetRequest {
	units := append([]ControlNetUnit(nil), r.ControlUnits...)
	if len(units) == 0 && r.ControlImagePath != "" {
		typ := ImageControlDepth
		if r.ControlType != "" {
			typ = r.ControlType
		}
		units = []ControlNetUnit{{Type: typ, ImagePath: r.ControlImagePath}}
	}
	return ImageControlNetRequest{
		Prompt: r.Prompt, ControlUnits: units, ControlMode: r.ControlMode,
		ImageCommonParams: commonFromImageRequest(r),
	}
}

func (r ImageRequest) ToDepthToImageRequest() ImageDepthToImageRequest {
	depthPath := r.ControlImagePath
	for _, u := range r.ControlUnits {
		if u.Type == ImageControlDepth {
			depthPath = u.ImagePath
			break
		}
	}
	return ImageDepthToImageRequest{
		Prompt: r.Prompt, DepthImagePath: depthPath,
		ImageCommonParams: commonFromImageRequest(r),
	}
}

func (r ImageRequest) ToVariationRequest() ImageVariationRequest {
	return ImageVariationRequest{
		InitImagePath:     r.InitImagePath,
		DenoisingStrength: r.DenoisingStrength,
		ImageCommonParams: commonFromImageRequest(r),
	}
}

func (r ImageRequest) ToStyleTransferRequest() ImageStyleTransferRequest {
	return ImageStyleTransferRequest{
		Prompt: r.Prompt, StyleImagePath: r.StyleImagePath, StyleWeight: r.StyleWeight,
		ImageCommonParams: commonFromImageRequest(r),
	}
}

func (r ImageRequest) ToSpriteRequest() ImageSpriteRequest {
	return ImageSpriteRequest{
		Prompt:             r.Prompt,
		InitImagePath:      r.InitImagePath,
		ReferenceImagePath: r.ReferenceImagePath,
		StyleImagePath:     r.StyleImagePath,
		StyleWeight:        r.StyleWeight,
		Action:             r.SpriteAction,
		View:               r.SpriteView,
		Directions:         r.SpriteDirections,
		Loop:               r.SpriteLoop,
		Padding:            r.SpritePadding,
		Transparent:        r.SpriteTransparent,
		FrameWidth:         r.FrameWidth,
		FrameHeight:        r.FrameHeight,
		Columns:            r.Columns,
		Rows:               r.Rows,
		FrameCount:         r.FrameCount,
		DenoisingStrength:  r.DenoisingStrength,
		ImageCommonParams:  commonFromImageRequest(r),
	}
}
