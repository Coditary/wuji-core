package sysmem

import (
	"bufio"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// TotalRAMMB returns total system RAM in megabytes, or 0 if unknown.
func TotalRAMMB() int {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				return 0
			}
			kb, err := strconv.Atoi(fields[1])
			if err != nil {
				return 0
			}
			return kb / 1024
		}
	}
	return 0
}

// TotalVRAMMB returns total GPU VRAM in megabytes from the first NVIDIA GPU, or 0.
func TotalVRAMMB() int {
	path, err := exec.LookPath("nvidia-smi")
	if err != nil {
		return 0
	}
	out, err := exec.Command(path,
		"--query-gpu=memory.total",
		"--format=csv,noheader,nounits",
	).Output()
	if err != nil {
		return 0
	}
	line := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	mb, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil || mb <= 0 {
		return 0
	}
	return mb
}

// SuggestedBudget returns conservative RAM and VRAM budgets (0 if undetected).
func SuggestedBudget() (ramMB, vramMB int) {
	totalRAM := TotalRAMMB()
	if totalRAM > 0 {
		ramMB = totalRAM * 80 / 100
	}
	totalVRAM := TotalVRAMMB()
	if totalVRAM > 0 {
		vramMB = totalVRAM * 90 / 100
	}
	return ramMB, vramMB
}
