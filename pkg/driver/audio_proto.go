package driver

import wujiv1 "github.com/coditary/wuji-core/api/proto/v1"

func AudioTasksToProto(tasks []AudioTask) []string {
	if len(tasks) == 0 {
		return nil
	}
	out := make([]string, 0, len(tasks))
	for _, task := range tasks {
		out = append(out, string(task))
	}
	return out
}

func AudioTasksFromProto(items []string) []AudioTask {
	if len(items) == 0 {
		return nil
	}
	out := make([]AudioTask, 0, len(items))
	for _, item := range items {
		task := AudioTask(item)
		if task.IsValid() {
			out = append(out, task)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func audioCommonToProto(p AudioCommonParams) (string, float32, string, string, string, float32, float32, float32, float32, int32, int32, *int32) {
	var seed *int32
	if p.Seed != nil {
		v := int32(*p.Seed)
		seed = &v
	}
	return p.Prompt, p.Duration, p.Lyrics, p.NegativePrompt, p.Model, p.Overlap, p.Temperature, p.CFGScale,
		p.TopP, int32(p.TopK), int32(p.SampleRate), seed
}

func audioCommonFromProto(prompt string, duration float32, lyrics, negativePrompt, model string, overlap, temperature, cfgScale, topP float32, topK int32, sampleRate int32, seed *int32) AudioCommonParams {
	out := AudioCommonParams{
		Prompt: prompt, Duration: duration, Lyrics: lyrics, NegativePrompt: negativePrompt, Model: model,
		Overlap: overlap, Temperature: temperature, CFGScale: cfgScale, TopP: topP, TopK: int(topK), SampleRate: int(sampleRate),
	}
	if seed != nil {
		v := int(*seed)
		out.Seed = &v
	}
	return out
}

func AudioMusicRequestFromProto(req *wujiv1.GenerateMusicRequest) AudioMusicRequest {
	if req == nil {
		return AudioMusicRequest{}
	}
	return AudioMusicRequest{
		AudioCommonParams: audioCommonFromProto(
			req.GetPrompt(), req.GetDuration(), req.GetLyrics(), req.GetNegativePrompt(), req.GetModel(),
			req.GetOverlap(), req.GetTemperature(), req.GetCfgScale(), req.GetTopP(), req.GetTopK(),
			req.GetSampleRate(), req.Seed,
		),
	}
}

func AudioMusicRequestToProto(req AudioMusicRequest) *wujiv1.GenerateMusicRequest {
	prompt, duration, lyrics, neg, model, overlap, temp, cfg, topP, topK, sampleRate, seed := audioCommonToProto(req.AudioCommonParams)
	return &wujiv1.GenerateMusicRequest{
		Prompt: prompt, Duration: duration, Lyrics: lyrics, NegativePrompt: neg, Model: model,
		Overlap: overlap, Temperature: temp, CfgScale: cfg, TopP: topP, TopK: int32(topK), SampleRate: sampleRate, Seed: seed,
	}
}

func AudioMelodyToMusicRequestFromProto(req *wujiv1.MelodyToMusicRequest) AudioMelodyToMusicRequest {
	if req == nil {
		return AudioMelodyToMusicRequest{}
	}
	return AudioMelodyToMusicRequest{
		ReferencePath: req.GetReferencePath(),
		AudioCommonParams: audioCommonFromProto(
			req.GetPrompt(), req.GetDuration(), req.GetLyrics(), req.GetNegativePrompt(), req.GetModel(),
			req.GetOverlap(), req.GetTemperature(), req.GetCfgScale(), req.GetTopP(), req.GetTopK(),
			req.GetSampleRate(), req.Seed,
		),
	}
}

func AudioMelodyToMusicRequestToProto(req AudioMelodyToMusicRequest) *wujiv1.MelodyToMusicRequest {
	prompt, duration, lyrics, neg, model, overlap, temp, cfg, topP, topK, sampleRate, seed := audioCommonToProto(req.AudioCommonParams)
	return &wujiv1.MelodyToMusicRequest{
		Prompt: prompt, Duration: duration, Lyrics: lyrics, NegativePrompt: neg, Model: model,
		Overlap: overlap, Temperature: temp, CfgScale: cfg, TopP: topP, TopK: int32(topK), SampleRate: sampleRate,
		Seed: seed, ReferencePath: req.ReferencePath,
	}
}

func AudioSFXRequestFromProto(req *wujiv1.GenerateSFXRequest) AudioSFXRequest {
	if req == nil {
		return AudioSFXRequest{}
	}
	out := AudioSFXRequest{
		Prompt: req.GetPrompt(), NegativePrompt: req.GetNegativePrompt(), Model: req.GetModel(),
		Duration: req.GetDuration(), Temperature: req.GetTemperature(), CFGScale: req.GetCfgScale(),
		TopP: req.GetTopP(), TopK: int(req.GetTopK()), SampleRate: int(req.GetSampleRate()),
	}
	if req.Seed != nil {
		v := int(*req.Seed)
		out.Seed = &v
	}
	return out
}

func AudioSFXRequestToProto(req AudioSFXRequest) *wujiv1.GenerateSFXRequest {
	out := &wujiv1.GenerateSFXRequest{
		Prompt: req.Prompt, Duration: req.Duration, NegativePrompt: req.NegativePrompt, Model: req.Model,
		Temperature: req.Temperature, CfgScale: req.CFGScale, TopP: req.TopP, TopK: int32(req.TopK), SampleRate: int32(req.SampleRate),
	}
	if req.Seed != nil {
		v := int32(*req.Seed)
		out.Seed = &v
	}
	return out
}

func AudioSpeechRequestFromProto(req *wujiv1.GenerateAudioSpeechRequest) AudioSpeechRequest {
	if req == nil {
		return AudioSpeechRequest{}
	}
	out := AudioSpeechRequest{
		Prompt: req.GetPrompt(), Model: req.GetModel(), Duration: req.GetDuration(), SampleRate: int(req.GetSampleRate()),
		Voice: req.GetVoice(), Language: req.GetLanguage(), Speed: req.GetSpeed(), Pitch: int(req.GetPitch()),
		Emotion: req.GetEmotion(), Style: req.GetStyle(), Energy: req.GetEnergy(),
	}
	if req.Seed != nil {
		v := int(*req.Seed)
		out.Seed = &v
	}
	return out
}

func AudioSpeechRequestToProto(req AudioSpeechRequest) *wujiv1.GenerateAudioSpeechRequest {
	out := &wujiv1.GenerateAudioSpeechRequest{
		Prompt: req.Prompt, Duration: req.Duration, Model: req.Model, SampleRate: int32(req.SampleRate),
		Voice: req.Voice, Language: req.Language, Speed: req.Speed, Pitch: int32(req.Pitch),
		Emotion: req.Emotion, Style: req.Style, Energy: req.Energy,
	}
	if req.Seed != nil {
		v := int32(*req.Seed)
		out.Seed = &v
	}
	return out
}

func AudioResponseFromProto(resp *wujiv1.GenerateAudioResponse) *AudioResponse {
	if resp == nil {
		return &AudioResponse{}
	}
	return &AudioResponse{
		Path: resp.GetPath(), Duration: resp.GetDuration(), SampleRate: int(resp.GetSampleRate()), Format: resp.GetFormat(),
	}
}

func AudioResponseToProto(resp *AudioResponse) *wujiv1.GenerateAudioResponse {
	if resp == nil {
		return &wujiv1.GenerateAudioResponse{}
	}
	return &wujiv1.GenerateAudioResponse{
		Path: resp.Path, Duration: resp.Duration, SampleRate: int32(resp.SampleRate), Format: resp.Format,
	}
}
