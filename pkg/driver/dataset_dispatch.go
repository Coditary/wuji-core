package driver

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/capability"
)

func HasDatasetTask(d Driver, task DatasetTask) bool {
	for _, t := range d.Info().DatasetTasks {
		if t == task {
			return true
		}
	}
	return false
}

func SupportsDataset(d Driver) bool {
	return len(d.Info().DatasetTasks) > 0
}

func RunDatasetTask(ctx context.Context, d Driver, req DatasetRequest) (*DatasetResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if !SupportsDataset(d) {
		return nil, fmt.Errorf("driver %q does not support dataset tasks", d.Info().ID)
	}

	task := req.TaskOrDefault()
	if !HasDatasetTask(d, task) {
		return nil, fmt.Errorf("driver %q does not support dataset task %q", d.Info().ID, task)
	}

	switch task {
	case DatasetTaskList:
		mgr, ok := d.(DatasetLister)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement DatasetLister", d.Info().ID)
		}
		return mgr.ListDatasets(ctx, req.ToListRequest())
	case DatasetTaskCreate:
		mgr, ok := d.(DatasetCreator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement DatasetCreator", d.Info().ID)
		}
		return mgr.CreateDataset(ctx, req.ToCreateRequest())
	case DatasetTaskDelete:
		mgr, ok := d.(DatasetDeleter)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement DatasetDeleter", d.Info().ID)
		}
		return mgr.DeleteDataset(ctx, req.ToDeleteRequest())
	case DatasetTaskIngest:
		mgr, ok := d.(DatasetIngester)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement DatasetIngester", d.Info().ID)
		}
		result, err := mgr.IngestDataset(ctx, req.ToIngestRequest())
		if err != nil {
			return nil, err
		}
		return &DatasetResponse{
			Message: result.Message, FilesAdded: result.FilesAdded,
			BytesAdded: result.BytesAdded, VersionID: result.VersionID,
		}, nil
	case DatasetTaskVersion:
		mgr, ok := d.(DatasetVersioner)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement DatasetVersioner", d.Info().ID)
		}
		result, err := mgr.CreateDatasetVersion(ctx, req.ToVersionRequest())
		if err != nil {
			return nil, err
		}
		return &DatasetResponse{
			Message: result.Message, Versions: []DatasetVersionEntry{result.Version},
		}, nil
	case DatasetTaskListVersions:
		mgr, ok := d.(DatasetVersionLister)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement DatasetVersionLister", d.Info().ID)
		}
		versions, err := mgr.ListDatasetVersions(ctx, req.ToListVersionsRequest())
		if err != nil {
			return nil, err
		}
		return &DatasetResponse{Versions: versions}, nil
	default:
		return nil, fmt.Errorf("unsupported dataset task %q", task)
	}
}

func DatasetCapabilitiesFromTasks(tasks []DatasetTask) []capability.Type {
	if len(tasks) == 0 {
		return nil
	}
	return []capability.Type{capability.DatasetMgmt}
}
