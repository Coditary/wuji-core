package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *Driver) VideoToAudio(ctx context.Context, req driver.Video2AudioRequest) (*driver.Video2AudioResponse, error) {
	if d.entry.Video2Audio == nil {
		return nil, fmt.Errorf("driver %q does not support video2audio", d.id)
	}
	spec := config.ResolveAPICapabilitySpec(d.entry.Video2Audio, "extract")
	raw, err := doHTTP(d.scoped(ctx), spec, video2AudioHTTPVars(req, spec))
	if err != nil {
		return nil, err
	}
	path, err := mediaPathFromHTTP(spec, raw, "audio")
	if err != nil {
		return nil, err
	}
	return &driver.Video2AudioResponse{AudioPath: path, Temporary: true}, nil
}

func (d *Driver) ProduceData(ctx context.Context, req driver.DataRequest) (data.DataShape, error) {
	if d.entry.Data == nil {
		return nil, fmt.Errorf("driver %q does not support data", d.id)
	}
	task := req.TaskOrDefault()
	spec := config.ResolveAPICapabilitySpec(d.entry.Data, string(task))
	raw, err := doHTTP(d.scoped(ctx), spec, dataHTTPVars(req, spec))
	if err != nil {
		return nil, err
	}
	return jsonShapeFromHTTP(raw, "data", string(task))
}

func (d *Driver) ProduceRAG(ctx context.Context, req driver.RAGRequest) (data.DataShape, error) {
	if d.entry.RAG == nil {
		return nil, fmt.Errorf("driver %q does not support rag", d.id)
	}
	spec := config.ResolveAPICapabilitySpec(d.entry.RAG, string(req.Task))
	raw, err := doHTTP(d.scoped(ctx), spec, ragHTTPVars(req, spec))
	if err != nil {
		return nil, err
	}
	return jsonShapeFromHTTP(raw, "rag", string(req.Task))
}

func (d *Driver) ListDatasets(ctx context.Context, _ driver.DatasetListRequest) (*driver.DatasetResponse, error) {
	return d.runDataset(ctx, driver.DatasetTaskList, driver.DatasetRequest{Task: driver.DatasetTaskList})
}

func (d *Driver) CreateDataset(ctx context.Context, req driver.DatasetCreateRequest) (*driver.DatasetResponse, error) {
	return d.runDataset(ctx, driver.DatasetTaskCreate, driver.DatasetRequest{
		Task: driver.DatasetTaskCreate, Name: req.Name, Path: req.Path, Description: req.Description,
	})
}

func (d *Driver) DeleteDataset(ctx context.Context, req driver.DatasetDeleteRequest) (*driver.DatasetResponse, error) {
	return d.runDataset(ctx, driver.DatasetTaskDelete, driver.DatasetRequest{
		Task: driver.DatasetTaskDelete, Name: req.Name,
	})
}

func (d *Driver) IngestDataset(ctx context.Context, req driver.DatasetIngestRequest) (*driver.DatasetIngestResult, error) {
	if d.entry.Dataset == nil {
		return nil, fmt.Errorf("driver %q does not support dataset", d.id)
	}
	spec := config.ResolveAPICapabilitySpec(d.entry.Dataset, string(driver.DatasetTaskIngest))
	union := driver.DatasetRequest{
		Task: driver.DatasetTaskIngest, DatasetID: req.DatasetID, SourcePath: req.SourcePath, Recursive: req.Recursive,
	}
	_, err := doHTTP(d.scoped(ctx), spec, datasetHTTPVars(union, spec))
	if err != nil {
		return nil, err
	}
	return &driver.DatasetIngestResult{FilesAdded: 1, Message: "ingested"}, nil
}

func (d *Driver) CreateDatasetVersion(ctx context.Context, req driver.DatasetVersionRequest) (*driver.DatasetVersionResult, error) {
	if d.entry.Dataset == nil {
		return nil, fmt.Errorf("driver %q does not support dataset", d.id)
	}
	spec := config.ResolveAPICapabilitySpec(d.entry.Dataset, string(driver.DatasetTaskVersion))
	union := driver.DatasetRequest{Task: driver.DatasetTaskVersion, DatasetID: req.DatasetID, Tag: req.Tag, Message: req.Message}
	raw, err := doHTTP(d.scoped(ctx), spec, datasetHTTPVars(union, spec))
	if err != nil {
		return nil, err
	}
	out := &driver.DatasetVersionResult{}
	if spec.Response != nil && spec.Response["version"] != "" {
		id, _ := extractJSONPath(raw, spec.Response["version"])
		out.Version = driver.DatasetVersionEntry{ID: id, Tag: req.Tag}
	}
	return out, nil
}

