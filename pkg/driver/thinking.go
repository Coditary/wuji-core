package driver

import "strings"

// ThinkingBlock is one reasoning segment from a model response.
// Most models emit a single block; interleaved or parsed output may produce several.
type ThinkingBlock struct {
	Index  int    `json:"index"`
	Text   string `json:"text"`
	Source string `json:"source,omitempty"`
}

// SingleThinkingBlock returns a one-element slice when text is non-empty.
func SingleThinkingBlock(text, source string) []ThinkingBlock {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	return []ThinkingBlock{{Index: 0, Text: text, Source: source}}
}

// JoinThinkingBlocks concatenates block text for display or legacy callers.
func JoinThinkingBlocks(blocks []ThinkingBlock) string {
	if len(blocks) == 0 {
		return ""
	}
	parts := make([]string, 0, len(blocks))
	for _, b := range blocks {
		if s := strings.TrimSpace(b.Text); s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, "\n\n")
}
