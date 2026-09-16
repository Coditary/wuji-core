package driver

import (
	"fmt"
	"strings"
)

// MeshTask identifies the mesh generation or transformation mode.
type MeshTask string

const (
	MeshTaskGenerate        MeshTask = "generate"
	MeshTaskImageToMesh     MeshTask = "i2m"
	MeshTaskMultiView       MeshTask = "multiview"
	MeshTaskVideoToMesh     MeshTask = "v2m"
	MeshTaskDepthToMesh     MeshTask = "depth2m"
	MeshTaskPointCloudToMesh MeshTask = "pcd2m"
	MeshTaskSplatToMesh     MeshTask = "splat2m"
	MeshTaskInpaint         MeshTask = "inpaint"
	MeshTaskImageToTexture  MeshTask = "img2tex"
	MeshTaskStyle           MeshTask = "style"
	MeshTaskTexture         MeshTask = "texture"
	MeshTaskPBR             MeshTask = "pbr"
	MeshTaskUV              MeshTask = "uv"
	MeshTaskDecimate        MeshTask = "decimate"
	MeshTaskRetopo          MeshTask = "retopo"
	MeshTaskRemesh          MeshTask = "remesh"
	MeshTaskRepair          MeshTask = "repair"
	MeshTaskSmooth          MeshTask = "smooth"
	MeshTaskRefine          MeshTask = "refine"
	MeshTaskSegment         MeshTask = "segment"
	MeshTaskRig             MeshTask = "rig"
	MeshTaskAnimate         MeshTask = "animate"
	MeshTaskRetarget        MeshTask = "retarget"
	MeshTaskBake            MeshTask = "bake"
	MeshTaskScene           MeshTask = "scene"
	MeshTaskVariation       MeshTask = "variation"
	MeshTaskEdit            MeshTask = "edit"
	MeshTaskUpscale         MeshTask = "upscale"
)

// MeshTaskInfo describes a task and its required inputs.
type MeshTaskInfo struct {
	Task                  MeshTask
	Description           string
	RequiresPrompt        bool
	RequiresInitMesh      bool
	RequiresInitImage     bool
	RequiresInitImages    bool
	RequiresInitVideo     bool
	RequiresInitDepth     bool
	RequiresInitPointCloud bool
	RequiresInitSplat     bool
	RequiresMask          bool
	RequiresStyleImage    bool
	RequiresTextureImage  bool
	RequiresHighMesh      bool
	RequiresAnimation     bool
	RequiresTargetTris    bool
	RequiresScale         bool
}

// AllInferrableMeshTasks returns mesh tasks that can be selected automatically from inputs.
func AllInferrableMeshTasks() []MeshTask {
	return []MeshTask{
		MeshTaskGenerate,
		MeshTaskImageToMesh,
		MeshTaskMultiView,
		MeshTaskVideoToMesh,
		MeshTaskDepthToMesh,
		MeshTaskPointCloudToMesh,
		MeshTaskSplatToMesh,
		MeshTaskInpaint,
		MeshTaskImageToTexture,
		MeshTaskStyle,
		MeshTaskTexture,
		MeshTaskPBR,
		MeshTaskUV,
		MeshTaskDecimate,
		MeshTaskRetopo,
		MeshTaskRemesh,
		MeshTaskRepair,
		MeshTaskSmooth,
		MeshTaskRefine,
		MeshTaskSegment,
		MeshTaskRig,
		MeshTaskAnimate,
		MeshTaskRetarget,
		MeshTaskBake,
		MeshTaskScene,
		MeshTaskVariation,
		MeshTaskEdit,
		MeshTaskUpscale,
	}
}

// AllMeshTasks returns every known mesh task (same as inferrable for now).
func AllMeshTasks() []MeshTask {
	return AllInferrableMeshTasks()
}

