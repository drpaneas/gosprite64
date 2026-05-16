package format

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestParseSheetRejectsBadMagic(t *testing.T) {
	_, err := ParseSheet([]byte("bad!"))
	if err == nil {
		t.Fatal("expected bad magic error")
	}
}

func TestParseMapSupports16BitCells(t *testing.T) {
	raw := mustReadFixture(t, "minimal.map")
	m, err := ParseMap(raw)
	if err != nil {
		t.Fatalf("ParseMap() error = %v", err)
	}
	if m.CellBits != 16 {
		t.Fatalf("CellBits = %d, want 16", m.CellBits)
	}
}

func TestBuildAndParseMapPreservesLayerCells(t *testing.T) {
	raw, err := BuildMap(MapConfig{
		Width:       2,
		Height:      2,
		LayerCount:  1,
		CellBits:    16,
		ChunkWidth:  2,
		ChunkHeight: 2,
		Layers: []MapLayerConfig{
			{Cells: []uint16{1, 2, 3, 4}},
		},
	})
	if err != nil {
		t.Fatalf("BuildMap() error = %v", err)
	}

	parsed, err := ParseMap(raw)
	if err != nil {
		t.Fatalf("ParseMap() error = %v", err)
	}
	if len(parsed.Layers) != 1 {
		t.Fatalf("len(Layers) = %d, want 1", len(parsed.Layers))
	}
	if got := parsed.Layers[0].Cells; len(got) != 4 || got[0] != 1 || got[3] != 4 {
		t.Fatalf("Cells = %v, want [1 2 3 4]", got)
	}
}

func TestBuildAndParseMapPreservesLayerSheetIDs(t *testing.T) {
	raw, err := BuildMap(MapConfig{
		Width:       1,
		Height:      1,
		LayerCount:  2,
		CellBits:    16,
		ChunkWidth:  1,
		ChunkHeight: 1,
		Layers: []MapLayerConfig{
			{SheetID: 1, Cells: []uint16{1}},
			{SheetID: 2, Cells: []uint16{1}},
		},
	})
	if err != nil {
		t.Fatalf("BuildMap() error = %v", err)
	}

	parsed, err := ParseMap(raw)
	if err != nil {
		t.Fatalf("ParseMap() error = %v", err)
	}
	if len(parsed.Layers) != 2 {
		t.Fatalf("len(Layers) = %d, want 2", len(parsed.Layers))
	}
	if parsed.Layers[0].SheetID != 1 || parsed.Layers[1].SheetID != 2 {
		t.Fatalf("sheet IDs = [%d %d], want [1 2]", parsed.Layers[0].SheetID, parsed.Layers[1].SheetID)
	}
}

func TestBuildAndParseSheetPreservesPixelData(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 16, 8))
	img.Set(0, 0, color.NRGBA{R: 255, A: 255})
	img.Set(8, 0, color.NRGBA{G: 255, A: 255})

	raw, err := BuildSheet(img, 8, 8)
	if err != nil {
		t.Fatalf("BuildSheet() error = %v", err)
	}

	sheet, err := ParseSheet(raw)
	if err != nil {
		t.Fatalf("ParseSheet() error = %v", err)
	}
	if sheet.AtlasWidth != 16 || sheet.AtlasHeight != 8 {
		t.Fatalf("atlas = %dx%d, want 16x8", sheet.AtlasWidth, sheet.AtlasHeight)
	}
	if len(sheet.Pixels) != 16*8*4 {
		t.Fatalf("len(Pixels) = %d, want %d", len(sheet.Pixels), 16*8*4)
	}
	if got := sheet.Pixels[:4]; got[0] != 255 || got[1] != 0 || got[2] != 0 || got[3] != 255 {
		t.Fatalf("first pixel = %v, want red", got)
	}
}

func TestBuildSheetRejectsNonDivisibleTileSize(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 10, 10))
	_, err := BuildSheet(img, 3, 3)
	if err == nil {
		t.Fatal("expected error for non-divisible tile size")
	}
}

func TestBuildSheetRejectsZeroTileSize(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	_, err := BuildSheet(img, 0, 8)
	if err == nil {
		t.Fatal("expected error for zero tile width")
	}
	_, err = BuildSheet(img, 8, 0)
	if err == nil {
		t.Fatal("expected error for zero tile height")
	}
}

func TestBuildSheetRejectsNilImage(t *testing.T) {
	_, err := BuildSheet(nil, 8, 8)
	if err == nil {
		t.Fatal("expected error for nil image")
	}
}

func TestBuildMapRejectsZeroDimensions(t *testing.T) {
	_, err := BuildMap(MapConfig{Width: 0, Height: 1, LayerCount: 1, CellBits: 16, ChunkWidth: 1, ChunkHeight: 1})
	if err == nil {
		t.Fatal("expected error for zero width")
	}
	_, err = BuildMap(MapConfig{Width: 1, Height: 0, LayerCount: 1, CellBits: 16, ChunkWidth: 1, ChunkHeight: 1})
	if err == nil {
		t.Fatal("expected error for zero height")
	}
}

func TestBuildMapRejectsUnsupportedCellBits(t *testing.T) {
	_, err := BuildMap(MapConfig{Width: 1, Height: 1, LayerCount: 1, CellBits: 32, ChunkWidth: 1, ChunkHeight: 1})
	if err == nil {
		t.Fatal("expected error for unsupported cell bits")
	}
}

func TestBuildMapRejectsLayerCountMismatch(t *testing.T) {
	_, err := BuildMap(MapConfig{
		Width: 1, Height: 1, LayerCount: 2, CellBits: 16, ChunkWidth: 1, ChunkHeight: 1,
		Layers: []MapLayerConfig{{Cells: []uint16{1}}},
	})
	if err == nil {
		t.Fatal("expected error for layer count mismatch")
	}
}

func TestBuildAnimRejectsClipNameTooLong(t *testing.T) {
	longName := string(make([]byte, 256))
	_, err := BuildAnim(AnimConfig{Clips: []AnimClipConfig{{Name: longName, FPS: 12, Frames: []uint16{0}}}})
	if err == nil {
		t.Fatal("expected error for clip name too long")
	}
}

func TestParseSheetRejectsTruncatedPayload(t *testing.T) {
	raw := encodeAsset("SHT2", []byte{1, 2})
	_, err := ParseSheet(raw)
	if err == nil {
		t.Fatal("expected error for truncated sheet payload")
	}
}

func TestParseMapRejectsTruncatedPayload(t *testing.T) {
	raw := encodeAsset("MAP2", []byte{1, 2})
	_, err := ParseMap(raw)
	if err == nil {
		t.Fatal("expected error for truncated map payload")
	}
}

func TestParseAnimRejectsTruncatedPayload(t *testing.T) {
	raw := encodeAsset("ANM2", []byte{})
	_, err := ParseAnim(raw)
	if err == nil {
		t.Fatal("expected error for truncated anim payload")
	}
}

func TestParseBundleRejectsTruncatedPayload(t *testing.T) {
	raw := encodeAsset("BND2", []byte{})
	_, err := ParseBundle(raw)
	if err == nil {
		t.Fatal("expected error for truncated bundle payload")
	}
}

func mustReadFixture(t *testing.T, name string) []byte {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	path := filepath.Join(filepath.Dir(file), "..", "testdata", name)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return raw
}
