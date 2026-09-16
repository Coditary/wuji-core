package capability

import (
	"fmt"
	"strings"
)

// Type identifies a capability that a driver can provide.
type Type string

const (
	TextGeneration  Type = "text"
	ImageGeneration Type = "image"
	ImageUpscale    Type = "upscale"
	ImageDownscale  Type = "downscale"
	ImageScale      Type = "scale"
	VideoGeneration Type = "video"
	AudioGeneration Type = "audio"
	Mesh            Type = "mesh"
	Video2Audio     Type = "video2audio"
	VoiceCloning    Type = "voice"
	Training        Type = "train"
	DatasetMgmt     Type = "dataset"
	Data            Type = "data"
	RAG             Type = "rag"
)

// All returns every known capability type.
func All() []Type {
	return []Type{
		TextGeneration,
		ImageGeneration,
		ImageUpscale,
		ImageDownscale,
		ImageScale,
		VideoGeneration,
		AudioGeneration,
		Mesh,
		Video2Audio,
		VoiceCloning,
		Training,
		DatasetMgmt,
		Data,
		RAG,
	}
}

func (t Type) String() string {
	return string(t)
}

// Parse normalizes a capability name. "3d" is accepted as a deprecated alias for mesh.
func Parse(name string) (Type, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "3d" {
		return Mesh, nil
	}
	for _, c := range All() {
		if string(c) == name {
			return c, nil
		}
	}
	return "", fmt.Errorf("unknown capability %q", name)
}
