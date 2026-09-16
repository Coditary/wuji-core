package data

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func parseDocumentMap(doc map[string]any) (DataShape, error) {
	kind, _ := doc["kind"].(string)
	switch ShapeKind(kind) {
	case ShapeVector:
		if batch, err := documentToVectorBatch(doc); err == nil {
			return batch, nil
		}
		rs, err := documentToRecordSet(doc)
		if err != nil {
			return nil, err
		}
		return RecordSetToVectorBatch(rs)
	case ShapeTable:
		rs, err := documentToRecordSet(doc)
		if err != nil {
			return nil, err
		}
		return RecordSetToTable(rs), nil
	case ShapeGraph:
		rs, err := documentToRecordSet(doc)
		if err != nil {
			return nil, err
		}
		return FromRecordSetGraph(rs), nil
	default:
		rs, err := documentToRecordSet(doc)
		if err != nil {
			return nil, err
		}
		return rs, nil
	}
}

func documentToRecordSet(doc map[string]any) (*RecordSet, error) {
	rs := &RecordSet{}
	if meta, ok := doc["meta"].(map[string]any); ok {
		rs.MetaData = metaFromMap(meta)
	}
	if schema, ok := doc["schema"].([]any); ok {
		for _, item := range schema {
			if m, ok := item.(map[string]any); ok {
				rs.Schema = append(rs.Schema, FieldDef{
					Name: str(m["name"]),
					Type: ValueKind(str(m["type"])),
					Dims: int(num(m["dims"])),
				})
			}
		}
	}
	recs, ok := doc["records"].([]any)
	if !ok {
		return rs, nil
	}
	for _, item := range recs {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rec := Record{
			ID:     str(m["id"]),
			Role:   str(m["role"]),
			Parent: str(m["parent"]),
			Seq:    int(num(m["seq"])),
			Fields: map[string]Value{},
		}
		fields, _ := m["fields"].(map[string]any)
		for k, v := range fields {
			rec.Fields[k] = anyToValue(v)
		}
		rs.Records = append(rs.Records, rec)
	}
	links, _ := doc["links"].([]any)
	for _, item := range links {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		link := Link{
			ID:   str(m["id"]),
			Type: str(m["type"]),
			From: str(m["from"]),
			To:   str(m["to"]),
		}
		if d, ok := m["directed"].(bool); ok {
			link.Directed = &d
		}
		if w := num(m["weight"]); m["weight"] != nil {
			weight := w
			link.Weight = &weight
		}
		if fields, ok := m["fields"].(map[string]any); ok {
			link.Fields = map[string]Value{}
			for k, v := range fields {
				link.Fields[k] = anyToValue(v)
			}
		}
		rs.Links = append(rs.Links, link)
	}
	return rs, nil
}

func documentToVectorBatch(doc map[string]any) (*VectorBatch, error) {
	tensor, ok := doc["tensor"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("no tensor block")
	}
	dims := int(num(tensor["dims"]))
	if dims <= 0 {
		return nil, fmt.Errorf("invalid tensor dims")
	}
	dataRaw, ok := tensor["data"]
	if !ok {
		return nil, fmt.Errorf("tensor data missing")
	}
	vecs, err := decodeFloat32Slice(dataRaw)
	if err != nil {
		return nil, err
	}
	rs, err := documentToRecordSet(doc)
	if err != nil {
		return nil, err
	}
	inputs := make([]string, 0, len(rs.Records))
	for _, rec := range rs.Records {
		in := rec.Fields["input"].S
		if in == "" {
			in = rec.Fields["text"].S
		}
		if in == "" {
			in = rec.ID
		}
		inputs = append(inputs, in)
	}
	if len(inputs) == 0 {
		count := len(vecs) / dims
		inputs = make([]string, count)
		for i := range inputs {
			inputs[i] = fmt.Sprintf("v%d", i)
		}
	}
	batch := &VectorBatch{Dims: dims, Inputs: inputs, Vectors: vecs, MetaData: rs.MetaData}
	if batch.Count()*dims != len(vecs) {
		return nil, fmt.Errorf("tensor size mismatch")
	}
	return batch, nil
}

