package loader

import (
	"image"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
)

func TestOpenBundleLoadsManifestOnly(t *testing.T) {
	l := NewMemoryLoader(fixtures(t))
	b, err := OpenBundle("level1.bundle", l)
	if err != nil {
		t.Fatalf("OpenBundle() error = %v", err)
	}
	if len(b.Entries) == 0 {
		t.Fatal("expected manifest entries")
	}
}

func TestLoadMapAndSheetMaterializeBundleAssets(t *testing.T) {
	l := NewMemoryLoader(fixtures(t))
	b, err := OpenBundle("level1.bundle", l)
	if err != nil {
		t.Fatalf("OpenBundle() error = %v", err)
	}

	mapPath := manifestPath(t, b, format.BundleKindMap)
	sheetPath := manifestPath(t, b, format.BundleKindSheet)

	m, err := LoadMap(mapPath, l)
	if err != nil {
		t.Fatalf("LoadMap() error = %v", err)
	}
	if m.CellBits != 16 {
		t.Fatalf("CellBits = %d, want 16", m.CellBits)
	}

	sheet, err := LoadSheet(sheetPath, l)
	if err != nil {
		t.Fatalf("LoadSheet() error = %v", err)
	}
	if sheet.TileCount == 0 {
		t.Fatal("expected non-empty sheet")
	}
}

func TestValidateSceneAssetsAllowsMultipleSheets(t *testing.T) {
	sheetA, err := format.BuildSheet(image.NewRGBA(image.Rect(0, 0, 8, 8)), 8, 8)
	if err != nil {
		t.Fatalf("BuildSheet() error = %v", err)
	}
	sheetB, err := format.BuildSheet(image.NewRGBA(image.Rect(0, 0, 8, 8)), 8, 8)
	if err != nil {
		t.Fatalf("BuildSheet() error = %v", err)
	}
	parsedSheetA, err := format.ParseSheet(sheetA)
	if err != nil {
		t.Fatalf("ParseSheet() error = %v", err)
	}
	parsedSheetB, err := format.ParseSheet(sheetB)
	if err != nil {
		t.Fatalf("ParseSheet() error = %v", err)
	}
	rawMap, err := format.BuildMap(format.MapConfig{
		Width:       1,
		Height:      1,
		LayerCount:  2,
		CellBits:    16,
		ChunkWidth:  1,
		ChunkHeight: 1,
		Layers: []format.MapLayerConfig{
			{SheetID: 1, Cells: []uint16{1}},
			{SheetID: 2, Cells: []uint16{1}},
		},
	})
	if err != nil {
		t.Fatalf("BuildMap() error = %v", err)
	}
	parsedMap, err := format.ParseMap(rawMap)
	if err != nil {
		t.Fatalf("ParseMap() error = %v", err)
	}

	if err := ValidateSceneAssets(parsedMap, []format.ParsedSheet{parsedSheetA, parsedSheetB}); err != nil {
		t.Fatalf("ValidateSceneAssets() error = %v", err)
	}
}

func TestValidateSceneAssetsRejectsMissingSheet(t *testing.T) {
	rawMap, err := format.BuildMap(format.MapConfig{
		Width:       1,
		Height:      1,
		LayerCount:  1,
		CellBits:    16,
		ChunkWidth:  1,
		ChunkHeight: 1,
		Layers: []format.MapLayerConfig{
			{SheetID: 2, Cells: []uint16{1}},
		},
	})
	if err != nil {
		t.Fatalf("BuildMap() error = %v", err)
	}
	parsedMap, err := format.ParseMap(rawMap)
	if err != nil {
		t.Fatalf("ParseMap() error = %v", err)
	}
	sheet, err := format.ParseSheet(mustBuildSheet(t))
	if err != nil {
		t.Fatalf("ParseSheet() error = %v", err)
	}

	if err := ValidateSceneAssets(parsedMap, []format.ParsedSheet{sheet}); err == nil {
		t.Fatal("expected missing-sheet validation error")
	}
}

func TestValidateSceneAssetsRejectsTileOutOfRange(t *testing.T) {
	rawMap, err := format.BuildMap(format.MapConfig{
		Width:       1,
		Height:      1,
		LayerCount:  1,
		CellBits:    16,
		ChunkWidth:  1,
		ChunkHeight: 1,
		Layers: []format.MapLayerConfig{
			{SheetID: 1, Cells: []uint16{2}},
		},
	})
	if err != nil {
		t.Fatalf("BuildMap() error = %v", err)
	}
	parsedMap, err := format.ParseMap(rawMap)
	if err != nil {
		t.Fatalf("ParseMap() error = %v", err)
	}
	sheet, err := format.ParseSheet(mustBuildSheet(t))
	if err != nil {
		t.Fatalf("ParseSheet() error = %v", err)
	}

	if err := ValidateSceneAssets(parsedMap, []format.ParsedSheet{sheet}); err == nil {
		t.Fatal("expected tile-range validation error")
	}
}

