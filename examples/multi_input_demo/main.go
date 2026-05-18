package main

import (
	"fmt"
	"image/color"

	"github.com/drpaneas/gosprite64"
)

type Game struct {
	posX [4]int
	posY [4]int
}

var playerColors = [4]color.Color{
	gosprite64.Red,
	gosprite64.Blue,
	gosprite64.Green,
	gosprite64.Yellow,
}

func (g *Game) Init() {
	g.posX = [4]int{60, 200, 60, 200}
	g.posY = [4]int{60, 60, 150, 150}
}

func (g *Game) Update() {
	for port := 0; port < gosprite64.MaxControllers; port++ {
		if !gosprite64.IsControllerConnected(port) {
			continue
		}
		if gosprite64.PlayerButtonDown(port, gosprite64.ButtonDPadUp) {
			g.posY[port] -= 2
		}
		if gosprite64.PlayerButtonDown(port, gosprite64.ButtonDPadDown) {
			g.posY[port] += 2
		}
		if gosprite64.PlayerButtonDown(port, gosprite64.ButtonDPadLeft) {
			g.posX[port] -= 2
		}
		if gosprite64.PlayerButtonDown(port, gosprite64.ButtonDPadRight) {
			g.posX[port] += 2
		}

		if g.posX[port] < 8 {
			g.posX[port] = 8
		}
		if g.posX[port] > 272 {
			g.posX[port] = 272
		}
		if g.posY[port] < 20 {
			g.posY[port] = 20
		}
		if g.posY[port] > 200 {
			g.posY[port] = 200
		}
	}
}

func (g *Game) Draw() {
	gosprite64.ClearScreen()

	gosprite64.DrawText("MULTI-CONTROLLER DEMO", 56, 4, gosprite64.White)

	for port := 0; port < gosprite64.MaxControllers; port++ {
		x := g.posX[port]
		y := g.posY[port]
		c := playerColors[port]

		if gosprite64.IsControllerConnected(port) {
			gosprite64.FillRect(x-6, y-6, x+6, y+6, c)
			gosprite64.DrawText(
				fmt.Sprintf("P%d", port+1),
				x-4, y-4, gosprite64.Black,
			)
		} else {
			gosprite64.DrawRect(x-6, y-6, x+6, y+6, gosprite64.DarkGray)
			gosprite64.DrawText(
				fmt.Sprintf("P%d", port+1),
				x-4, y-4, gosprite64.DarkGray,
			)
		}
	}

	gosprite64.DrawText(
		fmt.Sprintf("CONNECTED: %d", gosprite64.ConnectedControllers()),
		8, 210, gosprite64.LightGray,
	)
}

func main() {
	gosprite64.Run(&Game{})
}
