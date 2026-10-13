package driver

import (
	"fmt"
	"strings"
)

// DecisionTask identifies a typed decision question format.
type DecisionTask string

const (
	DecisionTaskChoice DecisionTask = "choice"
	DecisionTaskScore  DecisionTask = "score"
	DecisionTaskNoul   DecisionTask = "noul"
)

// AllDecisionTasks lists supported decision task names for driver metadata.
func AllDecisionTasks() []DecisionTask {
	return []DecisionTask{
		DecisionTaskChoice,
		DecisionTaskScore,
		DecisionTaskNoul,
	}
}

// DecisionTaskInfo describes one decision task for help text and validation.
type DecisionTaskInfo struct {
	Task        DecisionTask
	Description string
}

// AllDecisionTaskInfo returns metadata for every decision task.
func AllDecisionTaskInfo() []DecisionTaskInfo {
	return []DecisionTaskInfo{
		{Task: DecisionTaskChoice, Description: "Pick one option with calibrated probabilities"},
		{Task: DecisionTaskScore, Description: "Ordinal score on a rubric (expected level + per-level probs)"},
		{Task: DecisionTaskNoul, Description: "Calibrated probability that a statement is true"},
	}
}

// ParseDecisionTask normalizes a task name.
func ParseDecisionTask(name string) (DecisionTask, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return "", fmt.Errorf("decision task is required (choice, score, noul)")
	}
	for _, t := range AllDecisionTasks() {
		if string(t) == name {
			return t, nil
		}
	}
	return "", fmt.Errorf("unknown decision task %q", name)
}
