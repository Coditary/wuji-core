package config

// APIEntry defines a YAML-configured remote API backend.
// Each entry becomes a virtual in-process driver (no separate binary).
type APIEntry struct {
	Text        *APICapabilitySpec `yaml:"text,omitempty"`
	Image       *APICapabilitySpec `yaml:"image,omitempty"`
	Video       *APICapabilitySpec `yaml:"video,omitempty"`
	Audio       *APICapabilitySpec `yaml:"audio,omitempty"`
	Mesh        *APICapabilitySpec `yaml:"mesh,omitempty"`
	Voice       *APICapabilitySpec `yaml:"voice,omitempty"`
	Video2Audio *APICapabilitySpec `yaml:"video2audio,omitempty"`
	Train       *APICapabilitySpec `yaml:"train,omitempty"`
	Dataset     *APICapabilitySpec `yaml:"dataset,omitempty"`
	Data        *APICapabilitySpec `yaml:"data,omitempty"`
	RAG         *APICapabilitySpec `yaml:"rag,omitempty"`
}

// APICapabilitySpec configures one capability for an API driver.
type APICapabilitySpec struct {
	Protocol string `yaml:"protocol"`

	Method  string            `yaml:"method,omitempty"`
	URL     string            `yaml:"url,omitempty"`
	BaseURL string            `yaml:"base_url,omitempty"`
	Query   map[string]string `yaml:"query,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`

	// JSON body as key/value templates, or use body_raw for free-form JSON.
	Body    map[string]string `yaml:"body,omitempty"`
	BodyRaw string            `yaml:"body_raw,omitempty"`

	ContentType string `yaml:"content_type,omitempty"`

	// Response JSON paths: text, image_url, image_base64, audio_url, video_url, …
	Response  map[string]string `yaml:"response,omitempty"`
	ErrorPath string            `yaml:"error_path,omitempty"`

	Auth      *APIAuthConfig      `yaml:"auth,omitempty"`
	Defaults  *APIDefaults        `yaml:"defaults,omitempty"`
	Multipart []APIMultipartField `yaml:"multipart,omitempty"`
	Poll      *APIPollConfig      `yaml:"poll,omitempty"`

	Streaming bool `yaml:"streaming,omitempty"`
	// Think enables model reasoning for backends that support it.
	// For vLLM/Ollama/Gemma via OpenAI-compatible API, use "low", "medium", or "high"
	// (mapped to reasoning_effort). true maps to "low"; false/none disables thinking.
	Think any `yaml:"think,omitempty"`
	// EnsureDriver starts a managed inference driver (e.g. vllm) before HTTP calls.
	EnsureDriver   string `yaml:"ensure_driver,omitempty"`
	TimeoutSeconds int    `yaml:"timeout_seconds,omitempty"`

	APIKey    string `yaml:"api_key,omitempty"`
	APIKeyEnv string `yaml:"api_key_env,omitempty"`
	Model     string `yaml:"model,omitempty"`

	AnthropicVersion string  `yaml:"anthropic_version,omitempty"`
	MaxTokens        int     `yaml:"max_tokens,omitempty"`
	Temperature      float64 `yaml:"temperature,omitempty"`

	// Per-task overrides (merged over this spec). Keys are driver task names.
	Tasks map[string]*APICapabilitySpec `yaml:"tasks,omitempty"`
}

