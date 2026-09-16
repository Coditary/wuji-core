package driver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coditary/wuji-core/pkg/media"
)

// DefaultAudioSegmentDuration is the per-chunk length when stitching with --overlap.
const DefaultAudioSegmentDuration float32 = 10

// PlanAudioSegments returns how many generation chunks are needed for a target duration.
func PlanAudioSegments(totalDuration, overlap, segmentDuration float32) (int, error) {
	if totalDuration <= 0 {
		return 1, nil
	}
	if segmentDuration <= 0 {
		segmentDuration = DefaultAudioSegmentDuration
	}
	if overlap <= 0 || totalDuration <= segmentDuration {
		return 1, nil
	}
	if overlap >= segmentDuration {
		return 0, fmt.Errorf("--overlap must be less than %g seconds (segment length)", segmentDuration)
	}
	advance := segmentDuration - overlap
	n := 1
	covered := segmentDuration
	for covered < totalDuration {
		n++
		covered += advance
	}
	return n, nil
}

func audioTaskSupportsOverlap(task AudioTask) bool {
	return task == AudioTaskTextToMusic || task == AudioTaskMelodyToMusic
}

func runAudioWithOverlap(ctx context.Context, d Driver, req AudioRequest, task AudioTask, ffmpegBin string) (*AudioResponse, error) {
	segmentDuration := DefaultAudioSegmentDuration
	if req.Duration > 0 && req.Duration < segmentDuration {
		segmentDuration = req.Duration
	}

	segments, err := PlanAudioSegments(req.Duration, req.Overlap, segmentDuration)
	if err != nil {
		return nil, err
	}
	if segments <= 1 {
		return dispatchAudioTask(ctx, d, req, task)
	}

	chunkPaths := make([]string, 0, segments)
	sampleRate := req.SampleRate
	for i := 0; i < segments; i++ {
		chunkReq := req
		chunkReq.Duration = segmentDuration
		chunkReq.Overlap = 0
		if chunkReq.Seed != nil {
			seed := *chunkReq.Seed + i
			chunkReq.Seed = &seed
		}
		resp, err := dispatchAudioTask(ctx, d, chunkReq, task)
		if err != nil {
			return nil, err
		}
		if sampleRate <= 0 && resp.SampleRate > 0 {
			sampleRate = resp.SampleRate
		}
		chunkPaths = append(chunkPaths, resp.Path)
	}
	if sampleRate <= 0 {
		sampleRate = 44100
	}

	stitchedFile, err := os.CreateTemp("/tmp/wuji", "wuji-audio-stitch-*.wav")
	if err != nil {
		return nil, fmt.Errorf("create stitched audio file: %w", err)
	}
	stitchedPath := stitchedFile.Name()
	_ = stitchedFile.Close()
	_ = os.Remove(stitchedPath)

	if err := media.ConcatAudioWithCrossfade(ctx, ffmpegBin, chunkPaths, req.Overlap, stitchedPath); err != nil {
		return nil, err
	}

	finalPath := stitchedPath
	if req.Duration > 0 {
		trimFile, err := os.CreateTemp("/tmp/wuji", "wuji-audio-trim-*.wav")
		if err != nil {
			_ = os.Remove(stitchedPath)
			return nil, err
		}
		trimPath := trimFile.Name()
		_ = trimFile.Close()
		_ = os.Remove(trimPath)
		if err := media.TrimAudio(ctx, ffmpegBin, stitchedPath, trimPath, req.Duration); err != nil {
			_ = os.Remove(stitchedPath)
			return nil, err
		}
		_ = os.Remove(stitchedPath)
		finalPath = trimPath
	}

	native := media.AudioFormatFromPath(finalPath)
	if native == "" {
		native = "wav"
	}
	return &AudioResponse{
		Path:       finalPath,
		Duration:   req.Duration,
		SampleRate: sampleRate,
		Format:     native,
	}, nil
}

func dispatchAudioTask(ctx context.Context, d Driver, req AudioRequest, task AudioTask) (*AudioResponse, error) {
	switch task {
	case AudioTaskTextToMusic:
		gen, ok := d.(MusicGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement MusicGenerator", d.Info().ID)
		}
		return gen.GenerateMusic(ctx, req.ToMusicRequest())
	case AudioTaskMelodyToMusic:
		gen, ok := d.(MelodyToMusicGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement MelodyToMusicGenerator", d.Info().ID)
		}
		return gen.MelodyToMusic(ctx, req.ToMelodyToMusicRequest())
	case AudioTaskTextToSFX:
		gen, ok := d.(SFXGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement SFXGenerator", d.Info().ID)
		}
		return gen.GenerateSFX(ctx, req.ToSFXRequest())
	case AudioTaskTextToSpeech:
		gen, ok := d.(AudioSpeechGenerator)
		if !ok {
			return nil, fmt.Errorf("driver %q does not implement AudioSpeechGenerator", d.Info().ID)
		}
		return gen.GenerateAudioSpeech(ctx, req.ToSpeechRequest())
	default:
		return nil, fmt.Errorf("unsupported audio task %q", task)
	}
}

func finalizeAudioOutput(ctx context.Context, req AudioRequest, resp *AudioResponse, ffmpegBin string) (*AudioResponse, error) {
	if resp == nil {
		return nil, fmt.Errorf("audio response is nil")
	}

	target := strings.TrimSpace(req.Format)
	native := strings.TrimSpace(resp.Format)
	if native == "" {
		native = media.AudioFormatFromPath(resp.Path)
	}
	if native == "" {
		native = "wav"
	}
	if target == "" {
		target = native
	} else {
		parsed, err := media.ParseAudioFormat(target)
		if err != nil {
			return nil, err
		}
		target = parsed
	}

	if strings.EqualFold(target, native) {
		resp.Format = target
		return resp, nil
	}

	ext := media.AudioFormatExtension(target)
	base := strings.TrimSuffix(resp.Path, filepath.Ext(resp.Path))
	outPath := base + "." + ext
	if err := media.TranscodeAudio(ctx, ffmpegBin, resp.Path, outPath, target, resp.SampleRate); err != nil {
		return nil, err
	}
	if outPath != resp.Path {
		_ = os.Remove(resp.Path)
	}
	resp.Path = outPath
	resp.Format = target
	return resp, nil
}
