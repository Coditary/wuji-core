package loraweight

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		in   string
		want float32
	}{
		{"0.7", 0.7},
		{"1", 1},
		{"70%", 0.7},
		{"100%", 1},
	}
	for _, tc := range tests {
		got, err := Parse(tc.in)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("Parse(%q) = %g, want %g", tc.in, got, tc.want)
		}
	}
}
