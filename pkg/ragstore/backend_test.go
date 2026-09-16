package ragstore_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/coditary/wuji-core/pkg/ragstore"
)

func TestFilesystemBackendRoundTrip(t *testing.T) {
	root := t.TempDir()
	backend := ragstore.FilesystemBackend{}
	ctx := ragstore.WithBackend(context.Background(), backend)

	_ = ctx
	chunks := []ragstore.Chunk{
		{ID: "c1", Text: "hello world", Source: "a.txt", ChunkIdx: 0, Embedding: []float32{1, 0, 0}},
		{ID: "c2", Text: "foo bar", Source: "b.txt", ChunkIdx: 0, Embedding: []float32{0, 1, 0}},
	}
	manifest := ragstore.Manifest{
		Collection: "docs",
		EmbedModel: "test",
		Dims:       3,
		Sources:    []string{"a.txt", "b.txt"},
	}
	if err := backend.Save(context.Background(), root, manifest, chunks, nil); err != nil {
		t.Fatal(err)
	}
	if !backend.Exists(context.Background(), root, "docs") {
		t.Fatal("expected collection to exist")
	}

	col, err := backend.Load(context.Background(), root, "docs")
	if err != nil {
		t.Fatal(err)
	}
	if len(col.Chunks) != 2 || col.Manifest.Dims != 3 {
		t.Fatalf("chunks=%d dims=%d", len(col.Chunks), col.Manifest.Dims)
	}

	exportDir := filepath.Join(root, "export")
	if err := backend.ExportCollection(context.Background(), root, "docs", exportDir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(exportDir, "manifest.json")); err != nil {
		t.Fatal(err)
	}

	if err := backend.Delete(context.Background(), root, "docs"); err != nil {
		t.Fatal(err)
	}
	if backend.Exists(context.Background(), root, "docs") {
		t.Fatal("expected collection deleted")
	}
}

func TestFilesystemBackendTaijiStoreDir(t *testing.T) {
	root := t.TempDir()
	backend := ragstore.FilesystemBackend{}
	ctx := ragstore.WithStoreDir(context.Background(), ".taiji")
	chunks := []ragstore.Chunk{
		{ID: "c1", Text: "hello", Source: "a.md", ChunkIdx: 0, Embedding: []float32{1, 0}},
	}
	manifest := ragstore.Manifest{Collection: "book", Dims: 2}
	if err := backend.Save(ctx, root, manifest, chunks, nil); err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, ".taiji", "rag", "book", "manifest.json")
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("expected taiji store at %s: %v", want, err)
	}
}

func TestResolveBackendDefault(t *testing.T) {
	b := ragstore.ResolveBackend(context.Background())
	if _, ok := b.(ragstore.FilesystemBackend); !ok {
		t.Fatalf("expected filesystem backend, got %T", b)
	}
}
