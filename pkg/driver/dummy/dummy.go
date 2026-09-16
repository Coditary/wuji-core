package dummy

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/imageresize"
	"github.com/coditary/wuji-core/pkg/media"
	"github.com/coditary/wuji-core/pkg/modelformat"
)

const (
	DriverID   = "dummy"
	DriverName = "Dummy Driver"
)

// Driver is a placeholder backend for development and testing.
type Driver struct{}

func New() *Driver {
	return &Driver{}
}

func (d *Driver) Info() driver.Info {
	return driver.Info{
		ID:          DriverID,
		Name:        DriverName,
		Version:     "0.1.0",
		Description: "A dummy driver that returns placeholder responses for all capabilities.",
		Capabilities: []capability.Type{
			capability.TextGeneration,
			capability.ImageGeneration,
			capability.ImageUpscale,
			capability.ImageDownscale,
			capability.ImageScale,
			capability.VideoGeneration,
			capability.AudioGeneration,
			capability.Mesh,
			capability.Video2Audio,
			capability.VoiceCloning,
			capability.Training,
			capability.DatasetMgmt,
			capability.Data,
			capability.RAG,
		},
		ImageTasks: driver.AllImageTasks(),
		VideoTasks: driver.AllVideoTasks(),
		MeshTasks:  driver.AllMeshTasks(),
		DatasetTasks: driver.AllDatasetTasks(),
		AudioTasks:   driver.AllAudioTasks(),
		VoiceTasks:   driver.AllVoiceTasks(),
		DataTasks:    driver.AllDataTasks(),
		RAGTasks:     driver.AllRAGTasks(),
		FormatSupport: driver.CapabilityFormats{
			capability.TextGeneration: {
				modelformat.GGUF, modelformat.GGML, modelformat.SafeTensors, modelformat.HuggingFace,
				modelformat.Ollama, modelformat.PyTorch, modelformat.Bin, modelformat.ONNX,
				modelformat.MLX, modelformat.AWQ, modelformat.GPTQ, modelformat.LoRA,
			},
			capability.ImageGeneration: {
				modelformat.SafeTensors, modelformat.Diffusers, modelformat.CKPT,
				modelformat.PyTorch, modelformat.ONNX, modelformat.LoRA,
			},
			capability.VideoGeneration: {
				modelformat.SafeTensors, modelformat.Diffusers, modelformat.ONNX,
			},
			capability.AudioGeneration: {
				modelformat.SafeTensors, modelformat.ONNX, modelformat.PyTorch,
			},
			capability.Mesh: {
				modelformat.SafeTensors, modelformat.ONNX,
			},
			capability.VoiceCloning:    {modelformat.SafeTensors, modelformat.LoRA},
			capability.Training:        {modelformat.GGUF, modelformat.HuggingFace, modelformat.SafeTensors, modelformat.LoRA},
			capability.DatasetMgmt:     {},
		},
		Remote: false,
	}
}

func (d *Driver) Capabilities() []capability.Type {
	return d.Info().Capabilities
}

func formatLoras(loras []driver.LoRARef) string {
	if len(loras) == 0 {
		return "none"
	}
	parts := make([]string, 0, len(loras))
	for _, l := range loras {
		parts = append(parts, fmt.Sprintf("%s@%.2f", l.Path, l.LoRAWeight()))
	}
	return strings.Join(parts, ", ")
}

