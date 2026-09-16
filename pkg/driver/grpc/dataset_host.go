package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (s *Host) ListDatasets(ctx context.Context, _ *wujiv1.ListDatasetsRequest) (*wujiv1.ListDatasetsResponse, error) {
	mgr, ok := s.drv.(driver.DatasetLister)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "list datasets not supported")
	}
	resp, err := mgr.ListDatasets(ctx, driver.DatasetListRequest{})
	if err != nil {
		return nil, err
	}
	return driver.ListDatasetsResponseToProto(resp), nil
}

func (s *Host) CreateDataset(ctx context.Context, req *wujiv1.CreateDatasetRequest) (*wujiv1.DatasetMutationResponse, error) {
	mgr, ok := s.drv.(driver.DatasetCreator)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "create dataset not supported")
	}
	resp, err := mgr.CreateDataset(ctx, driver.DatasetCreateRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.DatasetMutationResponseToProto(resp), nil
}

func (s *Host) DeleteDataset(ctx context.Context, req *wujiv1.DeleteDatasetRequest) (*wujiv1.DatasetMutationResponse, error) {
	mgr, ok := s.drv.(driver.DatasetDeleter)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "delete dataset not supported")
	}
	resp, err := mgr.DeleteDataset(ctx, driver.DatasetDeleteRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.DatasetMutationResponseToProto(resp), nil
}

func (s *Host) IngestDataset(ctx context.Context, req *wujiv1.IngestDatasetRequest) (*wujiv1.IngestDatasetResponse, error) {
	mgr, ok := s.drv.(driver.DatasetIngester)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "ingest dataset not supported")
	}
	result, err := mgr.IngestDataset(ctx, driver.DatasetIngestRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.DatasetIngestResultToProto(result), nil
}

func (s *Host) CreateDatasetVersion(ctx context.Context, req *wujiv1.CreateDatasetVersionRequest) (*wujiv1.DatasetVersionResponse, error) {
	mgr, ok := s.drv.(driver.DatasetVersioner)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "create dataset version not supported")
	}
	result, err := mgr.CreateDatasetVersion(ctx, driver.DatasetVersionRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.DatasetVersionResultToProto(result), nil
}

func (s *Host) ListDatasetVersions(ctx context.Context, req *wujiv1.ListDatasetVersionsRequest) (*wujiv1.ListDatasetVersionsResponse, error) {
	mgr, ok := s.drv.(driver.DatasetVersionLister)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "list dataset versions not supported")
	}
	versions, err := mgr.ListDatasetVersions(ctx, driver.DatasetListVersionsRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.ListDatasetVersionsResponseToProto(&driver.DatasetResponse{Versions: versions}), nil
}
