package driver

// VideoCommonParams are shared generation parameters for text/image-to-video tasks.
type VideoCommonParams struct {
	Duration       float32
	FPS            int
	Frames         int
	MotionStrength float32
	ContextLength  int
	Sampler        string
	Scheduler      string
	NegativePrompt string
	Model          string
	Seed           *int
}

// VideoGenerateRequest is input for text-to-video (GenerateVideo RPC).
type VideoGenerateRequest struct {
	Prompt        string
	CameraControl CameraControl
	VideoCommonParams
}

// VideoImageToVideoRequest is input for image-to-video (ImageToVideo RPC).
type VideoImageToVideoRequest struct {
	Prompt        string
	InitImagePath string
	VideoCommonParams
}

// VideoInterpolateRequest is input for frame interpolation (InterpolateVideo RPC).
type VideoInterpolateRequest struct {
	InputVideoPath string
	TargetFPS      int
	Model          string
}

// VideoUpscaleRequest is input for video scaling (UpscaleVideo RPC).
type VideoUpscaleRequest struct {
	InputVideoPath string
	Scale          float32
	Model          string
}

func commonVideoFromRequest(r VideoRequest) VideoCommonParams {
	return VideoCommonParams{
		Duration: r.Duration, FPS: r.FPS, Frames: r.Frames, MotionStrength: r.MotionStrength,
		ContextLength: r.ContextLength, Sampler: r.Sampler, Scheduler: r.Scheduler,
		NegativePrompt: r.NegativePrompt, Model: r.Model, Seed: r.Seed,
	}
}

func (r VideoRequest) ToGenerateRequest() VideoGenerateRequest {
	return VideoGenerateRequest{
		Prompt: r.Prompt, CameraControl: r.CameraControl, VideoCommonParams: commonVideoFromRequest(r),
	}
}

func (r VideoRequest) ToImageToVideoRequest() VideoImageToVideoRequest {
	return VideoImageToVideoRequest{
		Prompt: r.Prompt, InitImagePath: r.InitImagePath, VideoCommonParams: commonVideoFromRequest(r),
	}
}

func (r VideoRequest) ToInterpolateRequest() VideoInterpolateRequest {
	targetFPS := r.FPS
	if targetFPS <= 0 {
		targetFPS = 0
	}
	return VideoInterpolateRequest{
		InputVideoPath: r.InitVideoPath,
		TargetFPS:      targetFPS,
		Model:          r.Model,
	}
}

func (r VideoRequest) ToUpscaleRequest() VideoUpscaleRequest {
	return VideoUpscaleRequest{InputVideoPath: r.InitVideoPath, Scale: r.Scale, Model: r.Model}
}
