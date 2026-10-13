package grpc

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
)

func TestRemoteDriverTextMediaInterfaces(t *testing.T) {
	rd := &RemoteDriver{info: driver.Info{ID: "test"}}
	modes := driver.TextInputModesFor(rd)
	want := []driver.TextInputMode{
		driver.TextInputPrompt,
		driver.TextInputImage,
		driver.TextInputVideo,
		driver.TextInputAudio,
		driver.TextInputDocument,
	}
	if len(modes) != len(want) {
		t.Fatalf("modes=%v want %v", modes, want)
	}
	for i, mode := range want {
		if modes[i] != mode {
			t.Fatalf("modes[%d]=%q want %q", i, modes[i], mode)
		}
	}
}
