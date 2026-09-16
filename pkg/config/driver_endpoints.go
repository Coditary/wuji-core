package config

import "strings"

// AddDriverEndpoint appends an endpoint and persists config.yaml.
func (c *Config) AddDriverEndpoint(endpoint string) error {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return nil
	}
	for _, ep := range c.DriverEndpoints {
		if ep == endpoint {
			return nil
		}
	}
	c.DriverEndpoints = append(c.DriverEndpoints, endpoint)
	return c.Save()
}

// ResolvedDriverEndpoints returns endpoints to connect at startup:
// explicit driver_endpoints plus grpc_addr for configured remote drivers.
func (c *Config) ResolvedDriverEndpoints() []string {
	if c == nil {
		return nil
	}
	seen := make(map[string]bool)
	var out []string
	add := func(ep string) {
		ep = strings.TrimSpace(ep)
		if ep == "" || seen[ep] {
			return
		}
		seen[ep] = true
		out = append(out, ep)
	}

	for _, ep := range c.DriverEndpoints {
		add(ep)
	}
	for _, id := range c.configuredDriverIDs() {
		if id == "dummy" || id == "local-rag" || id == "local" || id == "rag" || id == "ffmpeg" {
			continue // built-in drivers
		}
		add(c.DriverGRPCAddr(id))
	}
	return out
}

func (c *Config) configuredDriverIDs() []string {
	if c == nil {
		return nil
	}
	seen := make(map[string]bool)
	var ids []string
	add := func(id string) {
		id = strings.TrimSpace(id)
		if id == "" || seen[id] {
			return
		}
		seen[id] = true
		ids = append(ids, id)
	}
	add(c.DefaultDriver)
	for _, id := range c.CapabilityDrivers {
		add(id)
	}
	if c.Drivers.A1111 != nil {
		add("a1111")
	}
	if c.Drivers.VLLM != nil {
		add("vllm")
	}
	if c.Drivers.Llama != nil {
		add("llama")
	}
	if c.Drivers.LocalRAG != nil || c.Drivers.RAG != nil {
		add("local-rag")
	}
	if c.Drivers.Raggo != nil {
		add("raggo")
	}
	if c.Drivers.Gorag != nil {
		add("gorag")
	}
	return ids
}

func dedupeEndpoints(endpoints []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, ep := range endpoints {
		ep = strings.TrimSpace(ep)
		if ep == "" || seen[ep] {
			continue
		}
		seen[ep] = true
		out = append(out, ep)
	}
	return out
}
