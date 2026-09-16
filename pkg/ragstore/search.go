package ragstore

import "math"

// ScorePair holds a chunk index and similarity score.
type ScorePair struct {
	Index int
	Score float32
}

// TopKByCosine returns indices of the top matching chunks.
func TopKByCosine(query []float32, matrix []float32, dims, k int) []ScorePair {
	if dims <= 0 || len(matrix) == 0 || len(query) != dims {
		return nil
	}
	count := len(matrix) / dims
	if k <= 0 || k > count {
		k = count
	}
	scores := make([]ScorePair, count)
	for i := 0; i < count; i++ {
		vec := matrix[i*dims : (i+1)*dims]
		scores[i] = ScorePair{Index: i, Score: cosine(query, vec)}
	}
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].Score > scores[i].Score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}
	if k > len(scores) {
		k = len(scores)
	}
	return scores[:k]
}

// SelectTopK returns the best chunk indices after optional MMR diversification.
func SelectTopK(query []float32, matrix []float32, dims, searchK, finalK int, diverse bool) []ScorePair {
	if searchK <= 0 {
		searchK = finalK
	}
	pairs := TopKByCosine(query, matrix, dims, searchK)
	if !diverse || finalK <= 0 || len(pairs) <= finalK {
		if finalK > 0 && len(pairs) > finalK {
			return pairs[:finalK]
		}
		return pairs
	}
	return mmrSelect(query, matrix, dims, pairs, finalK)
}

func mmrSelect(query []float32, matrix []float32, dims int, candidates []ScorePair, k int) []ScorePair {
	const lambda = float32(0.7)
	selected := make([]ScorePair, 0, k)
	remaining := append([]ScorePair(nil), candidates...)
	for len(selected) < k && len(remaining) > 0 {
		bestIdx := -1
		var bestScore float32 = -1e9
		for i, cand := range remaining {
			relevance := cand.Score
			redundancy := float32(0)
			candVec := matrix[cand.Index*dims : (cand.Index+1)*dims]
			for _, sel := range selected {
				selVec := matrix[sel.Index*dims : (sel.Index+1)*dims]
				if sim := cosine(candVec, selVec); sim > redundancy {
					redundancy = sim
				}
			}
			score := lambda*relevance - (1-lambda)*redundancy
			if score > bestScore {
				bestScore = score
				bestIdx = i
			}
		}
		if bestIdx < 0 {
			break
		}
		selected = append(selected, remaining[bestIdx])
		remaining = append(remaining[:bestIdx], remaining[bestIdx+1:]...)
	}
	return selected
}

func cosine(a, b []float32) float32 {
	var dot, na, nb float32
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (float32(math.Sqrt(float64(na))) * float32(math.Sqrt(float64(nb))))
}
