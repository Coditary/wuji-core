package core

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/coditary/wuji-core/pkg/driver/dummy"
	ffmpegdrv "github.com/coditary/wuji-core/pkg/driver/ffmpeg"
	"github.com/coditary/wuji-core/pkg/driver/local"
	"github.com/coditary/wuji-core/pkg/netx"
)

const idleCheckInterval = 30 * time.Second

type driverActivity struct {
	lastUsed   time.Time
	endpoint   string
	activeOps  int // >0 while a wuji command holds the driver (generate, train, load, …)
	activeSince time.Time
}

type idleTracker struct {
	mu       sync.Mutex
	activity map[string]driverActivity
	stopCh   chan struct{}
}

func isBuiltinDriver(driverID string) bool {
	switch driverID {
	case dummy.DriverID, local.DriverID, ffmpegdrv.DriverID:
		return true
	default:
		return false
	}
}

func (c *Core) initIdleShutdown(lazy bool) {
	if !lazy || os.Getenv("WUJI_EMBEDDED") == "1" {
		return
	}
	secs := c.appConfig.ResolvedResources().IdleShutdownSeconds
	if secs <= 0 {
		return
	}

	c.idle = &idleTracker{
		activity: make(map[string]driverActivity),
		stopCh:   make(chan struct{}),
	}
	stopCh := c.idle.stopCh
	go c.idleMonitor(time.Duration(secs)*time.Second, stopCh)
}

func (c *Core) stopIdleShutdown() {
	if c.idle == nil {
		return
	}
	close(c.idle.stopCh)
	c.idle = nil
}

// beginDriverUse marks a driver as busy for the lifetime of a wuji operation.
// Idle shutdown is suppressed until every matching endDriverUse runs — including
// multi-hour training or generation calls.
func (c *Core) beginDriverUse(driverID string) {
	if c.idle == nil || isBuiltinDriver(driverID) {
		return
	}
	c.idle.mu.Lock()
	defer c.idle.mu.Unlock()
	act := c.idle.activity[driverID]
	act.activeOps++
	if act.activeOps == 1 {
		act.activeSince = time.Now()
	}
	act.lastUsed = time.Now()
	if act.endpoint == "" {
		act.endpoint = c.driverEndpoint(driverID)
	}
	c.idle.activity[driverID] = act
}

func (c *Core) endDriverUse(driverID string) {
	if c.idle == nil || isBuiltinDriver(driverID) {
		return
	}
	c.idle.mu.Lock()
	defer c.idle.mu.Unlock()
	act, ok := c.idle.activity[driverID]
	if !ok {
		return
	}
	if act.activeOps > 0 {
		act.activeOps--
	}
	if act.activeOps == 0 {
		act.activeSince = time.Time{}
	}
	act.lastUsed = time.Now()
	c.idle.activity[driverID] = act
}

func (c *Core) driverActiveOps(driverID string) int {
	if c.idle == nil {
		return 0
	}
	c.idle.mu.Lock()
	defer c.idle.mu.Unlock()
	return c.idle.activity[driverID].activeOps
}

func (c *Core) touchDriver(driverID string) {
	if c.idle == nil || isBuiltinDriver(driverID) {
		return
	}
	c.idle.mu.Lock()
	defer c.idle.mu.Unlock()
	act := c.idle.activity[driverID]
	act.lastUsed = time.Now()
	if act.endpoint == "" {
		act.endpoint = c.driverEndpoint(driverID)
	}
	c.idle.activity[driverID] = act
}

func (c *Core) driverEndpoint(driverID string) string {
	if d, err := c.registry.Get(driverID); err == nil && d.Info().Endpoint != "" {
		return d.Info().Endpoint
	}
	if c.appConfig != nil {
		return c.appConfig.DriverGRPCAddr(driverID)
	}
	return ""
}

func (c *Core) idleMonitor(idleTimeout time.Duration, stopCh <-chan struct{}) {
	ticker := time.NewTicker(idleCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			c.checkIdleDrivers(idleTimeout)
		}
	}
}

func (c *Core) checkIdleDrivers(idleTimeout time.Duration) {
	if c.idle == nil {
		return
	}
	now := time.Now()
	var candidates []string

	c.idle.mu.Lock()
	for id, act := range c.idle.activity {
		if act.activeOps > 0 {
			continue // command still running — may take hours
		}
		if now.Sub(act.lastUsed) >= idleTimeout {
			candidates = append(candidates, id)
		}
	}
	c.idle.mu.Unlock()

	for _, id := range candidates {
		c.shutdownIdleDriver(context.Background(), id)
	}
}

func (c *Core) shutdownIdleDriver(ctx context.Context, driverID string) {
	if isBuiltinDriver(driverID) {
		return
	}

	c.idle.mu.Lock()
	act, ok := c.idle.activity[driverID]
	if !ok || act.activeOps > 0 {
		c.idle.mu.Unlock()
		return
	}
	endpoint := act.endpoint
	if endpoint == "" {
		endpoint = c.driverEndpoint(driverID)
	}
	delete(c.idle.activity, driverID)
	c.idle.mu.Unlock()

	_ = c.unloadSlot(ctx, driverID)

	d, err := c.registry.Get(driverID)
	if err != nil {
		if endpoint != "" {
			_ = netx.StopListenerOnEndpoint(endpoint)
		}
		return
	}
	if !d.Info().Remote {
		return
	}
	if endpoint == "" {
		endpoint = d.Info().Endpoint
	}
	_ = d.Close()
	_ = c.registry.Unregister(driverID)
	if endpoint != "" {
		_ = netx.StopListenerOnEndpoint(endpoint)
	}
	fmt.Fprintf(os.Stderr, "idle shutdown: stopped driver %q after inactivity\n", driverID)
}
