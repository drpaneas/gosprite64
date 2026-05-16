package render

import (
	"image"
	"image/color"
	"testing"
)

func TestTilesetExtractsDistinctTilesFromAtlas(t *testing.T) {
	atlas := image.NewRGBA(image.Rect(0, 0, 16, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			atlas.Set(x, y, color.RGBA{R: 0xFF, A: 0xFF})
		}
		for x := 8; x < 16; x++ {
			atlas.Set(x, y, color.RGBA{G: 0xFF, A: 0xFF})
		}
	}

	tileset, err := NewTilesetFromAtlas(atlas, 8, 8)
	if err != nil {
		t.Fatalf("NewTilesetFromAtlas() error = %v", err)
	}

	tile1 := tileset.Tile(1)
	tile2 := tileset.Tile(2)
	if tile1 == nil || tile2 == nil {
		t.Fatal("expected both tiles to exist")
	}

	r1, g1, _, _ := tile1.At(tile1.Bounds().Min.X, tile1.Bounds().Min.Y).RGBA()
	r2, g2, _, _ := tile2.At(tile2.Bounds().Min.X, tile2.Bounds().Min.Y).RGBA()
	if r1 <= g1 {
		t.Fatalf("tile1 pixel = (%d,%d), want red-dominant", r1, g1)
	}
	if g2 <= r2 {
		t.Fatalf("tile2 pixel = (%d,%d), want green-dominant", r2, g2)
	}
}

func TestTilesetReturnsSameTextureInstance(t *testing.T) {
	atlas := image.NewRGBA(image.Rect(0, 0, 8, 8))
	tileset, err := NewTilesetFromAtlas(atlas, 8, 8)
	if err != nil {
		t.Fatalf("NewTilesetFromAtlas() error = %v", err)
	}

	if tileset.Tile(1) != tileset.Tile(1) {
		t.Fatal("expected cached tile texture instance")
	}
}
