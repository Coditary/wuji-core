package data

import (
	"io"
)

// Graph combines node records and edge links.
type Graph struct {
	MetaData Meta
	Schema   []FieldDef
	Nodes    []Record
	Edges    []Link
}

func (g *Graph) Kind() ShapeKind { return ShapeGraph }
func (g *Graph) Meta() Meta      { return g.MetaData }

func (g *Graph) Export(w io.Writer, opts ExportOpts) error {
	switch opts.Format {
	case FormatCSV:
		switch opts.CSVView {
		case CSVViewLinks:
			return writeLinksCSV(w, g.Edges)
		case CSVViewNodes:
			return writeRecordsCSV(w, g.Nodes, opts)
		default:
			return writeLinksCSV(w, g.Edges)
		}
	case FormatMsgpack:
		return WriteMsgpack(w, g)
	default:
		rs := g.ToRecordSet()
		return rs.Export(w, opts)
	}
}

// ToRecordSet flattens nodes and edges into a WDD envelope via RecordSet.
func (g *Graph) ToRecordSet() *RecordSet {
	nodes := make([]Record, len(g.Nodes))
	copy(nodes, g.Nodes)
	for i := range nodes {
		if nodes[i].Role == "" {
			nodes[i].Role = "node"
		}
	}
	links := make([]Link, len(g.Edges))
	copy(links, g.Edges)
	return &RecordSet{
		MetaData: g.MetaData,
		Schema:   g.Schema,
		Records:  nodes,
		Links:    links,
	}
}

// FromRecordSet extracts a graph when records have role node and links are present.
func FromRecordSetGraph(rs *RecordSet) *Graph {
	nodes := make([]Record, 0, len(rs.Records))
	for _, rec := range rs.Records {
		if rec.Role == "" || rec.Role == "node" {
			nodes = append(nodes, rec)
		}
	}
	return &Graph{
		MetaData: rs.MetaData,
		Schema:   rs.Schema,
		Nodes:    nodes,
		Edges:    rs.Links,
	}
}
