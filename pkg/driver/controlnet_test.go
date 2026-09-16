package driver_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
)

func TestCanonicalControlTypeSynonyms(t *testing.T) {
	tests := map[string]driver.ImageControlType{
		"openpose": driver.ImageControlPose,
		"pose":     driver.ImageControlPose,
		"edge":     driver.ImageControlCanny,
		"canny":    driver.ImageControlCanny,
	}
	for in, want := range tests {
		got, err := driver.CanonicalControlType(in)
		if err != nil {
			t.Fatalf("CanonicalControlType(%q): %v", in, err)
		}
		if got != want {
			t.Fatalf("CanonicalControlType(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestBuildControlUnitPoseVariant(t *testing.T) {
	unit, err := driver.BuildControlUnit(driver.ImageControlPose, []string{"full", "pose.png"}, driver.ControlUnitParams{})
	if err != nil {
		t.Fatalf("BuildControlUnit: %v", err)
	}
	if unit.Preprocessor != "full" {
		t.Fatalf("preprocessor = %q, want full", unit.Preprocessor)
	}
	if unit.ImagePath != "pose.png" {
		t.Fatalf("image = %q, want pose.png", unit.ImagePath)
	}
}

func TestMergeControlUnitsLegacy(t *testing.T) {
	units, err := driver.MergeControlUnits(nil, "depth.png", "")
	if err != nil {
		t.Fatalf("MergeControlUnits: %v", err)
	}
	if len(units) != 1 || units[0].Type != driver.ImageControlDepth {
		t.Fatalf("got %+v", units)
	}
}

func TestMergeControlUnitsDuplicate(t *testing.T) {
	_, err := driver.MergeControlUnits([]driver.ControlNetUnit{
		{Type: driver.ImageControlDepth, ImagePath: "a.png"},
	}, "b.png", "depth")
	if err == nil {
		t.Fatal("expected duplicate error")
	}
}

func TestInferImageTaskControlUnits(t *testing.T) {
	task := driver.InferImageTask(driver.ImageTaskInputs{
		Prompt: "person",
		ControlUnits: []driver.ControlNetUnit{
			{Type: driver.ImageControlPose, ImagePath: "pose.png"},
		},
	})
	if task != driver.ImageTaskControlNet {
		t.Fatalf("task = %q, want controlnet", task)
	}

	task = driver.InferImageTask(driver.ImageTaskInputs{
		Prompt: "landscape",
		ControlUnits: []driver.ControlNetUnit{
			{Type: driver.ImageControlDepth, ImagePath: "depth.png"},
		},
	})
	if task != driver.ImageTaskDepth2Img {
		t.Fatalf("task = %q, want depth2img", task)
	}

	task = driver.InferImageTask(driver.ImageTaskInputs{
		Prompt:        "portrait",
		InitImagePath: "photo.png",
		ControlUnits: []driver.ControlNetUnit{
			{Type: driver.ImageControlPose, ImagePath: "pose.png"},
		},
	})
	if task != driver.ImageTaskImg2Img {
		t.Fatalf("task = %q, want img2img", task)
	}
}

func TestControlNetUnitValidateDuplicateType(t *testing.T) {
	req := driver.ImageRequest{
		Task:   driver.ImageTaskControlNet,
		Prompt: "x",
		ControlUnits: []driver.ControlNetUnit{
			{Type: driver.ImageControlPose, ImagePath: "a.png"},
			{Type: driver.ImageControlPose, ImagePath: "b.png"},
		},
	}
	if err := req.Validate(); err == nil {
		t.Fatal("expected duplicate type error")
	}
}
