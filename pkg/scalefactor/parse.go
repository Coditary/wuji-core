package scalefactor

import (
	"fmt"
	"strconv"
	"strings"
)

// Parse normalizes user scale input to a multiplier (0.5 = half, 2 = double, 2.5 = 250%).
// Accepts whole numbers, decimals with ".", and percentages with a trailing "%".
func Parse(raw string) (float32, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, fmt.Errorf("scale value is empty")
	}
	if strings.Contains(raw, ",") {
		return 0, fmt.Errorf("use '.' as decimal separator, not ',' (got %q)", raw)
	}

	if strings.HasSuffix(raw, "%") {
		num := strings.TrimSpace(strings.TrimSuffix(raw, "%"))
		v, err := strconv.ParseFloat(num, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid scale percentage %q: %w", raw, err)
		}
		if v <= 0 {
			return 0, fmt.Errorf("scale percentage must be > 0 (got %g%%)", v)
		}
		return float32(v / 100), nil
	}

	v, err := strconv.ParseFloat(raw, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid scale factor %q: use 2, 2.5, 50%%, or 250%%", raw)
	}
	if v <= 0 {
		return 0, fmt.Errorf("scale factor must be > 0 (got %g)", v)
	}
	return float32(v), nil
}
