package main

import . "github.com/drpaneas/gosprite64"

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
	ClearScreen(DarkBlue)

	// Outline the public 288x216 canvas so the surrounding gutters stay obvious.
	DrawRect(0, 0, logicalWidth-1, logicalHeight-1, Pico8Palette[White])

	Rectfill(0, 0, markerSize, markerSize, Pico8Palette[Red])
	Rectfill(logicalWidth-1-markerSize, 0, logicalWidth-1, markerSize, Pico8Palette[Orange])
	Rectfill(0, logicalHeight-1-markerSize, markerSize, logicalHeight-1, Pico8Palette[Green])
	Rectfill(logicalWidth-1-markerSize, logicalHeight-1-markerSize, logicalWidth-1, logicalHeight-1, Pico8Palette[Blue])

	Rectfill(centerMinX, centerMinY, centerMaxX, centerMaxY, Pico8Palette[Yellow])

	DrawRect(squareLeft, squareTop, squareRight, squareBottom, Pico8Palette[Pink])
	Line(squareLeft, squareTop, squareRight, squareBottom, Pico8Palette[LightGray])
	Line(squareLeft, squareBottom, squareRight, squareTop, Pico8Palette[LightGray])

	Print("288x216", 112, 8, White)
	Print("TL", 6, 6, White)
	Print("TR", logicalWidth-22, 6, White)
	Print("BL", 6, logicalHeight-14, White)
	Print("BR", logicalWidth-22, logicalHeight-14, White)
}

func main() {
	Run(&Game{})
}
