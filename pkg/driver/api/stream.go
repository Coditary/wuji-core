package api

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type streamAccum struct {
	text       strings.Builder
	toolCalls  []streamToolCall
	finish     string
}

type streamToolCall struct {
	id   string
	name string
	args strings.Builder
}

// readSSE accumulates text deltas from Server-Sent Events (OpenAI-style).
func readSSE(r io.Reader) ([]byte, error) {
	var text strings.Builder
	err := streamSSE(r, func(delta string) error {
		text.WriteString(delta)
		return nil
	})
	if err != nil {
		return nil, err
	}
	out := map[string]string{"text": text.String()}
	return json.Marshal(out)
}

// streamSSE parses SSE chunks and invokes onDelta for each text fragment.
func streamSSE(r io.Reader, onDelta func(delta string) error) error {
	return streamSSEWithCallbacks(r, streamCallbacks{OnContent: onDelta})
}

// streamCallbacks handles optional content, thinking, and usage SSE events.
type streamCallbacks struct {
	OnContent func(delta string) error
	OnThink   func(delta string) error
	OnUsage   func(promptTokens, completionTokens int)
}

// streamSSEWithUsage parses SSE chunks, delivers text deltas, and captures usage when present.
func streamSSEWithUsage(r io.Reader, onDelta func(delta string) error, onUsage func(promptTokens, completionTokens int)) error {
	return streamSSEWithCallbacks(r, streamCallbacks{OnContent: onDelta, OnUsage: onUsage})
}

func streamSSEWithCallbacks(r io.Reader, cb streamCallbacks) error {
	if cb.OnContent == nil {
		cb.OnContent = func(string) error { return nil }
	}
	if cb.OnThink == nil {
		cb.OnThink = func(string) error { return nil }
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			break
		}
		if cb.OnUsage != nil {
			if prompt, completion, ok := openAISSEUsage(data); ok {
				cb.OnUsage(prompt, completion)
			}
		}
		if delta, ok := openAISSEThinkDelta(data); ok && delta != "" {
			if err := cb.OnThink(delta); err != nil {
				return err
			}
		}
		if delta, ok := openAISSEDelta(data); ok && delta != "" {
			if err := cb.OnContent(delta); err != nil {
				return err
			}
		}
		if delta, ok := anthropicSSEDelta(data); ok && delta != "" {
			if err := cb.OnContent(delta); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}

func openAISSEUsage(data string) (promptTokens, completionTokens int, ok bool) {
	var parsed struct {
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal([]byte(data), &parsed); err != nil {
		return 0, 0, false
	}
	if parsed.Usage.PromptTokens == 0 && parsed.Usage.CompletionTokens == 0 {
		return 0, 0, false
	}
	return parsed.Usage.PromptTokens, parsed.Usage.CompletionTokens, true
}

func openAISSEDelta(data string) (string, bool) {
	var openai struct {
		Choices []struct {
			Delta struct {
				Content string `json:"content"`
			} `json:"delta"`
		} `json:"choices"`
	}
	if err := json.Unmarshal([]byte(data), &openai); err != nil {
		return "", false
	}
	if len(openai.Choices) == 0 {
		return "", false
	}
	return openai.Choices[0].Delta.Content, true
}

func openAISSEThinkDelta(data string) (string, bool) {
	var openai struct {
		Choices []struct {
			Delta struct {
				ReasoningContent string `json:"reasoning_content"`
				Reasoning        string `json:"reasoning"`
				Thinking         string `json:"thinking"`
			} `json:"delta"`
			Message struct {
				ReasoningContent string `json:"reasoning_content"`
				Reasoning        string `json:"reasoning"`
				Thinking         string `json:"thinking"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal([]byte(data), &openai); err != nil {
		return "", false
	}
	if len(openai.Choices) == 0 {
		return "", false
	}
	d := openai.Choices[0].Delta
	if s := firstNonEmptyStr(d.ReasoningContent, d.Reasoning, d.Thinking); s != "" {
		return s, true
	}
	m := openai.Choices[0].Message
	if s := firstNonEmptyStr(m.ReasoningContent, m.Reasoning, m.Thinking); s != "" {
		return s, true
	}
	return "", false
}

func firstNonEmptyStr(parts ...string) string {
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			return p
		}
	}
	return ""
}

func anthropicSSEDelta(data string) (string, bool) {
	var anthropic struct {
		Type  string `json:"type"`
		Delta struct {
			Text string `json:"text"`
		} `json:"delta"`
	}
	if err := json.Unmarshal([]byte(data), &anthropic); err != nil {
		return "", false
	}
	if anthropic.Type != "content_block_delta" {
		return "", false
	}
	return anthropic.Delta.Text, true
}
