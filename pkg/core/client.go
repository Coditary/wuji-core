package core

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/core/remotemedia"
	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/netx"
	"github.com/coditary/wuji-core/pkg/scheduler"
)

// Client talks to a remote Wuji core daemon.
type Client struct {
	conn      *grpc.ClientConn
	client    wujiv1.WujiCoreClient
	endpoint  string
	broker    *remotemedia.Broker
	localRoot string
}

// DialClient connects to a running core daemon.
func DialClient(ctx context.Context, endpoint, localRoot string) (*Client, error) {
	conn, err := grpc.NewClient(netx.DialTarget(endpoint), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to core at %s: %w", endpoint, err)
	}
	grpcClient := wujiv1.NewWujiCoreClient(conn)
	c := &Client{
		conn:      conn,
		client:    grpcClient,
		endpoint:  endpoint,
		localRoot: localRoot,
	}
	if remotemedia.Enabled(endpoint) {
		c.broker = remotemedia.NewBroker(grpcClient, localRoot)
	}
	if _, err := c.client.Ping(ctx, &wujiv1.PingRequest{}); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping core at %s: %w", endpoint, err)
	}
	return c, nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) call(ctx context.Context, method string, payload any) ([]byte, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Call(ctx, &wujiv1.CallRequest{Method: method, PayloadJson: raw})
	if err != nil {
		return nil, err
	}
	if resp.GetError() != "" {
		return nil, fmt.Errorf("%s", resp.GetError())
	}
	return resp.GetPayloadJson(), nil
}

