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

// Draw game here
func (g *Game) Draw() {
	ClearScreen(DarkBlue)
	Pset(4, 0, Red)
	// Print(fmt.Sprintf("%v", Pget(4, 0)))
	Rect(0, 10, 20, 20)
	Rectfill(0, 60, 40, 80)
	Line(0, 8, 40, 8)
}

func main() {
	Run(&Game{})
}
