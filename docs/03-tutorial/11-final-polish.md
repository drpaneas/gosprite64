# Step 11: Final Polish

Review the complete game and look at what comes next.

## The complete game

Here is the final `main.go` with every feature from the tutorial. The `go:generate` lines are omitted for brevity - they are unchanged from Step 2.

```go
package main

import (
	"fmt"

	"github.com/drpaneas/gosprite64"
	"github.com/drpaneas/gosprite64/math2d"
)

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
		if s.fade.Done() { s.fade.Stop(); s.fade = nil }
	}
	s.blink++
	if s.blink%30 == 0 { s.show = !s.show }
	if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) || gosprite64.IsButtonJustPressed(gosprite64.ButtonStart) {
		s.fade = gosprite64.StartTransition(gosprite64.FadeToBlack, 20)
		s.sm.Switch(&PlayState{sm: s.sm})
	}
}

func (s *TitleState) Draw() {
	gosprite64.ClearScreenWith(gosprite64.DarkBlue)
	gosprite64.DrawText("PLATFORMER", 104, 60, gosprite64.White)
	gosprite64.DrawText("A GoSprite64 Tutorial", 60, 80, gosprite64.LightGray)
	if s.show { gosprite64.DrawText("PRESS START", 96, 140, gosprite64.Yellow) }
	if s.fade != nil { s.fade.Draw() }
}

func (s *TitleState) Exit() {}

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
	if err != nil { panic(err) }
	scene, err := gosprite64.LoadScene(bundle)
	if err != nil { panic(err) }
	s.scene = scene
	charSheet, err := gosprite64.LoadSpriteSheet("assets/character.sheet")
	if err != nil { panic(err) }
	s.charSS = charSheet
	animSet, err := bundle.LoadAnimation("anims")
	if err != nil { panic(err) }
	idleClip, ok := animSet.Clip("idle")
	if !ok { panic("idle clip not found") }
	s.idle = idleClip
	walkClip, ok := animSet.Clip("walk")
	if !ok { panic("walk clip not found") }
	s.walk = walkClip

	s.playerX, s.playerY = 80, 180
	s.camera = &gosprite64.Camera{Width: 288, Height: 216, FollowSpeed: 0.1}
	s.camera.FollowTarget = &math2d.Vec2{X: s.playerX, Y: s.playerY}
	s.camera.Bounds = &math2d.Rect{X: 0, Y: 0,
		W: float32(scene.Map().PixelWidth()), H: float32(scene.Map().PixelHeight())}
	s.player = gosprite64.NewAnimationPlayer()
	s.player.SetLoop(true)
	s.player.Play(s.idle)
	s.curClip = "idle"
	s.fade = gosprite64.StartTransition(gosprite64.FadeFromBlack, 30)
}

func (s *PlayState) Update() {
	if s.fade != nil {
		s.fade.Advance()
		if s.fade.Done() { s.fade.Stop(); s.fade = nil }
	}
	s.moving = false
	speed := float32(2)
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft)  { s.playerX -= speed; s.flipH = true; s.moving = true }
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadRight) { s.playerX += speed; s.flipH = false; s.moving = true }
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadUp)    { s.playerY -= speed; s.moving = true }
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadDown)  { s.playerY += speed; s.moving = true }
	if gosprite64.IsButtonJustPressed(gosprite64.ButtonA)  { s.score++ }
	if s.moving {
		if s.curClip != "walk" { s.player.Play(s.walk); s.curClip = "walk" }
	} else {
		if s.curClip != "idle" { s.player.Play(s.idle); s.curClip = "idle" }
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
	gosprite64.DrawWorldSpriteWithOptions(s.charSS, frame, s.playerX, s.playerY, s.camera, gosprite64.DrawSpriteOptions{FlipH: s.flipH})
	gosprite64.DrawText(fmt.Sprintf("SCORE:%d", s.score), 2, 2, gosprite64.White)
	if s.fade != nil { s.fade.Draw() }
}

func (s *PlayState) Exit() {}

type Game struct{ sm *gosprite64.StateMachine }

func (g *Game) Init() {
	title := &TitleState{}
	g.sm = gosprite64.NewStateMachine(title)
	title.sm = g.sm
	g.sm.Init()
}
func (g *Game) Update() { g.sm.Update() }
func (g *Game) Draw()   { g.sm.Draw() }

func main() { gosprite64.Run(&Game{}) }
```

## What you built

Over 11 steps you built a complete N64 game from scratch:

| Step | What you added |
|------|---------------|
| 1. Start the Engine | Game interface, Run loop, solid color screen |
| 2. Draw a Tilemap | Asset pipeline, bundles, tile rendering |
| 3. Add a Player Sprite | Sprite sheets, DrawWorldSprite |
| 4. Animate the Player | AnimationPlayer, clips, Play/Advance/Frame |
| 5. Move with D-Pad | IsButtonDown, movement, sprite flip |
| 6. Camera Following | FollowTarget, FollowSpeed, ClampToBounds |
| 7. Add Sound Effects | VADPCM pipeline, PlaySoundEffect |
| 8. Title Screen | GameState interface, StateMachine |
| 9. Screen Transitions | FadeToBlack, FadeFromBlack |
| 10. Score Display | DrawText, HUD overlay, IsButtonJustPressed |
| 11. Final Polish | Complete game, review, next steps |

## Suggested next steps

- **Add enemies** - create an enemy struct, draw it as a second sprite sheet, reverse direction at boundaries
- **Add collision** - use `math2d.AABBOverlap` to check player vs enemy or collectible rectangles
- **Add save data** - use `gosprite64.SaveData` and `gosprite64.LoadData` to persist the high score to SRAM
- **Add a pause menu** - use `sm.Push(&PauseState{})` to pause without destroying gameplay, `sm.Pop()` to resume
- **Add camera shake** - call `camera.AddTrauma(0.5)` on hit, `camera.UpdateShake()` each frame, apply `camera.ShakeOffset()` to draws
- **Try the other examples** - pong with audio, multi-layer tilemap, 3D triangle renderer, and more in `examples/`

## Try It

<iframe src="../emulator/play.html?rom=platformer.z64" width="640" height="480" frameborder="0" allow="autoplay" style="display:block;margin:0 auto;max-width:100%;"></iframe>

> **Controls:** Arrow keys = D-Pad, X = A button, C = B button, Enter = Start, Z = Z trigger

## Build and run

```bash
go generate ./examples/platformer
GOENV=n64.env go1.24.5-embedded build -o examples/platformer/game.elf ./examples/platformer
```

You now have a real N64 ROM with a title screen, smooth camera, animated character, score display, and fade transitions. From here, every new feature is just another struct, another state, and another call to the GoSprite64 API.
