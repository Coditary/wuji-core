package imageresize

import "testing"

func TestFormatResizeArg(t *testing.T) {
	tests := []struct {
		factor float32
		want   string
	}{
		{0.5, "50%"},
		{1, "100%"},
		{2, "200%"},
		{2.5, "250%"},
	}
	for _, tc := range tests {
		got := formatResizeArg(tc.factor)
		if got != tc.want {
			t.Fatalf("formatResizeArg(%g) = %q, want %q", tc.factor, got, tc.want)
		}
	}
}
