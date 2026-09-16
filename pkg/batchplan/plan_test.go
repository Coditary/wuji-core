package batchplan

import "testing"

func TestResolveExplicitBoth(t *testing.T) {
	got, err := Resolve(Input{Total: 12, Size: 4, Count: 3, SizeSet: true, CountSet: true})
	if err != nil {
		t.Fatal(err)
	}
	if got.Size != 4 || got.Count != 3 {
		t.Fatalf("got %dx%d, want 4x3", got.Size, got.Count)
	}
}

func TestResolveSizeOnly(t *testing.T) {
	got, err := Resolve(Input{Total: 12, Size: 4, SizeSet: true})
	if err != nil {
		t.Fatal(err)
	}
	if got.Size != 4 || got.Count != 3 {
		t.Fatalf("got %dx%d, want 4x3", got.Size, got.Count)
	}
}

func TestResolveCountOnly(t *testing.T) {
	got, err := Resolve(Input{Total: 12, Count: 3, CountSet: true})
	if err != nil {
		t.Fatal(err)
	}
	if got.Size != 4 || got.Count != 3 {
		t.Fatalf("got %dx%d, want 4x3", got.Size, got.Count)
	}
}

func TestResolveIndivisibleSize(t *testing.T) {
	_, err := Resolve(Input{Total: 12, Size: 5, SizeSet: true})
	if err == nil {
		t.Fatal("expected error for non-divisible batch-size")
	}
}

func TestResolveHeuristicVRAM(t *testing.T) {
	got, err := Resolve(Input{
		Total:           12,
		Width:           512,
		Height:          512,
		AvailableVRAMMB: 12000,
	})
	if err != nil {
		t.Fatal(err)
	}
	// ~12GB is a tight fit → prefer lower batch-size (2×6 over 4×3).
	if got.Size != 2 || got.Count != 6 {
		t.Fatalf("got %dx%d, want 2x6 with tight ~12GB VRAM", got.Size, got.Count)
	}
}

func TestResolveHeuristicComfortableVRAM(t *testing.T) {
	got, err := Resolve(Input{
		Total:           12,
		Width:           512,
		Height:          512,
		AvailableVRAMMB: 24000,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Size != 4 || got.Count != 3 {
		t.Fatalf("got %dx%d, want 4x3 with comfortable VRAM", got.Size, got.Count)
	}
}

func TestResolveHeuristicLowVRAM(t *testing.T) {
	got, err := Resolve(Input{Total: 12, AvailableVRAMMB: 4096})
	if err != nil {
		t.Fatal(err)
	}
	if got.Size != 1 || got.Count != 12 {
		t.Fatalf("got %dx%d, want 1x12 with low VRAM", got.Size, got.Count)
	}
}

func TestResolveNoBatchFlag(t *testing.T) {
	got, err := Resolve(Input{Size: 2, Count: 5})
	if err != nil {
		t.Fatal(err)
	}
	if got.Size != 2 || got.Count != 5 {
		t.Fatalf("got %dx%d, want 2x5", got.Size, got.Count)
	}
}

func TestSplitTotal(t *testing.T) {
	size, count := splitTotal(12, 4, false)
	if size != 4 || count != 3 {
		t.Fatalf("splitTotal(12,4,false) = %dx%d, want 4x3", size, count)
	}
	size, count = splitTotal(12, 4, true)
	if size != 2 || count != 6 {
		t.Fatalf("splitTotal(12,4,true) = %dx%d, want 2x6", size, count)
	}
	size, count = splitTotal(7, 4, false)
	if size != 1 || count != 7 {
		t.Fatalf("splitTotal(7,4,false) = %dx%d, want 1x7", size, count)
	}
	size, count = splitTotal(16, 8, false)
	if size != 8 || count != 2 {
		t.Fatalf("splitTotal(16,8,false) = %dx%d, want 8x2", size, count)
	}
}
