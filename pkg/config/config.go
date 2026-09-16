package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	dirName         = ".wuji"
	configFileName  = "config.yaml"
	driversFileName = "drivers.yaml"
)

// Config holds Wuji runtime configuration.
type Config struct {
	Root              string
	DefaultDriver     string                    `yaml:"default_driver,omitempty"`
	DefaultProvider   string                    `yaml:"default_provider,omitempty"`
	Providers         map[string]ProviderConfig `yaml:"providers,omitempty"`
	APIs              map[string]APIEntry       `yaml:"apis,omitempty"`
	CapabilityDrivers CapabilityDrivers         `yaml:"capability_drivers,omitempty"`
	LoRAs             LoRAConfig                `yaml:"loras,omitempty"`
	Drivers           DriversConfig             `yaml:"drivers,omitempty"`
	DriverEndpoints   []string
	Pricing           map[string]ModelPricing `yaml:"pricing,omitempty"`
	Catalog           *CatalogConfig          `yaml:"catalog,omitempty"`
}

// DriversConfig holds per-driver settings.
type DriversConfig struct {
	Resources *ResourcesConfig `yaml:"resources,omitempty"`
	VLLM      *VLLMConfig      `yaml:"vllm,omitempty"`
	Llama     *LlamaConfig     `yaml:"llama,omitempty"`
	A1111     *A1111Config     `yaml:"a1111,omitempty"`
	LocalRAG  *LocalRAGConfig  `yaml:"local_rag,omitempty"`
	RAG       *LocalRAGConfig  `yaml:"rag,omitempty"` // deprecated alias
	Raggo     *RaggoConfig     `yaml:"raggo,omitempty"`
	Gorag     *GoragConfig     `yaml:"gorag,omitempty"`
}

type driversEndpointsConfig struct {
	Drivers []driverEntry `yaml:"drivers"`
}

type driverEntry struct {
	Endpoint string `yaml:"endpoint"`
}

type fileConfig struct {
	DefaultDriver     string                    `yaml:"default_driver,omitempty"`
	DefaultProvider   string                    `yaml:"default_provider,omitempty"`
	Providers         map[string]ProviderConfig `yaml:"providers,omitempty"`
	APIs              map[string]APIEntry       `yaml:"apis,omitempty"`
	CapabilityDrivers CapabilityDrivers         `yaml:"capability_drivers,omitempty"`
	LoRAs             LoRAConfig                `yaml:"loras,omitempty"`
	Drivers           DriversConfig             `yaml:"drivers,omitempty"`
	DriverEndpoints   []string                  `yaml:"driver_endpoints,omitempty"`
	Pricing           map[string]ModelPricing   `yaml:"pricing,omitempty"`
	Catalog           *CatalogConfig            `yaml:"catalog,omitempty"`
}

// Load reads configuration from the project root.
func Load() (*Config, error) {
	root, err := findProjectRoot()
	if err != nil {
		return nil, err
	}

	cfg := &Config{Root: root}
	if v := os.Getenv("WUJI_ROOT"); v != "" {
		cfg.Root = v
	}

	configPath := filepath.Join(cfg.Root, dirName, configFileName)
	if data, err := os.ReadFile(configPath); err == nil {
		var fc fileConfig
		if err := yaml.Unmarshal(data, &fc); err != nil {
			return nil, fmt.Errorf("parse %s: %w", configPath, err)
		}
		cfg.DefaultDriver = fc.DefaultDriver
		cfg.DefaultProvider = fc.DefaultProvider
		cfg.Providers = fc.Providers
		cfg.APIs = fc.APIs
		cfg.CapabilityDrivers = fc.CapabilityDrivers
		cfg.LoRAs = fc.LoRAs
		cfg.Drivers = fc.Drivers
		cfg.DriverEndpoints = fc.DriverEndpoints
		cfg.Pricing = fc.Pricing
		cfg.Catalog = fc.Catalog
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("read %s: %w", configPath, err)
	}

	driversPath := filepath.Join(cfg.Root, dirName, driversFileName)
	legacyMigrated := false
	if data, err := os.ReadFile(driversPath); err == nil {
		var df driversEndpointsConfig
		if err := yaml.Unmarshal(data, &df); err != nil {
			return nil, fmt.Errorf("parse %s: %w", driversPath, err)
		}
		for _, d := range df.Drivers {
			if d.Endpoint != "" {
				cfg.DriverEndpoints = append(cfg.DriverEndpoints, d.Endpoint)
				legacyMigrated = true
			}
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("read %s: %w", driversPath, err)
	}

	cfg.DriverEndpoints = dedupeEndpoints(cfg.DriverEndpoints)
	cfg.migrateProvidersToAPIs()
	if cfg.DefaultDriver == "" && cfg.DefaultProvider != "" {
		cfg.DefaultDriver = cfg.DefaultProvider
	}

	if legacyMigrated {
		if err := cfg.Save(); err != nil {
			return nil, fmt.Errorf("migrate %s into %s: %w", driversFileName, configFileName, err)
		}
		_ = os.Remove(driversPath)
	}

	return cfg, nil
}

// Save writes the main config file, preserving unknown fields is not required for MVP.
func (c *Config) Save() error {
	dir := filepath.Join(c.Root, dirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	fc := fileConfig{
		DefaultDriver:     c.DefaultDriver,
		DefaultProvider:   c.DefaultProvider,
		Providers:         c.Providers,
		APIs:              c.APIs,
		CapabilityDrivers: c.CapabilityDrivers,
		LoRAs:             c.LoRAs,
		Drivers:           c.Drivers,
		DriverEndpoints:   c.DriverEndpoints,
		Pricing:           c.Pricing,
		Catalog:           c.Catalog,
	}
	data, err := yaml.Marshal(fc)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, configFileName), data, 0o644)
}

