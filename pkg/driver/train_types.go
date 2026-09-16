package driver

import "github.com/coditary/wuji-core/pkg/capability"

// TrainResponse is the output of any training job.
type TrainResponse struct {
	JobID      string
	Status     string
	Capability capability.Type
	OutputPath string
}

// TextTrainRequest configures LLM / text model fine-tuning.
type TextTrainRequest struct {
	Method                    TextTrainMethod
	Name                      string
	DatasetID                 string
	BaseModel                 string
	OutputPath                string
	Epochs                    int
	LearningRate              float32
	LoRARank                  int
	LoRAAlpha                 int
	ContextLength             int
	BatchSize                 int
	GradientAccumulationSteps int
	WarmupSteps               int
	SaveEveryEpoch            int
	Seed                      int
}

// ImageTrainRequest configures image / diffusion model training.
type ImageTrainRequest struct {
	Method            ImageTrainMethod
	Name              string
	DatasetID         string
	BaseModel         string
	OutputPath        string
	Epochs            int
	LearningRate      float32
	LoRARank          int
	LoRAAlpha         int
	Width             int
	Height            int
	ClassToken        string
	ControlType       ImageControlType
	Token             string
	Mode              ImageMode
	PriorPreservation bool
	RegDatasetID      string
	BatchSize         int
	Seed              int
}

// VideoTrainRequest configures video model training.
type VideoTrainRequest struct {
	Method        VideoTrainMethod
	Name          string
	DatasetID     string
	BaseModel     string
	OutputPath    string
	Epochs        int
	LearningRate  float32
	LoRARank      int
	Frames        int
	FPS           int
	ContextLength int
	BatchSize     int
	Seed          int
}

// AudioTrainRequest configures audio / music model training.
type AudioTrainRequest struct {
	Method       AudioTrainMethod
	Name         string
	DatasetID    string
	BaseModel    string
	OutputPath   string
	Epochs       int
	LearningRate float32
	SampleRate   int
	BatchSize    int
	TaskType     AudioTaskType
	Seed         int
}

// MeshTrainRequest configures mesh model training.
type MeshTrainRequest struct {
	Method       MeshTrainMethod
	Name         string
	DatasetID    string
	BaseModel    string
	OutputPath   string
	Epochs       int
	LearningRate float32
	LoRARank     int
	Format       string
	Mode         MeshMode
	Seed         int
}

// VoiceTrainRequest configures voice cloning model training (e.g. RVC).
type VoiceTrainRequest struct {
	Method          VoiceTrainMethod
	Name            string
	DatasetID       string
	OutputPath      string
	Epochs          int
	SamplePath      string
	PretrainedModel string
	PitchShift      int
	IndexRate       float32
	BatchSize       int
	SaveEveryEpoch  int
	Seed            int
}
