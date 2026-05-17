# Step 4: Animate the Player

Use the AnimationPlayer to cycle through sprite frames automatically.

## What you will learn

- Loading animation clips from a bundle
- Using `AnimationPlayer` to drive frame playback
- The Play / Advance / Frame pattern

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
	playerX float32
	playerY float32
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

	g.playerX = 144
	g.playerY = 108

	g.player = gosprite64.NewAnimationPlayer()
	g.player.SetLoop(true)
	g.player.Play(g.idle)
}

func (g *Game) Update() {
	g.player.Advance(1)
}

func (g *Game) Draw() {
	gosprite64.ClearScreen()
	g.scene.Draw(g.camera)

	frame := g.player.Frame()
	gosprite64.DrawWorldSprite(g.charSS, frame, g.playerX, g.playerY, g.camera)
}

func main() {
	gosprite64.Run(&Game{})
}
```

## How it works

### Animation data

The `anims.json` file defines two clips:

```json
{
  "clips": [
    {"name": "idle", "fps": 4, "frames": [0, 1]},
    {"name": "walk", "fps": 8, "frames": [0, 1, 2, 3]}
  ]
}
```

Each clip has a name, a playback rate (frames per second), and a list of sprite frame indices. The idle clip alternates between frames 0 and 1 at 4 FPS. The walk clip cycles all 4 frames at 8 FPS.

The `mk2danim` tool compiles this JSON into a binary `.anim` file, and `mk2dbundle` includes it in the bundle.

### Loading clips

```go
animSet, err := bundle.LoadAnimation("anims")
idleClip, ok := animSet.Clip("idle")
```

`LoadAnimation` loads the compiled animation data. `Clip("idle")` retrieves a specific named clip. An `AnimationClip` holds the FPS and frame sequence.

### The AnimationPlayer

```go
g.player = gosprite64.NewAnimationPlayer()
g.player.SetLoop(true)
g.player.Play(g.idle)
```

The `AnimationPlayer` tracks which clip is playing, the current frame index, and an internal accumulator that converts game ticks into animation frames at the clip's FPS rate.

| Method | What it does |
|--------|-------------|
| `Play(clip)` | Start playing a clip from the beginning |
| `SetLoop(true)` | Repeat the clip when it reaches the end |
| `Advance(ticks)` | Advance the internal clock by N game ticks |
| `Frame()` | Return the current sprite frame index |

### The pattern

In `Update()`, call `Advance(1)` every frame. In `Draw()`, call `Frame()` to get the current frame index and pass it to `DrawWorldSprite`. The player handles all the FPS math internally.

## Build and run

```bash
go generate ./examples/platformer
GOENV=n64.env go1.24.5-embedded build -o examples/platformer/game.elf ./examples/platformer
```

Build and run to see the character animating in place. The idle animation slowly alternates between two frames - a subtle breathing effect. The character still does not respond to input yet.
