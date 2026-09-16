package driver

// AudioCommonParams are shared generation parameters for music tasks.
type AudioCommonParams struct {
	Prompt         string
	Lyrics         string
	NegativePrompt string
	Model          string
	Duration       float32
	Overlap        float32
	Temperature    float32
	CFGScale       float32
	TopP           float32
	TopK           int
	SampleRate     int
	Seed           *int
}

// AudioMusicRequest is input for text-to-music (GenerateMusic RPC).
type AudioMusicRequest struct {
	AudioCommonParams
}

// AudioMelodyToMusicRequest is input for melody-to-music (MelodyToMusic RPC).
type AudioMelodyToMusicRequest struct {
	ReferencePath string
	AudioCommonParams
}

// AudioSFXRequest is input for text-to-sfx (GenerateSFX RPC).
type AudioSFXRequest struct {
	Prompt         string
	NegativePrompt string
	Model          string
	Duration       float32
	Temperature    float32
	CFGScale       float32
	TopP           float32
	TopK           int
	SampleRate     int
	Seed           *int
}

// AudioSpeechRequest is input for text-to-speech (GenerateAudioSpeech RPC).
type AudioSpeechRequest struct {
	Prompt     string
	Model      string
	Duration   float32
	SampleRate int
	Voice      string
	Language   string
	Speed      float32
	Pitch      int
	Emotion    string
	Style      string
	Energy     float32
	Seed       *int
}

func commonAudioFromRequest(r AudioRequest) AudioCommonParams {
	return AudioCommonParams{
		Prompt: r.Prompt, Lyrics: r.Lyrics, NegativePrompt: r.NegativePrompt, Model: r.Model,
		Duration: r.Duration, Overlap: r.Overlap, Temperature: r.Temperature, CFGScale: r.CFGScale,
		TopP: r.TopP, TopK: r.TopK, SampleRate: r.SampleRate, Seed: r.Seed,
	}
}

func (r AudioRequest) ToMusicRequest() AudioMusicRequest {
	return AudioMusicRequest{AudioCommonParams: commonAudioFromRequest(r)}
}

func (r AudioRequest) ToMelodyToMusicRequest() AudioMelodyToMusicRequest {
	return AudioMelodyToMusicRequest{
		ReferencePath: r.ReferencePath, AudioCommonParams: commonAudioFromRequest(r),
	}
}

func (r AudioRequest) ToSFXRequest() AudioSFXRequest {
	return AudioSFXRequest{
		Prompt: r.Prompt, NegativePrompt: r.NegativePrompt, Model: r.Model, Duration: r.Duration,
		Temperature: r.Temperature, CFGScale: r.CFGScale, TopP: r.TopP, TopK: r.TopK,
		SampleRate: r.SampleRate, Seed: r.Seed,
	}
}

func (r AudioRequest) ToSpeechRequest() AudioSpeechRequest {
	return AudioSpeechRequest{
		Prompt: r.Prompt, Model: r.Model, Duration: r.Duration, SampleRate: r.SampleRate,
		Voice: r.Voice, Language: r.Language, Speed: r.Speed, Pitch: r.Pitch,
		Emotion: r.Emotion, Style: r.Style, Energy: r.Energy, Seed: r.Seed,
	}
}