func (d *Driver) GenerateText(_ context.Context, req driver.TextRequest) (*driver.TextResponse, error) {
	meta := fmt.Sprintf(
		"(max_tokens=%d, temperature=%.1f, top_p=%.2f, top_k=%d, min_p=%.2f, freq_pen=%.2f, pres_pen=%.2f, rep_pen=%.2f, stops=%d, seed=%s, ctx=%d, loras=%s, translate=%t, target=%q)",
		req.MaxTokens, req.Temperature,
		req.TopP, req.TopK, req.MinP,
		req.FrequencyPenalty, req.PresencePenalty, req.RepetitionPenalty,
		len(req.StopSequences), formatSeed(req.Seed), req.ContextWindow, formatLoras(req.LoRAs),
		req.Translate, req.TargetLang,
	)
	if len(req.Messages) > 0 {
		msgs := req.Messages
		var b strings.Builder
		b.WriteString("[dummy:messages]")
		for _, m := range msgs {
			b.WriteString(fmt.Sprintf("\n[%s] %q", m.Role, m.Content))
		}
		b.WriteString("\n")
		b.WriteString(meta)
		return &driver.TextResponse{Text: b.String(), TokensUsed: b.Len() / 4, FinishReason: "stop"}, nil
	}
	text := fmt.Sprintf("[dummy:text] %q %s", req.Prompt, meta)
	if req.SystemPrompt != "" {
		text = fmt.Sprintf("[dummy:system] %q\n%s", req.SystemPrompt, text)
	}
	return &driver.TextResponse{Text: text, TokensUsed: len(text) / 4, FinishReason: "stop"}, nil
}

func (d *Driver) ImageToText(_ context.Context, req driver.TextRequest) (*driver.TextResponse, error) {
	prompt := strings.TrimSpace(req.Prompt)
	label := "caption"
	if prompt != "" {
		label = "vqa"
	}
	text := fmt.Sprintf("[dummy:image2text:%s] image=%q prompt=%q", label, req.MediaPath, prompt)
	return &driver.TextResponse{Text: text, FinishReason: "stop"}, nil
}

func (d *Driver) VideoToText(_ context.Context, req driver.TextRequest) (*driver.TextResponse, error) {
	prompt := strings.TrimSpace(req.Prompt)
	text := fmt.Sprintf("[dummy:video2text] video=%q prompt=%q", req.MediaPath, prompt)
	return &driver.TextResponse{Text: text, FinishReason: "stop"}, nil
}

func (d *Driver) DocumentToText(_ context.Context, req driver.TextRequest) (*driver.TextResponse, error) {
	text := fmt.Sprintf("[dummy:document2text] document=%q", req.MediaPath)
	return &driver.TextResponse{Text: text, FinishReason: "stop"}, nil
}

func (d *Driver) AudioToText(_ context.Context, req driver.TextRequest) (*driver.TextResponse, error) {
	text := fmt.Sprintf("[dummy:audio2text] transcribed from %q", req.MediaPath)
	if req.VADEnabled {
		text += fmt.Sprintf(" (trim_silence=%.2f)", req.VADThreshold)
	}
	return &driver.TextResponse{Text: text, FinishReason: "stop"}, nil
}

func (d *Driver) VideoToAudio(_ context.Context, req driver.Video2AudioRequest) (*driver.Video2AudioResponse, error) {
	if strings.TrimSpace(req.VideoPath) == "" {
		return nil, fmt.Errorf("video path is required for video2audio")
	}
	return &driver.Video2AudioResponse{
		AudioPath: fmt.Sprintf("/tmp/wuji/dummy-%s.wav", filepath.Base(req.VideoPath)),
		Temporary: false,
	}, nil
}

func formatSeed(seed *int) string {
	if seed == nil {
		return "random"
	}
	return fmt.Sprintf("%d", *seed)
}

func (d *Driver) GenerateImage(_ context.Context, req driver.ImageGenerateRequest) (*driver.ImageResponse, error) {
	if strings.TrimSpace(req.Prompt) == "" {
		return nil, fmt.Errorf("prompt is required for image task %q", driver.ImageTaskGenerate)
	}
	return dummyImageResponse(driver.ImageTaskGenerate, req.Width, req.Height, req.BatchSize, req.BatchCount, nil, req.LoRAs), nil
}

func (d *Driver) ImageToImage(_ context.Context, req driver.ImageToImageRequest) (*driver.ImageResponse, error) {
	if req.InitImagePath == "" {
		return nil, fmt.Errorf("--init-image is required for image task %q", driver.ImageTaskImg2Img)
	}
	return dummyImageResponse(driver.ImageTaskImg2Img, req.Width, req.Height, req.BatchSize, req.BatchCount,
		[]string{"init=" + req.InitImagePath}, req.LoRAs), nil
}

