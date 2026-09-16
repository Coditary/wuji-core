package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/chatformat"
)

func completeText(ctx context.Context, spec *config.APICapabilitySpec, req driver.TextRequest) (*driver.TextResponse, error) {
	if spec == nil {
		return nil, fmt.Errorf("text api spec is nil")
	}
	proto := strings.ToLower(strings.TrimSpace(spec.Protocol))
	switch proto {
	case "openai", "openai-chat", "openai-chat-completions":
		return completeOpenAI(ctx, spec, req)
	case "anthropic", "anthropic-messages":
		return completeAnthropic(ctx, spec, req)
	case "http":
		return completeTextHTTP(ctx, spec, req)
	default:
		return nil, fmt.Errorf("unsupported text protocol %q (use openai, anthropic, or http)", spec.Protocol)
	}
}

// applyThinkToOpenAIBody maps YAML think settings to OpenAI-compatible request fields.
// vLLM/Gemma 4 and Ollama ignore raw "think" on /v1/chat/completions; reasoning_effort is required.
func applyThinkToOpenAIBody(body map[string]any, think any) {
	if think == nil {
		return
	}
	switch v := think.(type) {
	case bool:
		if v {
			body["reasoning_effort"] = "low"
		}
	case string:
		s := strings.TrimSpace(strings.ToLower(v))
		if s == "" || s == "false" || s == "none" || s == "off" {
			return
		}
		if s == "true" || s == "on" {
			body["reasoning_effort"] = "low"
			return
		}
		body["reasoning_effort"] = s
	default:
		body["think"] = think
	}
}

func completeOpenAI(ctx context.Context, spec *config.APICapabilitySpec, req driver.TextRequest) (*driver.TextResponse, error) {
	req = mergeTextDefaults(spec, req)
	baseURL := strings.TrimRight(strings.TrimSpace(spec.BaseURL), "/")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if strings.TrimSpace(req.Model) == "" {
		return nil, fmt.Errorf("model is required")
	}
	body := map[string]any{"model": req.Model}
	if len(req.Messages) > 0 {
		body["messages"] = driver.MessagesOpenAI(req.Messages)
	} else {
		msgs := []map[string]any{}
		if strings.TrimSpace(req.SystemPrompt) != "" {
			msgs = append(msgs, map[string]any{"role": "system", "content": req.SystemPrompt})
		}
		msgs = append(msgs, map[string]any{"role": "user", "content": req.Prompt})
		body["messages"] = msgs
	}
	if len(req.Tools) > 0 {
		body["tools"] = driver.ToolsOpenAI(req.Tools)
		if tc := strings.TrimSpace(req.ToolChoice); tc != "" {
			body["tool_choice"] = chatformat.OpenAIToolChoice(chatformat.ParseToolChoice(tc), "")
		}
	}
	if req.MaxTokens > 0 {
		body["max_tokens"] = req.MaxTokens
	}
	if req.Temperature > 0 {
		body["temperature"] = req.Temperature
	}
	if req.TopP > 0 {
		body["top_p"] = req.TopP
	}
	applyThinkToOpenAIBody(body, spec.Think)
	stream := spec.Streaming || req.Stream
	if stream {
		body["stream"] = true
	}
	if req.OnDelta != nil && stream {
		return streamOpenAIChat(ctx, spec, baseURL, body, req)
	}
	return postOpenAIChat(ctx, spec, baseURL, body)
}

func completeAnthropic(ctx context.Context, spec *config.APICapabilitySpec, req driver.TextRequest) (*driver.TextResponse, error) {
	req = mergeTextDefaults(spec, req)
	baseURL := strings.TrimRight(strings.TrimSpace(spec.BaseURL), "/")
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}
	if strings.TrimSpace(req.Model) == "" {
		return nil, fmt.Errorf("model is required")
	}
	version := strings.TrimSpace(spec.AnthropicVersion)
	if version == "" {
		version = "2023-06-01"
	}
	var system string
	var messages []map[string]any
	nameMaps := chatformat.NameMaps{}
	var anthropicTools []map[string]any
	if len(req.Tools) > 0 {
		anthropicTools, nameMaps = driver.ToolsAnthropic(req.Tools)
	}
	if len(req.Messages) > 0 {
		system, messages = driver.MessagesAnthropic(req.Messages, nameMaps)
	} else {
		system = req.SystemPrompt
		messages = []map[string]any{{"role": "user", "content": req.Prompt}}
	}
	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 4096
	}
	body := map[string]any{"model": req.Model, "max_tokens": maxTokens, "messages": messages}
	if system != "" {
		body["system"] = system
	}
	if len(anthropicTools) > 0 {
		body["tools"] = anthropicTools
		if tc := strings.TrimSpace(req.ToolChoice); tc != "" {
			body["tool_choice"] = chatformat.AnthropicToolChoice(chatformat.ParseToolChoice(tc), "", nameMaps)
		}
	}
	if req.Temperature > 0 {
		body["temperature"] = req.Temperature
	}
	if req.TopP > 0 {
		body["top_p"] = req.TopP
	}
	if req.TopK > 0 {
		body["top_k"] = req.TopK
	}
	if spec.Streaming {
		body["stream"] = true
	}
	return postAnthropic(ctx, spec, baseURL, version, body, nameMaps)
}

