package api

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (d *Driver) GenerateMesh(ctx context.Context, req driver.MeshRequest) (*driver.MeshResponse, error) {
	if d.entry.Mesh == nil {
		return nil, fmt.Errorf("driver %q does not support mesh", d.id)
	}
	task := req.TaskOrDefault()
	spec := config.ResolveAPICapabilitySpec(d.entry.Mesh, string(task))
	raw, err := doHTTP(d.scoped(ctx), spec, meshHTTPVars(req, spec))
	if err != nil {
		return nil, err
	}
	path, format, err := meshPathFromHTTP(spec, raw)
	if err != nil {
		return nil, err
	}
	if format == "" {
		format = req.Format
	}
	if format == "" {
		format = strings.TrimPrefix(filepath.Ext(path), ".")
	}
	return &driver.MeshResponse{
		Path: path, Format: format, Task: task,
		Mode: req.ModeOrDefault(), Representation: req.RepresentationOrDefault(),
	}, nil
}

func meshPathFromHTTP(spec *config.APICapabilitySpec, raw []byte) (path, format string, err error) {
	if spec.Response == nil {
		return "", "", fmt.Errorf("mesh http response mapping is required")
	}
	for _, key := range []string{"mesh_url", "url", "path"} {
		if p := spec.Response[key]; p != "" {
			val, err := extractJSONPath(raw, p)
			if err != nil {
				return "", "", err
			}
			if key == "mesh_url" || key == "url" {
				path, err = downloadToTemp(val, "mesh")
				return path, "", err
			}
			return val, "", nil
		}
	}
	if p := spec.Response["mesh_base64"]; p != "" {
		b64, err := extractJSONPath(raw, p)
		if err != nil {
			return "", "", err
		}
		path, err = writeBase64File(b64, "mesh")
		return path, "", err
	}
	if p := spec.Response["format"]; p != "" {
		format, _ = extractJSONPath(raw, p)
	}
	return "", format, fmt.Errorf("mesh response needs mesh_url, mesh_base64, url, or path")
}