func (d *Driver) InpaintImage(_ context.Context, req driver.ImageInpaintRequest) (*driver.ImageResponse, error) {
	if req.InitImagePath == "" || req.MaskImagePath == "" {
		return nil, fmt.Errorf("init image and mask are required for inpaint")
	}
	return dummyImageResponse(driver.ImageTaskInpaint, req.Width, req.Height, req.BatchSize, req.BatchCount,
		[]string{"init=" + req.InitImagePath, "mask=" + req.MaskImagePath}, req.LoRAs), nil
}

func (d *Driver) UpscaleImage(_ context.Context, req driver.ImageUpscaleRequest) (*driver.ImageResponse, error) {
	if req.InitImagePath == "" {
		return nil, fmt.Errorf("--init-image is required for image task %q", driver.ImageTaskUpscale)
	}
	meta := []string{"init=" + req.InitImagePath}
	if req.Scale > 0 {
		meta = append(meta, fmt.Sprintf("scale=%g", req.Scale))
	}
	return dummyImageResponse(driver.ImageTaskUpscale, 0, 0, 1, 1, meta, req.LoRAs), nil
}

func (d *Driver) DownscaleImage(ctx context.Context, req driver.ImageDownscaleRequest) (*driver.ImageResponse, error) {
	if req.InitImagePath == "" {
		return nil, fmt.Errorf("--init-image is required for downscale")
	}
	path, err := imageresize.ScaleImage(ctx, req.InitImagePath, "", req.Scale)
	if err != nil {
		return nil, err
	}
	return &driver.ImageResponse{Path: path, Format: "png"}, nil
}

func (d *Driver) ScaleImage(ctx context.Context, req driver.ImageScaleRequest) (*driver.ImageResponse, error) {
	if req.InitImagePath == "" {
		return nil, fmt.Errorf("--init-image is required for scale")
	}
	if req.Scale <= 0 {
		return nil, fmt.Errorf("scale factor must be > 0")
	}
	drivers := driver.ScaleDrivers{Image: d, Upscale: d, Downscale: d}
	resp, err := driver.RunBuiltInHybridScale(ctx, drivers, driver.ImageRequest{
		Task:          driver.ImageTaskUpscale,
		InitImagePath: req.InitImagePath,
		Model:         req.Model,
		Scale:         req.Scale,
		LoRAs:         req.LoRAs,
	}, req.Scale)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *Driver) EditImage(_ context.Context, req driver.ImageEditRequest) (*driver.ImageResponse, error) {
	if req.InitImagePath == "" || strings.TrimSpace(req.Prompt) == "" {
		return nil, fmt.Errorf("prompt and init image are required for edit")
	}
	return dummyImageResponse(driver.ImageTaskEdit, req.Width, req.Height, req.BatchSize, req.BatchCount,
		[]string{"init=" + req.InitImagePath}, req.LoRAs), nil
}

func (d *Driver) ControlNetImage(_ context.Context, req driver.ImageControlNetRequest) (*driver.ImageResponse, error) {
	units := req.ControlUnits
	if len(units) == 0 {
		units = req.ImageCommonParams.ControlUnits
	}
	if len(units) == 0 {
		return nil, fmt.Errorf("at least one control flag (e.g. --pose, --depth) is required for controlnet")
	}
	meta := make([]string, 0, len(units)*2+1)
	meta = append(meta, "control_mode="+string(req.ControlMode))
	for _, u := range units {
		meta = append(meta, fmt.Sprintf("control=%s:%s", u.Type, u.ImagePath))
		if u.Preprocessor != "" {
			meta = append(meta, fmt.Sprintf("preprocessor=%s:%s", u.Type, u.Preprocessor))
		}
	}
	return dummyImageResponse(driver.ImageTaskControlNet, req.Width, req.Height, req.BatchSize, req.BatchCount, meta, req.LoRAs), nil
}

func (d *Driver) DepthToImage(_ context.Context, req driver.ImageDepthToImageRequest) (*driver.ImageResponse, error) {
	if req.DepthImagePath == "" {
		return nil, fmt.Errorf("depth image is required for depth2img")
	}
	return dummyImageResponse(driver.ImageTaskDepth2Img, req.Width, req.Height, req.BatchSize, req.BatchCount,
		[]string{"depth=" + req.DepthImagePath}, req.LoRAs), nil
}

