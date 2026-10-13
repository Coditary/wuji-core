package media

import "testing"

func TestRIFEPassesNeeded(t *testing.T) {
	tests := []struct {
		multiplier float64
		want       int
	}{
		{1.0, 0},
		{1.5, 1},
		{2.0, 1},
		{3.0, 2},
		{4.0, 2},
	}
	for _, tc := range tests {
		if got := RIFEPassesNeeded(tc.multiplier); got != tc.want {
			t.Fatalf("multiplier %g: got %d want %d", tc.multiplier, got, tc.want)
		}
	}
}
