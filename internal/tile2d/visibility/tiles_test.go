package visibility

import "testing"

func TestVisibleCellBoundsClampToMap(t *testing.T) {
	got := VisibleCellBounds(
		Camera{X: 16, Y: 8, Width: 288, Height: 216},
		MapInfo{Width: 128, Height: 64, TileWidth: 8, TileHeight: 8},
	)
	if got.MinX != 2 || got.MinY != 1 {
		t.Fatalf("unexpected min bounds: %+v", got)
	}
}
