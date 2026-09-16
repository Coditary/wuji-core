package driver

import (
	"encoding/json"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
)

// TextRequestFromProto converts a protobuf generate-text request.
func TextRequestFromProto(req *wujiv1.GenerateTextRequest) TextRequest {
	if req == nil {
		return TextRequest{}
	}

	out := TextRequest{
		InputMode:         TextInputMode(req.GetInputMode()),
		Prompt:            req.GetPrompt(),
		MediaPath:         req.GetMediaPath(),
		Translate:         req.GetTranslate(),
		TargetLang:        req.GetTargetLang(),
		Language:          req.GetLanguage(),
		BeamSize:          int(req.GetBeamSize()),
		WordTimestamps:    req.GetWordTimestamps(),
		VADEnabled:        req.GetVadEnabled(),
		VADThreshold:      req.GetVadThreshold(),
		Model:             req.GetModel(),
		SystemPrompt:      req.GetSystemPrompt(),
		MaxTokens:         int(req.GetMaxTokens()),
		Temperature:       req.GetTemperature(),
		TopP:              req.GetTopP(),
		TopK:              int(req.GetTopK()),
		MinP:              req.GetMinP(),
		FrequencyPenalty:  req.GetFrequencyPenalty(),
		PresencePenalty:   req.GetPresencePenalty(),
		RepetitionPenalty: req.GetRepetitionPenalty(),
		StopSequences:     append([]string(nil), req.GetStopSequences()...),
		ContextWindow:     int(req.GetContextWindow()),
		LoRAs:             lorasFromProto(req.GetLoras()),
		Messages:          chatMessagesFromProto(req.GetMessages()),
		Tools:             chatToolsFromProto(req.GetTools()),
		Stream:            req.GetStream(),
		ToolChoice:        req.GetToolChoice(),
	}
	if req.Seed != nil {
		seed := int(req.GetSeed())
		out.Seed = &seed
	}
	return out
}

// TextRequestToProto converts a TextRequest to protobuf.
func TextRequestToProto(req TextRequest) *wujiv1.GenerateTextRequest {
	out := &wujiv1.GenerateTextRequest{
		InputMode:         string(req.InputMode),
		Prompt:            req.Prompt,
		MediaPath:         req.MediaPath,
		Translate:         req.Translate,
		TargetLang:        req.TargetLang,
		Language:          req.Language,
		BeamSize:          int32(req.BeamSize),
		WordTimestamps:    req.WordTimestamps,
		VadEnabled:        req.VADEnabled,
		VadThreshold:      req.VADThreshold,
		Model:             req.Model,
		SystemPrompt:      req.SystemPrompt,
		MaxTokens:         int32(req.MaxTokens),
		Temperature:       req.Temperature,
		TopP:              req.TopP,
		TopK:              int32(req.TopK),
		MinP:              req.MinP,
		FrequencyPenalty:  req.FrequencyPenalty,
		PresencePenalty:   req.PresencePenalty,
		RepetitionPenalty: req.RepetitionPenalty,
		StopSequences:     append([]string(nil), req.StopSequences...),
		ContextWindow:     int32(req.ContextWindow),
		Loras:             lorasToProto(req.LoRAs),
		Messages:          chatMessagesToProto(req.Messages),
		Tools:             chatToolsToProto(req.Tools),
		Stream:            req.Stream,
		ToolChoice:        req.ToolChoice,
	}
	if req.Seed != nil {
		seed := int32(*req.Seed)
		out.Seed = &seed
	}
	return out
}

func chatMessagesToProto(msgs []ChatMessage) []*wujiv1.ChatMessage {
	if len(msgs) == 0 {
		return nil
	}
	out := make([]*wujiv1.ChatMessage, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, &wujiv1.ChatMessage{
			Role:       string(m.Role),
			Content:    m.Content,
			Name:       m.Name,
			ToolCallId: m.ToolCallID,
			ToolCalls:  chatToolCallsToProto(m.ToolCalls),
		})
	}
	return out
}

func chatMessagesFromProto(msgs []*wujiv1.ChatMessage) []ChatMessage {
	if len(msgs) == 0 {
		return nil
	}
	out := make([]ChatMessage, 0, len(msgs))
	for _, m := range msgs {
		if m == nil {
			continue
		}
		out = append(out, ChatMessage{
			Role:       ChatRole(m.GetRole()),
			Content:    m.GetContent(),
			Name:       m.GetName(),
			ToolCallID: m.GetToolCallId(),
			ToolCalls:  chatToolCallsFromProto(m.GetToolCalls()),
		})
	}
	return out
}

func chatToolsToProto(tools []ChatToolDef) []*wujiv1.ChatToolDef {
	if len(tools) == 0 {
		return nil
	}
	out := make([]*wujiv1.ChatToolDef, 0, len(tools))
	for _, t := range tools {
		paramsJSON := ""
		if len(t.Parameters) > 0 {
			raw, err := json.Marshal(t.Parameters)
			if err == nil {
				paramsJSON = string(raw)
			}
		}
		out = append(out, &wujiv1.ChatToolDef{
			Name:            t.Name,
			Description:     t.Description,
			ParametersJson: paramsJSON,
		})
	}
	return out
}

func chatToolCallsToProto(calls []ChatToolCall) []*wujiv1.ChatToolCall {
	if len(calls) == 0 {
		return nil
	}
	out := make([]*wujiv1.ChatToolCall, 0, len(calls))
	for _, c := range calls {
		out = append(out, &wujiv1.ChatToolCall{
			Id:            c.ID,
			Name:          c.Name,
			ArgumentsJson: c.Args,
		})
	}
	return out
}

func chatToolCallsFromProto(calls []*wujiv1.ChatToolCall) []ChatToolCall {
	if len(calls) == 0 {
		return nil
	}
	out := make([]ChatToolCall, 0, len(calls))
	for _, c := range calls {
		if c == nil {
			continue
		}
		out = append(out, ChatToolCall{
			ID:   c.GetId(),
			Name: c.GetName(),
			Args: c.GetArgumentsJson(),
		})
	}
	return out
}

// TextResponseToProto converts a TextResponse to protobuf.
func TextResponseToProto(resp TextResponse) *wujiv1.GenerateTextResponse {
	return &wujiv1.GenerateTextResponse{
		Text:         resp.Text,
		TokensUsed:   int32(resp.TokensUsed),
		FinishReason: resp.FinishReason,
		ToolCalls:    chatToolCallsToProto(resp.ToolCalls),
	}
}

// TextResponseFromProto converts a protobuf generate-text response.
func TextResponseFromProto(resp *wujiv1.GenerateTextResponse) TextResponse {
	if resp == nil {
		return TextResponse{}
	}
	return TextResponse{
		Text:         resp.GetText(),
		TokensUsed:   int(resp.GetTokensUsed()),
		FinishReason: resp.GetFinishReason(),
		ToolCalls:    chatToolCallsFromProto(resp.GetToolCalls()),
	}
}

func chatToolsFromProto(tools []*wujiv1.ChatToolDef) []ChatToolDef {
	if len(tools) == 0 {
		return nil
	}
	out := make([]ChatToolDef, 0, len(tools))
	for _, t := range tools {
		if t == nil {
			continue
		}
		var params map[string]any
		if raw := t.GetParametersJson(); raw != "" {
			_ = json.Unmarshal([]byte(raw), &params)
		}
		out = append(out, ChatToolDef{
			Name:        t.GetName(),
			Description: t.GetDescription(),
			Parameters:  params,
		})
	}
	return out
}
