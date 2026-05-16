package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
)

func TestMk2DMapProduces16BitCellMap(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "level.json")
	out := filepath.Join(dir, "level.map")

	input := []byte(`{"width":32,"height":18,"layer_count":1,"cell_bits":16,"chunk_width":8,"chunk_height":8}`)
	if err := os.WriteFile(in, input, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := run([]string{"-in", in, "-out", out}); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	parsed, err := format.ParseMap(raw)
	if err != nil {
		t.Fatalf("ParseMap() error = %v", err)
	}
	if parsed.CellBits != 16 {
		t.Fatalf("CellBits = %d, want 16", parsed.CellBits)
	}
}
