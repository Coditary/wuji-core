package scheduler_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/scheduler"
)

func TestSchedulerEvictsAndLazyRestores(t *testing.T) {
	root := t.TempDir()
	cfg := config.EffectiveResources{
		TotalVRAMMB:           10000,
		TotalRAMMB:            50000,
		EvictAndRestore:       true,
		LazyRestore:           true,
		UnloadTransientModels: true,
	}

	var unloadLog []string
	var warmLog []string
	s := scheduler.New(root, cfg,
		func(ctx context.Context, driverID string) error {
			unloadLog = append(unloadLog, driverID)
			return nil
		},
		func(ctx context.Context, snap scheduler.SessionSnapshot) error {
			warmLog = append(warmLog, snap.DriverID+":"+snap.Model)
			return nil
		},
		nil,
	)

	stateDir := filepath.Join(root, ".wuji", "runtime")
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	state := scheduler.PersistedState{
		Loaded: []scheduler.LoadedSlot{{
			DriverID: "llama", Model: "chat.gguf",
			Capability: capability.TextGeneration, VRAMMB: 6000, RAMMB: 12000,
		}},
	}
	data, _ := json.MarshalIndent(state, "", "  ")
	if err := os.WriteFile(filepath.Join(stateDir, "scheduler.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	err := s.Run(context.Background(), scheduler.JobSpec{
		DriverID: "dummy", Model: "sd-xl",
		Capability: capability.ImageGeneration, VRAMMB: 8000, RAMMB: 8000, Transient: true,
	}, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(unloadLog) == 0 {
		t.Fatal("expected unload during image job")
	}
	if len(warmLog) != 0 {
		t.Fatalf("lazy restore should not warm immediately, got %v", warmLog)
	}

	persisted, err := s.Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(persisted.RestoreStack) != 1 || persisted.RestoreStack[0].DriverID != "llama" {
		t.Fatalf("expected restore queue entry, got %+v", persisted.RestoreStack)
	}
}

func TestSchedulerEagerRestore(t *testing.T) {
	root := t.TempDir()
	cfg := config.EffectiveResources{
		TotalVRAMMB:           10000,
		EvictAndRestore:       true,
		LazyRestore:           false,
		UnloadTransientModels: true,
	}

	var warmLog []string
	s := scheduler.New(root, cfg, nil,
		func(ctx context.Context, snap scheduler.SessionSnapshot) error {
			warmLog = append(warmLog, snap.DriverID+":"+snap.Model)
			return nil
		},
		nil,
	)

	stateDir := filepath.Join(root, ".wuji", "runtime")
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	state := scheduler.PersistedState{
		Loaded: []scheduler.LoadedSlot{{
			DriverID: "llama", Model: "chat.gguf",
			Capability: capability.TextGeneration, VRAMMB: 6000,
		}},
	}
	data, _ := json.MarshalIndent(state, "", "  ")
	if err := os.WriteFile(filepath.Join(stateDir, "scheduler.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	err := s.Run(context.Background(), scheduler.JobSpec{
		DriverID: "dummy", Model: "sd-xl",
		Capability: capability.ImageGeneration, VRAMMB: 8000, Transient: true,
	}, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(warmLog) != 1 || warmLog[0] != "llama:chat.gguf" {
		t.Fatalf("expected eager restore, got %v", warmLog)
	}
}

func TestSchedulerRAMBudgetEviction(t *testing.T) {
	root := t.TempDir()
	cfg := config.EffectiveResources{
		TotalRAMMB:      20000,
		EvictAndRestore: true,
		LazyRestore:     true,
	}

	var unloaded []string
	s := scheduler.New(root, cfg,
		func(ctx context.Context, driverID string) error {
			unloaded = append(unloaded, driverID)
			return nil
		}, nil, nil)

	stateDir := filepath.Join(root, ".wuji", "runtime")
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	state := scheduler.PersistedState{
		Loaded: []scheduler.LoadedSlot{{
			DriverID: "vllm", Model: "big-model",
			Capability: capability.TextGeneration, RAMMB: 16000,
		}},
	}
	data, _ := json.MarshalIndent(state, "", "  ")
	if err := os.WriteFile(filepath.Join(stateDir, "scheduler.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	err := s.Run(context.Background(), scheduler.JobSpec{
		DriverID: "dummy", Model: "img",
		Capability: capability.ImageGeneration, RAMMB: 8000, Transient: true,
	}, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(unloaded) == 0 {
		t.Fatal("expected RAM-budget eviction")
	}
}

func TestSchedulerDisabledPassthrough(t *testing.T) {
	s := scheduler.New(t.TempDir(), config.EffectiveResources{}, nil, nil, nil)
	called := false
	if err := s.Run(context.Background(), scheduler.JobSpec{}, func(ctx context.Context) error {
		called = true
		return nil
	}); err != nil || !called {
		t.Fatalf("expected passthrough when disabled")
	}
}

func TestSchedulerReconcileDropsStaleSlots(t *testing.T) {
	root := t.TempDir()
	cfg := config.EffectiveResources{TotalRAMMB: 64000}
	probe := func(driverID string) bool { return false }
	s := scheduler.New(root, cfg, nil, nil, probe)

	stateDir := filepath.Join(root, ".wuji", "runtime")
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	state := scheduler.PersistedState{
		Loaded: []scheduler.LoadedSlot{{
			DriverID: "vllm", Model: "stale",
			Capability: capability.TextGeneration, RAMMB: 8000,
		}},
	}
	data, _ := json.MarshalIndent(state, "", "  ")
	if err := os.WriteFile(filepath.Join(stateDir, "scheduler.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	err := s.Run(context.Background(), scheduler.JobSpec{
		DriverID: "dummy", Model: "img",
		Capability: capability.ImageGeneration, RAMMB: 4000, Transient: true,
	}, func(ctx context.Context) error { return nil })
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	persisted, err := s.Status()
	if err != nil {
		t.Fatal(err)
	}
	for _, slot := range persisted.Loaded {
		if slot.DriverID == "vllm" {
			t.Fatal("stale vllm slot should have been reconciled away")
		}
	}
}

func TestSchedulerManualLoadUnload(t *testing.T) {
	root := t.TempDir()
	cfg := config.EffectiveResources{TotalRAMMB: 32000}

	var warmed []string
	var unloaded []string
	s := scheduler.New(root, cfg,
		func(ctx context.Context, driverID string) error {
			unloaded = append(unloaded, driverID)
			return nil
		},
		func(ctx context.Context, snap scheduler.SessionSnapshot) error {
			warmed = append(warmed, snap.DriverID+":"+snap.Model)
			return nil
		},
		nil,
	)

	err := s.LoadDriver(context.Background(), scheduler.JobSpec{
		DriverID: "vllm", Model: "chat",
		Capability: capability.TextGeneration, RAMMB: 8000,
	})
	if err != nil {
		t.Fatalf("LoadDriver: %v", err)
	}
	if len(warmed) != 1 {
		t.Fatalf("expected warm, got %v", warmed)
	}

	state, err := s.Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Loaded) != 1 || state.Loaded[0].DriverID != "vllm" {
		t.Fatalf("expected loaded slot, got %+v", state.Loaded)
	}

	if err := s.UnloadDriver(context.Background(), "vllm"); err != nil {
		t.Fatalf("UnloadDriver: %v", err)
	}
	if len(unloaded) != 1 {
		t.Fatalf("expected unload, got %v", unloaded)
	}
	state, err = s.Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Loaded) != 0 {
		t.Fatalf("expected empty loaded, got %+v", state.Loaded)
	}
}

func TestSchedulerConcurrentInference(t *testing.T) {
	root := t.TempDir()
	cfg := config.EffectiveResources{
		TotalRAMMB:      64000,
		EvictAndRestore: true,
		LazyRestore:     true,
		CapabilityBudgets: map[string]config.CapabilityBudget{
			string(capability.TextGeneration): {MaxConcurrent: 2},
		},
	}
	s := scheduler.New(root, cfg, nil, nil, nil)

	stateDir := filepath.Join(root, ".wuji", "runtime")
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	state := scheduler.PersistedState{
		Loaded: []scheduler.LoadedSlot{{
			DriverID: "vllm", Model: "chat",
			Capability: capability.TextGeneration, RAMMB: 8000,
		}},
	}
	data, _ := json.MarshalIndent(state, "", "  ")
	if err := os.WriteFile(filepath.Join(stateDir, "scheduler.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	job := scheduler.JobSpec{
		DriverID: "vllm", Model: "chat", Capability: capability.TextGeneration, RAMMB: 8000,
	}

	overlap := atomic.Bool{}
	var running atomic.Int32
	release := make(chan struct{})
	done := make(chan error, 2)

	runOne := func() {
		done <- s.Run(context.Background(), job, func(ctx context.Context) error {
			if running.Add(1) > 1 {
				overlap.Store(true)
			}
			<-release
			running.Add(-1)
			return nil
		})
	}

	go runOne()
	time.Sleep(20 * time.Millisecond)
	go runOne()

	deadline := time.After(2 * time.Second)
	for !overlap.Load() {
		select {
		case <-deadline:
			t.Fatal("expected two text jobs to run concurrently")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	close(release)

	for i := 0; i < 2; i++ {
		if err := <-done; err != nil {
			t.Fatalf("Run: %v", err)
		}
	}

	persisted, err := s.Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(persisted.Active) != 0 {
		t.Fatalf("expected no active jobs after completion, got %+v", persisted.Active)
	}
}

func TestSchedulerConcurrentWithActiveNotLoaded(t *testing.T) {
	root := t.TempDir()
	cfg := config.EffectiveResources{
		TotalRAMMB:      64000,
		EvictAndRestore: true,
		LazyRestore:     true,
		CapabilityBudgets: map[string]config.CapabilityBudget{
			string(capability.TextGeneration): {MaxConcurrent: 2},
		},
	}
	s := scheduler.New(root, cfg, nil, nil, nil)

	job := scheduler.JobSpec{
		DriverID: "vllm", Model: "chat", Capability: capability.TextGeneration, RAMMB: 8000,
	}

	overlap := atomic.Bool{}
	var running atomic.Int32
	release := make(chan struct{})
	done := make(chan error, 2)

	runOne := func() {
		done <- s.Run(context.Background(), job, func(ctx context.Context) error {
			if running.Add(1) > 1 {
				overlap.Store(true)
			}
			<-release
			running.Add(-1)
			return nil
		})
	}

	go runOne()
	time.Sleep(20 * time.Millisecond)
	go runOne()

	deadline := time.After(2 * time.Second)
	for !overlap.Load() {
		select {
		case <-deadline:
			t.Fatal("expected concurrent jobs without a loaded slot entry")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
	close(release)
	for i := 0; i < 2; i++ {
		if err := <-done; err != nil {
			t.Fatalf("Run: %v", err)
		}
	}
}
