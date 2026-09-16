package driver_test

import (
	"context"
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/driver/dummy"
)

func TestRunAudioTaskTextToSpeech(t *testing.T) {
	d := dummy.New()
	resp, err := driver.RunAudioTask(context.Background(), d, driver.AudioRequest{
		Task: driver.AudioTaskTextToSpeech, Prompt: "Guten Morgen", Voice: "myvoice",
		Language: "de", Speed: 1.0, Pitch: 0,
	})
	if err != nil {
		t.Fatalf("RunAudioTask: %v", err)
	}
	if resp.Format != "wav" || resp.Path == "" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}
