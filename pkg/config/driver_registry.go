package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// DriverSection holds install-time and runtime settings for an external driver.
// Any driver ID from drivers-reqpack.yaml or config.yaml maps here — no per-backend
// struct fields are required in wuji-core.
type DriverSection struct {
	DriverBin             string   `yaml:"driver_bin,omitempty"`
	AutoStartDriver       *bool    `yaml:"auto_start_driver,omitempty"`
	GRPCAddr              string   `yaml:"grpc_addr,omitempty"`
	ServerBin             string   `yaml:"server_bin,omitempty"`
	VllmBin               string   `yaml:"vllm_bin,omitempty"`
	DefaultModel          string   `yaml:"default_model,omitempty"`
	Host                  string   `yaml:"host,omitempty"`
	Port                  int      `yaml:"port,omitempty"`
	APIBase               string   `yaml:"api_base,omitempty"`
	APIKey                string   `yaml:"api_key,omitempty"`
	ManageServer          *bool    `yaml:"manage_server,omitempty"`
	TimeoutSeconds        int      `yaml:"timeout_seconds,omitempty"`
	StartupTimeoutSeconds int      `yaml:"startup_timeout_seconds,omitempty"`
	ExtraArgs             []string `yaml:"extra_args,omitempty"`
}

var reservedDriverConfigKeys = map[string]bool{
	"resources": true,
	"vllm":      true,
	"llama":     true,
	"a1111":     true,
	"local_rag": true,
	"rag":       true,
	"raggo":     true,
	"gorag":     true,
}

// Section returns registry settings for a driver ID, or nil when not configured.
func (d *DriversConfig) Section(driverID string) *DriverSection {
	if d == nil || driverID == "" {
		return nil
	}
	if d.Registry == nil {
		return nil
	}
	section, ok := d.Registry[driverID]
	if !ok {
		return nil
	}
	return &section
}

func (d *DriversConfig) ensureRegistry() {
	if d.Registry == nil {
		d.Registry = map[string]DriverSection{}
	}
}

func mergeDriverSection(dst *DriverSection, src DriverSection) {
	if src.DriverBin != "" {
		dst.DriverBin = src.DriverBin
	}
	if src.AutoStartDriver != nil {
		dst.AutoStartDriver = src.AutoStartDriver
	}
	if src.GRPCAddr != "" {
		dst.GRPCAddr = src.GRPCAddr
	}
	if src.ServerBin != "" {
		dst.ServerBin = src.ServerBin
	}
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
	if src.ManageServer != nil {
		dst.ManageServer = src.ManageServer
	}
	if src.TimeoutSeconds != 0 {
		dst.TimeoutSeconds = src.TimeoutSeconds
	}
	if src.StartupTimeoutSeconds != 0 {
		dst.StartupTimeoutSeconds = src.StartupTimeoutSeconds
	}
	if len(src.ExtraArgs) > 0 {
		dst.ExtraArgs = append([]string(nil), src.ExtraArgs...)
	}
}

func mergeDriverRegistry(dst, src map[string]DriverSection) map[string]DriverSection {
	if len(src) == 0 {
		return dst
	}
	if dst == nil {
		dst = map[string]DriverSection{}
	}
	for id, section := range src {
		existing := dst[id]
		mergeDriverSection(&existing, section)
		dst[id] = existing
	}
	return dst
}

func (d *DriversConfig) UnmarshalYAML(node *yaml.Node) error {
	*d = DriversConfig{Registry: map[string]DriverSection{}}
	if node == nil || node.Kind == 0 {
		return nil
	}
	var raw map[string]yaml.Node
	if err := node.Decode(&raw); err != nil {
		return err
	}
	for key, val := range raw {
		if reservedDriverConfigKeys[key] {
			switch key {
			case "resources":
				var section ResourcesConfig
				if err := val.Decode(&section); err != nil {
					return fmt.Errorf("drivers.resources: %w", err)
				}
				d.Resources = &section
			case "vllm":
				var section VLLMConfig
				if err := val.Decode(&section); err != nil {
					return fmt.Errorf("drivers.vllm: %w", err)
				}
				d.VLLM = &section
			case "llama":
				var section LlamaConfig
				if err := val.Decode(&section); err != nil {
					return fmt.Errorf("drivers.llama: %w", err)
				}
				d.Llama = &section
			case "a1111":
				var section A1111Config
				if err := val.Decode(&section); err != nil {
					return fmt.Errorf("drivers.a1111: %w", err)
				}
				d.A1111 = &section
			case "local_rag":
				var section LocalRAGConfig
				if err := val.Decode(&section); err != nil {
					return fmt.Errorf("drivers.local_rag: %w", err)
				}
				d.LocalRAG = &section
			case "rag":
				var section LocalRAGConfig
				if err := val.Decode(&section); err != nil {
					return fmt.Errorf("drivers.rag: %w", err)
				}
				d.RAG = &section
			case "raggo":
				var section RaggoConfig
				if err := val.Decode(&section); err != nil {
					return fmt.Errorf("drivers.raggo: %w", err)
				}
				d.Raggo = &section
			case "gorag":
				var section GoragConfig
				if err := val.Decode(&section); err != nil {
					return fmt.Errorf("drivers.gorag: %w", err)
				}
				d.Gorag = &section
			}
			continue
		}
		var section DriverSection
		if err := val.Decode(&section); err != nil {
			return fmt.Errorf("drivers.%s: %w", key, err)
		}
		d.Registry[key] = section
	}
	return nil
}

func (d DriversConfig) MarshalYAML() (interface{}, error) {
	out := map[string]interface{}{}
	if d.Resources != nil {
		out["resources"] = d.Resources
	}
	if d.VLLM != nil {
		out["vllm"] = d.VLLM
	}
	if d.Llama != nil {
		out["llama"] = d.Llama
	}
	if d.A1111 != nil {
		out["a1111"] = d.A1111
	}
	if d.LocalRAG != nil {
		out["local_rag"] = d.LocalRAG
	}
	if d.RAG != nil {
		out["rag"] = d.RAG
	}
	if d.Raggo != nil {
		out["raggo"] = d.Raggo
	}
	if d.Gorag != nil {
		out["gorag"] = d.Gorag
	}
	for id, section := range d.Registry {
		out[id] = section
	}
	return out, nil
}
