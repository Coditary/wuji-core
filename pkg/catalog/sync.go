package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// SyncResult describes a catalog sync operation.
type SyncResult struct {
	Source         string
	SyncedAt       time.Time
	ProviderCount  int
	ModelCount     int
	LabCount       int
	ProviderModels int
	IndexBytes     int
	Skipped        bool
	Message        string
}

// Sync downloads models.dev catalog.json and writes the slim index to disk.
func Sync(ctx context.Context, root, source string) (*SyncResult, error) {
	data, etag, err := Fetch(ctx, source)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	idx, meta, err := BuildIndex(data, source, now)
	if err != nil {
		return nil, err
	}
	meta.ETag = etag
	if err := Save(root, idx, meta); err != nil {
		return nil, err
	}

	indexData, _ := json.Marshal(idx)
	return &SyncResult{
		Source:         meta.Source,
		SyncedAt:       now,
		ProviderCount:  meta.ProviderCount,
		ModelCount:     meta.ModelCount,
		LabCount:       meta.LabCount,
		ProviderModels: meta.ProviderModels,
		IndexBytes:     len(indexData),
		Message:        fmt.Sprintf("synced %d providers, %d canonical models, %d labs", meta.ProviderCount, meta.ModelCount, meta.LabCount),
	}, nil
}
