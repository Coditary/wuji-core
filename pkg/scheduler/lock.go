package scheduler

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// acquireLock blocks until the scheduler lock is available.
func acquireLock(runtimeDir string) (func(), error) {
	if err := os.MkdirAll(runtimeDir, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(runtimeDir, lockFile)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		return nil, fmt.Errorf("acquire scheduler lock: %w", err)
	}

	return func() {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
	}, nil
}
