package config

import "os"

const (
	defaultRaggoVectorDB       = "file"
	defaultRaggoEmbedProvider  = "ollama"
	defaultRaggoHybridSearch   = true
	defaultRaggoOpenAIProvider = "openai"
)

// RaggoConfig holds settings for the raggo RAG driver.
type RaggoConfig struct {
	DefaultEmbedModel string `yaml:"default_embed_model,omitempty"`
	OllamaAPI         string `yaml:"ollama_api,omitempty"`
	OpenAIAPIKey      string `yaml:"openai_api_key,omitempty"`
	EmbedProvider     string `yaml:"embed_provider,omitempty"` // ollama | openai
	VectorDB          string `yaml:"vector_db,omitempty"`        // file | memory | chromem
	HybridSearch      *bool  `yaml:"hybrid_search,omitempty"`
	GRPCAddr          string `yaml:"grpc_addr,omitempty"`
	DriverBin         string `yaml:"driver_bin,omitempty"`
	AutoStartDriver   *bool  `yaml:"auto_start_driver,omitempty"`
	TimeoutSeconds    int    `yaml:"timeout_seconds,omitempty"`
}

// DefaultRaggoConfig returns built-in defaults.
func DefaultRaggoConfig() RaggoConfig {
	hybrid := defaultRaggoHybridSearch
	return RaggoConfig{
		DefaultEmbedModel: defaultLocalRAGEmbedModel,
		OllamaAPI:         defaultLocalRAGOllamaAPI,
		EmbedProvider:     defaultRaggoEmbedProvider,
		VectorDB:          defaultRaggoVectorDB,
		HybridSearch:      &hybrid,
	}
}

// ResolvedRaggo merges defaults with .wuji/config.yaml.
func (c *Config) ResolvedRaggo() RaggoConfig {
	out := DefaultRaggoConfig()
	if c == nil {
		return out
	}
	if c.Drivers.Raggo != nil {
		mergeRaggo(&out, *c.Drivers.Raggo)
	}
	local := c.ResolvedLocalRAG()
	if out.DefaultEmbedModel == defaultLocalRAGEmbedModel && local.DefaultEmbedModel != "" {
		out.DefaultEmbedModel = local.DefaultEmbedModel
	}
	if out.OllamaAPI == defaultLocalRAGOllamaAPI && local.OllamaAPI != "" {
		out.OllamaAPI = local.OllamaAPI
	}
	if out.OpenAIAPIKey == "" {
		out.OpenAIAPIKey = os.Getenv("OPENAI_API_KEY")
	}
	return out
}

func mergeRaggo(dst *RaggoConfig, src RaggoConfig) {
	if src.DefaultEmbedModel != "" {
		dst.DefaultEmbedModel = src.DefaultEmbedModel
	}
	if src.OllamaAPI != "" {
		dst.OllamaAPI = src.OllamaAPI
	}
	if src.OpenAIAPIKey != "" {
		dst.OpenAIAPIKey = src.OpenAIAPIKey
	}
	if src.EmbedProvider != "" {
		dst.EmbedProvider = src.EmbedProvider
	}
	if src.VectorDB != "" {
		dst.VectorDB = src.VectorDB
	}
	if src.HybridSearch != nil {
		dst.HybridSearch = src.HybridSearch
	}
	if src.GRPCAddr != "" {
		dst.GRPCAddr = src.GRPCAddr
	}
	if src.DriverBin != "" {
		dst.DriverBin = src.DriverBin
	}
	if src.AutoStartDriver != nil {
		dst.AutoStartDriver = src.AutoStartDriver
	}
	if src.TimeoutSeconds != 0 {
		dst.TimeoutSeconds = src.TimeoutSeconds
	}
}
