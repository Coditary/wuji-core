package core

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/catalog"
	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
	apidriver "github.com/coditary/wuji-core/pkg/driver/api"
	"github.com/coditary/wuji-core/pkg/driver/dummy"
	ffmpegdrv "github.com/coditary/wuji-core/pkg/driver/ffmpeg"
	localdrv "github.com/coditary/wuji-core/pkg/driver/local"
	grpcdriver "github.com/coditary/wuji-core/pkg/driver/grpc"
	"github.com/coditary/wuji-core/pkg/scheduler"
)

// Config holds core initialization options.
type Config struct {
	DefaultDriverID string
	AppConfig       *config.Config
	// Lazy skips connecting to remote drivers at startup (daemon mode).
	Lazy bool
}

// Core is the central orchestrator for all driver operations.
type Core struct {
	registry      *driver.Registry
	defaultDriver string
	appConfig     *config.Config
	resources     *scheduler.Scheduler
	idle          *idleTracker
}

// New creates a Core instance with built-in drivers registered.
func New(cfg Config) (*Core, error) {
	appCfg := cfg.AppConfig
	if appCfg == nil {
		var err error
		appCfg, err = config.Load()
		if err != nil {
			return nil, fmt.Errorf("load config: %w", err)
		}
	}

	c := &Core{
		registry:      driver.NewRegistry(),
		defaultDriver: cfg.DefaultDriverID,
		appConfig:     appCfg,
	}

	if err := c.registry.Register(dummy.New()); err != nil {
		return nil, fmt.Errorf("register dummy driver: %w", err)
	}
	if err := c.registry.Register(localdrv.NewCompositeFromConfig(appCfg)); err != nil {
		return nil, fmt.Errorf("register local rag driver: %w", err)
	}
	if err := c.registry.Register(ffmpegdrv.New()); err != nil {
		return nil, fmt.Errorf("register ffmpeg driver: %w", err)
	}
	if err := apidriver.RegisterAll(c.registry, appCfg.APIs); err != nil {
		return nil, fmt.Errorf("register api drivers: %w", err)
	}

	if !cfg.Lazy {
		for _, endpoint := range appCfg.ResolvedDriverEndpoints() {
			if err := c.ConnectRemote(context.Background(), endpoint, false); err != nil {
				driverID := appCfg.DriverIDForEndpoint(endpoint)
				if driverID != "" && appCfg.AutoStartDriverEnabled(driverID) {
					continue
				}
				fmt.Fprintf(os.Stderr, "warning: could not connect to driver at %s: %v\n", endpoint, err)
			}
		}
	}

	if c.defaultDriver == "" {
		if appCfg.DefaultDriver != "" {
			c.defaultDriver = appCfg.DefaultDriver
		} else {
			c.defaultDriver = dummy.DriverID
		}
	}

	if !cfg.Lazy {
		if _, err := c.registry.Get(c.defaultDriver); err != nil {
			return nil, fmt.Errorf("default driver %q: %w", c.defaultDriver, err)
		}
	}

	c.initScheduler()
	c.initIdleShutdown(cfg.Lazy)
	return c, nil
}

// ConnectRemote connects to a driver at the given gRPC endpoint and registers it.
func (c *Core) ConnectRemote(ctx context.Context, endpoint string, persist bool) error {
	remote, err := grpcdriver.Connect(ctx, endpoint)
	if err != nil {
		return err
	}

	if err := c.registry.Register(remote); err != nil {
		_ = remote.Close()
		return err
	}

	if persist {
		if err := c.appConfig.AddDriverEndpoint(endpoint); err != nil {
			return err
		}
	}
	c.touchDriver(remote.Info().ID)
	return nil
}

// AppConfig returns the application configuration.
func (c *Core) AppConfig() *config.Config {
	return c.appConfig
}

func (c *Core) resolve(driverID string) (driver.Driver, error) {
	driverID = localdrv.NormalizeDriverID(driverID)
	if driverID == "" {
		driverID = c.defaultDriver
	}
	driverID = localdrv.NormalizeDriverID(driverID)
	return c.registry.Get(driverID)
}

// Registry returns the driver registry.
func (c *Core) Registry() *driver.Registry {
	return c.registry
}

