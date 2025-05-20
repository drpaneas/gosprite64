package main

import (
	. "github.com/drpaneas/gosprite64"
)

// Game instances to store game state
type Game struct{}

var x, y int

// Init is called once at the start of the game
func (g *Game) Init() {
}

// Update game logic here
func (g *Game) Update() {}

// Draw game here
func (g *Game) Draw() {
	// ClearScreen(DarkBlue)

	if Btn(BtnUp) {
		y -= 1
	}
	if Btn(BtnDown) {
		y += 1
	}
	if Btn(BtnLeft) {
		x -= 1
	}
	if Btn(BtnRight) {
		x += 1
	}
	Rectfill(x, y, x+10, y+10, 7)
}

func main() {
	Run(&Game{})
}
