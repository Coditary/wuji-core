package driver

import (
	"fmt"
	"strings"
)

// TrainMethodInfo describes a training method and the data it expects.
type TrainMethodInfo struct {
	Method      string
	Description string
	Dataset     string // human-readable dataset requirements
}

// --- Text ---

type TextTrainMethod string

const (
	TextTrainMethodFull  TextTrainMethod = "full"
	TextTrainMethodLoRA  TextTrainMethod = "lora"
	TextTrainMethodQLoRA TextTrainMethod = "qlora"
)

func AllTextTrainMethods() []TextTrainMethod {
	return []TextTrainMethod{TextTrainMethodFull, TextTrainMethodLoRA, TextTrainMethodQLoRA}
}

func TextTrainMethodCatalog() []TrainMethodInfo {
	return []TrainMethodInfo{
		{Method: string(TextTrainMethodFull), Description: "Full fine-tune of the base LLM", Dataset: "JSONL/JSON with instruction–response or chat turns; one example per line"},
		{Method: string(TextTrainMethodLoRA), Description: "Low-rank adapter on top of a frozen base model", Dataset: "Same as full; smaller rank (--lora-rank) reduces VRAM"},
		{Method: string(TextTrainMethodQLoRA), Description: "Quantized LoRA (4-bit base + adapter)", Dataset: "Same as LoRA; lowest VRAM; set --lora-rank"},
	}
}

func ParseTextTrainMethod(raw string) (TextTrainMethod, error) {
	normalized := normalizeTrainMethod(raw)
	switch normalized {
	case "", "full", "finetune", "fine-tune", "sft":
		return TextTrainMethodFull, nil
	case "lora":
		return TextTrainMethodLoRA, nil
	case "qlora":
		return TextTrainMethodQLoRA, nil
	}
	return "", fmt.Errorf("unknown text train method %q (valid: %s)", raw, joinTextTrainMethods())
}

func joinTextTrainMethods() string {
	parts := make([]string, len(AllTextTrainMethods()))
	for i, m := range AllTextTrainMethods() {
		parts[i] = string(m)
	}
	return strings.Join(parts, ", ")
}

// --- Image ---

type ImageTrainMethod string

const (
	ImageTrainMethodFull             ImageTrainMethod = "full"
	ImageTrainMethodLoRA             ImageTrainMethod = "lora"
	ImageTrainMethodDreamBooth       ImageTrainMethod = "dreambooth"
	ImageTrainMethodTextualInversion ImageTrainMethod = "textual-inversion"
	ImageTrainMethodControlNet       ImageTrainMethod = "controlnet"
	ImageTrainMethodEmbedding        ImageTrainMethod = "embedding"
)

func AllImageTrainMethods() []ImageTrainMethod {
	return []ImageTrainMethod{
		ImageTrainMethodFull,
		ImageTrainMethodLoRA,
		ImageTrainMethodDreamBooth,
		ImageTrainMethodTextualInversion,
		ImageTrainMethodControlNet,
		ImageTrainMethodEmbedding,
	}
}

func ImageTrainMethodCatalog() []TrainMethodInfo {
	return []TrainMethodInfo{
		{Method: string(ImageTrainMethodFull), Description: "Fine-tune the full diffusion UNet (+ optional text encoder)", Dataset: "Folder of images with captions (.txt sidecars or metadata.json); 10–1000+ pairs"},
		{Method: string(ImageTrainMethodLoRA), Description: "Train a small LoRA adapter", Dataset: "Captioned image folder; 5–200 images typical"},
		{Method: string(ImageTrainMethodDreamBooth), Description: "DreamBooth subject fine-tune with a class token", Dataset: "3–30 subject photos + --class-token; optional regularization set"},
		{Method: string(ImageTrainMethodTextualInversion), Description: "Learn a new text embedding for a concept", Dataset: "3–20 concept images; set --token for the trigger word"},
		{Method: string(ImageTrainMethodControlNet), Description: "Train a ControlNet conditioning model", Dataset: "Image pairs: source frame + matching control map (pose, depth, canny, …); set --control-type"},
		{Method: string(ImageTrainMethodEmbedding), Description: "IP-Adapter / style embedding from reference images", Dataset: "Folder of style or reference images (no captions required for some backends)"},
	}
}

