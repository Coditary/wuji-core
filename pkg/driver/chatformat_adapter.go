package driver

import (
	"github.com/coditary/wuji-core/chatformat"
)

func toChatformatMessages(msgs []ChatMessage) []chatformat.Message {
	out := make([]chatformat.Message, 0, len(msgs))
	for _, m := range msgs {
		item := chatformat.Message{
			Role:       chatformat.Role(m.Role),
			Content:    m.Content,
			Name:       m.Name,
			ToolCallID: m.ToolCallID,
		}
		for _, tc := range m.ToolCalls {
			item.ToolCalls = append(item.ToolCalls, chatformat.ToolCall{
				ID: tc.ID, Name: tc.Name, Args: tc.Args,
			})
		}
		out = append(out, item)
	}
	return out
}

func toChatformatTools(tools []ChatToolDef) []chatformat.ToolDef {
	out := make([]chatformat.ToolDef, 0, len(tools))
	for _, t := range tools {
		out = append(out, chatformat.ToolDef{
			Name: t.Name, Description: t.Description, Parameters: t.Parameters,
		})
	}
	return out
}

// MessagesOpenAI encodes driver messages for OpenAI-compatible APIs.
func MessagesOpenAI(msgs []ChatMessage) []map[string]any {
	return chatformat.MessagesOpenAI(toChatformatMessages(msgs))
}

// ToolsOpenAI encodes driver tool defs for OpenAI-compatible APIs.
func ToolsOpenAI(tools []ChatToolDef) []map[string]any {
	return chatformat.ToolsOpenAI(toChatformatTools(tools))
}

// MessagesAnthropic encodes driver messages for Anthropic API.
func MessagesAnthropic(msgs []ChatMessage, maps chatformat.NameMaps) (string, []map[string]any) {
	return chatformat.MessagesAnthropic(toChatformatMessages(msgs), maps)
}

// ToolsAnthropic encodes driver tool defs for Anthropic API.
func ToolsAnthropic(tools []ChatToolDef) ([]map[string]any, chatformat.NameMaps) {
	return chatformat.ToolsAnthropic(toChatformatTools(tools))
}

// MessagesOllama encodes driver messages for Ollama /api/chat.
func MessagesOllama(msgs []ChatMessage) []chatformat.OllamaMessage {
	return chatformat.MessagesOllama(toChatformatMessages(msgs))
}

// SanitizeChatMessages fixes tool-calling history issues.
func SanitizeChatMessages(msgs []ChatMessage) []ChatMessage {
	canon := toChatformatMessages(msgs)
	canon = chatformat.SanitizeMessages(canon)
	out := make([]ChatMessage, 0, len(canon))
	for _, m := range canon {
		item := ChatMessage{
			Role: ChatRole(m.Role), Content: m.Content, Name: m.Name, ToolCallID: m.ToolCallID,
		}
		for _, tc := range m.ToolCalls {
			item.ToolCalls = append(item.ToolCalls, ChatToolCall{ID: tc.ID, Name: tc.Name, Args: tc.Args})
		}
		out = append(out, item)
	}
	return out
}
