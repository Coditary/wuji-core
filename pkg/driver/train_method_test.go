package driver

import "testing"

func TestParseImageTrainMethod(t *testing.T) {
	got, err := ParseImageTrainMethod("lora")
	if err != nil {
		t.Fatal(err)
	}
	if got != ImageTrainMethodLoRA {
		t.Fatalf("got %q", got)
	}

	got, err = ParseImageTrainMethod("")
	if err != nil || got != ImageTrainMethodFull {
		t.Fatalf("empty: got %q err %v", got, err)
	}

	_, err = ParseImageTrainMethod("unknown")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseAudioTrainMethodToTask(t *testing.T) {
	method, err := ParseAudioTrainMethod("speech")
	if err != nil {
		t.Fatal(err)
	}
	if method.ToAudioTask() != AudioTaskTextToSpeech {
		t.Fatalf("got %s", method.ToAudioTask())
	}
}

func TestParseTextTrainMethodQLoRA(t *testing.T) {
	got, err := ParseTextTrainMethod("qlora")
	if err != nil || got != TextTrainMethodQLoRA {
		t.Fatalf("got %q err %v", got, err)
	}
}

func TestParseVoiceTrainMethodDefault(t *testing.T) {
	got, err := ParseVoiceTrainMethod("")
	if err != nil || got != VoiceTrainMethodRVC {
		t.Fatalf("got %q err %v", got, err)
	}
}
