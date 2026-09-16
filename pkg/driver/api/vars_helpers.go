package api

import (
	"encoding/json"
	"fmt"
)

func boolStr(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func jsonStr(v any) string {
	if v == nil {
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func ptrIntStr(v *int) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%d", *v)
}

func optionalIntStr(v int) string {
	if v == 0 {
		return ""
	}
	return fmt.Sprintf("%d", v)
}

func optionalFloatStr(v float32) string {
	if v == 0 {
		return ""
	}
	return fmt.Sprintf("%g", float64(v))
}
