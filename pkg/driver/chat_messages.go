package driver

import (
	"fmt"
	"strings"
)

// ChatMessages returns the effective chat turns for a text request.
// Uses req.Messages when present, otherwise builds from SystemPrompt + Prompt.
func (r TextRequest) ChatMessages() []ChatMessage {
	if len(r.Messages) > 0 {
		return r.Messages
	}
	out := make([]ChatMessage, 0, 2)
	if strings.TrimSpace(r.SystemPrompt) != "" {
		out = append(out, ChatMessage{Role: ChatRoleSystem, Content: r.SystemPrompt})
	}
	if strings.TrimSpace(r.Prompt) != "" {
		out = append(out, ChatMessage{Role: ChatRoleUser, Content: r.Prompt})
	}
	return out
}

// FlattenChatMessages converts a chat history into system + prompt fields for
// completion-style backends (llama-server /completion) that do not accept messages[].
func FlattenChatMessages(msgs []ChatMessage) (systemPrompt, prompt string) {
	var systemParts []string
	var convo []string
	for _, m := range msgs {
		switch m.Role {
		case ChatRoleSystem:
			systemParts = append(systemParts, m.Content)
		case ChatRoleUser:
			convo = append(convo, "User: "+m.Content)
		case ChatRoleAssistant:
			convo = append(convo, "Assistant: "+m.Content)
		case ChatRoleTool:
			label := m.Name
			if label == "" {
				label = "tool"
			}
			convo = append(convo, fmt.Sprintf("Tool(%s): %s", label, m.Content))
		default:
			convo = append(convo, fmt.Sprintf("%s: %s", m.Role, m.Content))
		}
	}
	return strings.Join(systemParts, "\n\n"), strings.Join(convo, "\n")
}

// PromptFields returns system_prompt and prompt for backends that only accept strings.
func (r TextRequest) PromptFields() (systemPrompt, prompt string) {
	if len(r.Messages) > 0 {
		return FlattenChatMessages(r.Messages)
	}
	return r.SystemPrompt, r.Prompt
}
