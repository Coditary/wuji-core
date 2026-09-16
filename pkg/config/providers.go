package config

// ProviderConfig defines an LLM provider entry in .wuji/config.yaml.
// Taiji and other clients can sync this catalog via ListProviders.
type ProviderConfig struct {
	// driver | openai | anthropic
	Type string `yaml:"type"`

	// driver type: Wuji driver id (vllm, llama, echo, anthropic, …)
	Driver string `yaml:"driver,omitempty"`
	// driver type: GenerateText or ChatComplete (default ChatComplete)
	Method string `yaml:"method,omitempty"`

	Model string `yaml:"model,omitempty"`

	BaseURL   string `yaml:"base_url,omitempty"`
	APIKey    string `yaml:"api_key,omitempty"`
	APIKeyEnv string `yaml:"api_key_env,omitempty"`

	AnthropicVersion string  `yaml:"anthropic_version,omitempty"`
	MaxTokens        int     `yaml:"max_tokens,omitempty"`
	Temperature      float64 `yaml:"temperature,omitempty"`
}

// ProviderSummary is exposed to clients (Taiji, etc.) — no secrets or routing details.
type ProviderSummary struct {
	Name   string   `json:"name,omitempty"`
	Model  string   `json:"model,omitempty"`
	Driver string   `json:"driver,omitempty"` // Wuji driver id for GenerateText / ChatComplete
	Kind   string   `json:"kind,omitempty"`   // api | driver | provider | catalog
	Models []string `json:"models,omitempty"`
	API    string   `json:"api,omitempty"`
	NPM    string   `json:"npm,omitempty"`
	Env    []string `json:"env,omitempty"`
}

// ProviderCatalog is returned by the core ListProviders API.
type ProviderCatalog struct {
	DefaultProvider string                     `json:"default_provider"`
	Providers       map[string]ProviderSummary `json:"providers"`
}
