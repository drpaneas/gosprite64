package main

import (
	"fmt"

	"github.com/drpaneas/gosprite64"
)

type Game struct {
	menu     *gosprite64.Menu
	selected string
}

func (g *Game) Init() {
	g.menu = gosprite64.NewMenu([]gosprite64.MenuItem{
		{Label: "START GAME", OnConfirm: func() { g.selected = "START GAME" }},
		{Label: "OPTIONS", OnConfirm: func() { g.selected = "OPTIONS" }},
		{Label: "QUIT", OnConfirm: func() { g.selected = "QUIT" }},
	})
	g.menu.X = 100
	g.menu.Y = 80
	g.menu.Wrap = true
	g.selected = ""
}

func (g *Game) Update() {
	g.menu.HandleInput()

	if gosprite64.IsButtonJustPressed(gosprite64.ButtonB) {
		g.selected = ""
	}
}

func (g *Game) Draw() {
	gosprite64.ClearScreenWith(gosprite64.DarkBlue)

	gosprite64.DrawText("MENU DEMO", 108, 20, gosprite64.White)
	gosprite64.DrawText("D-PAD: NAVIGATE  A: SELECT", 40, 40, gosprite64.LightGray)

	g.menu.Draw()

	if g.selected != "" {
		gosprite64.DrawText(
			fmt.Sprintf("SELECTED: %s", g.selected),
			60, 160, gosprite64.Yellow,
		)
	}

	gosprite64.DrawText(
		fmt.Sprintf("CURSOR: %d", g.menu.Cursor()),
		8, 200, gosprite64.DarkGray,
	)
}

func main() {
	gosprite64.Run(&Game{})
}
