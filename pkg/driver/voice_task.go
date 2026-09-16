package driver

import (
	"fmt"
	"strings"
)

// VoiceTask identifies the voice cloning workflow.
type VoiceTask string

// VoiceCloneMode is a deprecated alias for VoiceTask values.
type VoiceCloneMode = VoiceTask

const (
	VoiceTaskZeroShot   VoiceTask = "zero-shot"
	VoiceTaskConversion VoiceTask = "conversion"
)

const (
	VoiceCloneModeZeroShot   = VoiceTaskZeroShot
	VoiceCloneModeConversion = VoiceTaskConversion
)

type VoiceTaskInfo struct {
	Task              VoiceTask
	Description       string
	RequiresName      bool
	RequiresSample    bool
	RequiresSource    bool
}

func AllVoiceTasks() []VoiceTask {
	return []VoiceTask{VoiceTaskZeroShot, VoiceTaskConversion}
}

func VoiceTaskCatalog() []VoiceTaskInfo {
	return []VoiceTaskInfo{
		{Task: VoiceTaskZeroShot, Description: "Create a voice profile from a reference sample", RequiresName: true, RequiresSample: true},
		{Task: VoiceTaskConversion, Description: "Convert source audio to a cloned voice", RequiresName: true, RequiresSample: true, RequiresSource: true},
	}
}

func (t VoiceTask) String() string { return string(t) }

func (t VoiceTask) IsValid() bool {
	for _, known := range AllVoiceTasks() {
		if t == known {
			return true
		}
	}
	return false
}

func ParseVoiceTask(raw string) (VoiceTask, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	switch normalized {
	case "", "zero-shot", "zeroshot", "clone", "create":
		return VoiceTaskZeroShot, nil
	case "conversion", "convert", "vc":
		return VoiceTaskConversion, nil
	}
	task := VoiceTask(normalized)
	if !task.IsValid() {
		return "", fmt.Errorf("unknown voice task %q (valid: %s)", raw, joinVoiceTasks())
	}
	return task, nil
}

func (r VoiceRequest) TaskOrDefault() VoiceTask {
	if r.Task != "" {
		return r.Task
	}
	if r.Mode != "" {
		return VoiceTask(r.Mode)
	}
	return VoiceTaskZeroShot
}

func (r VoiceRequest) Validate() error {
	task := r.TaskOrDefault()
	if !task.IsValid() {
		return fmt.Errorf("unknown voice task %q (valid: %s)", task, joinVoiceTasks())
	}
	info := voiceTaskInfo(task)
	if info.RequiresName && strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("voice name is required for task %q", task)
	}
	if info.RequiresSample && r.SamplePath == "" {
		return fmt.Errorf("--sample is required for voice task %q", task)
	}
	if info.RequiresSource && r.SourcePath == "" {
		return fmt.Errorf("--source is required for voice task %q", task)
	}
	return nil
}

func voiceTaskInfo(task VoiceTask) VoiceTaskInfo {
	for _, info := range VoiceTaskCatalog() {
		if info.Task == task {
			return info
		}
	}
	return VoiceTaskInfo{}
}

func joinVoiceTasks() string {
	parts := make([]string, len(AllVoiceTasks()))
	for i, t := range AllVoiceTasks() {
		parts[i] = string(t)
	}
	return strings.Join(parts, ", ")
}