func (d *Driver) ListDatasetVersions(ctx context.Context, req driver.DatasetListVersionsRequest) ([]driver.DatasetVersionEntry, error) {
	if d.entry.Dataset == nil {
		return nil, fmt.Errorf("driver %q does not support dataset", d.id)
	}
	spec := config.ResolveAPICapabilitySpec(d.entry.Dataset, string(driver.DatasetTaskListVersions))
	union := driver.DatasetRequest{Task: driver.DatasetTaskListVersions, DatasetID: req.DatasetID}
	raw, err := doHTTP(d.scoped(ctx), spec, datasetHTTPVars(union, spec))
	if err != nil {
		return nil, err
	}
	var entries []driver.DatasetVersionEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, fmt.Errorf("parse dataset versions: %w", err)
	}
	return entries, nil
}

func (d *Driver) runDataset(ctx context.Context, task driver.DatasetTask, req driver.DatasetRequest) (*driver.DatasetResponse, error) {
	if d.entry.Dataset == nil {
		return nil, fmt.Errorf("driver %q does not support dataset", d.id)
	}
	spec := config.ResolveAPICapabilitySpec(d.entry.Dataset, string(task))
	raw, err := doHTTP(d.scoped(ctx), spec, datasetHTTPVars(req, spec))
	if err != nil {
		return nil, err
	}
	resp := &driver.DatasetResponse{}
	if spec.Response != nil && spec.Response["datasets"] != "" {
		text, err := extractJSONPath(raw, spec.Response["datasets"])
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(text), &resp.Datasets); err != nil {
			if err2 := json.Unmarshal(raw, &resp); err2 != nil {
				return nil, fmt.Errorf("parse datasets: %w", err)
			}
		}
		return resp, nil
	}
	if err := json.Unmarshal(raw, resp); err != nil {
		return nil, fmt.Errorf("parse dataset response: %w", err)
	}
	return resp, nil
}

func (d *Driver) TrainText(ctx context.Context, req driver.TextTrainRequest) (*driver.TrainResponse, error) {
	return d.runTrain(ctx, capability.TextGeneration, "text", textTrainHTTPVars(req, d.trainSpec("text")))
}

func (d *Driver) TrainImage(ctx context.Context, req driver.ImageTrainRequest) (*driver.TrainResponse, error) {
	return d.runTrain(ctx, capability.ImageGeneration, "image", imageTrainHTTPVars(req, d.trainSpec("image")))
}

func (d *Driver) TrainVideo(ctx context.Context, req driver.VideoTrainRequest) (*driver.TrainResponse, error) {
	return d.runTrain(ctx, capability.VideoGeneration, "video", videoTrainHTTPVars(req, d.trainSpec("video")))
}

func (d *Driver) TrainAudio(ctx context.Context, req driver.AudioTrainRequest) (*driver.TrainResponse, error) {
	return d.runTrain(ctx, capability.AudioGeneration, "audio", audioTrainHTTPVars(req, d.trainSpec("audio")))
}

func (d *Driver) TrainMesh(ctx context.Context, req driver.MeshTrainRequest) (*driver.TrainResponse, error) {
	return d.runTrain(ctx, capability.Mesh, "mesh", meshTrainHTTPVars(req, d.trainSpec("mesh")))
}

func (d *Driver) TrainVoice(ctx context.Context, req driver.VoiceTrainRequest) (*driver.TrainResponse, error) {
	return d.runTrain(ctx, capability.VoiceCloning, "voice", voiceTrainHTTPVars(req, d.trainSpec("voice")))
}

func (d *Driver) trainSpec(task string) *config.APICapabilitySpec {
	return config.ResolveAPICapabilitySpec(d.entry.Train, task)
}

func (d *Driver) runTrain(ctx context.Context, cap capability.Type, task string, vars map[string]string) (*driver.TrainResponse, error) {
	if d.entry.Train == nil {
		return nil, fmt.Errorf("driver %q does not support train", d.id)
	}
	spec := config.ResolveAPICapabilitySpec(d.entry.Train, task)
	raw, err := doHTTP(d.scoped(ctx), spec, vars)
	if err != nil {
		return nil, err
	}
	resp := &driver.TrainResponse{Capability: cap, Status: "submitted"}
	if spec.Response != nil {
		if p := spec.Response["job_id"]; p != "" {
			resp.JobID, _ = extractJSONPath(raw, p)
		}
		if p := spec.Response["status"]; p != "" {
			resp.Status, _ = extractJSONPath(raw, p)
		}
		if p := spec.Response["output_path"]; p != "" {
			resp.OutputPath, _ = extractJSONPath(raw, p)
		}
	}
	return resp, nil
}

func jsonShapeFromHTTP(raw []byte, capabilityName, task string) (data.DataShape, error) {
	return &data.RecordSet{
		MetaData: data.Meta{Capability: capabilityName, Task: task},
		Records: []data.Record{{
			Fields: map[string]data.Value{
				"raw_json": data.NewString(string(raw)),
			},
		}},
	}, nil
}
