package api

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

func TestConfiguredImageTasksDefaultAll(t *testing.T) {
	spec := &config.APICapabilitySpec{Protocol: "http", URL: "https://example.com"}
	tasks := configuredImageTasks(spec)
	if len(tasks) != len(driver.AllImageTasks()) {
		t.Fatalf("got %d tasks", len(tasks))
	}
}

func TestConfiguredImageTasksFromYAML(t *testing.T) {
	spec := &config.APICapabilitySpec{
		Protocol: "http",
		Tasks: map[string]*config.APICapabilitySpec{
			"generate": {URL: "https://example.com/gen"},
			"img2img":  {URL: "https://example.com/i2i"},
		},
	}
	tasks := configuredImageTasks(spec)
	if len(tasks) != 2 {
		t.Fatalf("got %v", tasks)
	}
}

func TestDriverInfoAdvertisesImageTasks(t *testing.T) {
	d, err := New("test", config.APIEntry{
		Image: &config.APICapabilitySpec{
			Protocol: "http",
			Tasks:    map[string]*config.APICapabilitySpec{"generate": {URL: "https://x"}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Info().ImageTasks) != 1 {
		t.Fatalf("tasks=%v", d.Info().ImageTasks)
	}
}
