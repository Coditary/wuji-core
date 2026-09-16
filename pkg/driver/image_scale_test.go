package driver

import (
	"math"
	"testing"
)

func floatNear(a, b float32) bool {
	return math.Abs(float64(a-b)) < 1e-4
}

func TestPlanScale(t *testing.T) {
	tests := []struct {
		factor     float32
		wantAI     int
		wantMagick float32
	}{
		{0.5, 0, 0.5},
		{1, 0, 1},
		{1.5, 2, 0.75},
		{2, 2, 1},
		{2.5, 4, 0.625},
		{3, 4, 0.75},
		{4, 4, 1},
		{6, 8, 0.75},
		{8, 8, 1},
		{10, 8, 1.25},
	}

	for _, tc := range tests {
		plan := planScale(tc.factor)
		if plan.aiStep != tc.wantAI {
			t.Fatalf("planScale(%g).aiStep = %d, want %d", tc.factor, plan.aiStep, tc.wantAI)
		}
		if !floatNear(plan.magickFactor, tc.wantMagick) {
			t.Fatalf("planScale(%g).magickFactor = %g, want %g", tc.factor, plan.magickFactor, tc.wantMagick)
		}
	}
}
