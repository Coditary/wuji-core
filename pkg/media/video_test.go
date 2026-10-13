package media

import "testing"

func TestParseFrameRate(t *testing.T) {
	tests := []struct {
		raw  string
		want int
	}{
		{"24/1", 24},
		{"30000/1001", 30},
		{"60", 60},
	}
	for _, tc := range tests {
		got, err := parseFrameRate(tc.raw)
		if err != nil {
			t.Fatalf("%q: %v", tc.raw, err)
		}
		if got != tc.want {
			t.Fatalf("%q: got %d want %d", tc.raw, got, tc.want)
		}
	}
}
