package driver

import (
	"fmt"
	"strings"
)

// InputModeOrDefault returns the effective input mode.
func (r TextRequest) InputModeOrDefault() TextInputMode {
	if r.InputMode == "" {
		return TextInputPrompt
	}
	return r.InputMode
}

// Validate checks that the request has a valid input configuration.
func (r TextRequest) Validate() error {
	mode := r.InputModeOrDefault()
	switch mode {
	case TextInputPrompt:
		if strings.TrimSpace(r.Prompt) == "" && len(r.Messages) == 0 {
			return fmt.Errorf("text input required: provide a prompt, messages, or --text <file>")
		}
	case TextInputImage, TextInputVideo, TextInputAudio, TextInputDocument:
		if strings.TrimSpace(r.MediaPath) == "" {
			return fmt.Errorf("media path is required for --%s input", mode)
		}
	default:
		return fmt.Errorf("unknown text input mode %q", mode)
	}
	if r.Translate && mode != TextInputPrompt {
		return fmt.Errorf("--translate is only valid with text input (prompt or --text)")
	}
	return nil
}

// ForGeneration returns a copy prepared for LLM text generation (handles --translate).
func (r TextRequest) ForGeneration() TextRequest {
	out := r
	if !out.Translate {
		return out
	}
	target := strings.TrimSpace(out.TargetLang)
	if target == "" {
		target = "English"
	}
	source := strings.TrimSpace(out.Language)
	var instruction string
	if source != "" && source != "auto" {
		instruction = fmt.Sprintf(
			"Translate the following text from %s to %s. Output only the translation, with no explanation or commentary.",
			source, target,
		)
	} else {
		instruction = fmt.Sprintf(
			"Translate the following text to %s. Output only the translation, with no explanation or commentary.",
			target,
		)
	}
	if strings.TrimSpace(out.SystemPrompt) == "" {
		out.SystemPrompt = instruction
	} else {
		out.SystemPrompt = out.SystemPrompt + "\n\n" + instruction
	}
	return out
}
