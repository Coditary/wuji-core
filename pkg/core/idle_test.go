package core

import (
	"testing"
	"time"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/driver/dummy"
	ffmpegdrv "github.com/coditary/wuji-core/pkg/driver/ffmpeg"
	"github.com/coditary/wuji-core/pkg/driver/local"
)

func TestIsBuiltinDriver(t *testing.T) {
	if !isBuiltinDriver(dummy.DriverID) || !isBuiltinDriver(local.DriverID) || !isBuiltinDriver(ffmpegdrv.DriverID) {
		t.Fatal("expected built-in drivers")
	}
	if isBuiltinDriver("vllm") {
		t.Fatal("vllm should not be built-in")
	}
}

func TestIdleTrackerActivity(t *testing.T) {
	c := &Core{
		registry:  driver.NewRegistry(),
		appConfig: &config.Config{},
		idle: &idleTracker{
			activity: make(map[string]driverActivity),
			stopCh:   make(chan struct{}),
		},
	}

	c.beginDriverUse("vllm")
	act := c.idle.activity["vllm"]
	if act.activeOps != 1 {
		t.Fatalf("expected activeOps=1, got %d", act.activeOps)
	}

	c.beginDriverUse("vllm")
	act = c.idle.activity["vllm"]
	if act.activeOps != 2 {
		t.Fatalf("expected activeOps=2, got %d", act.activeOps)
	}

	c.endDriverUse("vllm")
	act = c.idle.activity["vllm"]
	if act.activeOps != 1 {
		t.Fatalf("expected activeOps=1 after end, got %d", act.activeOps)
	}

	c.touchDriver("vllm")
	if c.idle.activity["vllm"].lastUsed.IsZero() {
		t.Fatal("expected lastUsed to be set")
	}

	c.beginDriverUse(dummy.DriverID)
	if _, ok := c.idle.activity[dummy.DriverID]; ok {
		t.Fatal("built-in drivers should not be tracked")
	}
}

func TestCheckIdleDriversSkipsActiveOps(t *testing.T) {
	c := &Core{
		idle: &idleTracker{
			activity: map[string]driverActivity{
				"vllm": {lastUsed: time.Now().Add(-2 * time.Hour), activeOps: 1},
			},
			stopCh: make(chan struct{}),
		},
	}

	c.checkIdleDrivers(time.Second)
	if _, ok := c.idle.activity["vllm"]; !ok {
		t.Fatal("active driver should not be shut down")
	}
}

func TestCheckIdleDriversShutsDownAfterIdlePeriod(t *testing.T) {
	c := &Core{
		registry:  driver.NewRegistry(),
		appConfig: &config.Config{},
		idle: &idleTracker{
			activity: map[string]driverActivity{
				"vllm": {lastUsed: time.Now().Add(-2 * time.Hour), activeOps: 0, endpoint: "unix:///tmp/vllm.sock"},
			},
			stopCh: make(chan struct{}),
		},
	}

	c.checkIdleDrivers(time.Second)
	if _, ok := c.idle.activity["vllm"]; ok {
		t.Fatal("expected idle driver to be removed from activity map")
	}
}

func TestInitIdleShutdownEmbedded(t *testing.T) {
	t.Setenv("WUJI_EMBEDDED", "1")
	c := &Core{appConfig: &config.Config{}}
	c.initIdleShutdown(true)
	if c.idle != nil {
		t.Fatal("expected no idle tracker in embedded mode")
	}
}

func TestInitIdleShutdownDisabled(t *testing.T) {
	t.Setenv("WUJI_EMBEDDED", "")
	c := &Core{
		appConfig: &config.Config{
			Drivers: config.DriversConfig{
				Resources: &config.ResourcesConfig{
					IdleShutdownSeconds: intPtr(0),
				},
			},
		},
	}
	c.initIdleShutdown(true)
	if c.idle != nil {
		t.Fatal("expected no idle tracker when idle_shutdown_seconds=0")
	}
}

func TestInitIdleShutdownEnabled(t *testing.T) {
	t.Setenv("WUJI_EMBEDDED", "")
	c := &Core{appConfig: &config.Config{}}
	c.initIdleShutdown(true)
	if c.idle == nil {
		t.Fatal("expected idle tracker in daemon mode")
	}
	c.stopIdleShutdown()
}

func intPtr(n int) *int { return &n }
