package grpc

import (
	"context"
	"encoding/json"
	"fmt"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/core"
	"github.com/coditary/wuji-core/pkg/core/remotemedia"
	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/version"
)

// CoreService implements the WujiCore gRPC API.
type CoreService struct {
	wujiv1.UnimplementedWujiCoreServer
	core *core.Core
}

func NewCoreService(c *core.Core) *CoreService {
	return &CoreService{core: c}
}

func (s *CoreService) PutFiles(ctx context.Context, req *wujiv1.PutFilesRequest) (*wujiv1.PutFilesResponse, error) {
	root := "."
	if cfg := s.core.AppConfig(); cfg != nil && cfg.Root != "" {
		root = cfg.Root
	}
	paths, stagingRoot, err := remotemedia.WriteStaging(root, req.GetSessionId(), req.GetFiles())
	if err != nil {
		return nil, err
	}
	return &wujiv1.PutFilesResponse{Paths: paths, StagingRoot: stagingRoot}, nil
}

func (s *CoreService) GetFiles(ctx context.Context, req *wujiv1.GetFilesRequest) (*wujiv1.GetFilesResponse, error) {
	root := "."
	if cfg := s.core.AppConfig(); cfg != nil && cfg.Root != "" {
		root = cfg.Root
	}
	resolved := remotemedia.ResolvePaths(root, req.GetPaths())
	files, err := remotemedia.ReadFiles(resolved)
	if err != nil {
		return nil, err
	}
	return &wujiv1.GetFilesResponse{Files: files}, nil
}

func (s *CoreService) Ping(context.Context, *wujiv1.PingRequest) (*wujiv1.PingResponse, error) {
	return &wujiv1.PingResponse{Ok: true, Version: version.Version}, nil
}

func (s *CoreService) Call(ctx context.Context, req *wujiv1.CallRequest) (*wujiv1.CallResponse, error) {
	out, err := s.dispatch(ctx, req.GetMethod(), req.GetPayloadJson())
	if err != nil {
		return &wujiv1.CallResponse{Error: err.Error()}, nil
	}
	return &wujiv1.CallResponse{PayloadJson: out}, nil
}

