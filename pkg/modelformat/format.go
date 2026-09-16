package modelformat

// Type identifies a model or weight source format a driver can load.
type Type string

// Weight containers and serialization formats.
const (
	GGUF        Type = "gguf"
	GGML        Type = "ggml"
	GGJT        Type = "ggjt"
	SafeTensors Type = "safetensors"
	PyTorch     Type = "pytorch"
	Bin         Type = "bin"
	Pickle      Type = "pickle"
	CKPT        Type = "ckpt"
	HDF5        Type = "hdf5"
	ONNX        Type = "onnx"
	Llamafile   Type = "llamafile"
)

// Distribution layouts and model hubs.
const (
	HuggingFace Type = "huggingface"
	Ollama      Type = "ollama"
	Diffusers   Type = "diffusers"
	MLX         Type = "mlx"
	PEFT        Type = "peft"
)

// GPU-oriented quantized weight formats.
const (
	AWQ          Type = "awq"
	GPTQ         Type = "gptq"
	EXL2         Type = "exl2"
	HQQ          Type = "hqq"
	AQLM         Type = "aqlm"
	Marlin       Type = "marlin"
	BitsAndBytes Type = "bitsandbytes"
)

// Fine-tuning adapters and parameter-efficient methods.
const (
	LoRA  Type = "lora"
	QLoRA Type = "qlora"
	IA3   Type = "ia3"
	DoRA  Type = "dora"
	LoHa  Type = "loha"
)

// TensorFlow family.
const (
	TensorFlow     Type = "tensorflow"
	TensorFlowCKPT Type = "tensorflow_ckpt"
	TFLite         Type = "tflite"
	TensorFlowJS   Type = "tensorflowjs"
)

// Compiler IR and intermediate representations.
const (
	TorchScript Type = "torchscript"
	TVM         Type = "tvm"
	IREE        Type = "iree"
	MLIR        Type = "mlir"
	Glow        Type = "glow"
)

// Hardware/runtime-specific deployment formats.
const (
	TensorRT   Type = "tensorrt"
	CoreML     Type = "coreml"
	OpenVINO   Type = "openvino"
	NCNN       Type = "ncnn"
	TNN        Type = "tnn"
	MNN        Type = "mnn"
	RKNN       Type = "rknn"
	RKLLM      Type = "rkllm"
	SNPE       Type = "snpe"
	DLC        Type = "dlc"
	ExecuTorch Type = "executorch"
	MediaPipe  Type = "mediapipe"
)

// Framework-native and legacy deep learning formats.
const (
	Caffe   Type = "caffe"
	MXNet   Type = "mxnet"
	Paddle  Type = "paddle"
	JAX     Type = "jax"
	Flax    Type = "flax"
	Theano  Type = "theano"
	Chainer Type = "chainer"
	Darknet Type = "darknet"
	CNTK    Type = "cntk"
)

// Classical / tabular ML serialization.
const (
	XGBoost  Type = "xgboost"
	LightGBM Type = "lightgbm"
	Joblib   Type = "joblib"
	PMML     Type = "pmml"
	Sklearn  Type = "sklearn"
)

// All returns every known format type in a stable order.
func All() []Type {
	return []Type{
		// Weight containers
		GGUF, GGML, GGJT, SafeTensors, PyTorch, Bin, Pickle, CKPT, HDF5, ONNX, Llamafile,
		// Distribution
		HuggingFace, Ollama, Diffusers, MLX, PEFT,
		// Quantization
		AWQ, GPTQ, EXL2, HQQ, AQLM, Marlin, BitsAndBytes,
		// Adapters
		LoRA, QLoRA, IA3, DoRA, LoHa,
		// TensorFlow
		TensorFlow, TensorFlowCKPT, TFLite, TensorFlowJS,
		// Compiler IR
		TorchScript, TVM, IREE, MLIR, Glow,
		// Deployment / edge hardware
		TensorRT, CoreML, OpenVINO, NCNN, TNN, MNN, RKNN, RKLLM, SNPE, DLC, ExecuTorch, MediaPipe,
		// DL frameworks (incl. legacy)
		Caffe, MXNet, Paddle, JAX, Flax, Theano, Chainer, Darknet, CNTK,
		// Classical ML
		XGBoost, LightGBM, Joblib, PMML, Sklearn,
	}
}

// Description returns a short human-readable explanation of the format.
func (t Type) Description() string {
	if desc, ok := descriptions[t]; ok {
		return desc
	}
	return ""
}

func (t Type) String() string {
	return string(t)
}

// IsKnown reports whether t is a registered format identifier.
func IsKnown(t Type) bool {
	_, ok := descriptions[t]
	return ok
}

