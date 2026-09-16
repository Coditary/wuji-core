package chatformat

import (
	"fmt"
	"strings"
)

// MessagesOpenAI encodes canonical chat history for OpenAI-compatible APIs.
func MessagesOpenAI(msgs []Message) []map[string]any {
	msgs = SanitizeMessages(msgs)
	out := make([]map[string]any, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, OpenAIMessage(m))
	}
	return out
}

// OpenAIMessage encodes one canonical message for OpenAI-compatible APIs.
func OpenAIMessage(m Message) map[string]any {
	item := map[string]any{
		"role":    string(m.Role),
		"content": m.Content,
	}
	if m.Name != "" {
		item["name"] = m.Name
	}
	if m.ToolCallID != "" {
		item["tool_call_id"] = m.ToolCallID
	}
	if m.Role == RoleAssistant && len(m.ToolCalls) > 0 {
		item["tool_calls"] = OpenAIToolCallsPayload(m.ToolCalls)
	}
	return item
}

// OpenAIToolCallsPayload builds OpenAI tool_calls entries (arguments as JSON strings).
func OpenAIToolCallsPayload(calls []ToolCall) []map[string]any {
	out := make([]map[string]any, 0, len(calls))
	for i, tc := range calls {
		id := strings.TrimSpace(tc.ID)
		if id == "" {
			id = fmt.Sprintf("call_%d", i)
		}
		args := strings.TrimSpace(tc.Args)
		if args == "" {
			args = "{}"
		}
		out = append(out, map[string]any{
			"id":   id,
			"type": "function",
			"function": map[string]any{
				"name":      tc.Name,
				"arguments": args,
			},
		})
	}
	return out
}

// MessagesAnthropic splits system text and encodes messages for Anthropic API.
func MessagesAnthropic(msgs []Message, maps NameMaps) (string, []map[string]any) {
	msgs = SanitizeMessages(msgs)
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
			if len(m.ToolCalls) > 0 {
				content := make([]map[string]any, 0, 1+len(m.ToolCalls))
				if strings.TrimSpace(m.Content) != "" {
					content = append(content, map[string]any{"type": "text", "text": m.Content})
				}
				for i, tc := range m.ToolCalls {
					id := strings.TrimSpace(tc.ID)
					if id == "" {
						id = fmt.Sprintf("call_%d", i)
					}
					name := tc.Name
					if wire, ok := maps.Forward[name]; ok {
						name = wire
					}
					content = append(content, map[string]any{
						"type":  "tool_use",
						"id":    id,
						"name":  name,
						"input": ParseToolCallArgs(tc.Args),
					})
				}
				out = append(out, map[string]any{"role": "assistant", "content": content})
			} else {
				out = append(out, map[string]any{"role": "assistant", "content": m.Content})
			}
		default:
			out = append(out, map[string]any{"role": string(m.Role), "content": m.Content})
		}
	}
	return strings.Join(systemParts, "\n\n"), out
}

// OllamaMessage is the Ollama /api/chat message shape.
type OllamaMessage struct {
	Role      string           `json:"role"`
	Content   string           `json:"content,omitempty"`
	ToolName  string           `json:"tool_name,omitempty"`
	ToolCalls []OllamaToolCall `json:"tool_calls,omitempty"`
}

// OllamaToolCall mirrors Ollama's tool call wire format.
type OllamaToolCall struct {
	ID       string             `json:"id,omitempty"`
	Type     string             `json:"type,omitempty"`
	Function OllamaFunctionCall `json:"function"`
}

// OllamaFunctionCall is the function payload inside an Ollama tool call.
type OllamaFunctionCall struct {
	Index     int    `json:"index,omitempty"`
	Name      string `json:"name"`
	Arguments any    `json:"arguments"`
}

// MessagesOllama encodes canonical history for Ollama /api/chat.
func MessagesOllama(msgs []Message) []OllamaMessage {
	msgs = SanitizeMessages(msgs)
	out := make([]OllamaMessage, 0, len(msgs))
	for _, m := range msgs {
		if m.Role == RoleTool {
			toolName := strings.TrimSpace(m.Name)
			if toolName == "" {
				toolName = strings.TrimSpace(m.ToolCallID)
			}
			out = append(out, OllamaMessage{
				Role:     "tool",
				ToolName: toolName,
				Content:  m.Content,
			})
			continue
		}
		if m.Role == RoleAssistant && len(m.ToolCalls) > 0 {
			out = append(out, OllamaMessage{
				Role:      "assistant",
				Content:   m.Content,
				ToolCalls: OllamaToolCallsPayload(m.ToolCalls),
			})
			continue
		}
		out = append(out, OllamaMessage{Role: string(m.Role), Content: m.Content})
	}
	return out
}

// OllamaToolCallsPayload builds Ollama tool_calls from canonical calls.
func OllamaToolCallsPayload(calls []ToolCall) []OllamaToolCall {
	out := make([]OllamaToolCall, 0, len(calls))
	for i, tc := range calls {
		id := strings.TrimSpace(tc.ID)
		if id == "" {
			id = fmt.Sprintf("call_%d", i)
		}
		out = append(out, OllamaToolCall{
			ID:   id,
			Type: "function",
			Function: OllamaFunctionCall{
				Index:     i,
				Name:      tc.Name,
				Arguments: ParseToolCallArgs(tc.Args),
			},
		})
	}
	return out
}
