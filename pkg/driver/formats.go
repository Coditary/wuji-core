package driver

import (
	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/modelformat"
)

// CapabilityFormats maps a capability to the model/source formats a driver accepts.
type CapabilityFormats map[capability.Type][]modelformat.Type

// FormatsFor returns supported formats for a capability, or nil.
func (f CapabilityFormats) FormatsFor(cap capability.Type) []modelformat.Type {
	if f == nil {
		return nil
	}
	return f[cap]
}

func CapabilityFormatsToProto(f CapabilityFormats) []*wujiv1.CapabilityFormatSupport {
	return capabilityFormatsToProto(f)
}

func capabilityFormatsToProto(f CapabilityFormats) []*wujiv1.CapabilityFormatSupport {
	if len(f) == 0 {
		return nil
	}
	out := make([]*wujiv1.CapabilityFormatSupport, 0, len(f))
	for cap, formats := range f {
		if len(formats) == 0 {
			continue
		}
		names := make([]string, 0, len(formats))
		for _, format := range formats {
			names = append(names, format.String())
		}
		out = append(out, &wujiv1.CapabilityFormatSupport{
			Capability: cap.String(),
			Formats:    names,
		})
	}
	return out
}

func CapabilityFormatsFromProto(items []*wujiv1.CapabilityFormatSupport) CapabilityFormats {
	return capabilityFormatsFromProto(items)
}

func capabilityFormatsFromProto(items []*wujiv1.CapabilityFormatSupport) CapabilityFormats {
	if len(items) == 0 {
		return nil
	}
	out := make(CapabilityFormats, len(items))
	for _, item := range items {
		if item == nil || item.GetCapability() == "" || len(item.GetFormats()) == 0 {
			continue
		}
		formats := make([]modelformat.Type, 0, len(item.GetFormats()))
		for _, name := range item.GetFormats() {
			formats = append(formats, modelformat.Type(name))
		}
		out[capability.Type(item.GetCapability())] = formats
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