func (d *Driver) ImageVariation(_ context.Context, req driver.ImageVariationRequest) (*driver.ImageResponse, error) {
	if req.InitImagePath == "" {
		return nil, fmt.Errorf("--init-image is required for image task %q", driver.ImageTaskVariation)
	}
	meta := []string{"init=" + req.InitImagePath}
	for _, u := range req.ControlUnits {
		meta = append(meta, fmt.Sprintf("control=%s:%s", u.Type, u.ImagePath))
	}
	return dummyImageResponse(driver.ImageTaskVariation, req.Width, req.Height, req.BatchSize, req.BatchCount, meta, req.LoRAs), nil
}

func (d *Driver) StyleTransferImage(_ context.Context, req driver.ImageStyleTransferRequest) (*driver.ImageResponse, error) {
	if req.StyleImagePath == "" || strings.TrimSpace(req.Prompt) == "" {
		return nil, fmt.Errorf("prompt and style image are required for style-transfer")
	}
	return dummyImageResponse(driver.ImageTaskStyleTransfer, req.Width, req.Height, req.BatchSize, req.BatchCount,
		[]string{"style=" + req.StyleImagePath}, req.LoRAs), nil
}

func (d *Driver) SpriteImage(_ context.Context, req driver.ImageSpriteRequest) (*driver.ImageResponse, error) {
	if strings.TrimSpace(req.Prompt) == "" {
		return nil, fmt.Errorf("prompt is required for image task %q", driver.ImageTaskSprite)
	}
	cols, rows := req.Columns, req.Rows
	if cols <= 0 {
		cols = 1
	}
	if rows <= 0 {
		rows = 1
	}
	meta := []string{
		fmt.Sprintf("frame=%dx%d", req.FrameWidth, req.FrameHeight),
		fmt.Sprintf("grid=%dx%d", cols, rows),
		fmt.Sprintf("mode=%s", req.ImageCommonParams.ModeOrDefault()),
	}
	if req.InitImagePath != "" {
		meta = append(meta, "image="+req.InitImagePath)
	}
	if req.ReferenceImagePath != "" {
		meta = append(meta, "reference="+req.ReferenceImagePath)
	}
	if req.StyleImagePath != "" {
		meta = append(meta, "style="+req.StyleImagePath)
	}
	if req.Action != "" {
		meta = append(meta, "action="+string(req.Action))
	}
	if req.View != "" {
		meta = append(meta, "view="+string(req.View))
	}
	if req.Directions > 0 {
		meta = append(meta, fmt.Sprintf("directions=%d", req.Directions))
	}
	if req.Loop {
		meta = append(meta, "loop=true")
	}
	if req.Padding > 0 {
		meta = append(meta, fmt.Sprintf("padding=%d", req.Padding))
	}
	if req.Transparent {
		meta = append(meta, "transparent=true")
	}
	for _, u := range req.ControlUnits {
		meta = append(meta, fmt.Sprintf("control=%s:%s", u.Type, u.ImagePath))
	}
	return dummyImageResponse(driver.ImageTaskSprite, req.FrameWidth*cols, req.FrameHeight*rows, 1, 1, meta, req.LoRAs), nil
}

func dummyImageResponse(task driver.ImageTask, width, height, batchSize, batchCount int, extra []string, loras []driver.LoRARef) *driver.ImageResponse {
	if batchCount <= 0 {
		batchCount = 1
	}
	if batchSize <= 0 {
		batchSize = 1
	}
	if width <= 0 {
		width = 512
	}
	if height <= 0 {
		height = 512
	}

	paths := make([]string, 0, batchCount*batchSize)
	for b := 0; b < batchCount; b++ {
		for i := 0; i < batchSize; i++ {
			idx := b*batchSize + i
			paths = append(paths, fmt.Sprintf("/tmp/wuji/dummy-image-%s-%dx%d-%d.png", task, width, height, idx))
		}
	}

	meta := []string{fmt.Sprintf("task=%s", task), fmt.Sprintf("loras=%s", formatLoras(loras))}
	meta = append(meta, extra...)

	return &driver.ImageResponse{
		Path: paths[0], Paths: paths, Format: fmt.Sprintf("png (%s)", strings.Join(meta, ", ")),
	}
}

