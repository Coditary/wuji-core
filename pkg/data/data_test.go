package data

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestVectorBatchExportJSON(t *testing.T) {
	batch, err := NewVectorBatch([]string{"hello"}, 4, [][]float32{{0.1, 0.2, 0.3, 0.4}})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := batch.Export(&buf, DefaultExportOpts()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"wuji": "1"`) {
		t.Fatalf("expected wdd envelope, got %s", buf.String())
	}
}

func TestTableCSVRoundTrip(t *testing.T) {
	csv := "a,b\n1,2\n3,4\n"
	table, err := ParseCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}
	if table.RowCount() != 2 {
		t.Fatalf("rows = %d", table.RowCount())
	}
	var buf bytes.Buffer
	opts := DefaultExportOpts()
	opts.Format = FormatCSV
	if err := table.Export(&buf, opts); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "a,b") {
		t.Fatalf("csv header missing: %q", buf.String())
	}
}

func TestConvertShape(t *testing.T) {
	batch, err := NewVectorBatch([]string{"x"}, 2, [][]float32{{1, 2}})
	if err != nil {
		t.Fatal(err)
	}
	rs, err := ConvertShape(batch, ShapeRecordSet)
	if err != nil {
		t.Fatal(err)
	}
	if rs.Kind() != ShapeRecordSet {
		t.Fatalf("kind=%s", rs.Kind())
	}
}

func TestParseJSONRecordSet(t *testing.T) {
	raw := `{"wuji":"1","kind":"recordset","records":[{"id":"r0","fields":{"label":"car"}}]}`
	shape, err := ParseJSON(strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	rs, ok := shape.(*RecordSet)
	if !ok || len(rs.Records) != 1 {
		t.Fatalf("unexpected shape %#v", shape)
	}
}

func TestRecordSetCSVFlattenBBox(t *testing.T) {
	rs := &RecordSet{
		Records: []Record{{
			ID: "r0", Fields: map[string]Value{
				"bbox": NewBBox(BBox{X: 1, Y: 2, W: 3, H: 4}),
			},
		}},
	}
	var buf bytes.Buffer
	opts := DefaultExportOpts()
	opts.Format = FormatCSV
	if err := rs.Export(&buf, opts); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "bbox.x") || !strings.Contains(out, "bbox.y") {
		t.Fatalf("flatten failed: %q", out)
	}
}

func TestGraphLinksCSV(t *testing.T) {
	dir := true
	g := &Graph{
		Nodes: []Record{{ID: "n1", Role: "node", Fields: map[string]Value{"label": NewString("A")}}},
		Edges: []Link{{ID: "e1", Type: "edge", From: "n1", To: "n2", Directed: &dir}},
	}
	var buf bytes.Buffer
	opts := DefaultExportOpts()
	opts.Format = FormatCSV
	opts.CSVView = CSVViewLinks
	if err := g.Export(&buf, opts); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "n1,n2") {
		t.Fatalf("links csv: %q", buf.String())
	}
}

func TestMsgpackRoundTrip(t *testing.T) {
	batch, err := NewVectorBatch([]string{"a"}, 3, [][]float32{{1, 2, 3}})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := batch.Export(&buf, ExportOpts{Format: FormatMsgpack}); err != nil {
		t.Fatal(err)
	}
	shape, err := ParseMsgpack(&buf)
	if err != nil {
		t.Fatal(err)
	}
	vb, ok := shape.(*VectorBatch)
	if !ok || vb.Dims != 3 {
		t.Fatalf("shape=%T dims=%d", shape, vb.Dims)
	}
}

func TestTableToVectorBatch(t *testing.T) {
	table, err := ParseCSV(strings.NewReader("text,e0,e1\nhello,0.1,0.2\n"))
	if err != nil {
		t.Fatal(err)
	}
	batch, err := TableToVectorBatch(table)
	if err != nil {
		t.Fatal(err)
	}
	if batch.Dims != 2 || batch.Inputs[0] != "hello" {
		t.Fatalf("batch=%+v", batch)
	}
}

func TestTableToVectorBatchPlainCSVFails(t *testing.T) {
	table, err := ParseCSV(strings.NewReader("name,age\nalice,30\n"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = ConvertShape(table, ShapeVector)
	if err == nil {
		t.Fatal("expected error for non-embedding table")
	}
}

func TestEnvelopeJSONValid(t *testing.T) {
	rs := &RecordSet{Records: []Record{{ID: "r0", Fields: map[string]Value{"x": NewInt(1)}}}}
	doc := envelopeFromRecordSet(rs)
	if _, err := json.Marshal(doc); err != nil {
		t.Fatal(err)
	}
}
