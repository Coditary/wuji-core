package data

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// Record holds one structured item.
type Record struct {
	ID     string           `json:"id,omitempty"`
	Role   string           `json:"role,omitempty"`
	Parent string           `json:"parent,omitempty"`
	Seq    int              `json:"seq,omitempty"`
	Fields map[string]Value `json:"fields"`
}

// Link connects two records (graph edges, hierarchies, …).
type Link struct {
	ID       string           `json:"id,omitempty"`
	Type     string           `json:"type"`
	From     string           `json:"from"`
	To       string           `json:"to"`
	Directed *bool            `json:"directed,omitempty"`
	Weight   *float64         `json:"weight,omitempty"`
	Fields   map[string]Value `json:"fields,omitempty"`
}

// RecordSet is a list of records with optional links.
type RecordSet struct {
	MetaData Meta       `json:"-"`
	Schema   []FieldDef `json:"schema,omitempty"`
	Records  []Record   `json:"records"`
	Links    []Link     `json:"links,omitempty"`
}

func (r *RecordSet) Kind() ShapeKind { return ShapeRecordSet }
func (r *RecordSet) Meta() Meta      { return r.MetaData }

func (r *RecordSet) Export(w io.Writer, opts ExportOpts) error {
	switch opts.Format {
	case FormatCSV:
		return exportRecordSetCSV(w, r, opts)
	case FormatNDJSON:
		return exportRecordSetNDJSON(w, r, opts)
	case FormatMsgpack:
		return WriteMsgpack(w, r)
	default:
		return exportRecordSetJSON(w, r, opts)
	}
}

func exportRecordSetJSON(w io.Writer, rs *RecordSet, opts ExportOpts) error {
	doc := envelopeFromRecordSet(rs)
	enc := json.NewEncoder(w)
	if opts.PrettyJSON {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(doc)
}

func exportRecordSetNDJSON(w io.Writer, rs *RecordSet, opts ExportOpts) error {
	enc := json.NewEncoder(w)
	for _, rec := range rs.Records {
		row := map[string]any{"wuji": "1", "record": recordJSON(rec)}
		if err := enc.Encode(row); err != nil {
			return err
		}
	}
	return nil
}

func exportRecordSetCSV(w io.Writer, rs *RecordSet, opts ExportOpts) error {
	switch opts.CSVView {
	case CSVViewLinks:
		return writeLinksCSV(w, rs.Links)
	default:
		return writeRecordsCSV(w, rs.Records, opts)
	}
}

func recordJSON(rec Record) map[string]any {
	fields := map[string]any{}
	for k, v := range rec.Fields {
		fields[k] = v.JSONAny()
	}
	out := map[string]any{"fields": fields}
	if rec.ID != "" {
		out["id"] = rec.ID
	}
	if rec.Role != "" {
		out["role"] = rec.Role
	}
	if rec.Parent != "" {
		out["parent"] = rec.Parent
	}
	if rec.Seq != 0 {
		out["seq"] = rec.Seq
	}
	return out
}

func writeRecordsCSV(w io.Writer, records []Record, opts ExportOpts) error {
	columns := orderedRecordColumns(records, opts)
	if err := writeCSVRow(w, columns); err != nil {
		return err
	}
	for _, rec := range records {
		row := make([]string, len(columns))
		flat := map[string]string{"_id": rec.ID, "_role": rec.Role}
		for name, val := range rec.Fields {
			for k, v := range val.expandColumns(name, opts.VectorCSV) {
				if opts.PreviewLength > 0 && val.Kind == KindString && k == name && len(v) > opts.PreviewLength {
					v = v[:opts.PreviewLength]
				}
				flat[k] = v
			}
		}
		for i, col := range columns {
			row[i] = flat[col]
		}
		if err := writeCSVRow(w, row); err != nil {
			return err
		}
	}
	return nil
}

func orderedRecordColumns(records []Record, opts ExportOpts) []string {
	if len(opts.CSVFields) > 0 {
		cols := append([]string{"_id", "_role"}, opts.CSVFields...)
		return cols
	}
	seen := map[string]struct{}{"_id": {}, "_role": {}}
	for _, rec := range records {
		for name, val := range rec.Fields {
			if val.Kind == KindVector && opts.VectorCSV == VectorCSVExpand {
				for i := range val.Vec {
					seen[fmt.Sprintf("%s_%d", name, i)] = struct{}{}
				}
				continue
			}
			if val.Kind == KindVector && opts.VectorCSV == VectorCSVDims {
				seen[name+"_dims"] = struct{}{}
				seen[name] = struct{}{}
				continue
			}
			for k := range val.flattenColumns(name) {
				seen[k] = struct{}{}
			}
		}
	}
	cols := make([]string, 0, len(seen))
	for c := range seen {
		cols = append(cols, c)
	}
	sort.Strings(cols)
	return cols
}

func writeLinksCSV(w io.Writer, links []Link) error {
	header := []string{"id", "type", "from", "to", "directed", "weight"}
	if err := writeCSVRow(w, header); err != nil {
		return err
	}
	for _, link := range links {
		dir := ""
		if link.Directed != nil {
			dir = fmt.Sprintf("%t", *link.Directed)
		}
		weight := ""
		if link.Weight != nil {
			weight = fmt.Sprintf("%g", *link.Weight)
		}
		if err := writeCSVRow(w, []string{link.ID, link.Type, link.From, link.To, dir, weight}); err != nil {
			return err
		}
	}
	return nil
}

func envelopeFromRecordSet(rs *RecordSet) map[string]any {
	recs := make([]any, len(rs.Records))
	for i, rec := range rs.Records {
		recs[i] = recordJSON(rec)
	}
	links := make([]any, len(rs.Links))
	for i, link := range rs.Links {
		links[i] = linkJSON(link)
	}
	doc := map[string]any{
		"wuji":    "1",
		"kind":    string(ShapeRecordSet),
		"records": recs,
	}
	if optsIncludeMeta(rs.MetaData) {
		doc["meta"] = rs.MetaData
	}
	if len(rs.Schema) > 0 {
		doc["schema"] = rs.Schema
	}
	if len(rs.Links) > 0 {
		doc["links"] = links
	}
	return doc
}

func linkJSON(link Link) map[string]any {
	out := map[string]any{
		"type": link.Type,
		"from": link.From,
		"to":   link.To,
	}
	if link.ID != "" {
		out["id"] = link.ID
	}
	if link.Directed != nil {
		out["directed"] = *link.Directed
	}
	if link.Weight != nil {
		out["weight"] = *link.Weight
	}
	if len(link.Fields) > 0 {
		fields := map[string]any{}
		for k, v := range link.Fields {
			fields[k] = v.JSONAny()
		}
		out["fields"] = fields
	}
	return out
}

func optsIncludeMeta(m Meta) bool {
	return m.Task != "" || m.Capability != "" || m.Driver != "" || m.Model != "" || len(m.Source) > 0
}
