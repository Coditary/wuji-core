package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

func textHTTPVars(req driver.TextRequest, spec *config.APICapabilitySpec) map[string]string {
	msgs, _ := json.Marshal(req.Messages)
	tools, _ := json.Marshal(req.Tools)
	stops, _ := json.Marshal(req.StopSequences)
	loras, _ := json.Marshal(req.LoRAs)
	vars := map[string]string{
		"task":                "generate",
		"input_mode":          string(req.InputModeOrDefault()),
		"prompt":              req.Prompt,
		"system_prompt":       req.SystemPrompt,
		"media_path":          req.MediaPath,
		"translate":           boolStr(req.Translate),
		"target_lang":         req.TargetLang,
		"language":            req.Language,
		"beam_size":           fmt.Sprintf("%d", req.BeamSize),
		"word_timestamps":     boolStr(req.WordTimestamps),
		"vad_enabled":         boolStr(req.VADEnabled),
		"vad_threshold":       fmt.Sprintf("%g", float64(req.VADThreshold)),
		"loras_json":          string(loras),
		"seed":                seedStr(req.Seed),
		"model":               firstNonEmpty(req.Model, specModel(spec)),
		"max_tokens":          fmt.Sprintf("%d", firstInt(req.MaxTokens, specInt(spec, "max_tokens"))),
		"temperature":         fmt.Sprintf("%g", firstFloat(float64(req.Temperature), specFloat(spec, "temperature"))),
		"top_p":               fmt.Sprintf("%g", float64(req.TopP)),
		"top_k":               fmt.Sprintf("%d", req.TopK),
		"min_p":               fmt.Sprintf("%g", float64(req.MinP)),
		"frequency_penalty":   fmt.Sprintf("%g", float64(req.FrequencyPenalty)),
		"presence_penalty":    fmt.Sprintf("%g", float64(req.PresencePenalty)),
		"repetition_penalty":  fmt.Sprintf("%g", float64(req.RepetitionPenalty)),
		"stop_sequences_json": string(stops),
		"messages_json":       string(msgs),
		"tools_json":          string(tools),
		"context_window":      fmt.Sprintf("%d", req.ContextWindow),
	}
	applySpecDefaults(vars, spec)
	return vars
}

func textMediaHTTPVars(req driver.TextRequest, taskKey string, spec *config.APICapabilitySpec) map[string]string {
	vars := textHTTPVars(req, spec)
	vars["task"] = taskKey
	vars["media_path"] = req.MediaPath
	return vars
}

func imageHTTPVars(req driver.ImageRequest, spec *config.APICapabilitySpec) map[string]string {
	task := req.TaskOrDefault()
	units, _ := json.Marshal(req.ControlUnits)
	loras, _ := json.Marshal(req.LoRAs)
	vars := map[string]string{
		"task":                 string(task),
		"prompt":               req.Prompt,
		"negative_prompt":      req.NegativePrompt,
		"model":                firstNonEmpty(req.Model, specModel(spec)),
		"width":                fmt.Sprintf("%d", firstInt(req.Width, specDefaultInt(spec, "width"))),
		"height":               fmt.Sprintf("%d", firstInt(req.Height, specDefaultInt(spec, "height"))),
		"steps":                fmt.Sprintf("%d", firstInt(req.Steps, specDefaultInt(spec, "steps"))),
		"cfg_scale":            fmt.Sprintf("%g", float64(req.CFGScale)),
		"sampler":              req.Sampler,
		"batch_size":           fmt.Sprintf("%d", req.BatchSize),
		"batch_count":          fmt.Sprintf("%d", req.BatchCount),
		"seed":                 seedStr(req.Seed),
		"denoising_strength":   fmt.Sprintf("%g", float64(req.DenoisingStrength)),
		"init_image_path":      req.InitImagePath,
		"mask_image_path":      req.MaskImagePath,
		"control_image_path":   req.ControlImagePath,
		"control_type":         string(req.ControlType),
		"control_units_json":   string(units),
		"control_mode":         string(req.ControlMode),
		"style_image_path":     req.StyleImagePath,
		"style_weight":         fmt.Sprintf("%g", float64(req.StyleWeight)),
		"scale":                fmt.Sprintf("%g", float64(req.Scale)),
		"loras_json":           string(loras),
		"mode":                 string(req.Mode),
		"frame_width":          fmt.Sprintf("%d", req.FrameWidth),
		"frame_height":         fmt.Sprintf("%d", req.FrameHeight),
		"columns":              fmt.Sprintf("%d", req.Columns),
		"rows":                 fmt.Sprintf("%d", req.Rows),
		"frame_count":          fmt.Sprintf("%d", req.FrameCount),
		"reference_image_path": req.ReferenceImagePath,
		"sprite_action":        string(req.SpriteAction),
		"sprite_view":          string(req.SpriteView),
		"sprite_directions":    fmt.Sprintf("%d", req.SpriteDirections),
		"sprite_loop":          fmt.Sprintf("%t", req.SpriteLoop),
		"sprite_padding":       fmt.Sprintf("%d", req.SpritePadding),
		"sprite_transparent":   fmt.Sprintf("%t", req.SpriteTransparent),
		"depth_image_path":     req.ControlImagePath,
	}
	applySpecDefaults(vars, spec)
	return vars
}

