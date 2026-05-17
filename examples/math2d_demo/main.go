package main

import (
	"fmt"
	"image/color"

	"github.com/drpaneas/gosprite64"
	"github.com/drpaneas/gosprite64/math2d"
)

const (
	screenW = 288
	screenH = 216
	margin  = 8
	maxBalls = 12
)

var palette = []color.Color{
	gosprite64.Red,
	gosprite64.Orange,
	gosprite64.Yellow,
	gosprite64.Green,
	gosprite64.Blue,
	gosprite64.Pink,
	gosprite64.Indigo,
	gosprite64.Peach,
}

type Ball struct {
	pos   math2d.Vec2
	vel   math2d.Vec2
	size  float32
	color color.Color
}

func (b *Ball) Bounds() math2d.Rect {
	return math2d.Rect{X: b.pos.X, Y: b.pos.Y, W: b.size, H: b.size}
}

type Game struct {
	rng   *math2d.Rand
	balls []Ball
	court math2d.Rect
	frame int

	trailX    float32
	trailY    float32
	easeTimer float32
}

func (g *Game) Init() {
	g.rng = math2d.NewRand(64)
	g.court = math2d.Rect{X: margin, Y: margin + 12, W: screenW - margin*2, H: screenH - margin*2 - 12}

	for i := 0; i < 6; i++ {
		g.spawnBall()
	}
}

func (g *Game) spawnBall() {
	if len(g.balls) >= maxBalls {
		return
	}
	size := g.rng.RangeFloat32(3, 8)
	b := Ball{
		pos: math2d.Vec2{
			X: g.rng.RangeFloat32(g.court.X+size, g.court.Right()-size),
			Y: g.rng.RangeFloat32(g.court.Y+size, g.court.Bottom()-size),
		},
		vel: math2d.Vec2{
			X: g.rng.RangeFloat32(-2, 2),
			Y: g.rng.RangeFloat32(-2, 2),
		},
		size:  size,
		color: palette[g.rng.Intn(len(palette))],
	}
	if b.vel.LengthSq() < 0.5 {
		b.vel = math2d.Vec2{X: 1.5, Y: -1.0}
	}
	g.balls = append(g.balls, b)
}

func (g *Game) Update() {
	g.frame++

	if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) {
		g.spawnBall()
	}

	if gosprite64.IsButtonJustPressed(gosprite64.ButtonB) && len(g.balls) > 1 {
		g.balls = g.balls[:len(g.balls)-1]
	}

	g.easeTimer += 1.0 / 120.0
	if g.easeTimer > 1.0 {
		g.easeTimer = 0
	}

	sx, sy := gosprite64.StickPosition(0.2)
	g.trailX = math2d.MoveToward(g.trailX, float32(sx)*80+float32(screenW)/2, 2)
	g.trailY = math2d.MoveToward(g.trailY, float32(sy)*80+float32(screenH)/2, 2)

	for i := range g.balls {
		b := &g.balls[i]
		b.pos = b.pos.Add(b.vel)

		bounds := b.Bounds()

		if bounds.X <= g.court.X {
			b.pos.X = g.court.X
			b.vel.X = -b.vel.X
		}
		if bounds.Right() >= g.court.Right() {
			b.pos.X = g.court.Right() - b.size
			b.vel.X = -b.vel.X
		}
		if bounds.Y <= g.court.Y {
			b.pos.Y = g.court.Y
			b.vel.Y = -b.vel.Y
		}
		if bounds.Bottom() >= g.court.Bottom() {
			b.pos.Y = g.court.Bottom() - b.size
			b.vel.Y = -b.vel.Y
		}
	}

	for i := 0; i < len(g.balls); i++ {
		for j := i + 1; j < len(g.balls); j++ {
			bi := g.balls[i].Bounds()
			bj := g.balls[j].Bounds()
			if bi.Overlaps(bj) {
				g.balls[i].vel.X, g.balls[j].vel.X = g.balls[j].vel.X, g.balls[i].vel.X
				g.balls[i].vel.Y, g.balls[j].vel.Y = g.balls[j].vel.Y, g.balls[i].vel.Y

				sep := g.balls[i].pos.Sub(g.balls[j].pos).Normalize().Scale(1)
				g.balls[i].pos = g.balls[i].pos.Add(sep)
				g.balls[j].pos = g.balls[j].pos.Sub(sep)
			}
		}
	}
}

func (g *Game) Draw() {
	gosprite64.ClearScreen()

	gosprite64.DrawRect(
		int(g.court.X), int(g.court.Y),
		int(g.court.Right()), int(g.court.Bottom()),
		gosprite64.DarkGray,
	)

	easeVal := math2d.EaseInOutQuad(g.easeTimer)
	crosshairSize := int(math2d.Lerp(2, 6, easeVal))
	cx, cy := int(g.trailX), int(g.trailY)
	if g.court.ContainsPoint(math2d.Vec2{X: g.trailX, Y: g.trailY}) {
		gosprite64.DrawLine(cx-crosshairSize, cy, cx+crosshairSize, cy, gosprite64.LightGray)
		gosprite64.DrawLine(cx, cy-crosshairSize, cx, cy+crosshairSize, gosprite64.LightGray)
	}

	for _, b := range g.balls {
		gosprite64.FillRect(
			int(b.pos.X), int(b.pos.Y),
			int(b.pos.X+b.size), int(b.pos.Y+b.size),
			b.color,
		)
	}

	for i := 0; i < len(g.balls); i++ {
		for j := i + 1; j < len(g.balls); j++ {
			dist := g.balls[i].pos.Distance(g.balls[j].pos)
			if dist < 40 {
				gosprite64.DrawLine(
					int(g.balls[i].pos.X+g.balls[i].size/2),
					int(g.balls[i].pos.Y+g.balls[i].size/2),
					int(g.balls[j].pos.X+g.balls[j].size/2),
					int(g.balls[j].pos.Y+g.balls[j].size/2),
					gosprite64.DarkGray,
				)
			}
		}
	}

	gosprite64.DrawText(fmt.Sprintf("balls:%d A=add B=del", len(g.balls)), margin, 2, gosprite64.White)

	dist := math2d.Vec2{X: g.trailX, Y: g.trailY}.Distance(g.court.Center())
	gosprite64.DrawText(fmt.Sprintf("dist:%.0f ease:%.2f", dist, easeVal), margin+150, 2, gosprite64.Yellow)
}

func main() {
	gosprite64.Run(&Game{})
}
