package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/config"
)

func newJobID() string {
	return fmt.Sprintf("%d-%d", os.Getpid(), time.Now().UnixNano())
}

const runtimeDir = "runtime"
const stateFile = "scheduler.json"
const lockFile = "scheduler.lock"

// JobSpec describes a scheduled workload for resource planning.
type JobSpec struct {
	ID         string
	DriverID   string
	Model      string
	Capability capability.Type
	VRAMMB     int
	RAMMB      int
	Transient  bool // unload after completion (e.g. image generation)
}

// SessionSnapshot is a saved model slot to restore after a conflicting job.
type SessionSnapshot struct {
	DriverID   string          `json:"driver_id"`
	Model      string          `json:"model"`
	Capability capability.Type `json:"capability"`
	VRAMMB     int             `json:"vram_mb"`
	RAMMB      int             `json:"ram_mb"`
	SavedAt    time.Time       `json:"saved_at"`
}

// LoadedSlot tracks a model currently occupying memory.
type LoadedSlot struct {
	DriverID   string          `json:"driver_id"`
	Model      string          `json:"model"`
	Capability capability.Type `json:"capability"`
	VRAMMB     int             `json:"vram_mb"`
	RAMMB      int             `json:"ram_mb"`
	LoadedAt   time.Time       `json:"loaded_at"`
}

