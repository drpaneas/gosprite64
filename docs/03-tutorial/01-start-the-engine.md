# Step 1: Start the Engine

Set up a minimal GoSprite64 project that compiles, runs, and draws a solid blue screen.

## What you will learn

- The Game interface (Init, Update, Draw)
- How `gosprite64.Run` starts the engine
- How to clear the screen with a color
- How to build and run an N64 ROM

## The code

Create `examples/platformer/main.go`:

```go
package main

import (
	"github.com/drpaneas/gosprite64"
)

type Game struct{}

func (g *Game) Init()   {}
func (g *Game) Update() {}

func (g *Game) Draw() {
	gosprite64.ClearScreenWith(gosprite64.Blue)
}

func main() {
	gosprite64.Run(&Game{})
}
```

## How it works

Every GoSprite64 game is a struct that implements three methods:

| Method | Called | Purpose |
|--------|--------|---------|
| `Init()` | Once at startup | Load assets, set initial state |
| `Update()` | 60 times per second | Game logic, input handling |
| `Draw()` | Every frame after Update | Render everything to the screen |

`gosprite64.Run(&Game{})` boots the N64 hardware, initializes the display, calls your `Init()` once, then enters an infinite loop calling `Update()` and `Draw()` at 60 FPS.

`ClearScreenWith(color)` fills the entire 288x216 pixel framebuffer with a solid color. The engine provides 16 built-in colors (Black, DarkBlue, Blue, Green, Red, White, Yellow, etc.). Here we use `Blue` to get a sky-colored background.

## Build and run

You also need an asset embed file. Create `examples/platformer/assets_embed.go`:

```go
package main

import (
	"embed"

	"github.com/clktmr/n64/drivers/cartfs"
	"github.com/drpaneas/gosprite64"
)

//go:embed assets/*
var embeddedAssets embed.FS

var assetFS = cartfs.Embed(embeddedAssets)

func init() {
	gosprite64.RegisterAssetFS(assetFS)
}
```

For now, create a placeholder assets directory so the embed does not fail:

```bash
mkdir -p examples/platformer/assets
touch examples/platformer/assets/.gitkeep
```

Build and run:

```bash
GOENV=n64.env go1.24.5-embedded build -o examples/platformer/game.elf ./examples/platformer
n64go rom examples/platformer/game.elf
```

Open `game.z64` in the ares emulator. You should see a solid blue screen filling the display. That is your game running on N64 hardware - the simplest possible starting point.

## What comes next

In the next steps we will add a tile world, a player character, animation, and controls. But the structure will always be the same: load things in Init, update state in Update, draw everything in Draw.
