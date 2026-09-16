package scalefactor

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		in   string
		want float32
	}{
		{"2", 2},
		{"2.5", 2.5},
		{"50%", 0.5},
		{"250%", 2.5},
		{"200%", 2},
		{"100%", 1},
		{"300%", 3},
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

func TestParseRejectsInvalid(t *testing.T) {
	for _, in := range []string{"", "-1", "0", "0%", "-50%", "2,5", "250,5%"} {
		if _, err := Parse(in); err == nil {
			t.Fatalf("Parse(%q) expected error", in)
		}
	}
}
