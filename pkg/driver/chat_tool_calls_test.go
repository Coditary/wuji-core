package driver

import "testing"

func TestOpenAIMessageFromChatAssistantToolCalls(t *testing.T) {
	item := OpenAIMessageFromChat(ChatMessage{
		Role:    ChatRoleAssistant,
		Content: "",
		ToolCalls: []ChatToolCall{
			{ID: "call_1", Name: "list_dir", Args: `{"path":"."}`},
		},
	})
	calls, ok := item["tool_calls"].([]map[string]any)
	if !ok || len(calls) != 1 {
		t.Fatalf("tool_calls = %#v", item["tool_calls"])
	}
	fn, ok := calls[0]["function"].(map[string]any)
	if !ok || fn["name"] != "list_dir" || fn["arguments"] != `{"path":"."}` {
		t.Fatalf("function = %#v", calls[0]["function"])
	}
}
