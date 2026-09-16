package config

// migrateProvidersToAPIs copies legacy providers.* cloud entries into apis.* when absent.
func (c *Config) migrateProvidersToAPIs() {
	if c.Providers == nil {
		return
	}
	if c.APIs == nil {
		c.APIs = map[string]APIEntry{}
	}
	for id, p := range c.Providers {
		if _, exists := c.APIs[id]; exists {
			continue
		}
		entry := APIEntryFromProvider(p)
		if len(entry.Capabilities()) == 0 {
			continue
		}
		c.APIs[id] = entry
	}
}
