package sysmem_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/sysmem"
)

func TestSuggestedBudgetWithinSystem(t *testing.T) {
	ram, vram := sysmem.SuggestedBudget()
	totalRAM := sysmem.TotalRAMMB()
	if totalRAM > 0 && (ram <= 0 || ram > totalRAM) {
		t.Fatalf("suggested RAM %d out of range for total %d", ram, totalRAM)
	}
	totalVRAM := sysmem.TotalVRAMMB()
	if totalVRAM > 0 && (vram <= 0 || vram > totalVRAM) {
		t.Fatalf("suggested VRAM %d out of range for total %d", vram, totalVRAM)
	}
}
