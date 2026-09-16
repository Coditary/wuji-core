package config

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/loraweight"
)

const loraUseDefaultsSentinel = "@default@"

// LoRAConfig holds named LoRA aliases and default stacks for image/text commands.
type LoRAConfig struct {
	DefaultWeight float32                     `yaml:"default_weight,omitempty"`
	Default       []string                    `yaml:"default,omitempty"`
	Library       map[string]LoRALibraryEntry `yaml:"library,omitempty"`
}

// LoRALibraryEntry is one named LoRA in the config library.
type LoRALibraryEntry struct {
	Path   string
	Weight float32
}

func (e *LoRALibraryEntry) UnmarshalYAML(value *yaml.Node) error {
	if value == nil {
		return nil
	}
	switch value.Kind {
	case yaml.ScalarNode:
		return value.Decode(&e.Path)
	case yaml.MappingNode:
		var aux struct {
			Path   string  `yaml:"path"`
			Weight float32 `yaml:"weight"`
		}
		if err := value.Decode(&aux); err != nil {
			return err
		}
		e.Path = aux.Path
		e.Weight = aux.Weight
		return nil
	default:
		return fmt.Errorf("lora library entry must be a string or mapping")
	}
}

func (e LoRALibraryEntry) MarshalYAML() (interface{}, error) {
	if e.Weight <= 0 {
		return e.Path, nil
	}
	return map[string]interface{}{
		"path":   e.Path,
		"weight": e.Weight,
	}, nil
}

// ResolvedLoRAs returns LoRA settings with defaults applied.
func (c *Config) ResolvedLoRAs() LoRAConfig {
	out := LoRAConfig{DefaultWeight: 1}
	if c == nil {
		return out
	}
	if c.LoRAs.DefaultWeight > 0 {
		out.DefaultWeight = c.LoRAs.DefaultWeight
	}
	if len(c.LoRAs.Default) > 0 {
		out.Default = append([]string(nil), c.LoRAs.Default...)
	}
	if len(c.LoRAs.Library) > 0 {
		out.Library = make(map[string]LoRALibraryEntry, len(c.LoRAs.Library))
		for k, v := range c.LoRAs.Library {
			out.Library[k] = v
		}
	}
	return out
}

// ResolveLibraryName maps an alias to a LoRARef.
func (c LoRAConfig) ResolveLibraryName(name, weightRaw string, fallbackWeight float32) (driver.LoRARef, error) {
	entry, ok := c.Library[name]
	if !ok {
		return driver.LoRARef{}, fmt.Errorf("unknown lora alias %q (define it in .wuji/config.yaml under loras.library)", name)
	}
	if strings.TrimSpace(entry.Path) == "" {
		return driver.LoRARef{}, fmt.Errorf("lora alias %q has empty path in config", name)
	}
	weight, err := c.resolveWeight(weightRaw, entry.Weight, fallbackWeight)
	if err != nil {
		return driver.LoRARef{}, fmt.Errorf("lora %q: %w", name, err)
	}
	return driver.LoRARef{Path: entry.Path, Weight: weight}, nil
}

// ResolveDefault returns configured default LoRA stack.
func (c LoRAConfig) ResolveDefault(fallbackWeight float32) ([]driver.LoRARef, error) {
	if len(c.Default) == 0 {
		return nil, fmt.Errorf("no loras.default configured in .wuji/config.yaml")
	}
	out := make([]driver.LoRARef, 0, len(c.Default))
	for _, name := range c.Default {
		ref, err := c.ResolveLibraryName(strings.TrimSpace(name), "", fallbackWeight)
		if err != nil {
			return nil, err
		}
		out = append(out, ref)
	}
	return out, nil
}

// ResolveEntry resolves a --lora value (config alias or filesystem path).
func (c LoRAConfig) ResolveEntry(entry string, fallbackWeight float32) (driver.LoRARef, error) {
	entry = strings.TrimSpace(entry)
	if entry == "" {
		return driver.LoRARef{}, fmt.Errorf("empty --lora value")
	}
	if entry == loraUseDefaultsSentinel {
		return driver.LoRARef{}, fmt.Errorf("internal lora sentinel")
	}

	if lib, ok := c.Library[entry]; ok {
		if strings.TrimSpace(lib.Path) == "" {
			return driver.LoRARef{}, fmt.Errorf("lora alias %q has empty path in config", entry)
		}
		weight, err := c.resolveWeight("", lib.Weight, fallbackWeight)
		if err != nil {
			return driver.LoRARef{}, fmt.Errorf("lora %q: %w", entry, err)
		}
		return driver.LoRARef{Path: lib.Path, Weight: weight}, nil
	}

	weight, err := c.resolveWeight("", 0, fallbackWeight)
	if err != nil {
		return driver.LoRARef{}, err
	}
	return driver.LoRARef{Path: entry, Weight: weight}, nil
}

func (c LoRAConfig) resolveWeight(explicit string, libraryDefault, fallback float32) (float32, error) {
	switch {
	case strings.TrimSpace(explicit) != "":
		return loraweight.Parse(explicit)
	case libraryDefault > 0:
		return libraryDefault, nil
	case fallback > 0:
		return fallback, nil
	default:
		return 1, nil
	}
}

// LoRAUseDefaultsSentinel is the cobra NoOptDefVal for bare --lora.
func LoRAUseDefaultsSentinel() string {
	return loraUseDefaultsSentinel
}

// SetLoRALibrary sets or updates one library alias.
func SetLoRALibrary(root, alias, path string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	if cfg.Root == "" {
		cfg.Root = root
	}
	if cfg.LoRAs.Library == nil {
		cfg.LoRAs.Library = make(map[string]LoRALibraryEntry)
	}
	cfg.LoRAs.Library[alias] = LoRALibraryEntry{Path: path}
	return cfg.Save()
}

// SetLoRADefault sets the default LoRA alias list.
func SetLoRADefault(root, value string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	if cfg.Root == "" {
		cfg.Root = root
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return fmt.Errorf("lora-default must list at least one alias")
	}
	cfg.LoRAs.Default = out
	return cfg.Save()
}

// SetLoRADefaultWeight sets the fallback weight for LoRAs without explicit strength.
func SetLoRADefaultWeight(root, value string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	if cfg.Root == "" {
		cfg.Root = root
	}
	weight, err := loraweight.Parse(value)
	if err != nil {
		return err
	}
	cfg.LoRAs.DefaultWeight = weight
	return cfg.Save()
}
