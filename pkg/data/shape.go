package data

import "io"

// DataShape is the typed internal representation shared by all data tasks.
type DataShape interface {
	Kind() ShapeKind
	Meta() Meta
	Export(w io.Writer, opts ExportOpts) error
}
