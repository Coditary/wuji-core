package ragstore

import "strings"

// ChunksWithoutSources removes chunks whose source is in drop set.
func ChunksWithoutSources(chunks []Chunk, drop []string) []Chunk {
	if len(drop) == 0 {
		out := make([]Chunk, len(chunks))
		copy(out, chunks)
		return out
	}
	removed := map[string]struct{}{}
	for _, s := range drop {
		removed[s] = struct{}{}
	}
	out := make([]Chunk, 0, len(chunks))
	for _, c := range chunks {
		if _, ok := removed[c.Source]; ok {
			continue
		}
		out = append(out, c)
	}
	return out
}

// UniqueSources returns distinct source paths from chunks.
func UniqueSources(chunks []Chunk) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, c := range chunks {
		if _, ok := seen[c.Source]; ok {
			continue
		}
		seen[c.Source] = struct{}{}
		out = append(out, c.Source)
	}
	return out
}

// MergeSources combines kept chunk sources with newly indexed sources.
func MergeSources(kept []Chunk, incoming []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, c := range kept {
		if _, ok := seen[c.Source]; ok {
			continue
		}
		seen[c.Source] = struct{}{}
		out = append(out, c.Source)
	}
	for _, s := range incoming {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

// SplitTextSmart splits on markdown sections then fixed-size windows.
func SplitTextSmart(text string, size, overlap int) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	sections := splitSections(text)
	var out []string
	for _, sec := range sections {
		sec = strings.TrimSpace(sec)
		if sec == "" {
			continue
		}
		parts := SplitText(sec, size, overlap)
		out = append(out, parts...)
	}
	if len(out) == 0 {
		return SplitText(text, size, overlap)
	}
	return out
}

func splitSections(text string) []string {
	lines := strings.Split(text, "\n")
	var sections []string
	var b strings.Builder
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "#") && b.Len() > 0 {
			sections = append(sections, b.String())
			b.Reset()
		}
		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(line)
	}
	if b.Len() > 0 {
		sections = append(sections, b.String())
	}
	if len(sections) <= 1 {
		return strings.Split(text, "\n\n")
	}
	return sections
}
