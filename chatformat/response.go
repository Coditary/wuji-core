package chatformat

import (
	"encoding/json"
	"strings"
)

// ParseToolCallArgs decodes tool arguments JSON into a map for APIs that expect objects.
func ParseToolCallArgs(args string) map[string]any {
	args = strings.TrimSpace(args)
	if args == "" {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(args), &out); err != nil || out == nil {
		return map[string]any{}
	}
	return out
}

// ToolCallsFromOpenAI parses tool_calls from an OpenAI chat completion message.
func ToolCallsFromOpenAI(raw []struct {
	ID       string `json:"id"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}) []ToolCall {
	out := make([]ToolCall, 0, len(raw))
	for i, tc := range raw {
		id := strings.TrimSpace(tc.ID)
		if id == "" {
			id = "call_" + itoa(i)
		}
		out = append(out, ToolCall{
			ID:   id,
			Name: tc.Function.Name,
			Args: tc.Function.Arguments,
		})
	}
	return out
}

// ToolCallsFromAnthropic parses tool_use blocks from Anthropic message content.
func ToolCallsFromAnthropic(blocks []struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Input any    `json:"input"`
}, maps NameMaps) []ToolCall {
	out := make([]ToolCall, 0)
	for _, block := range blocks {
		if block.Type != "tool_use" {
			continue
		}
		args, _ := json.Marshal(block.Input)
		name := RestoreToolName(block.Name, maps)
		out = append(out, ToolCall{
			ID:   block.ID,
			Name: name,
			Args: string(args),
		})
	}
	return out
}

// ToolCallsFromOllama parses Ollama tool_calls into canonical form.
func ToolCallsFromOllama(calls []OllamaToolCall) []ToolCall {
	out := make([]ToolCall, 0, len(calls))
	for i, tc := range calls {
		id := strings.TrimSpace(tc.ID)
		if id == "" {
			id = "call_" + itoa(i)
		}
		args := "{}"
		switch v := tc.Function.Arguments.(type) {
		case nil:
		case string:
			args = strings.TrimSpace(v)
			if args == "" {
				args = "{}"
			}
		default:
			b, err := json.Marshal(v)
			if err == nil {
				args = string(b)
			}
		}
		out = append(out, ToolCall{
			ID:   id,
			Name: tc.Function.Name,
			Args: args,
		})
	}
	return out
}
