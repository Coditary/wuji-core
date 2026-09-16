package remotemedia_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/coditary/wuji-core/pkg/core/remotemedia"
	"github.com/coditary/wuji-core/pkg/driver"
)

func TestCollectAndRewriteImageRequest(t *testing.T) {
	dir := t.TempDir()
	initPath := filepath.Join(dir, "ref.png")
	loraPath := filepath.Join(dir, "style.safetensors")
	if err := os.WriteFile(initPath, []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(loraPath, []byte("lora"), 0o644); err != nil {
		t.Fatal(err)
	}
	req := driver.ImageRequest{
		InitImagePath: initPath,
		LoRAs:         []driver.LoRARef{{Path: loraPath}},
	}
	paths := remotemedia.CollectImageRequest(req)
	if len(paths) != 2 {
		t.Fatalf("expected 2 path refs, got %d", len(paths))
	}
	mapping := map[string]string{
		initPath: "/srv/staging/ref.png",
		loraPath: "/srv/staging/style.safetensors",
	}
	remotemedia.RewriteImageRequest(&req, mapping)
	if req.InitImagePath != "/srv/staging/ref.png" || req.LoRAs[0].Path != "/srv/staging/style.safetensors" {
		t.Fatalf("rewrite failed: %+v", req)
	}
}

func TestCollectImageTrainDatasetDir(t *testing.T) {
	dir := t.TempDir()
	dataset := filepath.Join(dir, "data")
	if err := os.MkdirAll(dataset, 0o755); err != nil {
		t.Fatal(err)
	}
	paths := remotemedia.CollectImageTrainRequest(driver.ImageTrainRequest{
		DatasetID:    dataset,
		RegDatasetID: dataset,
	})
	if len(paths) != 2 {
		t.Fatalf("expected 2 refs, got %d", len(paths))
	}
}

func TestEnabledForTCPOnly(t *testing.T) {
	if !remotemedia.Enabled("192.168.1.1:50051") {
		t.Fatal("expected TCP endpoint to enable transfer")
	}
	if remotemedia.Enabled("unix:///tmp/core.sock") {
		t.Fatal("expected unix endpoint to skip transfer")
	}
}
