package main

import (
	"image/color"

	"github.com/drpaneas/gosprite64"
)

const (
	screenW = 288
	screenH = 216
)

type Game struct {
	cameraX  int
	cameraY  int
	parallax gosprite64.ParallaxConfig
}

func (g *Game) Init() {
	g.parallax = gosprite64.NewParallaxConfig(
		gosprite64.ParallaxLayer{SpeedX: 0.2, SpeedY: 0.2},
		gosprite64.ParallaxLayer{SpeedX: 0.5, SpeedY: 0.5},
		gosprite64.ParallaxLayer{SpeedX: 1.0, SpeedY: 1.0},
	)
}

func (g *Game) Update() {
	sx, sy := gosprite64.StickPosition(0.15)
	g.cameraX += int(sx * 3)
	g.cameraY += int(sy * 3)
}

func (g *Game) Draw() {
	gosprite64.ClearScreenWith(gosprite64.DarkBlue)

	ox0, oy0 := g.parallax.LayerOffset(0, g.cameraX, g.cameraY)
	drawStarField(ox0, oy0, gosprite64.DarkGray)

	ox1, oy1 := g.parallax.LayerOffset(1, g.cameraX, g.cameraY)
	drawMountains(ox1, oy1, gosprite64.Indigo)

	ox2, oy2 := g.parallax.LayerOffset(2, g.cameraX, g.cameraY)
	drawGround(ox2, oy2, gosprite64.DarkGreen)

	gosprite64.DrawText("Parallax Demo", 92, 4, gosprite64.White)
	gosprite64.DrawText("Use stick to scroll", 72, 200, gosprite64.Yellow)
}

func drawStarField(ox, oy int, c color.Color) {
	stars := [][2]int{
		{40, 30}, {120, 50}, {200, 20}, {260, 60}, {80, 80},
		{180, 40}, {30, 70}, {240, 90}, {150, 25}, {60, 55},
	}
	for _, s := range stars {
		x := wrap(s[0]-ox, screenW)
		y := wrap(s[1]-oy, screenH/3)
		gosprite64.FillRect(x, y, x+1, y+1, c)
	}
}

func drawMountains(ox, oy int, c color.Color) {
	peaks := [][2]int{{50, 100}, {130, 85}, {220, 95}, {300, 80}, {370, 90}}
	for _, p := range peaks {
		px := wrap(p[0]-ox, screenW+200) - 100
		py := p[1] - oy%20
		for row := 0; row < 40; row++ {
			half := row * 2
			left := px - half
			right := px + half
			y := py + row
			if y >= 70 && y < 140 && left < screenW && right >= 0 {
				if left < 0 {
					left = 0
				}
				if right >= screenW {
					right = screenW - 1
				}
				gosprite64.FillRect(left, y, right, y, c)
			}
		}
	}
}

func drawGround(ox, oy int, c color.Color) {
	groundY := 140 - oy%40
	if groundY < screenH {
		gosprite64.FillRect(0, groundY, screenW-1, screenH-1, c)
	}
	stripeC := gosprite64.Brown
	for i := 0; i < 8; i++ {
		x := wrap(i*50-ox, screenW+100) - 50
		y := groundY + 10
		if y < screenH && x < screenW && x+20 >= 0 {
			left := x
			if left < 0 {
				left = 0
			}
			right := x + 20
			if right >= screenW {
				right = screenW - 1
			}
			gosprite64.FillRect(left, y, right, y+5, stripeC)
		}
	}
}

func wrap(v, mod int) int {
	if mod <= 0 {
		return v
	}
	v = v % mod
	if v < 0 {
		v += mod
	}
	return v
}

func main() {
	gosprite64.Run(&Game{})
}
