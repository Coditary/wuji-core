package core

import (
	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/driver"
	localdrv "github.com/coditary/wuji-core/pkg/driver/local"
)

// autoTTSOrder defines preferred brand aliases for text-to-speech auto-routing.
var autoTTSOrder = []string{
	"piper",
	"kokoro",
	"coqui",
	"melotts",
	"fishspeech",
	"chatterbox",
	"orpheus",
	"bark",
	"parler",
	"f5tts",
	"cosyvoice",
	"chattts",
	"styletts2",
	"gptsovits",
	"zonos",
}

// ResolveAudioDriverID picks the driver for an audio request.
// Brand names (piper, kokoro, coqui, …) resolve to runtime drivers (onnx-tts, pytorch-tts)
// when those are registered; otherwise the brand driver process is used directly.
func (c *Core) ResolveAudioDriverID(preferredID string, req *driver.AudioRequest) string {
	if preferredID != "" {
		id := localdrv.NormalizeDriverID(preferredID)
		if req != nil {
			if resolved, _ := driver.ApplyAudioBrandAlias(id, req, c.isDriverRegistered); resolved != id {
				return resolved
			}
		}
		return id
	}

	task := req.TaskOrDefault()
	if task == driver.AudioTaskTextToSpeech {
		for _, brandID := range autoTTSOrder {
			if req != nil {
				tmp := *req
				if resolved, ok := driver.ApplyAudioBrandAlias(brandID, &tmp, c.isDriverRegistered); ok {
					*req = tmp
					return resolved
				}
			}
			d, err := c.registry.Get(brandID)
			if err != nil {
				continue
			}
			if driver.HasAudioTask(d, task) {
				return brandID
			}
		}
		// Runtime drivers without a brand alias match.
		for _, runtimeID := range []string{driver.AudioRuntimeONNX, driver.AudioRuntimePyTorch} {
			d, err := c.registry.Get(runtimeID)
			if err != nil {
				continue
			}
			if driver.HasAudioTask(d, task) {
				return runtimeID
			}
		}
	}

	return c.ResolveDriverID("", capability.AudioGeneration)
}

func (c *Core) isDriverRegistered(id string) bool {
	_, err := c.registry.Get(id)
	return err == nil
}
