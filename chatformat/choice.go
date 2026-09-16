package chatformat

// OpenAIToolChoice encodes tool_choice for OpenAI / Ollama / Cursor-compatible APIs.
func OpenAIToolChoice(choice ToolChoice, toolName string) any {
	switch choice {
	case ToolChoiceRequired:
		return "required"
	case ToolChoiceNone:
		return "none"
	case ToolChoiceAuto:
		if toolName != "" {
			return map[string]any{
				"type": "function",
				"function": map[string]any{
					"name": toolName,
				},
			}
		}
		return "auto"
	default:
		return "auto"
	}
}

// AnthropicToolChoice encodes tool_choice for Anthropic Messages API.
func AnthropicToolChoice(choice ToolChoice, toolName string, maps NameMaps) map[string]any {
	switch choice {
	case ToolChoiceRequired:
		return map[string]any{"type": "any"}
	case ToolChoiceNone:
		return map[string]any{"type": "none"}
	case ToolChoiceAuto:
		if toolName != "" {
			if wire, ok := maps.Forward[toolName]; ok {
				toolName = wire
			}
			return map[string]any{"type": "tool", "name": toolName}
		}
		return map[string]any{"type": "auto"}
	default:
		return map[string]any{"type": "auto"}
	}
}
