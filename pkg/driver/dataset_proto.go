package driver

import wujiv1 "github.com/coditary/wuji-core/api/proto/v1"

func DatasetTasksToProto(tasks []DatasetTask) []string {
	if len(tasks) == 0 {
		return nil
	}
	out := make([]string, 0, len(tasks))
	for _, task := range tasks {
		out = append(out, string(task))
	}
	return out
}

func DatasetTasksFromProto(items []string) []DatasetTask {
	if len(items) == 0 {
		return nil
	}
	out := make([]DatasetTask, 0, len(items))
	for _, item := range items {
		task := DatasetTask(item)
		if task.IsValid() {
			out = append(out, task)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func datasetEntryFromProto(entry *wujiv1.DatasetEntry) DatasetEntry {
	if entry == nil {
		return DatasetEntry{}
	}
	return DatasetEntry{
		ID: entry.GetId(), Name: entry.GetName(), Path: entry.GetPath(), Size: entry.GetSize(),
		LatestVersion: entry.GetLatestVersion(), VersionCount: int(entry.GetVersionCount()),
	}
}

func datasetEntryToProto(entry DatasetEntry) *wujiv1.DatasetEntry {
	return &wujiv1.DatasetEntry{
		Id: entry.ID, Name: entry.Name, Path: entry.Path, Size: entry.Size,
		LatestVersion: entry.LatestVersion, VersionCount: int32(entry.VersionCount),
	}
}

func datasetEntriesFromProto(entries []*wujiv1.DatasetEntry) []DatasetEntry {
	if len(entries) == 0 {
		return nil
	}
	out := make([]DatasetEntry, 0, len(entries))
	for _, entry := range entries {
		out = append(out, datasetEntryFromProto(entry))
	}
	return out
}

func datasetEntriesToProto(entries []DatasetEntry) []*wujiv1.DatasetEntry {
	if len(entries) == 0 {
		return nil
	}
	out := make([]*wujiv1.DatasetEntry, 0, len(entries))
	for _, entry := range entries {
		out = append(out, datasetEntryToProto(entry))
	}
	return out
}

func datasetVersionFromProto(entry *wujiv1.DatasetVersionEntry) DatasetVersionEntry {
	if entry == nil {
		return DatasetVersionEntry{}
	}
	return DatasetVersionEntry{
		ID: entry.GetId(), DatasetID: entry.GetDatasetId(), Tag: entry.GetTag(), Path: entry.GetPath(),
		Size: entry.GetSize(), CreatedAtUnix: entry.GetCreatedAtUnix(), Message: entry.GetMessage(),
	}
}

func datasetVersionToProto(entry DatasetVersionEntry) *wujiv1.DatasetVersionEntry {
	return &wujiv1.DatasetVersionEntry{
		Id: entry.ID, DatasetId: entry.DatasetID, Tag: entry.Tag, Path: entry.Path,
		Size: entry.Size, CreatedAtUnix: entry.CreatedAtUnix, Message: entry.Message,
	}
}

func datasetVersionsFromProto(entries []*wujiv1.DatasetVersionEntry) []DatasetVersionEntry {
	if len(entries) == 0 {
		return nil
	}
	out := make([]DatasetVersionEntry, 0, len(entries))
	for _, entry := range entries {
		out = append(out, datasetVersionFromProto(entry))
	}
	return out
}

func datasetVersionsToProto(entries []DatasetVersionEntry) []*wujiv1.DatasetVersionEntry {
	if len(entries) == 0 {
		return nil
	}
	out := make([]*wujiv1.DatasetVersionEntry, 0, len(entries))
	for _, entry := range entries {
		out = append(out, datasetVersionToProto(entry))
	}
	return out
}

func DatasetCreateRequestFromProto(req *wujiv1.CreateDatasetRequest) DatasetCreateRequest {
	if req == nil {
		return DatasetCreateRequest{}
	}
	return DatasetCreateRequest{Name: req.GetName(), Path: req.GetPath(), Description: req.GetDescription()}
}

func DatasetCreateRequestToProto(req DatasetCreateRequest) *wujiv1.CreateDatasetRequest {
	return &wujiv1.CreateDatasetRequest{Name: req.Name, Path: req.Path, Description: req.Description}
}

func DatasetDeleteRequestFromProto(req *wujiv1.DeleteDatasetRequest) DatasetDeleteRequest {
	if req == nil {
		return DatasetDeleteRequest{}
	}
	return DatasetDeleteRequest{Name: req.GetName(), ID: req.GetId()}
}

func DatasetDeleteRequestToProto(req DatasetDeleteRequest) *wujiv1.DeleteDatasetRequest {
	return &wujiv1.DeleteDatasetRequest{Name: req.Name, Id: req.ID}
}

func DatasetIngestRequestFromProto(req *wujiv1.IngestDatasetRequest) DatasetIngestRequest {
	if req == nil {
		return DatasetIngestRequest{}
	}
	return DatasetIngestRequest{
		DatasetID: req.GetDatasetId(), Name: req.GetName(), SourcePath: req.GetSourcePath(),
		Recursive: req.GetRecursive(), Format: req.GetFormat(),
	}
}

func DatasetIngestRequestToProto(req DatasetIngestRequest) *wujiv1.IngestDatasetRequest {
	return &wujiv1.IngestDatasetRequest{
		DatasetId: req.DatasetID, Name: req.Name, SourcePath: req.SourcePath,
		Recursive: req.Recursive, Format: req.Format,
	}
}

func DatasetVersionRequestFromProto(req *wujiv1.CreateDatasetVersionRequest) DatasetVersionRequest {
	if req == nil {
		return DatasetVersionRequest{}
	}
	return DatasetVersionRequest{
		DatasetID: req.GetDatasetId(), Name: req.GetName(), Tag: req.GetTag(), Message: req.GetMessage(),
	}
}

func DatasetVersionRequestToProto(req DatasetVersionRequest) *wujiv1.CreateDatasetVersionRequest {
	return &wujiv1.CreateDatasetVersionRequest{
		DatasetId: req.DatasetID, Name: req.Name, Tag: req.Tag, Message: req.Message,
	}
}

func DatasetListVersionsRequestFromProto(req *wujiv1.ListDatasetVersionsRequest) DatasetListVersionsRequest {
	if req == nil {
		return DatasetListVersionsRequest{}
	}
	return DatasetListVersionsRequest{DatasetID: req.GetDatasetId(), Name: req.GetName()}
}

func DatasetListVersionsRequestToProto(req DatasetListVersionsRequest) *wujiv1.ListDatasetVersionsRequest {
	return &wujiv1.ListDatasetVersionsRequest{DatasetId: req.DatasetID, Name: req.Name}
}

func DatasetIngestResultFromProto(resp *wujiv1.IngestDatasetResponse) *DatasetIngestResult {
	if resp == nil {
		return &DatasetIngestResult{}
	}
	return &DatasetIngestResult{
		Message: resp.GetMessage(), FilesAdded: int(resp.GetFilesAdded()),
		BytesAdded: resp.GetBytesAdded(), VersionID: resp.GetVersionId(),
	}
}

func DatasetIngestResultToProto(result *DatasetIngestResult) *wujiv1.IngestDatasetResponse {
	if result == nil {
		return &wujiv1.IngestDatasetResponse{}
	}
	return &wujiv1.IngestDatasetResponse{
		Message: result.Message, FilesAdded: int32(result.FilesAdded),
		BytesAdded: result.BytesAdded, VersionId: result.VersionID,
	}
}

func DatasetVersionResultFromProto(resp *wujiv1.DatasetVersionResponse) *DatasetVersionResult {
	if resp == nil {
		return &DatasetVersionResult{}
	}
	return &DatasetVersionResult{
		Version: datasetVersionFromProto(resp.GetVersion()), Message: resp.GetMessage(),
	}
}

func DatasetVersionResultToProto(result *DatasetVersionResult) *wujiv1.DatasetVersionResponse {
	if result == nil {
		return &wujiv1.DatasetVersionResponse{}
	}
	return &wujiv1.DatasetVersionResponse{
		Version: datasetVersionToProto(result.Version), Message: result.Message,
	}
}

func DatasetMutationResponseFromProto(resp *wujiv1.DatasetMutationResponse) *DatasetResponse {
	if resp == nil {
		return &DatasetResponse{}
	}
	out := &DatasetResponse{Message: resp.GetMessage()}
	if resp.GetDataset() != nil {
		out.Datasets = []DatasetEntry{datasetEntryFromProto(resp.GetDataset())}
	}
	return out
}

func DatasetMutationResponseToProto(resp *DatasetResponse) *wujiv1.DatasetMutationResponse {
	if resp == nil {
		return &wujiv1.DatasetMutationResponse{}
	}
	out := &wujiv1.DatasetMutationResponse{Message: resp.Message}
	if len(resp.Datasets) > 0 {
		out.Dataset = datasetEntryToProto(resp.Datasets[0])
	}
	return out
}

func ListDatasetsResponseFromProto(resp *wujiv1.ListDatasetsResponse) *DatasetResponse {
	if resp == nil {
		return &DatasetResponse{}
	}
	return &DatasetResponse{Datasets: datasetEntriesFromProto(resp.GetDatasets())}
}

func ListDatasetsResponseToProto(resp *DatasetResponse) *wujiv1.ListDatasetsResponse {
	if resp == nil {
		return &wujiv1.ListDatasetsResponse{}
	}
	return &wujiv1.ListDatasetsResponse{Datasets: datasetEntriesToProto(resp.Datasets)}
}

func ListDatasetVersionsResponseFromProto(resp *wujiv1.ListDatasetVersionsResponse) *DatasetResponse {
	if resp == nil {
		return &DatasetResponse{}
	}
	return &DatasetResponse{Versions: datasetVersionsFromProto(resp.GetVersions())}
}

func ListDatasetVersionsResponseToProto(resp *DatasetResponse) *wujiv1.ListDatasetVersionsResponse {
	if resp == nil {
		return &wujiv1.ListDatasetVersionsResponse{}
	}
	return &wujiv1.ListDatasetVersionsResponse{Versions: datasetVersionsToProto(resp.Versions)}
}
