package ragstore

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const (
	dirName      = ".wuji"
	ragDirName   = "rag"
	manifestName = "manifest.json"
	chunksName   = "chunks.json"
	vectorsName  = "vectors.f32"
)

type storeDirKey struct{}

// WithStoreDir attaches a state directory name (e.g. ".wuji", ".taiji") to ctx for collection paths.
func WithStoreDir(ctx context.Context, dir string) context.Context {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return ctx
	}
	return context.WithValue(ctx, storeDirKey{}, filepath.Clean(dir))
}

func storeDirFrom(ctx context.Context) string {
	if ctx == nil {
		return dirName
	}
	if d, ok := ctx.Value(storeDirKey{}).(string); ok && strings.TrimSpace(d) != "" {
		return d
	}
	return dirName
}

// CollectionPath returns the on-disk path for a collection using the default state dir (.wuji).
func CollectionPath(root, collection string) string {
	return collectionPath(dirName, root, collection)
}

func collectionPath(stateDir, root, collection string) string {
	if strings.TrimSpace(stateDir) == "" {
		stateDir = dirName
	}
	return filepath.Join(root, stateDir, ragDirName, sanitizeName(collection))
}

func collectionPathCtx(ctx context.Context, root, collection string) string {
	return collectionPath(storeDirFrom(ctx), root, collection)
}

func collectionsBaseCtx(ctx context.Context, root string) string {
	return filepath.Join(root, storeDirFrom(ctx), ragDirName)
}

// Collection holds loaded chunks and vectors.
type Collection struct {
	Dir      string
	Manifest Manifest
	Chunks   []Chunk
	Vectors  []float32
}

// LoadManifest reads manifest.json for a collection directory.
func LoadManifest(dir string) (Manifest, error) {
	data, err := os.ReadFile(filepath.Join(dir, manifestName))
	if err != nil {
		return Manifest{}, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

func sanitizeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, string(os.PathSeparator), "_")
	return name
}

// Exists reports whether a collection directory exists.
func Exists(ctx context.Context, root, collection string) bool {
	return DefaultBackend.Exists(ctx, root, collection)
}

// Load opens a collection from disk.
func Load(ctx context.Context, root, name string) (*Collection, error) {
	return DefaultBackend.Load(ctx, root, name)
}

// Save persists chunks and vectors to disk.
func Save(ctx context.Context, root string, m Manifest, chunks []Chunk, vectors []float32) error {
	return DefaultBackend.Save(ctx, root, m, chunks, vectors)
}

// Delete removes a collection directory.
func Delete(ctx context.Context, root, collection string) error {
	return DefaultBackend.Delete(ctx, root, collection)
}

// ListCollections returns collection names under project root.
func ListCollections(ctx context.Context, root string) ([]Manifest, error) {
	return DefaultBackend.ListCollections(ctx, root)
}
