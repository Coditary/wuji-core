package driver

import wujiv1 "github.com/coditary/wuji-core/api/proto/v1"

func MeshTasksToProto(tasks []MeshTask) []string {
	return meshTasksToStrings(tasks)
}

func MeshTasksFromProto(items []string) []MeshTask {
	out := make([]MeshTask, 0, len(items))
	for _, item := range items {
		task, err := ParseMeshTask(item)
		if err != nil {
			continue
		}
		out = append(out, task)
	}
	return out
}

func MeshRequestFromProto(req *wujiv1.GenerateMeshRequest) MeshRequest {
	if req == nil {
		return MeshRequest{}
	}
	out := MeshRequest{
		Task:               MeshTask(req.GetTask()),
		Prompt:             req.GetPrompt(),
		Format:             req.GetFormat(),
		Model:              req.GetModel(),
		TargetTris:         int(req.GetTargetTris()),
		InitMeshPath:       req.GetInitMeshPath(),
		HighMeshPath:       req.GetHighMeshPath(),
		InitImagePath:      req.GetInitImagePath(),
		InitImagePaths:     append([]string(nil), req.GetInitImagePaths()...),
		InitVideoPath:      req.GetInitVideoPath(),
		InitDepthPath:      req.GetInitDepthPath(),
		InitPointCloudPath: req.GetInitPointcloudPath(),
		InitSplatPath:      req.GetInitSplatPath(),
		MaskPath:           req.GetMaskPath(),
		StyleImagePath:     req.GetStyleImagePath(),
		TextureImagePath:   req.GetTextureImagePath(),
		AnimationPath:      req.GetAnimationPath(),
		Mode:               MeshMode(req.GetMode()),
		Scale:              req.GetScale(),
		Representation:     MeshRepresentation(req.GetRepresentation()),
	}
	if req.Seed != nil {
		seed := int(req.GetSeed())
		out.Seed = &seed
	}
	return out
}

func MeshRequestToProto(req MeshRequest) *wujiv1.GenerateMeshRequest {
	out := &wujiv1.GenerateMeshRequest{
		Task:               string(req.TaskOrDefault()),
		Prompt:             req.Prompt,
		Format:             req.Format,
		Model:              req.Model,
		TargetTris:         int32(req.TargetTris),
		InitMeshPath:       req.InitMeshPath,
		HighMeshPath:       req.HighMeshPath,
		InitImagePath:      req.InitImagePath,
		InitImagePaths:     append([]string(nil), req.InitImagePaths...),
		InitVideoPath:      req.InitVideoPath,
		InitDepthPath:      req.InitDepthPath,
		InitPointcloudPath: req.InitPointCloudPath,
		InitSplatPath:      req.InitSplatPath,
		MaskPath:           req.MaskPath,
		StyleImagePath:     req.StyleImagePath,
		TextureImagePath:   req.TextureImagePath,
		AnimationPath:      req.AnimationPath,
		Mode:               string(req.ModeOrDefault()),
		Scale:              req.Scale,
		Representation:     string(req.RepresentationOrDefault()),
	}
	if req.Seed != nil {
		seed := int32(*req.Seed)
		out.Seed = &seed
	}
	return out
}

func MeshResponseFromProto(resp *wujiv1.GenerateMeshResponse) *MeshResponse {
	if resp == nil {
		return nil
	}
	return &MeshResponse{
		Path:           resp.GetPath(),
		Format:         resp.GetFormat(),
		Task:           MeshTask(resp.GetTask()),
		Mode:           AssetMode(resp.GetMode()),
		Representation: MeshRepresentation(resp.GetRepresentation()),
	}
}

func meshTasksToStrings(tasks []MeshTask) []string {
	out := make([]string, len(tasks))
	for i, t := range tasks {
		out[i] = string(t)
	}
	return out
}