func completeTextHTTP(ctx context.Context, spec *config.APICapabilitySpec, req driver.TextRequest) (*driver.TextResponse, error) {
	raw, err := doHTTP(ctx, spec, textHTTPVars(req, spec))
	if err != nil {
		return nil, err
	}
	out := &driver.TextResponse{}
	if spec.Streaming {
		var wrapped struct {
			Text string `json:"text"`
		}
		if err := json.Unmarshal(raw, &wrapped); err == nil && wrapped.Text != "" {
			out.Text = wrapped.Text
			return out, nil
		}
	}
	if spec.Response != nil {
		if p := spec.Response["text"]; p != "" {
			out.Text, err = extractJSONPath(raw, p)
			if err != nil {
				return nil, err
			}
		}
		if p := spec.Response["finish_reason"]; p != "" {
			out.FinishReason, _ = extractJSONPath(raw, p)
		}
	}
	if out.Text == "" {
		out.Text = strings.TrimSpace(string(raw))
	}
	return out, nil
}

func postOpenAIChat(ctx context.Context, spec *config.APICapabilitySpec, baseURL string, body map[string]any) (*driver.TextResponse, error) {
	respBody, err := postProtocolJSON(ctx, spec, http.MethodPost, baseURL+"/chat/completions", body, func(req *http.Request, apiKey string) {
		if spec.Auth == nil || strings.TrimSpace(spec.Auth.Type) == "" || strings.EqualFold(spec.Auth.Type, "bearer") {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}
	})
	if err != nil {
		return nil, err
	}
	if spec.Streaming {
		var wrapped struct {
			Text string `json:"text"`
		}
		if err := json.Unmarshal(respBody, &wrapped); err == nil {
			return &driver.TextResponse{Text: wrapped.Text}, nil
		}
	}
	return parseOpenAIResponse(respBody)
}

func streamOpenAIChat(ctx context.Context, spec *config.APICapabilitySpec, baseURL string, body map[string]any, req driver.TextRequest) (*driver.TextResponse, error) {
	apiKey, err := resolveAPIKey(ctx, spec.APIKey, spec.APIKeyEnv)
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	vars := map[string]string{"API_KEY": apiKey}
	url := appendQuery(baseURL+"/chat/completions", spec.Query, vars)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range expandHeaders(spec.Headers, vars, apiKey) {
		httpReq.Header.Set(k, v)
	}
	if apiKey != "" {
		if spec.Auth == nil || strings.TrimSpace(spec.Auth.Type) == "" || strings.EqualFold(spec.Auth.Type, "bearer") {
			httpReq.Header.Set("Authorization", "Bearer "+apiKey)
		} else {
			applyAuth(httpReq, spec, apiKey, vars)
		}
	}
	client := &http.Client{Timeout: httpTimeout(spec)}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, readErr
		}
		return nil, httpAPIError(spec, resp.StatusCode, respBody)
	}

	var text, thinking strings.Builder
	var promptTokens, completionTokens int
	if err := streamSSEWithCallbacks(resp.Body, streamCallbacks{
		OnThink: func(delta string) error {
			thinking.WriteString(delta)
			if req.OnThinkDelta != nil {
				return req.OnThinkDelta(delta)
			}
			return nil
		},
		OnContent: func(delta string) error {
			text.WriteString(delta)
			if req.OnDelta != nil {
				return req.OnDelta(delta)
			}
			return nil
		},
		OnUsage: func(prompt, completion int) {
			if prompt > 0 {
				promptTokens = prompt
			}
			if completion > 0 {
				completionTokens = completion
			}
		},
	}); err != nil {
		return nil, err
	}
	return &driver.TextResponse{
		Text:           text.String(),
		ThinkingBlocks: driver.SingleThinkingBlock(thinking.String(), "openai_reasoning"),
		FinishReason:   "stop",
		InputTokens:    promptTokens,
		TokensUsed:     completionTokens,
	}, nil
}

func postAnthropic(ctx context.Context, spec *config.APICapabilitySpec, baseURL, version string, body map[string]any, maps chatformat.NameMaps) (*driver.TextResponse, error) {
	respBody, err := postProtocolJSON(ctx, spec, http.MethodPost, baseURL+"/v1/messages", body, func(req *http.Request, apiKey string) {
		if spec.Auth == nil || strings.TrimSpace(spec.Auth.Type) == "" {
			req.Header.Set("x-api-key", apiKey)
			req.Header.Set("anthropic-version", version)
		}
	})
	if err != nil {
		return nil, err
	}
	if spec.Streaming {
		var wrapped struct {
			Text string `json:"text"`
		}
		if err := json.Unmarshal(respBody, &wrapped); err == nil {
			return &driver.TextResponse{Text: wrapped.Text}, nil
		}
	}
	return parseAnthropicResponse(respBody, maps)
}

