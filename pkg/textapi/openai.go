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

func completeOpenAI(ctx context.Context, pc config.ProviderConfig, req ChatRequest) (*ChatResponse, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(pc.BaseURL), "/")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	apiKey := strings.TrimSpace(pc.APIKey)
	if apiKey == "" && pc.APIKeyEnv != "" {
		apiKey = strings.TrimSpace(os.Getenv(pc.APIKeyEnv))
	}
	if apiKey == "" {
		return nil, fmt.Errorf("openai provider: api_key or api_key_env required")
	}
	model := firstNonEmpty(req.Model, pc.Model)
	if model == "" {
		return nil, fmt.Errorf("openai provider: model required")
	}

	body := map[string]any{
		"model":    model,
		"messages": toOpenAIMessages(req.Messages),
	}
	if len(req.Tools) > 0 {
		body["tools"] = toOpenAITools(req.Tools)
	}
	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = pc.MaxTokens
	}
	if maxTokens > 0 {
		body["max_tokens"] = maxTokens
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
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/chat/completions", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

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
		return nil, fmt.Errorf("openai api %s: %s", resp.Status, strings.TrimSpace(string(respBody)))
	}
	return parseOpenAIResponse(respBody)
}

func toOpenAIMessages(msgs []Message) []map[string]any {
	out := make([]map[string]any, 0, len(msgs))
	for _, m := range msgs {
		item := map[string]any{"role": string(m.Role), "content": m.Content}
		if m.Name != "" {
			item["name"] = m.Name
		}
		if m.ToolCallID != "" {
			item["tool_call_id"] = m.ToolCallID
		}
		out = append(out, item)
	}
	return out
}

func toOpenAITools(tools []ToolDef) []map[string]any {
	out := make([]map[string]any, 0, len(tools))
	for _, t := range tools {
		out = append(out, map[string]any{
			"type": "function",
			"function": map[string]any{
				"name":        t.Name,
				"description": t.Description,
				"parameters":  t.Parameters,
			},
		})
	}
	return out
}

func parseOpenAIResponse(raw []byte) (*ChatResponse, error) {
	var parsed struct {
		Choices []struct {
			Message struct {
				Content   string `json:"content"`
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
	out := &ChatResponse{
		Content:      choice.Message.Content,
		FinishReason: choice.FinishReason,
		Usage: Usage{
			InputTokens:  parsed.Usage.PromptTokens,
			OutputTokens: parsed.Usage.CompletionTokens,
		},
	}
	for _, tc := range choice.Message.ToolCalls {
		out.ToolCalls = append(out.ToolCalls, ToolCall{
			ID:   tc.ID,
			Name: tc.Function.Name,
			Args: tc.Function.Arguments,
		})
	}
	return out, nil
}

func firstNonEmpty(parts ...string) string {
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			return p
		}
	}
	return ""
}
