package driver

// AudioRuntimeID identifies a TTS runtime driver (inference backend).
const (
	AudioRuntimeONNX    = "onnx-tts"
	AudioRuntimePyTorch = "pytorch-tts"
)

// AudioBrandAlias maps a discoverable brand driver ID to a runtime + pipeline/engine.
type AudioBrandAlias struct {
	RuntimeID string
	Pipeline  string // onnx-tts: piper, kokoro
	Engine    string // pytorch-tts: coqui, bark, …
}

// AudioBrandAliases lists brand names that resolve to runtime drivers when the runtime is registered.
var AudioBrandAliases = map[string]AudioBrandAlias{
	"piper":       {RuntimeID: AudioRuntimeONNX, Pipeline: "piper"},
	"kokoro":      {RuntimeID: AudioRuntimeONNX, Pipeline: "kokoro"},
	"coqui":       {RuntimeID: AudioRuntimePyTorch, Engine: "coqui"},
	"melotts":     {RuntimeID: AudioRuntimePyTorch, Engine: "melotts"},
	"fishspeech":  {RuntimeID: AudioRuntimePyTorch, Engine: "fishspeech"},
	"bark":        {RuntimeID: AudioRuntimePyTorch, Engine: "bark"},
	"styletts2":   {RuntimeID: AudioRuntimePyTorch, Engine: "styletts2"},
	"chattts":     {RuntimeID: AudioRuntimePyTorch, Engine: "chattts"},
	"parler":      {RuntimeID: AudioRuntimePyTorch, Engine: "parler"},
	"f5tts":       {RuntimeID: AudioRuntimePyTorch, Engine: "f5tts"},
	"cosyvoice":   {RuntimeID: AudioRuntimePyTorch, Engine: "cosyvoice"},
	"gptsovits":   {RuntimeID: AudioRuntimePyTorch, Engine: "gptsovits"},
	"orpheus":     {RuntimeID: AudioRuntimePyTorch, Engine: "orpheus"},
	"chatterbox":  {RuntimeID: AudioRuntimePyTorch, Engine: "chatterbox"},
	"zonos":       {RuntimeID: AudioRuntimePyTorch, Engine: "zonos"},
}

// ApplyAudioBrandAlias rewrites driver ID and request when a brand alias targets a registered runtime.
// Returns the driver ID to use and whether the alias was applied.
func ApplyAudioBrandAlias(driverID string, req *AudioRequest, runtimeRegistered func(string) bool) (string, bool) {
	alias, ok := AudioBrandAliases[driverID]
	if !ok {
		return driverID, false
	}
	if runtimeRegistered == nil || !runtimeRegistered(alias.RuntimeID) {
		return driverID, false
	}
	if alias.Pipeline != "" {
		req.Pipeline = alias.Pipeline
	}
	if alias.Engine != "" {
		req.Engine = alias.Engine
	}
	return alias.RuntimeID, true
}
