package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *Driver) ImageToText(ctx context.Context, req driver.TextRequest) (*driver.TextResponse, error) {
	return d.runTextMedia(ctx, "image-to-text", req)
}

func (d *Driver) VideoToText(ctx context.Context, req driver.TextRequest) (*driver.TextResponse, error) {
	return d.runTextMedia(ctx, "video-to-text", req)
}

func (d *Driver) AudioToText(ctx context.Context, req driver.TextRequest) (*driver.TextResponse, error) {
	return d.runTextMedia(ctx, "audio-to-text", req)
}

func (d *Driver) DocumentToText(ctx context.Context, req driver.TextRequest) (*driver.TextResponse, error) {
	return d.runTextMedia(ctx, "document-to-text", req)
}

func (d *Driver) runTextMedia(ctx context.Context, taskKey string, req driver.TextRequest) (*driver.TextResponse, error) {
	if d.entry.Text == nil {
		return nil, fmt.Errorf("driver %q does not support text", d.id)
	}
	spec := config.ResolveAPICapabilitySpec(d.entry.Text, taskKey)
	proto := strings.ToLower(strings.TrimSpace(spec.Protocol))
	if proto == "openai" || proto == "anthropic" {
		return completeText(ctx, spec, req)
	}
	raw, err := doHTTP(d.scoped(ctx), spec, textMediaHTTPVars(req, taskKey, spec))
	if err != nil {
		return nil, err
	}
	return textResponseFromHTTP(spec, raw)
}

func textResponseFromHTTP(spec *config.APICapabilitySpec, raw []byte) (*driver.TextResponse, error) {
	out := &driver.TextResponse{}
	if spec.Response != nil {
		if p := spec.Response["text"]; p != "" {
			out.Text, _ = extractJSONPath(raw, p)
		}
		if p := spec.Response["finish_reason"]; p != "" {
			out.FinishReason, _ = extractJSONPath(raw, p)
		}
	}
	if out.Text == "" {
		out.Text = strings.TrimSpace(string(raw))
	}
	return out, nil
}
