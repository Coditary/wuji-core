package batchplan

import (
	"fmt"
)

// Input describes a batch planning request.
type Input struct {
	Total int // target image count (--batch)

	// Size and Count are CLI values; zero means unset unless Set is true.
	Size     int
	Count    int
	SizeSet  bool
	CountSet bool

	Width, Height   int
	ControlUnits    int
	AvailableVRAMMB int // 0 = unknown, use conservative defaults
}

// Result is a resolved batch_size × batch_count pair.
type Result struct {
	Size  int
	Count int
}

// Resolve computes batch_size and batch_count from --batch and optional overrides.
func Resolve(in Input) (Result, error) {
	if in.Total <= 0 {
		return Result{Size: in.Size, Count: normalizeCount(in.Count)}, nil
	}
	if in.Total > 1000 {
		return Result{}, fmt.Errorf("--batch must be at most 1000 (got %d)", in.Total)
	}

	if in.SizeSet && in.CountSet {
		if in.Size <= 0 {
			return Result{}, fmt.Errorf("--batch-size must be > 0 (got %d)", in.Size)
		}
		if in.Count <= 0 {
			return Result{}, fmt.Errorf("--batch-count must be > 0 (got %d)", in.Count)
		}
		return Result{Size: in.Size, Count: in.Count}, nil
	}

	if in.SizeSet {
		if in.Size <= 0 {
			return Result{}, fmt.Errorf("--batch-size must be > 0 (got %d)", in.Size)
		}
		if in.Total%in.Size != 0 {
			return Result{}, fmt.Errorf("--batch %d is not divisible by --batch-size %d", in.Total, in.Size)
		}
		return Result{Size: in.Size, Count: in.Total / in.Size}, nil
	}

	if in.CountSet {
		if in.Count <= 0 {
			return Result{}, fmt.Errorf("--batch-count must be > 0 (got %d)", in.Count)
		}
		if in.Total%in.Count != 0 {
			return Result{}, fmt.Errorf("--batch %d is not divisible by --batch-count %d", in.Total, in.Count)
		}
		return Result{Size: in.Total / in.Count, Count: in.Count}, nil
	}

	return heuristic(in)
}

func heuristic(in Input) (Result, error) {
	perImage := estimatePerImageMB(in.Width, in.Height, in.ControlUnits)
	maxParallel, tight := maxParallelBatchSize(in, perImage)
	size, count := splitTotal(in.Total, maxParallel, tight)
	if size <= 0 || count <= 0 {
		return Result{}, fmt.Errorf("could not plan batch for %d images", in.Total)
	}
	return Result{Size: size, Count: count}, nil
}

func normalizeCount(count int) int {
	if count <= 0 {
		return 1
	}
	return count
}

func maxParallelBatchSize(in Input, perImage int) (max int, tight bool) {
	if in.AvailableVRAMMB <= 0 {
		return 1, true
	}
	const modelReserveMB = 3500
	spare := in.AvailableVRAMMB - modelReserveMB
	if spare <= 0 {
		return 1, true
	}
	n := spare / perImage
	if n < 1 {
		return 1, true
	}
	if n > 8 {
		n = 8
	}
	used := modelReserveMB + n*perImage
	tight = used >= in.AvailableVRAMMB*85/100
	return n, tight
}

func estimatePerImageMB(width, height, controlUnits int) int {
	if width <= 0 {
		width = 512
	}
	if height <= 0 {
		height = 512
	}
	basePixels := 512 * 512
	pixels := width * height
	perSlot := 500 + (pixels*400)/basePixels
	perSlot += controlUnits * 150
	if perSlot < 400 {
		return 400
	}
	return perSlot
}

// splitTotal picks a power-of-two batch_size divisor of total that fits maxParallel.
// When preferLowSize is true (VRAM estimate is tight), it prefers the smallest
// viable batch_size above 1 instead of the largest one.
func splitTotal(total, maxParallel int, preferLowSize bool) (size, count int) {
	if total <= 0 {
		return 0, 0
	}
	if maxParallel < 1 {
		maxParallel = 1
	}

	var pow2 []int
	for d := 1; d <= total && d <= maxParallel; d++ {
		if total%d == 0 && d&(d-1) == 0 {
			pow2 = append(pow2, d)
		}
	}
	if len(pow2) == 0 {
		return 1, total
	}

	if preferLowSize {
		for _, d := range pow2 {
			if d > 1 {
				return d, total / d
			}
		}
		return 1, total
	}

	best := pow2[0]
	for _, d := range pow2 {
		if d > best {
			best = d
		}
	}
	return best, total / best
}
