package driver

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ControlNetMode sets prompt vs control priority for all units (backend-specific mapping).
type ControlNetMode string

const (
	ControlNetModeBalanced ControlNetMode = "balanced"
	ControlNetModePrompt   ControlNetMode = "prompt"
	ControlNetModeControl  ControlNetMode = "control"
)

// ControlNetUnit is one ControlNet conditioning input (multi-unit = slice).
type ControlNetUnit struct {
	Type          ImageControlType
	ImagePath     string
	Preprocessor  string
	Weight        float32
	GuidanceStart float32
	GuidanceEnd   float32
	ThresholdA    float32
	ThresholdB    float32
	Model         string
}

var controlTypeSynonyms = map[string]ImageControlType{
	"canny":         ImageControlCanny,
	"edge":          ImageControlCanny,
	"depth":         ImageControlDepth,
	"pose":          ImageControlPose,
	"openpose":      ImageControlPose,
	"lineart":       ImageControlLineart,
	"lineart-anime": ImageControlLineartAnime,
	"scribble":      ImageControlScribble,
	"normal":        ImageControlNormal,
	"segmentation":  ImageControlSegmentation,
	"mlsd":          ImageControlMLSD,
	"tile":          ImageControlTile,
}

var controlPreprocessorOptions = map[ImageControlType][]string{
	ImageControlDepth:        {"midas", "leres", "zoe"},
	ImageControlPose:         {"body", "hand", "face", "full"},
	ImageControlLineart:      {"realistic", "coarse", "standard"},
	ImageControlCanny:        {"hed", "pidinet"},
	ImageControlSegmentation: {"ade20k", "coco"},
	ImageControlScribble:     {"hed", "pidinet"},
}

// CanonicalControlType normalizes a CLI or API control type name (synonyms → canonical).
func CanonicalControlType(raw string) (ImageControlType, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	if normalized == "" {
		return "", fmt.Errorf("control type is required")
	}
	if ct, ok := controlTypeSynonyms[normalized]; ok {
		return ct, nil
	}
	ct := ImageControlType(normalized)
	if !ct.IsValid() {
		return "", fmt.Errorf("unknown control type %q", raw)
	}
	return ct, nil
}

// ParseImageControlType normalizes and validates a control type name (synonyms accepted).
func ParseImageControlType(raw string) (ImageControlType, error) {
	return CanonicalControlType(raw)
}

// ParseControlNetMode normalizes a control mode name.
func ParseControlNetMode(raw string) (ControlNetMode, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	if normalized == "" {
		return "", nil
	}
	mode := ControlNetMode(normalized)
	switch mode {
	case ControlNetModeBalanced, ControlNetModePrompt, ControlNetModeControl:
		return mode, nil
	default:
		return "", fmt.Errorf("unknown control mode %q (valid: balanced, prompt, control)", raw)
	}
}

// BuildControlUnit assembles a unit from image path, optional variant/preprocessor, and params.
func BuildControlUnit(typ ImageControlType, values []string, params ControlUnitParams) (ControlNetUnit, error) {
	unit := ControlNetUnit{Type: typ}
	allowed := controlPreprocessorOptions[typ]

	for _, raw := range values {
		v := strings.ToLower(strings.TrimSpace(raw))
		if v == "" {
			continue
		}
		if isAllowedPreprocessor(v, allowed) {
			if unit.Preprocessor != "" && unit.Preprocessor != v {
				return ControlNetUnit{}, fmt.Errorf("conflicting %s preprocessors %q and %q", typ, unit.Preprocessor, v)
			}
			unit.Preprocessor = v
			continue
		}
		// Image path: last value wins so synonym flags (--canny/--edge) can be used together.
		unit.ImagePath = v
		continue
	}

	unit.Weight = params.Weight
	unit.GuidanceStart = params.GuidanceStart
	unit.GuidanceEnd = params.GuidanceEnd
	unit.ThresholdA = params.ThresholdA
	unit.ThresholdB = params.ThresholdB
	unit.Model = params.Model

	if unit.ImagePath == "" {
		return ControlNetUnit{}, fmt.Errorf("--%s image path is required", typ)
	}
	if unit.Preprocessor != "" && !isAllowedPreprocessor(unit.Preprocessor, allowed) {
		return ControlNetUnit{}, fmt.Errorf("unknown %s preprocessor %q (valid: %s)", typ, unit.Preprocessor, joinStrings(allowed))
	}
	return unit, nil
}

// ControlUnitParams holds optional per-unit tuning from CLI flags.
type ControlUnitParams struct {
	Weight        float32
	GuidanceStart float32
	GuidanceEnd   float32
	ThresholdA    float32
	ThresholdB    float32
	Model         string
	WeightSet     bool
	StartSet      bool
	EndSet        bool
	ThresholdASet bool
	ThresholdBSet bool
}

// MergeControlUnits combines explicit units with legacy single control fields.
func MergeControlUnits(units []ControlNetUnit, legacyImage, legacyType string) ([]ControlNetUnit, error) {
	out := append([]ControlNetUnit(nil), units...)
	if legacyImage == "" {
		return out, nil
	}
	typ := ImageControlDepth
	if strings.TrimSpace(legacyType) != "" {
		parsed, err := ParseImageControlType(legacyType)
		if err != nil {
			return nil, err
		}
		typ = parsed
	}
	for _, u := range out {
		if u.Type == typ {
			return out, fmt.Errorf("duplicate control type %q from legacy --control-image and typed flags", typ)
		}
	}
	out = append(out, ControlNetUnit{Type: typ, ImagePath: legacyImage})
	return out, nil
}

func (u ControlNetUnit) Validate() error {
	if !u.Type.IsValid() {
		return fmt.Errorf("unknown control type %q", u.Type)
	}
	if strings.TrimSpace(u.ImagePath) == "" {
		return fmt.Errorf("control image path is required for type %q", u.Type)
	}
	if u.Preprocessor != "" {
		allowed := controlPreprocessorOptions[u.Type]
		if len(allowed) == 0 {
			return fmt.Errorf("preprocessor %q is not supported for control type %q", u.Preprocessor, u.Type)
		}
		if !isAllowedPreprocessor(u.Preprocessor, allowed) {
			return fmt.Errorf("unknown preprocessor %q for control type %q (valid: %s)", u.Preprocessor, u.Type, joinStrings(allowed))
		}
	}
	return nil
}

func validateControlUnits(units []ControlNetUnit) error {
	if len(units) == 0 {
		return nil
	}
	seen := make(map[ImageControlType]struct{}, len(units))
	for _, u := range units {
		if err := u.Validate(); err != nil {
			return err
		}
		if _, ok := seen[u.Type]; ok {
			return fmt.Errorf("duplicate control type %q (synonyms such as pose/openpose share one slot)", u.Type)
		}
		seen[u.Type] = struct{}{}
	}
	return nil
}

func isAllowedPreprocessor(v string, allowed []string) bool {
	for _, a := range allowed {
		if v == a {
			return true
		}
	}
	return false
}

func looksLikePath(v string) bool {
	if strings.ContainsAny(v, `/\`) {
		return true
	}
	ext := strings.ToLower(filepath.Ext(v))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp", ".bmp", ".tif", ".tiff":
		return true
	}
	return false
}

func joinStrings(items []string) string {
	if len(items) == 0 {
		return "—"
	}
	return strings.Join(items, ", ")
}
