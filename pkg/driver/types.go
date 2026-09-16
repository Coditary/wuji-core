package driver

// TextInputMode identifies how input is supplied to wuji text.
type TextInputMode string

const (
	TextInputPrompt   TextInputMode = "prompt"
	TextInputImage    TextInputMode = "image"
	TextInputVideo    TextInputMode = "video"
	TextInputAudio    TextInputMode = "audio"
	TextInputDocument TextInputMode = "document"
)

// ChatRole identifies a message author in multi-turn / agent requests.
type ChatRole string

const (
	ChatRoleSystem    ChatRole = "system"
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
	ChatRoleTool      ChatRole = "tool"
)

// ChatMessage is one turn in a chat-style text request.
type ChatMessage struct {
	Role       ChatRole       `json:"role"`
	Content    string         `json:"content"`
	Name       string         `json:"name,omitempty"`
	ToolCallID string         `json:"tool_call_id,omitempty"`
	ToolCalls  []ChatToolCall `json:"tool_calls,omitempty"`
}

// ChatToolDef describes a tool for chat completions.
type ChatToolDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

// ChatToolCall is a model-requested tool invocation in a chat response.
type ChatToolCall struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Args string `json:"args"`
}

// TextRequest is the input for text generation.
type TextRequest struct {
	InputMode  TextInputMode `json:"input_mode,omitempty"`
	Prompt     string        `json:"prompt,omitempty"`
	MediaPath  string        `json:"media_path,omitempty"`

	// Chat/agent mode (optional; when set, used by API drivers and agent clients)
	Messages []ChatMessage `json:"messages,omitempty"`
	Tools    []ChatToolDef `json:"tools,omitempty"`
	Translate  bool   `json:"translate,omitempty"`
	TargetLang string `json:"target_lang,omitempty"`
	Language   string `json:"language,omitempty"` // source language: --audio or --translate (text source)

	// Audio transcription options (--audio only)
	BeamSize       int     `json:"beam_size,omitempty"`
	WordTimestamps bool    `json:"word_timestamps,omitempty"`
	VADEnabled     bool    `json:"vad_enabled,omitempty"`
	VADThreshold   float32 `json:"vad_threshold,omitempty"`

	Model        string  `json:"model,omitempty"`
	SystemPrompt string  `json:"system_prompt,omitempty"`
	MaxTokens    int     `json:"max_tokens,omitempty"`
	Temperature  float32 `json:"temperature,omitempty"`

	// Sampling (zero values = use backend default)
	TopP  float32 `json:"top_p,omitempty"`
	TopK  int     `json:"top_k,omitempty"`
	MinP  float32 `json:"min_p,omitempty"`

	// Penalties (zero values = use backend default)
	FrequencyPenalty  float32 `json:"frequency_penalty,omitempty"`
	PresencePenalty   float32 `json:"presence_penalty,omitempty"`
	RepetitionPenalty float32 `json:"repetition_penalty,omitempty"`

	// Limits
	StopSequences []string `json:"stop_sequences,omitempty"`

	// Advanced (Seed nil = random; ContextWindow 0 = backend default)
	Seed          *int   `json:"seed,omitempty"`
	ContextWindow int    `json:"context_window,omitempty"`
	ToolChoice    string `json:"tool_choice,omitempty"` // auto, required, none

	LoRAs []LoRARef `json:"loras,omitempty"`

	// Stream requests token-by-token delivery via OnDelta (server-side only).
	Stream bool `json:"stream,omitempty"`
	// OnDelta is invoked for each streamed text chunk. Not serialized over RPC.
	OnDelta func(delta string) error `json:"-"`
	// OnThinkDelta is invoked for each streamed reasoning/thinking chunk. Not serialized over RPC.
	OnThinkDelta func(delta string) error `json:"-"`
}

// TextResponse is the output of text generation.
type TextResponse struct {
	Text             string          `json:"text"`
	ThinkingBlocks   []ThinkingBlock `json:"thinking_blocks,omitempty"`
	TokensUsed       int             `json:"tokens_used"`
	InputTokens      int            `json:"input_tokens,omitempty"`
	LatencyMs        int64          `json:"latency_ms,omitempty"`
	TokensPerSec     float64        `json:"tokens_per_sec,omitempty"`
	EstimatedCostUSD float64        `json:"estimated_cost_usd,omitempty"`
	FinishReason     string         `json:"finish_reason,omitempty"`
	ToolCalls        []ChatToolCall `json:"tool_calls,omitempty"`
}

// Video2AudioRequest is the input for extracting audio from a video file.
type Video2AudioRequest struct {
	VideoPath string
}

// Video2AudioResponse is the output of video-to-audio extraction.
type Video2AudioResponse struct {
	AudioPath string
	Temporary bool // caller should delete AudioPath when done
}

