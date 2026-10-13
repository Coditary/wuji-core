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
	ComfyUI     Type = "comfyui"
	MLX         Type = "mlx"
	PEFT        Type = "peft"
)

// Stable Diffusion / diffusion-model auxiliary weights.
const (
	ControlNet       Type = "controlnet"
	T2IAdapter       Type = "t2i_adapter"
	VAE              Type = "vae"
	TextualInversion Type = "textual_inversion"
	Hypernetwork     Type = "hypernetwork"
	IPAdapter        Type = "ip_adapter"
	ComfyUIWorkflow  Type = "comfyui_workflow"
	UpscalerModel    Type = "upscaler_model"
)

// Audio / speech model bundles.
const (
	WhisperCheckpoint Type = "whisper_checkpoint"
	BarkCheckpoint    Type = "bark_checkpoint"
	TTSPack           Type = "tts_pack"
)

// GPU-oriented quantized weight formats.
const (
	AWQ               Type = "awq"
	GPTQ              Type = "gptq"
	EXL2              Type = "exl2"
	EXL3              Type = "exl3"
	HQQ               Type = "hqq"
	HQQPlus           Type = "hqq_plus"
	AQLM              Type = "aqlm"
	Marlin            Type = "marlin"
	BitsAndBytes      Type = "bitsandbytes"
	NF4               Type = "nf4"
	FP4               Type = "fp4"
	CompressedTensors Type = "compressed_tensors"
	GGUFIQuant        Type = "gguf_iquant"
	SpQR              Type = "spqr"
	SqueezeLLM        Type = "squeezellm"
	QuIP              Type = "quip"
	BitNet            Type = "bitnet"
	VLLMCache         Type = "vllm_cache"
	GGUFShard         Type = "gguf_shard"
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
	TensorFlow       Type = "tensorflow"
	TensorFlowCKPT   Type = "tensorflow_ckpt"
	TFLite           Type = "tflite"
	EdgeTPU          Type = "edge_tpu"
	TFRecord         Type = "tfrecord"
	TensorFlowJS     Type = "tensorflowjs"
	Keras            Type = "keras"
)

// Compiler IR and intermediate representations.
const (
	TorchScript Type = "torchscript"
	TVM         Type = "tvm"
	IREE        Type = "iree"
	MLIR        Type = "mlir"
	Glow        Type = "glow"
	NNEF        Type = "nnef"
	SPIRV       Type = "spirv"
	LLVMBitcode Type = "llvm_bitcode"
)

// Pipeline, agent, and knowledge exchange formats.
const (
	OKF                Type = "okf"
	ChainSerialization Type = "chain_serialization"
	LangGraphState     Type = "langgraph_state"
	AgentYAML          Type = "agent_yaml"
	DSPyProgram        Type = "dspy_program"
	PromptfooConfig    Type = "promptfoo_config"
	SemanticKernel     Type = "semantic_kernel"
)

// Hardware/runtime-specific deployment formats.
const (
	TensorRT       Type = "tensorrt"
	CoreML         Type = "coreml"
	CoreMLPackage  Type = "coreml_package"
	MetalShader    Type = "metal_shader"
	ANE            Type = "ane"
	ORT            Type = "ort"
	OpenVINO       Type = "openvino"
	OpenVINOBlob   Type = "openvino_blob"
	NCNN           Type = "ncnn"
	TNN            Type = "tnn"
	MNN            Type = "mnn"
	RKNN           Type = "rknn"
	RKLLM          Type = "rkllm"
	SNPE           Type = "snpe"
	QNN            Type = "qnn"
	DLC            Type = "dlc"
	ExecuTorch     Type = "executorch"
	MediaPipe      Type = "mediapipe"
	ArmNN          Type = "armnn"
	NXPeIQ         Type = "nxp_eiq"
	NanoEdgeAI     Type = "nanoedge_ai"
	PopART         Type = "popart"
	Tenstorrent    Type = "tenstorrent"
	GroqModel      Type = "groq_model"
	Cerebras       Type = "cerebras"
	FPGABitstream  Type = "fpga_bitstream"
	MojoPackage    Type = "mojo_package"
)

