package main

import (
	"fmt"

	"github.com/drpaneas/gosprite64"
)

type Game struct {
	crossX float64
	crossY float64
}

func (g *Game) Init() {
	g.crossX = 144
	g.crossY = 108
}

func (g *Game) Update() {
	sx, sy := gosprite64.StickPosition(0.2)
	g.crossX = 144 + sx*120
	g.crossY = 108 + sy*90
}

func (g *Game) Draw() {
	gosprite64.ClearScreen()

	gosprite64.DrawText("ANALOG STICK DEMO", 72, 4, gosprite64.White)
	gosprite64.DrawText("MOVE THE STICK", 88, 16, gosprite64.LightGray)

	cx := int(g.crossX)
	cy := int(g.crossY)

	gosprite64.DrawRect(24, 28, 264, 208, gosprite64.DarkGray)

	gosprite64.FillRect(cx-8, cy, cx+8, cy+1, gosprite64.Green)
	gosprite64.FillRect(cx, cy-8, cx+1, cy+8, gosprite64.Green)
	gosprite64.FillRect(cx-2, cy-2, cx+2, cy+2, gosprite64.White)

	sx, sy := gosprite64.StickPosition(0.2)
	gosprite64.DrawText(fmt.Sprintf("X: %+.2f", sx), 8, 212, gosprite64.Yellow)
	gosprite64.DrawText(fmt.Sprintf("Y: %+.2f", sy), 160, 212, gosprite64.Yellow)
}

func main() {
	gosprite64.Run(&Game{})
}
