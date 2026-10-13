package driver

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// DecisionQuestion is one typed question in a decision request.
type DecisionQuestion struct {
	ID           string            `json:"id"`
	Type         DecisionTask      `json:"type"`
	Instructions string            `json:"instructions"`
	Criteria     []string          `json:"criteria,omitempty"`
	Options      map[string]string `json:"options,omitempty"`
}

// DecisionRequest is input for structured decision operations.
type DecisionRequest struct {
	Model         string
	State         map[string]string
	Questions     []DecisionQuestion
	Text          string
	ImagePath     string
	ImagePaths    []string
	QuestionsFile string
	UseStdin      bool
}

// DecisionAnswer is one calibrated answer from a decision backend.
type DecisionAnswer struct {
	QuestionID    string             `json:"question_id"`
	Type          DecisionTask       `json:"type"`
	Choice        string             `json:"choice,omitempty"`
	Score         float64            `json:"score,omitempty"`
	Noul          float64            `json:"noul,omitempty"`
	Probabilities map[string]float64 `json:"probabilities,omitempty"`
	Confidence    float64            `json:"confidence,omitempty"`
}

// DecisionResponse is the structured output of a decision call.
type DecisionResponse struct {
	Answers []DecisionAnswer `json:"answers"`
	Model   string           `json:"model,omitempty"`
	Driver  string           `json:"driver,omitempty"`
	Routing map[string]string `json:"routing,omitempty"`
}

// Validate checks required fields for a decision request.
func (r DecisionRequest) Validate() error {
	if len(r.Questions) == 0 && strings.TrimSpace(r.QuestionsFile) == "" && !r.UseStdin {
		return fmt.Errorf("decision requires --questions file or structured questions input")
	}
	for _, q := range r.Questions {
		if err := q.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate checks one decision question.
func (q DecisionQuestion) Validate() error {
	if strings.TrimSpace(q.ID) == "" {
		return fmt.Errorf("decision question id is required")
	}
	if _, err := ParseDecisionTask(string(q.Type)); err != nil {
		return err
	}
	if strings.TrimSpace(q.Instructions) == "" {
		return fmt.Errorf("decision question %q requires instructions", q.ID)
	}
	switch q.Type {
	case DecisionTaskChoice:
		if len(q.Criteria) == 0 && len(q.Options) == 0 {
			return fmt.Errorf("decision question %q requires criteria or options", q.ID)
		}
	case DecisionTaskScore:
		if len(q.Criteria) < 2 && len(q.Options) < 2 {
			return fmt.Errorf("score question %q requires at least two rubric levels", q.ID)
		}
	}
	return nil
}

// ResolvedState returns the state map with CLI shortcuts applied.
func (r DecisionRequest) ResolvedState() map[string]string {
	state := make(map[string]string, len(r.State)+3)
	for k, v := range r.State {
		state[k] = v
	}
	if t := strings.TrimSpace(r.Text); t != "" {
		if state["text"] == "" {
			state["text"] = t
		}
	}
	if p := strings.TrimSpace(r.ImagePath); p != "" {
		state["image"] = p
	}
	if len(r.ImagePaths) == 1 && state["image"] == "" {
		state["image"] = strings.TrimSpace(r.ImagePaths[0])
	}
	if len(r.ImagePaths) > 1 {
		state["images"] = strings.Join(r.ImagePaths, ",")
	}
	return state
}

// LoadQuestionsFile reads questions from JSON (object map or array).
func LoadQuestionsFile(path string) ([]DecisionQuestion, map[string]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	return ParseQuestionsJSON(raw)
}

// ParseQuestionsJSON accepts either a full decision payload or a questions-only document.
func ParseQuestionsJSON(raw []byte) ([]DecisionQuestion, map[string]string, error) {
	var envelope struct {
		State     map[string]string          `json:"state"`
		Questions map[string]json.RawMessage `json:"questions"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, nil, fmt.Errorf("parse questions json: %w", err)
	}
	if len(envelope.Questions) == 0 {
		var list []DecisionQuestion
		if err := json.Unmarshal(raw, &list); err == nil && len(list) > 0 {
			return list, nil, nil
		}
		var direct map[string]json.RawMessage
		if err := json.Unmarshal(raw, &direct); err != nil {
			return nil, nil, fmt.Errorf("questions must be a JSON object map or array")
		}
		out := make([]DecisionQuestion, 0, len(direct))
		for id, payload := range direct {
			q, err := parseQuestionJSON(id, payload)
			if err != nil {
				return nil, nil, err
			}
			out = append(out, q)
		}
		return out, nil, nil
	}
	out := make([]DecisionQuestion, 0, len(envelope.Questions))
	for id, payload := range envelope.Questions {
		q, err := parseQuestionJSON(id, payload)
		if err != nil {
			return nil, nil, err
		}
		out = append(out, q)
	}
	return out, envelope.State, nil
}

func parseQuestionJSON(id string, payload json.RawMessage) (DecisionQuestion, error) {
	var wire struct {
		Type         DecisionTask      `json:"type"`
		Instructions string            `json:"instructions"`
		Criteria     json.RawMessage   `json:"criteria"`
		Options      map[string]string `json:"options"`
	}
	if err := json.Unmarshal(payload, &wire); err != nil {
		return DecisionQuestion{}, fmt.Errorf("question %q: %w", id, err)
	}
	q := DecisionQuestion{
		ID:           id,
		Type:         wire.Type,
		Instructions: wire.Instructions,
		Options:      wire.Options,
	}
	if len(wire.Criteria) > 0 {
		if wire.Criteria[0] == '[' {
			if err := json.Unmarshal(wire.Criteria, &q.Criteria); err != nil {
				return DecisionQuestion{}, fmt.Errorf("question %q criteria: %w", id, err)
			}
		} else if wire.Criteria[0] == '{' {
			var opts map[string]string
			if err := json.Unmarshal(wire.Criteria, &opts); err != nil {
				return DecisionQuestion{}, fmt.Errorf("question %q criteria: %w", id, err)
			}
			q.Options = opts
		}
	}
	if err := q.Validate(); err != nil {
		return DecisionQuestion{}, err
	}
	return q, nil
}

// ResolveQuestions loads questions from inline fields or a file/stdin payload.
func (r *DecisionRequest) ResolveQuestions() error {
	if len(r.Questions) > 0 {
		return nil
	}
	var raw []byte
	switch {
	case r.UseStdin:
		data, err := os.ReadFile("/dev/stdin")
		if err != nil {
			return err
		}
		raw = data
	case strings.TrimSpace(r.QuestionsFile) != "":
		data, err := os.ReadFile(r.QuestionsFile)
		if err != nil {
			return err
		}
		raw = data
	default:
		return nil
	}
	questions, state, err := ParseQuestionsJSON(raw)
	if err != nil {
		return err
	}
	r.Questions = questions
	if r.State == nil {
		r.State = state
	} else {
		for k, v := range state {
			if r.State[k] == "" {
				r.State[k] = v
			}
		}
	}
	return nil
}
