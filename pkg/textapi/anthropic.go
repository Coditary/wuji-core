package textapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coditary/wuji-core/pkg/config"
)

func completeAnthropic(ctx context.Context, pc config.ProviderConfig, req ChatRequest) (*ChatResponse, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(pc.BaseURL), "/")
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}
	apiKey := strings.TrimSpace(pc.APIKey)
	if apiKey == "" && pc.APIKeyEnv != "" {
		apiKey = strings.TrimSpace(os.Getenv(pc.APIKeyEnv))
	}
	if apiKey == "" {
		return nil, fmt.Errorf("anthropic provider: api_key or api_key_env required")
	}
	model := firstNonEmpty(req.Model, pc.Model)
	if model == "" {
		return nil, fmt.Errorf("anthropic provider: model required")
	}
	version := strings.TrimSpace(pc.AnthropicVersion)
	if version == "" {
		version = "2023-06-01"
	}

	system, messages := splitAnthropicMessages(req.Messages)
	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = pc.MaxTokens
	}
	if maxTokens <= 0 {
		maxTokens = 4096
	}
	body := map[string]any{
		"model":      model,
		"max_tokens": maxTokens,
		"messages":   messages,
	}
	if system != "" {
		body["system"] = system
	}
	if len(req.Tools) > 0 {
		body["tools"] = toAnthropicTools(req.Tools)
	}
	temp := req.Temperature
	if temp <= 0 {
		temp = pc.Temperature
	}
	if temp > 0 {
		body["temperature"] = temp
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/v1/messages", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", apiKey)
	httpReq.Header.Set("anthropic-version", version)

	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("anthropic api %s: %s", resp.Status, strings.TrimSpace(string(respBody)))
	}
	return parseAnthropicResponse(respBody)
}

func splitAnthropicMessages(msgs []Message) (string, []map[string]any) {
	var systemParts []string
	out := make([]map[string]any, 0, len(msgs))
	for _, m := range msgs {
		switch m.Role {
		case RoleSystem:
			systemParts = append(systemParts, m.Content)
		case RoleTool:
			out = append(out, map[string]any{
				"role": "user",
				"content": []map[string]any{{
					"type":        "tool_result",
					"tool_use_id": m.ToolCallID,
					"content":     m.Content,
				}},
			})
		case RoleAssistant:
			out = append(out, map[string]any{"role": "assistant", "content": m.Content})
		default:
			out = append(out, map[string]any{"role": string(m.Role), "content": m.Content})
		}
	}
	return strings.Join(systemParts, "\n\n"), out
}

func toAnthropicTools(tools []ToolDef) []map[string]any {
	out := make([]map[string]any, 0, len(tools))
	for _, t := range tools {
		out = append(out, map[string]any{
			"name":         t.Name,
			"description":  t.Description,
			"input_schema": t.Parameters,
		})
	}
	return out
}

func parseAnthropicResponse(raw []byte) (*ChatResponse, error) {
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
	out := &ChatResponse{
		FinishReason: parsed.StopReason,
		Usage: Usage{
			InputTokens:  parsed.Usage.InputTokens,
			OutputTokens: parsed.Usage.OutputTokens,
		},
	}
	for _, block := range parsed.Content {
		switch block.Type {
		case "text":
			out.Content += block.Text
		case "tool_use":
			args, _ := json.Marshal(block.Input)
			out.ToolCalls = append(out.ToolCalls, ToolCall{
				ID:   block.ID,
				Name: block.Name,
				Args: string(args),
			})
		}
	}
	return out, nil
}
