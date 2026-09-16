package config

import (
	"fmt"
	"path/filepath"
)

const (
	defaultA1111Host                = "127.0.0.1"
	defaultA1111Port                = 7860
	defaultA1111TimeoutSeconds      = 600
	defaultA1111OutputDir           = "/tmp/wuji/a1111"
	defaultA1111VideoOutputDir      = "/tmp/wuji/a1111/video"
	defaultA1111DefaultVideoModel   = "t2v"
	defaultA1111DefaultInpaintModel = "sd-v1-5-inpainting.ckpt"
	defaultA1111VideoWidth          = 256
	defaultA1111VideoHeight         = 256
	defaultA1111VideoSteps          = 30
	defaultA1111VideoCFGScale       = 17
	defaultA1111VideoSampler        = "DDIM_Gaussian"
	defaultA1111VideoFPS            = 15
)

// A1111Config holds settings for the Automatic1111 image driver.
type A1111Config struct {
	Host                string `yaml:"host,omitempty"`
	Port                int    `yaml:"port,omitempty"`
	APIBase             string `yaml:"api_base,omitempty"`
	APIKey              string `yaml:"api_key,omitempty"`
	AuthUser            string `yaml:"auth_user,omitempty"`
	AuthPass            string `yaml:"auth_pass,omitempty"`
	DefaultModel        string `yaml:"default_model,omitempty"`
	DefaultInpaintModel string `yaml:"default_inpaint_model,omitempty"`
	DefaultVideoModel   string `yaml:"default_video_model,omitempty"`
	GRPCAddr            string `yaml:"grpc_addr,omitempty"`
	DriverBin           string `yaml:"driver_bin,omitempty"`
	WebUIDir            string `yaml:"webui_dir,omitempty"`
	TimeoutSeconds      int    `yaml:"timeout_seconds,omitempty"`
	OutputDir           string `yaml:"output_dir,omitempty"`
	VideoOutputDir      string `yaml:"video_output_dir,omitempty"`
	VideoWidth          int    `yaml:"video_width,omitempty"`
	VideoHeight         int    `yaml:"video_height,omitempty"`
	VideoSteps          int    `yaml:"video_steps,omitempty"`
	VideoCFGScale       int    `yaml:"video_cfg_scale,omitempty"`
	VideoSampler        string `yaml:"video_sampler,omitempty"`
	VideoFPS            int    `yaml:"video_fps,omitempty"`
	AutoStartDriver     *bool  `yaml:"auto_start_driver,omitempty"`
	ManageServer        *bool  `yaml:"manage_server,omitempty"`
}

// DefaultA1111Config returns built-in defaults (no file required).
func DefaultA1111Config() A1111Config {
	return A1111Config{
		Host:                defaultA1111Host,
		Port:                defaultA1111Port,
		TimeoutSeconds:      defaultA1111TimeoutSeconds,
		OutputDir:           defaultA1111OutputDir,
		VideoOutputDir:      defaultA1111VideoOutputDir,
		DefaultVideoModel:   defaultA1111DefaultVideoModel,
		DefaultInpaintModel: defaultA1111DefaultInpaintModel,
		VideoWidth:          defaultA1111VideoWidth,
		VideoHeight:         defaultA1111VideoHeight,
		VideoSteps:          defaultA1111VideoSteps,
		VideoCFGScale:       defaultA1111VideoCFGScale,
		VideoSampler:        defaultA1111VideoSampler,
		VideoFPS:            defaultA1111VideoFPS,
	}
}

// ResolvedA1111 returns A1111 settings merged from defaults and .wuji/config.yaml.
func (c *Config) ResolvedA1111() A1111Config {
	out := DefaultA1111Config()
	if c == nil || c.Drivers.A1111 == nil {
		return out
	}
	mergeA1111(&out, *c.Drivers.A1111)
	return out
}

func mergeA1111(dst *A1111Config, src A1111Config) {
	if src.Host != "" {
		dst.Host = src.Host
	}
	if src.Port != 0 {
		dst.Port = src.Port
	}
	if src.APIBase != "" {
		dst.APIBase = src.APIBase
	}
	if src.APIKey != "" {
		dst.APIKey = src.APIKey
	}
	if src.AuthUser != "" {
		dst.AuthUser = src.AuthUser
	}
	if src.AuthPass != "" {
		dst.AuthPass = src.AuthPass
	}
	if src.DefaultModel != "" {
		dst.DefaultModel = src.DefaultModel
	}
	if src.DefaultInpaintModel != "" {
		dst.DefaultInpaintModel = src.DefaultInpaintModel
	}
	if src.DefaultVideoModel != "" {
		dst.DefaultVideoModel = src.DefaultVideoModel
	}
	if src.GRPCAddr != "" {
		dst.GRPCAddr = src.GRPCAddr
	}
	if src.TimeoutSeconds != 0 {
		dst.TimeoutSeconds = src.TimeoutSeconds
	}
	if src.OutputDir != "" {
		dst.OutputDir = src.OutputDir
	}
	if src.VideoOutputDir != "" {
		dst.VideoOutputDir = src.VideoOutputDir
	}
	if src.VideoWidth != 0 {
		dst.VideoWidth = src.VideoWidth
	}
	if src.VideoHeight != 0 {
		dst.VideoHeight = src.VideoHeight
	}
	if src.VideoSteps != 0 {
		dst.VideoSteps = src.VideoSteps
	}
	if src.VideoCFGScale != 0 {
		dst.VideoCFGScale = src.VideoCFGScale
	}
	if src.VideoSampler != "" {
		dst.VideoSampler = src.VideoSampler
	}
	if src.VideoFPS != 0 {
		dst.VideoFPS = src.VideoFPS
	}
	if src.DriverBin != "" {
		dst.DriverBin = src.DriverBin
	}
	if src.WebUIDir != "" {
		dst.WebUIDir = src.WebUIDir
	}
	if src.AutoStartDriver != nil {
		dst.AutoStartDriver = src.AutoStartDriver
	}
	if src.ManageServer != nil {
		dst.ManageServer = src.ManageServer
	}
}

// ResolvedAPIBase returns the A1111 WebUI API base URL.
func (c A1111Config) ResolvedAPIBase() string {
	if c.APIBase != "" {
		if len(c.APIBase) > 0 && c.APIBase[len(c.APIBase)-1] == '/' {
			return c.APIBase[:len(c.APIBase)-1]
		}
		return c.APIBase
	}
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
}

func (c A1111Config) ManageServerEnabled() bool {
	if c.ManageServer == nil {
		return true
	}
	return *c.ManageServer
}

// ResolvedWebUIDir returns the Automatic1111 installation directory.
func (c *Config) ResolvedWebUIDir() string {
	if c != nil && c.Drivers.A1111 != nil && c.Drivers.A1111.WebUIDir != "" {
		return c.Drivers.A1111.WebUIDir
	}
	if c != nil && c.Root != "" {
		return filepath.Join(c.Root, "..", "..", "stable-diffusion-webui")
	}
	return "../../stable-diffusion-webui"
}
