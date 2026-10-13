package driver

import (
	"encoding/json"
	"fmt"
	"strings"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
)

func DecisionTasksToProto(tasks []DecisionTask) []string {
	if len(tasks) == 0 {
		return nil
	}
	out := make([]string, 0, len(tasks))
	for _, task := range tasks {
		out = append(out, string(task))
	}
	return out
}

func DecisionTasksFromProto(items []string) []DecisionTask {
	if len(items) == 0 {
		return nil
	}
	out := make([]DecisionTask, 0, len(items))
	for _, item := range items {
		task, err := ParseDecisionTask(item)
		if err != nil {
			continue
		}
		out = append(out, task)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func ProduceDecisionRequestToProto(req DecisionRequest) (*wujiv1.ProduceDecisionRequest, error) {
	stateJSON, err := json.Marshal(req.ResolvedState())
	if err != nil {
		return nil, fmt.Errorf("encode decision state: %w", err)
	}
	questionsJSON, err := json.Marshal(decisionQuestionsToWire(req.Questions))
	if err != nil {
		return nil, fmt.Errorf("encode decision questions: %w", err)
	}
	return &wujiv1.ProduceDecisionRequest{
		Model:         req.Model,
		Text:          req.Text,
		ImagePath:     req.ImagePath,
		ImagePaths:    append([]string(nil), req.ImagePaths...),
		QuestionsFile: req.QuestionsFile,
		UseStdin:      req.UseStdin,
		StateJson:     string(stateJSON),
		QuestionsJson: string(questionsJSON),
	}, nil
}

func ProduceDecisionRequestFromProto(p *wujiv1.ProduceDecisionRequest) (DecisionRequest, error) {
	if p == nil {
		return DecisionRequest{}, nil
	}
	req := DecisionRequest{
		Model:         p.GetModel(),
		Text:          p.GetText(),
		ImagePath:     p.GetImagePath(),
		ImagePaths:    append([]string(nil), p.GetImagePaths()...),
		QuestionsFile: p.GetQuestionsFile(),
		UseStdin:      p.GetUseStdin(),
	}
	if state := strings.TrimSpace(p.GetStateJson()); state != "" {
		if err := json.Unmarshal([]byte(state), &req.State); err != nil {
			return DecisionRequest{}, fmt.Errorf("decode decision state: %w", err)
		}
	}
	if questions := strings.TrimSpace(p.GetQuestionsJson()); questions != "" {
		parsed, state, err := ParseQuestionsJSON([]byte(questions))
		if err != nil {
			return DecisionRequest{}, err
		}
		req.Questions = parsed
		if req.State == nil {
			req.State = state
		} else {
			for k, v := range state {
				if req.State[k] == "" {
					req.State[k] = v
				}
			}
		}
	}
	return req, nil
}

func ProduceDecisionResponseToProto(resp *DecisionResponse) (*wujiv1.ProduceDecisionResponse, error) {
	if resp == nil {
		return &wujiv1.ProduceDecisionResponse{}, nil
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("encode decision response: %w", err)
	}
	return &wujiv1.ProduceDecisionResponse{PayloadJson: string(payload)}, nil
}

func ProduceDecisionResponseFromProto(p *wujiv1.ProduceDecisionResponse) (*DecisionResponse, error) {
	if p == nil || strings.TrimSpace(p.GetPayloadJson()) == "" {
		return nil, fmt.Errorf("empty decision response")
	}
	var resp DecisionResponse
	if err := json.Unmarshal([]byte(p.GetPayloadJson()), &resp); err != nil {
		return nil, fmt.Errorf("decode decision response: %w", err)
	}
	return &resp, nil
}

func decisionQuestionsToWire(questions []DecisionQuestion) map[string]any {
	out := make(map[string]any, len(questions))
	for _, q := range questions {
		wire := map[string]any{
			"type":         q.Type,
			"instructions": q.Instructions,
		}
		if len(q.Options) > 0 {
			wire["criteria"] = q.Options
		} else if len(q.Criteria) > 0 {
			wire["criteria"] = q.Criteria
		}
		out[q.ID] = wire
	}
	return out
}
