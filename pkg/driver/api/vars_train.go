package api

import (
	"fmt"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

func textTrainHTTPVars(req driver.TextTrainRequest, spec *config.APICapabilitySpec) map[string]string {
	vars := map[string]string{
		"task":                        "text",
		"capability":                  "text",
		"method":                      string(req.Method),
		"name":                        req.Name,
		"dataset_id":                  req.DatasetID,
		"base_model":                  req.BaseModel,
		"output_path":                 req.OutputPath,
		"epochs":                      fmt.Sprintf("%d", req.Epochs),
		"learning_rate":               fmt.Sprintf("%g", float64(req.LearningRate)),
		"lora_rank":                   fmt.Sprintf("%d", req.LoRARank),
		"lora_alpha":                  fmt.Sprintf("%d", req.LoRAAlpha),
		"context_length":              fmt.Sprintf("%d", req.ContextLength),
		"batch_size":                  fmt.Sprintf("%d", req.BatchSize),
		"gradient_accumulation_steps": fmt.Sprintf("%d", req.GradientAccumulationSteps),
		"warmup_steps":                fmt.Sprintf("%d", req.WarmupSteps),
		"save_every_epoch":            fmt.Sprintf("%d", req.SaveEveryEpoch),
		"seed":                        fmt.Sprintf("%d", req.Seed),
	}
	applySpecDefaults(vars, spec)
	return vars
}

func imageTrainHTTPVars(req driver.ImageTrainRequest, spec *config.APICapabilitySpec) map[string]string {
	vars := map[string]string{
		"task":               "image",
		"capability":         "image",
		"method":             string(req.Method),
		"name":               req.Name,
		"dataset_id":         req.DatasetID,
		"base_model":         req.BaseModel,
		"output_path":        req.OutputPath,
		"epochs":             fmt.Sprintf("%d", req.Epochs),
		"learning_rate":      fmt.Sprintf("%g", float64(req.LearningRate)),
		"lora_rank":          fmt.Sprintf("%d", req.LoRARank),
		"lora_alpha":         fmt.Sprintf("%d", req.LoRAAlpha),
		"width":              fmt.Sprintf("%d", req.Width),
		"height":             fmt.Sprintf("%d", req.Height),
		"class_token":        req.ClassToken,
		"control_type":       string(req.ControlType),
		"token":              req.Token,
		"mode":               string(req.Mode),
		"prior_preservation": boolStr(req.PriorPreservation),
		"reg_dataset_id":     req.RegDatasetID,
		"batch_size":         fmt.Sprintf("%d", req.BatchSize),
		"seed":               fmt.Sprintf("%d", req.Seed),
	}
	applySpecDefaults(vars, spec)
	return vars
}

func videoTrainHTTPVars(req driver.VideoTrainRequest, spec *config.APICapabilitySpec) map[string]string {
	vars := map[string]string{
		"task":           "video",
		"capability":     "video",
		"method":         string(req.Method),
		"name":           req.Name,
		"dataset_id":     req.DatasetID,
		"base_model":     req.BaseModel,
		"output_path":    req.OutputPath,
		"epochs":         fmt.Sprintf("%d", req.Epochs),
		"learning_rate":  fmt.Sprintf("%g", float64(req.LearningRate)),
		"lora_rank":      fmt.Sprintf("%d", req.LoRARank),
		"frames":         fmt.Sprintf("%d", req.Frames),
		"fps":            fmt.Sprintf("%d", req.FPS),
		"context_length": fmt.Sprintf("%d", req.ContextLength),
		"batch_size":     fmt.Sprintf("%d", req.BatchSize),
		"seed":           fmt.Sprintf("%d", req.Seed),
	}
	applySpecDefaults(vars, spec)
	return vars
}

func audioTrainHTTPVars(req driver.AudioTrainRequest, spec *config.APICapabilitySpec) map[string]string {
	vars := map[string]string{
		"task":          "audio",
		"capability":    "audio",
		"method":        string(req.Method),
		"name":          req.Name,
		"dataset_id":    req.DatasetID,
		"base_model":    req.BaseModel,
		"output_path":   req.OutputPath,
		"epochs":        fmt.Sprintf("%d", req.Epochs),
		"learning_rate": fmt.Sprintf("%g", float64(req.LearningRate)),
		"sample_rate":   fmt.Sprintf("%d", req.SampleRate),
		"batch_size":    fmt.Sprintf("%d", req.BatchSize),
		"audio_task":    string(req.TaskType),
		"seed":          fmt.Sprintf("%d", req.Seed),
	}
	applySpecDefaults(vars, spec)
	return vars
}

func meshTrainHTTPVars(req driver.MeshTrainRequest, spec *config.APICapabilitySpec) map[string]string {
	vars := map[string]string{
		"task":          "mesh",
		"capability":    "mesh",
		"method":        string(req.Method),
		"name":          req.Name,
		"dataset_id":    req.DatasetID,
		"base_model":    req.BaseModel,
		"output_path":   req.OutputPath,
		"epochs":        fmt.Sprintf("%d", req.Epochs),
		"learning_rate": fmt.Sprintf("%g", float64(req.LearningRate)),
		"lora_rank":     fmt.Sprintf("%d", req.LoRARank),
		"format":        req.Format,
		"mode":          string(req.Mode),
		"seed":          fmt.Sprintf("%d", req.Seed),
	}
	applySpecDefaults(vars, spec)
	return vars
}

func voiceTrainHTTPVars(req driver.VoiceTrainRequest, spec *config.APICapabilitySpec) map[string]string {
	vars := map[string]string{
		"task":             "voice",
		"capability":       "voice",
		"method":           string(req.Method),
		"name":             req.Name,
		"dataset_id":       req.DatasetID,
		"output_path":      req.OutputPath,
		"epochs":           fmt.Sprintf("%d", req.Epochs),
		"sample_path":      req.SamplePath,
		"pretrained_model": req.PretrainedModel,
		"pitch_shift":      fmt.Sprintf("%d", req.PitchShift),
		"index_rate":       fmt.Sprintf("%g", float64(req.IndexRate)),
		"batch_size":       fmt.Sprintf("%d", req.BatchSize),
		"save_every_epoch": fmt.Sprintf("%d", req.SaveEveryEpoch),
		"seed":             fmt.Sprintf("%d", req.Seed),
	}
	applySpecDefaults(vars, spec)
	return vars
}
