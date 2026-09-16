package driver_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
)

func TestInferImageTask(t *testing.T) {
	tests := []struct {
		name string
		in   driver.ImageTaskInputs
		want driver.ImageTask
	}{
		{
			name: "prompt only",
			in:   driver.ImageTaskInputs{Prompt: "cat"},
			want: driver.ImageTaskGenerate,
		},
		{
			name: "edit flag",
			in: driver.ImageTaskInputs{
				Prompt:        "make it winter",
				InitImagePath: "a.png",
				EditRequested: true,
			},
			want: driver.ImageTaskEdit,
		},
		{
			name: "edit flag loses to mask",
			in: driver.ImageTaskInputs{
				Prompt:        "boat",
				InitImagePath: "a.png",
				MaskImagePath: "m.png",
				EditRequested: true,
			},
			want: driver.ImageTaskInpaint,
		},
		{
			name: "init image and prompt",
			in: driver.ImageTaskInputs{
				Prompt:        "watercolor",
				InitImagePath: "a.png",
			},
			want: driver.ImageTaskImg2Img,
		},
		{
			name: "mask implies inpaint",
			in: driver.ImageTaskInputs{
				Prompt:        "boat",
				InitImagePath: "a.png",
				MaskImagePath: "m.png",
			},
			want: driver.ImageTaskInpaint,
		},
		{
			name: "style image",
			in: driver.ImageTaskInputs{
				Prompt:         "portrait",
				StyleImagePath: "style.png",
			},
			want: driver.ImageTaskStyleTransfer,
		},
		{
			name: "controlnet",
			in: driver.ImageTaskInputs{
				Prompt:           "person",
				ControlImagePath: "canny.png",
				ControlType:      "canny",
			},
			want: driver.ImageTaskControlNet,
		},
		{
			name: "depth2img",
			in: driver.ImageTaskInputs{
				Prompt:           "landscape",
				ControlImagePath: "depth.png",
			},
			want: driver.ImageTaskDepth2Img,
		},
		{
			name: "upscale",
			in: driver.ImageTaskInputs{
				InitImagePath: "a.png",
				Scale:         4,
			},
			want: driver.ImageTaskUpscale,
		},
		{
			name: "upscale flag",
			in: driver.ImageTaskInputs{
				InitImagePath:    "a.png",
				UpscaleRequested: true,
				Scale:            4,
			},
			want: driver.ImageTaskUpscale,
		},
		{
			name: "upscale flag with prompt",
			in: driver.ImageTaskInputs{
				Prompt:           "ignored for upscale",
				InitImagePath:    "a.png",
				UpscaleRequested: true,
				Scale:            4,
			},
			want: driver.ImageTaskUpscale,
		},
		{
			name: "variation",
			in: driver.ImageTaskInputs{
				InitImagePath: "a.png",
			},
			want: driver.ImageTaskVariation,
		},
		{
			name: "variation with control",
			in: driver.ImageTaskInputs{
				InitImagePath: "a.png",
				ControlUnits: []driver.ControlNetUnit{{
					Type: driver.ImageControlPose, ImagePath: "pose.png",
				}},
			},
			want: driver.ImageTaskVariation,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := driver.InferImageTask(tc.in)
			if got != tc.want {
				t.Fatalf("InferImageTask() = %q, want %q", got, tc.want)
			}
		})
	}
}
