package chatformat

import (
	"strings"
)

const emptyContentPlaceholder = "[empty message content omitted for protocol compatibility]"

// SanitizeMessages fixes common tool-calling history issues before encoding
// for strict providers (Anthropic) or after imperfect client replay.
func SanitizeMessages(msgs []Message) []Message {
	if len(msgs) == 0 {
		return msgs
	}
	out := make([]Message, 0, len(msgs))
	for _, m := range msgs {
		m = sanitizeEmptyContent(m)
		if m.Role == RoleTool && isOrphanedToolResult(out, m) {
			continue
		}
		out = append(out, m)
	}
	return addMissingToolResults(out)
}

func sanitizeEmptyContent(m Message) Message {
	if m.Role != RoleUser && m.Role != RoleAssistant {
		return m
	}
	if strings.TrimSpace(m.Content) != "" || len(m.ToolCalls) > 0 {
		return m
	}
	m.Content = emptyContentPlaceholder
	return m
}

func isOrphanedToolResult(history []Message, toolMsg Message) bool {
	id := strings.TrimSpace(toolMsg.ToolCallID)
	if id == "" {
		return true
	}
	for i := len(history) - 1; i >= 0; i-- {
		m := history[i]
		if m.Role != RoleAssistant || len(m.ToolCalls) == 0 {
			continue
		}
		for _, tc := range m.ToolCalls {
			if tc.ID == id {
				return false
			}
		}
		return true
	}
	return true
}

func addMissingToolResults(msgs []Message) []Message {
	out := make([]Message, 0, len(msgs)+4)
	for i := 0; i < len(msgs); i++ {
		m := msgs[i]
		out = append(out, m)
		if m.Role != RoleAssistant || len(m.ToolCalls) == 0 {
			continue
		}
		seen := map[string]bool{}
		for j := i + 1; j < len(msgs); j++ {
			if msgs[j].Role == RoleTool {
				seen[msgs[j].ToolCallID] = true
			}
		}
		for _, tc := range m.ToolCalls {
			if seen[tc.ID] {
				continue
			}
			out = append(out, Message{
				Role:       RoleTool,
				Name:       tc.Name,
				ToolCallID: tc.ID,
				Content:    "Tool result unavailable (synthesized to satisfy protocol).",
			})
		}
	}
	return out
}
