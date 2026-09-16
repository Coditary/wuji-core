package dummy

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *Driver) ProduceData(_ context.Context, req driver.DataRequest) (data.DataShape, error) {
	task := req.TaskOrDefault()
	meta := data.Meta{
		Capability: "data",
		Task:       string(task),
		Driver:     DriverID,
		Model:      req.Model,
		Source: map[string]string{
			"text":  req.Text,
			"image": req.ImagePath,
			"csv":   req.CSVPath,
			"graph": req.GraphPath,
		},
	}

	switch task {
	case driver.DataTaskEmbed, driver.DataTaskVector:
		return dummyEmbed(req, meta)
	case driver.DataTaskDetect:
		return dummyDetect(req, meta)
	case driver.DataTaskSegment:
		return dummySegment(req, meta)
	case driver.DataTaskNER:
		return dummyNER(req, meta)
	case driver.DataTaskKeypoint:
		return dummyKeypoints(req, meta)
	case driver.DataTaskClassify, driver.DataTaskRegress:
		return dummyTabular(req, meta, task)
	case driver.DataTaskForecast:
		return dummyForecast(req, meta)
	case driver.DataTaskGraph, driver.DataTaskNodeClass, driver.DataTaskLinkPredict:
		return dummyGraph(req, meta)
	default:
		return nil, fmt.Errorf("dummy driver: unsupported data task %q", task)
	}
}

func dummyEmbed(req driver.DataRequest, meta data.Meta) (*data.VectorBatch, error) {
	inputs := append([]string(nil), req.Texts...)
	if t := strings.TrimSpace(req.Text); t != "" {
		inputs = append(inputs, t)
	}
	if req.ImagePath != "" {
		inputs = append(inputs, "image:"+req.ImagePath)
	}
	if len(inputs) == 0 {
		return nil, fmt.Errorf("embed requires --text, batch texts, or --image")
	}
	const dims = 8
	rows := make([][]float32, len(inputs))
	for i, in := range inputs {
		rows[i] = pseudoVector(in, dims)
	}
	batch, err := data.NewVectorBatch(inputs, dims, rows)
	if err != nil {
		return nil, err
	}
	batch.MetaData = meta
	return batch, nil
}

func pseudoVector(seed string, dims int) []float32 {
	out := make([]float32, dims)
	h := fnv.New32a()
	_, _ = h.Write([]byte(seed))
	base := float32(h.Sum32()%1000) / 1000
	for i := range out {
		out[i] = base + float32(i)*0.01
	}
	return out
}

func dummyDetect(req driver.DataRequest, meta data.Meta) (*data.RecordSet, error) {
	if req.ImagePath == "" {
		return nil, fmt.Errorf("detect requires --image")
	}
	rs := &data.RecordSet{
		MetaData: meta,
		Schema: []data.FieldDef{
			{Name: "label", Type: data.KindString},
			{Name: "score", Type: data.KindFloat},
			{Name: "bbox", Type: data.KindBBox},
		},
		Records: []data.Record{{
			ID: "d0", Role: "detection",
			Fields: map[string]data.Value{
				"label": data.NewString("object"),
				"score": data.NewFloat(0.91),
				"bbox":  data.NewBBox(data.BBox{X: 10, Y: 20, W: 100, H: 80}),
			},
		}},
	}
	return rs, nil
}

func dummySegment(req driver.DataRequest, meta data.Meta) (*data.RecordSet, error) {
	if req.ImagePath == "" {
		return nil, fmt.Errorf("segment requires --image")
	}
	return &data.RecordSet{
		MetaData: meta,
		Records: []data.Record{{
			ID: "s0", Role: "segment",
			Fields: map[string]data.Value{
				"label": data.NewString("region"),
				"score": data.NewFloat(0.88),
				"bbox":  data.NewBBox(data.BBox{X: 0, Y: 0, W: 640, H: 480, Normalized: true}),
			},
		}},
	}, nil
}

