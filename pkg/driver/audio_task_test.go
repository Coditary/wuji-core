package driver_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
)

func TestParseAudioTask(t *testing.T) {
	tests := []struct {
		in   string
		want driver.AudioTask
	}{
		{"", driver.AudioTaskTextToMusic},
		{"text-to-music", driver.AudioTaskTextToMusic},
		{"melody-to-music", driver.AudioTaskMelodyToMusic},
		{"text-to-sfx", driver.AudioTaskTextToSFX},
		{"text-to-speech", driver.AudioTaskTextToSpeech},
		{"tts", driver.AudioTaskTextToSpeech},
	}

	for _, tc := range tests {
		got, err := driver.ParseAudioTask(tc.in)
		if err != nil {
			t.Fatalf("ParseAudioTask(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("ParseAudioTask(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestAudioRequestValidate(t *testing.T) {
	err := (driver.AudioRequest{Task: driver.AudioTaskMelodyToMusic, Prompt: "song"}).Validate()
	if err == nil {
		t.Fatal("expected error for missing reference")
	}
}

func TestInferAudioTask(t *testing.T) {
	got, err := driver.InferAudioTask(driver.AudioTaskInputs{})
	if err != nil || got != driver.AudioTaskTextToMusic {
		t.Fatalf("default: got %q err %v", got, err)
	}

	got, err = driver.InferAudioTask(driver.AudioTaskInputs{ReferencePath: "melody.mid"})
	if err != nil || got != driver.AudioTaskMelodyToMusic {
		t.Fatalf("melody: got %q err %v", got, err)
	}

	got, err = driver.InferAudioTask(driver.AudioTaskInputs{SFXRequested: true})
	if err != nil || got != driver.AudioTaskTextToSFX {
		t.Fatalf("sfx: got %q err %v", got, err)
	}

	if _, err := driver.InferAudioTask(driver.AudioTaskInputs{
		ReferencePath: "melody.mid", SFXRequested: true,
	}); err == nil {
		t.Fatal("expected error for reference and sfx together")
	}

	got, err = driver.InferAudioTask(driver.AudioTaskInputs{VoiceName: "myvoice"})
	if err != nil || got != driver.AudioTaskTextToSpeech {
		t.Fatalf("speech via voice: got %q err %v", got, err)
	}

	got, err = driver.InferAudioTask(driver.AudioTaskInputs{SpeechRequested: true})
	if err != nil || got != driver.AudioTaskTextToSpeech {
		t.Fatalf("speech flag: got %q err %v", got, err)
	}
}

func TestAudioSpeechProtoRoundTrip(t *testing.T) {
	req := driver.AudioSpeechRequest{
		Prompt: "hello", Voice: "myvoice", Language: "de", Speed: 1.1, Pitch: 2,
		Emotion: "happy", Style: "conversational", Energy: 0.8, SampleRate: 44100,
	}
	out := driver.AudioSpeechRequestFromProto(driver.AudioSpeechRequestToProto(req))
	if out.Voice != req.Voice || out.Language != req.Language || out.Speed != req.Speed || out.Pitch != req.Pitch ||
		out.Emotion != req.Emotion || out.Style != req.Style || out.Energy != req.Energy {
		t.Fatalf("round trip mismatch: %+v", out)
	}
}

func TestAudioRequestSpeechEnergyValidation(t *testing.T) {
	err := (driver.AudioRequest{
		Task: driver.AudioTaskTextToSpeech, Prompt: "hi", Energy: 1.5,
	}).Validate()
	if err == nil {
		t.Fatal("expected energy validation error")
	}
}

func TestParseVoiceTask(t *testing.T) {
	got, err := driver.ParseVoiceTask("conversion")
	if err != nil {
		t.Fatalf("ParseVoiceTask: %v", err)
	}
	if got != driver.VoiceTaskConversion {
		t.Fatalf("got %q want conversion", got)
	}
}

func TestVoiceRequestValidate(t *testing.T) {
	err := (driver.VoiceRequest{
		Task: driver.VoiceTaskConversion, Name: "demo", SamplePath: "ref.wav",
	}).Validate()
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestAudioMusicProtoRoundTrip(t *testing.T) {
	req := driver.AudioMusicRequest{
		AudioCommonParams: driver.AudioCommonParams{Prompt: "beat", Duration: 30, SampleRate: 44100},
	}
	out := driver.AudioMusicRequestFromProto(driver.AudioMusicRequestToProto(req))
	if out.Prompt != req.Prompt {
		t.Fatalf("prompt: got %q want %q", out.Prompt, req.Prompt)
	}
}
