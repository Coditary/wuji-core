package scheduler

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

// StopListenerOnPort terminates processes listening on the given TCP port.
func StopListenerOnPort(port int) error {
	if port <= 0 {
		return nil
	}
	portSpec := fmt.Sprintf("%d/tcp", port)
	if path, err := exec.LookPath("fuser"); err == nil {
		out, err := exec.Command(path, "-k", portSpec).CombinedOutput()
		if err != nil {
			msg := strings.TrimSpace(string(out))
			if msg != "" && !strings.Contains(msg, "No such file") {
				return fmt.Errorf("fuser -k %s: %v (%s)", portSpec, err, msg)
			}
		}
		return nil
	}
	return nil
}

// PortListening reports whether something accepts TCP connections on localhost:port.
func PortListening(port int) bool {
	if port <= 0 {
		return false
	}
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.DialTimeout("tcp", addr, 150*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
