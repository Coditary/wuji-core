package data

import (
	"encoding/json"
	"fmt"
	"io"
)

// VectorBatch stores embeddings in a contiguous float32 slab.
type VectorBatch struct {
	MetaData Meta
	Dims     int
	Inputs   []string
	Vectors  []float32 // len = len(Inputs) * Dims
}

func (v *VectorBatch) Kind() ShapeKind { return ShapeVector }
func (v *VectorBatch) Meta() Meta      { return v.MetaData }

func (v *VectorBatch) Count() int { return len(v.Inputs) }

func (v *VectorBatch) At(i int) []float32 {
	if v.Dims <= 0 || i < 0 || i >= len(v.Inputs) {
		return nil
	}
	start := i * v.Dims
	end := start + v.Dims
	if end > len(v.Vectors) {
		return nil
	}
	return v.Vectors[start:end]
}

func (v *VectorBatch) Export(w io.Writer, opts ExportOpts) error {
	switch opts.Format {
	case FormatCSV:
		return exportVectorCSV(w, v, opts)
	case FormatMsgpack:
		return WriteMsgpack(w, v)
	default:
		rs := v.ToRecordSet()
		return rs.Export(w, opts)
	}
}

func exportVectorCSV(w io.Writer, v *VectorBatch, opts ExportOpts) error {
	records := make([]Record, v.Count())
	for i := 0; i < v.Count(); i++ {
		fields := map[string]Value{
			"input":     NewString(v.Inputs[i]),
			"embedding": NewVector(v.At(i)),
		}
		records[i] = Record{
			ID:     fmt.Sprintf("v%d", i),
			Role:   "embedding",
			Seq:    i,
			Fields: fields,
		}
	}
	return writeRecordsCSV(w, records, opts)
}

// ToRecordSet materializes one record per embedding.
func (v *VectorBatch) ToRecordSet() *RecordSet {
	recs := make([]Record, v.Count())
	for i := 0; i < v.Count(); i++ {
		fields := map[string]Value{
			"input":     NewString(v.Inputs[i]),
			"embedding": NewVector(v.At(i)),
		}
		recs[i] = Record{
			ID:     fmt.Sprintf("v%d", i),
			Role:   "embedding",
			Seq:    i,
			Fields: fields,
		}
	}
	schema := []FieldDef{
		{Name: "input", Type: KindString},
		{Name: "embedding", Type: KindVector, Dims: v.Dims},
	}
	return &RecordSet{
		MetaData: v.MetaData,
		Schema:   schema,
		Records:  recs,
	}
}

// NewVectorBatch builds a batch from inputs and per-row vectors.
func NewVectorBatch(inputs []string, dims int, rows [][]float32) (*VectorBatch, error) {
	if dims <= 0 {
		return nil, fmt.Errorf("vector dims must be > 0")
	}
	if len(rows) != len(inputs) {
		return nil, fmt.Errorf("vector row count %d != input count %d", len(rows), len(inputs))
	}
	slab := make([]float32, 0, len(inputs)*dims)
	for i, row := range rows {
		if len(row) != dims {
			return nil, fmt.Errorf("vector row %d has %d dims, want %d", i, len(row), dims)
		}
		slab = append(slab, row...)
	}
	in := make([]string, len(inputs))
	copy(in, inputs)
	return &VectorBatch{Dims: dims, Inputs: in, Vectors: slab}, nil
}

func (v *VectorBatch) MarshalDocument() ([]byte, error) {
	rs := v.ToRecordSet()
	doc := envelopeFromRecordSet(rs)
	doc["kind"] = string(ShapeVector)
	doc["tensor"] = map[string]any{
		"dims":   v.Dims,
		"count":  v.Count(),
		"dtype":  "f32",
		"layout": "row_major",
	}
	return json.Marshal(doc)
}
