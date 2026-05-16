//go:generate sh -c "mkdir -p assets && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/tiles.png -out assets/tiles.sheet -tile-width 8 -tile-height 8 && go run github.com/drpaneas/gosprite64/cmd/mk2dmap -in assets-src/level.json -out assets/level.map && go run github.com/drpaneas/gosprite64/cmd/mk2dbundle -sheet assets/tiles.sheet -map assets/level.map -out assets/level.bundle"

package main

import (
	"fmt"

	"github.com/drpaneas/gosprite64"
)

type Game struct {
	scene  *gosprite64.Scene
	camera *gosprite64.Camera
}

func (g *Game) Init() {
	bundle, err := gosprite64.OpenBundle("assets/level.bundle")
	if err != nil {
		panic(err)
	}

	scene, err := gosprite64.LoadScene(bundle)
	if err != nil {
		panic(err)
	}

	g.scene = scene
	g.camera = &gosprite64.Camera{Width: 128, Height: 96}
}

func (g *Game) Update() {
	if g.camera == nil {
		return
	}

	speed := 1

	if gosprite64.IsButtonDown(gosprite64.ButtonDPadUp) {
		g.camera.Y -= speed
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadDown) {
		g.camera.Y += speed
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft) {
		g.camera.X -= speed
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadRight) {
		g.camera.X += speed
	}

	m := g.scene.Map()
	if m == nil {
		return
	}

	maxX := m.PixelWidth() - g.camera.Width
	maxY := m.PixelHeight() - g.camera.Height
	if g.camera.X < 0 {
		g.camera.X = 0
	}
	if g.camera.Y < 0 {
		g.camera.Y = 0
	}
	if g.camera.X > maxX {
		g.camera.X = maxX
	}
	if g.camera.Y > maxY {
		g.camera.Y = maxY
	}
}

func (g *Game) Draw() {
	gosprite64.ClearScreen()

	if g.scene != nil && g.camera != nil {
		g.scene.Draw(g.camera)
	}

	if g.scene != nil {
		stats := g.scene.Stats()
		gosprite64.DrawText(fmt.Sprintf("vis:%d", stats.VisibleTiles), 2, 2, gosprite64.White)
	}
}

func main() {
	gosprite64.Run(&Game{})
}
