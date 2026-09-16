package driver

import "testing"

func TestImageVariationRequestProtoRoundtrip(t *testing.T) {
	seed := 42
	in := ImageVariationRequest{
		InitImagePath:     "photo.png",
		DenoisingStrength: 0.35,
		ImageCommonParams: ImageCommonParams{
			NegativePrompt: "blur",
			Model:          "model.safetensors",
			Width:          768,
			Height:         512,
			Steps:          25,
			Sampler:        "euler_a",
			CFGScale:       7,
			BatchSize:      2,
			BatchCount:     3,
			Seed:           &seed,
			ControlUnits: []ControlNetUnit{{
				Type: ImageControlPose, ImagePath: "pose.png", Weight: 0.8,
			}},
			ControlMode: ControlNetModeBalanced,
		},
	}

	out := ImageVariationRequestFromProto(ImageVariationRequestToProto(in))
	if out.InitImagePath != in.InitImagePath {
		t.Fatalf("init = %q", out.InitImagePath)
	}
	if out.DenoisingStrength != in.DenoisingStrength {
		t.Fatalf("denoise = %g", out.DenoisingStrength)
	}
	if out.Width != in.Width || out.Height != in.Height {
		t.Fatalf("size = %dx%d", out.Width, out.Height)
	}
	if out.BatchSize != in.BatchSize || out.BatchCount != in.BatchCount {
		t.Fatalf("batch = %dx%d", out.BatchSize, out.BatchCount)
	}
	if len(out.ControlUnits) != 1 || out.ControlUnits[0].Type != ImageControlPose {
		t.Fatalf("control units = %+v", out.ControlUnits)
	}
	if out.ControlMode != ControlNetModeBalanced {
		t.Fatalf("control mode = %q", out.ControlMode)
	}
}
