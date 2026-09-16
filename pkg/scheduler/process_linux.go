//go:build linux

package scheduler

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// BinaryNewerThanListener reports whether binPath was modified after the
// process listening on port started (typical after rebuilding a driver).
func BinaryNewerThanListener(binPath string, port int) bool {
	fi, err := os.Stat(binPath)
	if err != nil {
		return false
	}
	pids, err := pidsListeningOnPort(port)
	if err != nil || len(pids) == 0 {
		return false
	}
	start, err := processStartTime(pids[0])
	if err != nil {
		return false
	}
	return fi.ModTime().After(start)
}

func pidsListeningOnPort(port int) ([]int, error) {
	path, err := exec.LookPath("fuser")
	if err != nil {
		return nil, err
	}
	out, err := exec.Command(path, fmt.Sprintf("%d/tcp", port)).CombinedOutput()
	text := strings.TrimSpace(string(out))
	if text == "" {
		return nil, nil
	}
	if err != nil && !strings.Contains(text, "/tcp") {
		return nil, err
	}

	var pids []int
	for _, field := range strings.Fields(text) {
		field = strings.TrimSuffix(field, "/tcp")
		if i := strings.Index(field, ":"); i >= 0 {
			field = field[i+1:]
		}
		if pid, convErr := strconv.Atoi(field); convErr == nil {
			pids = append(pids, pid)
		}
	}
	return pids, nil
}

func processStartTime(pid int) (time.Time, error) {
	boot, err := bootTime()
	if err != nil {
		return time.Time{}, err
	}

	stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return time.Time{}, err
	}
	closeParen := strings.LastIndex(string(stat), ")")
	if closeParen < 0 {
		return time.Time{}, fmt.Errorf("parse /proc/%d/stat", pid)
	}
	fields := strings.Fields(string(stat[closeParen+1:]))
	if len(fields) < 20 {
		return time.Time{}, fmt.Errorf("short /proc/%d/stat", pid)
	}
	ticks, err := strconv.ParseInt(fields[19], 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return boot.Add(time.Duration(ticks) * (time.Second / 100)), nil
}

func bootTime() (time.Time, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return time.Time{}, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "btime ") {
			sec, err := strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(line, "btime")), 10, 64)
			if err != nil {
				return time.Time{}, err
			}
			return time.Unix(sec, 0), nil
		}
	}
	return time.Time{}, fmt.Errorf("btime not found in /proc/stat")
}
