// Package chatformat normalizes tool-calling messages and schemas across
// OpenAI, Anthropic, Ollama, and other backends. The canonical wire shape is
// Message / ToolDef (OpenAI-compatible history).
package chatformat

// WireFormat identifies a provider's tool-calling wire encoding.
type WireFormat string

const (
	FormatOpenAI    WireFormat = "openai"
	FormatAnthropic WireFormat = "anthropic"
	FormatOllama    WireFormat = "ollama"
)

// ToolChoice is the normalized tool selection policy.
type ToolChoice string

const (
	ToolChoiceAuto     ToolChoice = "auto"
	ToolChoiceRequired ToolChoice = "required"
	ToolChoiceNone     ToolChoice = "none"
)

// ParseToolChoice normalizes common spellings.
func ParseToolChoice(raw string) ToolChoice {
	switch raw {
	case "required", "any", "true":
		return ToolChoiceRequired
	case "none", "false":
		return ToolChoiceNone
	default:
		return ToolChoiceAuto
	}
}
