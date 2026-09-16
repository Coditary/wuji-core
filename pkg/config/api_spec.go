package config

import "strings"

// ResolveAPICapabilitySpec returns the effective spec for a task, merging task overrides onto base.
func ResolveAPICapabilitySpec(base *APICapabilitySpec, task string) *APICapabilitySpec {
	if base == nil {
		return nil
	}
	task = strings.TrimSpace(task)
	if task == "" || base.Tasks == nil {
		return base
	}
	child, ok := base.Tasks[task]
	if !ok || child == nil {
		return base
	}
	return MergeAPICapabilitySpec(base, child)
}

// MergeAPICapabilitySpec overlays override onto base (child wins for set fields).
func MergeAPICapabilitySpec(base, override *APICapabilitySpec) *APICapabilitySpec {
	if base == nil {
		return override
	}
	if override == nil {
		return base
	}
	out := *base
	if s := strings.TrimSpace(override.Protocol); s != "" {
		out.Protocol = s
	}
	if s := strings.TrimSpace(override.Method); s != "" {
		out.Method = s
	}
	if s := strings.TrimSpace(override.URL); s != "" {
		out.URL = s
	}
	if s := strings.TrimSpace(override.BaseURL); s != "" {
		out.BaseURL = s
	}
	if s := strings.TrimSpace(override.BodyRaw); s != "" {
		out.BodyRaw = s
	}
	if s := strings.TrimSpace(override.ContentType); s != "" {
		out.ContentType = s
	}
	if s := strings.TrimSpace(override.ErrorPath); s != "" {
		out.ErrorPath = s
	}
	if s := strings.TrimSpace(override.APIKey); s != "" {
		out.APIKey = s
	}
	if s := strings.TrimSpace(override.APIKeyEnv); s != "" {
		out.APIKeyEnv = s
	}
	if s := strings.TrimSpace(override.Model); s != "" {
		out.Model = s
	}
	if s := strings.TrimSpace(override.AnthropicVersion); s != "" {
		out.AnthropicVersion = s
	}
	if override.MaxTokens > 0 {
		out.MaxTokens = override.MaxTokens
	}
	if override.Temperature > 0 {
		out.Temperature = override.Temperature
	}
	if override.TimeoutSeconds > 0 {
		out.TimeoutSeconds = override.TimeoutSeconds
	}
	if override.Streaming {
		out.Streaming = true
	}
	if len(override.Query) > 0 {
		out.Query = mergeStringMap(base.Query, override.Query)
	}
	if len(override.Headers) > 0 {
		out.Headers = mergeStringMap(base.Headers, override.Headers)
	}
	if len(override.Body) > 0 {
		out.Body = mergeStringMap(base.Body, override.Body)
	}
	if len(override.Response) > 0 {
		out.Response = mergeStringMap(base.Response, override.Response)
	}
	if override.Auth != nil {
		out.Auth = override.Auth
	}
	if override.Defaults != nil {
		out.Defaults = mergeAPIDefaults(base.Defaults, override.Defaults)
	}
	if len(override.Multipart) > 0 {
		out.Multipart = override.Multipart
	}
	if override.Poll != nil {
		out.Poll = override.Poll
	}
	return &out
}

func mergeStringMap(base, override map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range base {
		out[k] = v
	}
	for k, v := range override {
		out[k] = v
	}
	return out
}

func mergeAPIDefaults(base, override *APIDefaults) *APIDefaults {
	if override == nil {
		return base
	}
	if base == nil {
		return override
	}
	out := *base
	if override.Model != "" {
		out.Model = override.Model
	}
	if override.MaxTokens > 0 {
		out.MaxTokens = override.MaxTokens
	}
	if override.Temperature > 0 {
		out.Temperature = override.Temperature
	}
	if override.TopP > 0 {
		out.TopP = override.TopP
	}
	if override.TopK > 0 {
		out.TopK = override.TopK
	}
	if override.MinP > 0 {
		out.MinP = override.MinP
	}
	if override.CFGScale > 0 {
		out.CFGScale = override.CFGScale
	}
	if override.FrequencyPenalty > 0 {
		out.FrequencyPenalty = override.FrequencyPenalty
	}
	if override.PresencePenalty > 0 {
		out.PresencePenalty = override.PresencePenalty
	}
	if override.RepetitionPenalty > 0 {
		out.RepetitionPenalty = override.RepetitionPenalty
	}
	if override.Steps > 0 {
		out.Steps = override.Steps
	}
	if override.Width > 0 {
		out.Width = override.Width
	}
	if override.Height > 0 {
		out.Height = override.Height
	}
	return &out
}
