//go:generate sh -c "mkdir -p assets && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/hero.png -out assets/hero.sheet -tile-width 32 -tile-height 32"

package main

const (
	heroWidth      = 64
	heroHeight     = 64
	heroTileWidth  = 32
	heroTileHeight = 32
)

type heroTile struct {
	frame int
	x     int
	y     int
}

func centeredSpritePosition(screenW, screenH, spriteW, spriteH int) (int, int) {
	return (screenW - spriteW) / 2, (screenH - spriteH) / 2
}

func heroCompositeTiles(x, y int) []heroTile {
	return []heroTile{
		{frame: 0, x: x, y: y},
		{frame: 1, x: x + heroTileWidth, y: y},
		{frame: 2, x: x, y: y + heroTileHeight},
		{frame: 3, x: x + heroTileWidth, y: y + heroTileHeight},
	}
}
