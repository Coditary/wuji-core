package modelformat

// MergeTypes returns a de-duplicated list preserving first-seen order.
func MergeTypes(groups ...[]Type) []Type {
	seen := make(map[Type]struct{})
	out := make([]Type, 0)
	for _, group := range groups {
		for _, t := range group {
			if _, ok := seen[t]; ok {
				continue
			}
			seen[t] = struct{}{}
			out = append(out, t)
		}
	}
	return out
}

// DiffusionAuxFormats lists Stable Diffusion auxiliary weight types.
func DiffusionAuxFormats() []Type {
	return []Type{ControlNet, T2IAdapter, VAE, TextualInversion, Hypernetwork, IPAdapter}
}

// DiffusionCheckpointFormats lists common SD base checkpoint layouts.
func DiffusionCheckpointFormats() []Type {
	return []Type{SafeTensors, CKPT, Diffusers, PyTorch, LoRA}
}

// DiffusionImageFormats lists SD WebUI-style image generation model sources.
func DiffusionImageFormats() []Type {
	return MergeTypes(DiffusionCheckpointFormats(), DiffusionAuxFormats())
}

// ComfyUIFormats extends diffusion image formats with ComfyUI-specific layouts.
func ComfyUIFormats() []Type {
	return MergeTypes(DiffusionImageFormats(), []Type{ComfyUI, ComfyUIWorkflow})
}

// WebUIVideoFormats lists common SD video model sources.
func WebUIVideoFormats() []Type {
	return []Type{SafeTensors, CKPT, Diffusers, PyTorch}
}

// GGUFEcosystemFormats lists llama.cpp / Ollama local LLM sources.
func GGUFEcosystemFormats() []Type {
	return []Type{GGUF, GGML, GGUFIQuant, GGUFShard, Ollama, LoRA}
}

// HuggingFaceServerFormats lists common HF inference-server model sources.
func HuggingFaceServerFormats() []Type {
	return []Type{HuggingFace, SafeTensors, PyTorch, Bin, AWQ, GPTQ, CompressedTensors, LoRA}
}

// ExLlamaFormats lists ExLlama-family local GPU LLM sources.
func ExLlamaFormats() []Type {
	return []Type{EXL2, EXL3, AWQ, GPTQ, SafeTensors, HuggingFace, LoRA}
}

// UpscalerFormats lists AI image upscaler model sources.
func UpscalerFormats() []Type {
	return []Type{ONNX, PyTorch, UpscalerModel}
}

// TTSONNXFormats lists ONNX-based TTS voice/model sources.
func TTSONNXFormats() []Type {
	return []Type{ONNX, TTSPack}
}

// TTSPyTorchFormats lists PyTorch-based TTS model sources.
func TTSPyTorchFormats() []Type {
	return []Type{PyTorch, GGUF, TTSPack, BarkCheckpoint}
}

// RAGEmbedFormats lists embedding backends commonly used for RAG indexing.
func RAGEmbedFormats() []Type {
	return []Type{Ollama, GGUF, HuggingFace, SafeTensors}
}
