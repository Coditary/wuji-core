package core

import (
	"context"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/scheduler"
)

// Runtime is the core API used by the CLI and external clients.
type Runtime interface {
	GenerateText(ctx context.Context, driverID string, req driver.TextRequest) (*driver.TextResponse, error)
	GenerateTextStream(ctx context.Context, driverID string, req driver.TextRequest, onDelta func(string) error) (*driver.TextResponse, error)
	GenerateImage(ctx context.Context, driverID string, req driver.ImageRequest) (*driver.ImageResponse, error)
	GenerateVideo(ctx context.Context, driverID string, req driver.VideoRequest) (*driver.VideoResponse, error)
	GenerateAudio(ctx context.Context, driverID string, req driver.AudioRequest) (*driver.AudioResponse, error)
	GenerateMesh(ctx context.Context, driverID string, req driver.MeshRequest) (*driver.MeshResponse, error)
	ManageDataset(ctx context.Context, driverID string, req driver.DatasetRequest) (*driver.DatasetResponse, error)
	TrainText(ctx context.Context, driverID string, req driver.TextTrainRequest) (*driver.TrainResponse, error)
	TrainImage(ctx context.Context, driverID string, req driver.ImageTrainRequest) (*driver.TrainResponse, error)
	TrainVideo(ctx context.Context, driverID string, req driver.VideoTrainRequest) (*driver.TrainResponse, error)
	TrainAudio(ctx context.Context, driverID string, req driver.AudioTrainRequest) (*driver.TrainResponse, error)
	TrainMesh(ctx context.Context, driverID string, req driver.MeshTrainRequest) (*driver.TrainResponse, error)
	TrainVoice(ctx context.Context, driverID string, req driver.VoiceTrainRequest) (*driver.TrainResponse, error)
	RunRAG(ctx context.Context, driverID string, req driver.RAGRequest) (data.DataShape, error)
	ProduceData(ctx context.Context, driverID string, req driver.DataRequest) (data.DataShape, error)
	ConnectRemote(ctx context.Context, endpoint string, persist bool) error
	ListDrivers() []driver.Info
	ResourceStatus() (scheduler.PersistedState, error)
	ResourcesConfig() config.EffectiveResources
	LoadInference(ctx context.Context, driverID, model string, cap capability.Type) error
	UnloadInferenceDriver(ctx context.Context, driverID string) error
	UnloadAllInference(ctx context.Context) error
	DefaultDriverID() string
	Close() error
}
