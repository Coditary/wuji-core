package grpc

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/netx"
)

// RemoteDriver wraps a driver exposed over gRPC.
type RemoteDriver struct {
	info   driver.Info
	client wujiv1.DriverServiceClient
	conn   *grpc.ClientConn
}

// Connect creates a remote driver client for the given endpoint.
func Connect(ctx context.Context, endpoint string) (*RemoteDriver, error) {
	conn, err := grpc.NewClient(netx.DialTarget(endpoint), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to driver at %s: %w", endpoint, err)
	}

	client := wujiv1.NewDriverServiceClient(conn)

	resp, err := client.GetInfo(ctx, &wujiv1.GetInfoRequest{})
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("get driver info: %w", err)
	}

	meta := resp.GetMetadata()
	caps := make([]capability.Type, 0, len(meta.GetCapabilities()))
	for _, c := range meta.GetCapabilities() {
		caps = append(caps, capability.Type(c))
	}

	return &RemoteDriver{
		info: driver.Info{
			ID:            meta.GetId(),
			Name:          meta.GetName(),
			Version:       meta.GetVersion(),
			Description:   meta.GetDescription(),
			Capabilities:  caps,
			ImageTasks:    driver.ImageTasksFromProto(meta.GetImageTasks()),
			VideoTasks:    driver.VideoTasksFromProto(meta.GetVideoTasks()),
			DatasetTasks:  driver.DatasetTasksFromProto(meta.GetDatasetTasks()),
			AudioTasks:    driver.AudioTasksFromProto(meta.GetAudioTasks()),
			VoiceTasks:    driver.VoiceTasksFromProto(meta.GetVoiceTasks()),
			MeshTasks:     driver.MeshTasksFromProto(meta.GetMeshTasks()),
			DataTasks:     driver.DataTasksFromProto(meta.GetDataTasks()),
			RAGTasks:      driver.RAGTasksFromProto(meta.GetRagTasks()),
			FormatSupport: driver.CapabilityFormatsFromProto(meta.GetFormatSupport()),
			Remote:        true,
			Endpoint:      endpoint,
		},
		client: client,
		conn:   conn,
	}, nil
}

func (d *RemoteDriver) Info() driver.Info               { return d.info }
func (d *RemoteDriver) Capabilities() []capability.Type { return d.info.Capabilities }

func (d *RemoteDriver) GenerateText(ctx context.Context, req driver.TextRequest) (*driver.TextResponse, error) {
	if req.OnDelta != nil || req.OnThinkDelta != nil || req.Stream {
		return d.generateTextStream(ctx, req)
	}
	resp, err := d.client.GenerateText(ctx, driver.TextRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote generate text: %w", err)
	}
	out := driver.TextResponseFromProto(resp)
	return &out, nil
}

func (d *RemoteDriver) generateTextStream(ctx context.Context, req driver.TextRequest) (*driver.TextResponse, error) {
	protoReq := driver.TextRequestToProto(req)
	protoReq.Stream = true
	stream, err := d.client.GenerateTextStream(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("remote generate text stream: %w", err)
	}
	for {
		chunk, err := stream.Recv()
		if err != nil {
			return nil, fmt.Errorf("remote text stream recv: %w", err)
		}
		if chunk.GetError() != "" {
			return nil, fmt.Errorf("%s", chunk.GetError())
		}
		if delta := chunk.GetDelta(); delta != "" && req.OnDelta != nil {
			if err := req.OnDelta(delta); err != nil {
				return nil, err
			}
		}
		if think := chunk.GetThinkDelta(); think != "" && req.OnThinkDelta != nil {
			if err := req.OnThinkDelta(think); err != nil {
				return nil, err
			}
		}
		if chunk.GetDone() {
			if len(chunk.GetPayloadJson()) == 0 {
				return &driver.TextResponse{FinishReason: "stop"}, nil
			}
			var resp driver.TextResponse
			if err := json.Unmarshal(chunk.GetPayloadJson(), &resp); err != nil {
				return nil, fmt.Errorf("decode stream payload: %w", err)
			}
			return &resp, nil
		}
	}
}

// WarmModel preloads remote inference backends (vLLM, llama) by triggering a minimal request.
func (d *RemoteDriver) WarmModel(ctx context.Context, model string) error {
	_, err := d.GenerateText(ctx, driver.TextRequest{
		Model:     model,
		MaxTokens: 1,
		Prompt:    "ping",
	})
	return err
}

var _ driver.ModelWarmer = (*RemoteDriver)(nil)

func (d *RemoteDriver) GenerateMesh(ctx context.Context, req driver.MeshRequest) (*driver.MeshResponse, error) {
	resp, err := d.client.GenerateMesh(ctx, driver.MeshRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote generate mesh: %w", err)
	}
	return driver.MeshResponseFromProto(resp), nil
}

func (d *RemoteDriver) AudioToText(ctx context.Context, req driver.TextRequest) (*driver.TextResponse, error) {
	return d.GenerateText(ctx, req)
}

func (d *RemoteDriver) TrainText(ctx context.Context, req driver.TextTrainRequest) (*driver.TrainResponse, error) {
	resp, err := d.client.TrainText(ctx, driver.TextTrainRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote train text: %w", err)
	}
	return driver.TrainResponseFromProto(resp), nil
}

func (d *RemoteDriver) TrainImage(ctx context.Context, req driver.ImageTrainRequest) (*driver.TrainResponse, error) {
	resp, err := d.client.TrainImage(ctx, driver.ImageTrainRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote train image: %w", err)
	}
	return driver.TrainResponseFromProto(resp), nil
}

func (d *RemoteDriver) TrainVideo(ctx context.Context, req driver.VideoTrainRequest) (*driver.TrainResponse, error) {
	resp, err := d.client.TrainVideo(ctx, driver.VideoTrainRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote train video: %w", err)
	}
	return driver.TrainResponseFromProto(resp), nil
}

func (d *RemoteDriver) TrainAudio(ctx context.Context, req driver.AudioTrainRequest) (*driver.TrainResponse, error) {
	resp, err := d.client.TrainAudio(ctx, driver.AudioTrainRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote train audio: %w", err)
	}
	return driver.TrainResponseFromProto(resp), nil
}

func (d *RemoteDriver) TrainMesh(ctx context.Context, req driver.MeshTrainRequest) (*driver.TrainResponse, error) {
	resp, err := d.client.TrainMesh(ctx, driver.MeshTrainRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote train mesh: %w", err)
	}
	return driver.TrainResponseFromProto(resp), nil
}

func (d *RemoteDriver) TrainVoice(ctx context.Context, req driver.VoiceTrainRequest) (*driver.TrainResponse, error) {
	resp, err := d.client.TrainVoice(ctx, driver.VoiceTrainRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("remote train voice: %w", err)
	}
	return driver.TrainResponseFromProto(resp), nil
}

func (d *RemoteDriver) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}