func decodeFloat32Slice(raw any) ([]float32, error) {
	switch t := raw.(type) {
	case []float32:
		out := make([]float32, len(t))
		copy(out, t)
		return out, nil
	case []float64:
		out := make([]float32, len(t))
		for i, v := range t {
			out[i] = float32(v)
		}
		return out, nil
	case []any:
		out := make([]float32, len(t))
		for i, v := range t {
			out[i] = float32(num(v))
		}
		return out, nil
	case []byte:
		if len(t)%4 != 0 {
			return nil, fmt.Errorf("invalid float32 byte length")
		}
		out := make([]float32, len(t)/4)
		for i := range out {
			bits := uint32(t[i*4]) | uint32(t[i*4+1])<<8 | uint32(t[i*4+2])<<16 | uint32(t[i*4+3])<<24
			out[i] = math.Float32frombits(bits)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported tensor data type %T", raw)
	}
}

func metaFromMap(m map[string]any) Meta {
	meta := Meta{Source: map[string]string{}, Extra: map[string]string{}}
	meta.Capability = str(m["capability"])
	meta.Task = str(m["task"])
	meta.Driver = str(m["driver"])
	meta.Model = str(m["model"])
	if src, ok := m["source"].(map[string]any); ok {
		for k, v := range src {
			meta.Source[k] = str(v)
		}
	}
	return meta
}

// TableToVectorBatch converts tabular data with embedding columns to vectors.
func TableToVectorBatch(t *Table) (*VectorBatch, error) {
	if rsBatch, err := RecordSetToVectorBatch(t.ToRecordSet()); err == nil {
		return rsBatch, nil
	}
	if t.RowCount() == 0 {
		return nil, fmt.Errorf("cannot convert empty table to vector")
	}

	inputCol := tableInputColumn(t)
	vecCols := tableVectorColumns(t, inputCol)
	if len(vecCols) == 0 {
		return nil, fmt.Errorf("cannot convert table to vector: need embedding column(s) or numeric columns beside input/text (or run --embed to compute vectors)")
	}

	dims := len(vecCols)
	inputs := make([]string, t.RowCount())
	rows := make([][]float32, t.RowCount())
	for r := 0; r < t.RowCount(); r++ {
		if inputCol >= 0 && r < len(t.Columns[inputCol].Cells) {
			inputs[r] = cellString(t.Columns[inputCol].Cells[r])
		} else {
			inputs[r] = fmt.Sprintf("row%d", r)
		}
		row := make([]float32, dims)
		for i, colIdx := range vecCols {
			if r < len(t.Columns[colIdx].Cells) {
				row[i] = valueToFloat32(t.Columns[colIdx].Cells[r])
			}
		}
		rows[r] = row
	}
	batch, err := NewVectorBatch(inputs, dims, rows)
	if err != nil {
		return nil, err
	}
	batch.MetaData = t.MetaData
	return batch, nil
}

func tableInputColumn(t *Table) int {
	for i, col := range t.Columns {
		switch col.Name {
		case "input", "text":
			return i
		}
	}
	for i, col := range t.Columns {
		if col.Kind == KindString {
			return i
		}
	}
	return -1
}

func tableVectorColumns(t *Table, inputCol int) []int {
	var explicit []int
	for i, col := range t.Columns {
		if i == inputCol {
			continue
		}
		if strings.HasPrefix(col.Name, "embedding_") {
			if n, err := strconv.Atoi(strings.TrimPrefix(col.Name, "embedding_")); err == nil {
				for len(explicit) <= n {
					explicit = append(explicit, -1)
				}
				explicit[n] = i
			}
		}
	}
	if len(explicit) > 0 {
		out := make([]int, 0, len(explicit))
		for _, idx := range explicit {
			if idx >= 0 {
				out = append(out, idx)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	var numeric []int
	for i, col := range t.Columns {
		if i == inputCol {
			continue
		}
		if col.Name == "embedding" {
			continue
		}
		if col.Kind == KindFloat || col.Kind == KindInt {
			numeric = append(numeric, i)
		}
	}
	if len(numeric) >= 2 {
		return numeric
	}
	return nil
}

func valueToFloat32(v Value) float32 {
	switch v.Kind {
	case KindFloat:
		return float32(v.F)
	case KindInt:
		return float32(v.I)
	default:
		f, _ := strconv.ParseFloat(v.S, 32)
		return float32(f)
	}
}
