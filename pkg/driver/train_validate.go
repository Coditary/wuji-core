package driver

import (
	"fmt"
	"strings"
)

func ValidateTextTrainRequest(req TextTrainRequest) error {
	method := req.Method
	if method == "" {
		method = TextTrainMethodFull
	}
	if _, err := ParseTextTrainMethod(string(method)); err != nil {
		return err
	}
	if strings.TrimSpace(req.DatasetID) == "" {
		return fmt.Errorf("--dataset is required")
	}
	if method == TextTrainMethodLoRA && req.LoRARank <= 0 {
		return fmt.Errorf("--lora-rank is required for lora training (or set a positive rank)")
	}
	if method == TextTrainMethodQLoRA && req.LoRARank <= 0 {
		return fmt.Errorf("--lora-rank is required for qlora training")
	}
	return nil
}

func ValidateImageTrainRequest(req ImageTrainRequest) error {
	method := req.Method
	if method == "" {
		method = ImageTrainMethodFull
	}
	if _, err := ParseImageTrainMethod(string(method)); err != nil {
		return err
	}
	if strings.TrimSpace(req.DatasetID) == "" {
		return fmt.Errorf("--dataset is required")
	}
	switch method {
	case ImageTrainMethodDreamBooth:
		if strings.TrimSpace(req.ClassToken) == "" {
			return fmt.Errorf("--class-token is required for dreambooth training")
		}
	case ImageTrainMethodTextualInversion:
		if strings.TrimSpace(req.Token) == "" {
			return fmt.Errorf("--token is required for textual-inversion training")
		}
	case ImageTrainMethodControlNet:
		if req.ControlType == "" {
			return fmt.Errorf("--control-type or a typed control flag (e.g. --pose) is required for controlnet training")
		}
	case ImageTrainMethodLoRA:
		if req.LoRARank <= 0 {
			return fmt.Errorf("--lora-rank is required for lora training (or set a positive rank)")
		}
	case ImageTrainMethodEmbedding:
		// dataset of reference images is sufficient
	case ImageTrainMethodFull:
		// no extra requirements
	}
	if req.Mode != "" {
		if _, err := ParseImageMode(string(req.Mode)); err != nil {
			return err
		}
	}
	return nil
}

func ValidateVideoTrainRequest(req VideoTrainRequest) error {
	method := req.Method
	if method == "" {
		method = VideoTrainMethodFull
	}
	if _, err := ParseVideoTrainMethod(string(method)); err != nil {
		return err
	}
	if strings.TrimSpace(req.DatasetID) == "" {
		return fmt.Errorf("--dataset is required")
	}
	if method == VideoTrainMethodLoRA && req.LoRARank <= 0 {
		return fmt.Errorf("--lora-rank is required for lora training")
	}
	return nil
}

func ValidateAudioTrainRequest(req AudioTrainRequest) error {
	method := req.Method
	if method == "" {
		method = AudioTrainMethodMusic
	}
	if _, err := ParseAudioTrainMethod(string(method)); err != nil {
		return err
	}
	if strings.TrimSpace(req.DatasetID) == "" {
		return fmt.Errorf("--dataset is required")
	}
	return nil
}

func ValidateMeshTrainRequest(req MeshTrainRequest) error {
	method := req.Method
	if method == "" {
		method = MeshTrainMethodFull
	}
	if _, err := ParseMeshTrainMethod(string(method)); err != nil {
		return err
	}
	if strings.TrimSpace(req.DatasetID) == "" {
		return fmt.Errorf("--dataset is required")
	}
	if method == MeshTrainMethodLoRA && req.LoRARank <= 0 {
		return fmt.Errorf("--lora-rank is required for lora training")
	}
	if req.Mode != "" {
		if _, err := ParseMeshMode(string(req.Mode)); err != nil {
			return err
		}
	}
	return nil
}

func ValidateVoiceTrainRequest(req VoiceTrainRequest) error {
	method := req.Method
	if method == "" {
		method = VoiceTrainMethodRVC
	}
	if _, err := ParseVoiceTrainMethod(string(method)); err != nil {
		return err
	}
	if strings.TrimSpace(req.DatasetID) == "" {
		return fmt.Errorf("--dataset is required")
	}
	return nil
}
