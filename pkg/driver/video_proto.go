package driver

import wujiv1 "github.com/coditary/wuji-core/api/proto/v1"

func VideoTasksToProto(tasks []VideoTask) []string {
	return videoTasksToProto(tasks)
}

func VideoTasksFromProto(items []string) []VideoTask {
	return videoTasksFromProto(items)
}

func videoTasksToProto(tasks []VideoTask) []string {
	if len(tasks) == 0 {
		return nil
	}
	out := make([]string, 0, len(tasks))
	for _, task := range tasks {
		out = append(out, string(task))
	}
	return out
}

func videoTasksFromProto(items []string) []VideoTask {
	if len(items) == 0 {
		return nil
	}
	out := make([]VideoTask, 0, len(items))
	for _, item := range items {
		task := VideoTask(item)
		if task.IsValid() {
			out = append(out, task)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func videoCommonToProto(p VideoCommonParams) (float32, int32, int32, float32, int32, string, string, string, string, *int32) {
	var seed *int32
	if p.Seed != nil {
		v := int32(*p.Seed)
		seed = &v
	}
	return p.Duration, int32(p.FPS), int32(p.Frames), p.MotionStrength, int32(p.ContextLength),
		p.Sampler, p.Scheduler, p.NegativePrompt, p.Model, seed
}

func videoCommonFromProto(duration float32, fps, frames int32, motionStrength float32, contextLength int32, sampler, scheduler, negativePrompt, model string, seed *int32) VideoCommonParams {
	out := VideoCommonParams{
		Duration: duration, FPS: int(fps), Frames: int(frames), MotionStrength: motionStrength,
		ContextLength: int(contextLength), Sampler: sampler, Scheduler: scheduler,
		NegativePrompt: negativePrompt, Model: model,
	}
	if seed != nil {
		v := int(*seed)
		out.Seed = &v
	}
	return out
}

func VideoGenerateRequestFromProto(req *wujiv1.GenerateVideoRequest) VideoGenerateRequest {
	if req == nil {
		return VideoGenerateRequest{}
	}
	return VideoGenerateRequest{
		Prompt:        req.GetPrompt(),
		CameraControl: CameraControl(req.GetCameraControl()),
		VideoCommonParams: videoCommonFromProto(
			req.GetDuration(), req.GetFps(), req.GetFrames(), req.GetMotionStrength(), req.GetContextLength(),
			req.GetSampler(), req.GetScheduler(), req.GetNegativePrompt(), req.GetModel(), req.Seed,
		),
	}
}

func VideoGenerateRequestToProto(req VideoGenerateRequest) *wujiv1.GenerateVideoRequest {
	dur, fps, frames, motion, ctxLen, sampler, scheduler, neg, model, seed := videoCommonToProto(req.VideoCommonParams)
	return &wujiv1.GenerateVideoRequest{
		Prompt: req.Prompt, Duration: dur, Fps: fps, Frames: frames, MotionStrength: motion,
		ContextLength: ctxLen, Sampler: sampler, Scheduler: scheduler, CameraControl: string(req.CameraControl),
		NegativePrompt: neg, Model: model, Seed: seed,
	}
}

func VideoImageToVideoRequestFromProto(req *wujiv1.ImageToVideoRequest) VideoImageToVideoRequest {
	if req == nil {
		return VideoImageToVideoRequest{}
	}
	return VideoImageToVideoRequest{
		Prompt: req.GetPrompt(), InitImagePath: req.GetInitImagePath(),
		VideoCommonParams: videoCommonFromProto(
			req.GetDuration(), req.GetFps(), req.GetFrames(), req.GetMotionStrength(), req.GetContextLength(),
			req.GetSampler(), req.GetScheduler(), req.GetNegativePrompt(), req.GetModel(), req.Seed,
		),
	}
}

func VideoImageToVideoRequestToProto(req VideoImageToVideoRequest) *wujiv1.ImageToVideoRequest {
	dur, fps, frames, motion, ctxLen, sampler, scheduler, neg, model, seed := videoCommonToProto(req.VideoCommonParams)
	return &wujiv1.ImageToVideoRequest{
		Prompt: req.Prompt, InitImagePath: req.InitImagePath, Duration: dur, Fps: fps, Frames: frames,
		MotionStrength: motion, ContextLength: ctxLen, Sampler: sampler, Scheduler: scheduler,
		NegativePrompt: neg, Model: model, Seed: seed,
	}
}

func VideoInterpolateRequestFromProto(req *wujiv1.InterpolateVideoRequest) VideoInterpolateRequest {
	if req == nil {
		return VideoInterpolateRequest{}
	}
	return VideoInterpolateRequest{
		InputVideoPath: req.GetInputVideoPath(),
		TargetFPS:      int(req.GetTargetFps()),
		Model:          req.GetModel(),
	}
}

func VideoInterpolateRequestToProto(req VideoInterpolateRequest) *wujiv1.InterpolateVideoRequest {
	return &wujiv1.InterpolateVideoRequest{
		InputVideoPath: req.InputVideoPath,
		TargetFps:      int32(req.TargetFPS),
		Model:          req.Model,
	}
}

func VideoUpscaleRequestFromProto(req *wujiv1.UpscaleVideoRequest) VideoUpscaleRequest {
	if req == nil {
		return VideoUpscaleRequest{}
	}
	return VideoUpscaleRequest{
		InputVideoPath: req.GetInputVideoPath(), Scale: float32(req.GetScale()), Model: req.GetModel(),
	}
}

func VideoUpscaleRequestToProto(req VideoUpscaleRequest) *wujiv1.UpscaleVideoRequest {
	return &wujiv1.UpscaleVideoRequest{
		InputVideoPath: req.InputVideoPath, Scale: int32(req.Scale), Model: req.Model,
	}
}

func VideoResponseFromProto(resp *wujiv1.GenerateVideoResponse) *VideoResponse {
	if resp == nil {
		return &VideoResponse{}
	}
	return &VideoResponse{
		Path: resp.GetPath(), Duration: resp.GetDuration(),
		Frames: int(resp.GetFrames()), FPS: int(resp.GetFps()),
	}
}

func VideoResponseToProto(resp *VideoResponse) *wujiv1.GenerateVideoResponse {
	if resp == nil {
		return &wujiv1.GenerateVideoResponse{}
	}
	return &wujiv1.GenerateVideoResponse{
		Path: resp.Path, Duration: resp.Duration,
		Frames: int32(resp.Frames), Fps: int32(resp.FPS),
	}
}
