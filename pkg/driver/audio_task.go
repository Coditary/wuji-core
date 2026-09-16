package driver

import (
	"fmt"
	"strings"
)

// AudioTask identifies the audio generation mode.
type AudioTask string

// AudioTaskType is a deprecated alias for AudioTask.
type AudioTaskType = AudioTask

const (
	AudioTaskTextToMusic   AudioTask = "text-to-music"
	AudioTaskMelodyToMusic AudioTask = "melody-to-music"
	AudioTaskTextToSFX     AudioTask = "text-to-sfx"
	AudioTaskTextToSpeech  AudioTask = "text-to-speech"
)

type AudioTaskInfo struct {
	Task                 AudioTask
	Description          string
	RequiresPrompt       bool
	RequiresReference    bool
}

func AllAudioTasks() []AudioTask {
	return []AudioTask{
		AudioTaskTextToMusic,
		AudioTaskMelodyToMusic,
		AudioTaskTextToSFX,
		AudioTaskTextToSpeech,
	}
}

func AudioTaskCatalog() []AudioTaskInfo {
	return []AudioTaskInfo{
		{Task: AudioTaskTextToMusic, Description: "Generate music from a text prompt", RequiresPrompt: true},
		{Task: AudioTaskMelodyToMusic, Description: "Generate music guided by a reference melody", RequiresPrompt: true, RequiresReference: true},
		{Task: AudioTaskTextToSFX, Description: "Generate sound effects from a text prompt", RequiresPrompt: true},
		{Task: AudioTaskTextToSpeech, Description: "Speak text aloud (text-to-speech)", RequiresPrompt: true},
	}
}

func (t AudioTask) String() string { return string(t) }

func (t AudioTask) IsValid() bool {
	for _, known := range AllAudioTasks() {
		if t == known {
			return true
		}
	}
	return false
}

func ParseAudioTask(raw string) (AudioTask, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	switch normalized {
	case "", "text-to-music", "music", "txt2music":
		return AudioTaskTextToMusic, nil
	case "melody-to-music", "melody2music", "melody":
		return AudioTaskMelodyToMusic, nil
	case "text-to-sfx", "sfx", "sound-effects":
		return AudioTaskTextToSFX, nil
	case "text-to-speech", "tts", "speech", "speak":
		return AudioTaskTextToSpeech, nil
	}
	task := AudioTask(normalized)
	if !task.IsValid() {
		return "", fmt.Errorf("unknown audio task %q (valid: %s)", raw, joinAudioTasks())
	}
	return task, nil
}

// AudioTaskInputs collects CLI fields used to infer an audio task from flags.
type AudioTaskInputs struct {
	ReferencePath   string
	SFXRequested    bool
	SpeechRequested bool
	VoiceName       string
}

// InferAudioTask selects a task from input flags when --task is omitted.
func InferAudioTask(in AudioTaskInputs) (AudioTask, error) {
	hasReference := strings.TrimSpace(in.ReferencePath) != ""
	hasVoice := strings.TrimSpace(in.VoiceName) != ""
	speech := in.SpeechRequested || hasVoice

	if speech && in.SFXRequested {
		return "", fmt.Errorf("use either --voice/--speech for text-to-speech or --sfx for sound effects, not both")
	}
	if speech && hasReference {
		return "", fmt.Errorf("use either --voice/--speech for text-to-speech or --reference for melody-to-music, not both")
	}
	if hasReference && in.SFXRequested {
		return "", fmt.Errorf("use either --reference for melody-to-music or --sfx for sound effects, not both")
	}
	if speech {
		return AudioTaskTextToSpeech, nil
	}
	if hasReference {
		return AudioTaskMelodyToMusic, nil
	}
	if in.SFXRequested {
		return AudioTaskTextToSFX, nil
	}
	return AudioTaskTextToMusic, nil
}

func (r AudioRequest) TaskOrDefault() AudioTask {
	if r.Task != "" {
		return r.Task
	}
	if r.TaskType != "" {
		return AudioTask(r.TaskType)
	}
	return AudioTaskTextToMusic
}

func (r AudioRequest) Validate() error {
	task := r.TaskOrDefault()
	if !task.IsValid() {
		return fmt.Errorf("unknown audio task %q (valid: %s)", task, joinAudioTasks())
	}
	info := audioTaskInfo(task)
	if info.RequiresPrompt && strings.TrimSpace(r.Prompt) == "" {
		return fmt.Errorf("prompt is required for audio task %q", task)
	}
	if info.RequiresReference && r.ReferencePath == "" {
		return fmt.Errorf("--reference is required for audio task %q", task)
	}
	if err := r.validateOverlap(); err != nil {
		return err
	}
	if err := r.validateFormat(); err != nil {
		return err
	}
	if err := r.validateSpeechFields(); err != nil {
		return err
	}
	return nil
}

func (r AudioRequest) validateSpeechFields() error {
	task := r.TaskOrDefault()
	speechFieldsSet := strings.TrimSpace(r.Voice) != "" || strings.TrimSpace(r.Language) != "" ||
		strings.TrimSpace(r.Emotion) != "" || strings.TrimSpace(r.Style) != "" ||
		r.Speed != 0 || r.Pitch != 0 || r.Energy != 0

	if task == AudioTaskTextToSpeech {
		if r.Speed < 0 {
			return fmt.Errorf("--speed must be >= 0 for text-to-speech (got %g)", r.Speed)
		}
		if r.Speed > 0 && (r.Speed < 0.25 || r.Speed > 4) {
			return fmt.Errorf("--speed must be between 0.25 and 4.0 for text-to-speech (got %g)", r.Speed)
		}
		if r.Pitch < -12 || r.Pitch > 12 {
			return fmt.Errorf("--pitch must be between -12 and +12 semitones (got %d)", r.Pitch)
		}
		if r.Energy < 0 || r.Energy > 1 {
			return fmt.Errorf("--energy must be between 0.0 and 1.0 for text-to-speech (got %g)", r.Energy)
		}
		return nil
	}
	if speechFieldsSet {
		return fmt.Errorf("speech flags (--voice, --lang, --speed, --pitch, --emotion, --style, --energy) are only valid for text-to-speech, not %q", task)
	}
	return nil
}

func audioTaskInfo(task AudioTask) AudioTaskInfo {
	for _, info := range AudioTaskCatalog() {
		if info.Task == task {
			return info
		}
	}
	return AudioTaskInfo{}
}

func joinAudioTasks() string {
	parts := make([]string, len(AllAudioTasks()))
	for i, t := range AllAudioTasks() {
		parts[i] = string(t)
	}
	return strings.Join(parts, ", ")
}
