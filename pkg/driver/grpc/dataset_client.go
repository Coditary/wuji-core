package grpc

import (
	"context"
	"fmt"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *RemoteDriver) ListDatasets(ctx context.Context, req driver.DatasetListRequest) (*driver.DatasetResponse, error) {
	resp, err := d.client.ListDatasets(ctx, &wujiv1.ListDatasetsRequest{})
	if err != nil {
		return nil, fmt.Errorf("remote list datasets: %w", err)
	}
	return driver.ListDatasetsResponseFromProto(resp), nil
}

func (d *RemoteDriver) CreateDataset(ctx context.Context, req driver.DatasetCreateRequest) (*driver.DatasetResponse, error) {
	resp, err := d.client.CreateDataset(ctx, driver.DatasetCreateRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote create dataset: %w", err)
	}
	return driver.DatasetMutationResponseFromProto(resp), nil
}

func (d *RemoteDriver) DeleteDataset(ctx context.Context, req driver.DatasetDeleteRequest) (*driver.DatasetResponse, error) {
	resp, err := d.client.DeleteDataset(ctx, driver.DatasetDeleteRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote delete dataset: %w", err)
	}
	return driver.DatasetMutationResponseFromProto(resp), nil
}

func (d *RemoteDriver) IngestDataset(ctx context.Context, req driver.DatasetIngestRequest) (*driver.DatasetIngestResult, error) {
	resp, err := d.client.IngestDataset(ctx, driver.DatasetIngestRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote ingest dataset: %w", err)
	}
	return driver.DatasetIngestResultFromProto(resp), nil
}

func (d *RemoteDriver) CreateDatasetVersion(ctx context.Context, req driver.DatasetVersionRequest) (*driver.DatasetVersionResult, error) {
	resp, err := d.client.CreateDatasetVersion(ctx, driver.DatasetVersionRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote create dataset version: %w", err)
	}
	return driver.DatasetVersionResultFromProto(resp), nil
}

func (d *RemoteDriver) ListDatasetVersions(ctx context.Context, req driver.DatasetListVersionsRequest) ([]driver.DatasetVersionEntry, error) {
	resp, err := d.client.ListDatasetVersions(ctx, driver.DatasetListVersionsRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote list dataset versions: %w", err)
	}
	return driver.ListDatasetVersionsResponseFromProto(resp).Versions, nil
}