func mergeTextDefaults(spec *config.APICapabilitySpec, req driver.TextRequest) driver.TextRequest {
	if spec == nil {
		return req
	}
	if strings.TrimSpace(req.Model) == "" {
		req.Model = specModel(spec)
	}
	if req.MaxTokens <= 0 {
		req.MaxTokens = specInt(spec, "max_tokens")
	}
	if req.Temperature <= 0 {
		req.Temperature = float32(specFloat(spec, "temperature"))
	}
	if spec.Defaults != nil {
		if req.TopP <= 0 && spec.Defaults.TopP > 0 {
			req.TopP = float32(spec.Defaults.TopP)
		}
		if req.TopK <= 0 && spec.Defaults.TopK > 0 {
			req.TopK = spec.Defaults.TopK
		}
		if req.MinP <= 0 && spec.Defaults.MinP > 0 {
			req.MinP = float32(spec.Defaults.MinP)
		}
		if req.FrequencyPenalty <= 0 && spec.Defaults.FrequencyPenalty > 0 {
			req.FrequencyPenalty = float32(spec.Defaults.FrequencyPenalty)
		}
		if req.PresencePenalty <= 0 && spec.Defaults.PresencePenalty > 0 {
			req.PresencePenalty = float32(spec.Defaults.PresencePenalty)
		}
		if req.RepetitionPenalty <= 0 && spec.Defaults.RepetitionPenalty > 0 {
			req.RepetitionPenalty = float32(spec.Defaults.RepetitionPenalty)
		}
	}
	return req
}

func postProtocolJSON(ctx context.Context, spec *config.APICapabilitySpec, method, url string, body map[string]any, defaultAuth func(*http.Request, string)) ([]byte, error) {
	apiKey, err := resolveAPIKey(ctx, spec.APIKey, spec.APIKeyEnv)
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	vars := map[string]string{"API_KEY": apiKey}
	url = appendQuery(url, spec.Query, vars)
	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range expandHeaders(spec.Headers, vars, apiKey) {
		httpReq.Header.Set(k, v)
	}
	if apiKey != "" {
		if spec.Auth == nil || strings.TrimSpace(spec.Auth.Type) == "" {
			if defaultAuth != nil {
				defaultAuth(httpReq, apiKey)
			} else {
				applyAuth(httpReq, spec, apiKey, vars)
			}
		} else {
			applyAuth(httpReq, spec, apiKey, vars)
		}
	}
	client := &http.Client{Timeout: httpTimeout(spec)}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if spec.Streaming {
		out, err := readSSE(resp.Body)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 400 {
			return nil, httpAPIError(spec, resp.StatusCode, out)
		}
		return out, nil
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, httpAPIError(spec, resp.StatusCode, respBody)
	}
	return respBody, nil
}

func parseOpenAIResponse(raw []byte) (*driver.TextResponse, error) {
	var parsed struct {
		Choices []struct {
			Message struct {
				Content          string `json:"content"`
				ReasoningContent string `json:"reasoning_content"`
				Reasoning        string `json:"reasoning"`
				Thinking         string `json:"thinking"`
				ToolCalls []struct {
					ID       string `json:"id"`
					Function struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	if len(parsed.Choices) == 0 {
		return nil, fmt.Errorf("openai: empty choices")
	}
	choice := parsed.Choices[0]
	resp := &driver.TextResponse{
		Text:         choice.Message.Content,
		FinishReason: choice.FinishReason,
		InputTokens:  parsed.Usage.PromptTokens,
		TokensUsed:   parsed.Usage.CompletionTokens,
	}
	if think := firstNonEmpty(choice.Message.ReasoningContent, choice.Message.Reasoning, choice.Message.Thinking); think != "" {
		resp.ThinkingBlocks = driver.SingleThinkingBlock(think, "openai_reasoning")
	}
	for _, tc := range choice.Message.ToolCalls {
		resp.ToolCalls = append(resp.ToolCalls, driver.ChatToolCall{
			ID: tc.ID, Name: tc.Function.Name, Args: tc.Function.Arguments,
		})
	}
	return resp, nil
}

func parseAnthropicResponse(raw []byte, maps chatformat.NameMaps) (*driver.TextResponse, error) {
	var parsed struct {
		Content []struct {
			Type  string `json:"type"`
			Text  string `json:"text"`
			ID    string `json:"id"`
			Name  string `json:"name"`
			Input any    `json:"input"`
		} `json:"content"`
		StopReason string `json:"stop_reason"`
		Usage      struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	resp := &driver.TextResponse{
		FinishReason: parsed.StopReason,
		InputTokens:  parsed.Usage.InputTokens,
		TokensUsed:   parsed.Usage.OutputTokens,
	}
	for _, block := range parsed.Content {
		switch block.Type {
		case "text":
			resp.Text += block.Text
		case "tool_use":
			args, _ := json.Marshal(block.Input)
			resp.ToolCalls = append(resp.ToolCalls, driver.ChatToolCall{
				ID: block.ID, Name: chatformat.RestoreToolName(block.Name, maps), Args: string(args),
			})
		}
	}
	return resp, nil
}

func firstNonEmpty(parts ...string) string {
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			return p
		}
	}
	return ""
}
