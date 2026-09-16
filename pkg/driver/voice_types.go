package driver

// VoiceCreateRequest is input for zero-shot voice creation (CreateVoice RPC).
type VoiceCreateRequest struct {
	Name            string
	SamplePath      string
	TargetModel     string
	Denoise         bool
	DenoiseStrength float32
}

// VoiceConvertRequest is input for voice conversion (ConvertVoice RPC).
type VoiceConvertRequest struct {
	Name            string
	SamplePath      string
	TargetModel     string
	SourcePath      string
	PitchShift      int
	IndexRate       float32
	Protect         float32
	Vocoder         string
	Denoise         bool
	DenoiseStrength float32
	ChunkSize       int
	Crossfade       float32
}

func (r VoiceRequest) ToCreateRequest() VoiceCreateRequest {
	return VoiceCreateRequest{
		Name: r.Name, SamplePath: r.SamplePath, TargetModel: r.TargetModel,
		Denoise: r.Denoise, DenoiseStrength: r.DenoiseStrength,
	}
}

func (r VoiceRequest) ToConvertRequest() VoiceConvertRequest {
	return VoiceConvertRequest{
		Name: r.Name, SamplePath: r.SamplePath, TargetModel: r.TargetModel, SourcePath: r.SourcePath,
		PitchShift: r.PitchShift, IndexRate: r.IndexRate, Protect: r.Protect, Vocoder: r.Vocoder,
		Denoise: r.Denoise, DenoiseStrength: r.DenoiseStrength, ChunkSize: r.ChunkSize, Crossfade: r.Crossfade,
	}
}
