package ragstore

import "time"

// Manifest describes a persisted RAG collection.
type Manifest struct {
	Collection string    `json:"collection"`
	EmbedModel string    `json:"embed_model,omitempty"`
	Dims       int       `json:"dims"`
	ChunkCount int       `json:"chunk_count"`
	ChunkSize  int       `json:"chunk_size,omitempty"`
	Overlap    int       `json:"overlap,omitempty"`
	Sources    []string  `json:"sources,omitempty"`
	UpdatedAt  time.Time `json:"updated_at"`
}
