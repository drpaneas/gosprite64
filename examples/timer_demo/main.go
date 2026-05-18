package main

import (
	"fmt"
	"image/color"

	"github.com/drpaneas/gosprite64"
)

type Game struct {
	colorTimer *gosprite64.RepeatingTimer
	shotTimer  *gosprite64.Timer
	colorIdx   int
	message    string
}

var palette = []color.Color{
	gosprite64.Red,
	gosprite64.Orange,
	gosprite64.Yellow,
	gosprite64.Green,
	gosprite64.Blue,
	gosprite64.Indigo,
}

func (g *Game) Init() {
	g.colorTimer = gosprite64.NewRepeatingTimer(60)
	g.shotTimer = gosprite64.NewTimer(180)
	g.message = "WAITING..."
}

func (g *Game) Update() {
	if g.colorTimer.Tick() {
		g.colorIdx = g.colorTimer.Count() % len(palette)
	}
	if g.shotTimer.Tick() {
		g.message = "TIMER FIRED!"
	}
}

func (g *Game) Draw() {
	gosprite64.ClearScreen()

	gosprite64.DrawText("TIMER DEMO", 100, 4, gosprite64.White)

	gosprite64.DrawText("REPEATING TIMER (60 FRAMES):", 20, 30, gosprite64.LightGray)
	gosprite64.FillRect(20, 44, 268, 84, palette[g.colorIdx])
	gosprite64.DrawText(
		fmt.Sprintf("CYCLE: %d", g.colorTimer.Count()),
		100, 58, gosprite64.Black,
	)

	gosprite64.DrawText("ONE-SHOT TIMER (180 FRAMES):", 20, 100, gosprite64.LightGray)
	gosprite64.DrawText(g.message, 20, 118, gosprite64.Yellow)

	progress := g.shotTimer.Progress()
	barW := int(progress * 200)
	gosprite64.FillRect(20, 134, 20+barW, 146, gosprite64.Green)
	gosprite64.DrawRect(20, 134, 220, 146, gosprite64.DarkGray)
	gosprite64.DrawText(
		fmt.Sprintf("PROGRESS: %.0f%%", progress*100),
		20, 154, gosprite64.White,
	)

	gosprite64.DrawText(
		fmt.Sprintf("ELAPSED: %d  REMAINING: %d",
			g.shotTimer.Elapsed(), g.shotTimer.Remaining()),
		20, 180, gosprite64.LightGray,
	)
}

func main() {
	gosprite64.Run(&Game{})
}
