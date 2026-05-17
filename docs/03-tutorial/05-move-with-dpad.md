# Step 5: Move with D-Pad

Read controller input to move the player and switch between idle and walk animations.

## What you will learn

- Reading D-pad input with `IsButtonDown`
- Moving the player by changing world position
- Flipping the sprite horizontally with `DrawSpriteOptions`
- Switching animation clips based on state

## The code

```go
//go:generate sh -c "cd assets-src && go run gen_assets.go"
//go:generate sh -c "mkdir -p assets && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/character.png -out assets/character.sheet -tile-width 16 -tile-height 16 && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/tiles.png -out assets/tiles.sheet -tile-width 8 -tile-height 8 && go run github.com/drpaneas/gosprite64/cmd/mk2dmap -in assets-src/level.json -out assets/level.map && go run github.com/drpaneas/gosprite64/cmd/mk2danim -in assets-src/anims.json -out assets/anims.anim && go run github.com/drpaneas/gosprite64/cmd/mk2dbundle -sheet assets/tiles.sheet -map assets/level.map -anim assets/anims.anim -out assets/level.bundle"

package main

import (
	"github.com/drpaneas/gosprite64"
)

type Game struct {
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
	g.camera = &gosprite64.Camera{Width: 288, Height: 216}

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

	g.playerX = 144
	g.playerY = 108

	g.player = gosprite64.NewAnimationPlayer()
	g.player.SetLoop(true)
	g.player.Play(g.idle)
	g.curClip = "idle"
}

func (g *Game) Update() {
	g.moving = false
	speed := float32(2)

	if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft) {
		g.playerX -= speed
		g.flipH = true
		g.moving = true
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadRight) {
		g.playerX += speed
		g.flipH = false
		g.moving = true
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadUp) {
		g.playerY -= speed
		g.moving = true
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadDown) {
		g.playerY += speed
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
}

func (g *Game) Draw() {
	gosprite64.ClearScreen()
	g.scene.Draw(g.camera)

	frame := g.player.Frame()
	gosprite64.DrawWorldSpriteWithOptions(g.charSS, frame, g.playerX, g.playerY, g.camera, gosprite64.DrawSpriteOptions{
		FlipH: g.flipH,
	})
}

func main() {
	gosprite64.Run(&Game{})
}
```

## How it works

### Reading input

```go
if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft) {
    g.playerX -= speed
    g.flipH = true
    g.moving = true
}
```

`IsButtonDown` returns true every frame the button is held. The D-pad constants are:

| Constant | Button |
|----------|--------|
| `ButtonDPadUp` | D-pad up |
| `ButtonDPadDown` | D-pad down |
| `ButtonDPadLeft` | D-pad left |
| `ButtonDPadRight` | D-pad right |

Each frame, we check all four directions. Multiple directions can be pressed simultaneously for diagonal movement. We track whether any direction was pressed in `g.moving`.

### Horizontal flip

When moving left, we set `g.flipH = true`. When moving right, `g.flipH = false`. In Draw, we pass this to `DrawSpriteOptions`:

```go
gosprite64.DrawWorldSpriteWithOptions(g.charSS, frame, g.playerX, g.playerY, g.camera, gosprite64.DrawSpriteOptions{
    FlipH: g.flipH,
})
```

`DrawWorldSpriteWithOptions` works like `DrawWorldSprite` but accepts extra options. `FlipH: true` mirrors the sprite horizontally so the character faces left without needing separate left-facing art.

### Switching animations

```go
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
```

We only call `Play()` when the clip actually changes. Calling `Play()` resets the animation to frame 0, so calling it every frame would prevent the animation from advancing. The `curClip` string tracks which clip is currently active.

The walk clip plays at 8 FPS (faster leg movement), while idle plays at 4 FPS (slow breathing). The AnimationPlayer handles the FPS conversion internally.

## Build and run

```bash
go generate ./examples/platformer
GOENV=n64.env go1.24.5-embedded build -o examples/platformer/game.elf ./examples/platformer
```

Build and run to see the character respond to the D-pad. Push left or right to walk (the sprite flips direction). Push up or down to move vertically. Release all buttons to see the idle animation. The walk animation is noticeably faster than idle.
