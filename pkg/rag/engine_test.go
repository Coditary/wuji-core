package rag_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/rag"
	"github.com/coditary/wuji-core/pkg/ragstore"
)

func TestIndexIncrementalReplaceSource(t *testing.T) {
	root := t.TempDir()
	pathA := filepath.Join(root, "a.txt")
	pathB := filepath.Join(root, "b.txt")
	writeFile(t, pathA, "alpha document about cats")
	writeFile(t, pathB, "beta document about dogs")

	ctx := rag.WithEmbedder(context.Background(), rag.PseudoEmbedder)
	req := driver.RAGRequest{StoreRoot: root, Collection: "inc", ChunkSize: 200}

	req.SourcePaths = []string{pathA}
	if _, err := rag.Index(ctx, req); err != nil {
		t.Fatal(err)
	}

	writeFile(t, pathA, "alpha updated content about cats and kittens")
	req.SourcePaths = []string{pathA}
	if _, err := rag.Index(ctx, req); err != nil {
		t.Fatal(err)
	}
	col, err := ragstore.Load(context.Background(), root, "inc")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range col.Chunks {
		if c.Source == pathA && c.Text == "alpha document about cats" {
			t.Fatalf("old chunk text still present: %+v", c)
		}
	}

	req.SourcePaths = []string{pathB}
	if _, err := rag.Index(ctx, req); err != nil {
		t.Fatal(err)
	}
	col2, err := ragstore.Load(context.Background(), root, "inc")
	if err != nil {
		t.Fatal(err)
	}
	hasA, hasB := false, false
	for _, c := range col2.Chunks {
		if c.Source == pathA {
			hasA = true
		}
		if c.Source == pathB {
			hasB = true
		}
	}
	if !hasA || !hasB {
		t.Fatalf("expected both sources, hasA=%t hasB=%t chunks=%d", hasA, hasB, len(col2.Chunks))
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
