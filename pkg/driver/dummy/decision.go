package dummy

import (
	"context"
	"hash/fnv"
	"strings"

	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *Driver) Decide(_ context.Context, req driver.DecisionRequest) (*driver.DecisionResponse, error) {
	if err := req.ResolveQuestions(); err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	answers := make([]driver.DecisionAnswer, 0, len(req.Questions))
	for _, q := range req.Questions {
		answers = append(answers, dummyAnswer(q, req.ResolvedState()))
	}
	return &driver.DecisionResponse{
		Driver:  DriverID,
		Model:   req.Model,
		Answers: answers,
	}, nil
}

func dummyAnswer(q driver.DecisionQuestion, state map[string]string) driver.DecisionAnswer {
	seed := q.ID + "|" + stateText(state)
	h := fnv.New32a()
	_, _ = h.Write([]byte(seed))
	n := h.Sum32()

	ans := driver.DecisionAnswer{
		QuestionID: q.ID,
		Type:       q.Type,
		Confidence: 0.5 + float64(n%500)/1000,
	}
	switch q.Type {
	case driver.DecisionTaskChoice:
		opts := choiceOptions(q)
		idx := int(n % uint32(len(opts)))
		ans.Choice = opts[idx]
		ans.Probabilities = uniformProbs(opts, idx)
	case driver.DecisionTaskScore:
		levels := scoreLevels(q)
		idx := int(n % uint32(len(levels)))
		ans.Score = float64(idx)
		ans.Probabilities = uniformProbs(levels, idx)
	case driver.DecisionTaskNoul:
		ans.Noul = float64(n%1000) / 1000
	}
	return ans
}

func stateText(state map[string]string) string {
	if t := strings.TrimSpace(state["text"]); t != "" {
		return t
	}
	parts := make([]string, 0, len(state))
	for k, v := range state {
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, ";")
}

func choiceOptions(q driver.DecisionQuestion) []string {
	if len(q.Criteria) > 0 {
		return q.Criteria
	}
	out := make([]string, 0, len(q.Options))
	for label := range q.Options {
		out = append(out, label)
	}
	return out
}

func scoreLevels(q driver.DecisionQuestion) []string {
	if len(q.Criteria) > 0 {
		return q.Criteria
	}
	out := make([]string, 0, len(q.Options))
	for label := range q.Options {
		out = append(out, label)
	}
	return out
}

func uniformProbs(labels []string, peak int) map[string]float64 {
	if len(labels) == 0 {
		return nil
	}
	out := make(map[string]float64, len(labels))
	rest := 1.0
	if peak >= 0 && peak < len(labels) {
		out[labels[peak]] = 0.7
		rest = 0.3
	}
	share := rest / float64(len(labels))
	for _, label := range labels {
		if out[label] == 0 {
			out[label] = share
		}
	}
	return out
}
