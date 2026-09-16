package data

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
)

func writeCSVRow(w io.Writer, fields []string) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(fields); err != nil {
		return err
	}
	cw.Flush()
	return cw.Error()
}

// ParseJSON reads a WDD JSON document into the closest shape.
func ParseJSON(r io.Reader) (DataShape, error) {
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r).Decode(&raw); err != nil {
		return nil, err
	}
	kind := ""
	if k, ok := raw["kind"]; ok {
		_ = json.Unmarshal(k, &kind)
	}
	switch ShapeKind(kind) {
	case ShapeVector:
		return parseJSONVector(raw)
	case ShapeTable:
		return parseJSONTable(raw)
	case ShapeGraph:
		return parseJSONGraph(raw)
	default:
		return parseJSONRecordSet(raw)
	}
}

func parseJSONRecordSet(raw map[string]json.RawMessage) (*RecordSet, error) {
	rs := &RecordSet{}
	if m, ok := raw["meta"]; ok {
		_ = json.Unmarshal(m, &rs.MetaData)
	}
	if s, ok := raw["schema"]; ok {
		_ = json.Unmarshal(s, &rs.Schema)
	}
	if recs, ok := raw["records"]; ok {
		var wire []recordWire
		if err := json.Unmarshal(recs, &wire); err != nil {
			return nil, err
		}
		rs.Records = decodeRecords(wire)
	}
	if links, ok := raw["links"]; ok {
		var wire []linkWire
		if err := json.Unmarshal(links, &wire); err != nil {
			return nil, err
		}
		rs.Links = decodeLinks(wire)
	}
	return rs, nil
}

func parseJSONGraph(raw map[string]json.RawMessage) (*Graph, error) {
	rs, err := parseJSONRecordSet(raw)
	if err != nil {
		return nil, err
	}
	return FromRecordSetGraph(rs), nil
}

func parseJSONTable(raw map[string]json.RawMessage) (*Table, error) {
	rs, err := parseJSONRecordSet(raw)
	if err != nil {
		return nil, err
	}
	return RecordSetToTable(rs), nil
}

func parseJSONVector(raw map[string]json.RawMessage) (*VectorBatch, error) {
	rs, err := parseJSONRecordSet(raw)
	if err != nil {
		return nil, err
	}
	return RecordSetToVectorBatch(rs)
}

type recordWire struct {
	ID     string                 `json:"id"`
	Role   string                 `json:"role"`
	Parent string                 `json:"parent"`
	Seq    int                    `json:"seq"`
	Fields map[string]json.RawMessage `json:"fields"`
}

type linkWire struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	From     string                 `json:"from"`
	To       string                 `json:"to"`
	Directed *bool                  `json:"directed"`
	Weight   *float64               `json:"weight"`
	Fields   map[string]json.RawMessage `json:"fields"`
}

func decodeRecords(wire []recordWire) []Record {
	out := make([]Record, len(wire))
	for i, w := range wire {
		fields := map[string]Value{}
		for k, raw := range w.Fields {
			fields[k] = decodeValue(raw)
		}
		out[i] = Record{
			ID: w.ID, Role: w.Role, Parent: w.Parent, Seq: w.Seq, Fields: fields,
		}
	}
	return out
}

func decodeLinks(wire []linkWire) []Link {
	out := make([]Link, len(wire))
	for i, w := range wire {
		fields := map[string]Value{}
		for k, raw := range w.Fields {
			fields[k] = decodeValue(raw)
		}
		out[i] = Link{
			ID: w.ID, Type: w.Type, From: w.From, To: w.To,
			Directed: w.Directed, Weight: w.Weight, Fields: fields,
		}
	}
	return out
}

func decodeValue(raw json.RawMessage) Value {
	var generic any
	if err := json.Unmarshal(raw, &generic); err != nil {
		return Value{Kind: KindString, S: string(raw)}
	}
	return anyToValue(generic)
}