// SetVLLMField updates a single vLLM config key and persists config.yaml.
func SetVLLMField(root, key, value string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	if cfg.Root == "" {
		cfg.Root = root
	}
	if cfg.Drivers.VLLM == nil {
		cfg.Drivers.VLLM = &VLLMConfig{}
	}

	switch key {
	case "vllm_bin":
		cfg.Drivers.VLLM.VllmBin = value
	case "default_model":
		cfg.Drivers.VLLM.DefaultModel = value
	case "host":
		cfg.Drivers.VLLM.Host = value
	case "port":
		var port int
		if _, err := fmt.Sscanf(value, "%d", &port); err != nil {
			return fmt.Errorf("port must be an integer")
		}
		cfg.Drivers.VLLM.Port = port
	case "api_base":
		cfg.Drivers.VLLM.APIBase = value
	case "api_key":
		cfg.Drivers.VLLM.APIKey = value
	case "grpc_addr":
		cfg.Drivers.VLLM.GRPCAddr = value
	case "timeout_seconds":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil {
			return fmt.Errorf("timeout_seconds must be an integer")
		}
		cfg.Drivers.VLLM.TimeoutSeconds = n
	case "startup_timeout_seconds":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil {
			return fmt.Errorf("startup_timeout_seconds must be an integer")
		}
		cfg.Drivers.VLLM.StartupTimeoutSeconds = n
	case "manage_server":
		switch strings.ToLower(value) {
		case "true", "1", "yes":
			v := true
			cfg.Drivers.VLLM.ManageServer = &v
		case "false", "0", "no":
			v := false
			cfg.Drivers.VLLM.ManageServer = &v
		default:
			return fmt.Errorf("manage_server must be true or false")
		}
	default:
		return fmt.Errorf("unknown vLLM config key %q", key)
	}

	return cfg.Save()
}

// SetLlamaField updates a single llama driver config key and persists config.yaml.
func SetLlamaField(root, key, value string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	if cfg.Root == "" {
		cfg.Root = root
	}
	if cfg.Drivers.Llama == nil {
		cfg.Drivers.Llama = &LlamaConfig{}
	}

	switch key {
	case "server_bin":
		cfg.Drivers.Llama.ServerBin = value
	case "server_run_as":
		cfg.Drivers.Llama.ServerRunAs = value
	case "models_dir":
		cfg.Drivers.Llama.ModelsDir = value
	case "default_model":
		cfg.Drivers.Llama.DefaultModel = value
	case "inference_host":
		cfg.Drivers.Llama.InferenceHost = value
	case "inference_port":
		var port int
		if _, err := fmt.Sscanf(value, "%d", &port); err != nil {
			return fmt.Errorf("inference_port must be an integer")
		}
		cfg.Drivers.Llama.InferencePort = port
	case "grpc_addr":
		cfg.Drivers.Llama.GRPCAddr = value
	case "ollama_api":
		cfg.Drivers.Llama.OllamaAPI = value
	case "ollama_think":
		switch strings.ToLower(value) {
		case "true", "1", "yes":
			v := true
			cfg.Drivers.Llama.OllamaThink = &v
		case "false", "0", "no":
			v := false
			cfg.Drivers.Llama.OllamaThink = &v
		default:
			return fmt.Errorf("ollama_think must be true or false")
		}
	case "use_ollama_api":
		switch strings.ToLower(value) {
		case "true", "1", "yes":
			v := true
			cfg.Drivers.Llama.UseOllamaAPI = &v
		case "false", "0", "no":
			v := false
			cfg.Drivers.Llama.UseOllamaAPI = &v
		default:
			return fmt.Errorf("use_ollama_api must be true or false")
		}
	default:
		return fmt.Errorf("unknown llama config key %q", key)
	}

	return cfg.Save()
}