func audioHTTPVars(req driver.AudioRequest, spec *config.APICapabilitySpec) map[string]string {
	task := req.TaskOrDefault()
	vars := map[string]string{
		"task":            string(task),
		"prompt":          req.Prompt,
		"lyrics":          req.Lyrics,
		"negative_prompt": req.NegativePrompt,
		"model":           firstNonEmpty(req.Model, specModel(spec)),
		"voice":           req.Voice,
		"language":        req.Language,
		"duration":        fmt.Sprintf("%g", req.Duration),
		"overlap":         fmt.Sprintf("%g", req.Overlap),
		"temperature":     fmt.Sprintf("%g", float64(req.Temperature)),
		"cfg_scale":       fmt.Sprintf("%g", float64(req.CFGScale)),
		"top_p":           fmt.Sprintf("%g", float64(req.TopP)),
		"top_k":           fmt.Sprintf("%d", req.TopK),
		"sample_rate":     fmt.Sprintf("%d", req.SampleRate),
		"reference_path":  req.ReferencePath,
		"format":          req.Format,
		"speed":           fmt.Sprintf("%g", float64(req.Speed)),
		"pitch":           fmt.Sprintf("%d", req.Pitch),
		"emotion":         req.Emotion,
		"style":           req.Style,
		"energy":          fmt.Sprintf("%g", float64(req.Energy)),
		"seed":            seedStr(req.Seed),
	}
	applySpecDefaults(vars, spec)
	return vars
}

func audioSpeechHTTPVars(req driver.AudioSpeechRequest, spec *config.APICapabilitySpec) map[string]string {
	audio := driver.AudioRequest{
		Task: driver.AudioTaskTextToSpeech, Prompt: req.Prompt, Model: req.Model,
		Voice: req.Voice, Language: req.Language, Duration: req.Duration,
		Speed: req.Speed, Pitch: req.Pitch, Emotion: req.Emotion, Style: req.Style,
		Energy: req.Energy, Seed: req.Seed,
	}
	return audioHTTPVars(audio, spec)
}

func videoHTTPVars(req driver.VideoRequest, spec *config.APICapabilitySpec) map[string]string {
	task := req.TaskOrDefault()
	vars := map[string]string{
		"task":              string(task),
		"prompt":            req.Prompt,
		"negative_prompt":   req.NegativePrompt,
		"model":             firstNonEmpty(req.Model, specModel(spec)),
		"frames":            fmt.Sprintf("%d", req.Frames),
		"fps":               fmt.Sprintf("%d", req.FPS),
		"duration":          fmt.Sprintf("%g", req.Duration),
		"init_image_path":   req.InitImagePath,
		"init_video_path":   req.InitVideoPath,
		"motion_strength":   fmt.Sprintf("%g", float64(req.MotionStrength)),
		"context_length":    fmt.Sprintf("%d", req.ContextLength),
		"sampler":           req.Sampler,
		"scheduler":         req.Scheduler,
		"camera_control":    string(req.CameraControl),
		"scale":             fmt.Sprintf("%g", float64(req.Scale)),
		"seed":              seedStr(req.Seed),
		"target_fps":        fmt.Sprintf("%d", req.FPS),
		"input_video_path":  req.InitVideoPath,
	}
	applySpecDefaults(vars, spec)
	return vars
}

func meshHTTPVars(req driver.MeshRequest, spec *config.APICapabilitySpec) map[string]string {
	task := req.TaskOrDefault()
	paths, _ := json.Marshal(req.InitImagePaths)
	vars := map[string]string{
		"task":                  string(task),
		"prompt":                req.Prompt,
		"format":                req.Format,
		"model":                 firstNonEmpty(req.Model, specModel(spec)),
		"seed":                  seedStr(req.Seed),
		"target_tris":           fmt.Sprintf("%d", req.TargetTris),
		"scale":                 fmt.Sprintf("%g", float64(req.Scale)),
		"mode":                  string(req.Mode),
		"representation":        string(req.Representation),
		"init_mesh_path":        req.InitMeshPath,
		"high_mesh_path":        req.HighMeshPath,
		"init_image_path":       req.InitImagePath,
		"init_image_paths_json": string(paths),
		"init_video_path":       req.InitVideoPath,
		"init_depth_path":       req.InitDepthPath,
		"init_point_cloud_path": req.InitPointCloudPath,
		"init_splat_path":       req.InitSplatPath,
		"mask_path":             req.MaskPath,
		"style_image_path":      req.StyleImagePath,
		"texture_image_path":    req.TextureImagePath,
		"animation_path":        req.AnimationPath,
	}
	applySpecDefaults(vars, spec)
	return vars
}