func anyToValue(v any) Value {
	switch t := v.(type) {
	case nil:
		return Value{Kind: KindNull}
	case bool:
		return Value{Kind: KindBool, B: t}
	case float64:
		if float64(int64(t)) == t {
			return Value{Kind: KindInt, I: int64(t)}
		}
		return Value{Kind: KindFloat, F: t}
	case string:
		return Value{Kind: KindString, S: t}
	case []any:
		arr := make([]Value, len(t))
		for i, item := range t {
			arr[i] = anyToValue(item)
		}
		return Value{Kind: KindArray, A: arr}
	case map[string]any:
		typ, _ := t["type"].(string)
		if typ == "vector" {
			vals, _ := t["values"].([]any)
			vec := make([]float32, len(vals))
			for i, x := range vals {
				if f, ok := x.(float64); ok {
					vec[i] = float32(f)
				}
			}
			return Value{Kind: KindVector, Vec: vec}
		}
		if typ == "bbox" {
			return Value{Kind: KindBBox, BBox: BBox{
				X: num(t["x"]), Y: num(t["y"]), W: num(t["w"]), H: num(t["h"]),
				Normalized: boolVal(t["normalized"]),
			}}
		}
		if typ == "span" {
			return Value{Kind: KindSpan, Span: Span{
				Start: int(num(t["start"])),
				End:   int(num(t["end"])),
				Text:  str(t["text"]),
			}}
		}
		obj := map[string]Value{}
		for k, item := range t {
			obj[k] = anyToValue(item)
		}
		return Value{Kind: KindObject, O: obj}
	default:
		return Value{Kind: KindString, S: fmt.Sprint(v)}
	}
}

func num(v any) float64 {
	f, _ := v.(float64)
	return f
}

func str(v any) string {
	s, _ := v.(string)
	return s
}

func boolVal(v any) bool {
	b, _ := v.(bool)
	return b
}

// RecordSetToTable builds a columnar table from row records.
func RecordSetToTable(rs *RecordSet) *Table {
	if len(rs.Records) == 0 {
		return &Table{MetaData: rs.MetaData, Schema: rs.Schema}
	}
	names := orderedRecordColumns(rs.Records, DefaultExportOpts())
	cols := make([]Column, 0, len(names))
	for _, name := range names {
		if name == "_id" || name == "_role" {
			continue
		}
		col := Column{Name: name, Kind: KindString}
		for _, rec := range rs.Records {
			if val, ok := rec.Fields[name]; ok {
				col.Cells = append(col.Cells, val)
				col.Kind = val.Kind
			} else {
				col.Cells = append(col.Cells, Value{Kind: KindNull})
			}
		}
		cols = append(cols, col)
	}
	return &Table{MetaData: rs.MetaData, Schema: rs.Schema, Columns: cols}
}

// RecordSetToVectorBatch extracts embeddings from records.
func RecordSetToVectorBatch(rs *RecordSet) (*VectorBatch, error) {
	inputs := make([]string, 0, len(rs.Records))
	rows := make([][]float32, 0, len(rs.Records))
	dims := 0
	for _, rec := range rs.Records {
		input := rec.Fields["input"].S
		if input == "" {
			input = rec.Fields["text"].S
		}
		emb, ok := rec.Fields["embedding"]
		if !ok || emb.Kind != KindVector {
			continue
		}
		if dims == 0 {
			dims = len(emb.Vec)
		}
		inputs = append(inputs, input)
		rows = append(rows, emb.Vec)
	}
	if dims == 0 {
		return nil, fmt.Errorf("no embeddings found in record set")
	}
	batch, err := NewVectorBatch(inputs, dims, rows)
	if err != nil {
		return nil, err
	}
	batch.MetaData = rs.MetaData
	return batch, nil
}

// AsTable converts supported shapes to Table.
func AsTable(s DataShape) (*Table, bool) {
	switch v := s.(type) {
	case *Table:
		return v, true
	case *RecordSet:
		return RecordSetToTable(v), true
	case *VectorBatch:
		return RecordSetToTable(v.ToRecordSet()), true
	case *Graph:
		return RecordSetToTable(v.ToRecordSet()), true
	default:
		return nil, false
	}
}

// ConvertShape transforms between compatible shapes.
func ConvertShape(s DataShape, target ShapeKind) (DataShape, error) {
	if s.Kind() == target {
		return s, nil
	}
	switch target {
	case ShapeRecordSet:
		switch v := s.(type) {
		case *RecordSet:
			return v, nil
		case *Table:
			return v.ToRecordSet(), nil
		case *VectorBatch:
			return v.ToRecordSet(), nil
		case *Graph:
			return v.ToRecordSet(), nil
		}
	case ShapeTable:
		switch v := s.(type) {
		case *Table:
			return v, nil
		case *RecordSet:
			return RecordSetToTable(v), nil
		case *VectorBatch:
			return RecordSetToTable(v.ToRecordSet()), nil
		case *Graph:
			return RecordSetToTable(v.ToRecordSet()), nil
		}
	case ShapeVector:
		switch v := s.(type) {
		case *VectorBatch:
			return v, nil
		case *RecordSet:
			return RecordSetToVectorBatch(v)
		case *Table:
			return TableToVectorBatch(v)
		}
	case ShapeGraph:
		switch v := s.(type) {
		case *Graph:
			return v, nil
		case *RecordSet:
			return FromRecordSetGraph(v), nil
		}
	}
	return nil, fmt.Errorf("cannot convert %s to %s", s.Kind(), target)
}
