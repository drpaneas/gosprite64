package main

import (
	"image/color"

	"github.com/drpaneas/gosprite64"
	"github.com/drpaneas/gosprite64/math3d"
)

const (
	screenW = 288
	screenH = 216
)

type Game struct {
	angle float32
	verts [3]math3d.Vec3
}

func (g *Game) Init() {
	g.verts = [3]math3d.Vec3{
		{X: 144, Y: 40, Z: 0},
		{X: 60, Y: 180, Z: 0},
		{X: 228, Y: 180, Z: 0},
	}
}

func (g *Game) Update() {
	g.angle += 0.5
	if g.angle >= 360 {
		g.angle -= 360
	}
}

func (g *Game) Draw() {
	gosprite64.ClearScreenWith(gosprite64.DarkBlue)

	cx := float32(screenW / 2)
	cy := float32(screenH / 2)

	rot := math3d.Rotate(g.angle, 0, 0, 1)

	var rotated [3]math3d.Vec3
	for i, v := range g.verts {
		rel := math3d.Vec4{X: v.X - cx, Y: v.Y - cy, Z: 0, W: 1}
		out := rot.MulVec4(rel)
		rotated[i] = math3d.Vec3{X: out.X + cx, Y: out.Y + cy, Z: 0}
	}

	fillTriangle(rotated[0], rotated[1], rotated[2], gosprite64.Red)

	gosprite64.DrawText("GoSprite64 3D", 96, 4, gosprite64.White)
	gosprite64.DrawText("Triangle", 116, 200, gosprite64.Yellow)
}

func fillTriangle(a, b, c math3d.Vec3, col color.Color) {
	if a.Y > b.Y {
		a, b = b, a
	}
	if a.Y > c.Y {
		a, c = c, a
	}
	if b.Y > c.Y {
		b, c = c, b
	}

	yStart := int(a.Y)
	yEnd := int(c.Y)
	if yStart < 0 {
		yStart = 0
	}
	if yEnd >= screenH {
		yEnd = screenH - 1
	}

	for y := yStart; y <= yEnd; y++ {
		fy := float32(y) + 0.5

		x1 := edgeX(a, c, fy)

		var x2 float32
		if fy < b.Y {
			x2 = edgeX(a, b, fy)
		} else {
			x2 = edgeX(b, c, fy)
		}

		if x1 > x2 {
			x1, x2 = x2, x1
		}

		left := int(x1)
		right := int(x2)
		if left < 0 {
			left = 0
		}
		if right >= screenW {
			right = screenW - 1
		}
		if left <= right {
			gosprite64.FillRect(left, y, right, y, col)
		}
	}
}

func edgeX(from, to math3d.Vec3, y float32) float32 {
	dy := to.Y - from.Y
	if dy == 0 {
		return from.X
	}
	t := (y - from.Y) / dy
	return from.X + t*(to.X-from.X)
}

func main() {
	gosprite64.Run(&Game{})
}
