package driver

import "testing"

func TestClassifyScale(t *testing.T) {
	tests := []struct {
		factor float32
		want   ScaleKind
	}{
		{0.5, ScaleKindDownscale},
		{1, ScaleKindUpscaleExact},
		{2, ScaleKindUpscaleExact},
		{2.5, ScaleKindHybrid},
		{4, ScaleKindUpscaleExact},
		{6, ScaleKindHybrid},
	}
	for _, tc := range tests {
		if got := ClassifyScale(tc.factor); got != tc.want {
			t.Fatalf("ClassifyScale(%g) = %v, want %v", tc.factor, got, tc.want)
		}
	}
}
