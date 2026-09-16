package remotemedia

import (
	"path/filepath"
	"strings"

	"github.com/coditary/wuji-core/pkg/driver"
)

func appendPath(paths []string, p string) []string {
	p = strings.TrimSpace(p)
	if p == "" || isRemoteURL(p) {
		return paths
	}
	return append(paths, p)
}

func appendPaths(paths []string, ps []string) []string {
	for _, p := range ps {
		paths = appendPath(paths, p)
	}
	return paths
}

func isRemoteURL(p string) bool {
	lower := strings.ToLower(p)
	return strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")
}

func appendLoRAs(paths []string, loras []driver.LoRARef) []string {
	for _, l := range loras {
		paths = appendPath(paths, l.Path)
	}
	return paths
}

func CollectTextRequest(req driver.TextRequest) []string {
	paths := appendPath(nil, req.MediaPath)
	return appendLoRAs(paths, req.LoRAs)
}

func RewriteTextRequest(req *driver.TextRequest, mapping map[string]string) {
	req.MediaPath = remap(req.MediaPath, mapping)
	for i := range req.LoRAs {
		req.LoRAs[i].Path = remap(req.LoRAs[i].Path, mapping)
	}
}

func CollectImageRequest(req driver.ImageRequest) []string {
	paths := appendPath(nil, req.InitImagePath)
	paths = appendPath(paths, req.MaskImagePath)
	paths = appendPath(paths, req.ControlImagePath)
	paths = appendPath(paths, req.StyleImagePath)
	paths = appendPath(paths, req.ReferenceImagePath)
	for _, u := range req.ControlUnits {
		paths = appendPath(paths, u.ImagePath)
	}
	return appendLoRAs(paths, req.LoRAs)
}

func RewriteImageRequest(req *driver.ImageRequest, mapping map[string]string) {
	req.InitImagePath = remap(req.InitImagePath, mapping)
	req.MaskImagePath = remap(req.MaskImagePath, mapping)
	req.ControlImagePath = remap(req.ControlImagePath, mapping)
	req.StyleImagePath = remap(req.StyleImagePath, mapping)
	req.ReferenceImagePath = remap(req.ReferenceImagePath, mapping)
	for i := range req.ControlUnits {
		req.ControlUnits[i].ImagePath = remap(req.ControlUnits[i].ImagePath, mapping)
	}
	for i := range req.LoRAs {
		req.LoRAs[i].Path = remap(req.LoRAs[i].Path, mapping)
	}
}

func CollectVideoRequest(req driver.VideoRequest) []string {
	paths := appendPath(nil, req.InitImagePath)
	return appendPath(paths, req.InitVideoPath)
}

func RewriteVideoRequest(req *driver.VideoRequest, mapping map[string]string) {
	req.InitImagePath = remap(req.InitImagePath, mapping)
	req.InitVideoPath = remap(req.InitVideoPath, mapping)
}

func CollectAudioRequest(req driver.AudioRequest) []string {
	return appendPath(nil, req.ReferencePath)
}

func RewriteAudioRequest(req *driver.AudioRequest, mapping map[string]string) {
	req.ReferencePath = remap(req.ReferencePath, mapping)
}

func CollectMeshRequest(req driver.MeshRequest) []string {
	paths := appendPath(nil, req.InitMeshPath)
	paths = appendPath(paths, req.HighMeshPath)
	paths = appendPath(paths, req.InitImagePath)
	paths = appendPaths(paths, req.InitImagePaths)
	paths = appendPath(paths, req.InitVideoPath)
	paths = appendPath(paths, req.InitDepthPath)
	paths = appendPath(paths, req.InitPointCloudPath)
	paths = appendPath(paths, req.InitSplatPath)
	paths = appendPath(paths, req.MaskPath)
	paths = appendPath(paths, req.StyleImagePath)
	paths = appendPath(paths, req.TextureImagePath)
	paths = appendPath(paths, req.AnimationPath)
	return paths
}

func RewriteMeshRequest(req *driver.MeshRequest, mapping map[string]string) {
	req.InitMeshPath = remap(req.InitMeshPath, mapping)
	req.HighMeshPath = remap(req.HighMeshPath, mapping)
	req.InitImagePath = remap(req.InitImagePath, mapping)
	for i := range req.InitImagePaths {
		req.InitImagePaths[i] = remap(req.InitImagePaths[i], mapping)
	}
	req.InitVideoPath = remap(req.InitVideoPath, mapping)
	req.InitDepthPath = remap(req.InitDepthPath, mapping)
	req.InitPointCloudPath = remap(req.InitPointCloudPath, mapping)
	req.InitSplatPath = remap(req.InitSplatPath, mapping)
	req.MaskPath = remap(req.MaskPath, mapping)
	req.StyleImagePath = remap(req.StyleImagePath, mapping)
	req.TextureImagePath = remap(req.TextureImagePath, mapping)
	req.AnimationPath = remap(req.AnimationPath, mapping)
}

func CollectDataRequest(req driver.DataRequest) []string {
	paths := appendPath(nil, req.TextFile)
	paths = appendPath(paths, req.ImagePath)
	paths = appendPath(paths, req.CSVPath)
	paths = appendPath(paths, req.GraphPath)
	return paths
}

func RewriteDataRequest(req *driver.DataRequest, mapping map[string]string) {
	req.TextFile = remap(req.TextFile, mapping)
	req.ImagePath = remap(req.ImagePath, mapping)
	req.CSVPath = remap(req.CSVPath, mapping)
	req.GraphPath = remap(req.GraphPath, mapping)
}

