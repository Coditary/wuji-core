package driver_test

import (
	"context"
	"testing"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
)

type img2imgOnlyDriver struct {
	called bool
	last   driver.ImageToImageRequest
}

func (d *img2imgOnlyDriver) Info() driver.Info {
	return driver.Info{
		ID:         "img2img-only",
		ImageTasks: []driver.ImageTask{driver.ImageTaskImg2Img},
	}
}

func (d *img2imgOnlyDriver) Capabilities() []capability.Type {
	return []capability.Type{capability.ImageGeneration}
}

func (d *img2imgOnlyDriver) Close() error { return nil }

func (d *img2imgOnlyDriver) ImageToImage(_ context.Context, req driver.ImageToImageRequest) (*driver.ImageResponse, error) {
	d.called = true
	d.last = req
	return &driver.ImageResponse{Path: "/tmp/out.png", Format: "png"}, nil
}

func TestRunImageTaskEditFallsBackToImg2Img(t *testing.T) {
	d := &img2imgOnlyDriver{}
	resp, err := driver.RunImageTask(context.Background(), d, driver.ImageRequest{
		Task:          driver.ImageTaskEdit,
		Prompt:        "make it winter",
		InitImagePath: "photo.png",
	})
	if err != nil {
		t.Fatalf("RunImageTask: %v", err)
	}
	if !d.called {
		t.Fatal("expected ImageToImage fallback")
	}
	if d.last.Prompt != "make it winter" || d.last.InitImagePath != "photo.png" {
		t.Fatalf("unexpected img2img request: %+v", d.last)
	}
	if resp.Path != "/tmp/out.png" {
		t.Fatalf("path = %q", resp.Path)
	}
}

func TestResolveImageTaskForDriver(t *testing.T) {
	d := &img2imgOnlyDriver{}
	got := driver.ResolveImageTaskForDriver(d, driver.ImageTaskEdit)
	if got != driver.ImageTaskImg2Img {
		t.Fatalf("got %q, want img2img", got)
	}
}
