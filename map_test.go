package gosprite64

import "testing"

func TestMapTileWidthDefault(t *testing.T) {
	m := &Map{}
	if m.TileWidth() != 8 {
		t.Fatalf("expected default 8, got %d", m.TileWidth())
	}
}

func TestMapTileWidthCustom(t *testing.T) {
	m := &Map{tileW: 16, tileH: 16}
	if m.TileWidth() != 16 {
		t.Fatalf("expected 16, got %d", m.TileWidth())
	}
	if m.TileHeight() != 16 {
		t.Fatalf("expected 16, got %d", m.TileHeight())
	}
}