func MeshTaskCatalog() []MeshTaskInfo {
	return []MeshTaskInfo{
		{Task: MeshTaskGenerate, Description: "Text-to-mesh generation from a prompt", RequiresPrompt: true},
		{Task: MeshTaskImageToMesh, Description: "Reconstruct a mesh from a single image", RequiresInitImage: true},
		{Task: MeshTaskMultiView, Description: "Reconstruct a mesh from multiple views", RequiresInitImages: true},
		{Task: MeshTaskVideoToMesh, Description: "Reconstruct a mesh from video", RequiresInitVideo: true},
		{Task: MeshTaskDepthToMesh, Description: "Build a mesh from a depth or normal map", RequiresInitDepth: true},
		{Task: MeshTaskPointCloudToMesh, Description: "Surface a point cloud into a mesh", RequiresInitPointCloud: true},
		{Task: MeshTaskSplatToMesh, Description: "Export a Gaussian splat as a mesh", RequiresInitSplat: true},
		{Task: MeshTaskInpaint, Description: "Replace a masked region on a mesh", RequiresInitMesh: true, RequiresMask: true, RequiresPrompt: true},
		{Task: MeshTaskImageToTexture, Description: "Project a reference image onto a mesh", RequiresInitMesh: true, RequiresTextureImage: true},
		{Task: MeshTaskStyle, Description: "Apply a style reference to a mesh", RequiresInitMesh: true, RequiresStyleImage: true},
		{Task: MeshTaskTexture, Description: "Generate textures for a mesh from a prompt", RequiresInitMesh: true, RequiresPrompt: true},
		{Task: MeshTaskPBR, Description: "Generate PBR material maps for a mesh", RequiresInitMesh: true},
		{Task: MeshTaskUV, Description: "Auto-unwrap UVs for a mesh", RequiresInitMesh: true},
		{Task: MeshTaskDecimate, Description: "Reduce polygon count", RequiresInitMesh: true, RequiresTargetTris: true},
		{Task: MeshTaskRetopo, Description: "Retopologize for game-ready topology", RequiresInitMesh: true},
		{Task: MeshTaskRemesh, Description: "Rebuild mesh with uniform topology", RequiresInitMesh: true},
		{Task: MeshTaskRepair, Description: "Fix holes and non-manifold geometry", RequiresInitMesh: true},
		{Task: MeshTaskSmooth, Description: "Smooth noisy mesh geometry", RequiresInitMesh: true},
		{Task: MeshTaskRefine, Description: "Add geometric detail to a mesh", RequiresInitMesh: true},
		{Task: MeshTaskSegment, Description: "Semantic part segmentation", RequiresInitMesh: true},
		{Task: MeshTaskRig, Description: "Auto-rig a mesh", RequiresInitMesh: true},
		{Task: MeshTaskAnimate, Description: "Generate animation on a rigged mesh", RequiresInitMesh: true},
		{Task: MeshTaskRetarget, Description: "Retarget animation onto a mesh", RequiresInitMesh: true, RequiresAnimation: true},
		{Task: MeshTaskBake, Description: "Bake maps from a high-poly to low-poly mesh", RequiresInitMesh: true, RequiresHighMesh: true},
		{Task: MeshTaskScene, Description: "Generate a multi-object scene from a prompt", RequiresPrompt: true},
		{Task: MeshTaskVariation, Description: "Create a variant of an existing mesh", RequiresInitMesh: true},
		{Task: MeshTaskEdit, Description: "Edit mesh geometry from a prompt", RequiresInitMesh: true, RequiresPrompt: true},
		{Task: MeshTaskUpscale, Description: "Increase mesh geometric resolution", RequiresInitMesh: true, RequiresScale: true},
	}
}

func (t MeshTask) String() string { return string(t) }

func (t MeshTask) IsValid() bool {
	for _, known := range AllMeshTasks() {
		if t == known {
			return true
		}
	}
	return false
}

func (t MeshTask) Info() *MeshTaskInfo {
	for _, info := range MeshTaskCatalog() {
		if info.Task == t {
			copy := info
			return &copy
		}
	}
	return nil
}

