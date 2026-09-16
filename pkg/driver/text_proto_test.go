package driver

import (
	"testing"
)

func TestTextRequestProtoRoundTripMessages(t *testing.T) {
	orig := TextRequest{
		Model: "test-model",
		Messages: []ChatMessage{
			{Role: ChatRoleSystem, Content: "You are helpful."},
			{Role: ChatRoleUser, Content: "hey"},
		},
		Tools: []ChatToolDef{
			{Name: "search", Description: "search the web", Parameters: map[string]any{"type": "object"}},
		},
		MaxTokens: 100,
	}
	if err := orig.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}

	back := TextRequestFromProto(TextRequestToProto(orig))
	if err := back.Validate(); err != nil {
		t.Fatalf("validate round-trip: %v", err)
	}
	if len(back.Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(back.Messages))
	}
	if back.Messages[1].Content != "hey" {
		t.Fatalf("unexpected user message: %q", back.Messages[1].Content)
	}
	if len(back.Tools) != 1 || back.Tools[0].Name != "search" {
		t.Fatalf("unexpected tools: %+v", back.Tools)
	}
}

func TestTextRequestProtoRoundTripToolChoiceAndToolCalls(t *testing.T) {
	orig := TextRequest{
		Model:      "test-model",
		ToolChoice: "required",
		Messages: []ChatMessage{
			{Role: ChatRoleAssistant, ToolCalls: []ChatToolCall{{ID: "1", Name: "list_dir", Args: `{"path":"."}`}}},
			{Role: ChatRoleTool, Name: "list_dir", ToolCallID: "1", Content: "a.txt"},
		},
	}
	back := TextRequestFromProto(TextRequestToProto(orig))
	if back.ToolChoice != "required" {
		t.Fatalf("tool_choice = %q", back.ToolChoice)
	}
	if len(back.Messages) != 2 || len(back.Messages[0].ToolCalls) != 1 {
		t.Fatalf("unexpected messages: %+v", back.Messages)
	}
	if back.Messages[0].ToolCalls[0].Name != "list_dir" {
		t.Fatalf("tool call = %+v", back.Messages[0].ToolCalls[0])
	}
}

func TestTextResponseProtoRoundTripToolCalls(t *testing.T) {
	orig := TextResponse{
		Text:         "",
		TokensUsed:   12,
		FinishReason: "tool_calls",
		ToolCalls:    []ChatToolCall{{ID: "call_1", Name: "list_dir", Args: `{"path":"."}`}},
	}
	back := TextResponseFromProto(TextResponseToProto(orig))
	if len(back.ToolCalls) != 1 || back.ToolCalls[0].Name != "list_dir" {
		t.Fatalf("tool_calls = %+v", back.ToolCalls)
	}
}
