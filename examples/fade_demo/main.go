package main

import (
	"image/color"

	"github.com/drpaneas/gosprite64"
)

type Game struct {
	frame      int
	transition *gosprite64.Transition
	fadeOut    bool
	bgColor    color.Color
}

func (g *Game) Init() {
	g.bgColor = gosprite64.DarkBlue
	g.transition = gosprite64.StartTransition(gosprite64.FadeFromBlack, 60)
}

func (g *Game) Update() {
	g.frame++

	if g.transition != nil {
		g.transition.Advance()
		if g.transition.Done() {
			g.transition.Stop()
			g.transition = nil
		}
	}

	if g.frame%180 == 0 {
		if g.fadeOut {
			g.transition = gosprite64.StartTransition(gosprite64.FadeFromBlack, 60)
			g.fadeOut = false
			switch g.bgColor {
			case gosprite64.DarkBlue:
				g.bgColor = gosprite64.DarkGreen
			case gosprite64.DarkGreen:
				g.bgColor = gosprite64.DarkPurple
			default:
				g.bgColor = gosprite64.DarkBlue
			}
		} else {
			g.transition = gosprite64.StartTransition(gosprite64.FadeToBlack, 60)
			g.fadeOut = true
		}
	}
}

func (g *Game) Draw() {
	gosprite64.ClearScreenWith(g.bgColor)

	gosprite64.FillRect(44, 58, 144, 108, gosprite64.Red)
	gosprite64.FillRect(144, 108, 244, 158, gosprite64.Green)
	gosprite64.FillRect(94, 28, 194, 78, gosprite64.Blue)

	gosprite64.DrawText("Fade Demo", 108, 180, gosprite64.White)

	if g.frame%60 < 30 {
		gosprite64.DrawText("Press START", 100, 196, gosprite64.Yellow)
	}

	if g.transition != nil {
		g.transition.Draw()
	}
}

func main() {
	gosprite64.Run(&Game{})
}
