package main

import (
	. "github.com/drpaneas/gosprite64"
)

// Game instances to store game state
type Game struct{}

// var x, y int

// Init is called once at the start of the game
func (g *Game) Init() {
}

// Update game logic here
func (g *Game) Update() {}

// Draw game here
func (g *Game) Draw() {

	// if Btn(BtnUp) {
	// 	y -= 1
	// }
	// if Btn(BtnDown) {
	// 	y += 1
	// }
	// if Btn(BtnLeft) {
	// 	x -= 1
	// }
	// if Btn(BtnRight) {
	// 	x += 1
	// }
	// Rectfill(x, y, x+10, y+10, 7)
	// Rectfill(1, 1, ScreenWidth-1, ScreenHeight-1, 7)
	// Draw something at virtual (50,50)
	// ClearScreen(Red)
	// Rectfill(x, y, x+20, y+20, color.White)
}

func main() {
	Run(&Game{})
}
