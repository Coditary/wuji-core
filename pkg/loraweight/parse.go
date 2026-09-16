package loraweight

import (
	"fmt"
	"strconv"
	"strings"
)

// Parse normalizes a LoRA strength to a multiplier (0.7, 1, 70%, 100%).
// Percentages use a trailing "%" (same convention as wuji --scale).
func Parse(raw string) (float32, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, fmt.Errorf("lora weight is empty")
	}
	if strings.Contains(raw, ",") {
		return 0, fmt.Errorf("use '.' as decimal separator, not ',' (got %q)", raw)
	}

	if strings.HasSuffix(raw, "%") {
		num := strings.TrimSpace(strings.TrimSuffix(raw, "%"))
		v, err := strconv.ParseFloat(num, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid lora weight percentage %q: %w", raw, err)
		}
		if v <= 0 {
			return 0, fmt.Errorf("lora weight percentage must be > 0 (got %g%%)", v)
		}
		return float32(v / 100), nil
	}

	v, err := strconv.ParseFloat(raw, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid lora weight %q: use 0.7, 1, or 70%%", raw)
	}
	if v <= 0 {
		return 0, fmt.Errorf("lora weight must be > 0 (got %g)", v)
	}
	return float32(v), nil
}