// Framework-native and legacy deep learning formats.
const (
	Caffe           Type = "caffe"
	Caffe2          Type = "caffe2"
	MXNet           Type = "mxnet"
	Paddle          Type = "paddle"
	PaddleInference Type = "paddle_inference"
	PaddleParams    Type = "paddle_params"
	JAX             Type = "jax"
	Flax            Type = "flax"
	FlaxMsgpack     Type = "flax_msgpack"
	Orbax           Type = "orbax"
	Theano          Type = "theano"
	Lasagne         Type = "lasagne"
	Chainer         Type = "chainer"
	Darknet         Type = "darknet"
	CNTK            Type = "cntk"
	BigDL           Type = "bigdl"
	DL4J            Type = "dl4j"
	MatConvNet      Type = "matconvnet"
	WolframNet      Type = "wolfram_net"
	Torch7          Type = "torch7"
	DeepSpeed       Type = "deepspeed"
	MegatronLM      Type = "megatron_lm"
	Fairseq         Type = "fairseq"
	HFTrainerState  Type = "hf_trainer_state"
)

// Vector database and embedding index formats.
const (
	FAISSIndex   Type = "faiss_index"
	ChromaDB     Type = "chromadb"
	Milvus       Type = "milvus"
	Qdrant       Type = "qdrant"
	HNSWLib      Type = "hnswlib"
	Annoy        Type = "annoy"
	Lance        Type = "lance"
	DuckDB       Type = "duckdb"
)

// Training and dataset storage formats.
const (
	Parquet   Type = "parquet"
	Arrow     Type = "arrow"
	WebDataset Type = "webdataset"
	JSONL     Type = "jsonl"
	CSV       Type = "csv"
	TSV       Type = "tsv"
	Zarr      Type = "zarr"
)

// Classical / tabular ML serialization.
const (
	XGBoost     Type = "xgboost"
	LightGBM    Type = "lightgbm"
	CatBoost    Type = "catboost"
	Treelite    Type = "treelite"
	Joblib      Type = "joblib"
	PMML        Type = "pmml"
	PFA         Type = "pfa"
	Sklearn     Type = "sklearn"
	Statsmodels Type = "statsmodels"
	RapidMiner  Type = "rapidminer"
	KNIME       Type = "knime"
	Weka        Type = "weka"
	SPSS        Type = "spss"
	SAS         Type = "sas"
	Orange      Type = "orange"
)