func CollectRAGRequest(req driver.RAGRequest) []string {
	paths := appendPaths(nil, req.SourcePaths)
	paths = appendPath(paths, req.QueryFile)
	paths = appendPath(paths, req.ImportPath)
	return paths
}

func RewriteRAGRequest(req *driver.RAGRequest, mapping map[string]string) {
	for i := range req.SourcePaths {
		req.SourcePaths[i] = remap(req.SourcePaths[i], mapping)
	}
	req.QueryFile = remap(req.QueryFile, mapping)
	req.ImportPath = remap(req.ImportPath, mapping)
}

func CollectDatasetRequest(req driver.DatasetRequest) []string {
	paths := appendPath(nil, req.Path)
	paths = appendPath(paths, req.SourcePath)
	return paths
}

func RewriteDatasetRequest(req *driver.DatasetRequest, mapping map[string]string) {
	req.Path = remap(req.Path, mapping)
	req.SourcePath = remap(req.SourcePath, mapping)
}

func RewriteVoiceTrainRequest(req *driver.VoiceTrainRequest, mapping map[string]string) {
	req.DatasetID = remap(req.DatasetID, mapping)
	req.SamplePath = remap(req.SamplePath, mapping)
}

func CollectTextTrainRequest(req driver.TextTrainRequest) []string {
	return appendPath(nil, req.DatasetID)
}

func RewriteTextTrainRequest(req *driver.TextTrainRequest, mapping map[string]string) {
	req.DatasetID = remap(req.DatasetID, mapping)
}

func CollectImageTrainRequest(req driver.ImageTrainRequest) []string {
	paths := appendPath(nil, req.DatasetID)
	return appendPath(paths, req.RegDatasetID)
}

func RewriteImageTrainRequest(req *driver.ImageTrainRequest, mapping map[string]string) {
	req.DatasetID = remap(req.DatasetID, mapping)
	req.RegDatasetID = remap(req.RegDatasetID, mapping)
}

func CollectVideoTrainRequest(req driver.VideoTrainRequest) []string {
	return appendPath(nil, req.DatasetID)
}

func RewriteVideoTrainRequest(req *driver.VideoTrainRequest, mapping map[string]string) {
	req.DatasetID = remap(req.DatasetID, mapping)
}

func CollectAudioTrainRequest(req driver.AudioTrainRequest) []string {
	return appendPath(nil, req.DatasetID)
}

func RewriteAudioTrainRequest(req *driver.AudioTrainRequest, mapping map[string]string) {
	req.DatasetID = remap(req.DatasetID, mapping)
}

func CollectMeshTrainRequest(req driver.MeshTrainRequest) []string {
	return appendPath(nil, req.DatasetID)
}

func RewriteMeshTrainRequest(req *driver.MeshTrainRequest, mapping map[string]string) {
	req.DatasetID = remap(req.DatasetID, mapping)
}

func CollectVoiceTrainRequest(req driver.VoiceTrainRequest) []string {
	paths := appendPath(nil, req.DatasetID)
	return appendPath(paths, req.SamplePath)
}

func CollectImageResponse(resp driver.ImageResponse) []string {
	paths := appendPath(nil, resp.Path)
	return appendPaths(paths, resp.Paths)
}

func RewriteImageResponse(resp *driver.ImageResponse, mapping map[string]string) {
	resp.Path = remap(resp.Path, mapping)
	for i := range resp.Paths {
		resp.Paths[i] = remap(resp.Paths[i], mapping)
	}
}

func CollectVideoResponse(resp driver.VideoResponse) []string {
	return appendPath(nil, resp.Path)
}

func RewriteVideoResponse(resp *driver.VideoResponse, mapping map[string]string) {
	resp.Path = remap(resp.Path, mapping)
}

func CollectAudioResponse(resp driver.AudioResponse) []string {
	return appendPath(nil, resp.Path)
}

func RewriteAudioResponse(resp *driver.AudioResponse, mapping map[string]string) {
	resp.Path = remap(resp.Path, mapping)
}

func CollectMeshResponse(resp driver.MeshResponse) []string {
	return appendPath(nil, resp.Path)
}

func RewriteMeshResponse(resp *driver.MeshResponse, mapping map[string]string) {
	resp.Path = remap(resp.Path, mapping)
}

func CollectTrainResponse(resp driver.TrainResponse) []string {
	return appendPath(nil, resp.OutputPath)
}

func RewriteTrainResponse(resp *driver.TrainResponse, mapping map[string]string) {
	resp.OutputPath = remap(resp.OutputPath, mapping)
}

func remap(path string, mapping map[string]string) string {
	if path == "" {
		return path
	}
	if mapped, ok := mapping[path]; ok && mapped != "" {
		return mapped
	}
	clean := filepath.Clean(path)
	if mapped, ok := mapping[clean]; ok && mapped != "" {
		return mapped
	}
	return path
}

// RAGExportDownload returns server paths to fetch after an export task.
func RAGExportDownload(req driver.RAGRequest) []FetchTarget {
	if req.Task != driver.RAGTaskExport || strings.TrimSpace(req.ExportPath) == "" {
		return nil
	}
	return []FetchTarget{{
		Server: req.ExportPath,
		Local:  req.ExportPath,
	}}
}

// LocalOutputPath saves downloaded media under <root>/.wuji/remote-out/.
func LocalOutputPath(root, serverPath, sessionID string) string {
	base := filepath.Join(root, ".wuji", "remote-out", sanitizeSessionID(sessionID))
	name := filepath.Base(serverPath)
	if name == "" || name == "." || name == string(filepath.Separator) {
		name = "output"
	}
	return filepath.Join(base, name)
}