func voiceCreateHTTPVars(req driver.VoiceCreateRequest, spec *config.APICapabilitySpec) map[string]string {
	vars := map[string]string{
		"task":             string(driver.VoiceTaskZeroShot),
		"name":             req.Name,
		"sample_path":      req.SamplePath,
		"target_model":     req.TargetModel,
		"denoise":          fmt.Sprintf("%t", req.Denoise),
		"denoise_strength": fmt.Sprintf("%g", float64(req.DenoiseStrength)),
	}
	applySpecDefaults(vars, spec)
	return vars
}

func voiceConvertHTTPVars(req driver.VoiceConvertRequest, spec *config.APICapabilitySpec) map[string]string {
	vars := map[string]string{
		"task":             string(driver.VoiceTaskConversion),
		"name":             req.Name,
		"sample_path":      req.SamplePath,
		"target_model":     req.TargetModel,
		"source_path":      req.SourcePath,
		"pitch_shift":      fmt.Sprintf("%d", req.PitchShift),
		"index_rate":       fmt.Sprintf("%g", float64(req.IndexRate)),
		"protect":          fmt.Sprintf("%g", float64(req.Protect)),
		"vocoder":          req.Vocoder,
		"denoise":          fmt.Sprintf("%t", req.Denoise),
		"denoise_strength": fmt.Sprintf("%g", float64(req.DenoiseStrength)),
		"chunk_size":       fmt.Sprintf("%d", req.ChunkSize),
		"crossfade":        fmt.Sprintf("%g", float64(req.Crossfade)),
	}
	applySpecDefaults(vars, spec)
	return vars
}

func dataHTTPVars(req driver.DataRequest, spec *config.APICapabilitySpec) map[string]string {
	texts, _ := json.Marshal(req.Texts)
	vars := map[string]string{
		"task":          string(req.TaskOrDefault()),
		"model":         req.Model,
		"text":          req.Text,
		"texts_json":    string(texts),
		"text_file":     req.TextFile,
		"image_path":    req.ImagePath,
		"csv_path":      req.CSVPath,
		"graph_path":    req.GraphPath,
		"use_stdin":     boolStr(req.UseStdin),
		"input_format":  string(req.InputFormat),
		"input_shape":   string(driver.InferInputShape(req)),
		"output_shape":  string(driver.InferOutputShape(req)),
		"horizon":       fmt.Sprintf("%d", req.Horizon),
	}
	applySpecDefaults(vars, spec)
	return vars
}

func ragHTTPVars(req driver.RAGRequest, spec *config.APICapabilitySpec) map[string]string {
	sources, _ := json.Marshal(req.SourcePaths)
	excludes, _ := json.Marshal(req.Exclude)
	filter, _ := json.Marshal(req.Filter)
	meta, _ := json.Marshal(req.IndexMetadata)
	vars := map[string]string{
		"task":               string(req.Task),
		"store_root":         req.StoreRoot,
		"collection":         req.Collection,
		"query":              req.Query,
		"source_paths_json":  string(sources),
		"source_paths":       string(sources),
		"use_stdin":          boolStr(req.UseStdin),
		"recursive":          boolStr(req.Recursive),
		"chunk_size":         fmt.Sprintf("%d", req.ChunkSize),
		"chunk_overlap":      fmt.Sprintf("%d", req.ChunkOverlap),
		"embed_model":        req.EmbedModel,
		"text_model":         req.TextModel,
		"system_prompt":      req.SystemPrompt,
		"top_k":              fmt.Sprintf("%d", req.TopK),
		"min_score":          fmt.Sprintf("%g", float64(req.MinScore)),
		"max_tokens":         fmt.Sprintf("%d", req.MaxTokens),
		"filter_json":        string(filter),
		"purge_source":       req.PurgeSource,
		"index_mode":         string(req.IndexMode),
		"force":              boolStr(req.Force),
		"dry_run":            boolStr(req.DryRun),
		"index_metadata_json": string(meta),
		"glob":               req.Glob,
		"exclude_json":       string(excludes),
		"url":                req.URL,
		"query_file":         req.QueryFile,
		"query_stdin":        boolStr(req.QueryStdin),
		"rerank_top":         fmt.Sprintf("%d", req.RerankTop),
		"diverse":            boolStr(req.Diverse),
		"context_max_chars":  fmt.Sprintf("%d", req.ContextMaxChars),
		"cite":               boolStr(req.Cite),
		"temperature":        fmt.Sprintf("%g", float64(req.Temperature)),
		"top_p":              fmt.Sprintf("%g", float64(req.TopP)),
		"no_context":         boolStr(req.NoContext),
		"include_scores":     boolStr(req.IncludeScores),
		"rename_to":          req.RenameTo,
		"export_path":        req.ExportPath,
		"import_path":        req.ImportPath,
	}
	applySpecDefaults(vars, spec)
	return vars
}

