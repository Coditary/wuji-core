package driver

import "github.com/coditary/wuji-core/chatformat"

// OpenAIMessageFromChat converts one chat message to an OpenAI-compatible object.
func OpenAIMessageFromChat(m ChatMessage) map[string]any {
	return chatformat.OpenAIMessage(toChatformatMessages([]ChatMessage{m})[0])
}

// OpenAIToolCallsPayload builds OpenAI tool_calls entries (arguments as JSON strings).
func OpenAIToolCallsPayload(calls []ChatToolCall) []map[string]any {
	out := make([]chatformat.ToolCall, 0, len(calls))
	for _, tc := range calls {
		out = append(out, chatformat.ToolCall{ID: tc.ID, Name: tc.Name, Args: tc.Args})
	}
	return chatformat.OpenAIToolCallsPayload(out)
}

// ParseToolCallArgs decodes tool arguments JSON into a map for APIs that expect objects.
func ParseToolCallArgs(args string) map[string]any {
	return chatformat.ParseToolCallArgs(args)
}
