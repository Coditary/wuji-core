package config

const (
	defaultLocalRAGEmbedModel = "nomic-embed-text"
	defaultLocalRAGOllamaAPI  = defaultOllamaAPI
)

// LocalRAGConfig holds settings for the local-rag driver.
type LocalRAGConfig struct {
	DefaultEmbedModel string `yaml:"default_embed_model,omitempty"`
	OllamaAPI         string `yaml:"ollama_api,omitempty"`
	GRPCAddr          string `yaml:"grpc_addr,omitempty"`
	DriverBin         string `yaml:"driver_bin,omitempty"`
	AutoStartDriver   *bool  `yaml:"auto_start_driver,omitempty"`
	TimeoutSeconds    int    `yaml:"timeout_seconds,omitempty"`
}

// RAGConfig is a deprecated alias for LocalRAGConfig (drivers.rag in old configs).
type RAGConfig = LocalRAGConfig

// DefaultLocalRAGConfig returns built-in defaults (no file required).
func DefaultLocalRAGConfig() LocalRAGConfig {
	return LocalRAGConfig{
		DefaultEmbedModel: defaultLocalRAGEmbedModel,
		OllamaAPI:         defaultLocalRAGOllamaAPI,
	}
}

// ResolvedLocalRAG returns local-rag settings merged from defaults and .wuji/config.yaml.
func (c *Config) ResolvedLocalRAG() LocalRAGConfig {
	out := DefaultLocalRAGConfig()
	if c == nil {
		return out
	}
	if c.Drivers.LocalRAG != nil {
		mergeLocalRAG(&out, *c.Drivers.LocalRAG)
	} else if c.Drivers.RAG != nil {
		mergeLocalRAG(&out, *c.Drivers.RAG)
	}
	return out
}

// ResolvedRAG is an alias for ResolvedLocalRAG.
func (c *Config) ResolvedRAG() LocalRAGConfig {
	return c.ResolvedLocalRAG()
}

func mergeLocalRAG(dst *LocalRAGConfig, src LocalRAGConfig) {
	if src.DefaultEmbedModel != "" {
		dst.DefaultEmbedModel = src.DefaultEmbedModel
	}
	if src.OllamaAPI != "" {
		dst.OllamaAPI = src.OllamaAPI
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
