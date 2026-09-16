package driver

import (
	"fmt"
	"strings"
)

// DataTask identifies a structured-data operation.
type DataTask string

const (
	DataTaskEmbed       DataTask = "embed"
	DataTaskVector      DataTask = "vector"
	DataTaskDetect      DataTask = "detect"
	DataTaskSegment     DataTask = "segment"
	DataTaskNER         DataTask = "ner"
	DataTaskKeypoint    DataTask = "keypoint"
	DataTaskClassify    DataTask = "classify"
	DataTaskRegress     DataTask = "regress"
	DataTaskForecast    DataTask = "forecast"
	DataTaskGraph       DataTask = "graph"
	DataTaskNodeClass   DataTask = "node_classify"
	DataTaskLinkPredict DataTask = "link_predict"
)

// AllDataTasks lists supported task names for driver metadata.
func AllDataTasks() []DataTask {
	return []DataTask{
		DataTaskEmbed, DataTaskVector, DataTaskDetect, DataTaskSegment, DataTaskNER,
		DataTaskKeypoint, DataTaskClassify, DataTaskRegress, DataTaskForecast,
		DataTaskGraph, DataTaskNodeClass, DataTaskLinkPredict,
	}
}

// DataTaskInfo describes one task for help text and validation.
type DataTaskInfo struct {
	Task        DataTask
	Description string
	InputShape  string
	OutputShape string
}

// AllDataTaskInfo returns metadata for every task.
func AllDataTaskInfo() []DataTaskInfo {
	return []DataTaskInfo{
		{Task: DataTaskEmbed, Description: "Text/image embeddings", InputShape: "text|image", OutputShape: "vector"},
		{Task: DataTaskVector, Description: "Alias for embed", InputShape: "text|image", OutputShape: "vector"},
		{Task: DataTaskDetect, Description: "Object detection", InputShape: "image", OutputShape: "recordset"},
		{Task: DataTaskSegment, Description: "Image segmentation", InputShape: "image", OutputShape: "recordset"},
		{Task: DataTaskNER, Description: "Named entity recognition", InputShape: "text", OutputShape: "recordset"},
		{Task: DataTaskKeypoint, Description: "Keypoint detection", InputShape: "image", OutputShape: "recordset"},
		{Task: DataTaskClassify, Description: "Tabular classification", InputShape: "table", OutputShape: "table"},
		{Task: DataTaskRegress, Description: "Tabular regression", InputShape: "table", OutputShape: "table"},
		{Task: DataTaskForecast, Description: "Time series forecast", InputShape: "table", OutputShape: "table"},
		{Task: DataTaskGraph, Description: "Graph inference", InputShape: "graph", OutputShape: "graph"},
		{Task: DataTaskNodeClass, Description: "Node classification", InputShape: "graph", OutputShape: "graph"},
		{Task: DataTaskLinkPredict, Description: "Link prediction", InputShape: "graph", OutputShape: "graph"},
	}
}

// ParseDataTask normalizes a task name.
func ParseDataTask(name string) (DataTask, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return "", fmt.Errorf("task is required (embed, detect, forecast, graph, …)")
	}
	for _, t := range AllDataTasks() {
		if string(t) == name {
			return t, nil
		}
	}
	return "", fmt.Errorf("unknown data task %q", name)
}

// DefaultOutputShape returns the natural shape for a task.
func (t DataTask) DefaultOutputShape() string {
	switch t {
	case DataTaskEmbed, DataTaskVector:
		return "vector"
	case DataTaskClassify, DataTaskRegress, DataTaskForecast:
		return "table"
	case DataTaskGraph, DataTaskNodeClass, DataTaskLinkPredict:
		return "graph"
	default:
		return "recordset"
	}
}

// Normalize aliases (vector → embed internally).
func (t DataTask) Normalize() DataTask {
	if t == DataTaskVector {
		return DataTaskEmbed
	}
	return t
}