func (d *Driver) GenerateVideo(_ context.Context, req driver.VideoGenerateRequest) (*driver.VideoResponse, error) {
	return dummyVideoResponse(req.Duration, req.FPS, req.Frames)
}

func (d *Driver) ImageToVideo(_ context.Context, req driver.VideoImageToVideoRequest) (*driver.VideoResponse, error) {
	return dummyVideoResponse(req.Duration, req.FPS, req.Frames)
}

func (d *Driver) InterpolateVideo(_ context.Context, req driver.VideoInterpolateRequest) (*driver.VideoResponse, error) {
	fps := req.TargetFPS
	if fps <= 0 {
		fps = 48
	}
	return dummyVideoResponse(0, fps, fps*5)
}

func (d *Driver) UpscaleVideo(_ context.Context, req driver.VideoUpscaleRequest) (*driver.VideoResponse, error) {
	scale := req.Scale
	if scale <= 0 {
		scale = 2
	}
	return &driver.VideoResponse{
		Path:     fmt.Sprintf("/tmp/wuji/dummy-video-scale-x%g.mp4", scale),
		Duration: 5,
		Frames:   120,
		FPS:      24,
	}, nil
}

func dummyVideoResponse(duration float32, fps, frames int) (*driver.VideoResponse, error) {
	if fps <= 0 {
		fps = 24
	}
	if frames <= 0 && duration > 0 {
		frames = int(duration * float32(fps))
	}
	if frames <= 0 {
		frames = fps * 5
	}
	if duration <= 0 {
		duration = float32(frames) / float32(fps)
	}
	return &driver.VideoResponse{
		Path:     fmt.Sprintf("/tmp/wuji/dummy-video-%dfr-%dfps.mp4", frames, fps),
		Duration: duration,
		Frames:   frames,
		FPS:      fps,
	}, nil
}

func (d *Driver) GenerateMusic(_ context.Context, req driver.AudioMusicRequest) (*driver.AudioResponse, error) {
	return dummyAudioResponse(req.Duration, req.SampleRate, "music")
}

func (d *Driver) MelodyToMusic(_ context.Context, req driver.AudioMelodyToMusicRequest) (*driver.AudioResponse, error) {
	return dummyAudioResponse(req.Duration, req.SampleRate, "melody")
}

func (d *Driver) GenerateSFX(_ context.Context, req driver.AudioSFXRequest) (*driver.AudioResponse, error) {
	return dummyAudioResponse(req.Duration, req.SampleRate, "sfx")
}

func (d *Driver) GenerateAudioSpeech(_ context.Context, req driver.AudioSpeechRequest) (*driver.AudioResponse, error) {
	duration := req.Duration
	if duration <= 0 {
		duration = estimateSpeechDuration(req.Prompt, req.Speed)
	}
	return dummyAudioResponse(duration, req.SampleRate, "speech")
}

func estimateSpeechDuration(prompt string, speed float32) float32 {
	words := len(strings.Fields(strings.TrimSpace(prompt)))
	if words == 0 {
		return 2
	}
	rate := speed
	if rate <= 0 {
		rate = 1
	}
	seconds := float32(words) / (2.5 * rate)
	if seconds < 1 {
		return 1
	}
	if seconds > 30 {
		return 30
	}
	return seconds
}

func dummyAudioResponse(duration float32, sampleRate int, kind string) (*driver.AudioResponse, error) {
	if duration <= 0 {
		duration = 10
	}
	if sampleRate <= 0 {
		sampleRate = 44100
	}
	if err := os.MkdirAll("/tmp/wuji", 0o755); err != nil {
		return nil, err
	}
	f, err := os.CreateTemp("/tmp/wuji", fmt.Sprintf("dummy-audio-%s-*.wav", kind))
	if err != nil {
		return nil, err
	}
	path := f.Name()
	_ = f.Close()
	if err := media.WritePlaceholderWAV(path, duration, sampleRate); err != nil {
		_ = os.Remove(path)
		return nil, err
	}
	return &driver.AudioResponse{
		Path:       path,
		Duration:   duration,
		SampleRate: sampleRate,
		Format:     "wav",
	}, nil
}

