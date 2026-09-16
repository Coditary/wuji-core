package driver_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
)

func TestParseVideoTask(t *testing.T) {
	tests := []struct {
		in   string
		want driver.VideoTask
	}{
		{"", driver.VideoTaskGenerate},
		{"generate", driver.VideoTaskGenerate},
		{"t2v", driver.VideoTaskGenerate},
		{"i2v", driver.VideoTaskImageToVideo},
		{"interpolate", driver.VideoTaskInterpolate},
		{"scale", driver.VideoTaskScale},
		{"upscale", driver.VideoTaskScale},
	}

	for _, tc := range tests {
		got, err := driver.ParseVideoTask(tc.in)
		if err != nil {
			t.Fatalf("ParseVideoTask(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("ParseVideoTask(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestVideoRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     driver.VideoRequest
		wantErr bool
	}{
		{
			name:    "generate requires prompt",
			req:     driver.VideoRequest{Task: driver.VideoTaskGenerate},
			wantErr: true,
		},
		{
			name: "generate ok",
			req:  driver.VideoRequest{Task: driver.VideoTaskGenerate, Prompt: "ocean waves"},
		},
		{
			name: "i2v requires init image",
			req: driver.VideoRequest{
				Task: driver.VideoTaskImageToVideo, Prompt: "animate",
			},
			wantErr: true,
		},
		{
			name: "i2v ok",
			req: driver.VideoRequest{
				Task: driver.VideoTaskImageToVideo, Prompt: "animate", InitImagePath: "frame.png",
			},
		},
		{
			name: "interpolate requires input video",
			req: driver.VideoRequest{
				Task: driver.VideoTaskInterpolate,
			},
			wantErr: true,
		},
		{
			name: "interpolate ok",
			req: driver.VideoRequest{
				Task: driver.VideoTaskInterpolate, InitVideoPath: "clip.mp4",
			},
		},
		{
			name: "scale invalid scale",
			req: driver.VideoRequest{
				Task: driver.VideoTaskScale, InitVideoPath: "clip.mp4", Scale: -1,
			},
			wantErr: true,
		},
		{
			name: "scale requires factor",
			req: driver.VideoRequest{
				Task: driver.VideoTaskScale, InitVideoPath: "clip.mp4",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.Validate()
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestParseCameraControl(t *testing.T) {
	got, err := driver.ParseCameraControl("zoom_in")
	if err != nil || got != driver.CameraControlZoomIn {
		t.Fatalf("zoom_in: got %q err %v", got, err)
	}
	if _, err := driver.ParseCameraControl("orbit"); err == nil {
		t.Fatal("expected error for unknown preset")
	}
}

func TestInferVideoTask(t *testing.T) {
	got, err := driver.InferVideoTask(driver.VideoTaskInputs{})
	if err != nil || got != driver.VideoTaskGenerate {
		t.Fatalf("empty: got %q err %v", got, err)
	}

	got, err = driver.InferVideoTask(driver.VideoTaskInputs{ImagePath: "f.png"})
	if err != nil || got != driver.VideoTaskImageToVideo {
		t.Fatalf("image: got %q err %v", got, err)
	}

	got, err = driver.InferVideoTask(driver.VideoTaskInputs{VideoPath: "c.mp4", FPSSet: true, ScaleSet: false})
	if err != nil || got != driver.VideoTaskInterpolate {
		t.Fatalf("interpolate: got %q err %v", got, err)
	}

	got, err = driver.InferVideoTask(driver.VideoTaskInputs{VideoPath: "c.mp4", ScaleSet: true, Scale: 4})
	if err != nil || got != driver.VideoTaskScale {
		t.Fatalf("scale: got %q err %v", got, err)
	}

	if _, err := driver.InferVideoTask(driver.VideoTaskInputs{ImagePath: "a.png", VideoPath: "c.mp4"}); err == nil {
		t.Fatal("expected error for image and video together")
	}
	if _, err := driver.InferVideoTask(driver.VideoTaskInputs{VideoPath: "c.mp4"}); err == nil {
		t.Fatal("expected error for video without fps or scale")
	}
}

func TestVideoInterpolateProtoRoundTrip(t *testing.T) {
	req := driver.VideoInterpolateRequest{
		InputVideoPath: "clip.mp4", TargetFPS: 60, Model: "rife",
	}

	out := driver.VideoInterpolateRequestFromProto(driver.VideoInterpolateRequestToProto(req))
	if out.InputVideoPath != req.InputVideoPath {
		t.Fatalf("path: got %q want %q", out.InputVideoPath, req.InputVideoPath)
	}
	if out.TargetFPS != req.TargetFPS {
		t.Fatalf("fps: got %d want %d", out.TargetFPS, req.TargetFPS)
	}
}