func TestValidateSceneAssetsAcceptsNon8x8TileSize(t *testing.T) {
	rawMap, err := format.BuildMap(format.MapConfig{
		Width:       1,
		Height:      1,
		LayerCount:  1,
		CellBits:    16,
		ChunkWidth:  1,
		ChunkHeight: 1,
		Layers: []format.MapLayerConfig{
			{SheetID: 1, Cells: []uint16{1}},
		},
	})
	if err != nil {
		t.Fatalf("BuildMap() error = %v", err)
	}
	parsedMap, err := format.ParseMap(rawMap)
	if err != nil {
		t.Fatalf("ParseMap() error = %v", err)
	}

	rawSheet, err := format.BuildSheet(image.NewRGBA(image.Rect(0, 0, 16, 8)), 16, 8)
	if err != nil {
		t.Fatalf("BuildSheet() error = %v", err)
	}
	sheet, err := format.ParseSheet(rawSheet)
	if err != nil {
		t.Fatalf("ParseSheet() error = %v", err)
	}

	err = ValidateSceneAssets(parsedMap, []format.ParsedSheet{sheet})
	if err != nil {
		t.Fatalf("should accept non-8x8 tile size: %v", err)
	}
}

func TestOpenBundleRejectsInvalidMagic(t *testing.T) {
	l := NewMemoryLoader(map[string][]byte{"bad.bundle": []byte("not a bundle")})
	_, err := OpenBundle("bad.bundle", l)
	if err == nil {
		t.Fatal("expected error for invalid bundle magic")
	}
}

func TestOpenBundleRejectsMissingFile(t *testing.T) {
	l := NewMemoryLoader(map[string][]byte{})
	_, err := OpenBundle("missing.bundle", l)
	if err == nil {
		t.Fatal("expected error for missing bundle file")
	}
}

func TestLoadSheetRejectsMalformedData(t *testing.T) {
	l := NewMemoryLoader(map[string][]byte{"bad.sheet": []byte("bad!")})
	_, err := LoadSheet("bad.sheet", l)
	if err == nil {
		t.Fatal("expected error for malformed sheet")
	}
}

func TestLoadMapRejectsMalformedData(t *testing.T) {
	l := NewMemoryLoader(map[string][]byte{"bad.map": []byte("bad!")})
	_, err := LoadMap("bad.map", l)
	if err == nil {
		t.Fatal("expected error for malformed map")
	}
}

func TestLoadAnimRejectsMalformedData(t *testing.T) {
	l := NewMemoryLoader(map[string][]byte{"bad.anim": []byte("bad!")})
	_, err := LoadAnim("bad.anim", l)
	if err == nil {
		t.Fatal("expected error for malformed anim")
	}
}

func TestBundleManifestIsNotMonolithicBlob(t *testing.T) {
	l := NewMemoryLoader(fixtures(t))
	b, err := OpenBundle("level1.bundle", l)
	if err != nil {
		t.Fatalf("OpenBundle() error = %v", err)
	}
	for _, entry := range b.Entries {
		if entry.Path == "" {
			t.Fatalf("bundle entry %q has empty path", entry.Name)
		}
		_, err := l.ReadAsset(entry.Path)
		if err != nil {
			t.Fatalf("bundle entry %q path %q not separately loadable: %v", entry.Name, entry.Path, err)
		}
	}
}

func mustBuildSheet(t *testing.T) []byte {
	t.Helper()
	raw, err := format.BuildSheet(image.NewRGBA(image.Rect(0, 0, 8, 8)), 8, 8)
	if err != nil {
		t.Fatalf("BuildSheet() error = %v", err)
	}
	return raw
}

func fixtures(t *testing.T) map[string][]byte {
	t.Helper()

	bundle := mustReadFixture(t, "minimal.bundle")
	sheet := mustReadFixture(t, "minimal.sheet")
	m := mustReadFixture(t, "minimal.map")
	anim := mustReadFixture(t, "minimal.anim")

	return map[string][]byte{
		"level1.bundle":  bundle,
		"minimal.bundle": bundle,
		"minimal.sheet":  sheet,
		"minimal.map":    m,
		"minimal.anim":   anim,
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

func manifestPath(t *testing.T, b format.ParsedBundle, kind uint8) string {
	t.Helper()

	for _, entry := range b.Entries {
		if entry.Kind == kind {
			return entry.Path
		}
	}

	t.Fatalf("missing manifest entry kind %d", kind)
	return ""
}
