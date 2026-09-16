package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *Driver) CreateVoice(ctx context.Context, req driver.VoiceCreateRequest) (*driver.VoiceResponse, error) {
	if d.entry.Voice == nil {
		return nil, fmt.Errorf("driver %q does not support voice", d.id)
	}
	spec := config.ResolveAPICapabilitySpec(d.entry.Voice, string(driver.VoiceTaskZeroShot))
	raw, err := doHTTP(d.scoped(ctx), spec, voiceCreateHTTPVars(req, spec))
	if err != nil {
		return nil, err
	}
	return voiceResponseFromHTTP(spec, raw, req.Name)
}

func (d *Driver) ConvertVoice(ctx context.Context, req driver.VoiceConvertRequest) (*driver.VoiceResponse, error) {
	if d.entry.Voice == nil {
		return nil, fmt.Errorf("driver %q does not support voice", d.id)
	}
	spec := config.ResolveAPICapabilitySpec(d.entry.Voice, string(driver.VoiceTaskConversion))
	raw, err := doHTTP(d.scoped(ctx), spec, voiceConvertHTTPVars(req, spec))
	if err != nil {
		return nil, err
	}
	resp, err := voiceResponseFromHTTP(spec, raw, req.Name)
	if err != nil {
		return nil, err
	}
	if resp.OutputPath == "" {
		path, err := mediaPathFromHTTP(spec, raw, "audio")
		if err == nil {
			resp.OutputPath = path
		}
	}
	return resp, nil
}

func (d *Driver) ListVoiceProfiles(ctx context.Context) ([]driver.VoiceProfileInfo, error) {
	if d.entry.Voice == nil {
		return nil, fmt.Errorf("driver %q does not support voice", d.id)
	}
	spec := config.ResolveAPICapabilitySpec(d.entry.Voice, "list")
	if spec.URL == "" && spec.BaseURL == "" {
		return nil, fmt.Errorf("driver %q: voice list requires tasks.list or url in voice config", d.id)
	}
	raw, err := doHTTP(d.scoped(ctx), spec, map[string]string{"task": "list"})
	if err != nil {
		return nil, err
	}
	return voiceProfilesFromHTTP(spec, raw)
}

func voiceResponseFromHTTP(spec *config.APICapabilitySpec, raw []byte, fallbackName string) (*driver.VoiceResponse, error) {
	resp := &driver.VoiceResponse{Name: fallbackName}
	if spec.Response == nil {
		return resp, nil
	}
	if p := spec.Response["voice_id"]; p != "" {
		resp.VoiceID, _ = extractJSONPath(raw, p)
	}
	if p := spec.Response["name"]; p != "" {
		resp.Name, _ = extractJSONPath(raw, p)
	}
	if p := spec.Response["output_path"]; p != "" {
		resp.OutputPath, _ = extractJSONPath(raw, p)
	}
	return resp, nil
}

func voiceProfilesFromHTTP(spec *config.APICapabilitySpec, raw []byte) ([]driver.VoiceProfileInfo, error) {
	path := "$"
	if spec.Response != nil && spec.Response["profiles"] != "" {
		path = spec.Response["profiles"]
	}
	text, err := extractJSONPath(raw, path)
	if err != nil {
		return nil, err
	}
	var profiles []driver.VoiceProfileInfo
	if err := json.Unmarshal([]byte(text), &profiles); err != nil {
		// try parsing whole body as array
		if err2 := json.Unmarshal(raw, &profiles); err2 != nil {
			return nil, fmt.Errorf("parse voice profiles: %w", err)
		}
	}
	return profiles, nil
}
