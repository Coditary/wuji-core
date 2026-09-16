package api

import (
	"sort"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

func configuredImageTasks(spec *config.APICapabilitySpec) []driver.ImageTask {
	if spec == nil {
		return nil
	}
	if len(spec.Tasks) > 0 {
		return parseImageTasks(spec.Tasks)
	}
	return driver.AllImageTasks()
}

func configuredVideoTasks(spec *config.APICapabilitySpec) []driver.VideoTask {
	if spec == nil {
		return nil
	}
	if len(spec.Tasks) > 0 {
		return parseVideoTasks(spec.Tasks)
	}
	return driver.AllVideoTasks()
}

func configuredAudioTasks(spec *config.APICapabilitySpec) []driver.AudioTask {
	if spec == nil {
		return nil
	}
	if len(spec.Tasks) > 0 {
		return parseAudioTasks(spec.Tasks)
	}
	return driver.AllAudioTasks()
}

func configuredMeshTasks(spec *config.APICapabilitySpec) []driver.MeshTask {
	if spec == nil {
		return nil
	}
	if len(spec.Tasks) > 0 {
		return parseMeshTasks(spec.Tasks)
	}
	return driver.AllMeshTasks()
}

func configuredVoiceTasks(spec *config.APICapabilitySpec) []driver.VoiceTask {
	if spec == nil {
		return nil
	}
	if len(spec.Tasks) > 0 {
		return parseVoiceTasks(spec.Tasks)
	}
	return driver.AllVoiceTasks()
}

func configuredDataTasks(spec *config.APICapabilitySpec) []driver.DataTask {
	if spec == nil {
		return nil
	}
	if len(spec.Tasks) > 0 {
		return parseDataTasks(spec.Tasks)
	}
	return driver.AllDataTasks()
}

func configuredRAGTasks(spec *config.APICapabilitySpec) []driver.RAGTask {
	if spec == nil {
		return nil
	}
	if len(spec.Tasks) > 0 {
		return parseRAGTasks(spec.Tasks)
	}
	return driver.AllRAGTasks()
}

func configuredDatasetTasks(spec *config.APICapabilitySpec) []driver.DatasetTask {
	if spec == nil {
		return nil
	}
	if len(spec.Tasks) > 0 {
		return parseDatasetTasks(spec.Tasks)
	}
	return driver.AllDatasetTasks()
}

func parseImageTasks(tasks map[string]*config.APICapabilitySpec) []driver.ImageTask {
	out := make([]driver.ImageTask, 0, len(tasks))
	for name := range tasks {
		t, err := driver.ParseImageTask(name)
		if err == nil {
			out = append(out, t)
		}
	}
	sortImageTasks(out)
	return out
}

func parseVideoTasks(tasks map[string]*config.APICapabilitySpec) []driver.VideoTask {
	out := make([]driver.VideoTask, 0, len(tasks))
	for name := range tasks {
		t, err := driver.ParseVideoTask(name)
		if err == nil {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func parseAudioTasks(tasks map[string]*config.APICapabilitySpec) []driver.AudioTask {
	out := make([]driver.AudioTask, 0, len(tasks))
	for name := range tasks {
		t, err := driver.ParseAudioTask(name)
		if err == nil {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func parseMeshTasks(tasks map[string]*config.APICapabilitySpec) []driver.MeshTask {
	out := make([]driver.MeshTask, 0, len(tasks))
	for name := range tasks {
		t, err := driver.ParseMeshTask(name)
		if err == nil {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func parseVoiceTasks(tasks map[string]*config.APICapabilitySpec) []driver.VoiceTask {
	out := make([]driver.VoiceTask, 0, len(tasks))
	for name := range tasks {
		t, err := driver.ParseVoiceTask(name)
		if err == nil {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func parseDataTasks(tasks map[string]*config.APICapabilitySpec) []driver.DataTask {
	out := make([]driver.DataTask, 0, len(tasks))
	for name := range tasks {
		t, err := driver.ParseDataTask(name)
		if err == nil {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func parseRAGTasks(tasks map[string]*config.APICapabilitySpec) []driver.RAGTask {
	out := make([]driver.RAGTask, 0, len(tasks))
	for name := range tasks {
		t, err := driver.ParseRAGTask(name)
		if err == nil {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func parseDatasetTasks(tasks map[string]*config.APICapabilitySpec) []driver.DatasetTask {
	out := make([]driver.DatasetTask, 0, len(tasks))
	for name := range tasks {
		t, err := driver.ParseDatasetTask(name)
		if err == nil {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func sortImageTasks(tasks []driver.ImageTask) {
	order := driver.AllImageTasks()
	sort.Slice(tasks, func(i, j int) bool {
		return taskIndex(order, tasks[i]) < taskIndex(order, tasks[j])
	})
}

func taskIndex[T comparable](order []T, item T) int {
	for i, v := range order {
		if v == item {
			return i
		}
	}
	return len(order)
}
