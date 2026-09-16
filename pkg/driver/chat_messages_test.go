package driver

import "testing"

func TestTextRequestChatMessagesFromPrompt(t *testing.T) {
	msgs := (TextRequest{
		SystemPrompt: "be helpful",
		Prompt:       "hello",
	}).ChatMessages()
	if len(msgs) != 2 || msgs[0].Role != ChatRoleSystem || msgs[1].Content != "hello" {
		t.Fatalf("unexpected messages: %+v", msgs)
	}
}

func TestTextRequestChatMessagesPreservesHistory(t *testing.T) {
	msgs := (TextRequest{
		Messages: []ChatMessage{
			{Role: ChatRoleSystem, Content: "sys"},
			{Role: ChatRoleUser, Content: "hi"},
			{Role: ChatRoleAssistant, Content: "hello"},
			{Role: ChatRoleTool, Name: "list_dir", ToolCallID: "1", Content: "a.txt"},
		},
	}).ChatMessages()
	if len(msgs) != 4 {
		t.Fatalf("expected 4 messages, got %d", len(msgs))
	}
}

func TestFlattenChatMessages(t *testing.T) {
	system, prompt := FlattenChatMessages([]ChatMessage{
		{Role: ChatRoleSystem, Content: "sys"},
		{Role: ChatRoleUser, Content: "question"},
		{Role: ChatRoleAssistant, Content: "answer"},
	})
	if system != "sys" {
		t.Fatalf("system: %q", system)
	}
	if prompt != "User: question\nAssistant: answer" {
		t.Fatalf("prompt: %q", prompt)
	}
}