func (d *Driver) GenerateMesh(_ context.Context, req driver.MeshRequest) (*driver.MeshResponse, error) {
	format := req.Format
	if format == "" {
		format = "glb"
	}
	task := req.TaskOrDefault()
	mode := req.ModeOrDefault()
	rep := req.RepresentationOrDefault()
	if rep != driver.MeshRepresentationMesh {
		switch rep {
		case driver.MeshRepresentationNeRF:
			format = "nerf"
		case driver.MeshRepresentationSplat:
			format = "splat"
		case driver.MeshRepresentationPointCloud:
			format = "ply"
		}
	}
	return &driver.MeshResponse{
		Path:           fmt.Sprintf("/tmp/wuji/dummy-mesh-%s-%s-%s.%s", task, mode, rep, format),
		Format:         format,
		Task:           task,
		Mode:           mode,
		Representation: rep,
	}, nil
}

func (d *Driver) ListDatasets(_ context.Context, _ driver.DatasetListRequest) (*driver.DatasetResponse, error) {
	return &driver.DatasetResponse{
		Datasets: []driver.DatasetEntry{
			{ID: "ds-1", Name: "example", Path: "/data/example", Size: 1024, LatestVersion: "v1", VersionCount: 1},
		},
	}, nil
}

func (d *Driver) CreateDataset(_ context.Context, req driver.DatasetCreateRequest) (*driver.DatasetResponse, error) {
	return &driver.DatasetResponse{
		Message: fmt.Sprintf("created dataset %q at %s", req.Name, req.Path),
		Datasets: []driver.DatasetEntry{
			{ID: "ds-new", Name: req.Name, Path: req.Path, LatestVersion: "v1", VersionCount: 1},
		},
	}, nil
}

func (d *Driver) DeleteDataset(_ context.Context, req driver.DatasetDeleteRequest) (*driver.DatasetResponse, error) {
	name := req.Name
	if name == "" {
		name = req.ID
	}
	return &driver.DatasetResponse{Message: fmt.Sprintf("deleted dataset %q", name)}, nil
}

func (d *Driver) IngestDataset(_ context.Context, req driver.DatasetIngestRequest) (*driver.DatasetIngestResult, error) {
	ref := req.Name
	if ref == "" {
		ref = req.DatasetID
	}
	return &driver.DatasetIngestResult{
		Message: fmt.Sprintf("ingested %q into dataset %q", req.SourcePath, ref),
		FilesAdded: 1, BytesAdded: 4096, VersionID: "v2",
	}, nil
}

func (d *Driver) CreateDatasetVersion(_ context.Context, req driver.DatasetVersionRequest) (*driver.DatasetVersionResult, error) {
	ref := req.Name
	if ref == "" {
		ref = req.DatasetID
	}
	tag := req.Tag
	if tag == "" {
		tag = "v2"
	}
	return &driver.DatasetVersionResult{
		Version: driver.DatasetVersionEntry{
			ID: "ver-2", DatasetID: ref, Tag: tag, Path: fmt.Sprintf("/data/%s/%s", ref, tag),
			Size: 4096, CreatedAtUnix: 1_700_000_000, Message: req.Message,
		},
		Message: fmt.Sprintf("created version %q for dataset %q", tag, ref),
	}, nil
}

func (d *Driver) ListDatasetVersions(_ context.Context, req driver.DatasetListVersionsRequest) ([]driver.DatasetVersionEntry, error) {
	ref := req.Name
	if ref == "" {
		ref = req.DatasetID
	}
	return []driver.DatasetVersionEntry{
		{ID: "ver-1", DatasetID: ref, Tag: "v1", Path: fmt.Sprintf("/data/%s/v1", ref), Size: 1024, CreatedAtUnix: 1_690_000_000},
	}, nil
}

func (d *Driver) CreateVoice(_ context.Context, req driver.VoiceCreateRequest) (*driver.VoiceResponse, error) {
	return &driver.VoiceResponse{
		VoiceID: fmt.Sprintf("voice-%s", req.Name),
		Name:    req.Name,
	}, nil
}

