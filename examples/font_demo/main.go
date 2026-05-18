//go:generate sh -c "cd assets-src && go run gen_assets.go"
//go:generate sh -c "mkdir -p assets && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/font.png -out assets/font.sheet -tile-width 8 -tile-height 8"

package main

import (
	"github.com/drpaneas/gosprite64"
)

type Game struct {
	font *gosprite64.Font
}

func (g *Game) Init() {
	sheet, err := gosprite64.LoadSpriteSheet("assets/font.sheet")
	if err != nil {
		panic(err)
	}

	glyphs := make(map[rune]gosprite64.Glyph)
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i, r := range charset {
		glyphs[r] = gosprite64.Glyph{
			Frame:   i,
			Width:   8,
			Advance: 9,
		}
	}
	glyphs[' '] = gosprite64.Glyph{Frame: 0, Width: 0, Advance: 6}

	g.font = gosprite64.NewFont(sheet, glyphs, 10)
	g.font.Fallback = '?'
}

func (g *Game) Update() {}

func (g *Game) Draw() {
	gosprite64.ClearScreenWith(gosprite64.DarkBlue)

	gosprite64.DrawText("CUSTOM FONT DEMO", 72, 10, gosprite64.White)

	g.font.DrawTextEx("HELLO WORLD", 80, 50, gosprite64.AlignLeft)
	g.font.DrawTextEx("ABCDEFGHIJKLM", 40, 80, gosprite64.AlignLeft)
	g.font.DrawTextEx("NOPQRSTUVWXYZ", 40, 100, gosprite64.AlignLeft)
	g.font.DrawTextEx("0123456789", 80, 130, gosprite64.AlignCenter)

	gosprite64.DrawText("BUILT IN FONT FOR COMPARISON", 32, 180, gosprite64.LightGray)
}

func main() {
	gosprite64.Run(&Game{})
}
