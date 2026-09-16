package driver_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
)

func TestPlanAudioSegments(t *testing.T) {
	n, err := driver.PlanAudioSegments(10, 2, 10)
	if err != nil || n != 1 {
		t.Fatalf("short track: n=%d err=%v", n, err)
	}

	n, err = driver.PlanAudioSegments(25, 3, 10)
	if err != nil || n != 4 {
		t.Fatalf("25s with 3s overlap: n=%d want 4, err=%v", n, err)
	}

	_, err = driver.PlanAudioSegments(30, 10, 10)
	if err == nil {
		t.Fatal("expected error when overlap equals segment length")
	}
}

func TestAudioRequestOverlapValidation(t *testing.T) {
	err := (driver.AudioRequest{
		Task: driver.AudioTaskTextToSFX, Overlap: 2, Duration: 10,
	}).Validate()
	if err == nil {
		t.Fatal("expected overlap error for sfx task")
	}

	err = (driver.AudioRequest{
		Task: driver.AudioTaskTextToMusic, Prompt: "beat", Overlap: 2, Duration: 25,
	}).Validate()
	if err != nil {
		t.Fatalf("valid overlap request: %v", err)
	}
}
