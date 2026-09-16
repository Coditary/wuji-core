package data

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Column holds one typed column in a table.
type Column struct {
	Name  string
	Kind  ValueKind
	Cells []Value
}

// Table is a columnar representation for tabular / forecast data.
type Table struct {
	MetaData Meta     `json:"-"`
	Schema   []FieldDef `json:"schema,omitempty"`
	Columns  []Column   `json:"-"`
}

func (t *Table) Kind() ShapeKind { return ShapeTable }
func (t *Table) Meta() Meta      { return t.MetaData }

func (t *Table) RowCount() int {
	if len(t.Columns) == 0 {
		return 0
	}
	return len(t.Columns[0].Cells)
}

func (t *Table) Export(w io.Writer, opts ExportOpts) error {
	switch opts.Format {
	case FormatCSV:
		return t.exportCSV(w)
	case FormatMsgpack:
		return WriteMsgpack(w, t)
	default:
		rs := t.ToRecordSet()
		return rs.Export(w, opts)
	}
}

func (t *Table) exportCSV(w io.Writer) error {
	if len(t.Columns) == 0 {
		return writeCSVRow(w, nil)
	}
	header := make([]string, len(t.Columns))
	for i, col := range t.Columns {
		header[i] = col.Name
	}
	if err := writeCSVRow(w, header); err != nil {
		return err
	}
	rows := t.RowCount()
	for r := 0; r < rows; r++ {
		row := make([]string, len(t.Columns))
		for c, col := range t.Columns {
			if r < len(col.Cells) {
				row[c] = cellString(col.Cells[r])
			}
		}
		if err := writeCSVRow(w, row); err != nil {
			return err
		}
	}
	return nil
}

func cellString(v Value) string {
	switch v.Kind {
	case KindNull:
		return ""
	case KindBool:
		return fmt.Sprintf("%t", v.B)
	case KindInt:
		return fmt.Sprintf("%d", v.I)
	case KindFloat:
		return fmt.Sprintf("%g", v.F)
	case KindString, KindTimestamp:
		return v.S
	default:
		raw, _ := json.Marshal(v.JSONAny())
		return string(raw)
	}
}

// ToRecordSet converts columns to row records.
func (t *Table) ToRecordSet() *RecordSet {
	rows := t.RowCount()
	recs := make([]Record, rows)
	for i := 0; i < rows; i++ {
		fields := map[string]Value{}
		for _, col := range t.Columns {
			if i < len(col.Cells) {
				fields[col.Name] = col.Cells[i]
			}
		}
		recs[i] = Record{
			ID:     fmt.Sprintf("row%d", i),
			Role:   "row",
			Seq:    i,
			Fields: fields,
		}
	}
	return &RecordSet{
		MetaData: t.MetaData,
		Schema:   t.Schema,
		Records:  recs,
	}
}

// ParseCSVFile loads a CSV file into a Table shape.
func ParseCSVFile(path string) (*Table, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseCSV(f)
}

// ParseCSV reads CSV from r into a Table.
func ParseCSV(r io.Reader) (*Table, error) {
	reader := csv.NewReader(r)
	reader.ReuseRecord = true
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return &Table{}, nil
	}
	header := rows[0]
	cols := make([]Column, len(header))
	schema := make([]FieldDef, len(header))
	for i, name := range header {
		cols[i] = Column{Name: name, Kind: KindString, Cells: make([]Value, 0, len(rows)-1)}
		schema[i] = FieldDef{Name: name, Type: KindString}
	}
	for _, row := range rows[1:] {
		for i := range cols {
			val := KindString
			cell := ""
			if i < len(row) {
				cell = row[i]
				if f, err := strconv.ParseFloat(cell, 64); err == nil {
					cols[i].Cells = append(cols[i].Cells, NewFloat(f))
					cols[i].Kind = KindFloat
					schema[i].Type = KindFloat
					continue
				}
				if n, err := strconv.ParseInt(cell, 10, 64); err == nil {
					cols[i].Cells = append(cols[i].Cells, NewInt(n))
					cols[i].Kind = KindInt
					schema[i].Type = KindInt
					continue
				}
			}
			_ = val
			cols[i].Cells = append(cols[i].Cells, NewString(cell))
		}
	}
	return &Table{Schema: schema, Columns: cols}, nil
}
