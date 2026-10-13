package driver_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
)

func TestProduceDecisionProtoRoundTrip(t *testing.T) {
	req := driver.DecisionRequest{
		Model: "laya-multilingual",
		Text:  "Package arrived broken",
		Questions: []driver.DecisionQuestion{
			{
				ID:           "intent",
				Type:         driver.DecisionTaskChoice,
				Instructions: "What is this about?",
				Criteria:     []string{"billing", "damage", "other"},
			},
		},
	}
	pb, err := driver.ProduceDecisionRequestToProto(req)
	if err != nil {
		t.Fatal(err)
	}
	back, err := driver.ProduceDecisionRequestFromProto(pb)
	if err != nil {
		t.Fatal(err)
	}
	if len(back.Questions) != 1 || back.Questions[0].ID != "intent" {
		t.Fatalf("unexpected questions: %+v", back.Questions)
	}

	resp := &driver.DecisionResponse{
		Model: "laya-multilingual",
		Answers: []driver.DecisionAnswer{
			{QuestionID: "intent", Type: driver.DecisionTaskChoice, Choice: "damage", Confidence: 0.9},
		},
	}
	pbResp, err := driver.ProduceDecisionResponseToProto(resp)
	if err != nil {
		t.Fatal(err)
	}
	restored, err := driver.ProduceDecisionResponseFromProto(pbResp)
	if err != nil {
		t.Fatal(err)
	}
	if restored.Answers[0].Choice != "damage" {
		t.Fatalf("unexpected choice: %+v", restored)
	}
}
