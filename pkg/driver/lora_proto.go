package driver

import wujiv1 "github.com/coditary/wuji-core/api/proto/v1"

func lorasFromProto(items []*wujiv1.LoraAttachment) []LoRARef {
	if len(items) == 0 {
		return nil
	}
	out := make([]LoRARef, 0, len(items))
	for _, item := range items {
		if item == nil || item.GetPath() == "" {
			continue
		}
		out = append(out, LoRARef{Path: item.GetPath(), Weight: item.GetWeight()})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func lorasToProto(items []LoRARef) []*wujiv1.LoraAttachment {
	if len(items) == 0 {
		return nil
	}
	out := make([]*wujiv1.LoraAttachment, 0, len(items))
	for _, item := range items {
		if item.Path == "" {
			continue
		}
		out = append(out, &wujiv1.LoraAttachment{Path: item.Path, Weight: item.LoRAWeight()})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
