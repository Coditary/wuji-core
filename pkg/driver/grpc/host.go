package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/netx"
)

// HostOptions configures a standalone driver gRPC server.
type HostOptions struct {
	Addr     string
	CoreAddr string
}

// Host implements DriverService by delegating to a local driver.Driver.
type Host struct {
	wujiv1.UnimplementedDriverServiceServer
	drv driver.Driver
}

// NewHost wraps a driver for gRPC serving.
func NewHost(drv driver.Driver) *Host {
	return &Host{drv: drv}
}

func (s *Host) GetInfo(_ context.Context, _ *wujiv1.GetInfoRequest) (*wujiv1.GetInfoResponse, error) {
	info := s.drv.Info()
	caps := make([]string, 0, len(info.Capabilities))
	for _, c := range info.Capabilities {
		caps = append(caps, c.String())
	}
	return &wujiv1.GetInfoResponse{
		Metadata: &wujiv1.DriverMetadata{
			Id: info.ID, Name: info.Name, Version: info.Version,
			Description: info.Description, Capabilities: caps,
			FormatSupport: driver.CapabilityFormatsToProto(info.FormatSupport),
			ImageTasks:    driver.ImageTasksToProto(info.ImageTasks),
			VideoTasks:    driver.VideoTasksToProto(info.VideoTasks),
			DatasetTasks:  driver.DatasetTasksToProto(info.DatasetTasks),
			AudioTasks:    driver.AudioTasksToProto(info.AudioTasks),
			VoiceTasks:    driver.VoiceTasksToProto(info.VoiceTasks),
			MeshTasks:     driver.MeshTasksToProto(info.MeshTasks),
			DataTasks:     driver.DataTasksToProto(info.DataTasks),
			RagTasks:      driver.RAGTasksToProto(info.RAGTasks),
		},
	}, nil
}

