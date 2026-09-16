package driver

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/capability"
)

// RunTextInput dispatches a text request to the appropriate generator on a single driver.
func RunTextInput(ctx context.Context, d Driver, req TextRequest) (*TextResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	mode := req.InputModeOrDefault()
	switch mode {
	case TextInputPrompt:
		gen, err := As[TextGenerator](d, capability.TextGeneration)
		if err != nil {
			return nil, err
		}
		return gen.GenerateText(ctx, req.ForGeneration())
	case TextInputImage:
		gen, ok := d.(ImageToTextGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not support --image input", d.Info().ID)
		}
		return gen.ImageToText(ctx, req)
	case TextInputVideo:
		gen, ok := d.(VideoToTextGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not support --video input", d.Info().ID)
		}
		return gen.VideoToText(ctx, req)
	case TextInputDocument:
		gen, ok := d.(DocumentToTextGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not support --document input", d.Info().ID)
		}
		return gen.DocumentToText(ctx, req)
	case TextInputAudio:
		gen, ok := d.(AudioToTextGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not support --audio input", d.Info().ID)
		}
		return gen.AudioToText(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported text input mode %q", mode)
	}
}
