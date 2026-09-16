package driver

import (
	"fmt"
	"strings"

	"github.com/coditary/wuji-core/pkg/data"
)

// DataInputFormat selects how CLI input files are interpreted.
type DataInputFormat string

const (
	DataInputAuto    DataInputFormat = "auto"
	DataInputJSON    DataInputFormat = "json"
	DataInputCSV     DataInputFormat = "csv"
	DataInputMsgpack DataInputFormat = "msgpack"
)

// DataRequest is input for structured data operations.
type DataRequest struct {
	Task        DataTask
	Model       string
	Text        string
	Texts       []string
	TextFile    string
	ImagePath   string
	CSVPath     string
	GraphPath   string
	UseStdin    bool
	InputFormat DataInputFormat
	InputShape  data.ShapeKind
	OutputShape data.ShapeKind
	Input       data.DataShape
	Horizon     int
}

// Validate checks required fields for a data request.
func (r DataRequest) Validate() error {
	if _, err := ParseDataTask(string(r.Task)); err != nil && r.Task != "" {
		return err
	}
	if r.Task == "" && r.Input == nil {
		return fmt.Errorf("either a task flag (--embed, --detect, …) or structured --input is required")
	}
	hasSource := r.Input != nil || strings.TrimSpace(r.Text) != "" || len(r.Texts) > 0 || r.TextFile != "" ||
		r.ImagePath != "" || r.CSVPath != "" || r.GraphPath != "" || r.UseStdin
	if r.Task != "" && !hasSource && r.Task != DataTaskGraph {
		return fmt.Errorf("data task %q requires input (--text, --image, --csv-file, --graph-file, or stdin)", r.Task)
	}
	return nil
}

// TaskOrDefault returns the parsed task or empty for convert-only runs.
func (r DataRequest) TaskOrDefault() DataTask {
	return r.Task.Normalize()
}

// InferInputShape guesses input shape from flags when --input-shape auto.
func InferInputShape(req DataRequest) data.ShapeKind {
	if req.InputShape != "" && req.InputShape != "auto" {
		return req.InputShape
	}
	if req.CSVPath != "" || req.InputFormat == DataInputCSV {
		return data.ShapeTable
	}
	if req.GraphPath != "" {
		return data.ShapeGraph
	}
	switch req.TaskOrDefault() {
	case DataTaskEmbed, DataTaskVector:
		return data.ShapeVector
	case DataTaskClassify, DataTaskRegress, DataTaskForecast:
		return data.ShapeTable
	case DataTaskGraph, DataTaskNodeClass, DataTaskLinkPredict:
		return data.ShapeGraph
	default:
		return data.ShapeRecordSet
	}
}

// InferOutputShape picks output shape from task or explicit flag.
func InferOutputShape(req DataRequest) data.ShapeKind {
	if req.OutputShape != "" && req.OutputShape != "auto" {
		return req.OutputShape
	}
	if req.Task == "" {
		return data.ShapeRecordSet
	}
	switch req.TaskOrDefault().DefaultOutputShape() {
	case "vector":
		return data.ShapeVector
	case "table":
		return data.ShapeTable
	case "graph":
		return data.ShapeGraph
	default:
		return data.ShapeRecordSet
	}
}
