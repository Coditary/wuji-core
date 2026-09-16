package driver

import (
	"fmt"
	"strings"
)

// MeshRepresentation is an intermediate 3D output format before optional mesh export.
type MeshRepresentation string

const (
	MeshRepresentationMesh       MeshRepresentation = "mesh"
	MeshRepresentationNeRF       MeshRepresentation = "nerf"
	MeshRepresentationSplat      MeshRepresentation = "splat"
	MeshRepresentationPointCloud MeshRepresentation = "pointcloud"
)

func AllMeshRepresentations() []MeshRepresentation {
	return []MeshRepresentation{
		MeshRepresentationMesh,
		MeshRepresentationNeRF,
		MeshRepresentationSplat,
		MeshRepresentationPointCloud,
	}
}

func (r MeshRepresentation) String() string { return string(r) }

func (r MeshRepresentation) IsValid() bool {
	for _, known := range AllMeshRepresentations() {
		if r == known {
			return true
		}
	}
	return false
}

// ParseMeshRepresentation normalizes a representation name. Empty defaults to mesh.
func ParseMeshRepresentation(raw string) (MeshRepresentation, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	if normalized == "" {
		return MeshRepresentationMesh, nil
	}
	switch normalized {
	case "mesh", "glb":
		return MeshRepresentationMesh, nil
	case "nerf", "radiance-field":
		return MeshRepresentationNeRF, nil
	case "splat", "gaussian-splat", "gaussians":
		return MeshRepresentationSplat, nil
	case "pointcloud", "point-cloud", "pcd", "ply":
		return MeshRepresentationPointCloud, nil
	}
	rep := MeshRepresentation(normalized)
	if !rep.IsValid() {
		return "", fmt.Errorf("unknown mesh representation %q (valid: %s)", raw, joinMeshRepresentations())
	}
	return rep, nil
}

func (r MeshRequest) RepresentationOrDefault() MeshRepresentation {
	if r.Representation == "" {
		return MeshRepresentationMesh
	}
	return r.Representation
}

func meshTaskSupportsRepresentation(task MeshTask) bool {
	switch task {
	case MeshTaskGenerate, MeshTaskImageToMesh, MeshTaskMultiView,
		MeshTaskVideoToMesh, MeshTaskDepthToMesh, MeshTaskScene:
		return true
	default:
		return false
	}
}

func joinMeshRepresentations() string {
	parts := make([]string, len(AllMeshRepresentations()))
	for i, r := range AllMeshRepresentations() {
		parts[i] = string(r)
	}
	return strings.Join(parts, ", ")
}