func dummyNER(req driver.DataRequest, meta data.Meta) (*data.RecordSet, error) {
	text := strings.TrimSpace(req.Text)
	if text == "" {
		return nil, fmt.Errorf("ner requires text input")
	}
	word := text
	if idx := strings.IndexByte(text, ' '); idx > 0 {
		word = text[:idx]
	}
	return &data.RecordSet{
		MetaData: meta,
		Records: []data.Record{{
			ID: "e0", Role: "entity",
			Fields: map[string]data.Value{
				"label": data.NewString("MISC"),
				"text":  data.NewString(word),
				"span":  data.NewSpan(data.Span{Start: 0, End: len(word), Text: word}),
			},
		}},
	}, nil
}

func dummyKeypoints(req driver.DataRequest, meta data.Meta) (*data.RecordSet, error) {
	if req.ImagePath == "" {
		return nil, fmt.Errorf("keypoint requires --image")
	}
	return &data.RecordSet{
		MetaData: meta,
		Records: []data.Record{{
			ID: "k0", Role: "keypoint",
			Fields: map[string]data.Value{
				"label": data.NewString("joint"),
				"x":     data.NewFloat(120),
				"y":     data.NewFloat(240),
				"score": data.NewFloat(0.77),
			},
		}},
	}, nil
}

func dummyTabular(req driver.DataRequest, meta data.Meta, task driver.DataTask) (data.DataShape, error) {
	var table *data.Table
	var err error
	if req.Input != nil {
		table, _ = data.AsTable(req.Input)
	}
	if table == nil && req.CSVPath != "" {
		table, err = data.ParseCSVFile(req.CSVPath)
		if err != nil {
			return nil, err
		}
	}
	if table == nil {
		return nil, fmt.Errorf("%s requires --csv-file input", task)
	}
	out := *table
	out.MetaData = meta
	colName := "prediction"
	if task == driver.DataTaskRegress {
		colName = "prediction_value"
	}
	vals := make([]data.Value, out.RowCount())
	for i := range vals {
		if task == driver.DataTaskRegress {
			vals[i] = data.NewFloat(float64(i) + 0.5)
		} else {
			vals[i] = data.NewString("positive")
		}
	}
	out.Columns = append(out.Columns, data.Column{Name: colName, Kind: kindForTask(task), Cells: vals})
	return &out, nil
}

func dummyForecast(req driver.DataRequest, meta data.Meta) (data.DataShape, error) {
	horizon := req.Horizon
	if horizon <= 0 {
		horizon = 7
	}
	cells := make([]data.Value, horizon)
	for i := range cells {
		cells[i] = data.NewFloat(float64(100 + i))
	}
	table := &data.Table{
		MetaData: meta,
		Schema:   []data.FieldDef{{Name: "value", Type: data.KindFloat}, {Name: "kind", Type: data.KindString}},
		Columns: []data.Column{
			{Name: "value", Kind: data.KindFloat, Cells: cells},
			{Name: "kind", Kind: data.KindString, Cells: repeatString("predicted", horizon)},
		},
	}
	return table, nil
}

func repeatString(s string, n int) []data.Value {
	out := make([]data.Value, n)
	for i := range out {
		out[i] = data.NewString(s)
	}
	return out
}

func dummyGraph(req driver.DataRequest, meta data.Meta) (*data.Graph, error) {
	if req.Input != nil {
		if g, ok := req.Input.(*data.Graph); ok {
			g.MetaData = meta
			return g, nil
		}
		if rs, ok := req.Input.(*data.RecordSet); ok {
			g := data.FromRecordSetGraph(rs)
			g.MetaData = meta
			return g, nil
		}
	}
	dir := true
	return &data.Graph{
		MetaData: meta,
		Nodes: []data.Record{
			{ID: "n1", Role: "node", Fields: map[string]data.Value{"label": data.NewString("A"), "score": data.NewFloat(0.9)}},
			{ID: "n2", Role: "node", Fields: map[string]data.Value{"label": data.NewString("B"), "score": data.NewFloat(0.4)}},
		},
		Edges: []data.Link{
			{ID: "e1", Type: "edge", From: "n1", To: "n2", Directed: &dir, Weight: fp(0.8)},
		},
	}, nil
}

func fp(v float64) *float64 { return &v }

func kindForTask(task driver.DataTask) data.ValueKind {
	if task == driver.DataTaskRegress {
		return data.KindFloat
	}
	return data.KindString
}
