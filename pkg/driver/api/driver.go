package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

// Driver is a YAML-configured HTTP API backend.
type Driver struct {
	id    string
	entry config.APIEntry
	caps  []capability.Type
}

// New creates an API driver from config.
func New(id string, entry config.APIEntry) (*Driver, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("api driver id is required")
	}
	if len(entry.Capabilities()) == 0 {
		return nil, fmt.Errorf("api %q: at least one capability block is required", id)
	}
	d := &Driver{id: id, entry: entry}
	for _, name := range entry.Capabilities() {
		c, err := capability.Parse(name)
		if err != nil {
			return nil, fmt.Errorf("api %q: %w", id, err)
		}
		d.caps = appendUniqueCap(d.caps, c)
	}
	d.caps = augmentCapabilities(entry, d.caps)
	return d, nil
}

func appendUniqueCap(caps []capability.Type, c capability.Type) []capability.Type {
	for _, existing := range caps {
		if existing == c {
			return caps
		}
	}
	return append(caps, c)
}

func augmentCapabilities(entry config.APIEntry, caps []capability.Type) []capability.Type {
	if entry.Image != nil {
		for _, task := range configuredImageTasks(entry.Image) {
			if task == driver.ImageTaskUpscale {
				caps = appendUniqueCap(caps, capability.ImageUpscale)
				caps = appendUniqueCap(caps, capability.ImageDownscale)
				caps = appendUniqueCap(caps, capability.ImageScale)
				break
			}
		}
	}
	return caps
}

func (d *Driver) Info() driver.Info {
	info := driver.Info{
		ID:           d.id,
		Name:         "API " + d.id,
		Version:      "0.1.0",
		Description:  "YAML-configured HTTP API backend",
		Capabilities: d.caps,
		Remote:       false,
	}
	if d.entry.Image != nil {
		info.ImageTasks = configuredImageTasks(d.entry.Image)
	}
	if d.entry.Video != nil {
		info.VideoTasks = configuredVideoTasks(d.entry.Video)
	}
	if d.entry.Audio != nil {
		info.AudioTasks = configuredAudioTasks(d.entry.Audio)
	}
	if d.entry.Mesh != nil {
		info.MeshTasks = configuredMeshTasks(d.entry.Mesh)
	}
	if d.entry.Voice != nil {
		info.VoiceTasks = configuredVoiceTasks(d.entry.Voice)
	}
	if d.entry.Data != nil {
		info.DataTasks = configuredDataTasks(d.entry.Data)
	}
	if d.entry.RAG != nil {
		info.RAGTasks = configuredRAGTasks(d.entry.RAG)
	}
	if d.entry.Dataset != nil {
		info.DatasetTasks = configuredDatasetTasks(d.entry.Dataset)
	}
	return info
}

func (d *Driver) Capabilities() []capability.Type { return d.caps }

func (d *Driver) Close() error { return nil }

func (d *Driver) scoped(ctx context.Context) context.Context {
	return WithProviderID(ctx, d.id)
}

func (d *Driver) GenerateText(ctx context.Context, req driver.TextRequest) (*driver.TextResponse, error) {
	if d.entry.Text == nil {
		return nil, fmt.Errorf("driver %q does not support text", d.id)
	}
	if len(req.Messages) == 0 && strings.TrimSpace(req.Prompt) == "" && req.InputModeOrDefault() == driver.TextInputPrompt {
		return nil, fmt.Errorf("prompt or messages required")
	}
	return completeText(d.scoped(ctx), d.entry.Text, req)
}

// RegisterAll registers YAML api drivers from config.
func RegisterAll(reg *driver.Registry, apis map[string]config.APIEntry) error {
	for id, entry := range apis {
		drv, err := New(id, entry)
		if err != nil {
			return err
		}
		if err := reg.Register(drv); err != nil {
			return err
		}
	}
	return nil
}