func (d *Driver) ConvertVoice(_ context.Context, req driver.VoiceConvertRequest) (*driver.VoiceResponse, error) {
	return &driver.VoiceResponse{
		VoiceID:    fmt.Sprintf("voice-%s", req.Name),
		Name:       req.Name,
		OutputPath: fmt.Sprintf("/tmp/wuji/dummy-vc-%s.wav", req.Name),
	}, nil
}

func (d *Driver) ListVoiceProfiles(_ context.Context) ([]driver.VoiceProfileInfo, error) {
	return []driver.VoiceProfileInfo{
		{
			Name: "narrator", VoiceID: "voice-narrator", Source: "builtin",
			Description: "Built-in neutral narrator for generic text-to-speech",
			Speech: &driver.VoiceSpeechCapabilities{
				Languages: []string{"de", "en", "auto"}, LanguagesSupported: true,
				Emotions: []string{"neutral", "happy", "sad"}, EmotionsSupported: true,
				Styles: []string{"conversational", "news"}, StylesSupported: true,
				SpeedMin: 0.5, SpeedMax: 2.0, PitchMin: -6, PitchMax: 6, EnergySupported: true,
			},
		},
		{
			Name: "demo", VoiceID: "voice-demo", Source: "cloned",
			Description: "Example cloned voice profile",
			Speech: &driver.VoiceSpeechCapabilities{
				Languages: []string{"de", "en"}, LanguagesSupported: true,
				Emotions: []string{"neutral", "happy"}, EmotionsSupported: true,
				Styles: []string{"conversational"}, StylesSupported: true,
				SpeedMin: 0.75, SpeedMax: 1.5, PitchMin: -12, PitchMax: 12, EnergySupported: false,
			},
		},
	}, nil
}

func (d *Driver) makeTrainResp(cap capability.Type, name, dataset, suffix string, epochs int, output string) *driver.TrainResponse {
	if output == "" && name != "" {
		output = fmt.Sprintf("/tmp/wuji/trained-%s-%s", cap, name)
	}
	return &driver.TrainResponse{
		JobID:      fmt.Sprintf("job-%s-%s-%s", cap, dataset, suffix),
		Status:     fmt.Sprintf("queued %s training (%d epochs)", cap, epochs),
		Capability: cap,
		OutputPath: output,
	}
}

func (d *Driver) TrainText(_ context.Context, req driver.TextTrainRequest) (*driver.TrainResponse, error) {
	return d.makeTrainResp(capability.TextGeneration, req.Name, req.DatasetID, req.BaseModel, req.Epochs, req.OutputPath), nil
}

func (d *Driver) TrainImage(_ context.Context, req driver.ImageTrainRequest) (*driver.TrainResponse, error) {
	return d.makeTrainResp(capability.ImageGeneration, req.Name, req.DatasetID, req.BaseModel, req.Epochs, req.OutputPath), nil
}

func (d *Driver) TrainVideo(_ context.Context, req driver.VideoTrainRequest) (*driver.TrainResponse, error) {
	return d.makeTrainResp(capability.VideoGeneration, req.Name, req.DatasetID, req.BaseModel, req.Epochs, req.OutputPath), nil
}

func (d *Driver) TrainAudio(_ context.Context, req driver.AudioTrainRequest) (*driver.TrainResponse, error) {
	return d.makeTrainResp(capability.AudioGeneration, req.Name, req.DatasetID, req.BaseModel, req.Epochs, req.OutputPath), nil
}

func (d *Driver) TrainMesh(_ context.Context, req driver.MeshTrainRequest) (*driver.TrainResponse, error) {
	return d.makeTrainResp(capability.Mesh, req.Name, req.DatasetID, req.BaseModel, req.Epochs, req.OutputPath), nil
}

func (d *Driver) TrainVoice(_ context.Context, req driver.VoiceTrainRequest) (*driver.TrainResponse, error) {
	suffix := req.PretrainedModel
	if suffix == "" {
		suffix = "rvc"
	}
	return d.makeTrainResp(capability.VoiceCloning, req.Name, req.DatasetID, suffix, req.Epochs, req.OutputPath), nil
}

func (d *Driver) Close() error {
	return nil
}