func ParseImageTrainMethod(raw string) (ImageTrainMethod, error) {
	normalized := normalizeTrainMethod(raw)
	switch normalized {
	case "", "full", "finetune", "fine-tune":
		return ImageTrainMethodFull, nil
	case "lora":
		return ImageTrainMethodLoRA, nil
	case "dreambooth", "db":
		return ImageTrainMethodDreamBooth, nil
	case "textual-inversion", "ti", "inversion":
		return ImageTrainMethodTextualInversion, nil
	case "controlnet", "control":
		return ImageTrainMethodControlNet, nil
	case "embedding", "ip-adapter", "ipadapter", "style":
		return ImageTrainMethodEmbedding, nil
	}
	return "", fmt.Errorf("unknown image train method %q (valid: %s)", raw, joinImageTrainMethods())
}

func joinImageTrainMethods() string {
	parts := make([]string, len(AllImageTrainMethods()))
	for i, m := range AllImageTrainMethods() {
		parts[i] = string(m)
	}
	return strings.Join(parts, ", ")
}

// --- Video ---

type VideoTrainMethod string

const (
	VideoTrainMethodFull VideoTrainMethod = "full"
	VideoTrainMethodLoRA VideoTrainMethod = "lora"
)

func AllVideoTrainMethods() []VideoTrainMethod {
	return []VideoTrainMethod{VideoTrainMethodFull, VideoTrainMethodLoRA}
}

func VideoTrainMethodCatalog() []TrainMethodInfo {
	return []TrainMethodInfo{
		{Method: string(VideoTrainMethodFull), Description: "Fine-tune a video diffusion model", Dataset: "Clips folder with captions; set --train-frames and --train-fps to match clip length"},
		{Method: string(VideoTrainMethodLoRA), Description: "Temporal LoRA on a frozen video model", Dataset: "Short captioned clips (2–8 s); motion/style LoRA"},
	}
}

func ParseVideoTrainMethod(raw string) (VideoTrainMethod, error) {
	normalized := normalizeTrainMethod(raw)
	switch normalized {
	case "", "full", "finetune", "fine-tune":
		return VideoTrainMethodFull, nil
	case "lora":
		return VideoTrainMethodLoRA, nil
	}
	return "", fmt.Errorf("unknown video train method %q (valid: %s)", raw, joinVideoTrainMethods())
}

func joinVideoTrainMethods() string {
	parts := make([]string, len(AllVideoTrainMethods()))
	for i, m := range AllVideoTrainMethods() {
		parts[i] = string(m)
	}
	return strings.Join(parts, ", ")
}

// --- Audio ---

type AudioTrainMethod string

const (
	AudioTrainMethodMusic  AudioTrainMethod = "music"
	AudioTrainMethodMelody AudioTrainMethod = "melody"
	AudioTrainMethodSFX    AudioTrainMethod = "sfx"
	AudioTrainMethodSpeech AudioTrainMethod = "speech"
)

func AllAudioTrainMethods() []AudioTrainMethod {
	return []AudioTrainMethod{
		AudioTrainMethodMusic,
		AudioTrainMethodMelody,
		AudioTrainMethodSFX,
		AudioTrainMethodSpeech,
	}
}

func AudioTrainMethodCatalog() []TrainMethodInfo {
	return []TrainMethodInfo{
		{Method: string(AudioTrainMethodMusic), Description: "Text-to-music generation model", Dataset: "Audio files (WAV/MP3) with text captions or tags"},
		{Method: string(AudioTrainMethodMelody), Description: "Melody-conditioned music model", Dataset: "Melody audio + target music pairs"},
		{Method: string(AudioTrainMethodSFX), Description: "Sound-effects model", Dataset: "Short SFX clips with text descriptions"},
		{Method: string(AudioTrainMethodSpeech), Description: "Text-to-speech / voice model", Dataset: "Clean speech recordings with transcripts (one speaker per dataset)"},
	}
}

