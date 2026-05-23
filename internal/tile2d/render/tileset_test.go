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

func TestTilesetExtractsDistinctTilesFromMultipleRows(t *testing.T) {
	atlas := image.NewRGBA(image.Rect(0, 0, 64, 64))

	fill := func(x0, y0, x1, y1 int, c color.RGBA) {
		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				atlas.Set(x, y, c)
			}
		}
	}

	fill(0, 0, 32, 32, color.RGBA{R: 0xFF, A: 0xFF})
	fill(32, 0, 64, 32, color.RGBA{G: 0xFF, A: 0xFF})
	fill(0, 32, 32, 64, color.RGBA{B: 0xFF, A: 0xFF})
	fill(32, 32, 64, 64, color.RGBA{R: 0xFF, G: 0xFF, A: 0xFF})

	tileset, err := NewTilesetFromAtlas(atlas, 32, 32)
	if err != nil {
		t.Fatalf("NewTilesetFromAtlas() error = %v", err)
	}

	wantDominant := []struct {
		tileID uint16
		check  func(r, g, b uint32) bool
		label  string
	}{
		{tileID: 1, check: func(r, g, b uint32) bool { return r > g && r > b }, label: "red"},
		{tileID: 2, check: func(r, g, b uint32) bool { return g > r && g > b }, label: "green"},
		{tileID: 3, check: func(r, g, b uint32) bool { return b > r && b > g }, label: "blue"},
		{tileID: 4, check: func(r, g, b uint32) bool { return r > 0 && g > 0 && b == 0 }, label: "yellow"},
	}

	for _, tc := range wantDominant {
		tile := tileset.Tile(tc.tileID)
		if tile == nil {
			t.Fatalf("tile %d missing", tc.tileID)
		}
		r, g, b, _ := tile.At(tile.Bounds().Min.X, tile.Bounds().Min.Y).RGBA()
		if !tc.check(r, g, b) {
			t.Fatalf("tile %d pixel = (%d,%d,%d), want %s-dominant", tc.tileID, r, g, b, tc.label)
		}
	}
}
