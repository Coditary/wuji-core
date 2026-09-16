package catalog

const (
	// IndexVersion is bumped when the on-disk index schema changes.
	IndexVersion = 1
	DefaultSource = "https://models.dev"
)

// Index is the slim runtime catalog stored under .wuji/catalog/index.json.
type Index struct {
	Version   int                       `json:"version"`
	Source    string                    `json:"source"`
	SyncedAt  string                    `json:"synced_at"`
	Providers map[string]ProviderEntry  `json:"providers"`
	Models    map[string]ModelEntry     `json:"models"`
	Labs      map[string]LabEntry       `json:"labs"`
}

// Meta describes the last catalog sync stored in meta.json.
type Meta struct {
	Source         string `json:"source"`
	SyncedAt       string `json:"synced_at"`
	ETag           string `json:"etag,omitempty"`
	ProviderCount  int    `json:"provider_count"`
	ModelCount     int    `json:"model_count"`
	LabCount       int    `json:"lab_count"`
	ProviderModels int    `json:"provider_models"`
}

// ProviderEntry is a models.dev API provider (github-copilot, anthropic, …).
type ProviderEntry struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	API    string            `json:"api,omitempty"`
	NPM    string            `json:"npm,omitempty"`
	Doc    string            `json:"doc,omitempty"`
	Env    []string          `json:"env,omitempty"`
	Models map[string]string `json:"models"` // provider model id -> display name
}

// ModelEntry is a canonical models.dev model (google/gemini-3.8-flash, …).
type ModelEntry struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Lab       string `json:"lab"`
	Family    string `json:"family,omitempty"`
	Reasoning bool   `json:"reasoning,omitempty"`
	ToolCall  bool   `json:"tool_call,omitempty"`
	Context   int    `json:"context,omitempty"`
	Output    int    `json:"output,omitempty"`
	Open      bool   `json:"open_weights,omitempty"`
}

// LabEntry groups canonical models by lab/author prefix.
type LabEntry struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Models []string `json:"models"`
}