// ImageRequest is the input for image generation.
type ImageRequest struct {
	Task           ImageTask
	Prompt         string
	NegativePrompt string
	Model          string
	Width          int
	Height         int
	Steps          int
	Sampler        string
	CFGScale       float32
	BatchSize      int
	BatchCount     int

	// Advanced (Seed nil = random; zero values = backend default)
	Seed              *int
	DenoisingStrength float32
	InitImagePath     string
	MaskImagePath     string
	ControlImagePath  string
	ControlType       ImageControlType
	ControlUnits      []ControlNetUnit
	ControlMode       ControlNetMode
	StyleImagePath    string
	StyleWeight       float32
	Scale             float32
	LoRAs             []LoRARef
	Mode              AssetMode
	FrameWidth         int
	FrameHeight        int
	Columns            int
	Rows               int
	FrameCount         int
	ReferenceImagePath string
	SpriteAction       SpriteAction
	SpriteView         SpriteView
	SpriteDirections   int
	SpriteLoop         bool
	SpritePadding      int
	SpriteTransparent  bool
}

// ImageResponse is the output of image generation.
type ImageResponse struct {
	Path   string
	Paths  []string
	Format string
}

// CameraControl identifies a preset camera movement.
type CameraControl string

const (
	CameraControlNone     CameraControl = ""
	CameraControlZoomIn   CameraControl = "zoom_in"
	CameraControlZoomOut  CameraControl = "zoom_out"
	CameraControlPanLeft  CameraControl = "pan_left"
	CameraControlPanRight CameraControl = "pan_right"
	CameraControlPanUp    CameraControl = "pan_up"
	CameraControlPanDown  CameraControl = "pan_down"
)

// VideoRequest is the CLI union input for all video tasks.
type VideoRequest struct {
	Task   VideoTask
	Prompt string
	Duration float32
	FPS      int
	Frames   int

	// Motion & dynamics (zero values = backend default)
	MotionStrength float32
	ContextLength  int

	// Methods & quality
	Sampler   string
	Scheduler string

	// Task-specific inputs
	InitImagePath string
	InitVideoPath string
	CameraControl CameraControl
	Scale         float32

	// Optional extras for backends
	NegativePrompt string
	Model          string
	Seed           *int
}

// EffectiveDuration returns video length in seconds from frames/fps or duration.
func (r VideoRequest) EffectiveDuration() float32 {
	if r.Frames > 0 && r.FPS > 0 {
		return float32(r.Frames) / float32(r.FPS)
	}
	return r.Duration
}

// VideoResponse is the output of video generation.
type VideoResponse struct {
	Path     string
	Duration float32
	Frames   int
	FPS      int
}

// AudioRequest is the CLI union input for all audio tasks.
type AudioRequest struct {
	Task           AudioTask
	Prompt         string
	Lyrics         string
	NegativePrompt string
	Model          string
	Duration       float32
	Overlap        float32
	Temperature    float32
	CFGScale       float32
	TopP           float32
	TopK           int
	SampleRate     int
	ReferencePath  string
	Voice          string
	Language       string
	Speed          float32
	Pitch          int
	Emotion        string
	Style          string
	Energy         float32
	Format         string
	Seed           *int

	// Deprecated: use Task instead.
	TaskType AudioTaskType
}

// AudioResponse is the output of audio generation.
type AudioResponse struct {
	Path       string
	Duration   float32
	SampleRate int
	Format     string
}

// MeshRequest is the CLI union input for all mesh tasks.
type MeshRequest struct {
	Task   MeshTask
	Prompt string
	Format string
	Model  string
	Seed   *int

	TargetTris int
	Scale      float32
	Mode       AssetMode
	Representation MeshRepresentation

	InitMeshPath       string
	HighMeshPath       string
	InitImagePath      string
	InitImagePaths     []string
	InitVideoPath      string
	InitDepthPath      string
	InitPointCloudPath string
	InitSplatPath      string
	MaskPath           string
	StyleImagePath     string
	TextureImagePath   string
	AnimationPath      string
}

// MeshResponse is the output of mesh generation.
type MeshResponse struct {
	Path           string
	Format         string
	Task           MeshTask
	Mode           AssetMode
	Representation MeshRepresentation
}

// VoiceRequest is the CLI union input for all voice tasks.
type VoiceRequest struct {
	Task            VoiceTask
	Name            string
	SamplePath      string
	TargetModel     string
	SourcePath      string
	PitchShift      int
	IndexRate       float32
	Protect         float32
	Vocoder         string
	Denoise         bool
	DenoiseStrength float32
	ChunkSize       int
	Crossfade       float32

	// Deprecated: use Task instead.
	Mode VoiceCloneMode
}

// VoiceResponse is the output of voice cloning.
type VoiceResponse struct {
	VoiceID    string
	Name       string
	OutputPath string
}

// DatasetRequest is the CLI union input for all dataset tasks.
type DatasetRequest struct {
	Task        DatasetTask
	Action      DatasetAction // deprecated alias
	Name        string
	Path        string
	DatasetID   string
	SourcePath  string
	Recursive   bool
	Format      string
	Description string
	Tag         string
	Message     string
}

// DatasetEntry describes a single dataset.
type DatasetEntry struct {
	ID            string
	Name          string
	Path          string
	Size          int64
	LatestVersion string
	VersionCount  int
}

// DatasetResponse is the output of dataset management.
type DatasetResponse struct {
	Datasets   []DatasetEntry
	Versions   []DatasetVersionEntry
	Message    string
	FilesAdded int
	BytesAdded int64
	VersionID  string
}
