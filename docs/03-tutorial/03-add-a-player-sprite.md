# Step 3: Add a Player Sprite

Load a sprite sheet and draw a character on top of the tile world.

## What you will learn

- How to generate a character sprite sheet programmatically
- Loading a sprite sheet with `LoadSpriteSheet`
- Drawing a sprite at a world position with `DrawWorldSprite`

## The code

At this step your `main.go` looks like this:

```go
//go:generate sh -c "cd assets-src && go run gen_assets.go"
//go:generate sh -c "mkdir -p assets && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/character.png -out assets/character.sheet -tile-width 16 -tile-height 16 && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/tiles.png -out assets/tiles.sheet -tile-width 8 -tile-height 8 && go run github.com/drpaneas/gosprite64/cmd/mk2dmap -in assets-src/level.json -out assets/level.map && go run github.com/drpaneas/gosprite64/cmd/mk2dbundle -sheet assets/tiles.sheet -map assets/level.map -out assets/level.bundle"

package main

import (
	"github.com/drpaneas/gosprite64"
)

type Game struct {
	scene   *gosprite64.Scene
	camera  *gosprite64.Camera
	charSS  *gosprite64.SpriteSheet
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

	g.playerX = 144
	g.playerY = 108
}

func (g *Game) Update() {}

func (g *Game) Draw() {
	gosprite64.ClearScreen()
	g.scene.Draw(g.camera)
	gosprite64.DrawWorldSprite(g.charSS, 0, g.playerX, g.playerY, g.camera)
}

func main() {
	gosprite64.Run(&Game{})
}
```

## How it works

### Generating the character sprite

The `gen_assets.go` file (run by the first `go:generate` line) creates a 64x16 PNG containing four 16x16 frames of a simple character. Each frame is a colored stick figure - head, body, and legs drawn with basic rectangles.

The pipeline tool `mk2dsheet` slices this image into a sprite sheet with 4 frames of 16x16 pixels each.

### Loading the sprite sheet

```go
charSheet, err := gosprite64.LoadSpriteSheet("assets/character.sheet")
```

`LoadSpriteSheet` reads the compiled `.sheet` file from the embedded cartridge filesystem. It returns a `*SpriteSheet` that knows the frame count and dimensions of each frame.

### Drawing the sprite

```go
gosprite64.DrawWorldSprite(g.charSS, 0, g.playerX, g.playerY, g.camera)
```

| Parameter | Meaning |
|-----------|---------|
| `g.charSS` | Which sprite sheet to draw from |
| `0` | Frame index (first frame) |
| `g.playerX, g.playerY` | World position |
| `g.camera` | Camera (converts world coords to screen coords) |

The sprite is drawn at its world position, offset by the camera. Since we placed the player at 144, 108 (center of the 288x216 screen) and the camera is at 0, 0 - the sprite appears in the center.

### Screen vs World coordinates

`DrawSprite` uses screen coordinates (top-left of the display is 0,0). `DrawWorldSprite` uses world coordinates and subtracts the camera position automatically. Use world coordinates when your game world is larger than the screen.

## Build and run

```bash
go generate ./examples/platformer
GOENV=n64.env go1.24.5-embedded build -o examples/platformer/game.elf ./examples/platformer
```

Build and run to see the tile world with a small character sprite sitting motionless in the center of the screen. The character does not move yet - we will add animation and input in the next steps.
