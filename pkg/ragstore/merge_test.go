package ragstore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/coditary/wuji-core/pkg/ragstore"
)

func TestChunksWithoutSources(t *testing.T) {
	chunks := []ragstore.Chunk{
		{Source: "a.txt", Text: "a"},
		{Source: "b.txt", Text: "b"},
	}
	out := ragstore.ChunksWithoutSources(chunks, []string{"a.txt"})
	if len(out) != 1 || out[0].Source != "b.txt" {
		t.Fatalf("out=%+v", out)
	}
}

func TestReadHTMLSource(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "page.html")
	if err := os.WriteFile(path, []byte("<html><body><p>Hello <b>world</b></p></body></html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	docs, err := ragstore.ReadSources([]string{path}, ragstore.ReadOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(docs) != 1 || docs[0].Content != "Hello world" {
		t.Fatalf("docs=%+v", docs)
	}
}

func TestSplitTextSmartSections(t *testing.T) {
	text := "# Title\n\nIntro paragraph.\n\n## Section\n\nMore text here."
	parts := ragstore.SplitTextSmart(text, 200, 20)
	if len(parts) == 0 {
		t.Fatal("expected chunks")
	}
}
