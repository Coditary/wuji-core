package chatformat

import "testing"

func TestSanitizeAddsMissingToolResults(t *testing.T) {
	msgs := []Message{
		{Role: RoleAssistant, ToolCalls: []ToolCall{{ID: "1", Name: "list_dir", Args: `{"path":"."}`}}},
		{Role: RoleUser, Content: "thanks"},
	}
	out := SanitizeMessages(msgs)
	found := false
	for _, m := range out {
		if m.Role == RoleTool && m.ToolCallID == "1" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected synthesized tool result")
	}
}

func TestAnthropicToolNameSanitization(t *testing.T) {
	tools := []ToolDef{{Name: "actions/download", Description: "x", Parameters: map[string]any{"type": "object"}}}
	_, maps := ToolsAnthropic(tools)
	if maps.Forward["actions/download"] == "" {
		t.Fatal("expected sanitized forward mapping")
	}
}

func TestOpenAIToolChoiceMapping(t *testing.T) {
	if OpenAIToolChoice(ToolChoiceRequired, "") != "required" {
		t.Fatal("required")
	}
	if AnthropicToolChoice(ToolChoiceRequired, "", NameMaps{})["type"] != "any" {
		t.Fatal("anthropic required -> any")
	}
}
