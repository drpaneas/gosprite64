//go:generate sh -c "cd assets-src && go run gen_assets.go"
//go:generate sh -c "mkdir -p assets && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/character.png -out assets/character.sheet -tile-width 16 -tile-height 16 && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/tiles.png -out assets/tiles.sheet -tile-width 8 -tile-height 8 && go run github.com/drpaneas/gosprite64/cmd/mk2dmap -in assets-src/level.json -out assets/level.map && go run github.com/drpaneas/gosprite64/cmd/mk2danim -in assets-src/anims.json -out assets/anims.anim && go run github.com/drpaneas/gosprite64/cmd/mk2dbundle -sheet assets/tiles.sheet -map assets/level.map -anim assets/anims.anim -out assets/level.bundle"

package main

import (
	"fmt"

	"github.com/drpaneas/gosprite64"
	"github.com/drpaneas/gosprite64/math2d"
)

// TitleState shows a start screen before gameplay begins.
type TitleState struct {
	sm    *gosprite64.StateMachine
	blink int
	show  bool
	fade  *gosprite64.Transition
}

func (s *TitleState) Enter() {
	s.show = true
	s.fade = gosprite64.StartTransition(gosprite64.FadeFromBlack, 30)
}

func (s *TitleState) Update() {
	if s.fade != nil {
		s.fade.Advance()
		if s.fade.Done() {
			s.fade.Stop()
			s.fade = nil
		}
	}
	s.blink++
	if s.blink%30 == 0 {
		s.show = !s.show
	}
	if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) || gosprite64.IsButtonJustPressed(gosprite64.ButtonStart) {
		s.fade = gosprite64.StartTransition(gosprite64.FadeToBlack, 20)
		s.sm.Switch(&PlayState{sm: s.sm})
	}
}

func (s *TitleState) Draw() {
	gosprite64.ClearScreenWith(gosprite64.DarkBlue)
	gosprite64.DrawText("PLATFORMER", 104, 60, gosprite64.White)
	gosprite64.DrawText("A GoSprite64 Tutorial", 60, 80, gosprite64.LightGray)
	if s.show {
		gosprite64.DrawText("PRESS START", 96, 140, gosprite64.Yellow)
	}
	if s.fade != nil {
		s.fade.Draw()
	}
}

func (s *TitleState) Exit() {}

// PlayState holds all gameplay state.
type PlayState struct {
	sm      *gosprite64.StateMachine
	scene   *gosprite64.Scene
	camera  *gosprite64.Camera
	charSS  *gosprite64.SpriteSheet
	player  *gosprite64.AnimationPlayer
	idle    gosprite64.AnimationClip
	walk    gosprite64.AnimationClip
	playerX float32
	playerY float32
	flipH   bool
	moving  bool
	curClip string
	score   int
	fade    *gosprite64.Transition
}

func (s *PlayState) Enter() {
	bundle, err := gosprite64.OpenBundle("assets/level.bundle")
	if err != nil {
		panic(err)
	}

	scene, err := gosprite64.LoadScene(bundle)
	if err != nil {
		panic(err)
	}
	s.scene = scene

	charSheet, err := gosprite64.LoadSpriteSheet("assets/character.sheet")
	if err != nil {
		panic(err)
	}
	s.charSS = charSheet

	animSet, err := bundle.LoadAnimation("anims")
	if err != nil {
		panic(err)
	}

	idleClip, ok := animSet.Clip("idle")
	if !ok {
		panic("idle clip not found")
	}
	s.idle = idleClip

	walkClip, ok := animSet.Clip("walk")
	if !ok {
		panic("walk clip not found")
	}
	s.walk = walkClip

	s.playerX = 80
	s.playerY = 180

	s.camera = &gosprite64.Camera{
		Width:       288,
		Height:      216,
		FollowSpeed: 0.1,
	}
	s.camera.FollowTarget = &math2d.Vec2{X: s.playerX, Y: s.playerY}
	s.camera.Bounds = &math2d.Rect{
		X: 0, Y: 0,
		W: float32(scene.Map().PixelWidth()),
		H: float32(scene.Map().PixelHeight()),
	}

	s.player = gosprite64.NewAnimationPlayer()
	s.player.SetLoop(true)
	s.player.Play(s.idle)
	s.curClip = "idle"

	s.fade = gosprite64.StartTransition(gosprite64.FadeFromBlack, 30)
}

func (s *PlayState) Update() {
	if s.fade != nil {
		s.fade.Advance()
		if s.fade.Done() {
			s.fade.Stop()
			s.fade = nil
		}
	}

	s.moving = false
	speed := float32(2)

	if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft) {
		s.playerX -= speed
		s.flipH = true
		s.moving = true
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadRight) {
		s.playerX += speed
		s.flipH = false
		s.moving = true
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadUp) {
		s.playerY -= speed
		s.moving = true
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadDown) {
		s.playerY += speed
		s.moving = true
	}

	if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) {
		s.score++
	}

	if s.moving {
		if s.curClip != "walk" {
			s.player.Play(s.walk)
			s.curClip = "walk"
		}
	} else {
		if s.curClip != "idle" {
			s.player.Play(s.idle)
			s.curClip = "idle"
		}
	}
	s.player.Advance(1)

	s.camera.FollowTarget.X = s.playerX
	s.camera.FollowTarget.Y = s.playerY
	s.camera.UpdateFollow()
	s.camera.ClampToBounds()
}

func (s *PlayState) Draw() {
	gosprite64.ClearScreen()

	s.scene.Draw(s.camera)

	frame := s.player.Frame()
	gosprite64.DrawWorldSpriteWithOptions(s.charSS, frame, s.playerX, s.playerY, s.camera, gosprite64.DrawSpriteOptions{
		FlipH: s.flipH,
	})

	gosprite64.DrawText(fmt.Sprintf("SCORE:%d", s.score), 2, 2, gosprite64.White)

	if s.fade != nil {
		s.fade.Draw()
	}
}

func (s *PlayState) Exit() {}

// Game is the top-level wrapper that holds the state machine.
type Game struct {
	sm *gosprite64.StateMachine
}

func (g *Game) Init() {
	title := &TitleState{}
	g.sm = gosprite64.NewStateMachine(title)
	title.sm = g.sm
	g.sm.Init()
}

func (g *Game) Update() { g.sm.Update() }
func (g *Game) Draw()   { g.sm.Draw() }

func main() {
	gosprite64.Run(&Game{})
}