// DefaultDriverID returns the configured default driver.
func (c *Core) DefaultDriverID() string {
	return c.defaultDriver
}

// ListDrivers returns metadata for all registered drivers.
func (c *Core) ListDrivers() []driver.Info {
	return c.registry.List()
}

// ListProviders returns the text provider catalog for clients (Taiji sync).
// Includes apis.* entries, legacy providers.*, and registered text-capable drivers.
func (c *Core) ListProviders() config.ProviderCatalog {
	summaries := map[string]config.ProviderSummary{}
	if c.appConfig != nil {
		for id, entry := range c.appConfig.APIs {
			summaries[id] = config.ProviderSummary{
				Model:  defaultModelForAPI(entry),
				Driver: id,
				Kind:   "api",
			}
		}
		for id, p := range c.appConfig.Providers {
			if _, exists := summaries[id]; exists {
				continue
			}
			driverID := strings.TrimSpace(p.Driver)
			if driverID == "" {
				driverID = id
			}
			summaries[id] = config.ProviderSummary{
				Model:  p.Model,
				Driver: driverID,
				Kind:   "provider",
			}
		}
	}
	for _, d := range c.registry.List() {
		if !driverSupportsText(d) {
			continue
		}
		if _, exists := summaries[d.ID]; exists {
			continue
		}
		summaries[d.ID] = config.ProviderSummary{
			Model:  c.defaultModelForDriver(d.ID),
			Driver: d.ID,
			Kind:   "driver",
		}
	}

	defaultID := ""
	if c.appConfig != nil {
		defaultID = c.appConfig.DefaultProvider
		if defaultID == "" {
			defaultID = c.appConfig.DefaultDriver
		}
		if defaultID == "" && c.appConfig.CapabilityDrivers != nil {
			defaultID = c.appConfig.CapabilityDrivers["text"]
		}
	}
	if c.appConfig != nil {
		catalog.MergeIntoProviderCatalog(c.appConfig.Root, summaries)
	}
	return config.ProviderCatalog{
		DefaultProvider: defaultID,
		Providers:       summaries,
	}
}

func driverSupportsText(d driver.Info) bool {
	for _, cap := range d.Capabilities {
		if cap == capability.TextGeneration {
			return true
		}
	}
	return false
}

func (c *Core) defaultModelForDriver(driverID string) string {
	if c.appConfig == nil {
		return ""
	}
	switch driverID {
	case "vllm":
		return c.appConfig.ResolvedVLLM().DefaultModel
	case "llama":
		return c.appConfig.ResolvedLlama().DefaultModel
	default:
		return ""
	}
}

func defaultModelForAPI(entry config.APIEntry) string {
	if entry.Text != nil && entry.Text.Model != "" {
		return entry.Text.Model
	}
	if entry.Image != nil && entry.Image.Model != "" {
		return entry.Image.Model
	}
	if entry.Audio != nil && entry.Audio.Model != "" {
		return entry.Audio.Model
	}
	return ""
}

// DefaultProviderID returns the configured default LLM provider id.
func (c *Core) DefaultProviderID() string {
	if c.appConfig == nil {
		return ""
	}
	return c.appConfig.DefaultProvider
}

func (c *Core) GenerateText(ctx context.Context, driverID string, req driver.TextRequest) (*driver.TextResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	start := time.Now()
	streamed := req.Stream
	resolvedID := c.ResolveDriverID(driverID, capability.TextGeneration)
	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = c.defaultModelFor(resolvedID)
	}

	resp, err := scheduleGenerate(c, ctx, driverID, req.Model, capability.TextGeneration, func(ctx context.Context) (*driver.TextResponse, error) {
		mode := req.InputModeOrDefault()
		if mode == driver.TextInputAudio {
			return c.generateTextFromAudio(ctx, driverID, req)
		}
		d, err := c.resolveForCapability(driverID, capability.TextGeneration)
		if err != nil {
			return nil, err
		}
		return driver.RunTextInput(ctx, d, req)
	})
	if err != nil {
		return nil, err
	}
	c.enrichAndLogTextMetrics(resolvedID, model, req, resp, start, streamed)
	return resp, nil
}

