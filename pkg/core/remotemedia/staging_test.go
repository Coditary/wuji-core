package remotemedia_test

import (
	"os"
	"path/filepath"
	"testing"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/core/remotemedia"
)

func TestWriteAndReadStaging(t *testing.T) {
	root := t.TempDir()
	raw := []byte("hello world repeated for compression test " + string(make([]byte, 600)))
	compressed, zstdFlag, err := remotemedia.BlobData(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !zstdFlag {
		t.Fatal("expected zstd for large payload")
	}

	paths, stagingRoot, err := remotemedia.WriteStaging(root, "sess1", []*wujiv1.FileBlob{
		{Name: "subdir/a.png", Data: compressed, Zstd: zstdFlag},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(paths))
	}
	want := filepath.Join(root, ".wuji", "staging", "sess1", "subdir", "a.png")
	if paths[0] != want {
		t.Fatalf("path = %q want %q", paths[0], want)
	}
	if stagingRoot != filepath.Join(root, ".wuji", "staging", "sess1") {
		t.Fatalf("staging root = %q", stagingRoot)
	}

	files, err := remotemedia.ReadFiles([]string{want})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	out, err := remotemedia.MaybeDecompress(files[0].GetData(), files[0].GetZstd())
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != string(raw) {
		t.Fatal("round-trip mismatch")
	}
}

func TestBuildUploadPlanDirectory(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "captions")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "a.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	entries, dirs, err := remotemedia.BuildUploadPlan([]string{sub})
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 1 || len(entries) != 1 {
		t.Fatalf("entries=%d dirs=%d", len(entries), len(dirs))
	}
	if entries[0].RootPath != sub {
		t.Fatalf("root path = %q", entries[0].RootPath)
	}
}
