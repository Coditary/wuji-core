package config

import "fmt"

const (
	defaultVLLMBin               = "vllm"
	defaultVLLMHost              = "127.0.0.1"
	defaultVLLMPort              = 8000
	defaultVLLMTimeoutSeconds    = 300
	defaultVLLMStartupTimeoutSec = 300
)

// VLLMConfig holds settings for the vLLM driver.
type VLLMConfig struct {
	VllmBin               string   `yaml:"vllm_bin,omitempty"`
	DefaultModel          string   `yaml:"default_model,omitempty"`
	Host                  string   `yaml:"host,omitempty"`
	Port                  int      `yaml:"port,omitempty"`
	APIBase               string   `yaml:"api_base,omitempty"`
	APIKey                string   `yaml:"api_key,omitempty"`
	GRPCAddr              string   `yaml:"grpc_addr,omitempty"`
	DriverBin             string   `yaml:"driver_bin,omitempty"`
	AutoStartDriver       *bool    `yaml:"auto_start_driver,omitempty"`
	TimeoutSeconds        int      `yaml:"timeout_seconds,omitempty"`
	StartupTimeoutSeconds int      `yaml:"startup_timeout_seconds,omitempty"`
	ManageServer          *bool    `yaml:"manage_server,omitempty"`
	ExtraArgs             []string `yaml:"extra_args,omitempty"`
}

// DefaultVLLMConfig returns built-in defaults (no file required).
func DefaultVLLMConfig() VLLMConfig {
	manage := true
	return VLLMConfig{
		VllmBin:               defaultVLLMBin,
		Host:                  defaultVLLMHost,
		Port:                  defaultVLLMPort,
		TimeoutSeconds:        defaultVLLMTimeoutSeconds,
		StartupTimeoutSeconds: defaultVLLMStartupTimeoutSec,
		ManageServer:          &manage,
	}
}

// ResolvedVLLM returns vLLM settings merged from defaults and .wuji/config.yaml.
func (c *Config) ResolvedVLLM() VLLMConfig {
	out := DefaultVLLMConfig()
	if c == nil || c.Drivers.VLLM == nil {
		return out
	}
	mergeVLLM(&out, *c.Drivers.VLLM)
	return out
}

func mergeVLLM(dst *VLLMConfig, src VLLMConfig) {
	if src.VllmBin != "" {
		dst.VllmBin = src.VllmBin
	}
	if src.DefaultModel != "" {
		dst.DefaultModel = src.DefaultModel
	}
	if src.Host != "" {
		dst.Host = src.Host
	}
	if src.Port != 0 {
		dst.Port = src.Port
	}
	if src.APIBase != "" {
		dst.APIBase = src.APIBase
	}
	if src.APIKey != "" {
		dst.APIKey = src.APIKey
	}
	if src.GRPCAddr != "" {
		dst.GRPCAddr = src.GRPCAddr
	}
	if src.TimeoutSeconds != 0 {
		dst.TimeoutSeconds = src.TimeoutSeconds
	}
	if src.StartupTimeoutSeconds != 0 {
		dst.StartupTimeoutSeconds = src.StartupTimeoutSeconds
	}
	if src.ManageServer != nil {
		dst.ManageServer = src.ManageServer
	}
	if len(src.ExtraArgs) > 0 {
		dst.ExtraArgs = append([]string(nil), src.ExtraArgs...)
	}
	if src.DriverBin != "" {
		dst.DriverBin = src.DriverBin
	}
	if src.AutoStartDriver != nil {
		dst.AutoStartDriver = src.AutoStartDriver
	}
}

// ResolvedAPIBase returns the OpenAI API base URL for vLLM.
func (c VLLMConfig) ResolvedAPIBase() string {
	if c.APIBase != "" {
		if len(c.APIBase) > 0 && c.APIBase[len(c.APIBase)-1] == '/' {
			return c.APIBase[:len(c.APIBase)-1]
		}
		return c.APIBase
	}
	return fmt.Sprintf("http://%s:%d/v1", c.Host, c.Port)
}

func (c VLLMConfig) ManageServerEnabled() bool {
	if c.ManageServer == nil {
		return true
	}
	return *c.ManageServer
}
