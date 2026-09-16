package api

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *Driver) runAudio(ctx context.Context, req driver.AudioRequest) (*driver.AudioResponse, error) {
	if d.entry.Audio == nil {
		return nil, fmt.Errorf("driver %q does not support audio", d.id)
	}
	task := req.TaskOrDefault()
	spec := config.ResolveAPICapabilitySpec(d.entry.Audio, string(task))
	raw, err := doHTTP(d.scoped(ctx), spec, audioHTTPVars(req, spec))
	if err != nil {
		return nil, err
	}
	path, err := mediaPathFromHTTP(spec, raw, "audio")
	if err != nil {
		return nil, err
	}
	return &driver.AudioResponse{Path: path, Format: filepath.Ext(path)}, nil
}

func (d *Driver) GenerateMusic(ctx context.Context, req driver.AudioMusicRequest) (*driver.AudioResponse, error) {
	union := driver.AudioRequest{
		Task: driver.AudioTaskTextToMusic, Prompt: req.Prompt, Lyrics: req.Lyrics,
		NegativePrompt: req.NegativePrompt, Model: req.Model, Duration: req.Duration,
		Overlap: req.Overlap, Temperature: req.Temperature, CFGScale: req.CFGScale,
		TopP: req.TopP, TopK: req.TopK, SampleRate: req.SampleRate, Seed: req.Seed,
	}
	return d.runAudio(ctx, union)
}

func (d *Driver) MelodyToMusic(ctx context.Context, req driver.AudioMelodyToMusicRequest) (*driver.AudioResponse, error) {
	union := driver.AudioRequest{
		Task: driver.AudioTaskMelodyToMusic, Prompt: req.Prompt, Lyrics: req.Lyrics,
		NegativePrompt: req.NegativePrompt, Model: req.Model, Duration: req.Duration,
		Overlap: req.Overlap, Temperature: req.Temperature, CFGScale: req.CFGScale,
		TopP: req.TopP, TopK: req.TopK, SampleRate: req.SampleRate, Seed: req.Seed,
		ReferencePath: req.ReferencePath,
	}
	return d.runAudio(ctx, union)
}

func (d *Driver) GenerateSFX(ctx context.Context, req driver.AudioSFXRequest) (*driver.AudioResponse, error) {
	union := driver.AudioRequest{
		Task: driver.AudioTaskTextToSFX, Prompt: req.Prompt, NegativePrompt: req.NegativePrompt,
		Model: req.Model, Duration: req.Duration, Temperature: req.Temperature, CFGScale: req.CFGScale,
		TopP: req.TopP, TopK: req.TopK, SampleRate: req.SampleRate, Seed: req.Seed,
	}
	return d.runAudio(ctx, union)
}

func (d *Driver) GenerateAudioSpeech(ctx context.Context, req driver.AudioSpeechRequest) (*driver.AudioResponse, error) {
	if d.entry.Audio == nil {
		return nil, fmt.Errorf("driver %q does not support audio", d.id)
	}
	spec := config.ResolveAPICapabilitySpec(d.entry.Audio, string(driver.AudioTaskTextToSpeech))
	raw, err := doHTTP(d.scoped(ctx), spec, audioSpeechHTTPVars(req, spec))
	if err != nil {
		return nil, err
	}
	path, err := mediaPathFromHTTP(spec, raw, "audio")
	if err != nil {
		return nil, err
	}
	return &driver.AudioResponse{Path: path, Format: filepath.Ext(path)}, nil
}
