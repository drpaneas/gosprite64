package main

import (
	"image/color"

	gospr64 "github.com/drpaneas/gosprite64"
)

var Azure = color.RGBA{0xf0, 0xff, 0xff, 0xff} // rgb(240, 255, 255)

type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *gospr64.Screen) {
	screen.Clear(Azure)
}

func main() {
	// Initialize the game
	game := &Game{}

	// Run the game
	if err := gospr64.Run(game); err != nil {
		panic(err)
	}

}
