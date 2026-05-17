package loader

import (
	"testing"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
)

func TestValidateAccepts16x16(t *testing.T) {
	m := format.ParsedMap{Layers: []format.ParsedMapLayer{{Cells: []uint16{1}}}}
	sheets := []format.ParsedSheet{{TileWidth: 16, TileHeight: 16, TileCount: 4}}
	if err := ValidateSceneAssets(m, sheets); err != nil {
		t.Fatalf("should accept 16x16: %v", err)
	}
}

func TestValidateRejectsMixedSizes(t *testing.T) {
	m := format.ParsedMap{Layers: []format.ParsedMapLayer{{Cells: []uint16{1}}}}
	sheets := []format.ParsedSheet{
		{TileWidth: 8, TileHeight: 8, TileCount: 4},
		{TileWidth: 16, TileHeight: 16, TileCount: 4},
	}
	if err := ValidateSceneAssets(m, sheets); err == nil {
		t.Fatal("should reject mixed tile sizes")
	}
}

func TestValidateRejectsZeroSize(t *testing.T) {
	m := format.ParsedMap{Layers: []format.ParsedMapLayer{{Cells: []uint16{1}}}}
	sheets := []format.ParsedSheet{{TileWidth: 0, TileHeight: 0, TileCount: 4}}
	if err := ValidateSceneAssets(m, sheets); err == nil {
		t.Fatal("should reject zero tile size")
	}
}

func TestValidateAccepts8x8(t *testing.T) {
	m := format.ParsedMap{Layers: []format.ParsedMapLayer{{Cells: []uint16{1}}}}
	sheets := []format.ParsedSheet{{TileWidth: 8, TileHeight: 8, TileCount: 4}}
	if err := ValidateSceneAssets(m, sheets); err != nil {
		t.Fatalf("should still accept 8x8: %v", err)
	}
}
