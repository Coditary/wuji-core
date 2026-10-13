package driver

// TextInputModesFor reports which text input modes a driver implements.
func TextInputModesFor(d Driver) []TextInputMode {
	modes := []TextInputMode{TextInputPrompt}
	if _, ok := d.(ImageToTextGenerator); ok {
		modes = append(modes, TextInputImage)
	}
	if _, ok := d.(VideoToTextGenerator); ok {
		modes = append(modes, TextInputVideo)
	}
	if _, ok := d.(AudioToTextGenerator); ok {
		modes = append(modes, TextInputAudio)
	}
	if _, ok := d.(DocumentToTextGenerator); ok {
		modes = append(modes, TextInputDocument)
	}
	return modes
}

// WithTextInputModes fills TextInputModes when unset.
func (info Info) WithTextInputModes(d Driver) Info {
	if len(info.TextInputModes) == 0 && d != nil {
		info.TextInputModes = TextInputModesFor(d)
	}
	return info
}