func (c *Core) GenerateTextStream(ctx context.Context, driverID string, req driver.TextRequest, onDelta func(string) error) (*driver.TextResponse, error) {
	req.Stream = true
	req.OnDelta = onDelta
	return c.GenerateText(ctx, driverID, req)
}

func (c *Core) GenerateImage(ctx context.Context, driverID string, req driver.ImageRequest) (*driver.ImageResponse, error) {
	return scheduleGenerate(c, ctx, driverID, req.Model, capability.ImageGeneration, func(ctx context.Context) (*driver.ImageResponse, error) {
		if req.TaskOrDefault() == driver.ImageTaskUpscale {
			imgD, err := c.resolveForCapability(driverID, capability.ImageGeneration)
			if err != nil {
				return nil, err
			}
			upD, err := c.resolveForCapability(driverID, capability.ImageUpscale)
			if err != nil {
				upD = imgD
			}
			downD, err := c.resolveForCapability(driverID, capability.ImageDownscale)
			if err != nil {
				downD = imgD
			}
			scaleD, err := c.resolveForCapability(driverID, capability.ImageScale)
			if err != nil {
				scaleD = imgD
			}
			return driver.RunImageScale(ctx, driver.ScaleDrivers{
				Image: imgD, Upscale: upD, Downscale: downD, Scale: scaleD,
			}, req)
		}

		d, err := c.resolveForCapability(driverID, capability.ImageGeneration)
		if err != nil {
			return nil, err
		}
		return driver.RunImageTask(ctx, d, req)
	})
}

func (c *Core) GenerateVideo(ctx context.Context, driverID string, req driver.VideoRequest) (*driver.VideoResponse, error) {
	return scheduleGenerate(c, ctx, driverID, req.Model, capability.VideoGeneration, func(ctx context.Context) (*driver.VideoResponse, error) {
		d, err := c.resolveForCapability(driverID, capability.VideoGeneration)
		if err != nil {
			return nil, err
		}
		return driver.RunVideoTask(ctx, d, req)
	})
}

func (c *Core) GenerateAudio(ctx context.Context, driverID string, req driver.AudioRequest) (*driver.AudioResponse, error) {
	return scheduleGenerate(c, ctx, driverID, req.Model, capability.AudioGeneration, func(ctx context.Context) (*driver.AudioResponse, error) {
		d, err := c.resolveForCapability(driverID, capability.AudioGeneration)
		if err != nil {
			return nil, err
		}
		return driver.RunAudioTask(ctx, d, req)
	})
}

func (c *Core) GenerateMesh(ctx context.Context, driverID string, req driver.MeshRequest) (*driver.MeshResponse, error) {
	return scheduleGenerate(c, ctx, driverID, req.Model, capability.Mesh, func(ctx context.Context) (*driver.MeshResponse, error) {
		d, err := c.resolveForCapability(driverID, capability.Mesh)
		if err != nil {
			return nil, err
		}
		return driver.RunMeshTask(ctx, d, req)
	})
}

func (c *Core) CloneVoice(ctx context.Context, driverID string, req driver.VoiceRequest) (*driver.VoiceResponse, error) {
	return scheduleGenerate(c, ctx, driverID, req.TargetModel, capability.VoiceCloning, func(ctx context.Context) (*driver.VoiceResponse, error) {
		d, err := c.resolveForCapability(driverID, capability.VoiceCloning)
		if err != nil {
			return nil, err
		}
		return driver.RunVoiceTask(ctx, d, req)
	})
}

func (c *Core) ManageDataset(ctx context.Context, driverID string, req driver.DatasetRequest) (*driver.DatasetResponse, error) {
	resolvedID := c.ResolveDriverID(driverID, capability.DatasetMgmt)
	c.beginDriverUse(resolvedID)
	defer c.endDriverUse(resolvedID)
	if err := c.EnsureDriver(ctx, resolvedID); err != nil {
		return nil, err
	}
	d, err := c.resolveForCapability(driverID, capability.DatasetMgmt)
	if err != nil {
		return nil, err
	}
	return driver.RunDatasetTask(ctx, d, req)
}

// Close shuts down the core and all registered drivers.
func (c *Core) Close() error {
	c.stopIdleShutdown()
	return c.registry.Close()
}