var descriptions = map[Type]string{
	GGUF:           "Single-file LLM container with optional quantization (llama.cpp, Ollama, LM Studio)",
	GGML:           "Legacy LLM format, superseded by GGUF",
	GGJT:           "Early GGML tensor format (GGJT/GGJT_V2/V3), obsolete",
	SafeTensors:    "Secure tensor-only storage, Hugging Face default (.safetensors)",
	PyTorch:        "PyTorch native serialization (.pt, .pth)",
	Bin:            "Legacy weight blobs, often PyTorch pickle-based (.bin)",
	Pickle:         "Python pickle serialization (security risk on untrusted files)",
	CKPT:           "Classic monolithic checkpoints, common in older Stable Diffusion",
	HDF5:           "Keras / TensorFlow HDF5 full models (.h5)",
	ONNX:           "Cross-framework graph + weights for portable deployment (.onnx)",
	Llamafile:      "Self-contained executable bundling GGUF weights and runtime",
	HuggingFace:    "Hugging Face Hub repo (sharded safetensors + config + tokenizer)",
	Ollama:         "Ollama registry tags (name:tag), not a direct file format",
	Diffusers:      "Hugging Face Diffusers pipeline directory layout",
	MLX:            "Apple MLX model packages (usually safetensors + mlx config)",
	PEFT:           "Hugging Face PEFT adapter package layout (LoRA, QLoRA, etc.)",
	AWQ:            "Activation-aware 4-bit quantization for GPU inference servers",
	GPTQ:           "Post-training 4-bit quantization for GPU inference servers",
	EXL2:           "ExLlamaV2 mixed-bit quantization format for local GPU LLMs",
	HQQ:            "Half-Quadratic Quantization for efficient LLM compression",
	AQLM:           "Additive Quantization LM format for extreme compression",
	Marlin:         "Marlin-optimized INT4/FP8 weight layout for fast GPU kernels",
	BitsAndBytes:   "bitsandbytes 4-bit/8-bit quantized weights (NF4, FP4)",
	LoRA:           "Low-rank adaptation adapter weights",
	QLoRA:          "Quantized LoRA fine-tuning adapters",
	IA3:            "Infused Adapter by Inhibiting and Amplifying Inner Activations",
	DoRA:           "Weight-Decomposed Low-Rank Adaptation",
	LoHa:           "Low-Rank Hadamard product adaptation",
	TensorFlow:     "TensorFlow SavedModel or frozen GraphDef (.pb)",
	TensorFlowCKPT: "TensorFlow training checkpoints (.ckpt, .data, .index)",
	TFLite:         "TensorFlow Lite flatbuffers for mobile and embedded (.tflite)",
	TensorFlowJS:   "TensorFlow.js web model format (graph-model / layers-model)",
	TorchScript:    "Traced or scripted PyTorch models for C++ deployment (.pt)",
	TVM:            "Apache TVM compiled runtime modules",
	IREE:           "IREE compiler VM modules for portable deployment",
	MLIR:           "Multi-Level IR compiler representation",
	Glow:           "Facebook Glow graph compiler format",
	TensorRT:       "NVIDIA TensorRT optimized engines (.engine)",
	CoreML:         "Apple Core ML models for iOS and macOS (.mlmodel)",
	OpenVINO:       "Intel OpenVINO intermediate representation (.xml + .bin)",
	NCNN:           "Tencent NCNN param/bin format for efficient mobile inference",
	TNN:            "Tencent Neural Network inference format",
	MNN:            "Alibaba MNN mobile neural network format",
	RKNN:           "Rockchip NPU model format (.rknn)",
	RKLLM:          "Rockchip LLM deployment format",
	SNPE:           "Qualcomm Snapdragon Neural Processing Engine format",
	DLC:            "Qualcomm Deep Learning Container (.dlc)",
	ExecuTorch:     "Meta ExecuTorch programs for on-device PyTorch deployment",
	MediaPipe:      "Google MediaPipe task model bundles",
	Caffe:          "Berkeley Caffe models (.caffemodel + .prototxt)",
	MXNet:          "Apache MXNet parameter files (.params)",
	Paddle:         "PaddlePaddle inference model format",
	JAX:            "JAX/Flax checkpoint trees (often msgpack or orbax)",
	Flax:           "Flax linen checkpoint format",
	Theano:         "Legacy Theano shared variables (.pkl, .npz)",
	Chainer:        "Legacy Chainer NPZ checkpoints",
	Darknet:        "Darknet YOLO weights (.weights)",
	CNTK:           "Microsoft Cognitive Toolkit model format (deprecated)",
	XGBoost:        "XGBoost model binaries (.model, .ubj)",
	LightGBM:       "LightGBM text model format (.txt)",
	Joblib:         "scikit-learn joblib persistence (.joblib)",
	PMML:           "Predictive Model Markup Language XML standard",
	Sklearn:        "scikit-learn pickle-based estimators",
}