func (c *Client) GenerateText(ctx context.Context, driverID string, req driver.TextRequest) (*driver.TextResponse, error) {
	if err := c.prepareText(ctx, &req); err != nil {
		return nil, err
	}
	out, err := c.call(ctx, "GenerateText", map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	var resp driver.TextResponse
	return &resp, json.Unmarshal(out, &resp)
}

func (c *Client) GenerateTextStream(ctx context.Context, driverID string, req driver.TextRequest, onDelta func(string) error) (*driver.TextResponse, error) {
	if err := c.prepareText(ctx, &req); err != nil {
		return nil, err
	}
	onThinkDelta := req.OnThinkDelta
	req.Stream = true
	req.OnDelta = onDelta
	raw, err := json.Marshal(map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	stream, err := c.client.GenerateTextStream(ctx, &wujiv1.CallRequest{PayloadJson: raw})
	if err != nil {
		return nil, err
	}
	var resp driver.TextResponse
	for {
		chunk, err := stream.Recv()
		if err != nil {
			return nil, err
		}
		if chunk.GetError() != "" {
			return nil, fmt.Errorf("%s", chunk.GetError())
		}
		if delta := chunk.GetDelta(); delta != "" && onDelta != nil {
			if err := onDelta(delta); err != nil {
				return nil, err
			}
		}
		if think := chunk.GetThinkDelta(); think != "" && onThinkDelta != nil {
			if err := onThinkDelta(think); err != nil {
				return nil, err
			}
		}
		if chunk.GetDone() {
			if len(chunk.GetPayloadJson()) > 0 {
				if err := json.Unmarshal(chunk.GetPayloadJson(), &resp); err != nil {
					return nil, err
				}
			}
			return &resp, nil
		}
	}
}

func (c *Client) GenerateImage(ctx context.Context, driverID string, req driver.ImageRequest) (*driver.ImageResponse, error) {
	if err := c.prepareImage(ctx, &req); err != nil {
		return nil, err
	}
	out, err := c.call(ctx, "GenerateImage", map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	var resp driver.ImageResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, err
	}
	if err := c.finalizeImage(ctx, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GenerateVideo(ctx context.Context, driverID string, req driver.VideoRequest) (*driver.VideoResponse, error) {
	if err := c.prepareVideo(ctx, &req); err != nil {
		return nil, err
	}
	out, err := c.call(ctx, "GenerateVideo", map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	var resp driver.VideoResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, err
	}
	if err := c.finalizeVideo(ctx, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GenerateAudio(ctx context.Context, driverID string, req driver.AudioRequest) (*driver.AudioResponse, error) {
	if err := c.prepareAudio(ctx, &req); err != nil {
		return nil, err
	}
	out, err := c.call(ctx, "GenerateAudio", map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	var resp driver.AudioResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, err
	}
	if err := c.finalizeAudio(ctx, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GenerateMesh(ctx context.Context, driverID string, req driver.MeshRequest) (*driver.MeshResponse, error) {
	if err := c.prepareMesh(ctx, &req); err != nil {
		return nil, err
	}
	out, err := c.call(ctx, "GenerateMesh", map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	var resp driver.MeshResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, err
	}
	if err := c.finalizeMesh(ctx, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) ManageDataset(ctx context.Context, driverID string, req driver.DatasetRequest) (*driver.DatasetResponse, error) {
	if err := c.prepareDataset(ctx, &req); err != nil {
		return nil, err
	}
	out, err := c.call(ctx, "ManageDataset", map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	var resp driver.DatasetResponse
	return &resp, json.Unmarshal(out, &resp)
}

func (c *Client) TrainText(ctx context.Context, driverID string, req driver.TextTrainRequest) (*driver.TrainResponse, error) {
	if err := c.prepareTextTrain(ctx, &req); err != nil {
		return nil, err
	}
	out, err := c.call(ctx, "TrainText", map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	var resp driver.TrainResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, err
	}
	if err := c.finalizeTrain(ctx, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) TrainImage(ctx context.Context, driverID string, req driver.ImageTrainRequest) (*driver.TrainResponse, error) {
	if err := c.prepareImageTrain(ctx, &req); err != nil {
		return nil, err
	}
	out, err := c.call(ctx, "TrainImage", map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	var resp driver.TrainResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, err
	}
	if err := c.finalizeTrain(ctx, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) TrainVideo(ctx context.Context, driverID string, req driver.VideoTrainRequest) (*driver.TrainResponse, error) {
	if err := c.prepareVideoTrain(ctx, &req); err != nil {
		return nil, err
	}
	out, err := c.call(ctx, "TrainVideo", map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	var resp driver.TrainResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, err
	}
	if err := c.finalizeTrain(ctx, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) TrainAudio(ctx context.Context, driverID string, req driver.AudioTrainRequest) (*driver.TrainResponse, error) {
	if err := c.prepareAudioTrain(ctx, &req); err != nil {
		return nil, err
	}
	out, err := c.call(ctx, "TrainAudio", map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	var resp driver.TrainResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, err
	}
	if err := c.finalizeTrain(ctx, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) TrainMesh(ctx context.Context, driverID string, req driver.MeshTrainRequest) (*driver.TrainResponse, error) {
	if err := c.prepareMeshTrain(ctx, &req); err != nil {
		return nil, err
	}
	out, err := c.call(ctx, "TrainMesh", map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	var resp driver.TrainResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, err
	}
	if err := c.finalizeTrain(ctx, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) TrainVoice(ctx context.Context, driverID string, req driver.VoiceTrainRequest) (*driver.TrainResponse, error) {
	if err := c.prepareVoiceTrain(ctx, &req); err != nil {
		return nil, err
	}
	out, err := c.call(ctx, "TrainVoice", map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	var resp driver.TrainResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, err
	}
	if err := c.finalizeTrain(ctx, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) RunRAG(ctx context.Context, driverID string, req driver.RAGRequest) (data.DataShape, error) {
	if err := c.prepareRAG(ctx, &req); err != nil {
		return nil, err
	}
	out, err := c.call(ctx, "RunRAG", map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Shape json.RawMessage `json:"shape"`
	}
	if err := json.Unmarshal(out, &wrapper); err != nil {
		return nil, err
	}
	if err := c.finalizeRAG(ctx, req); err != nil {
		return nil, err
	}
	return data.ParseJSONShape(wrapper.Shape)
}

func (c *Client) ProduceData(ctx context.Context, driverID string, req driver.DataRequest) (data.DataShape, error) {
	if err := c.prepareData(ctx, &req); err != nil {
		return nil, err
	}
	out, err := c.call(ctx, "ProduceData", map[string]any{"driver_id": driverID, "request": req})
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Shape json.RawMessage `json:"shape"`
	}
	if err := json.Unmarshal(out, &wrapper); err != nil {
		return nil, err
	}
	return data.ParseJSONShape(wrapper.Shape)
}

func (c *Client) ConnectRemote(ctx context.Context, endpoint string, persist bool) error {
	_, err := c.call(ctx, "ConnectRemote", map[string]any{"endpoint": endpoint, "persist": persist})
	return err
}

func (c *Client) ListDrivers() []driver.Info {
	out, err := c.call(context.Background(), "ListDrivers", map[string]any{})
	if err != nil {
		return nil
	}
	var drivers []driver.Info
	if err := json.Unmarshal(out, &drivers); err != nil {
		return nil
	}
	return drivers
}

func (c *Client) ResourceStatus() (scheduler.PersistedState, error) {
	out, err := c.call(context.Background(), "ResourceStatus", map[string]any{})
	if err != nil {
		return scheduler.PersistedState{}, err
	}
	var state scheduler.PersistedState
	return state, json.Unmarshal(out, &state)
}

func (c *Client) ResourcesConfig() config.EffectiveResources {
	out, err := c.call(context.Background(), "ResourcesConfig", map[string]any{})
	if err != nil {
		return config.EffectiveResources{}
	}
	var cfg config.EffectiveResources
	_ = json.Unmarshal(out, &cfg)
	return cfg
}

func (c *Client) LoadInference(ctx context.Context, driverID, model string, cap capability.Type) error {
	_, err := c.call(ctx, "LoadInference", map[string]any{"driver_id": driverID, "model": model, "capability": cap})
	return err
}

func (c *Client) UnloadInferenceDriver(ctx context.Context, driverID string) error {
	_, err := c.call(ctx, "UnloadInferenceDriver", map[string]any{"driver_id": driverID})
	return err
}

func (c *Client) UnloadAllInference(ctx context.Context) error {
	_, err := c.call(ctx, "UnloadAllInference", map[string]any{})
	return err
}

func (c *Client) DefaultDriverID() string {
	out, err := c.call(context.Background(), "DefaultDriverID", map[string]any{})
	if err != nil {
		return ""
	}
	var body struct {
		DriverID string `json:"driver_id"`
	}
	if err := json.Unmarshal(out, &body); err != nil {
		return ""
	}
	return body.DriverID
}