// PersistedState is shared across CLI invocations via a lock file.
type PersistedState struct {
	Loaded       []LoadedSlot      `json:"loaded"`
	Active       []ActiveJob       `json:"active,omitempty"`
	RestoreStack []SessionSnapshot `json:"restore_stack"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// ActiveJob tracks an in-flight inference job while the scheduler lock is released.
type ActiveJob struct {
	ID         string          `json:"id"`
	DriverID   string          `json:"driver_id"`
	Model      string          `json:"model"`
	Capability capability.Type `json:"capability"`
	StartedAt  time.Time       `json:"started_at"`
}

// UnloadFunc releases VRAM for a driver.
type UnloadFunc func(ctx context.Context, driverID string) error

// WarmFunc preloads a model after restore.
type WarmFunc func(ctx context.Context, snap SessionSnapshot) error

// LoadedProbe reports whether a driver's inference server is actually running.
type LoadedProbe func(driverID string) bool

// Scheduler coordinates VRAM usage and job ordering.
type Scheduler struct {
	cfg    config.EffectiveResources
	root   string
	mu     sync.Mutex
	unload UnloadFunc
	warm   WarmFunc
	probe  LoadedProbe
}

// New creates a scheduler. Pass nil unload/warm/probe when resources are disabled.
func New(root string, cfg config.EffectiveResources, unload UnloadFunc, warm WarmFunc, probe LoadedProbe) *Scheduler {
	return &Scheduler{cfg: cfg, root: root, unload: unload, warm: warm, probe: probe}
}

// Enabled reports whether scheduling is active.
func (s *Scheduler) Enabled() bool {
	return s != nil && s.cfg.Enabled()
}

// Run executes fn with resource planning. The global lock is held only for
// planning/bookkeeping; inference runs concurrently up to MaxConcurrent per capability.
func (s *Scheduler) Run(ctx context.Context, job JobSpec, fn func(context.Context) error) error {
	if !s.Enabled() {
		return fn(ctx)
	}
	if strings.TrimSpace(job.ID) == "" {
		job.ID = newJobID()
	}

	var evicted []SessionSnapshot
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		ev, wait, err := s.reserveJob(ctx, job)
		if wait {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		if err != nil {
			return err
		}
		evicted = ev
		break
	}

	runErr := fn(ctx)

	if postErr := s.completeJob(ctx, job, evicted); postErr != nil {
		if runErr == nil {
			return postErr
		}
		return fmt.Errorf("%v; %w", runErr, postErr)
	}
	return runErr
}

func (s *Scheduler) reserveJob(ctx context.Context, job JobSpec) (evicted []SessionSnapshot, wait bool, err error) {
	release, err := acquireLock(s.runtimePath())
	if err != nil {
		return nil, false, err
	}

	state, err := loadState(s.statePath())
	if err != nil {
		release()
		return nil, false, err
	}
	s.reconcileActive(&state)
	s.reconcile(&state)

	if s.activeCount(state, job.Capability) >= s.maxConcurrent(job.Capability) {
		release()
		return nil, true, nil
	}

	s.dropRestoreMatch(&state, job)
	if s.modelInUse(state, job) {
		evicted = nil
	} else {
		evicted, err = s.prepare(ctx, &state, job)
		if err != nil {
			release()
			return nil, false, err
		}
	}

	state.Active = append(state.Active, ActiveJob{
		ID:         job.ID,
		DriverID:   job.DriverID,
		Model:      job.Model,
		Capability: job.Capability,
		StartedAt:  time.Now(),
	})
	if err := saveState(s.statePath(), state); err != nil {
		release()
		return nil, false, err
	}
	release()
	return evicted, false, nil
}

func (s *Scheduler) completeJob(ctx context.Context, job JobSpec, evicted []SessionSnapshot) error {
	release, err := acquireLock(s.runtimePath())
	if err != nil {
		return err
	}
	defer release()

	state, err := loadState(s.statePath())
	if err != nil {
		return err
	}
	s.reconcileActive(&state)
	s.reconcile(&state)
	s.removeActive(&state, job.ID)

	s.mu.Lock()
	defer s.mu.Unlock()

	var finishErr error

	s.upsertSlot(&state, LoadedSlot{
		DriverID:   job.DriverID,
		Model:      job.Model,
		Capability: job.Capability,
		VRAMMB:     job.VRAMMB,
		RAMMB:      job.RAMMB,
		LoadedAt:   time.Now(),
	})

	if job.Transient && s.cfg.UnloadTransientModels {
		if s.unload != nil {
			if err := s.unload(ctx, job.DriverID); err != nil {
				finishErr = fmt.Errorf("unload transient model: %w", err)
			}
		}
		s.removeSlot(&state, job.DriverID, job.Model)
	}

	if s.cfg.EvictAndRestore && len(evicted) > 0 && !s.cfg.LazyRestore {
		for i := len(evicted) - 1; i >= 0; i-- {
			snap := evicted[i]
			if s.warm != nil {
				if err := s.warm(ctx, snap); err != nil {
					if finishErr == nil {
						finishErr = fmt.Errorf("restore %s/%s: %w", snap.DriverID, snap.Model, err)
					} else {
						finishErr = fmt.Errorf("%v; restore %s/%s: %w", finishErr, snap.DriverID, snap.Model, err)
					}
					continue
				}
			}
			s.upsertSlot(&state, LoadedSlot{
				DriverID:   snap.DriverID,
				Model:      snap.Model,
				Capability: snap.Capability,
				VRAMMB:     snap.VRAMMB,
				RAMMB:      snap.RAMMB,
				LoadedAt:   time.Now(),
			})
		}
		state.RestoreStack = nil
	}

	state.UpdatedAt = time.Now()
	if err := saveState(s.statePath(), state); err != nil {
		if finishErr == nil {
			return err
		}
		return fmt.Errorf("%v; save scheduler state: %w", finishErr, err)
	}
	return finishErr
}

func (s *Scheduler) maxConcurrent(cap capability.Type) int {
	return s.cfg.MaxConcurrent(cap)
}

func (s *Scheduler) activeCount(state PersistedState, cap capability.Type) int {
	n := 0
	for _, job := range state.Active {
		if job.Capability == cap {
			n++
		}
	}
	return n
}

func (s *Scheduler) removeActive(state *PersistedState, id string) {
	if id == "" {
		return
	}
	filtered := make([]ActiveJob, 0, len(state.Active))
	for _, job := range state.Active {
		if job.ID != id {
			filtered = append(filtered, job)
		}
	}
	state.Active = filtered
}

func (s *Scheduler) reconcileActive(state *PersistedState) {
	if len(state.Active) == 0 {
		return
	}
	cutoff := time.Now().Add(-2 * time.Hour)
	filtered := make([]ActiveJob, 0, len(state.Active))
	for _, job := range state.Active {
		if job.StartedAt.IsZero() || job.StartedAt.Before(cutoff) {
			continue
		}
		filtered = append(filtered, job)
	}
	state.Active = filtered
}

// Status returns a snapshot of persisted scheduler state.
func (s *Scheduler) Status() (PersistedState, error) {
	if s == nil {
		return PersistedState{}, nil
	}
	return loadState(s.statePath())
}

// LoadDriver preloads a model, evicting others if the budget requires it.
func (s *Scheduler) LoadDriver(ctx context.Context, job JobSpec) error {
	if s == nil {
		return fmt.Errorf("scheduler not initialized")
	}
	if job.DriverID == "" {
		return fmt.Errorf("driver id required")
	}

	if !s.Enabled() {
		if s.warm == nil {
			return fmt.Errorf("resource management disabled")
		}
		return s.warm(ctx, s.snapshotFromJob(job))
	}

	release, err := acquireLock(s.runtimePath())
	if err != nil {
		return err
	}
	defer release()

	state, err := loadState(s.statePath())
	if err != nil {
		return err
	}
	s.reconcile(&state)
	s.dropRestoreMatch(&state, job)

	if _, err := s.prepare(ctx, &state, job); err != nil {
		return err
	}

	if s.warm != nil {
		if err := s.warm(ctx, s.snapshotFromJob(job)); err != nil {
			return err
		}
	}

	s.upsertSlot(&state, LoadedSlot{
		DriverID:   job.DriverID,
		Model:      job.Model,
		Capability: job.Capability,
		VRAMMB:     job.VRAMMB,
		RAMMB:      job.RAMMB,
		LoadedAt:   time.Now(),
	})
	state.UpdatedAt = time.Now()
	return saveState(s.statePath(), state)
}

// UnloadDriver stops inference for one driver and clears its scheduler slots.
func (s *Scheduler) UnloadDriver(ctx context.Context, driverID string) error {
	if s == nil {
		return fmt.Errorf("scheduler not initialized")
	}
	if driverID == "" {
		return fmt.Errorf("driver id required")
	}

	release, err := acquireLock(s.runtimePath())
	if err != nil {
		return err
	}
	defer release()

	state, err := loadState(s.statePath())
	if err != nil {
		return err
	}
	s.reconcile(&state)

	loaded := s.driverTracked(state, driverID)
	running := s.probe != nil && s.probe(driverID)
	if !loaded && !running {
		return nil
	}

	if s.unload != nil {
		if err := s.unload(ctx, driverID); err != nil {
			return err
		}
	}
	s.removeDriverFromState(&state, driverID)
	state.UpdatedAt = time.Now()
	return saveState(s.statePath(), state)
}

func (s *Scheduler) snapshotFromJob(job JobSpec) SessionSnapshot {
	vram := job.VRAMMB
	if vram <= 0 {
		vram = s.cfg.EstimateModelVRAM(job.Capability, job.Model)
	}
	ram := job.RAMMB
	if ram <= 0 {
		ram = s.cfg.EstimateModelRAM(job.Capability, job.Model)
	}
	return SessionSnapshot{
		DriverID:   job.DriverID,
		Model:      job.Model,
		Capability: job.Capability,
		VRAMMB:     vram,
		RAMMB:      ram,
		SavedAt:    time.Now(),
	}
}

func (s *Scheduler) driverTracked(state PersistedState, driverID string) bool {
	for _, slot := range state.Loaded {
		if slot.DriverID == driverID {
			return true
		}
	}
	return false
}

func (s *Scheduler) removeDriverFromState(state *PersistedState, driverID string) {
	filtered := make([]LoadedSlot, 0, len(state.Loaded))
	for _, slot := range state.Loaded {
		if slot.DriverID != driverID {
			filtered = append(filtered, slot)
		}
	}
	state.Loaded = filtered

	kept := make([]SessionSnapshot, 0, len(state.RestoreStack))
	for _, snap := range state.RestoreStack {
		if snap.DriverID != driverID {
			kept = append(kept, snap)
		}
	}
	state.RestoreStack = kept
}

func (s *Scheduler) runtimePath() string {
	return filepath.Join(s.root, ".wuji", runtimeDir)
}

func (s *Scheduler) statePath() string {
	return filepath.Join(s.runtimePath(), stateFile)
}

func (s *Scheduler) reconcile(state *PersistedState) {
	if s.probe == nil {
		return
	}
	filtered := make([]LoadedSlot, 0, len(state.Loaded))
	for _, slot := range state.Loaded {
		if s.probe(slot.DriverID) {
			filtered = append(filtered, slot)
		}
	}
	state.Loaded = filtered
}

func (s *Scheduler) dropRestoreMatch(state *PersistedState, job JobSpec) {
	if len(state.RestoreStack) == 0 {
		return
	}
	kept := make([]SessionSnapshot, 0, len(state.RestoreStack))
	for _, snap := range state.RestoreStack {
		if snap.DriverID == job.DriverID && snap.Model == job.Model {
			continue
		}
		kept = append(kept, snap)
	}
	state.RestoreStack = kept
}

func (s *Scheduler) prepare(ctx context.Context, state *PersistedState, job JobSpec) ([]SessionSnapshot, error) {
	if s.modelInUse(*state, job) {
		return nil, nil
	}

	vramNeeded := job.VRAMMB
	if vramNeeded <= 0 {
		vramNeeded = s.cfg.EstimateModelVRAM(job.Capability, job.Model)
	}
	ramNeeded := job.RAMMB
	if ramNeeded <= 0 {
		ramNeeded = s.cfg.EstimateModelRAM(job.Capability, job.Model)
	}
	job.VRAMMB = vramNeeded
	job.RAMMB = ramNeeded

	var evicted []SessionSnapshot
	for !s.canFit(state, job, vramNeeded, ramNeeded) {
		slot, ok := s.pickEvictionCandidate(state, job)
		if !ok {
			return evicted, fmt.Errorf(
				"insufficient memory: need %d MB VRAM and %d MB RAM for %s/%s (budgets %d MB VRAM, %d MB RAM)",
				vramNeeded, ramNeeded, job.DriverID, displayModel(job.Model),
				s.cfg.TotalVRAMMB, s.cfg.TotalRAMMB,
			)
		}
		snap := SessionSnapshot{
			DriverID:   slot.DriverID,
			Model:      slot.Model,
			Capability: slot.Capability,
			VRAMMB:     slot.VRAMMB,
			RAMMB:      s.slotRAM(slot),
			SavedAt:    time.Now(),
		}
		evicted = append(evicted, snap)
		state.RestoreStack = append(state.RestoreStack, snap)
		s.removeSlot(state, slot.DriverID, slot.Model)
		if s.unload != nil {
			if err := s.unload(ctx, slot.DriverID); err != nil {
				return evicted, fmt.Errorf("unload %s: %w", slot.DriverID, err)
			}
		}
	}

	// Same driver but different model — unload first.
	for _, slot := range state.Loaded {
		if slot.DriverID == job.DriverID && slot.Model != job.Model {
			evicted = append(evicted, SessionSnapshot{
				DriverID: slot.DriverID, Model: slot.Model,
				Capability: slot.Capability, VRAMMB: slot.VRAMMB,
				RAMMB: s.slotRAM(slot), SavedAt: time.Now(),
			})
			state.RestoreStack = append(state.RestoreStack, evicted[len(evicted)-1])
			s.removeSlot(state, slot.DriverID, slot.Model)
			if s.unload != nil {
				if err := s.unload(ctx, slot.DriverID); err != nil {
					return evicted, err
				}
			}
		}
	}

	return evicted, nil
}

func (s *Scheduler) slotRAM(slot LoadedSlot) int {
	if slot.RAMMB > 0 {
		return slot.RAMMB
	}
	return s.cfg.EstimateModelRAM(slot.Capability, slot.Model)
}

func (s *Scheduler) canFit(state *PersistedState, job JobSpec, vramNeeded, ramNeeded int) bool {
	vramUsed, ramUsed := s.usedMemory(state, job)
	if s.cfg.TotalVRAMMB > 0 && vramUsed+vramNeeded > s.cfg.TotalVRAMMB {
		return false
	}
	if s.cfg.TotalRAMMB > 0 && ramUsed+ramNeeded > s.cfg.TotalRAMMB {
		return false
	}
	return true
}

func (s *Scheduler) usedMemory(state *PersistedState, job JobSpec) (vramUsed, ramUsed int) {
	for _, slot := range state.Loaded {
		if slot.DriverID == job.DriverID && slot.Model == job.Model {
			continue
		}
		vramUsed += slot.VRAMMB
		ramUsed += s.slotRAM(slot)
	}
	return vramUsed, ramUsed
}

func (s *Scheduler) pickEvictionCandidate(state *PersistedState, job JobSpec) (LoadedSlot, bool) {
	// Prefer evicting other drivers, then oldest slot.
	var fallback LoadedSlot
	var found bool
	oldest := time.Now()
	for _, slot := range state.Loaded {
		if slot.DriverID == job.DriverID && slot.Model == job.Model {
			continue
		}
		if slot.DriverID != job.DriverID {
			return slot, true
		}
		if !found || slot.LoadedAt.Before(oldest) {
			fallback = slot
			oldest = slot.LoadedAt
			found = true
		}
	}
	return fallback, found
}

func (s *Scheduler) upsertSlot(state *PersistedState, slot LoadedSlot) {
	for i, existing := range state.Loaded {
		if existing.DriverID == slot.DriverID && existing.Model == slot.Model {
			state.Loaded[i] = slot
			return
		}
	}
	state.Loaded = append(state.Loaded, slot)
}

func (s *Scheduler) removeSlot(state *PersistedState, driverID, model string) {
	filtered := make([]LoadedSlot, 0, len(state.Loaded))
	for _, slot := range state.Loaded {
		if slot.DriverID == driverID && slot.Model == model {
			continue
		}
		filtered = append(filtered, slot)
	}
	state.Loaded = filtered
}

func (s *Scheduler) slotLoaded(state PersistedState, job JobSpec) bool {
	for _, slot := range state.Loaded {
		if slot.DriverID == job.DriverID && slot.Model == job.Model {
			return true
		}
	}
	return false
}

func (s *Scheduler) modelInUse(state PersistedState, job JobSpec) bool {
	if s.slotLoaded(state, job) {
		return true
	}
	for _, active := range state.Active {
		if active.DriverID == job.DriverID &&
			active.Model == job.Model &&
			active.Capability == job.Capability {
			return true
		}
	}
	return false
}

func displayModel(model string) string {
	if model == "" {
		return "(default)"
	}
	return model
}

func loadState(path string) (PersistedState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return PersistedState{}, nil
		}
		return PersistedState{}, err
	}
	var state PersistedState
	if err := json.Unmarshal(data, &state); err != nil {
		return PersistedState{}, fmt.Errorf("parse scheduler state: %w", err)
	}
	return state, nil
}

func saveState(path string, state PersistedState) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	state.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
