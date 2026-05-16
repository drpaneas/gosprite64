package stats

import (
	"image"
	"testing"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
	tilerender "github.com/drpaneas/gosprite64/internal/tile2d/render"
)

func TestFromSceneAssetsSummarizesSceneMemoryAndUploads(t *testing.T) {
	rawSheet, err := format.BuildSheet(image.NewRGBA(image.Rect(0, 0, 16, 8)), 8, 8)
	if err != nil {
		t.Fatalf("BuildSheet() error = %v", err)
	}
	sheet, err := format.ParseSheet(rawSheet)
	if err != nil {
		t.Fatalf("ParseSheet() error = %v", err)
	}

	rawMap, err := format.BuildMap(format.MapConfig{
		Width:       2,
		Height:      1,
		LayerCount:  2,
		CellBits:    16,
		ChunkWidth:  1,
		ChunkHeight: 1,
		Layers: []format.MapLayerConfig{
			{SheetID: 1, Cells: []uint16{1, 0}},
			{SheetID: 1, Cells: []uint16{0, 2}},
		},
	})
	if err != nil {
		t.Fatalf("BuildMap() error = %v", err)
	}
	m, err := format.ParseMap(rawMap)
	if err != nil {
		t.Fatalf("ParseMap() error = %v", err)
	}

	got := FromSceneAssets(m, []format.ParsedSheet{sheet}, tilerender.DrawStats{Uploads: 3})
	if got.SheetRAMBytes != len(sheet.Pixels) {
		t.Fatalf("SheetRAMBytes = %d, want %d", got.SheetRAMBytes, len(sheet.Pixels))
	}
	if got.MapRAMBytes != 8 {
		t.Fatalf("MapRAMBytes = %d, want 8", got.MapRAMBytes)
	}
	if got.VisibleTiles != 0 {
		t.Fatalf("VisibleTiles = %d, want 0", got.VisibleTiles)
	}
	if got.SheetCount != 1 {
		t.Fatalf("SheetCount = %d, want 1", got.SheetCount)
	}
	if got.LayerCount != 2 {
		t.Fatalf("LayerCount = %d, want 2", got.LayerCount)
	}
	if got.UploadCount != 3 {
		t.Fatalf("UploadCount = %d, want 3", got.UploadCount)
	}
}
