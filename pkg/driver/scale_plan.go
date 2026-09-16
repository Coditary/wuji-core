package driver

// ScaleKind classifies a requested scale factor for orchestration.
type ScaleKind int

const (
	ScaleKindDownscale ScaleKind = iota
	ScaleKindUpscaleExact
	ScaleKindHybrid
)

// ClassifyScale reports whether the factor is pure downscale, exact AI upscale (2/4/8), or hybrid.
func ClassifyScale(factor float32) ScaleKind {
	if factor < 1 {
		return ScaleKindDownscale
	}
	if nearOne(factor) {
		return ScaleKindUpscaleExact
	}
	switch factor {
	case 2, 4, 8:
		return ScaleKindUpscaleExact
	default:
		return ScaleKindHybrid
	}
}

// ScaleDrivers holds backends for coordinated scaling.
type ScaleDrivers struct {
	Image     Driver
	Upscale   Driver
	Downscale Driver
	Scale     Driver
}
