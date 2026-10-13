package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReqpackDriverBinaryPathMissing(t *testing.T) {
	if path := ReqpackDriverBinaryPath("missing-driver"); path != "" {
		t.Fatalf("expected empty path, got %q", path)
	}
}

func TestReqpackDriverBinaryPathFound(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	binDir := filepath.Join(home, ".local", "share", "wuji", "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	binPath := filepath.Join(binDir, "wuji-driver-ffmpeg")
	if err := os.WriteFile(binPath, []byte("bin"), 0o755); err != nil {
		t.Fatal(err)
	}

	if got := ReqpackDriverBinaryPath("ffmpeg"); got != binPath {
		t.Fatalf("expected %q, got %q", binPath, got)
	}
}

func TestDriverBinaryPathPrefersReqpackInstall(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	binDir := filepath.Join(home, ".local", "share", "wuji", "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	reqpackBin := filepath.Join(binDir, "wuji-driver-ffmpeg")
	if err := os.WriteFile(reqpackBin, []byte("bin"), 0o755); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{Root: home}
	if got := cfg.DriverBinaryPath("ffmpeg"); got != reqpackBin {
		t.Fatalf("expected reqpack binary %q, got %q", reqpackBin, got)
	}
}
