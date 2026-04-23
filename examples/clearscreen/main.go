package main

import (
	. "github.com/drpaneas/gosprite64"
)

// Game instances to store game state
type Game struct{}

// Init is called once at the start of the game
func (g *Game) Init() {}

// Update game logic here
func (g *Game) Update() {}

// ClearScreen is the simplest render sanity check. Other drawing APIs use the
// fixed 288x216 logical canvas described in the docs and calibration example.
func (g *Game) Draw() {
	ClearScreen(Red)
}

func main() {
	Run(&Game{})
}
