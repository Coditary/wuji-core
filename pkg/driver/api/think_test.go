package api

import "testing"

func TestApplyThinkToOpenAIBody(t *testing.T) {
	t.Run("true maps to reasoning_effort low", func(t *testing.T) {
		body := map[string]any{}
		applyThinkToOpenAIBody(body, true)
		if body["reasoning_effort"] != "low" {
			t.Fatalf("reasoning_effort = %#v", body["reasoning_effort"])
		}
		if _, ok := body["think"]; ok {
			t.Fatalf("unexpected think field: %#v", body["think"])
		}
	})

	t.Run("level string", func(t *testing.T) {
		body := map[string]any{}
		applyThinkToOpenAIBody(body, "high")
		if body["reasoning_effort"] != "high" {
			t.Fatalf("reasoning_effort = %#v", body["reasoning_effort"])
		}
	})

	t.Run("false disables thinking", func(t *testing.T) {
		body := map[string]any{}
		applyThinkToOpenAIBody(body, false)
		if len(body) != 0 {
			t.Fatalf("body = %#v", body)
		}
	})

	t.Run("none disables thinking", func(t *testing.T) {
		body := map[string]any{}
		applyThinkToOpenAIBody(body, "none")
		if len(body) != 0 {
			t.Fatalf("body = %#v", body)
		}
	})
}