func (s *Host) GenerateText(ctx context.Context, req *wujiv1.GenerateTextRequest) (*wujiv1.GenerateTextResponse, error) {
	resp, err := driver.RunTextInput(ctx, s.drv, driver.TextRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.TextResponseToProto(*resp), nil
}

func (s *Host) GenerateTextStream(req *wujiv1.GenerateTextRequest, stream wujiv1.DriverService_GenerateTextStreamServer) error {
	ctx := stream.Context()
	textReq := driver.TextRequestFromProto(req)
	textReq.Stream = true
	textReq.OnDelta = func(delta string) error {
		return stream.Send(&wujiv1.TextStreamChunk{Delta: delta})
	}
	textReq.OnThinkDelta = func(delta string) error {
		return stream.Send(&wujiv1.TextStreamChunk{ThinkDelta: delta})
	}

	resp, err := driver.RunTextInput(ctx, s.drv, textReq)
	if err != nil {
		return stream.Send(&wujiv1.TextStreamChunk{Error: err.Error(), Done: true})
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		return stream.Send(&wujiv1.TextStreamChunk{Error: err.Error(), Done: true})
	}
	return stream.Send(&wujiv1.TextStreamChunk{PayloadJson: payload, Done: true})
}

func (s *Host) GenerateMesh(ctx context.Context, req *wujiv1.GenerateMeshRequest) (*wujiv1.GenerateMeshResponse, error) {
	meshReq := driver.MeshRequestFromProto(req)
	resp, err := driver.RunMeshTask(ctx, s.drv, meshReq)
	if err != nil {
		return nil, err
	}
	return &wujiv1.GenerateMeshResponse{
		Path: resp.Path, Format: resp.Format, Task: string(resp.Task),
		Mode: string(resp.ModeOrDefault()), Representation: string(resp.RepresentationOrDefault()),
	}, nil
}

func (s *Host) TrainText(ctx context.Context, req *wujiv1.TrainTextRequest) (*wujiv1.TrainResponse, error) {
	tr, err := driver.As[driver.TextTrainer](s.drv, capability.Training)
	if err != nil {
		return nil, err
	}
	resp, err := tr.TrainText(ctx, driver.TextTrainRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.TrainResponseToProto(resp), nil
}

func (s *Host) TrainImage(ctx context.Context, req *wujiv1.TrainImageRequest) (*wujiv1.TrainResponse, error) {
	tr, err := driver.As[driver.ImageTrainer](s.drv, capability.Training)
	if err != nil {
		return nil, err
	}
	resp, err := tr.TrainImage(ctx, driver.ImageTrainRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.TrainResponseToProto(resp), nil
}

func (s *Host) TrainVideo(ctx context.Context, req *wujiv1.TrainVideoRequest) (*wujiv1.TrainResponse, error) {
	tr, err := driver.As[driver.VideoTrainer](s.drv, capability.Training)
	if err != nil {
		return nil, err
	}
	resp, err := tr.TrainVideo(ctx, driver.VideoTrainRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.TrainResponseToProto(resp), nil
}

func (s *Host) TrainAudio(ctx context.Context, req *wujiv1.TrainAudioRequest) (*wujiv1.TrainResponse, error) {
	tr, err := driver.As[driver.AudioTrainer](s.drv, capability.Training)
	if err != nil {
		return nil, err
	}
	resp, err := tr.TrainAudio(ctx, driver.AudioTrainRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.TrainResponseToProto(resp), nil
}

func (s *Host) TrainMesh(ctx context.Context, req *wujiv1.TrainMeshRequest) (*wujiv1.TrainResponse, error) {
	tr, err := driver.As[driver.MeshTrainer](s.drv, capability.Training)
	if err != nil {
		return nil, err
	}
	resp, err := tr.TrainMesh(ctx, driver.MeshTrainRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.TrainResponseToProto(resp), nil
}

func (s *Host) TrainVoice(ctx context.Context, req *wujiv1.TrainVoiceRequest) (*wujiv1.TrainResponse, error) {
	tr, err := driver.As[driver.VoiceTrainer](s.drv, capability.Training)
	if err != nil {
		return nil, err
	}
	resp, err := tr.TrainVoice(ctx, driver.VoiceTrainRequestFromProto(req))
	if err != nil {
		return nil, err
	}
	return driver.TrainResponseToProto(resp), nil
}

// Serve listens, optionally registers with a Wuji core, and blocks until SIGINT/SIGTERM.
func Serve(drv driver.Driver, opts HostOptions) error {
	lis, dialTarget, err := netx.Listen(opts.Addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	srv := grpc.NewServer()
	wujiv1.RegisterDriverServiceServer(srv, NewHost(drv))

	if opts.CoreAddr != "" {
		if err := registerWithCore(drv, opts.CoreAddr, dialTarget); err != nil {
			return err
		}
	}

	fmt.Fprintf(os.Stderr, "Driver %q listening on %s\n", drv.Info().ID, dialTarget)

	errCh := make(chan error, 1)
	go func() { errCh <- srv.Serve(lis) }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("server: %w", err)
	case <-sigCh:
		srv.GracefulStop()
		return nil
	}
}

func registerWithCore(drv driver.Driver, coreAddr, endpoint string) error {
	conn, err := grpc.NewClient(coreAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("connect to core: %w", err)
	}
	defer conn.Close()

	info := drv.Info()
	caps := make([]string, 0, len(info.Capabilities))
	for _, c := range info.Capabilities {
		caps = append(caps, c.String())
	}

	client := wujiv1.NewDriverRegistryClient(conn)
	resp, err := client.Register(context.Background(), &wujiv1.RegisterRequest{
		Metadata: &wujiv1.DriverMetadata{
			Id: info.ID, Name: info.Name, Version: info.Version,
			Description: info.Description, Capabilities: caps,
			FormatSupport: driver.CapabilityFormatsToProto(info.FormatSupport),
			ImageTasks:    driver.ImageTasksToProto(info.ImageTasks),
			VideoTasks:    driver.VideoTasksToProto(info.VideoTasks),
			DatasetTasks:  driver.DatasetTasksToProto(info.DatasetTasks),
			AudioTasks:    driver.AudioTasksToProto(info.AudioTasks),
			VoiceTasks:    driver.VoiceTasksToProto(info.VoiceTasks),
			MeshTasks:     driver.MeshTasksToProto(info.MeshTasks),
			DataTasks:     driver.DataTasksToProto(info.DataTasks),
			RagTasks:      driver.RAGTasksToProto(info.RAGTasks),
		},
		Endpoint: endpoint,
	})
	if err != nil {
		return fmt.Errorf("register with core: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Registered with core: %s\n", resp.GetMessage())
	return nil
}
