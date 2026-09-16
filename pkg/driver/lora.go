package driver

// LoRARef attaches a LoRA adapter to a base model at inference time.
type LoRARef struct {
	Path   string
	Weight float32 // 0 = default strength (1.0)
}

// LoRAWeight returns the effective mixing weight (defaults to 1.0).
func (l LoRARef) LoRAWeight() float32 {
	if l.Weight <= 0 {
		return 1
	}
	return l.Weight
}

// CloneLoRAs returns a deep copy of LoRA references.
func CloneLoRAs(in []LoRARef) []LoRARef {
	if len(in) == 0 {
		return nil
	}
	return append([]LoRARef(nil), in...)
}
