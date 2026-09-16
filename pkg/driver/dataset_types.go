package driver

// DatasetListRequest lists datasets.
type DatasetListRequest struct{}

// DatasetCreateRequest creates a dataset.
type DatasetCreateRequest struct {
	Name        string
	Path        string
	Description string
}

// DatasetDeleteRequest deletes a dataset.
type DatasetDeleteRequest struct {
	Name string
	ID   string
}

// DatasetIngestRequest imports files into a dataset.
type DatasetIngestRequest struct {
	DatasetID  string
	Name       string
	SourcePath string
	Recursive  bool
	Format     string
}

// DatasetVersionRequest creates a dataset snapshot.
type DatasetVersionRequest struct {
	DatasetID string
	Name      string
	Tag       string
	Message   string
}

// DatasetListVersionsRequest lists dataset versions.
type DatasetListVersionsRequest struct {
	DatasetID string
	Name      string
}

// DatasetVersionEntry describes an immutable dataset snapshot.
type DatasetVersionEntry struct {
	ID            string
	DatasetID     string
	Tag           string
	Path          string
	Size          int64
	CreatedAtUnix int64
	Message       string
}

// DatasetIngestResult is the output of dataset ingest.
type DatasetIngestResult struct {
	Message    string
	FilesAdded int
	BytesAdded int64
	VersionID  string
}

// DatasetVersionResult is the output of version creation.
type DatasetVersionResult struct {
	Version DatasetVersionEntry
	Message string
}

func (r DatasetRequest) ToListRequest() DatasetListRequest { return DatasetListRequest{} }

func (r DatasetRequest) ToCreateRequest() DatasetCreateRequest {
	return DatasetCreateRequest{Name: r.Name, Path: r.Path, Description: r.Description}
}

func (r DatasetRequest) ToDeleteRequest() DatasetDeleteRequest {
	return DatasetDeleteRequest{Name: r.Name, ID: r.DatasetID}
}

func (r DatasetRequest) ToIngestRequest() DatasetIngestRequest {
	return DatasetIngestRequest{
		DatasetID: r.DatasetID, Name: r.Name, SourcePath: r.SourcePath,
		Recursive: r.Recursive, Format: r.Format,
	}
}

func (r DatasetRequest) ToVersionRequest() DatasetVersionRequest {
	return DatasetVersionRequest{
		DatasetID: r.DatasetID, Name: r.Name, Tag: r.Tag, Message: r.Message,
	}
}

func (r DatasetRequest) ToListVersionsRequest() DatasetListVersionsRequest {
	return DatasetListVersionsRequest{DatasetID: r.DatasetID, Name: r.Name}
}
