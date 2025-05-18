package main

import (
	"image/color"
	"github.com/drpaneas/gosprite64"
)

var Azure = color.RGBA{0xf0, 0xff, 0xff, 0xff} // rgb(240, 255, 255)

type Game struct{}

// Init is called once at the start of the game
func (g *Game) Init() {
	// Initialize game state, load resources, etc.
}

func (g *Game) Update() {
	// Update game logic here
}

func (g *Game) Draw() {
	gosprite64.Clear(Azure)
}

func main() {
	// Run the game
	gosprite64.Run(&Game{})
}
