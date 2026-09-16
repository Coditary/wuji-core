package netx_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/coditary/wuji-core/pkg/netx"
)

func TestListenUnixSocket(t *testing.T) {
	dir := t.TempDir()
	socket := filepath.Join(dir, "test.sock")
	endpoint := "unix://" + socket

	lis, dialTarget, err := netx.Listen(endpoint)
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	defer lis.Close()

	if dialTarget != endpoint {
		t.Fatalf("dial target %q want %q", dialTarget, endpoint)
	}
	if !netx.EndpointListening(endpoint) {
		t.Fatal("expected socket to be listening")
	}
}

func TestDialTargetNormalizesUnix(t *testing.T) {
	got := netx.DialTarget("unix:/tmp/foo.sock")
	if !strings.HasPrefix(got, "unix://") {
		t.Fatalf("expected unix:// prefix, got %q", got)
	}
}
