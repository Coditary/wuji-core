package textapi

import (
	"fmt"
	"strings"

	"github.com/coditary/wuji-core/pkg/driver"
)

func flattenMessages(msgs []Message) (system string, prompt string) {
	var systemParts []string
	var convo []string
	for _, m := range msgs {
		switch m.Role {
		case RoleSystem:
			systemParts = append(systemParts, m.Content)
		case RoleUser:
			convo = append(convo, "User: "+m.Content)
		case RoleAssistant:
			convo = append(convo, "Assistant: "+m.Content)
		case RoleTool:
			label := m.Name
			if label == "" {
				label = "tool"
			}
			convo = append(convo, fmt.Sprintf("Tool(%s): %s", label, m.Content))
		}
	}
	return strings.Join(systemParts, "\n\n"), strings.Join(convo, "\n")
}

func toTextRequest(req ChatRequest, defaultModel string) driver.TextRequest {
	system, prompt := flattenMessages(req.Messages)
	out := driver.TextRequest{
		Prompt:       prompt,
		SystemPrompt: system,
		Model:        firstNonEmpty(req.Model, defaultModel),
	}
	if req.MaxTokens > 0 {
		out.MaxTokens = req.MaxTokens
	}
	if req.Temperature > 0 {
		out.Temperature = float32(req.Temperature)
	}
	return out
}
