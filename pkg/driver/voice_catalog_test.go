package driver

import (
	"context"
	"testing"

	"github.com/coditary/wuji-core/pkg/capability"
)

func TestVoiceSpeechCapabilitiesSupports(t *testing.T) {
	cap := VoiceSpeechCapabilities{
		Languages: []string{"de", "en"}, LanguagesSupported: true,
		Emotions: []string{"happy", "sad"}, EmotionsSupported: true,
		Styles: []string{"conversational"}, StylesSupported: true,
		SpeedMin: 0.5, SpeedMax: 2, PitchMin: -6, PitchMax: 6,
		EnergySupported: true,
	}

	if !cap.SupportsLanguage("de") || !cap.SupportsLanguage("auto") {
		t.Fatal("expected supported language")
	}
	if cap.SupportsLanguage("fr") {
		t.Fatal("expected unsupported language")
	}
	if !cap.SupportsEmotion("happy") || cap.SupportsEmotion("angry") {
		t.Fatal("unexpected emotion support")
	}
	if !cap.SupportsStyle("conversational") || cap.SupportsStyle("news") {
		t.Fatal("unexpected style support")
	}
	if !cap.SupportsSpeed(1) || cap.SupportsSpeed(3) {
		t.Fatal("unexpected speed support")
	}
	if !cap.SupportsPitch(3) || cap.SupportsPitch(9) {
		t.Fatal("unexpected pitch support")
	}
	if !cap.SupportsEnergy(0.8) || !cap.SupportsEnergy(0.5) {
		t.Fatal("expected energy support")
	}

	disabled := VoiceSpeechCapabilities{EnergySupported: false}
	if disabled.SupportsEmotion("happy") || disabled.SupportsStyle("news") || disabled.SupportsLanguage("de") {
		t.Fatal("expected disabled dimensions to reject explicit values")
	}
	if !disabled.SupportsEmotion("") || !disabled.SupportsStyle("") || !disabled.SupportsLanguage("auto") {
		t.Fatal("expected omitted dimensions to stay valid")
	}

	anyEmotion := VoiceSpeechCapabilities{EmotionsSupported: true}
	if !anyEmotion.SupportsEmotion("angry") || !anyEmotion.EmotionsEnabled() {
		t.Fatal("expected supported dimension with open list to accept any value")
	}
}

type voiceCatalogStub struct {
	profiles []VoiceProfileInfo
}

func (s voiceCatalogStub) Info() Info                      { return Info{ID: "stub"} }
func (s voiceCatalogStub) Capabilities() []capability.Type { return nil }
func (s voiceCatalogStub) Close() error                    { return nil }
func (s voiceCatalogStub) ListVoiceProfiles(context.Context) ([]VoiceProfileInfo, error) {
	return s.profiles, nil
}

func TestListVoiceProfiles(t *testing.T) {
	stub := voiceCatalogStub{profiles: []VoiceProfileInfo{
		{Name: "narrator", VoiceID: "voice-narrator", Speech: &VoiceSpeechCapabilities{Emotions: []string{"happy"}, EmotionsSupported: true}},
	}}
	profiles, err := ListVoiceProfiles(context.Background(), stub)
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}

	profile, err := FindVoiceProfile(context.Background(), stub, "narrator")
	if err != nil {
		t.Fatal(err)
	}
	if !profile.SupportsSpeech() || !profile.Speech.SupportsEmotion("happy") {
		t.Fatal("expected narrator speech capabilities")
	}
}
