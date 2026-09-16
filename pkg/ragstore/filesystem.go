package ragstore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FilesystemBackend stores collections as JSON + raw vectors on disk.
type FilesystemBackend struct{}

func (FilesystemBackend) Exists(ctx context.Context, root, collection string) bool {
	dir := collectionPathCtx(ctx, root, collection)
	st, err := os.Stat(dir)
	return err == nil && st.IsDir()
}

func (FilesystemBackend) Load(ctx context.Context, root, name string) (*Collection, error) {
	dir := collectionPathCtx(ctx, root, name)
	m, err := LoadManifest(dir)
	if err != nil {
		return nil, fmt.Errorf("load collection %q: %w", name, err)
	}
	chunkData, err := os.ReadFile(filepath.Join(dir, chunksName))
	if err != nil {
		return nil, err
	}
	var chunks []Chunk
	if err := json.Unmarshal(chunkData, &chunks); err != nil {
		return nil, err
	}
	vecData, err := os.ReadFile(filepath.Join(dir, vectorsName))
	if err != nil {
		return nil, err
	}
	vecs := bytesToFloat32(vecData)
	if m.Dims > 0 && len(vecs) != len(chunks)*m.Dims {
		return nil, fmt.Errorf("vector count mismatch")
	}
	if m.Dims > 0 {
		for i := range chunks {
			start := i * m.Dims
			end := start + m.Dims
			if end <= len(vecs) {
				chunks[i].Embedding = append([]float32(nil), vecs[start:end]...)
			}
		}
	}
	return &Collection{Dir: dir, Manifest: m, Chunks: chunks, Vectors: vecs}, nil
}

func (FilesystemBackend) Save(ctx context.Context, root string, m Manifest, chunks []Chunk, vectors []float32) error {
	if strings.TrimSpace(m.Collection) == "" {
		return fmt.Errorf("collection name required")
	}
	dir := collectionPathCtx(ctx, root, m.Collection)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	m.UpdatedAt = time.Now().UTC()
	m.ChunkCount = len(chunks)
	m.Dims = 0
	if len(chunks) > 0 && len(chunks[0].Embedding) > 0 {
		m.Dims = len(chunks[0].Embedding)
	}
	if m.Dims == 0 && len(vectors) > 0 && len(chunks) > 0 {
		m.Dims = len(vectors) / len(chunks)
	}
	chunkJSON, err := json.MarshalIndent(chunks, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, chunksName), chunkJSON, 0o644); err != nil {
		return err
	}
	if len(vectors) == 0 {
		for _, c := range chunks {
			vectors = append(vectors, c.Embedding...)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, vectorsName), float32ToBytes(vectors), 0o644); err != nil {
		return err
	}
	manifestJSON, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, manifestName), manifestJSON, 0o644)
}

func (FilesystemBackend) Delete(ctx context.Context, root, collection string) error {
	return os.RemoveAll(collectionPathCtx(ctx, root, collection))
}

func (FilesystemBackend) ListCollections(ctx context.Context, root string) ([]Manifest, error) {
	base := collectionsBaseCtx(ctx, root)
	entries, err := os.ReadDir(base)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	out := make([]Manifest, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		m, err := LoadManifest(filepath.Join(base, e.Name()))
		if err != nil {
			continue
		}
		out = append(out, m)
	}
	return out, nil
}

func (FilesystemBackend) ExportCollection(ctx context.Context, root, collection, destPath string) error {
	src := collectionPathCtx(ctx, root, collection)
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("export collection %q: %w", collection, err)
	}
	if err := os.MkdirAll(destPath, 0o755); err != nil {
		return err
	}
	return copyDir(src, destPath)
}

func (FilesystemBackend) ImportCollection(ctx context.Context, root, collection, srcPath string) error {
	st, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("import from %q: %w", srcPath, err)
	}
	if !st.IsDir() {
		return fmt.Errorf("import path %q is not a directory", srcPath)
	}
	dest := collectionPathCtx(ctx, root, collection)
	if err := os.RemoveAll(dest); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	return copyDir(srcPath, dest)
}

func (FilesystemBackend) RenameCollection(ctx context.Context, root, collection, newName string) error {
	src := collectionPathCtx(ctx, root, collection)
	dest := collectionPathCtx(ctx, root, newName)
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("rename collection %q: %w", collection, err)
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("collection %q already exists", newName)
	}
	return os.Rename(src, dest)
}

func float32ToBytes(v []float32) []byte {
	b := make([]byte, len(v)*4)
	for i, f := range v {
		u := math.Float32bits(f)
		b[i*4] = byte(u)
		b[i*4+1] = byte(u >> 8)
		b[i*4+2] = byte(u >> 16)
		b[i*4+3] = byte(u >> 24)
	}
	return b
}

func bytesToFloat32(b []byte) []float32 {
	if len(b)%4 != 0 {
		return nil
	}
	out := make([]float32, len(b)/4)
	for i := range out {
		u := uint32(b[i*4]) | uint32(b[i*4+1])<<8 | uint32(b[i*4+2])<<16 | uint32(b[i*4+3])<<24
		out[i] = math.Float32frombits(u)
	}
	return out
}

func copyDir(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
