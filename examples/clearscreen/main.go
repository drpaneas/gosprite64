package main

import (
	"github.com/drpaneas/gosprite64"
)

// Game instances to store game state
type Game struct{}

// Init is called once at the start of the game
func (g *Game) Init() {}

// Update game logic here
func (g *Game) Update() {}

// Draw renders a solid blue screen.
func (g *Game) Draw() {
	gosprite64.ClearScreenWith(gosprite64.Blue)
}

func main() {
	gosprite64.Run(&Game{})
}
