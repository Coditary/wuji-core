package driver

import (
	"bytes"
	"fmt"
	"strings"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/data"
)

const (
	DataPayloadMsgpack = "msgpack"
	DataPayloadJSON    = "json"
)

func DataTasksToProto(tasks []DataTask) []string {
	if len(tasks) == 0 {
		return nil
	}
	out := make([]string, 0, len(tasks))
	for _, task := range tasks {
		out = append(out, string(task))
	}
	return out
}

func DataTasksFromProto(items []string) []DataTask {
	if len(items) == 0 {
		return nil
	}
	out := make([]DataTask, 0, len(items))
	for _, item := range items {
		task, err := ParseDataTask(item)
		if err != nil {
			continue
		}
		out = append(out, task)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func ProduceDataRequestToProto(req DataRequest) (*wujiv1.ProduceDataRequest, error) {
	out := &wujiv1.ProduceDataRequest{
		Task:        string(req.Task),
		Model:       req.Model,
		Text:        req.Text,
		Texts:       append([]string(nil), req.Texts...),
		TextFile:    req.TextFile,
		ImagePath:  req.ImagePath,
		CsvPath:    req.CSVPath,
		GraphPath:  req.GraphPath,
		InputFormat: string(req.InputFormat),
		InputShape:  string(req.InputShape),
		OutputShape: string(req.OutputShape),
		Horizon:     int32(req.Horizon),
	}
	if req.Input != nil {
		payload, err := data.MarshalMsgpack(req.Input)
		if err != nil {
			return nil, fmt.Errorf("encode input payload: %w", err)
		}
		out.InputPayload = payload
		out.InputPayloadFormat = DataPayloadMsgpack
	}
	return out, nil
}

func ProduceDataRequestFromProto(p *wujiv1.ProduceDataRequest) (DataRequest, error) {
	if p == nil {
		return DataRequest{}, nil
	}
	req := DataRequest{
		Model:       p.GetModel(),
		Text:        p.GetText(),
		Texts:       append([]string(nil), p.GetTexts()...),
		TextFile:    p.GetTextFile(),
		ImagePath:   p.GetImagePath(),
		CSVPath:     p.GetCsvPath(),
		GraphPath:   p.GetGraphPath(),
		InputFormat: DataInputFormat(p.GetInputFormat()),
		InputShape:  data.ShapeKind(p.GetInputShape()),
		OutputShape: data.ShapeKind(p.GetOutputShape()),
		Horizon:     int(p.GetHorizon()),
	}
	if task := strings.TrimSpace(p.GetTask()); task != "" {
		parsed, err := ParseDataTask(task)
		if err != nil {
			return DataRequest{}, err
		}
		req.Task = parsed
	}
	if len(p.GetInputPayload()) > 0 {
		format := p.GetInputPayloadFormat()
		if format == "" {
			format = DataPayloadMsgpack
		}
		shape, err := dataShapeFromPayload(p.GetInputPayload(), format)
		if err != nil {
			return DataRequest{}, err
		}
		req.Input = shape
	}
	return req, nil
}

func ProduceDataResponseToProto(shape data.DataShape) (*wujiv1.ProduceDataResponse, error) {
	if shape == nil {
		return &wujiv1.ProduceDataResponse{}, nil
	}
	payload, err := data.MarshalMsgpack(shape)
	if err != nil {
		return nil, fmt.Errorf("encode data payload: %w", err)
	}
	return &wujiv1.ProduceDataResponse{
		Payload:       payload,
		PayloadFormat: DataPayloadMsgpack,
		Kind:          string(shape.Kind()),
	}, nil
}

func ProduceDataResponseFromProto(p *wujiv1.ProduceDataResponse) (data.DataShape, error) {
	if p == nil {
		return nil, fmt.Errorf("empty produce data response")
	}
	format := p.GetPayloadFormat()
	if format == "" {
		format = DataPayloadMsgpack
	}
	return dataShapeFromPayload(p.GetPayload(), format)
}

func dataShapeFromPayload(payload []byte, format string) (data.DataShape, error) {
	switch strings.ToLower(format) {
	case DataPayloadJSON:
		return data.ParseJSON(bytes.NewReader(payload))
	case DataPayloadMsgpack, "":
		return data.UnmarshalMsgpack(payload)
	default:
		return nil, fmt.Errorf("unsupported data payload format %q", format)
	}
}
