package driver

import wujiv1 "github.com/coditary/wuji-core/api/proto/v1"

func VoiceTasksToProto(tasks []VoiceTask) []string {
	if len(tasks) == 0 {
		return nil
	}
	out := make([]string, 0, len(tasks))
	for _, task := range tasks {
		out = append(out, string(task))
	}
	return out
}

func VoiceTasksFromProto(items []string) []VoiceTask {
	if len(items) == 0 {
		return nil
	}
	out := make([]VoiceTask, 0, len(items))
	for _, item := range items {
		task := VoiceTask(item)
		if task.IsValid() {
			out = append(out, task)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func VoiceCreateRequestFromProto(req *wujiv1.CreateVoiceRequest) VoiceCreateRequest {
	if req == nil {
		return VoiceCreateRequest{}
	}
	return VoiceCreateRequest{
		Name: req.GetName(), SamplePath: req.GetSamplePath(), TargetModel: req.GetTargetModel(),
		Denoise: req.GetDenoise(), DenoiseStrength: req.GetDenoiseStrength(),
	}
}

func VoiceCreateRequestToProto(req VoiceCreateRequest) *wujiv1.CreateVoiceRequest {
	return &wujiv1.CreateVoiceRequest{
		Name: req.Name, SamplePath: req.SamplePath, TargetModel: req.TargetModel,
		Denoise: req.Denoise, DenoiseStrength: req.DenoiseStrength,
	}
}

func VoiceConvertRequestFromProto(req *wujiv1.ConvertVoiceRequest) VoiceConvertRequest {
	if req == nil {
		return VoiceConvertRequest{}
	}
	return VoiceConvertRequest{
		Name: req.GetName(), SamplePath: req.GetSamplePath(), TargetModel: req.GetTargetModel(), SourcePath: req.GetSourcePath(),
		PitchShift: int(req.GetPitchShift()), IndexRate: req.GetIndexRate(), Protect: req.GetProtect(), Vocoder: req.GetVocoder(),
		Denoise: req.GetDenoise(), DenoiseStrength: req.GetDenoiseStrength(), ChunkSize: int(req.GetChunkSize()), Crossfade: req.GetCrossfade(),
	}
}

func VoiceConvertRequestToProto(req VoiceConvertRequest) *wujiv1.ConvertVoiceRequest {
	return &wujiv1.ConvertVoiceRequest{
		Name: req.Name, SamplePath: req.SamplePath, TargetModel: req.TargetModel, SourcePath: req.SourcePath,
		PitchShift: int32(req.PitchShift), IndexRate: req.IndexRate, Protect: req.Protect, Vocoder: req.Vocoder,
		Denoise: req.Denoise, DenoiseStrength: req.DenoiseStrength, ChunkSize: int32(req.ChunkSize), Crossfade: req.Crossfade,
	}
}

func VoiceResponseFromProto(resp *wujiv1.CloneVoiceResponse) *VoiceResponse {
	if resp == nil {
		return &VoiceResponse{}
	}
	return &VoiceResponse{VoiceID: resp.GetVoiceId(), Name: resp.GetName(), OutputPath: resp.GetOutputPath()}
}

func VoiceResponseToProto(resp *VoiceResponse) *wujiv1.CloneVoiceResponse {
	if resp == nil {
		return &wujiv1.CloneVoiceResponse{}
	}
	return &wujiv1.CloneVoiceResponse{VoiceId: resp.VoiceID, Name: resp.Name, OutputPath: resp.OutputPath}
}
