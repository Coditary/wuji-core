package config

import "testing"

func TestLoRAConfigResolveEntry(t *testing.T) {
	cfg := LoRAConfig{
		DefaultWeight: 1,
		Library: map[string]LoRALibraryEntry{
			"style":  {Path: "anime_style.safetensors", Weight: 0.8},
			"detail": {Path: "/models/Lora/detail.safetensors"},
		},
	}

	ref, err := cfg.ResolveEntry("style", 1)
	if err != nil {
		t.Fatal(err)
	}
	if ref.Path != "anime_style.safetensors" || ref.Weight != 0.8 {
		t.Fatalf("style = %+v", ref)
	}

	ref, err = cfg.ResolveEntry("/tmp/custom.safetensors", 0.5)
	if err != nil {
		t.Fatal(err)
	}
	if ref.Path != "/tmp/custom.safetensors" || ref.Weight != 0.5 {
		t.Fatalf("custom = %+v", ref)
	}
}

func TestLoRAConfigResolveDefault(t *testing.T) {
	cfg := LoRAConfig{
		Library: map[string]LoRALibraryEntry{
			"style":  {Path: "anime_style.safetensors"},
			"detail": {Path: "detail.safetensors"},
		},
		Default: []string{"style", "detail"},
	}
	refs, err := cfg.ResolveDefault(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(refs) != 2 {
		t.Fatalf("got %d refs", len(refs))
	}
}