func ParseAudioTrainMethod(raw string) (AudioTrainMethod, error) {
	normalized := normalizeTrainMethod(raw)
	switch normalized {
	case "", "music", "text-to-music", "txt2music":
		return AudioTrainMethodMusic, nil
	case "melody", "melody-to-music", "melody2music":
		return AudioTrainMethodMelody, nil
	case "sfx", "text-to-sfx", "sound-effects":
		return AudioTrainMethodSFX, nil
	case "speech", "tts", "text-to-speech":
		return AudioTrainMethodSpeech, nil
	}
	return "", fmt.Errorf("unknown audio train method %q (valid: %s)", raw, joinAudioTrainMethods())
}

func (m AudioTrainMethod) ToAudioTask() AudioTask {
	switch m {
	case AudioTrainMethodMelody:
		return AudioTaskMelodyToMusic
	case AudioTrainMethodSFX:
		return AudioTaskTextToSFX
	case AudioTrainMethodSpeech:
		return AudioTaskTextToSpeech
	default:
		return AudioTaskTextToMusic
	}
}

func joinAudioTrainMethods() string {
	parts := make([]string, len(AllAudioTrainMethods()))
	for i, m := range AllAudioTrainMethods() {
		parts[i] = string(m)
	}
	return strings.Join(parts, ", ")
}

// --- Mesh ---

type MeshTrainMethod string

const (
	MeshTrainMethodFull MeshTrainMethod = "full"
	MeshTrainMethodLoRA MeshTrainMethod = "lora"
)

func AllMeshTrainMethods() []MeshTrainMethod {
	return []MeshTrainMethod{MeshTrainMethodFull, MeshTrainMethodLoRA}
}

func MeshTrainMethodCatalog() []TrainMethodInfo {
	return []TrainMethodInfo{
		{Method: string(MeshTrainMethodFull), Description: "Fine-tune a 3D generation model", Dataset: "Meshes (OBJ/GLB/PLY) with optional text captions or multi-view renders"},
		{Method: string(MeshTrainMethodLoRA), Description: "LoRA adapter for mesh/style generation", Dataset: "Small curated mesh set with captions; set --train-format for output"},
	}
}

func ParseMeshTrainMethod(raw string) (MeshTrainMethod, error) {
	normalized := normalizeTrainMethod(raw)
	switch normalized {
	case "", "full", "finetune", "fine-tune":
		return MeshTrainMethodFull, nil
	case "lora":
		return MeshTrainMethodLoRA, nil
	}
	return "", fmt.Errorf("unknown mesh train method %q (valid: %s)", raw, joinMeshTrainMethods())
}

func joinMeshTrainMethods() string {
	parts := make([]string, len(AllMeshTrainMethods()))
	for i, m := range AllMeshTrainMethods() {
		parts[i] = string(m)
	}
	return strings.Join(parts, ", ")
}

// --- Voice ---

type VoiceTrainMethod string

const (
	VoiceTrainMethodRVC  VoiceTrainMethod = "rvc"
	VoiceTrainMethodFull VoiceTrainMethod = "full"
)

func AllVoiceTrainMethods() []VoiceTrainMethod {
	return []VoiceTrainMethod{VoiceTrainMethodRVC, VoiceTrainMethodFull}
}

func VoiceTrainMethodCatalog() []TrainMethodInfo {
	return []TrainMethodInfo{
		{Method: string(VoiceTrainMethodRVC), Description: "RVC voice model from samples", Dataset: "5–30 min clean speech WAVs in --dataset; optional --sample for a reference clip"},
		{Method: string(VoiceTrainMethodFull), Description: "Full voice-cloning model fine-tune", Dataset: "Larger multi-speaker or single-speaker corpus with transcripts"},
	}
}

func ParseVoiceTrainMethod(raw string) (VoiceTrainMethod, error) {
	normalized := normalizeTrainMethod(raw)
	switch normalized {
	case "", "rvc", "retrieval-based-voice-conversion":
		return VoiceTrainMethodRVC, nil
	case "full", "finetune", "fine-tune":
		return VoiceTrainMethodFull, nil
	}
	return "", fmt.Errorf("unknown voice train method %q (valid: %s)", raw, joinVoiceTrainMethods())
}

func joinVoiceTrainMethods() string {
	parts := make([]string, len(AllVoiceTrainMethods()))
	for i, m := range AllVoiceTrainMethods() {
		parts[i] = string(m)
	}
	return strings.Join(parts, ", ")
}

func normalizeTrainMethod(raw string) string {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	return strings.ReplaceAll(normalized, "_", "-")
}
