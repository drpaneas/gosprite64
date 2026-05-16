//go:generate sh -c "cd assets-src && go run gen_assets.go"
//go:generate sh -c "mkdir -p assets && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/character.png -out assets/character.sheet -tile-width 16 -tile-height 16 && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/tiles.png -out assets/tiles.sheet -tile-width 8 -tile-height 8 && go run github.com/drpaneas/gosprite64/cmd/mk2dmap -in assets-src/level.json -out assets/level.map && go run github.com/drpaneas/gosprite64/cmd/mk2danim -in assets-src/anims.json -out assets/anims.anim && go run github.com/drpaneas/gosprite64/cmd/mk2dbundle -sheet assets/tiles.sheet -map assets/level.map -anim assets/anims.anim -out assets/level.bundle"

package main

import (
	"fmt"

	"github.com/drpaneas/gosprite64"
)

type Game struct {
	scene    *gosprite64.Scene
	camera   *gosprite64.Camera
	charSS   *gosprite64.SpriteSheet
	player   *gosprite64.AnimationPlayer
	idle     gosprite64.AnimationClip
	walk     gosprite64.AnimationClip
	playerX  float32
	playerY  float32
	prevX    float32
	prevY    float32
	flipH    bool
	moving   bool
	curClip  string
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

	charSheet, err := gosprite64.LoadSpriteSheet("assets/character.sheet")
	if err != nil {
		panic(err)
	}
	g.charSS = charSheet

	animSet, err := bundle.LoadAnimation("anims")
	if err != nil {
		panic(err)
	}

	idleClip, ok := animSet.Clip("idle")
	if !ok {
		panic("idle clip not found")
	}
	g.idle = idleClip

	walkClip, ok := animSet.Clip("walk")
	if !ok {
		panic("walk clip not found")
	}
	g.walk = walkClip

	m := g.scene.Map()
	g.playerX = float32(m.PixelWidth()) / 2
	g.playerY = float32(m.PixelHeight()) / 2

	g.camera = &gosprite64.Camera{Width: 288, Height: 216}
	g.player = gosprite64.NewAnimationPlayer()
	g.player.SetLoop(true)
	g.player.Play(g.idle)
	g.curClip = "idle"
}

func (g *Game) Update() {
	g.prevX = g.playerX
	g.prevY = g.playerY
	g.moving = false

	if gosprite64.IsButtonDown(gosprite64.ButtonDPadUp) {
		g.playerY--
		g.moving = true
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadDown) {
		g.playerY++
		g.moving = true
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft) {
		g.playerX--
		g.flipH = true
		g.moving = true
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadRight) {
		g.playerX++
		g.flipH = false
		g.moving = true
	}

	if g.moving {
		if g.curClip != "walk" {
			g.player.Play(g.walk)
			g.curClip = "walk"
		}
	} else {
		if g.curClip != "idle" {
			g.player.Play(g.idle)
			g.curClip = "idle"
		}
	}
	g.player.Advance(1)

	g.camera.X = int(g.playerX) - 144
	g.camera.Y = int(g.playerY) - 108

	m := g.scene.Map()
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

	g.scene.Draw(g.camera)

	frame := g.player.Frame()

	gosprite64.DrawWorldSpriteWithOptions(g.charSS, frame, g.playerX, g.playerY+12, g.camera, gosprite64.DrawSpriteOptions{
		ScaleX: 1.5,
		ScaleY: 0.3,
		Blend:  gosprite64.BlendAlpha,
		Alpha:  0.3,
	})

	gosprite64.DrawWorldSpriteWithOptions(g.charSS, frame, g.playerX, g.playerY, g.camera, gosprite64.DrawSpriteOptions{
		FlipH: g.flipH,
	})

	gosprite64.DrawWorldSpriteWithOptions(g.charSS, frame, g.prevX, g.prevY, g.camera, gosprite64.DrawSpriteOptions{
		Blend: gosprite64.BlendAlpha,
		Alpha: 0.5,
	})

	gosprite64.DrawSpriteWithOptions(g.charSS, frame, 2, 208, gosprite64.DrawSpriteOptions{
		Blend: gosprite64.BlendAlpha,
		Alpha: 0.7,
	})

	gosprite64.DrawText(fmt.Sprintf("x:%.0f y:%.0f", g.playerX, g.playerY), 2, 2, gosprite64.White)
}

func main() {
	gosprite64.Run(&Game{})
}