// MeshTaskInputs holds CLI fields used to infer a mesh task implicitly.
type MeshTaskInputs struct {
	Prompt           string
	MeshPath         string
	HighMeshPath     string
	ImagePath        string
	Images           []string
	VideoPath        string
	DepthPath        string
	PointCloudPath   string
	SplatPath        string
	MaskPath         string
	StyleImagePath   string
	TextureImagePath string
	AnimationPath    string
	TargetTris       int
	TargetTrisSet    bool
	RetopoRequested  bool
	RemeshRequested  bool
	RepairRequested  bool
	SmoothRequested  bool
	RefineRequested  bool
	SegmentRequested bool
	RigRequested     bool
	AnimateRequested bool
	RetargetRequested bool
	UVRequested       bool
	PBRRequested      bool
	SceneRequested    bool
	VariationRequested bool
	EditRequested     bool
	UpscaleRequested  bool
	Scale             float32
	ScaleSet          bool
}

// InferMeshTask selects a mesh task from provided inputs when --task is omitted.
func InferMeshTask(in MeshTaskInputs) (MeshTask, error) {
	if err := validateMeshPrimaryInputs(in); err != nil {
		return "", err
	}
	if err := validateMeshOpFlags(in); err != nil {
		return "", err
	}
	if err := validateMeshSceneInputs(in); err != nil {
		return "", err
	}

	hasMesh := strings.TrimSpace(in.MeshPath) != ""
	if hasMesh {
		if strings.TrimSpace(in.MaskPath) != "" {
			return MeshTaskInpaint, nil
		}
		if in.RetargetRequested || strings.TrimSpace(in.AnimationPath) != "" {
			return MeshTaskRetarget, nil
		}
		if strings.TrimSpace(in.HighMeshPath) != "" {
			return MeshTaskBake, nil
		}
		if strings.TrimSpace(in.StyleImagePath) != "" {
			return MeshTaskStyle, nil
		}
		if strings.TrimSpace(in.TextureImagePath) != "" {
			return MeshTaskImageToTexture, nil
		}
		if in.TargetTrisSet {
			return MeshTaskDecimate, nil
		}
		if in.RetopoRequested {
			return MeshTaskRetopo, nil
		}
		if in.RemeshRequested {
			return MeshTaskRemesh, nil
		}
		if in.RepairRequested {
			return MeshTaskRepair, nil
		}
		if in.SmoothRequested {
			return MeshTaskSmooth, nil
		}
		if in.RefineRequested {
			return MeshTaskRefine, nil
		}
		if in.SegmentRequested {
			return MeshTaskSegment, nil
		}
		if in.RigRequested {
			return MeshTaskRig, nil
		}
		if in.AnimateRequested {
			return MeshTaskAnimate, nil
		}
		if in.UVRequested {
			return MeshTaskUV, nil
		}
		if in.PBRRequested {
			return MeshTaskPBR, nil
		}
		if in.EditRequested && strings.TrimSpace(in.Prompt) != "" {
			return MeshTaskEdit, nil
		}
		if in.VariationRequested {
			return MeshTaskVariation, nil
		}
		if in.UpscaleRequested || in.ScaleSet {
			return MeshTaskUpscale, nil
		}
		if strings.TrimSpace(in.Prompt) != "" {
			return MeshTaskTexture, nil
		}
		return "", fmt.Errorf("with --mesh, set an operation flag (--target-tris, --repair, --variation, --edit, …) or a prompt for texturing")
	}

	if len(in.Images) >= 2 {
		return MeshTaskMultiView, nil
	}
	if strings.TrimSpace(in.VideoPath) != "" {
		return MeshTaskVideoToMesh, nil
	}
	if strings.TrimSpace(in.DepthPath) != "" {
		return MeshTaskDepthToMesh, nil
	}
	if strings.TrimSpace(in.PointCloudPath) != "" {
		return MeshTaskPointCloudToMesh, nil
	}
	if strings.TrimSpace(in.SplatPath) != "" {
		return MeshTaskSplatToMesh, nil
	}
	if strings.TrimSpace(in.ImagePath) != "" || len(in.Images) == 1 {
		return MeshTaskImageToMesh, nil
	}
	if in.SceneRequested && strings.TrimSpace(in.Prompt) != "" {
		return MeshTaskScene, nil
	}
	if strings.TrimSpace(in.Prompt) != "" {
		return MeshTaskGenerate, nil
	}
	return "", fmt.Errorf("prompt or an input file is required (--mesh, --image, --video, --depth, --pointcloud, --splat)")
}

