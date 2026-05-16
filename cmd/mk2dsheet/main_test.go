package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
)

func TestMk2DSheetProducesStableHeader(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "tiles.png")
	out := filepath.Join(dir, "tiles.sheet")

	img := image.NewNRGBA(image.Rect(0, 0, 16, 8))
	img.Set(0, 0, color.NRGBA{R: 255, A: 255})
	img.Set(8, 0, color.NRGBA{G: 255, A: 255})

	f, err := os.Create(in)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	if err := run([]string{"-in", in, "-out", out, "-tile-width", "8", "-tile-height", "8"}); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if got := string(raw[:4]); got != "SHT2" {
		t.Fatalf("magic = %q, want SHT2", got)
	}

	sheet, err := format.ParseSheet(raw)
	if err != nil {
		t.Fatalf("ParseSheet() error = %v", err)
	}
	if sheet.TileCount != 2 {
		t.Fatalf("TileCount = %d, want 2", sheet.TileCount)
	}
}