// All returns every known format type in a stable order.
func All() []Type {
	return []Type{
		// Weight containers
		GGUF, GGML, GGJT, SafeTensors, PyTorch, Bin, Pickle, CKPT, HDF5, ONNX, Llamafile,
		// Distribution
		HuggingFace, Ollama, Diffusers, ComfyUI, MLX, PEFT,
		// Diffusion aux weights
		ControlNet, T2IAdapter, VAE, TextualInversion, Hypernetwork, IPAdapter, ComfyUIWorkflow, UpscalerModel,
		// Audio / speech
		WhisperCheckpoint, BarkCheckpoint, TTSPack,
		// Quantization
		AWQ, GPTQ, EXL2, EXL3, HQQ, HQQPlus, AQLM, Marlin, BitsAndBytes, NF4, FP4, CompressedTensors, GGUFIQuant, SpQR, SqueezeLLM, QuIP, BitNet, VLLMCache, GGUFShard,
		// Adapters
		LoRA, QLoRA, IA3, DoRA, LoHa,
		// TensorFlow
		TensorFlow, TensorFlowCKPT, TFLite, EdgeTPU, TFRecord, TensorFlowJS, Keras,
		// Compiler IR
		TorchScript, TVM, IREE, MLIR, Glow, NNEF, SPIRV, LLVMBitcode,
		// Pipeline / knowledge exchange
		OKF, ChainSerialization, LangGraphState, AgentYAML, DSPyProgram, PromptfooConfig, SemanticKernel,
		// Deployment / edge hardware
		TensorRT, CoreML, CoreMLPackage, MetalShader, ANE, ORT, OpenVINO, OpenVINOBlob, NCNN, TNN, MNN, RKNN, RKLLM, SNPE, QNN, DLC, ExecuTorch, MediaPipe, ArmNN, NXPeIQ, NanoEdgeAI, PopART, Tenstorrent, GroqModel, Cerebras, FPGABitstream, MojoPackage,
		// DL frameworks (incl. legacy)
		Caffe, Caffe2, MXNet, Paddle, PaddleInference, PaddleParams, JAX, Flax, FlaxMsgpack, Orbax, Theano, Lasagne, Chainer, Darknet, CNTK, BigDL, DL4J, MatConvNet, WolframNet, Torch7, DeepSpeed, MegatronLM, Fairseq, HFTrainerState,
		// Vector DB / embedding indexes
		FAISSIndex, ChromaDB, Milvus, Qdrant, HNSWLib, Annoy, Lance, DuckDB,
		// Training / dataset storage
		Parquet, Arrow, WebDataset, JSONL, CSV, TSV, Zarr,
		// Classical ML
		XGBoost, LightGBM, CatBoost, Treelite, Joblib, PMML, PFA, Sklearn, Statsmodels, RapidMiner, KNIME, Weka, SPSS, SAS, Orange,
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
	GGUF:              "Single-file LLM container with optional quantization (llama.cpp, Ollama, LM Studio)",
	GGML:              "Legacy LLM format, superseded by GGUF",
	GGJT:              "Early GGML tensor format (GGJT/GGJT_V2/V3), obsolete",
	SafeTensors:       "Secure tensor-only storage, Hugging Face default (.safetensors)",
	PyTorch:           "PyTorch native serialization (.pt, .pth)",
	Bin:               "Legacy weight blobs, often PyTorch pickle-based (.bin)",
	Pickle:            "Python pickle serialization (security risk on untrusted files)",
	CKPT:              "Classic monolithic checkpoints, common in older Stable Diffusion",
	HDF5:              "Keras / TensorFlow HDF5 full models and scientific arrays (.h5, .he5)",
	ONNX:              "Cross-framework graph + weights for portable deployment (.onnx)",
	Llamafile:         "Self-contained executable bundling GGUF weights and runtime",
	HuggingFace:       "Hugging Face Hub repo (sharded safetensors + config + tokenizer)",
	Ollama:            "Ollama registry tags (name:tag), not a direct file format",
	Diffusers:         "Hugging Face Diffusers pipeline directory layout",
	ComfyUI:           "ComfyUI models directory layout (checkpoints, LoRAs, ControlNet, VAE, embeddings)",
	MLX:               "Apple MLX model packages (usually safetensors + mlx config)",
	ControlNet:        "Stable Diffusion ControlNet adapter weights (.controlnet.safetensors, .pth)",
	T2IAdapter:        "Text-to-Image adapter conditioning modules (.safetensors)",
	VAE:               "Stable Diffusion VAE weights (.vae.safetensors, .vae.pt)",
	TextualInversion:  "Stable Diffusion textual inversion / embedding files (.pt, .safetensors, .bin)",
	Hypernetwork:      "Stable Diffusion hypernetwork fine-tuning weights (.pt)",
	IPAdapter:         "IP-Adapter image-prompt conditioning weights (.safetensors, .bin)",
	ComfyUIWorkflow:   "ComfyUI workflow graph definitions exported as JSON",
	UpscalerModel:     "AI upscaler weights such as Real-ESRGAN (.pth, .onnx)",
	WhisperCheckpoint: "OpenAI Whisper speech recognition checkpoints (.bin, .safetensors)",
	BarkCheckpoint:    "Suno Bark text-to-audio model checkpoints (.bin, .safetensors)",
	TTSPack:           "Text-to-speech voice packs (Piper ONNX, Kokoro PyTorch, etc.)",
	PEFT:              "Hugging Face PEFT adapter package layout (LoRA, QLoRA, etc.)",
	AWQ:               "Activation-aware 4-bit quantization for GPU inference servers",
	GPTQ:              "Post-training 4-bit quantization for GPU inference servers",
	EXL2:              "ExLlamaV2 mixed-bit quantization format for local GPU LLMs",
	EXL3:              "ExLlamaV3 mixed-bit quantization, optimized for fast token generation on NVIDIA GPUs",
	HQQ:               "Half-Quadratic Quantization for efficient LLM compression",
	HQQPlus:           "HQQ+ extended Half-Quadratic Quantization with optimized inference kernels",
	AQLM:              "Additive Quantization LM format (1–2 bit variants) for extreme compression",
	Marlin:            "Marlin-optimized INT4/FP8 weight layout for fast GPU kernels",
	BitsAndBytes:      "bitsandbytes 4-bit/8-bit quantized weights runtime wrapper",
	NF4:               "NormalFloat 4-bit quantization basis for normally distributed neural weights (QLoRA, bitsandbytes)",
	FP4:               "4-bit floating-point quantization basis used in bitsandbytes and QLoRA fine-tuning",
	CompressedTensors: "Hugging Face compressed-tensors format for FP8 and structured sparsity (llm-compressor, vLLM)",
	GGUFIQuant:        "GGUF importance-matrix quantizations (IQ1_S, IQ2_XXS, IQ3_S, …) for ultra-low-bit inference",
	SpQR:              "SpQR sparse-plus-quantized LLM weights for loss-aware compression",
	SqueezeLLM:        "SqueezeLLM sensitivity-based sparse quantization format",
	QuIP:              "QuIP / QuIP# lattice-based 2-bit and 4-bit LLM quantization",
	BitNet:            "BitNet 1.58-bit ternary LLM weight containers (-1, 0, +1)",
	VLLMCache:         "vLLM PagedAttention KV-cache and block memory layouts for inference servers",
	GGUFShard:         "Sharded GGUF weight parts (.gguf.part1, .gguf.part2, …) for large downloads",
	LoRA:              "Low-rank adaptation adapter weights",
	QLoRA:             "Quantized LoRA fine-tuning adapters",
	IA3:               "Infused Adapter by Inhibiting and Amplifying Inner Activations",
	DoRA:              "Weight-Decomposed Low-Rank Adaptation",
	LoHa:              "Low-Rank Hadamard product adaptation",
	TensorFlow:        "TensorFlow SavedModel protocol buffers (saved_model.pb)",
	TensorFlowCKPT:    "TensorFlow training checkpoints (.index, .data-00000-of-00001)",
	TFLite:            "TensorFlow Lite flatbuffers for mobile and embedded (.tflite)",
	EdgeTPU:           "Google Coral Edge TPU compiled TFLite models with EdgeTPU signature",
	TFRecord:          "TensorFlow training data record container format",
	TensorFlowJS:      "TensorFlow.js web model format (graph-model / layers-model)",
	Keras:             "Modern Keras v3 single-file ZIP model container (.keras)",
	TorchScript:       "Traced or scripted PyTorch models for C++ deployment (.pt)",
	TVM:               "Apache TVM compiled runtime modules",
	IREE:              "IREE compiler VM modules for portable deployment",
	MLIR:              "Multi-Level IR compiler representation",
	Glow:              "Facebook Glow graph compiler format",
	NNEF:              "Khronos Neural Network Exchange Format (.nnef), vendor-neutral alternative to ONNX",
	SPIRV:             "Standard Portable Intermediate Representation for GPU compute shaders (.spv)",
	LLVMBitcode:       "LLVM IR bitcode for compiler-generated AI kernel intermediate representation (.bc)",
	OKF:               "Open Knowledge Format for framework-independent RAG knowledge bases and document structures",
	ChainSerialization: "LangChain / LlamaIndex serialized chains, agent graphs, and prompt templates (JSON/YAML)",
	LangGraphState:    "LangGraph serialized agent state graphs (JSON)",
	AgentYAML:         "AutoGPT-style structured YAML configurations for autonomous agents",
	DSPyProgram:       "DSPy compiled prompt-optimization program state (.json)",
	PromptfooConfig:   "Promptfoo YAML test and evaluation configurations for LLM outputs",
	SemanticKernel:    "Microsoft Semantic Kernel plugin and planner configuration files",
	TensorRT:          "NVIDIA TensorRT optimized inference plans (.engine, .plan)",
	CoreML:            "Apple Core ML flat models for iOS and macOS (.mlmodel)",
	CoreMLPackage:     "Apple Core ML compiled model bundles (.mlpackage directory)",
	MetalShader:       "Apple Metal compute shaders and graph optimizations for Apple Silicon",
	ANE:               "Apple Neural Engine optimized model graphs and runtime artifacts",
	ORT:               "ONNX Runtime compiled models optimized for mobile and server inference (.ort)",
	OpenVINO:          "Intel OpenVINO intermediate representation (.xml + .bin)",
	OpenVINOBlob:      "Intel OpenVINO compiled blob for Movidius and Neural Compute Stick (.blob)",
	NCNN:              "Tencent NCNN param/bin format for efficient mobile inference",
	TNN:               "Tencent Neural Network inference format",
	MNN:               "Alibaba MNN mobile neural network format",
	RKNN:              "Rockchip NPU model format (.rknn)",
	RKLLM:             "Rockchip LLM deployment format",
	SNPE:              "Qualcomm Snapdragon Neural Processing Engine format",
	QNN:               "Qualcomm QNN (Qualcomm Neural Network) runtime container format",
	DLC:               "Qualcomm Deep Learning Container (.dlc)",
	ExecuTorch:        "Meta ExecuTorch programs for on-device PyTorch deployment",
	MediaPipe:         "Google MediaPipe task model bundles (.task)",
	ArmNN:             "Arm NN serialized inference graphs for ARM-optimized runtimes (.armnn)",
	NXPeIQ:            "NXP eIQ embedded inference data files for microcontrollers",
	NanoEdgeAI:        "STMicroelectronics NanoEdge AI models exported as C headers for IoT (.h)",
	PopART:            "Graphcore PopART / Poplar executables for IPU inference",
	Tenstorrent:       "Tenstorrent netlist and hardware accelerator configuration format",
	GroqModel:         "Groq LPU-optimized model binary sequences",
	Cerebras:          "Cerebras wafer-scale engine weight stream format",
	FPGABitstream:     "FPGA bitstreams for directly synthesized neural networks (.bit, .bin)",
	MojoPackage:       "Modular Mojo compiled AI modules (.mojopkg)",
	Caffe:             "Berkeley Caffe models (.caffemodel + .prototxt)",
	Caffe2:            "Facebook Caffe2 NetDef models (.pb, mostly legacy embedded deployments)",
	MXNet:             "Apache MXNet parameter files (.params)",
	Paddle:            "Baidu PaddlePaddle training checkpoints and export bundles",
	PaddleInference:   "PaddlePaddle inference graph (.pdmodel)",
	PaddleParams:      "PaddlePaddle inference weights (.pdiparams)",
	JAX:               "JAX checkpoint trees and export bundles",
	Flax:              "Flax linen checkpoint format",
	FlaxMsgpack:       "Flax msgpack-serialized weight checkpoints (.msgpack)",
	Orbax:             "Orbax checkpoint tree directories for modern JAX training",
	Theano:            "Legacy Theano shared variables (.pkl, .npz)",
	Lasagne:           "Lasagne / Theano layer stack pickle checkpoints",
	Chainer:           "Legacy Chainer NPZ checkpoints",
	Darknet:           "Darknet YOLO weights (.weights)",
	CNTK:              "Microsoft Cognitive Toolkit model format (deprecated)",
	BigDL:             "Apache Spark BigDL distributed deep learning models",
	DL4J:              "DeepLearning4J Java-native neural network bundles (.zip + JSON)",
	MatConvNet:        "MATLAB MatConvNet deep learning checkpoints (.mat)",
	WolframNet:        "Wolfram Neural Net Repository model representation (.wl)",
	Torch7:            "Original Lua Torch serialized models (.t7) before PyTorch",
	DeepSpeed:         "DeepSpeed ZeRO sharded distributed training checkpoint layout",
	MegatronLM:        "NVIDIA Megatron-LM tensor-parallel checkpoint directory layout",
	Fairseq:           "Meta Fairseq NLP research archives (.pt)",
	HFTrainerState:    "Hugging Face Trainer training metadata and history (trainer_state.json)",
	FAISSIndex:        "Facebook AI Similarity Search vector index files (.index)",
	ChromaDB:          "ChromaDB local vector store snapshots (SQLite / Parquet)",
	Milvus:            "Milvus distributed vector segment and binlog files",
	Qdrant:            "Qdrant vector database snapshot backups (.snapshot)",
	HNSWLib:           "Hierarchical Navigable Small World approximate nearest neighbor index",
	Annoy:             "Spotify Annoy approximate nearest neighbor index files (.ann)",
	Lance:             "Lance columnar format for vectors and ML datasets (.lance, LanceDB)",
	DuckDB:            "DuckDB embedded database files for metadata and embedding storage (.duckdb)",
	Parquet:           "Apache Parquet columnar format for LLM pretraining and tabular data",
	Arrow:             "Apache Arrow in-memory columnar format for fast data transfer (.arrow)",
	WebDataset:        "Tar-based WebDataset containers for large-scale vision and multimodal training",
	JSONL:             "JSON Lines fine-tuning and evaluation datasets (.jsonl)",
	CSV:               "Comma-separated tabular training and evaluation data",
	TSV:               "Tab-separated tabular training and evaluation data",
	Zarr:              "Chunked compressed N-dimensional arrays for large scientific datasets",
	XGBoost:           "XGBoost model binaries (.model, universal .ubj JSON format)",
	LightGBM:          "LightGBM model format (.txt text or .bin binary)",
	CatBoost:          "Yandex CatBoost native model binaries (.cbm)",
	Treelite:          "Treelite compiled decision-tree models as native libraries (.so, .dll, .dylib)",
	Joblib:            "scikit-learn joblib persistence (.joblib)",
	PMML:              "Predictive Model Markup Language XML standard for analytics interchange",
	PFA:               "Portable Format for Analytics JSON successor to PMML for realtime scoring",
	Sklearn:           "scikit-learn pickle-based estimators",
	Statsmodels:       "Statsmodels econometric and statistical model results (.pickle, .sats)",
	RapidMiner:        "RapidMiner process and model files (.rmp)",
	KNIME:             "KNIME workflow and PMML export bundles (.knime)",
	Weka:              "University of Waikato Weka model containers (.model)",
	SPSS:              "IBM SPSS predictive analytics model files (.spe, .xml)",
	SAS:               "SAS enterprise scoring code and model exports (.sas)",
	Orange:            "Orange Data Mining visual workflow canvases (.ows)",
}
