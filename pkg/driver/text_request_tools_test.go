package driver

import (
	"encoding/json"
	"testing"
)

func TestTextRequestUnmarshalToolsAndToolChoice(t *testing.T) {
	raw := `{
		"driver_id": "llama",
		"request": {
			"messages": [{"role":"user","content":"list files"}],
			"tools": [{"name":"list_dir","description":"list","parameters":{"type":"object"}}],
			"tool_choice": "required",
			"context_window": 32768
		}
	}`
	var body struct {
		DriverID string      `json:"driver_id"`
		Request  TextRequest `json:"request"`
	}
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Request.Tools) != 1 || body.Request.Tools[0].Name != "list_dir" {
		t.Fatalf("tools = %#v", body.Request.Tools)
	}
	if body.Request.ToolChoice != "required" {
		t.Fatalf("tool_choice = %q", body.Request.ToolChoice)
	}
	if body.Request.ContextWindow != 32768 {
		t.Fatalf("context_window = %d", body.Request.ContextWindow)
	}
}
