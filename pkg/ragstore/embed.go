package ragstore

import (
	"fmt"
	"hash/fnv"
)

const defaultEmbedDims = 32

// DefaultEmbedDims returns the pseudo-embedding width for dev drivers.
func DefaultEmbedDims() int { return defaultEmbedDims }

// EmbedTexts returns deterministic pseudo-embeddings for development drivers.
func EmbedTexts(texts []string, dims int) [][]float32 {
	if dims <= 0 {
		dims = defaultEmbedDims
	}
	out := make([][]float32, len(texts))
	for i, text := range texts {
		out[i] = pseudoVector(text, dims)
	}
	return out
}

func pseudoVector(seed string, dims int) []float32 {
	out := make([]float32, dims)
	h := fnv.New32a()
	_, _ = h.Write([]byte(seed))
	base := float32(h.Sum32()%1000) / 1000
	for i := range out {
		out[i] = base + float32(i)*0.003
	}
	return out
}

// BuildChunks splits documents and attaches pseudo-embeddings (legacy helper).
func BuildChunks(docs []SourceDoc, chunkSize, overlap int, dims int) ([]Chunk, []float32) {
	chunks := SplitDocs(docs, chunkSize, overlap)
	if dims <= 0 {
		dims = defaultEmbedDims
	}
	texts := make([]string, len(chunks))
	for i, c := range chunks {
		texts[i] = c.Text
	}
	vecs := EmbedTexts(texts, dims)
	_ = AttachEmbeddings(chunks, vecs, dims)
	return chunks, FlattenEmbeddings(chunks)
}

// SplitDocs turns source documents into chunks without embeddings.
func SplitDocs(docs []SourceDoc, chunkSize, overlap int) []Chunk {
	var chunks []Chunk
	for _, doc := range docs {
		parts := SplitTextSmart(doc.Content, chunkSize, overlap)
		for i, part := range parts {
			id := fmtChunkID(doc.Path, i)
			chunks = append(chunks, Chunk{
				ID: id, Text: part, Source: doc.Path, ChunkIdx: i,
				Metadata: map[string]string{"source": doc.Path},
			})
		}
	}
	return chunks
}

// AttachEmbeddings copies vectors onto chunks and validates dims.
func AttachEmbeddings(chunks []Chunk, vecs [][]float32, dims int) error {
	if len(chunks) != len(vecs) {
		return fmt.Errorf("embedding count %d != chunk count %d", len(vecs), len(chunks))
	}
	for i := range chunks {
		if len(vecs[i]) != dims {
			return fmt.Errorf("chunk %d embedding has %d dims, want %d", i, len(vecs[i]), dims)
		}
		chunks[i].Embedding = append([]float32(nil), vecs[i]...)
	}
	return nil
}

// FlattenEmbeddings returns row-major vectors for all chunks.
func FlattenEmbeddings(chunks []Chunk) []float32 {
	var out []float32
	for _, c := range chunks {
		out = append(out, c.Embedding...)
	}
	return out
}

func fmtChunkID(source string, idx int) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(source))
	return fmt.Sprintf("c%x_%d", h.Sum32(), idx)
}
