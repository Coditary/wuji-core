package driver

import "context"

// ModelHost is implemented by drivers that load heavy inference backends into VRAM.
type ModelHost interface {
	UnloadInference(ctx context.Context) error
	InferenceLoaded() bool
}

// ModelWarmer can preload a model after eviction restore.
type ModelWarmer interface {
	WarmModel(ctx context.Context, model string) error
}

// UnloadInference stops any loaded inference backend for the driver.
func UnloadInference(ctx context.Context, d Driver) error {
	if h, ok := d.(ModelHost); ok {
		return h.UnloadInference(ctx)
	}
	return nil
}

// InferenceLoaded reports whether the driver currently holds a loaded model.
func InferenceLoaded(d Driver) bool {
	if h, ok := d.(ModelHost); ok {
		return h.InferenceLoaded()
	}
	return false
}
