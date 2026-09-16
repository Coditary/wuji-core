package driver

import (
	"context"
	"fmt"
	"strings"
)

// VoiceSpeechCapabilities describes text-to-speech options a voice profile supports.
// When a dimension is supported and its list is empty, any value is accepted.
type VoiceSpeechCapabilities struct {
	Languages          []string
	Emotions           []string
	Styles             []string
	SpeedMin           float32
	SpeedMax           float32
	PitchMin           int
	PitchMax           int
	EnergySupported    bool
	LanguagesSupported bool
	EmotionsSupported  bool
	StylesSupported    bool
}

// VoiceProfileInfo describes a speakable voice and its supported speech parameters.
type VoiceProfileInfo struct {
	Name        string
	VoiceID     string
	Description string
	Source      string // builtin, cloned, model
	Speech      *VoiceSpeechCapabilities
}

// SupportsSpeech reports whether the profile can be used with text-to-speech.
func (p VoiceProfileInfo) SupportsSpeech() bool {
	return p.Speech != nil
}

func (c VoiceSpeechCapabilities) languagesEnabled() bool {
	return c.LanguagesSupported || len(c.Languages) > 0
}

func (c VoiceSpeechCapabilities) emotionsEnabled() bool {
	return c.EmotionsSupported || len(c.Emotions) > 0
}

func (c VoiceSpeechCapabilities) stylesEnabled() bool {
	return c.StylesSupported || len(c.Styles) > 0
}

// LanguagesEnabled reports whether language selection is available for this voice.
func (c VoiceSpeechCapabilities) LanguagesEnabled() bool { return c.languagesEnabled() }

// EmotionsEnabled reports whether emotion presets are available for this voice.
func (c VoiceSpeechCapabilities) EmotionsEnabled() bool { return c.emotionsEnabled() }

// StylesEnabled reports whether style presets are available for this voice.
func (c VoiceSpeechCapabilities) StylesEnabled() bool { return c.stylesEnabled() }

func (c VoiceSpeechCapabilities) supportsValue(requested string, enabled bool, allowed []string) bool {
	requested = strings.ToLower(strings.TrimSpace(requested))
	if requested == "" || requested == "auto" {
		return true
	}
	if !enabled {
		return false
	}
	if len(allowed) == 0 {
		return true
	}
	for _, item := range allowed {
		if strings.EqualFold(strings.TrimSpace(item), requested) {
			return true
		}
	}
	return false
}

// SupportsLanguage reports whether the voice accepts the requested language code.
func (c VoiceSpeechCapabilities) SupportsLanguage(lang string) bool {
	return c.supportsValue(lang, c.languagesEnabled(), c.Languages)
}

// SupportsEmotion reports whether the voice accepts the requested emotion preset.
func (c VoiceSpeechCapabilities) SupportsEmotion(emotion string) bool {
	return c.supportsValue(emotion, c.emotionsEnabled(), c.Emotions)
}

// SupportsStyle reports whether the voice accepts the requested style preset.
func (c VoiceSpeechCapabilities) SupportsStyle(style string) bool {
	return c.supportsValue(style, c.stylesEnabled(), c.Styles)
}

// SupportsSpeed reports whether the requested speed multiplier is in range.
func (c VoiceSpeechCapabilities) SupportsSpeed(speed float32) bool {
	if speed <= 0 {
		return true
	}
	min, max := c.SpeedMin, c.SpeedMax
	if min <= 0 {
		min = 0.25
	}
	if max <= 0 {
		max = 4
	}
	return speed >= min && speed <= max
}

// SupportsPitch reports whether the requested pitch shift is in range.
func (c VoiceSpeechCapabilities) SupportsPitch(pitch int) bool {
	if pitch == 0 {
		return true
	}
	min, max := c.PitchMin, c.PitchMax
	if min == 0 && max == 0 {
		min, max = -12, 12
	}
	return pitch >= min && pitch <= max
}

// SupportsEnergy reports whether energy control is available for this voice.
func (c VoiceSpeechCapabilities) SupportsEnergy(energy float32) bool {
	if energy <= 0 {
		return true
	}
	return c.EnergySupported
}

// ListVoiceProfiles returns speakable voice metadata from a driver when supported.
func ListVoiceProfiles(ctx context.Context, d Driver) ([]VoiceProfileInfo, error) {
	provider, ok := d.(VoiceCatalogProvider)
	if !ok {
		return nil, fmt.Errorf("driver %q does not support voice catalog queries", d.Info().ID)
	}
	return provider.ListVoiceProfiles(ctx)
}

// FindVoiceProfile returns one profile by name or voice id.
func FindVoiceProfile(ctx context.Context, d Driver, nameOrID string) (*VoiceProfileInfo, error) {
	nameOrID = strings.TrimSpace(nameOrID)
	if nameOrID == "" {
		return nil, fmt.Errorf("voice name or id is required")
	}
	profiles, err := ListVoiceProfiles(ctx, d)
	if err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		if strings.EqualFold(profile.Name, nameOrID) || strings.EqualFold(profile.VoiceID, nameOrID) {
			copy := profile
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("voice profile %q not found", nameOrID)
}
