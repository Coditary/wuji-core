package ragstore

import "unicode/utf8"

// SplitText splits text into overlapping chunks measured in runes.
func SplitText(text string, size, overlap int) []string {
	if size <= 0 {
		size = 800
	}
	if overlap < 0 {
		overlap = 0
	}
	if overlap >= size {
		overlap = size / 4
	}
	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}
	if len(runes) <= size {
		return []string{text}
	}
	step := size - overlap
	if step <= 0 {
		step = size
	}
	var out []string
	for start := 0; start < len(runes); start += step {
		end := start + size
		if end > len(runes) {
			end = len(runes)
		}
		out = append(out, string(runes[start:end]))
		if end >= len(runes) {
			break
		}
	}
	return out
}

// EstimateRunes returns rune count for logging/stats.
func EstimateRunes(s string) int {
	return utf8.RuneCountInString(s)
}
