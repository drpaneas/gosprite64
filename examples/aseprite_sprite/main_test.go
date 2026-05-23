package main

import (
	"os"
	"strings"
	"testing"
)

func TestCenteredSpritePosition(t *testing.T) {
	x, y := centeredSpritePosition(288, 216, heroWidth, heroHeight)

	if x != 112 || y != 76 {
		t.Fatalf("centeredSpritePosition(288, 216, %d, %d) = (%d, %d), want (112, 76)", heroWidth, heroHeight, x, y)
	}
}

func TestHeroCompositeTiles(t *testing.T) {
	tiles := heroCompositeTiles(112, 76)

	want := []heroTile{
		{frame: 0, x: 112, y: 76},
		{frame: 1, x: 144, y: 76},
		{frame: 2, x: 112, y: 108},
		{frame: 3, x: 144, y: 108},
	}

	if len(tiles) != len(want) {
		t.Fatalf("heroCompositeTiles tile count = %d, want %d", len(tiles), len(want))
	}

	for i := range want {
		if tiles[i] != want[i] {
			t.Fatalf("heroCompositeTiles()[%d] = %+v, want %+v", i, tiles[i], want[i])
		}
	}
}

func TestGenerateDirectiveUses32x32Tiles(t *testing.T) {
	source, err := os.ReadFile("centered_position.go")
	if err != nil {
		t.Fatalf("read centered_position.go: %v", err)
	}

	if !strings.Contains(string(source), "-tile-width 32 -tile-height 32") {
		t.Fatalf("centered_position.go go:generate directive must use 32x32 tiles")
	}
}
