package driver_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
)

func TestParseSpriteAction(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want driver.SpriteAction
	}{
		{"walk", driver.SpriteActionWalk},
		{"walking", driver.SpriteActionWalk},
		{"", ""},
	}
	for _, tc := range cases {
		got, err := driver.ParseSpriteAction(tc.in)
		if err != nil {
			t.Fatalf("ParseSpriteAction(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("ParseSpriteAction(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
	if _, err := driver.ParseSpriteAction("fly"); err == nil {
		t.Fatal("expected error for unknown action")
	}
}

func TestParseSpriteView(t *testing.T) {
	t.Parallel()
	got, err := driver.ParseSpriteView("side-view")
	if err != nil {
		t.Fatal(err)
	}
	if got != driver.SpriteViewSide {
		t.Fatalf("got %q", got)
	}
}

func TestInferImageTaskSpriteWithStyleAndImage(t *testing.T) {
	t.Parallel()
	task := driver.InferImageTask(driver.ImageTaskInputs{
		Prompt:           "walk cycle",
		InitImagePath:    "hero.png",
		StyleImagePath:   "style.png",
		SpriteRequested:  true,
	})
	if task != driver.ImageTaskSprite {
		t.Fatalf("task = %q, want sprite", task)
	}
}

func TestImageRequestSpriteValidateReferencesOptional(t *testing.T) {
	t.Parallel()
	req := driver.ImageRequest{
		Task:               driver.ImageTaskSprite,
		Prompt:             "walk",
		InitImagePath:      "hero.png",
		ReferenceImagePath: "sheet.png",
		StyleImagePath:     "style.png",
		FrameWidth:         64,
		FrameHeight:        64,
		Columns:            4,
		Rows:               2,
		SpriteAction:       driver.SpriteActionWalk,
		SpriteView:         driver.SpriteViewSide,
		SpriteDirections:   4,
	}
	if err := req.Validate(); err != nil {
		t.Fatal(err)
	}
	spriteReq := req.ToSpriteRequest()
	if spriteReq.InitImagePath != "hero.png" || spriteReq.ReferenceImagePath != "sheet.png" || spriteReq.StyleImagePath != "style.png" {
		t.Fatalf("unexpected sprite request: %+v", spriteReq)
	}
}

func TestParseImageTaskSpriteAliases(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{"sprite", "spritesheet", "sprite-sheet"} {
		task, err := driver.ParseImageTask(raw)
		if err != nil {
			t.Fatalf("ParseImageTask(%q): %v", raw, err)
		}
		if task != driver.ImageTaskSprite {
			t.Fatalf("ParseImageTask(%q) = %q", raw, task)
		}
	}
}
