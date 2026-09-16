package driver

import wujiv1 "github.com/coditary/wuji-core/api/proto/v1"

func ImageTasksToProto(tasks []ImageTask) []string {
	return imageTasksToProto(tasks)
}

func ImageTasksFromProto(items []string) []ImageTask {
	return imageTasksFromProto(items)
}

func imageTasksToProto(tasks []ImageTask) []string {
	if len(tasks) == 0 {
		return nil
	}
	out := make([]string, 0, len(tasks))
	for _, task := range tasks {
		out = append(out, string(task))
	}
	return out
}

func imageTasksFromProto(items []string) []ImageTask {
	if len(items) == 0 {
		return nil
	}
	out := make([]ImageTask, 0, len(items))
	for _, item := range items {
		task := ImageTask(item)
		if task.IsValid() {
			out = append(out, task)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func commonParamsToProto(p ImageCommonParams) (int32, int32, int32, string, string, string, float32, int32, int32, *int32, []*wujiv1.LoraAttachment, []*wujiv1.ControlNetUnit, string) {
	var seed *int32
	if p.Seed != nil {
		v := int32(*p.Seed)
		seed = &v
	}
	return int32(p.Width), int32(p.Height), int32(p.Steps), p.NegativePrompt, p.Model, p.Sampler,
		p.CFGScale, int32(p.BatchSize), int32(p.BatchCount), seed, lorasToProto(p.LoRAs),
		controlUnitsToProto(p.ControlUnits), string(p.ControlMode)
}

func commonParamsFromProto(width, height, steps int32, negativePrompt, model, sampler string, cfgScale float32, batchSize, batchCount int32, seed *int32, loras []*wujiv1.LoraAttachment, units []*wujiv1.ControlNetUnit, controlMode, mode string) ImageCommonParams {
	out := ImageCommonParams{
		NegativePrompt: negativePrompt,
		Model:          model,
		Width:          int(width),
		Height:         int(height),
		Steps:          int(steps),
		Sampler:        sampler,
		CFGScale:       cfgScale,
		BatchSize:      int(batchSize),
		BatchCount:     int(batchCount),
		LoRAs:          lorasFromProto(loras),
		ControlUnits:   controlUnitsFromProto(units),
		ControlMode:    ControlNetMode(controlMode),
		Mode:           AssetMode(mode),
	}
	if seed != nil {
		v := int(*seed)
		out.Seed = &v
	}
	return out
}

func controlUnitsToProto(units []ControlNetUnit) []*wujiv1.ControlNetUnit {
	if len(units) == 0 {
		return nil
	}
	out := make([]*wujiv1.ControlNetUnit, 0, len(units))
	for _, u := range units {
		out = append(out, &wujiv1.ControlNetUnit{
			Type: u.Type.String(), ImagePath: u.ImagePath, Preprocessor: u.Preprocessor,
			Weight: u.Weight, GuidanceStart: u.GuidanceStart, GuidanceEnd: u.GuidanceEnd,
			ThresholdA: u.ThresholdA, ThresholdB: u.ThresholdB, Model: u.Model,
		})
	}
	return out
}

func controlUnitsFromProto(units []*wujiv1.ControlNetUnit) []ControlNetUnit {
	if len(units) == 0 {
		return nil
	}
	out := make([]ControlNetUnit, 0, len(units))
	for _, u := range units {
		if u == nil {
			continue
		}
		typ, err := ParseImageControlType(u.GetType())
		if err != nil {
			typ = ImageControlType(u.GetType())
		}
		out = append(out, ControlNetUnit{
			Type: typ, ImagePath: u.GetImagePath(), Preprocessor: u.GetPreprocessor(),
			Weight: u.GetWeight(), GuidanceStart: u.GetGuidanceStart(), GuidanceEnd: u.GetGuidanceEnd(),
			ThresholdA: u.GetThresholdA(), ThresholdB: u.GetThresholdB(), Model: u.GetModel(),
		})
	}
	return out
}

func ImageGenerateRequestFromProto(req *wujiv1.GenerateImageRequest) ImageGenerateRequest {
	if req == nil {
		return ImageGenerateRequest{}
	}
	return ImageGenerateRequest{
		Prompt: req.GetPrompt(),
		ImageCommonParams: commonParamsFromProto(
			req.GetWidth(), req.GetHeight(), req.GetSteps(), req.GetNegativePrompt(), req.GetModel(),
			req.GetSampler(), req.GetCfgScale(), req.GetBatchSize(), req.GetBatchCount(), req.Seed, req.GetLoras(),
			req.GetControlUnits(), req.GetControlMode(), req.GetMode(),
		),
	}
}

func ImageGenerateRequestToProto(req ImageGenerateRequest) *wujiv1.GenerateImageRequest {
	w, h, steps, neg, model, sampler, cfg, bs, bc, seed, loras, units, controlMode := commonParamsToProto(req.ImageCommonParams)
	return &wujiv1.GenerateImageRequest{
		Prompt: req.Prompt, Width: w, Height: h, Steps: steps, NegativePrompt: neg, Model: model,
		Sampler: sampler, CfgScale: cfg, BatchSize: bs, BatchCount: bc, Seed: seed, Loras: loras,
		ControlUnits: units, ControlMode: controlMode,
		Mode: string(req.ImageCommonParams.ModeOrDefault()),
	}
}

func ImageToImageRequestFromProto(req *wujiv1.ImageToImageRequest) ImageToImageRequest {
	if req == nil {
		return ImageToImageRequest{}
	}
	return ImageToImageRequest{
		Prompt: req.GetPrompt(), InitImagePath: req.GetInitImagePath(), DenoisingStrength: req.GetDenoisingStrength(),
		ImageCommonParams: commonParamsFromProto(
			req.GetWidth(), req.GetHeight(), req.GetSteps(), req.GetNegativePrompt(), req.GetModel(),
			req.GetSampler(), req.GetCfgScale(), req.GetBatchSize(), req.GetBatchCount(), req.Seed, req.GetLoras(),
			req.GetControlUnits(), req.GetControlMode(), "",
		),
	}
}

func ImageToImageRequestToProto(req ImageToImageRequest) *wujiv1.ImageToImageRequest {
	w, h, steps, neg, model, sampler, cfg, bs, bc, seed, loras, units, mode := commonParamsToProto(req.ImageCommonParams)
	return &wujiv1.ImageToImageRequest{
		Prompt: req.Prompt, InitImagePath: req.InitImagePath, DenoisingStrength: req.DenoisingStrength,
		NegativePrompt: neg, Model: model, Width: w, Height: h, Steps: steps, Sampler: sampler,
		CfgScale: cfg, BatchSize: bs, BatchCount: bc, Seed: seed, Loras: loras,
		ControlUnits: units, ControlMode: mode,
	}
}

func ImageInpaintRequestFromProto(req *wujiv1.InpaintImageRequest) ImageInpaintRequest {
	if req == nil {
		return ImageInpaintRequest{}
	}
	return ImageInpaintRequest{
		Prompt: req.GetPrompt(), InitImagePath: req.GetInitImagePath(), MaskImagePath: req.GetMaskImagePath(),
		DenoisingStrength: req.GetDenoisingStrength(),
		ImageCommonParams: commonParamsFromProto(
			req.GetWidth(), req.GetHeight(), req.GetSteps(), req.GetNegativePrompt(), req.GetModel(),
			req.GetSampler(), req.GetCfgScale(), req.GetBatchSize(), req.GetBatchCount(), req.Seed, req.GetLoras(),
			req.GetControlUnits(), req.GetControlMode(), "",
		),
	}
}

func ImageInpaintRequestToProto(req ImageInpaintRequest) *wujiv1.InpaintImageRequest {
	w, h, steps, neg, model, sampler, cfg, bs, bc, seed, loras, units, mode := commonParamsToProto(req.ImageCommonParams)
	return &wujiv1.InpaintImageRequest{
		Prompt: req.Prompt, InitImagePath: req.InitImagePath, MaskImagePath: req.MaskImagePath,
		DenoisingStrength: req.DenoisingStrength, NegativePrompt: neg, Model: model, Width: w, Height: h,
		Steps: steps, Sampler: sampler, CfgScale: cfg, BatchSize: bs, BatchCount: bc, Seed: seed, Loras: loras,
		ControlUnits: units, ControlMode: mode,
	}
}

func ImageUpscaleRequestFromProto(req *wujiv1.UpscaleImageRequest) ImageUpscaleRequest {
	if req == nil {
		return ImageUpscaleRequest{}
	}
	return ImageUpscaleRequest{
		InitImagePath: req.GetInitImagePath(), Model: req.GetModel(), Scale: float32(req.GetScale()),
		LoRAs: lorasFromProto(req.GetLoras()),
	}
}

func ImageUpscaleRequestToProto(req ImageUpscaleRequest) *wujiv1.UpscaleImageRequest {
	return &wujiv1.UpscaleImageRequest{
		InitImagePath: req.InitImagePath, Model: req.Model, Scale: int32(req.Scale), Loras: lorasToProto(req.LoRAs),
	}
}

func ImageEditRequestFromProto(req *wujiv1.EditImageRequest) ImageEditRequest {
	if req == nil {
		return ImageEditRequest{}
	}
	return ImageEditRequest{
		Prompt: req.GetPrompt(), InitImagePath: req.GetInitImagePath(), DenoisingStrength: req.GetDenoisingStrength(),
		ImageCommonParams: commonParamsFromProto(
			req.GetWidth(), req.GetHeight(), req.GetSteps(), req.GetNegativePrompt(), req.GetModel(),
			req.GetSampler(), req.GetCfgScale(), req.GetBatchSize(), req.GetBatchCount(), req.Seed, req.GetLoras(),
			req.GetControlUnits(), req.GetControlMode(), "",
		),
	}
}

func ImageEditRequestToProto(req ImageEditRequest) *wujiv1.EditImageRequest {
	w, h, steps, neg, model, sampler, cfg, bs, bc, seed, loras, units, mode := commonParamsToProto(req.ImageCommonParams)
	return &wujiv1.EditImageRequest{
		Prompt: req.Prompt, InitImagePath: req.InitImagePath, DenoisingStrength: req.DenoisingStrength,
		NegativePrompt: neg, Model: model, Width: w, Height: h, Steps: steps, Sampler: sampler,
		CfgScale: cfg, BatchSize: bs, BatchCount: bc, Seed: seed, Loras: loras,
		ControlUnits: units, ControlMode: mode,
	}
}

func ImageControlNetRequestFromProto(req *wujiv1.ControlNetImageRequest) ImageControlNetRequest {
	if req == nil {
		return ImageControlNetRequest{}
	}
	units := controlUnitsFromProto(req.GetControlUnits())
	if len(units) == 0 && req.GetControlImagePath() != "" {
		typ := ImageControlDepth
		if req.GetControlType() != "" {
			if parsed, err := ParseImageControlType(req.GetControlType()); err == nil {
				typ = parsed
			}
		}
		units = []ControlNetUnit{{Type: typ, ImagePath: req.GetControlImagePath()}}
	}
	cp := commonParamsFromProto(
		req.GetWidth(), req.GetHeight(), req.GetSteps(), req.GetNegativePrompt(), req.GetModel(),
		req.GetSampler(), req.GetCfgScale(), req.GetBatchSize(), req.GetBatchCount(), req.Seed, req.GetLoras(),
		controlUnitsToProto(units), req.GetControlMode(), "",
	)
	return ImageControlNetRequest{
		Prompt: req.GetPrompt(), ControlUnits: units, ControlMode: cp.ControlMode,
		ImageCommonParams: cp,
	}
}

func ImageControlNetRequestToProto(req ImageControlNetRequest) *wujiv1.ControlNetImageRequest {
	w, h, steps, neg, model, sampler, cfg, bs, bc, seed, loras, _, mode := commonParamsToProto(req.ImageCommonParams)
	out := &wujiv1.ControlNetImageRequest{
		Prompt: req.Prompt, NegativePrompt: neg, Model: model, Width: w, Height: h, Steps: steps, Sampler: sampler,
		CfgScale: cfg, BatchSize: bs, BatchCount: bc, Seed: seed, Loras: loras,
		ControlUnits: controlUnitsToProto(req.ControlUnits), ControlMode: mode,
	}
	if len(req.ControlUnits) == 1 {
		out.ControlImagePath = req.ControlUnits[0].ImagePath
		out.ControlType = string(req.ControlUnits[0].Type)
	}
	return out
}

func ImageDepthToImageRequestFromProto(req *wujiv1.DepthToImageRequest) ImageDepthToImageRequest {
	if req == nil {
		return ImageDepthToImageRequest{}
	}
	return ImageDepthToImageRequest{
		Prompt: req.GetPrompt(), DepthImagePath: req.GetDepthImagePath(),
		ImageCommonParams: commonParamsFromProto(
			req.GetWidth(), req.GetHeight(), req.GetSteps(), req.GetNegativePrompt(), req.GetModel(),
			req.GetSampler(), req.GetCfgScale(), req.GetBatchSize(), req.GetBatchCount(), req.Seed, req.GetLoras(),
			req.GetControlUnits(), req.GetControlMode(), "",
		),
	}
}

func ImageDepthToImageRequestToProto(req ImageDepthToImageRequest) *wujiv1.DepthToImageRequest {
	w, h, steps, neg, model, sampler, cfg, bs, bc, seed, loras, units, mode := commonParamsToProto(req.ImageCommonParams)
	return &wujiv1.DepthToImageRequest{
		Prompt: req.Prompt, DepthImagePath: req.DepthImagePath, NegativePrompt: neg, Model: model,
		Width: w, Height: h, Steps: steps, Sampler: sampler, CfgScale: cfg, BatchSize: bs, BatchCount: bc,
		Seed: seed, Loras: loras, ControlUnits: units, ControlMode: mode,
	}
}

func ImageVariationRequestFromProto(req *wujiv1.ImageVariationRequest) ImageVariationRequest {
	if req == nil {
		return ImageVariationRequest{}
	}
	return ImageVariationRequest{
		InitImagePath:     req.GetInitImagePath(),
		DenoisingStrength: req.GetDenoisingStrength(),
		ImageCommonParams: commonParamsFromProto(
			req.GetWidth(), req.GetHeight(), req.GetSteps(), req.GetNegativePrompt(), req.GetModel(),
			req.GetSampler(), req.GetCfgScale(), req.GetBatchSize(), req.GetBatchCount(), req.Seed, req.GetLoras(),
			req.GetControlUnits(), req.GetControlMode(), "",
		),
	}
}

func ImageVariationRequestToProto(req ImageVariationRequest) *wujiv1.ImageVariationRequest {
	w, h, steps, neg, model, sampler, cfg, bs, bc, seed, loras, units, mode := commonParamsToProto(req.ImageCommonParams)
	out := &wujiv1.ImageVariationRequest{
		InitImagePath: req.InitImagePath, DenoisingStrength: req.DenoisingStrength,
		Model: model, NegativePrompt: neg, Width: w, Height: h, Steps: steps,
		Sampler: sampler, CfgScale: cfg, BatchSize: bs, BatchCount: bc,
		Loras: loras, ControlUnits: units, ControlMode: mode,
	}
	if seed != nil {
		out.Seed = seed
	}
	return out
}

func ImageStyleTransferRequestFromProto(req *wujiv1.StyleTransferImageRequest) ImageStyleTransferRequest {
	if req == nil {
		return ImageStyleTransferRequest{}
	}
	return ImageStyleTransferRequest{
		Prompt: req.GetPrompt(), StyleImagePath: req.GetStyleImagePath(), StyleWeight: req.GetStyleWeight(),
		ImageCommonParams: commonParamsFromProto(
			req.GetWidth(), req.GetHeight(), req.GetSteps(), req.GetNegativePrompt(), req.GetModel(),
			req.GetSampler(), req.GetCfgScale(), req.GetBatchSize(), req.GetBatchCount(), req.Seed, req.GetLoras(),
			req.GetControlUnits(), req.GetControlMode(), req.GetMode(),
		),
	}
}

func ImageStyleTransferRequestToProto(req ImageStyleTransferRequest) *wujiv1.StyleTransferImageRequest {
	w, h, steps, neg, model, sampler, cfg, bs, bc, seed, loras, units, controlMode := commonParamsToProto(req.ImageCommonParams)
	return &wujiv1.StyleTransferImageRequest{
		Prompt: req.Prompt, StyleImagePath: req.StyleImagePath, StyleWeight: req.StyleWeight,
		NegativePrompt: neg, Model: model, Width: w, Height: h, Steps: steps, Sampler: sampler,
		CfgScale: cfg, BatchSize: bs, BatchCount: bc, Seed: seed, Loras: loras,
		ControlUnits: units, ControlMode: controlMode,
		Mode: string(req.ImageCommonParams.ModeOrDefault()),
	}
}

func ImageSpriteRequestFromProto(req *wujiv1.SpriteImageRequest) ImageSpriteRequest {
	if req == nil {
		return ImageSpriteRequest{}
	}
	out := ImageSpriteRequest{
		Prompt: req.GetPrompt(), InitImagePath: req.GetInitImagePath(),
		ReferenceImagePath: req.GetReferenceImagePath(),
		StyleImagePath:     req.GetStyleImagePath(),
		StyleWeight:        req.GetStyleWeight(),
		Action:             SpriteAction(req.GetAction()),
		View:               SpriteView(req.GetView()),
		Directions:         int(req.GetDirections()),
		Loop:               req.GetLoop(),
		Padding:            int(req.GetPadding()),
		Transparent:        req.GetTransparent(),
		FrameWidth: int(req.GetFrameWidth()), FrameHeight: int(req.GetFrameHeight()),
		Columns: int(req.GetColumns()), Rows: int(req.GetRows()), FrameCount: int(req.GetFrameCount()),
		DenoisingStrength: req.GetDenoisingStrength(),
		ImageCommonParams: ImageCommonParams{
			NegativePrompt: req.GetNegativePrompt(),
			Model:          req.GetModel(),
			Steps:          int(req.GetSteps()),
			Sampler:        req.GetSampler(),
			CFGScale:       req.GetCfgScale(),
			LoRAs:          lorasFromProto(req.GetLoras()),
			ControlUnits:   controlUnitsFromProto(req.GetControlUnits()),
			ControlMode:    ControlNetMode(req.GetControlMode()),
			Mode:           AssetMode(req.GetMode()),
		},
	}
	if req.Seed != nil {
		seed := int(*req.Seed)
		out.Seed = &seed
	}
	return out
}

func ImageSpriteRequestToProto(req ImageSpriteRequest) *wujiv1.SpriteImageRequest {
	var seed *int32
	if req.Seed != nil {
		v := int32(*req.Seed)
		seed = &v
	}
	return &wujiv1.SpriteImageRequest{
		Prompt: req.Prompt, InitImagePath: req.InitImagePath,
		ReferenceImagePath: req.ReferenceImagePath,
		StyleImagePath:     req.StyleImagePath,
		StyleWeight:        req.StyleWeight,
		Action:             string(req.Action),
		View:               string(req.View),
		Directions:         int32(req.Directions),
		Loop:               req.Loop,
		Padding:            int32(req.Padding),
		Transparent:        req.Transparent,
		FrameWidth: int32(req.FrameWidth), FrameHeight: int32(req.FrameHeight),
		Columns: int32(req.Columns), Rows: int32(req.Rows), FrameCount: int32(req.FrameCount),
		Mode: string(req.ImageCommonParams.ModeOrDefault()),
		NegativePrompt: req.NegativePrompt, Model: req.Model, Steps: int32(req.Steps),
		Sampler: req.Sampler, CfgScale: req.CFGScale, Seed: seed, Loras: lorasToProto(req.LoRAs),
		ControlUnits: controlUnitsToProto(req.ControlUnits), ControlMode: string(req.ControlMode),
		DenoisingStrength: req.DenoisingStrength,
	}
}

func ImageResponseFromProto(resp *wujiv1.GenerateImageResponse) *ImageResponse {
	if resp == nil {
		return &ImageResponse{}
	}
	return &ImageResponse{
		Path: resp.GetPath(), Paths: append([]string(nil), resp.GetPaths()...), Format: resp.GetFormat(),
	}
}

func ImageResponseToProto(resp *ImageResponse) *wujiv1.GenerateImageResponse {
	if resp == nil {
		return &wujiv1.GenerateImageResponse{}
	}
	return &wujiv1.GenerateImageResponse{
		Path: resp.Path, Paths: append([]string(nil), resp.Paths...), Format: resp.Format,
	}
}

// Deprecated: use task-specific From/To helpers. Kept for transitional callers.
func ImageRequestFromProto(req *wujiv1.GenerateImageRequest) ImageRequest {
	gen := ImageGenerateRequestFromProto(req)
	r := ImageRequest{Task: ImageTaskGenerate, Prompt: gen.Prompt}
	cp := gen.ImageCommonParams
	r.NegativePrompt = cp.NegativePrompt
	r.Model = cp.Model
	r.Width = cp.Width
	r.Height = cp.Height
	r.Steps = cp.Steps
	r.Sampler = cp.Sampler
	r.CFGScale = cp.CFGScale
	r.BatchSize = cp.BatchSize
	r.BatchCount = cp.BatchCount
	r.Seed = cp.Seed
	r.LoRAs = cp.LoRAs
	r.ControlUnits = cp.ControlUnits
	r.ControlMode = cp.ControlMode
	r.Mode = cp.Mode
	return r
}

// Deprecated: use ImageGenerateRequestToProto.
func ImageRequestToProto(req ImageRequest) *wujiv1.GenerateImageRequest {
	return ImageGenerateRequestToProto(req.ToGenerateRequest())
}
