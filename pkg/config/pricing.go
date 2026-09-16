package config

// ModelPricing holds per-million-token pricing in USD.
type ModelPricing struct {
	InputPerMillion  float64 `yaml:"input_per_million"`
	OutputPerMillion float64 `yaml:"output_per_million"`
}

// PricingFor returns pricing for a model, falling back to "default" then zero cost.
func (c *Config) PricingFor(model string) ModelPricing {
	if c == nil || len(c.Pricing) == 0 {
		return ModelPricing{}
	}
	if p, ok := c.Pricing[model]; ok {
		return p
	}
	return c.Pricing["default"]
}
