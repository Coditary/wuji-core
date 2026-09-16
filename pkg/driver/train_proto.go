package driver

import (
	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/capability"
)

func TextTrainRequestFromProto(req *wujiv1.TrainTextRequest) TextTrainRequest {
	if req == nil {
		return TextTrainRequest{}
	}
	method, _ := ParseTextTrainMethod(req.GetMethod())
	return TextTrainRequest{
		Method: method, Name: req.GetName(), DatasetID: req.GetDatasetId(), BaseModel: req.GetBaseModel(),
		OutputPath: req.GetOutputPath(), Epochs: int(req.GetEpochs()),
		LearningRate: req.GetLearningRate(), LoRARank: int(req.GetLoraRank()), LoRAAlpha: int(req.GetLoraAlpha()),
		ContextLength: int(req.GetContextLength()), BatchSize: int(req.GetBatchSize()),
		GradientAccumulationSteps: int(req.GetGradientAccumulationSteps()),
		WarmupSteps:               int(req.GetWarmupSteps()), SaveEveryEpoch: int(req.GetSaveEveryEpoch()),
		Seed: int(req.GetSeed()),
	}
}

func TextTrainRequestToProto(req TextTrainRequest) *wujiv1.TrainTextRequest {
	method := string(req.Method)
	if method == "" {
		method = string(TextTrainMethodFull)
	}
	return &wujiv1.TrainTextRequest{
		Method: method, Name: req.Name, DatasetId: req.DatasetID, BaseModel: req.BaseModel,
		OutputPath: req.OutputPath, Epochs: int32(req.Epochs),
		LearningRate: req.LearningRate, LoraRank: int32(req.LoRARank), LoraAlpha: int32(req.LoRAAlpha),
		ContextLength: int32(req.ContextLength), BatchSize: int32(req.BatchSize),
		GradientAccumulationSteps: int32(req.GradientAccumulationSteps),
		WarmupSteps:               int32(req.WarmupSteps), SaveEveryEpoch: int32(req.SaveEveryEpoch),
		Seed: int32(req.Seed),
	}
}

func ImageTrainRequestFromProto(req *wujiv1.TrainImageRequest) ImageTrainRequest {
	if req == nil {
		return ImageTrainRequest{}
	}
	method, _ := ParseImageTrainMethod(req.GetMethod())
	var controlType ImageControlType
	if req.GetControlType() != "" {
		controlType, _ = ParseImageControlType(req.GetControlType())
	}
	var mode ImageMode
	if req.GetMode() != "" {
		mode, _ = ParseImageMode(req.GetMode())
	}
	return ImageTrainRequest{
		Method: method, Name: req.GetName(), DatasetID: req.GetDatasetId(), BaseModel: req.GetBaseModel(),
		OutputPath: req.GetOutputPath(), Epochs: int(req.GetEpochs()),
		LearningRate: req.GetLearningRate(), LoRARank: int(req.GetLoraRank()), LoRAAlpha: int(req.GetLoraAlpha()),
		Width: int(req.GetWidth()), Height: int(req.GetHeight()), ClassToken: req.GetClassToken(),
		ControlType: controlType, Token: req.GetToken(), Mode: mode,
		PriorPreservation: req.GetPriorPreservation(), RegDatasetID: req.GetRegDatasetId(),
		BatchSize: int(req.GetBatchSize()), Seed: int(req.GetSeed()),
	}
}

func ImageTrainRequestToProto(req ImageTrainRequest) *wujiv1.TrainImageRequest {
	method := string(req.Method)
	if method == "" {
		method = string(ImageTrainMethodFull)
	}
	return &wujiv1.TrainImageRequest{
		Method: method, Name: req.Name, DatasetId: req.DatasetID, BaseModel: req.BaseModel,
		OutputPath: req.OutputPath, Epochs: int32(req.Epochs),
		LearningRate: req.LearningRate, LoraRank: int32(req.LoRARank), LoraAlpha: int32(req.LoRAAlpha),
		Width: int32(req.Width), Height: int32(req.Height), ClassToken: req.ClassToken,
		ControlType: string(req.ControlType), Token: req.Token, Mode: string(req.Mode),
		PriorPreservation: req.PriorPreservation, RegDatasetId: req.RegDatasetID,
		BatchSize: int32(req.BatchSize), Seed: int32(req.Seed),
	}
}

func VideoTrainRequestFromProto(req *wujiv1.TrainVideoRequest) VideoTrainRequest {
	if req == nil {
		return VideoTrainRequest{}
	}
	method, _ := ParseVideoTrainMethod(req.GetMethod())
	return VideoTrainRequest{
		Method: method, Name: req.GetName(), DatasetID: req.GetDatasetId(), BaseModel: req.GetBaseModel(),
		OutputPath: req.GetOutputPath(), Epochs: int(req.GetEpochs()),
		LearningRate: req.GetLearningRate(), LoRARank: int(req.GetLoraRank()),
		Frames: int(req.GetFrames()), FPS: int(req.GetFps()), ContextLength: int(req.GetContextLength()),
		BatchSize: int(req.GetBatchSize()), Seed: int(req.GetSeed()),
	}
}

func VideoTrainRequestToProto(req VideoTrainRequest) *wujiv1.TrainVideoRequest {
	method := string(req.Method)
	if method == "" {
		method = string(VideoTrainMethodFull)
	}
	return &wujiv1.TrainVideoRequest{
		Method: method, Name: req.Name, DatasetId: req.DatasetID, BaseModel: req.BaseModel,
		OutputPath: req.OutputPath, Epochs: int32(req.Epochs),
		LearningRate: req.LearningRate, LoraRank: int32(req.LoRARank),
		Frames: int32(req.Frames), Fps: int32(req.FPS), ContextLength: int32(req.ContextLength),
		BatchSize: int32(req.BatchSize), Seed: int32(req.Seed),
	}
}

