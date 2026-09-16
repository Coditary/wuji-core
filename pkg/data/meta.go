package data

// Meta holds provenance and task context for a payload.
type Meta struct {
	Capability string            `json:"capability,omitempty"`
	Task       string            `json:"task,omitempty"`
	Driver     string            `json:"driver,omitempty"`
	Model      string            `json:"model,omitempty"`
	Source     map[string]string `json:"source,omitempty"`
	Extra      map[string]string `json:"extra,omitempty"`
}
