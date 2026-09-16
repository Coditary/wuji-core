package data

import (
	"encoding/json"
	"fmt"
)

// ValueKind identifies a typed field value.
type ValueKind string

const (
	KindNull      ValueKind = "null"
	KindBool      ValueKind = "bool"
	KindInt       ValueKind = "int"
	KindFloat     ValueKind = "float"
	KindString    ValueKind = "string"
	KindArray     ValueKind = "array"
	KindObject    ValueKind = "object"
	KindVector    ValueKind = "vector"
	KindBBox      ValueKind = "bbox"
	KindSpan      ValueKind = "span"
	KindTimestamp ValueKind = "timestamp"
)

// Value is a typed dynamic field (no interface{} boxing in hot paths).
type Value struct {
	Kind ValueKind
	B    bool
	I    int64
	F    float64
	S    string
	A    []Value
	O    map[string]Value
	Vec  []float32
	BBox BBox
	Span Span
}

// BBox is an axis-aligned rectangle.
type BBox struct {
	X, Y, W, H float64
	Normalized bool
}

// Span is a text span (e.g. NER).
type Span struct {
	Start int
	End   int
	Text  string
}

// FieldDef describes one column/field in a schema.
type FieldDef struct {
	Name string    `json:"name"`
	Type ValueKind `json:"type"`
	Dims int       `json:"dims,omitempty"`
}

// NewString returns a string value.
func NewString(s string) Value { return Value{Kind: KindString, S: s} }

// NewFloat returns a float value.
func NewFloat(f float64) Value { return Value{Kind: KindFloat, F: f} }

// NewInt returns an int value.
func NewInt(i int64) Value { return Value{Kind: KindInt, I: i} }

// NewBool returns a bool value.
func NewBool(b bool) Value { return Value{Kind: KindBool, B: b} }

// NewVector returns a copy of embedding values.
func NewVector(v []float32) Value {
	out := make([]float32, len(v))
	copy(out, v)
	return Value{Kind: KindVector, Vec: out}
}

// NewBBox returns a bounding box value.
func NewBBox(b BBox) Value { return Value{Kind: KindBBox, BBox: b} }

// NewSpan returns a span value.
func NewSpan(s Span) Value { return Value{Kind: KindSpan, Span: s} }

func (v Value) JSONAny() any {
	switch v.Kind {
	case KindNull:
		return nil
	case KindBool:
		return v.B
	case KindInt:
		return v.I
	case KindFloat:
		return v.F
	case KindString, KindTimestamp:
		return v.S
	case KindArray:
		items := make([]any, len(v.A))
		for i := range v.A {
			items[i] = v.A[i].JSONAny()
		}
		return items
	case KindObject:
		m := make(map[string]any, len(v.O))
		for k, val := range v.O {
			m[k] = val.JSONAny()
		}
		return m
	case KindVector:
		vals := make([]float64, len(v.Vec))
		for i, x := range v.Vec {
			vals[i] = float64(x)
		}
		return map[string]any{
			"type":   "vector",
			"dims":   len(v.Vec),
			"dtype":  "f32",
			"values": vals,
		}
	case KindBBox:
		return map[string]any{
			"type":       "bbox",
			"x":          v.BBox.X,
			"y":          v.BBox.Y,
			"w":          v.BBox.W,
			"h":          v.BBox.H,
			"normalized": v.BBox.Normalized,
		}
	case KindSpan:
		return map[string]any{
			"type":  "span",
			"start": v.Span.Start,
			"end":   v.Span.End,
			"text":  v.Span.Text,
		}
	default:
		return nil
	}
}

func (v Value) flattenColumns(prefix string) map[string]string {
	out := map[string]string{}
	switch v.Kind {
	case KindNull:
		if prefix != "" {
			out[prefix] = ""
		}
	case KindBool:
		out[prefix] = fmt.Sprintf("%t", v.B)
	case KindInt:
		out[prefix] = fmt.Sprintf("%d", v.I)
	case KindFloat:
		out[prefix] = fmt.Sprintf("%g", v.F)
	case KindString, KindTimestamp:
		out[prefix] = v.S
	case KindVector:
		raw, _ := json.Marshal(v.JSONAny())
		out[prefix] = string(raw)
	case KindBBox:
		out[prefix+".x"] = fmt.Sprintf("%g", v.BBox.X)
		out[prefix+".y"] = fmt.Sprintf("%g", v.BBox.Y)
		out[prefix+".w"] = fmt.Sprintf("%g", v.BBox.W)
		out[prefix+".h"] = fmt.Sprintf("%g", v.BBox.H)
	case KindSpan:
		out[prefix+".start"] = fmt.Sprintf("%d", v.Span.Start)
		out[prefix+".end"] = fmt.Sprintf("%d", v.Span.End)
		out[prefix+".text"] = v.Span.Text
	case KindArray, KindObject:
		raw, _ := json.Marshal(v.JSONAny())
		out[prefix] = string(raw)
	}
	return out
}

func (v Value) expandColumns(prefix string, mode VectorCSVMode) map[string]string {
	if v.Kind == KindVector && mode == VectorCSVExpand {
		out := map[string]string{}
		for i, x := range v.Vec {
			out[fmt.Sprintf("%s_%d", prefix, i)] = fmt.Sprintf("%g", x)
		}
		return out
	}
	if v.Kind == KindVector && mode == VectorCSVDims {
		out := map[string]string{}
		out[prefix+"_dims"] = fmt.Sprintf("%d", len(v.Vec))
		raw, _ := json.Marshal(v.JSONAny())
		out[prefix] = string(raw)
		return out
	}
	return v.flattenColumns(prefix)
}