func validateMeshOpFlags(in MeshTaskInputs) error {
	type flag struct {
		name string
		set  bool
	}
	flags := []flag{
		{"--retopo", in.RetopoRequested},
		{"--remesh", in.RemeshRequested},
		{"--repair", in.RepairRequested},
		{"--smooth", in.SmoothRequested},
		{"--refine", in.RefineRequested},
		{"--segment", in.SegmentRequested},
		{"--rig", in.RigRequested},
		{"--animate", in.AnimateRequested},
		{"--retarget", in.RetargetRequested},
		{"--uv", in.UVRequested},
		{"--pbr", in.PBRRequested},
		{"--variation", in.VariationRequested},
		{"--edit", in.EditRequested},
		{"--upscale", in.UpscaleRequested},
		{"--target-tris", in.TargetTrisSet},
		{"--scene", in.SceneRequested},
	}
	active := make([]string, 0, len(flags))
	for _, f := range flags {
		if f.set {
			active = append(active, f.name)
		}
	}
	if len(active) > 1 {
		return fmt.Errorf("use only one mesh operation flag (%s)", strings.Join(active, ", "))
	}
	if in.EditRequested && strings.TrimSpace(in.Prompt) == "" {
		return fmt.Errorf("prompt is required with --edit")
	}
	if in.SceneRequested && strings.TrimSpace(in.Prompt) == "" {
		return fmt.Errorf("prompt is required with --scene")
	}
	return nil
}

func validateMeshSceneInputs(in MeshTaskInputs) error {
	if !in.SceneRequested {
		return nil
	}
	hasOtherInput := strings.TrimSpace(in.MeshPath) != "" ||
		strings.TrimSpace(in.ImagePath) != "" || len(in.Images) > 0 ||
		strings.TrimSpace(in.VideoPath) != "" ||
		strings.TrimSpace(in.DepthPath) != "" ||
		strings.TrimSpace(in.PointCloudPath) != "" ||
		strings.TrimSpace(in.SplatPath) != ""
	if hasOtherInput {
		return fmt.Errorf("--scene is only valid with a prompt (no --mesh, --image, --video, …)")
	}
	return nil
}

func validateMeshPrimaryInputs(in MeshTaskInputs) error {
	type slot struct {
		name string
		set  bool
	}
	slots := []slot{
		{"--mesh", strings.TrimSpace(in.MeshPath) != ""},
		{"--image/--images", strings.TrimSpace(in.ImagePath) != "" || len(in.Images) > 0},
		{"--video", strings.TrimSpace(in.VideoPath) != ""},
		{"--depth", strings.TrimSpace(in.DepthPath) != ""},
		{"--pointcloud", strings.TrimSpace(in.PointCloudPath) != ""},
		{"--splat", strings.TrimSpace(in.SplatPath) != ""},
	}
	active := make([]string, 0, len(slots))
	for _, s := range slots {
		if s.set {
			active = append(active, s.name)
		}
	}
	if len(active) > 1 {
		return fmt.Errorf("use only one primary input (%s)", strings.Join(active, ", "))
	}
	return nil
}

// ParseMeshTask normalizes and validates a task name from CLI or config.
func ParseMeshTask(raw string) (MeshTask, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	switch normalized {
	case "", "txt2mesh", "text-to-mesh", "text2mesh":
		return MeshTaskGenerate, nil
	case "image-to-mesh", "img2mesh":
		return MeshTaskImageToMesh, nil
	case "multi-view", "mv":
		return MeshTaskMultiView, nil
	case "video-to-mesh", "v2m":
		return MeshTaskVideoToMesh, nil
	case "depth-to-mesh", "depth2mesh":
		return MeshTaskDepthToMesh, nil
	case "pointcloud-to-mesh", "pc2mesh":
		return MeshTaskPointCloudToMesh, nil
	case "splat-to-mesh":
		return MeshTaskSplatToMesh, nil
	case "image-to-texture", "tex2mesh":
		return MeshTaskImageToTexture, nil
	case "style-transfer":
		return MeshTaskStyle, nil
	}

	task := MeshTask(normalized)
	if !task.IsValid() {
		return "", fmt.Errorf("unknown mesh task %q (valid: %s)", raw, joinMeshTasks())
	}
	return task, nil
}

