package config

import "os"

const (
	defaultGoragVectorDB   = "qdrant"
	defaultGoragQdrantHost = "localhost"
	defaultGoragQdrantPort = 6334
)

// GoragConfig holds settings for the gorag RAG driver.
type GoragConfig struct {
	DefaultEmbedModel string `yaml:"default_embed_model,omitempty"`
	OllamaAPI         string `yaml:"ollama_api,omitempty"`
	VectorDB          string `yaml:"vector_db,omitempty"` // qdrant | pgvector
	QdrantHost        string `yaml:"qdrant_host,omitempty"`
	QdrantPort        int    `yaml:"qdrant_port,omitempty"`
	PostgresURL       string `yaml:"postgres_url,omitempty"`
	GRPCAddr          string `yaml:"grpc_addr,omitempty"`
	DriverBin         string `yaml:"driver_bin,omitempty"`
	AutoStartDriver   *bool  `yaml:"auto_start_driver,omitempty"`
	TimeoutSeconds    int    `yaml:"timeout_seconds,omitempty"`
}

// DefaultGoragConfig returns built-in defaults.
func DefaultGoragConfig() GoragConfig {
	return GoragConfig{
		DefaultEmbedModel: defaultLocalRAGEmbedModel,
		OllamaAPI:         defaultLocalRAGOllamaAPI,
		VectorDB:          defaultGoragVectorDB,
		QdrantHost:        defaultGoragQdrantHost,
		QdrantPort:        defaultGoragQdrantPort,
	}
}

// ResolvedGorag merges defaults with .wuji/config.yaml.
func (c *Config) ResolvedGorag() GoragConfig {
	out := DefaultGoragConfig()
	if c == nil {
		return out
	}
	if c.Drivers.Gorag != nil {
		mergeGorag(&out, *c.Drivers.Gorag)
	}
	local := c.ResolvedLocalRAG()
	if out.DefaultEmbedModel == defaultLocalRAGEmbedModel && local.DefaultEmbedModel != "" {
		out.DefaultEmbedModel = local.DefaultEmbedModel
	}
	if out.OllamaAPI == defaultLocalRAGOllamaAPI && local.OllamaAPI != "" {
		out.OllamaAPI = local.OllamaAPI
	}
	if out.PostgresURL == "" {
		out.PostgresURL = os.Getenv("GORAG_POSTGRES_URL")
	}
	return out
}

func mergeGorag(dst *GoragConfig, src GoragConfig) {
	if src.DefaultEmbedModel != "" {
		dst.DefaultEmbedModel = src.DefaultEmbedModel
	}
	if src.OllamaAPI != "" {
		dst.OllamaAPI = src.OllamaAPI
	}
	if src.VectorDB != "" {
		dst.VectorDB = src.VectorDB
	}
	if src.QdrantHost != "" {
		dst.QdrantHost = src.QdrantHost
	}
	if src.QdrantPort != 0 {
		dst.QdrantPort = src.QdrantPort
	}
	if src.PostgresURL != "" {
		dst.PostgresURL = src.PostgresURL
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
