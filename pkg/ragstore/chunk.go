package ragstore

// Chunk is one indexed text segment with optional embedding.
type Chunk struct {
	ID        string            `json:"id"`
	Text      string            `json:"text"`
	Source    string            `json:"source"`
	ChunkIdx  int               `json:"chunk_index"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Embedding []float32         `json:"-"`
}
