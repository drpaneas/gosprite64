//go:build n64

package main

import "github.com/drpaneas/gosprite64"

const heroSheetPath = "assets/hero.sheet"

func loadHeroSpriteSheet() (*gosprite64.SpriteSheet, error) {
	return gosprite64.LoadSpriteSheet(heroSheetPath)
}

type Game struct {
	hero *gosprite64.SpriteSheet
	x    float32
	y    float32
}

func (g *Game) Init() {
	sheet, err := loadHeroSpriteSheet()
	if err != nil {
		panic(err)
	}
	g.hero = sheet
	g.x, g.y = centeredHeroPosition(sheet.FrameWidth(), sheet.FrameHeight())
}

func (g *Game) Update() {}

func (g *Game) Draw() {
	gosprite64.ClearScreenWith(gosprite64.DarkBlue)
	gosprite64.DrawSpriteWithOptions(g.hero, 0, g.x, g.y, gosprite64.DrawSpriteOptions{
		Blend: gosprite64.BlendMasked,
	})
}

func main() {
	gosprite64.Run(&Game{})
}
