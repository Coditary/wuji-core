package driver

// TextInputModesToProto serializes text input modes for driver metadata.
func TextInputModesToProto(modes []TextInputMode) []string {
	if len(modes) == 0 {
		return nil
	}
	out := make([]string, 0, len(modes))
	for _, mode := range modes {
		if mode != "" {
			out = append(out, string(mode))
		}
	}
	return out
}

// TextInputModesFromProto deserializes text input modes from driver metadata.
func TextInputModesFromProto(items []string) []TextInputMode {
	if len(items) == 0 {
		return nil
	}
	out := make([]TextInputMode, 0, len(items))
	for _, item := range items {
		switch TextInputMode(item) {
		case TextInputPrompt, TextInputImage, TextInputVideo, TextInputAudio, TextInputDocument:
			out = append(out, TextInputMode(item))
		}
	}
	return out
}
