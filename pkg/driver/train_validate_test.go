package driver

import "testing"

func TestValidateImageTrainRequestDreamBooth(t *testing.T) {
	err := ValidateImageTrainRequest(ImageTrainRequest{
		Method: ImageTrainMethodDreamBooth, DatasetID: "ds", Name: "hero",
	})
	if err == nil {
		t.Fatal("expected class-token error")
	}
}

func TestValidateImageTrainRequestLoRA(t *testing.T) {
	err := ValidateImageTrainRequest(ImageTrainRequest{
		Method: ImageTrainMethodLoRA, DatasetID: "ds", Name: "style", LoRARank: 16,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidateTextTrainRequestQLoRA(t *testing.T) {
	err := ValidateTextTrainRequest(TextTrainRequest{
		Method: TextTrainMethodQLoRA, DatasetID: "ds", LoRARank: 8,
	})
	if err != nil {
		t.Fatal(err)
	}
}
