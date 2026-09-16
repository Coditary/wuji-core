package data

// ShapeKind identifies the internal representation of a data payload.
type ShapeKind string

const (
	ShapeRecordSet ShapeKind = "recordset"
	ShapeTable     ShapeKind = "table"
	ShapeVector    ShapeKind = "vector"
	ShapeGraph     ShapeKind = "graph"
)

// ExportFormat selects stdout serialization.
type ExportFormat string

const (
	FormatJSON    ExportFormat = "json"
	FormatNDJSON  ExportFormat = "ndjson"
	FormatCSV     ExportFormat = "csv"
	FormatMsgpack ExportFormat = "msgpack"
)

// CSVView selects which projection to emit for CSV export.
type CSVView string

const (
	CSVViewRecords CSVView = "records"
	CSVViewLinks   CSVView = "links"
	CSVViewNodes   CSVView = "nodes"
)

// VectorCSVMode controls how embeddings appear in CSV output.
type VectorCSVMode string

const (
	VectorCSVJSON   VectorCSVMode = "json"
	VectorCSVExpand VectorCSVMode = "expand"
	VectorCSVDims   VectorCSVMode = "dims"
)

// ExportOpts configures shape export.
type ExportOpts struct {
	Format        ExportFormat
	CSVView       CSVView
	VectorCSV     VectorCSVMode
	IncludeMeta   bool
	PrettyJSON    bool
	CSVFields     []string
	PreviewLength int
}

// DefaultExportOpts returns sensible CLI defaults.
func DefaultExportOpts() ExportOpts {
	return ExportOpts{
		Format:      FormatJSON,
		CSVView:     CSVViewRecords,
		VectorCSV:   VectorCSVDims,
		IncludeMeta: true,
		PrettyJSON:  true,
	}
}
