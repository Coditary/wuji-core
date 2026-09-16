package chatformat

// Role identifies a message author in tool-calling chat history.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message is one turn in canonical (OpenAI-shaped) chat history.
type Message struct {
	Role       Role
	Content    string
	Name       string
	ToolCallID string
	ToolCalls  []ToolCall
}

// ToolDef describes a callable tool.
type ToolDef struct {
	Name        string
	Description string
	Parameters  map[string]any
}

// ToolCall is a model-requested tool invocation.
type ToolCall struct {
	ID   string
	Name string
	Args string
}
