package config

import "fmt"

const (
	defaultLlamaServerBinRel      = "../driver/llama/vendor/llama/llama-server"
	defaultLlamaModelsDirRel      = "../driver/llama/models"
	defaultLlamaHost              = "127.0.0.1"
	defaultLlamaPort              = 8080
	defaultOllamaAPI              = "http://127.0.0.1:11434"
	defaultLlamaStartupTimeoutSec = 600
)

// LlamaConfig holds settings for the llama.cpp driver.
type LlamaConfig struct {
	ServerBin             string   `yaml:"server_bin,omitempty"`
	ServerRunAs           string   `yaml:"server_run_as,omitempty"`
	ModelsDir             string   `yaml:"models_dir,omitempty"`
	DefaultModel          string   `yaml:"default_model,omitempty"`
	InferenceHost         string   `yaml:"inference_host,omitempty"`
	InferencePort         int      `yaml:"inference_port,omitempty"`
	StartupTimeoutSeconds int      `yaml:"startup_timeout_seconds,omitempty"`
	GRPCAddr              string   `yaml:"grpc_addr,omitempty"`
	DriverBin             string   `yaml:"driver_bin,omitempty"`
	AutoStartDriver       *bool    `yaml:"auto_start_driver,omitempty"`
	OllamaAPI             string   `yaml:"ollama_api,omitempty"`
	OllamaThink           *bool    `yaml:"ollama_think,omitempty"`
	UseOllamaAPI          *bool    `yaml:"use_ollama_api,omitempty"`
	ExtraArgs             []string `yaml:"extra_args,omitempty"`
}

// DefaultLlamaConfig returns built-in defaults (no file required).
func DefaultLlamaConfig() LlamaConfig {
	return LlamaConfig{
		ServerBin:             defaultLlamaServerBinRel,
		ModelsDir:             defaultLlamaModelsDirRel,
		InferenceHost:         defaultLlamaHost,
		InferencePort:         defaultLlamaPort,
		StartupTimeoutSeconds: defaultLlamaStartupTimeoutSec,
		OllamaAPI:             defaultOllamaAPI,
	}
}

// ResolvedLlama returns llama settings merged from defaults and .wuji/config.yaml.
func (c *Config) ResolvedLlama() LlamaConfig {
	out := DefaultLlamaConfig()
	if c == nil || c.Drivers.Llama == nil {
		return out
	}
	mergeLlama(&out, *c.Drivers.Llama)
	return out
}

func mergeLlama(dst *LlamaConfig, src LlamaConfig) {
	if src.ServerBin != "" {
		dst.ServerBin = src.ServerBin
	}
	if src.ServerRunAs != "" {
		dst.ServerRunAs = src.ServerRunAs
	}
	if src.ModelsDir != "" {
		dst.ModelsDir = src.ModelsDir
	}
	if src.DefaultModel != "" {
		dst.DefaultModel = src.DefaultModel
	}
	if src.InferenceHost != "" {
		dst.InferenceHost = src.InferenceHost
	}
	if src.InferencePort != 0 {
		dst.InferencePort = src.InferencePort
	}
	if src.StartupTimeoutSeconds != 0 {
		dst.StartupTimeoutSeconds = src.StartupTimeoutSeconds
	}
	if src.GRPCAddr != "" {
		dst.GRPCAddr = src.GRPCAddr
	}
	if src.OllamaAPI != "" {
		dst.OllamaAPI = src.OllamaAPI
	}
	if src.OllamaThink != nil {
		dst.OllamaThink = src.OllamaThink
	}
	if src.UseOllamaAPI != nil {
		dst.UseOllamaAPI = src.UseOllamaAPI
	}
	if src.DriverBin != "" {
		dst.DriverBin = src.DriverBin
	}
	if src.AutoStartDriver != nil {
		dst.AutoStartDriver = src.AutoStartDriver
	}
	if len(src.ExtraArgs) > 0 {
		dst.ExtraArgs = append([]string(nil), src.ExtraArgs...)
	}
}

func (c LlamaConfig) OllamaThinkEnabled() bool {
	if c.OllamaThink != nil {
		return *c.OllamaThink
	}
	// Native Ollama thinking when routing through the HTTP API (previous default behavior).
	return c.UseOllamaAPIEnabled()
}

// UseOllamaAPIEnabled reports whether unreadable Ollama blob symlinks may fall back to the Ollama HTTP API.
func (c LlamaConfig) UseOllamaAPIEnabled() bool {
	if c.UseOllamaAPI == nil {
		return false
	}
	return *c.UseOllamaAPI
}

// InferenceBaseURL returns the llama-server HTTP base URL.
func (c LlamaConfig) InferenceBaseURL() string {
	return fmt.Sprintf("http://%s:%d", c.InferenceHost, c.InferencePort)
}