// APIDefaults are YAML defaults merged when request fields are zero/empty.
type APIDefaults struct {
	Model             string  `yaml:"model,omitempty"`
	MaxTokens         int     `yaml:"max_tokens,omitempty"`
	Temperature       float64 `yaml:"temperature,omitempty"`
	TopP              float64 `yaml:"top_p,omitempty"`
	TopK              int     `yaml:"top_k,omitempty"`
	MinP              float64 `yaml:"min_p,omitempty"`
	CFGScale          float64 `yaml:"cfg_scale,omitempty"`
	FrequencyPenalty  float64 `yaml:"frequency_penalty,omitempty"`
	PresencePenalty   float64 `yaml:"presence_penalty,omitempty"`
	RepetitionPenalty float64 `yaml:"repetition_penalty,omitempty"`
	Steps             int     `yaml:"steps,omitempty"`
	Width             int     `yaml:"width,omitempty"`
	Height            int     `yaml:"height,omitempty"`
	Duration          float64 `yaml:"duration,omitempty"`
	FPS               int     `yaml:"fps,omitempty"`
	SampleRate        int     `yaml:"sample_rate,omitempty"`
	BatchSize         int     `yaml:"batch_size,omitempty"`
	Overlap           float64 `yaml:"overlap,omitempty"`
	Scale             float64 `yaml:"scale,omitempty"`
	LoRARank          int     `yaml:"lora_rank,omitempty"`
	LearningRate      float64 `yaml:"learning_rate,omitempty"`
	Epochs            int     `yaml:"epochs,omitempty"`
}

// APIAuthConfig controls how API keys are sent.
// type: bearer (default when api_key set), header, query, none
type APIAuthConfig struct {
	Type       string `yaml:"type,omitempty"`
	Header     string `yaml:"header,omitempty"`
	QueryParam string `yaml:"query_param,omitempty"`
	Prefix     string `yaml:"prefix,omitempty"`
}

// APIMultipartField is one part of a multipart/form-data request.
type APIMultipartField struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value,omitempty"`
	File  string `yaml:"file,omitempty"`
}

// APIPollConfig polls async job APIs until completion.
type APIPollConfig struct {
	IntervalSeconds int      `yaml:"interval_seconds,omitempty"`
	TimeoutSeconds  int      `yaml:"timeout_seconds,omitempty"`
	JobIDPath       string   `yaml:"job_id_path,omitempty"`
	PollURL         string   `yaml:"poll_url,omitempty"`
	StatusPath      string   `yaml:"status_path,omitempty"`
	DoneValues      []string `yaml:"done_values,omitempty"`
	FailValues      []string `yaml:"fail_values,omitempty"`
}

// Capabilities returns capability names configured on this entry.
func (e APIEntry) Capabilities() []string {
	var out []string
	if e.Text != nil {
		out = append(out, "text")
	}
	if e.Image != nil {
		out = append(out, "image")
	}
	if e.Video != nil {
		out = append(out, "video")
	}
	if e.Audio != nil {
		out = append(out, "audio")
	}
	if e.Mesh != nil {
		out = append(out, "mesh")
	}
	if e.Voice != nil {
		out = append(out, "voice")
	}
	if e.Video2Audio != nil {
		out = append(out, "video2audio")
	}
	if e.Train != nil {
		out = append(out, "train")
	}
	if e.Dataset != nil {
		out = append(out, "dataset")
	}
	if e.Data != nil {
		out = append(out, "data")
	}
	if e.RAG != nil {
		out = append(out, "rag")
	}
	return out
}

// APIEntryFromProvider migrates a legacy cloud provider config to an API entry.
func APIEntryFromProvider(p ProviderConfig) APIEntry {
	switch p.Type {
	case "openai":
		return APIEntry{
			Text: &APICapabilitySpec{
				Protocol:    "openai",
				BaseURL:     p.BaseURL,
				APIKey:      p.APIKey,
				APIKeyEnv:   p.APIKeyEnv,
				Model:       p.Model,
				MaxTokens:   p.MaxTokens,
				Temperature: p.Temperature,
			},
		}
	case "anthropic":
		return APIEntry{
			Text: &APICapabilitySpec{
				Protocol:         "anthropic",
				BaseURL:          p.BaseURL,
				APIKey:           p.APIKey,
				APIKeyEnv:        p.APIKeyEnv,
				Model:            p.Model,
				AnthropicVersion: p.AnthropicVersion,
				MaxTokens:        p.MaxTokens,
				Temperature:      p.Temperature,
			},
		}
	default:
		return APIEntry{}
	}
}