// SetA1111Field updates a single A1111 driver config key and persists config.yaml.
func SetA1111Field(root, key, value string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	if cfg.Root == "" {
		cfg.Root = root
	}
	if cfg.Drivers.A1111 == nil {
		cfg.Drivers.A1111 = &A1111Config{}
	}

	switch key {
	case "host":
		cfg.Drivers.A1111.Host = value
	case "port":
		var port int
		if _, err := fmt.Sscanf(value, "%d", &port); err != nil {
			return fmt.Errorf("port must be an integer")
		}
		cfg.Drivers.A1111.Port = port
	case "api_base":
		cfg.Drivers.A1111.APIBase = value
	case "api_key":
		cfg.Drivers.A1111.APIKey = value
	case "auth_user":
		cfg.Drivers.A1111.AuthUser = value
	case "auth_pass":
		cfg.Drivers.A1111.AuthPass = value
	case "default_model":
		cfg.Drivers.A1111.DefaultModel = value
	case "default_inpaint_model":
		cfg.Drivers.A1111.DefaultInpaintModel = value
	case "default_video_model":
		cfg.Drivers.A1111.DefaultVideoModel = value
	case "grpc_addr":
		cfg.Drivers.A1111.GRPCAddr = value
	case "timeout_seconds":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil {
			return fmt.Errorf("timeout_seconds must be an integer")
		}
		cfg.Drivers.A1111.TimeoutSeconds = n
	case "output_dir":
		cfg.Drivers.A1111.OutputDir = value
	case "video_output_dir":
		cfg.Drivers.A1111.VideoOutputDir = value
	case "video_width":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil {
			return fmt.Errorf("video_width must be an integer")
		}
		cfg.Drivers.A1111.VideoWidth = n
	case "video_height":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil {
			return fmt.Errorf("video_height must be an integer")
		}
		cfg.Drivers.A1111.VideoHeight = n
	case "video_steps":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil {
			return fmt.Errorf("video_steps must be an integer")
		}
		cfg.Drivers.A1111.VideoSteps = n
	case "video_cfg_scale":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil {
			return fmt.Errorf("video_cfg_scale must be an integer")
		}
		cfg.Drivers.A1111.VideoCFGScale = n
	case "video_sampler":
		cfg.Drivers.A1111.VideoSampler = value
	case "video_fps":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil {
			return fmt.Errorf("video_fps must be an integer")
		}
		cfg.Drivers.A1111.VideoFPS = n
	case "webui_dir":
		cfg.Drivers.A1111.WebUIDir = value
	case "auto_start_driver":
		switch strings.ToLower(value) {
		case "true", "1", "yes":
			v := true
			cfg.Drivers.A1111.AutoStartDriver = &v
		case "false", "0", "no":
			v := false
			cfg.Drivers.A1111.AutoStartDriver = &v
		default:
			return fmt.Errorf("auto_start_driver must be true or false")
		}
	case "manage_server":
		switch strings.ToLower(value) {
		case "true", "1", "yes":
			v := true
			cfg.Drivers.A1111.ManageServer = &v
		case "false", "0", "no":
			v := false
			cfg.Drivers.A1111.ManageServer = &v
		default:
			return fmt.Errorf("manage_server must be true or false")
		}
	default:
		return fmt.Errorf("unknown A1111 config key %q", key)
	}

	return cfg.Save()
}

// SaveDriverEndpoint appends an endpoint to config.yaml (driver_endpoints).
func SaveDriverEndpoint(root, endpoint string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	if cfg.Root == "" {
		cfg.Root = root
	}
	return cfg.AddDriverEndpoint(endpoint)
}

func findProjectRoot() (string, error) {
	if v := strings.TrimSpace(os.Getenv("WUJI_ROOT")); v != "" {
		return filepath.Clean(v), nil
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			if sibling := siblingCoreRoot(dir); sibling != "" {
				return sibling, nil
			}
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("could not find project root (no go.mod); set WUJI_ROOT to your wuji-core checkout (e.g. ~/Dev/Coditary/core/wuji-core)")
			}
			return home, nil
		}
		dir = parent
	}
}

func siblingCoreRoot(dir string) string {
	mod, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return ""
	}
	if !strings.Contains(string(mod), "module github.com/coditary/wuji-ai") {
		return ""
	}
	parent := filepath.Dir(dir)
	candidates := []string{
		filepath.Join(parent, "wuji-core"),
		filepath.Join(parent, "..", "core", "wuji-core"),
		filepath.Join(parent, "..", "..", "Coditary", "core", "wuji-core"),
	}
	for _, sibling := range candidates {
		if _, err := os.Stat(filepath.Join(sibling, ".wuji")); err == nil {
			return sibling
		}
	}
	return ""
}
