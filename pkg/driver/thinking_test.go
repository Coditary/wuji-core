package driver

import "testing"

func TestSingleThinkingBlock(t *testing.T) {
	if got := SingleThinkingBlock("", "x"); got != nil {
		t.Fatalf("empty text = %#v", got)
	}
	blocks := SingleThinkingBlock("  hello  ", "openai_reasoning")
	if len(blocks) != 1 || blocks[0].Text != "hello" || blocks[0].Index != 0 {
		t.Fatalf("blocks = %#v", blocks)
	}
}

func TestJoinThinkingBlocks(t *testing.T) {
	joined := JoinThinkingBlocks([]ThinkingBlock{
		{Index: 0, Text: "a"},
		{Index: 1, Text: "b"},
	})
	if joined != "a\n\nb" {
		t.Fatalf("joined = %q", joined)
	}
}