func datasetHTTPVars(req driver.DatasetRequest, spec *config.APICapabilitySpec) map[string]string {
	task := req.Task
	if task == "" {
		task = req.Action
	}
	vars := map[string]string{
		"task":        string(task),
		"name":        req.Name,
		"path":        req.Path,
		"dataset_id":  req.DatasetID,
		"source_path": req.SourcePath,
		"recursive":   boolStr(req.Recursive),
		"format":      req.Format,
		"description": req.Description,
		"tag":         req.Tag,
		"message":     req.Message,
	}
	applySpecDefaults(vars, spec)
	return vars
}

func video2AudioHTTPVars(req driver.Video2AudioRequest, spec *config.APICapabilitySpec) map[string]string {
	vars := map[string]string{
		"task":       "extract",
		"video_path": req.VideoPath,
	}
	applySpecDefaults(vars, spec)
	return vars
}

func applySpecDefaults(vars map[string]string, spec *config.APICapabilitySpec) {
	if spec == nil || spec.Defaults == nil {
		return
	}
	d := spec.Defaults
	setDefault(vars, "model", d.Model)
	setDefault(vars, "max_tokens", intStr(d.MaxTokens))
	setDefault(vars, "temperature", floatStr(d.Temperature))
	setDefault(vars, "top_p", floatStr(d.TopP))
	setDefault(vars, "top_k", intStr(d.TopK))
	setDefault(vars, "min_p", floatStr(d.MinP))
	setDefault(vars, "cfg_scale", floatStr(d.CFGScale))
	setDefault(vars, "steps", intStr(d.Steps))
	setDefault(vars, "width", intStr(d.Width))
	setDefault(vars, "height", intStr(d.Height))
	setDefault(vars, "duration", floatStr(d.Duration))
	setDefault(vars, "fps", intStr(d.FPS))
	setDefault(vars, "sample_rate", intStr(d.SampleRate))
	setDefault(vars, "batch_size", intStr(d.BatchSize))
	setDefault(vars, "overlap", floatStr(d.Overlap))
	setDefault(vars, "scale", floatStr(d.Scale))
	setDefault(vars, "lora_rank", intStr(d.LoRARank))
	setDefault(vars, "learning_rate", floatStr(d.LearningRate))
	setDefault(vars, "epochs", intStr(d.Epochs))
}

func setDefault(vars map[string]string, key, val string) {
	if strings.TrimSpace(val) == "" || val == "0" || val == "0.000000" {
		return
	}
	if v, ok := vars[key]; !ok || strings.TrimSpace(v) == "" || v == "0" || v == "0.000000" {
		vars[key] = val
	}
}

func specModel(spec *config.APICapabilitySpec) string {
	if spec == nil {
		return ""
	}
	if spec.Model != "" {
		return spec.Model
	}
	if spec.Defaults != nil {
		return spec.Defaults.Model
	}
	return ""
}

func specInt(spec *config.APICapabilitySpec, field string) int {
	if spec == nil {
		return 0
	}
	if field == "max_tokens" && spec.MaxTokens > 0 {
		return spec.MaxTokens
	}
	if spec.Defaults != nil && spec.Defaults.MaxTokens > 0 {
		return spec.Defaults.MaxTokens
	}
	return 0
}

func specFloat(spec *config.APICapabilitySpec, field string) float64 {
	if spec == nil {
		return 0
	}
	if field == "temperature" && spec.Temperature > 0 {
		return spec.Temperature
	}
	if spec.Defaults != nil && spec.Defaults.Temperature > 0 {
		return spec.Defaults.Temperature
	}
	return 0
}

func specDefaultInt(spec *config.APICapabilitySpec, field string) int {
	if spec == nil || spec.Defaults == nil {
		return 0
	}
	switch field {
	case "width":
		return spec.Defaults.Width
	case "height":
		return spec.Defaults.Height
	case "steps":
		return spec.Defaults.Steps
	}
	return 0
}

func firstInt(a, b int) int {
	if a > 0 {
		return a
	}
	return b
}

func firstFloat(a, b float64) float64 {
	if a > 0 {
		return a
	}
	return b
}

func intStr(v int) string {
	if v == 0 {
		return ""
	}
	return fmt.Sprintf("%d", v)
}

func floatStr(v float64) string {
	if v == 0 {
		return ""
	}
	return fmt.Sprintf("%g", v)
}

func seedStr(seed *int) string {
	if seed == nil {
		return ""
	}
	return fmt.Sprintf("%d", *seed)
}
