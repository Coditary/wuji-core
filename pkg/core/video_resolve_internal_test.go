package core

import "testing"

func TestAIScaleDriverOrder(t *testing.T) {
	cases := []struct {
		scale float32
		want  []string
	}{
		{1, []string{waifu2xDriverID, chainnerDriverID}},
		{1.5, []string{waifu2xDriverID, chainnerDriverID}},
		{2, []string{realesrganDriverID, waifu2xDriverID, chainnerDriverID}},
		{3, []string{realesrganDriverID, chainnerDriverID}},
		{4, []string{realesrganDriverID, chainnerDriverID}},
		{0.5, nil},
		{2.5, nil},
	}
	for _, tc := range cases {
		got := aiScaleDriverOrder(tc.scale)
		if len(got) != len(tc.want) {
			t.Fatalf("scale %g: got %v want %v", tc.scale, got, tc.want)
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Fatalf("scale %g: got %v want %v", tc.scale, got, tc.want)
			}
		}
	}
}
