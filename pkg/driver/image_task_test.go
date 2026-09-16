package driver_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
)

func TestParseImageTask(t *testing.T) {
	tests := []struct {
		in   string
		want driver.ImageTask
	}{
		{"", driver.ImageTaskGenerate},
		{"generate", driver.ImageTaskGenerate},
		{"txt2img", driver.ImageTaskGenerate},
		{"img2img", driver.ImageTaskImg2Img},
		{"style_transfer", driver.ImageTaskStyleTransfer},
		{"ip-adapter", driver.ImageTaskStyleTransfer},
	}

	for _, tc := range tests {
		got, err := driver.ParseImageTask(tc.in)
		if err != nil {
			t.Fatalf("ParseImageTask(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("ParseImageTask(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestImageRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     driver.ImageRequest
		wantErr bool
	}{
		{
			name: "generate requires prompt",
			req:  driver.ImageRequest{Task: driver.ImageTaskGenerate},
			wantErr: true,
		},
		{
			name: "generate ok",
			req:  driver.ImageRequest{Task: driver.ImageTaskGenerate, Prompt: "cat"},
		},
		{
			name: "inpaint requires mask",
			req: driver.ImageRequest{
				Task: driver.ImageTaskInpaint, Prompt: "fill", InitImagePath: "a.png",
			},
			wantErr: true,
		},
		{
			name: "inpaint ok",
			req: driver.ImageRequest{
				Task: driver.ImageTaskInpaint, Prompt: "fill",
				InitImagePath: "a.png", MaskImagePath: "m.png",
			},
		},
		{
			name: "controlnet requires unit",
			req: driver.ImageRequest{
				Task: driver.ImageTaskControlNet, Prompt: "person",
			},
			wantErr: true,
		},
		{
			name: "controlnet ok",
			req: driver.ImageRequest{
				Task: driver.ImageTaskControlNet, Prompt: "person",
				ControlUnits: []driver.ControlNetUnit{
					{Type: driver.ImageControlCanny, ImagePath: "c.png"},
				},
			},
		},
		{
			name: "upscale invalid scale",
			req: driver.ImageRequest{
				Task: driver.ImageTaskUpscale, InitImagePath: "a.png", Scale: -1,
			},
			wantErr: true,
		},
		{
			name: "style-transfer requires style image",
			req: driver.ImageRequest{
				Task: driver.ImageTaskStyleTransfer, Prompt: "portrait",
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

func TestImageInpaintProtoRoundTrip(t *testing.T) {
	req := driver.ImageInpaintRequest{
		Prompt: "fix sky", InitImagePath: "photo.png", MaskImagePath: "mask.png",
		ImageCommonParams: driver.ImageCommonParams{Width: 512, Height: 512},
	}

	out := driver.ImageInpaintRequestFromProto(driver.ImageInpaintRequestToProto(req))
	if out.MaskImagePath != req.MaskImagePath {
		t.Fatalf("mask: got %q want %q", out.MaskImagePath, req.MaskImagePath)
	}
	if out.Width != req.Width {
		t.Fatalf("width: got %d want %d", out.Width, req.Width)
	}
}