func AudioTrainRequestFromProto(req *wujiv1.TrainAudioRequest) AudioTrainRequest {
	if req == nil {
		return AudioTrainRequest{}
	}
	method, _ := ParseAudioTrainMethod(req.GetMethod())
	if method == "" && req.GetTaskType() != "" {
		method, _ = ParseAudioTrainMethod(req.GetTaskType())
	}
	taskType := method.ToAudioTask()
	if req.GetTaskType() != "" {
		taskType, _ = ParseAudioTask(req.GetTaskType())
	}
	return AudioTrainRequest{
		Method: method, Name: req.GetName(), DatasetID: req.GetDatasetId(), BaseModel: req.GetBaseModel(),
		OutputPath: req.GetOutputPath(), Epochs: int(req.GetEpochs()),
		LearningRate: req.GetLearningRate(), SampleRate: int(req.GetSampleRate()),
		BatchSize: int(req.GetBatchSize()), TaskType: taskType, Seed: int(req.GetSeed()),
	}
}

func AudioTrainRequestToProto(req AudioTrainRequest) *wujiv1.TrainAudioRequest {
	method := string(req.Method)
	if method == "" {
		method = string(AudioTrainMethodMusic)
	}
	return &wujiv1.TrainAudioRequest{
		Method: method, Name: req.Name, DatasetId: req.DatasetID, BaseModel: req.BaseModel,
		OutputPath: req.OutputPath, Epochs: int32(req.Epochs),
		LearningRate: req.LearningRate, SampleRate: int32(req.SampleRate),
		BatchSize: int32(req.BatchSize), TaskType: string(req.TaskType), Seed: int32(req.Seed),
	}
}

func MeshTrainRequestFromProto(req *wujiv1.TrainMeshRequest) MeshTrainRequest {
	if req == nil {
		return MeshTrainRequest{}
	}
	method, _ := ParseMeshTrainMethod(req.GetMethod())
	var mode MeshMode
	if req.GetMode() != "" {
		mode, _ = ParseMeshMode(req.GetMode())
	}
	return MeshTrainRequest{
		Method: method, Name: req.GetName(), DatasetID: req.GetDatasetId(), BaseModel: req.GetBaseModel(),
		OutputPath: req.GetOutputPath(), Epochs: int(req.GetEpochs()),
		LearningRate: req.GetLearningRate(), LoRARank: int(req.GetLoraRank()), Format: req.GetFormat(),
		Mode: mode, Seed: int(req.GetSeed()),
	}
}

func MeshTrainRequestToProto(req MeshTrainRequest) *wujiv1.TrainMeshRequest {
	method := string(req.Method)
	if method == "" {
		method = string(MeshTrainMethodFull)
	}
	return &wujiv1.TrainMeshRequest{
		Method: method, Name: req.Name, DatasetId: req.DatasetID, BaseModel: req.BaseModel,
		OutputPath: req.OutputPath, Epochs: int32(req.Epochs),
		LearningRate: req.LearningRate, LoraRank: int32(req.LoRARank), Format: req.Format,
		Mode: string(req.Mode), Seed: int32(req.Seed),
	}
}

func VoiceTrainRequestFromProto(req *wujiv1.TrainVoiceRequest) VoiceTrainRequest {
	if req == nil {
		return VoiceTrainRequest{}
	}
	method, _ := ParseVoiceTrainMethod(req.GetMethod())
	return VoiceTrainRequest{
		Method: method, Name: req.GetName(), DatasetID: req.GetDatasetId(), OutputPath: req.GetOutputPath(),
		Epochs: int(req.GetEpochs()), SamplePath: req.GetSamplePath(),
		PretrainedModel: req.GetPretrainedModel(), PitchShift: int(req.GetPitchShift()),
		IndexRate: req.GetIndexRate(), BatchSize: int(req.GetBatchSize()),
		SaveEveryEpoch: int(req.GetSaveEveryEpoch()), Seed: int(req.GetSeed()),
	}
}

func VoiceTrainRequestToProto(req VoiceTrainRequest) *wujiv1.TrainVoiceRequest {
	method := string(req.Method)
	if method == "" {
		method = string(VoiceTrainMethodRVC)
	}
	return &wujiv1.TrainVoiceRequest{
		Method: method, Name: req.Name, DatasetId: req.DatasetID, OutputPath: req.OutputPath,
		Epochs: int32(req.Epochs), SamplePath: req.SamplePath,
		PretrainedModel: req.PretrainedModel, PitchShift: int32(req.PitchShift),
		IndexRate: req.IndexRate, BatchSize: int32(req.BatchSize),
		SaveEveryEpoch: int32(req.SaveEveryEpoch), Seed: int32(req.Seed),
	}
}

func TrainResponseFromProto(resp *wujiv1.TrainResponse) *TrainResponse {
	if resp == nil {
		return &TrainResponse{}
	}
	return &TrainResponse{
		JobID: resp.GetJobId(), Status: resp.GetStatus(),
		Capability: capability.Type(resp.GetCapability()), OutputPath: resp.GetOutputPath(),
	}
}

func TrainResponseToProto(resp *TrainResponse) *wujiv1.TrainResponse {
	if resp == nil {
		return &wujiv1.TrainResponse{}
	}
	return &wujiv1.TrainResponse{
		JobId: resp.JobID, Status: resp.Status,
		Capability: string(resp.Capability), OutputPath: resp.OutputPath,
	}
}
