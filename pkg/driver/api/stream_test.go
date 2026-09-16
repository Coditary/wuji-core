package api

import (
	"strings"
	"testing"
)

func TestStreamSSEOpenAI(t *testing.T) {
	raw := strings.Join([]string{
		"data: {\"choices\":[{\"delta\":{\"content\":\"Hello\"}}]}",
		"data: {\"choices\":[{\"delta\":{\"content\":\" world\"}}]}",
		"data: [DONE]",
	}, "\n")
	var got strings.Builder
	if err := streamSSE(strings.NewReader(raw), func(delta string) error {
		got.WriteString(delta)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if got.String() != "Hello world" {
		t.Fatalf("got %q", got.String())
	}
}

func TestReadSSEOpenAI(t *testing.T) {
	raw := strings.Join([]string{
		"data: {\"choices\":[{\"delta\":{\"content\":\"Hello\"}}]}",
		"data: {\"choices\":[{\"delta\":{\"content\":\" world\"}}]}",
		"data: [DONE]",
	}, "\n")
	out, err := readSSE(strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "Hello world") {
		t.Fatalf("got %s", string(out))
	}
}

func TestReadSSEAnthropic(t *testing.T) {
	raw := "data: {\"type\":\"content_block_delta\",\"delta\":{\"text\":\"Hi\"}}\n"
	out, err := readSSE(strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "Hi") {
		t.Fatalf("got %s", string(out))
	}
}

func TestStreamSSEOpenAIThinking(t *testing.T) {
	raw := strings.Join([]string{
		"data: {\"choices\":[{\"delta\":{\"reasoning\":\"Let me think\"}}]}",
		"data: {\"choices\":[{\"delta\":{\"content\":\"Answer\"}}]}",
		"data: [DONE]",
	}, "\n")
	var thinking, content strings.Builder
	if err := streamSSEWithCallbacks(strings.NewReader(raw), streamCallbacks{
		OnThink: func(delta string) error {
			thinking.WriteString(delta)
			return nil
		},
		OnContent: func(delta string) error {
			content.WriteString(delta)
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}
	if thinking.String() != "Let me think" {
		t.Fatalf("thinking = %q", thinking.String())
	}
	if content.String() != "Answer" {
		t.Fatalf("content = %q", content.String())
	}
}
