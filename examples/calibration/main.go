package main

import "github.com/drpaneas/gosprite64"

const (
	logicalWidth  = 288
	logicalHeight = 216
	markerSize    = 3

	centerMinX = 143
	centerMinY = 107
	centerMaxX = 144
	centerMaxY = 108

	squareLeft   = 124
	squareTop    = 88
	squareRight  = 163
	squareBottom = 127
)

type Game struct{}

func (g *Game) Init() {}

func (g *Game) Update() {}

func (g *Game) Draw() {
	gosprite64.ClearScreenWith(gosprite64.DarkBlue)

	// Outline the public 288x216 canvas so the surrounding gutters stay obvious.
	gosprite64.DrawRect(0, 0, logicalWidth-1, logicalHeight-1, gosprite64.White)

	gosprite64.FillRect(0, 0, markerSize, markerSize, gosprite64.Red)
	gosprite64.FillRect(logicalWidth-1-markerSize, 0, logicalWidth-1, markerSize, gosprite64.Orange)
	gosprite64.FillRect(0, logicalHeight-1-markerSize, markerSize, logicalHeight-1, gosprite64.Green)
	gosprite64.FillRect(logicalWidth-1-markerSize, logicalHeight-1-markerSize, logicalWidth-1, logicalHeight-1, gosprite64.Blue)

	gosprite64.FillRect(centerMinX, centerMinY, centerMaxX, centerMaxY, gosprite64.Yellow)

	gosprite64.DrawRect(squareLeft, squareTop, squareRight, squareBottom, gosprite64.Pink)
	gosprite64.DrawLine(squareLeft, squareTop, squareRight, squareBottom, gosprite64.LightGray)
	gosprite64.DrawLine(squareLeft, squareBottom, squareRight, squareTop, gosprite64.LightGray)

	gosprite64.DrawText("288x216", 112, 8, gosprite64.White)
	gosprite64.DrawText("TL", 6, 6, gosprite64.White)
	gosprite64.DrawText("TR", logicalWidth-22, 6, gosprite64.White)
	gosprite64.DrawText("BL", 6, logicalHeight-14, gosprite64.White)
	gosprite64.DrawText("BR", logicalWidth-22, logicalHeight-14, gosprite64.White)
}

func main() {
	gosprite64.Run(&Game{})
}
