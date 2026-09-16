package driver

import (
	"fmt"
	"strings"
)

// AssetMode selects a subject/pipeline profile for mesh and image asset generation.
type AssetMode string

const (
	AssetModeProp        AssetMode = "prop"
	AssetModeHuman       AssetMode = "human"
	AssetModeHardSurface AssetMode = "hardsurface"
	AssetModeTerrain     AssetMode = "terrain"
	AssetModeEnvironment AssetMode = "environment"
	AssetModeFace        AssetMode = "face"
	AssetModePixel       AssetMode = "pixel"
	AssetModeTile        AssetMode = "tile"
	AssetModeUI          AssetMode = "ui"
	AssetModeTexture     AssetMode = "texture"
)

// MeshMode is an alias for AssetMode in mesh commands.
type MeshMode = AssetMode

// ImageMode is an alias for AssetMode in image commands.
type ImageMode = AssetMode

const (
	MeshModeProp        = AssetModeProp
	MeshModeHuman       = AssetModeHuman
	MeshModeHardSurface = AssetModeHardSurface
	MeshModeTerrain     = AssetModeTerrain
	MeshModeEnvironment = AssetModeEnvironment
	MeshModeFace        = AssetModeFace
)

func AllAssetModes() []AssetMode {
	return []AssetMode{
		AssetModeProp,
		AssetModeHuman,
		AssetModeHardSurface,
		AssetModeTerrain,
		AssetModeEnvironment,
		AssetModeFace,
		AssetModePixel,
		AssetModeTile,
		AssetModeUI,
		AssetModeTexture,
	}
}

func (m AssetMode) String() string { return string(m) }

func (m AssetMode) IsValid() bool {
	for _, known := range AllAssetModes() {
		if m == known {
			return true
		}
	}
	return false
}

// ParseAssetMode normalizes a mode name. Empty string defaults to prop.
func ParseAssetMode(raw string) (AssetMode, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	if normalized == "" {
		return AssetModeProp, nil
	}
	switch normalized {
	case "prop", "object", "asset":
		return AssetModeProp, nil
	case "human", "character", "humanoid":
		return AssetModeHuman, nil
	case "hardsurface", "hard-surface", "cad":
		return AssetModeHardSurface, nil
	case "terrain", "landscape":
		return AssetModeTerrain, nil
	case "environment", "env", "scene-env":
		return AssetModeEnvironment, nil
	case "face", "head":
		return AssetModeFace, nil
	case "pixel", "pixelart", "pixel-art":
		return AssetModePixel, nil
	case "tile", "tileset":
		return AssetModeTile, nil
	case "ui", "icon", "icons":
		return AssetModeUI, nil
	case "texture", "material":
		return AssetModeTexture, nil
	}
	mode := AssetMode(normalized)
	if !mode.IsValid() {
		return "", fmt.Errorf("unknown asset mode %q (valid: %s)", raw, joinAssetModes())
	}
	return mode, nil
}

// ParseMeshMode parses a mesh --mode value.
func ParseMeshMode(raw string) (MeshMode, error) {
	return ParseAssetMode(raw)
}

// ParseImageMode parses an image --mode value.
func ParseImageMode(raw string) (ImageMode, error) {
	return ParseAssetMode(raw)
}

func (r MeshRequest) ModeOrDefault() AssetMode {
	if r.Mode == "" {
		return AssetModeProp
	}
	return r.Mode
}

func (r ImageRequest) ModeOrDefault() AssetMode {
	if r.Mode == "" {
		return AssetModeProp
	}
	return r.Mode
}

func (p ImageCommonParams) ModeOrDefault() AssetMode {
	if p.Mode == "" {
		return AssetModeProp
	}
	return p.Mode
}

func (r MeshResponse) ModeOrDefault() AssetMode {
	if r.Mode == "" {
		return AssetModeProp
	}
	return r.Mode
}

func (r MeshResponse) RepresentationOrDefault() MeshRepresentation {
	if r.Representation == "" {
		return MeshRepresentationMesh
	}
	return r.Representation
}

func joinAssetModes() string {
	parts := make([]string, len(AllAssetModes()))
	for i, m := range AllAssetModes() {
		parts[i] = string(m)
	}
	return strings.Join(parts, ", ")
}