func (s *CoreService) dispatch(ctx context.Context, method string, payload []byte) ([]byte, error) {
	switch method {
	case "GenerateText":
		return callJSON(ctx, payload, s.core.GenerateText)
	case "ChatComplete":
		var body struct {
			ProviderID string `json:"provider_id"`
			DriverID   string `json:"driver_id"`
			Request    struct {
				Messages    []driver.ChatMessage `json:"messages"`
				Tools       []driver.ChatToolDef `json:"tools"`
				Model       string               `json:"model"`
				MaxTokens   int                  `json:"max_tokens"`
				Temperature float64              `json:"temperature"`
			} `json:"request"`
		}
		if err := json.Unmarshal(payload, &body); err != nil {
			return nil, err
		}
		driverID := body.DriverID
		if driverID == "" {
			driverID = body.ProviderID
		}
		textReq := driver.TextRequest{
			Messages:    body.Request.Messages,
			Tools:       body.Request.Tools,
			Model:       body.Request.Model,
			MaxTokens:   body.Request.MaxTokens,
			Temperature: float32(body.Request.Temperature),
			Prompt:      "",
		}
		raw, err := json.Marshal(map[string]any{"driver_id": driverID, "request": textReq})
		if err != nil {
			return nil, err
		}
		return callJSON(ctx, raw, s.core.GenerateText)
	case "GenerateImage":
		return callJSON(ctx, payload, s.core.GenerateImage)
	case "GenerateVideo":
		return callJSON(ctx, payload, s.core.GenerateVideo)
	case "GenerateAudio":
		return callJSON(ctx, payload, s.core.GenerateAudio)
	case "GenerateMesh":
		return callJSON(ctx, payload, s.core.GenerateMesh)
	case "ManageDataset":
		return callJSON(ctx, payload, s.core.ManageDataset)
	case "TrainText":
		return callJSON(ctx, payload, s.core.TrainText)
	case "TrainImage":
		return callJSON(ctx, payload, s.core.TrainImage)
	case "TrainVideo":
		return callJSON(ctx, payload, s.core.TrainVideo)
	case "TrainAudio":
		return callJSON(ctx, payload, s.core.TrainAudio)
	case "TrainMesh":
		return callJSON(ctx, payload, s.core.TrainMesh)
	case "TrainVoice":
		return callJSON(ctx, payload, s.core.TrainVoice)
	case "RunRAG":
		var body struct {
			DriverID string            `json:"driver_id"`
			Request  driver.RAGRequest `json:"request"`
		}
		if err := json.Unmarshal(payload, &body); err != nil {
			return nil, err
		}
		shape, err := s.core.RunRAG(ctx, body.DriverID, body.Request)
		if err != nil {
			return nil, err
		}
		raw, err := data.MarshalJSONShape(shape)
		if err != nil {
			return nil, err
		}
		return json.Marshal(map[string]any{"shape": json.RawMessage(raw)})
	case "ProduceData":
		var body struct {
			DriverID string             `json:"driver_id"`
			Request  driver.DataRequest `json:"request"`
		}
		if err := json.Unmarshal(payload, &body); err != nil {
			return nil, err
		}
		shape, err := s.core.ProduceData(ctx, body.DriverID, body.Request)
		if err != nil {
			return nil, err
		}
		raw, err := data.MarshalJSONShape(shape)
		if err != nil {
			return nil, err
		}
		return json.Marshal(map[string]any{"shape": json.RawMessage(raw)})
	case "ConnectRemote":
		var body struct {
			Endpoint string `json:"endpoint"`
			Persist  bool   `json:"persist"`
		}
		if err := json.Unmarshal(payload, &body); err != nil {
			return nil, err
		}
		if err := s.core.ConnectRemote(ctx, body.Endpoint, body.Persist); err != nil {
			return nil, err
		}
		return []byte("{}"), nil
	case "ListDrivers":
		return json.Marshal(s.core.ListDrivers())
	case "ResourceStatus":
		state, err := s.core.ResourceStatus()
		if err != nil {
			return nil, err
		}
		return json.Marshal(state)
	case "ResourcesConfig":
		return json.Marshal(s.core.ResourcesConfig())
	case "LoadInference":
		var body struct {
			DriverID string          `json:"driver_id"`
			Model    string          `json:"model"`
			Cap      capability.Type `json:"capability"`
		}
		if err := json.Unmarshal(payload, &body); err != nil {
			return nil, err
		}
		if err := s.core.LoadInference(ctx, body.DriverID, body.Model, body.Cap); err != nil {
			return nil, err
		}
		return []byte("{}"), nil
	case "UnloadInferenceDriver":
		var body struct {
			DriverID string `json:"driver_id"`
		}
		if err := json.Unmarshal(payload, &body); err != nil {
			return nil, err
		}
		if err := s.core.UnloadInferenceDriver(ctx, body.DriverID); err != nil {
			return nil, err
		}
		return []byte("{}"), nil
	case "UnloadAllInference":
		if err := s.core.UnloadAllInference(ctx); err != nil {
			return nil, err
		}
		return []byte("{}"), nil
	case "DefaultDriverID":
		return json.Marshal(map[string]string{"driver_id": s.core.DefaultDriverID()})
	case "ListProviders":
		return json.Marshal(s.core.ListProviders())
	case "DefaultProviderID":
		return json.Marshal(map[string]string{"provider_id": s.core.DefaultProviderID()})
	default:
		return nil, fmt.Errorf("unknown method %q", method)
	}
}

func callJSON[T any, R any](ctx context.Context, payload []byte, fn func(context.Context, string, T) (R, error)) ([]byte, error) {
	var body struct {
		DriverID string `json:"driver_id"`
		Request  T      `json:"request"`
	}
	if err := json.Unmarshal(payload, &body); err != nil {
		return nil, err
	}
	resp, err := fn(ctx, body.DriverID, body.Request)
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp)
}
