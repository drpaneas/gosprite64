//go:build n64

package main

import "github.com/drpaneas/gosprite64"

const (
	screenWidth  = 288
	screenHeight = 216
)

type Game struct {
	hero *gosprite64.SpriteSheet
	x    int
	y    int
}

func (g *Game) Init() {
	hero, err := gosprite64.LoadSpriteSheet("assets/hero.sheet")
	if err != nil {
		panic(err)
	}

	g.hero = hero
	g.x, g.y = centeredSpritePosition(screenWidth, screenHeight, heroWidth, heroHeight)
}

func (*Game) Update() {}

func (g *Game) Draw() {
	gosprite64.ClearScreenWith(gosprite64.DarkBlue)
	gosprite64.DrawText("Aseprite sprite loaded", 56, 24, gosprite64.White)
	for _, tile := range heroCompositeTiles(g.x, g.y) {
		gosprite64.DrawSpriteWithOptions(g.hero, tile.frame, float32(tile.x), float32(tile.y), gosprite64.DrawSpriteOptions{
			Blend: gosprite64.BlendMasked,
		})
	}
}

func main() {
	gosprite64.Run(&Game{})
}