func (r MeshRequest) TaskOrDefault() MeshTask {
	if r.Task == "" {
		return MeshTaskGenerate
	}
	return r.Task
}

// Validate checks task-specific required fields.
func (r MeshRequest) Validate() error {
	task := r.TaskOrDefault()
	if !task.IsValid() {
		return fmt.Errorf("unknown mesh task %q (valid: %s)", task, joinMeshTasks())
	}
	info := task.Info()
	if info == nil {
		return fmt.Errorf("unknown mesh task %q", task)
	}

	if info.RequiresPrompt && strings.TrimSpace(r.Prompt) == "" {
		return fmt.Errorf("prompt is required for mesh task %q", task)
	}
	if info.RequiresInitMesh && r.InitMeshPath == "" {
		return fmt.Errorf("--mesh is required for mesh task %q", task)
	}
	if info.RequiresInitImage && r.InitImagePath == "" && len(r.InitImagePaths) == 0 {
		return fmt.Errorf("--image is required for mesh task %q", task)
	}
	if info.RequiresInitImages && len(r.InitImagePaths) < 2 {
		return fmt.Errorf("at least two --images are required for mesh task %q", task)
	}
	if info.RequiresInitVideo && r.InitVideoPath == "" {
		return fmt.Errorf("--video is required for mesh task %q", task)
	}
	if info.RequiresInitDepth && r.InitDepthPath == "" {
		return fmt.Errorf("--depth is required for mesh task %q", task)
	}
	if info.RequiresInitPointCloud && r.InitPointCloudPath == "" {
		return fmt.Errorf("--pointcloud is required for mesh task %q", task)
	}
	if info.RequiresInitSplat && r.InitSplatPath == "" {
		return fmt.Errorf("--splat is required for mesh task %q", task)
	}
	if info.RequiresMask && r.MaskPath == "" {
		return fmt.Errorf("--mask is required for mesh task %q", task)
	}
	if info.RequiresStyleImage && r.StyleImagePath == "" {
		return fmt.Errorf("--style-image is required for mesh task %q", task)
	}
	if info.RequiresTextureImage && r.TextureImagePath == "" {
		return fmt.Errorf("--texture-image is required for mesh task %q", task)
	}
	if info.RequiresHighMesh && r.HighMeshPath == "" {
		return fmt.Errorf("--high-mesh is required for mesh task %q", task)
	}
	if info.RequiresAnimation && r.AnimationPath == "" {
		return fmt.Errorf("--animation is required for mesh task %q", task)
	}
	if info.RequiresTargetTris && r.TargetTris <= 0 {
		return fmt.Errorf("--target-tris is required for mesh task %q", task)
	}
	if info.RequiresScale && r.Scale <= 0 {
		return fmt.Errorf("--scale is required for mesh task %q (e.g. 2, 2.5, 200%%)", task)
	}
	if r.Mode != "" {
		if _, err := ParseAssetMode(string(r.Mode)); err != nil {
			return err
		}
	}
	rep := r.RepresentationOrDefault()
	if r.Representation != "" && !r.Representation.IsValid() {
		return fmt.Errorf("unknown mesh representation %q", r.Representation)
	}
	if rep != MeshRepresentationMesh && !meshTaskSupportsRepresentation(task) {
		return fmt.Errorf("--representation %q is only valid for generation tasks (generate, i2m, multiview, v2m, depth2m, scene)", rep)
	}
	return nil
}

func joinMeshTasks() string {
	parts := make([]string, len(AllMeshTasks()))
	for i, t := range AllMeshTasks() {
		parts[i] = string(t)
	}
	return strings.Join(parts, ", ")
}
