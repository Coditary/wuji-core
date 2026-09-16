package config

// CatalogConfig controls models.dev catalog sync and cache behavior.
type CatalogConfig struct {
	// Source is the models.dev base URL (default https://models.dev).
	Source string `yaml:"source,omitempty"`
	// AutoSyncHours refreshes the catalog when older than this many hours (0 = manual only).
	AutoSyncHours int `yaml:"auto_sync_hours,omitempty"`
}

// ResolvedCatalog returns catalog settings with defaults applied.
func (c *Config) ResolvedCatalog() CatalogConfig {
	out := CatalogConfig{
		Source:        "https://models.dev",
		AutoSyncHours: 0,
	}
	if c == nil || c.Catalog == nil {
		return out
	}
	if s := c.Catalog.Source; s != "" {
		out.Source = s
	}
	if c.Catalog.AutoSyncHours > 0 {
		out.AutoSyncHours = c.Catalog.AutoSyncHours
	}
	return out
}
