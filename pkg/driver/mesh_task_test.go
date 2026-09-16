package driver

import (
	"context"
	"testing"

	"github.com/coditary/wuji-core/pkg/capability"
)

func TestInferMeshTaskGenerate(t *testing.T) {
	task, err := InferMeshTask(MeshTaskInputs{Prompt: "crate"})
	if err != nil {
		t.Fatalf("InferMeshTask: %v", err)
	}
	if task != MeshTaskGenerate {
		t.Fatalf("task = %q, want generate", task)
	}
}

func TestInferMeshTaskImageToMesh(t *testing.T) {
	task, err := InferMeshTask(MeshTaskInputs{ImagePath: "ref.png"})
	if err != nil {
		t.Fatalf("InferMeshTask: %v", err)
	}
	if task != MeshTaskImageToMesh {
		t.Fatalf("task = %q, want i2m", task)
	}
}

func TestInferMeshTaskMultiView(t *testing.T) {
	task, err := InferMeshTask(MeshTaskInputs{Images: []string{"a.png", "b.png"}})
	if err != nil {
		t.Fatalf("InferMeshTask: %v", err)
	}
	if task != MeshTaskMultiView {
		t.Fatalf("task = %q, want multiview", task)
	}
}

func TestInferMeshTaskTexture(t *testing.T) {
	task, err := InferMeshTask(MeshTaskInputs{MeshPath: "hero.glb", Prompt: "rusty"})
	if err != nil {
		t.Fatalf("InferMeshTask: %v", err)
	}
	if task != MeshTaskTexture {
		t.Fatalf("task = %q, want texture", task)
	}
}

func TestInferMeshTaskDecimate(t *testing.T) {
	task, err := InferMeshTask(MeshTaskInputs{MeshPath: "hero.glb", TargetTris: 5000, TargetTrisSet: true})
	if err != nil {
		t.Fatalf("InferMeshTask: %v", err)
	}
	if task != MeshTaskDecimate {
		t.Fatalf("task = %q, want decimate", task)
	}
}

func TestInferMeshTaskAmbiguousMeshOnly(t *testing.T) {
	_, err := InferMeshTask(MeshTaskInputs{MeshPath: "hero.glb"})
	if err == nil {
		t.Fatal("expected error for mesh-only input")
	}
}

func TestInferMeshTaskPrimaryInputCollision(t *testing.T) {
	_, err := InferMeshTask(MeshTaskInputs{MeshPath: "a.glb", ImagePath: "b.png"})
	if err == nil {
		t.Fatal("expected error for multiple primary inputs")
	}
}

func TestInferMeshTaskInpaint(t *testing.T) {
	task, err := InferMeshTask(MeshTaskInputs{
		MeshPath: "hero.glb",
		MaskPath: "mask.png",
		Prompt:   "sword",
	})
	if err != nil {
		t.Fatalf("InferMeshTask: %v", err)
	}
	if task != MeshTaskInpaint {
		t.Fatalf("task = %q, want inpaint", task)
	}
}

func TestMeshRequestValidateGenerateRequiresPrompt(t *testing.T) {
	req := MeshRequest{Task: MeshTaskGenerate}
	if err := req.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestParseMeshTaskAliases(t *testing.T) {
	task, err := ParseMeshTask("text-to-mesh")
	if err != nil {
		t.Fatalf("ParseMeshTask: %v", err)
	}
	if task != MeshTaskGenerate {
		t.Fatalf("task = %q, want generate", task)
	}
}

func TestInferMeshTaskScene(t *testing.T) {
	task, err := InferMeshTask(MeshTaskInputs{Prompt: "medieval room", SceneRequested: true})
	if err != nil {
		t.Fatalf("InferMeshTask: %v", err)
	}
	if task != MeshTaskScene {
		t.Fatalf("task = %q, want scene", task)
	}
}

func TestInferMeshTaskVariation(t *testing.T) {
	task, err := InferMeshTask(MeshTaskInputs{MeshPath: "hero.glb", VariationRequested: true})
	if err != nil {
		t.Fatalf("InferMeshTask: %v", err)
	}
	if task != MeshTaskVariation {
		t.Fatalf("task = %q, want variation", task)
	}
}

func TestInferMeshTaskEdit(t *testing.T) {
	task, err := InferMeshTask(MeshTaskInputs{MeshPath: "hero.glb", Prompt: "add horns", EditRequested: true})
	if err != nil {
		t.Fatalf("InferMeshTask: %v", err)
	}
	if task != MeshTaskEdit {
		t.Fatalf("task = %q, want edit", task)
	}
}

func TestInferMeshTaskUpscaleDefaultScaleValidation(t *testing.T) {
	req := MeshRequest{Task: MeshTaskUpscale, InitMeshPath: "hero.glb", Scale: 2}
	if err := req.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestInferMeshTaskConflictingOpFlags(t *testing.T) {
	_, err := InferMeshTask(MeshTaskInputs{MeshPath: "a.glb", RepairRequested: true, VariationRequested: true})
	if err == nil {
		t.Fatal("expected conflicting op flag error")
	}
}

func TestParseMeshModeAliases(t *testing.T) {
	mode, err := ParseMeshMode("character")
	if err != nil {
		t.Fatalf("ParseMeshMode: %v", err)
	}
	if mode != MeshModeHuman {
		t.Fatalf("mode = %q, want human", mode)
	}
}

func TestRunMeshTaskStub(t *testing.T) {
	d := &meshTestDriver{info: Info{ID: "mesh-test", MeshTasks: AllMeshTasks()}}
	resp, err := RunMeshTask(context.Background(), d, MeshRequest{
		Task:   MeshTaskGenerate,
		Prompt: "crate",
		Format: "glb",
	})
	if err != nil {
		t.Fatalf("RunMeshTask: %v", err)
	}
	if resp.Task != MeshTaskGenerate {
		t.Fatalf("resp.Task = %q", resp.Task)
	}
}

type meshTestDriver struct {
	info Info
}

func (d *meshTestDriver) Info() Info            { return d.info }
func (d *meshTestDriver) Capabilities() []capability.Type { return nil }
func (d *meshTestDriver) Close() error        { return nil }

func (d *meshTestDriver) GenerateMesh(_ context.Context, req MeshRequest) (*MeshResponse, error) {
	return &MeshResponse{Path: "/tmp/test.glb", Format: "glb", Task: req.TaskOrDefault()}, nil
}
