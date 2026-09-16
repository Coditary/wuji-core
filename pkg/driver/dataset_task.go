package driver

import (
	"fmt"
	"strings"
)

// DatasetTask identifies a dataset management operation.
type DatasetTask string

const (
	DatasetTaskList          DatasetTask = "list"
	DatasetTaskCreate        DatasetTask = "create"
	DatasetTaskDelete        DatasetTask = "delete"
	DatasetTaskIngest        DatasetTask = "ingest"
	DatasetTaskVersion       DatasetTask = "version"
	DatasetTaskListVersions  DatasetTask = "list-versions"
)

// DatasetAction is a deprecated alias for DatasetTask.
type DatasetAction = DatasetTask

const (
	DatasetList         = DatasetTaskList
	DatasetCreate       = DatasetTaskCreate
	DatasetDelete       = DatasetTaskDelete
	DatasetIngest       = DatasetTaskIngest
	DatasetVersion      = DatasetTaskVersion
	DatasetListVersions = DatasetTaskListVersions
)

type DatasetTaskInfo struct {
	Task                DatasetTask
	Description         string
	RequiresName        bool
	RequiresPath        bool
	RequiresSourcePath  bool
	RequiresDatasetRef  bool
}

func AllDatasetTasks() []DatasetTask {
	return []DatasetTask{
		DatasetTaskList,
		DatasetTaskCreate,
		DatasetTaskDelete,
		DatasetTaskIngest,
		DatasetTaskVersion,
		DatasetTaskListVersions,
	}
}

func DatasetTaskCatalog() []DatasetTaskInfo {
	return []DatasetTaskInfo{
		{Task: DatasetTaskList, Description: "List registered datasets"},
		{Task: DatasetTaskCreate, Description: "Create a new dataset", RequiresName: true, RequiresPath: true},
		{Task: DatasetTaskDelete, Description: "Delete a dataset", RequiresName: true},
		{Task: DatasetTaskIngest, Description: "Import files into a dataset", RequiresDatasetRef: true, RequiresSourcePath: true},
		{Task: DatasetTaskVersion, Description: "Create an immutable dataset snapshot", RequiresDatasetRef: true},
		{Task: DatasetTaskListVersions, Description: "List versions of a dataset", RequiresDatasetRef: true},
	}
}

func (t DatasetTask) String() string { return string(t) }

func (t DatasetTask) IsValid() bool {
	for _, known := range AllDatasetTasks() {
		if t == known {
			return true
		}
	}
	return false
}

func ParseDatasetTask(raw string) (DatasetTask, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	task := DatasetTask(normalized)
	if !task.IsValid() {
		return "", fmt.Errorf("unknown dataset task %q (valid: %s)", raw, joinDatasetTasks())
	}
	return task, nil
}

func (r DatasetRequest) TaskOrDefault() DatasetTask {
	if r.Task != "" {
		return r.Task
	}
	if r.Action != "" {
		return DatasetTask(r.Action)
	}
	return DatasetTaskList
}

func (r DatasetRequest) datasetRef() bool {
	return r.DatasetID != "" || r.Name != ""
}

func (r DatasetRequest) Validate() error {
	task := r.TaskOrDefault()
	if !task.IsValid() {
		return fmt.Errorf("unknown dataset task %q (valid: %s)", task, joinDatasetTasks())
	}
	info := datasetTaskInfo(task)
	if info.RequiresName && r.Name == "" {
		return fmt.Errorf("dataset name is required for task %q", task)
	}
	if info.RequiresPath && r.Path == "" {
		return fmt.Errorf("dataset path is required for task %q", task)
	}
	if info.RequiresSourcePath && r.SourcePath == "" {
		return fmt.Errorf("source path is required for task %q", task)
	}
	if info.RequiresDatasetRef && !r.datasetRef() {
		return fmt.Errorf("dataset name or id is required for task %q", task)
	}
	return nil
}

func datasetTaskInfo(task DatasetTask) DatasetTaskInfo {
	for _, info := range DatasetTaskCatalog() {
		if info.Task == task {
			return info
		}
	}
	return DatasetTaskInfo{}
}

func joinDatasetTasks() string {
	parts := make([]string, len(AllDatasetTasks()))
	for i, t := range AllDatasetTasks() {
		parts[i] = string(t)
	}
	return strings.Join(parts, ", ")
}
