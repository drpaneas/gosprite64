# Your First Tile Game

This tutorial walks you through building a complete tile-based game for the Nintendo 64 from scratch. By the end you will have a scrollable tile world with a green grass border, brown dirt patches, controller input, and a debug overlay - all running as a real N64 ROM.

No prior game development experience is required. You should be comfortable with basic Go (variables, structs, functions) and the command line.

## What you will build

A simple top-down tile world where:

- The screen shows a portion of a larger tile map
- You scroll the camera around with the D-pad
- Green tiles form walls around the edges
- Brown tiles form dirt patches scattered inside
- A debug overlay shows how many tiles are visible

This is the simplest possible tile game, but it uses the same pipeline you would use for a real project: authored assets, offline compilation, bundle loading, scene rendering, and camera control.

## Prerequisites

Complete the [Installation](../02-getting-started/installation.md) guide first. If you want a smaller warm-up before this longer tutorial, run through [Hello World](../02-getting-started/hello-world.md) too. You need these tools installed:

| Tool | Purpose |
|------|---------|
| `go` | Standard Go (dependency resolution, code generation) |
| `go1.24.5-embedded` | EmbeddedGo toolchain (cross-compiles for N64) |
| `n64go` | Converts compiled ELF binaries into N64 ROM files |
| An emulator | [ares](https://ares-emu.net/) is recommended for testing |

Verify your tools work by running `./build_examples.sh` in the GoSprite64 repository. If that prints "All examples built successfully!", you are ready.

## Step 1: Create the project

Create a new directory for your game and initialize a Go module:

```bash
mkdir -p ~/gocode/src/github.com/yourname/myfirstgame
cd ~/gocode/src/github.com/yourname/myfirstgame
go mod init github.com/yourname/myfirstgame
```

Replace `yourname` with your GitHub username or any name you like.

## Step 2: Create the toolchain config

Every GoSprite64 project needs an `n64.env` file that tells Go how to cross-compile for the N64. Create it in your project root:

```
GOTOOLCHAIN=go1.24.5-embedded
GOOS=noos
GOARCH=mips64
GOFLAGS='-tags=n64' '-trimpath' '-ldflags=-M=0x00000000:8M -F=0x00000400:8M -stripfn=1'
```

You will never need to edit this file. It is the same for every GoSprite64 project.

## Step 3: Draw your tilesheet

A tilesheet is a small PNG image that contains all the tile graphics your game uses, arranged in a grid. Each tile is a fixed-size square (8x8 pixels by default).

Create a `assets-src/` directory and draw a `tiles.png` file:

```bash
mkdir -p assets-src
```

Using any pixel editor (Aseprite, GIMP, Pixelorama, or even MS Paint), create a **16x8 pixel PNG** with two 8x8 tiles side by side:

```
+--------+--------+
| tile 1 | tile 2 |
| green  | brown  |
| (grass)| (dirt) |
+--------+--------+
  8x8 px   8x8 px
```

- **Tile 1** (left half): Fill with a green color like `#228B22`
- **Tile 2** (right half): Fill with a brown color like `#8B5A2B`

Save it as `assets-src/tiles.png`.

The important rules for tilesheets:

- The image width must be divisible by the tile width (8)
- The image height must be divisible by the tile height (8)
- Tile IDs start at 1 (tile 1 is top-left, tile 2 is next to the right, and so on)
- Tile ID 0 means "empty" - nothing is drawn

## Step 4: Design your map

A map is a grid of tile IDs that describes what your world looks like. Create `assets-src/level.json`:

```json
{
  "width": 48,
  "height": 36,
  "layer_count": 1,
  "cell_bits": 16,
  "chunk_width": 8,
  "chunk_height": 8,
  "layers": [
    {
      "sheet_id": 1,
      "cells": [
        1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,0,0,0,0,0,0,0,1,
        1,0,0,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,2,0,0,0,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,2,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,2,0,0,0,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,0,2,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,2,0,0,0,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,0,2,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,2,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,
        1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,0,0,0,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,0,0,0,2,0,0,0,0,0,0,0,0,2,2,2,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,2,2,2,0,0,0,0,0,0,0,0,2,0,2,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,2,2,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,2,2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
        1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1
      ]
    }
  ]
}
```

Here is what the numbers mean:

| Tile ID | What it draws |
|---------|--------------|
| `0` | Nothing (empty space, black background) |
| `1` | Green grass tile (from the left half of your tilesheet) |
| `2` | Brown dirt tile (from the right half of your tilesheet) |

The map is 48 tiles wide and 36 tiles tall. At 8 pixels per tile, that is 384x288 pixels - larger than the 288x216 screen, so you will be able to scroll around.

Understanding the JSON fields:

| Field | What it means |
|-------|--------------|
| `width` | How many tiles wide the map is |
| `height` | How many tiles tall the map is |
| `layer_count` | How many layers (we use 1 for simplicity) |
| `cell_bits` | How many bits per tile ID (16 allows up to 65535 tile types) |
| `chunk_width` / `chunk_height` | Internal rendering optimization (8 is a good default) |
| `sheet_id` | Which tilesheet this layer uses (1 = first sheet) |
| `cells` | The actual tile data, read left-to-right, top-to-bottom |

## Step 5: Write the game code

Your game needs three files. Here is the first and most important one.

Create `main.go`:

```go
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
	g.camera = &gosprite64.Camera{Width: 288, Height: 216}
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
```

Let's walk through what each part does:

### The `//go:generate` line

```go
//go:generate sh -c "mkdir -p assets && go run .../mk2dsheet ... && go run .../mk2dmap ... && go run .../mk2dbundle ..."
```

This tells Go's `go generate` command to run three tools in sequence:

1. `mk2dsheet` converts your `tiles.png` into a compiled `.sheet` file
2. `mk2dmap` converts your `level.json` into a compiled `.map` file
3. `mk2dbundle` packages the `.sheet` and `.map` into one `.bundle` manifest

You run this once after changing your assets. You do not run it every build.

### The `Game` struct

```go
type Game struct {
	scene  *gosprite64.Scene
	camera *gosprite64.Camera
}
```

Every GoSprite64 game is a struct that implements three methods: `Init`, `Update`, and `Draw`. The `scene` holds your loaded tile world. The `camera` defines which portion of the world is visible on screen.

### Init - load your world

```go
func (g *Game) Init() {
	bundle, err := gosprite64.OpenBundle("assets/level.bundle")
	scene, err := gosprite64.LoadScene(bundle)
	g.scene = scene
	g.camera = &gosprite64.Camera{Width: 288, Height: 216}
}
```

`OpenBundle` reads the bundle manifest. `LoadScene` loads all the sheets and map data into memory and prepares them for rendering. The camera is set to 288x216 - the full screen size - so the scene fills the entire display.

### Update - handle input

```go
func (g *Game) Update() {
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadUp) {
		g.camera.Y -= speed
	}
	// ... same for Down, Left, Right
}
```

`Update` runs every frame (60 times per second). Here we check the D-pad and move the camera. The camera position is then clamped so it cannot scroll past the edges of the map.

### Draw - render the frame

```go
func (g *Game) Draw() {
	gosprite64.ClearScreen()
	g.scene.Draw(g.camera)
	gosprite64.DrawText(fmt.Sprintf("vis:%d", stats.VisibleTiles), 2, 2, gosprite64.White)
}
```

`ClearScreen` fills the screen with black. `scene.Draw(camera)` renders only the tiles visible through the camera viewport - not the entire map. `DrawText` overlays debug info showing how many tiles are currently visible.

## Step 6: Write the asset embed file

The N64 loads assets from cartridge storage. Go's `//go:embed` directive bakes your compiled assets into the ROM binary. Create `assets_embed.go`:

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

This file does one thing: it makes your `.sheet`, `.map`, and `.bundle` files available to `OpenBundle` at runtime. Without it, the game cannot find its assets.

You write this file once. You do not need to edit it when you change your assets.

## Step 7: Resolve dependencies

```bash
env -u GOENV -u GOOS -u GOARCH -u GOFLAGS -u GOTOOLCHAIN go mod tidy
```

The `env -u ...` prefix clears any N64-specific environment variables so `go mod tidy` runs with your normal Go toolchain. This downloads GoSprite64 and all its dependencies.

## Step 8: Generate the compiled assets

```bash
go generate ./...
```

This runs the `//go:generate` line from your `main.go` and produces three files:

```
assets/
  tiles.sheet    # compiled tilesheet (binary)
  level.map      # compiled map (binary)
  level.bundle   # manifest that ties them together
```

You should see no output if everything works. If you see an error, check that your `tiles.png` is exactly 16x8 pixels and your `level.json` is valid JSON.

## Step 9: Build the ROM

```bash
GOENV=n64.env go1.24.5-embedded build -o game.elf .
GOENV=n64.env n64go rom game.elf
```

The first command cross-compiles your game for the N64 (MIPS64, no operating system). The second converts the binary into a `.z64` ROM file.

## Step 10: Run it

Open `game.z64` in the ares emulator.

You should see a tile world filling the entire screen: green grass forming a border with brown dirt patches scattered inside. Use the D-pad (arrow keys in most emulator keybindings) to scroll around. The `vis:` counter in the top-left shows how many tiles are being drawn each frame.

If the D-pad does not work, check your emulator's input settings. In ares, go to Settings > Input and make sure the N64 D-pad buttons are mapped to your keyboard arrow keys.

## Project structure

Your project should now look like this:

```
myfirstgame/
  main.go              # game code (Init, Update, Draw)
  assets_embed.go      # embeds compiled assets into the ROM
  assets-src/
    tiles.png          # your hand-drawn tilesheet (source)
    level.json         # your map layout (source)
  assets/
    tiles.sheet        # compiled tilesheet (generated)
    level.map          # compiled map (generated)
    level.bundle       # bundle manifest (generated)
  n64.env              # N64 build settings
  go.mod               # Go module file
  go.sum               # dependency checksums
  game.elf             # compiled binary (after build)
  game.z64             # N64 ROM (after build)
```

Files under `assets-src/` are your source files that you edit by hand. Files under `assets/` are generated and should not be edited. Regenerate them with `go generate ./...` whenever you change your source assets.

## What to try next

Now that you have a working tile game, here are some things to experiment with:

**Change the map layout.** Edit `level.json` and rearrange the tile IDs. Run `go generate ./...` and rebuild to see your changes.

**Add more tile types.** Make your `tiles.png` wider (for example 32x8 for 4 tiles, or 16x16 for a 4-tile grid). Each new 8x8 region becomes a new tile ID. Update your map to use the new tile IDs.

**Change the scroll speed.** In `Update()`, change `speed := 1` to `speed := 2` for faster scrolling.

**Add a second layer.** Change `layer_count` to 2 in your JSON, add a second layer with its own `sheet_id` and `cells`, and create a second tilesheet for overlay tiles (like trees or rocks on top of the ground).

**Look at the advanced example.** The `examples/tilemap` directory in the GoSprite64 repository shows a more complex setup with multiple layers, overlay sheets, and tile animations.

## Key concepts

| Concept | What it means |
|---------|--------------|
| **Tilesheet** | A PNG grid of small tile images (8x8 pixels each) |
| **Map** | A grid of tile IDs that describes your world layout |
| **Bundle** | A manifest that packages sheets and maps together |
| **Scene** | A loaded, renderable world assembled from a bundle |
| **Camera** | A viewport that controls which part of the world is visible |
| **Tile ID** | A number identifying which tile graphic to draw (0 = empty) |

## Troubleshooting

**"tiles.png not found"** - Make sure the file is at `assets-src/tiles.png`, not `assets/tiles.png`. Source assets go in `assets-src/`.

**"image size not divisible by tile size"** - Your PNG dimensions must be exact multiples of 8. A 16x8 image works. A 15x8 image does not.

**"bundle has no map"** - The bundle needs at least one `.map` file. Check that your `go generate` command includes the `mk2dmap` step.

**Scene only fills part of the screen** - Make sure your camera is `{Width: 288, Height: 216}`. A smaller camera means a smaller viewport.

**D-pad does nothing** - Check your emulator input settings. Also make sure your map is larger than 288x216 pixels (larger than 36x27 tiles), otherwise there is nowhere to scroll.

**Black screen** - Check that `assets_embed.go` exists and contains `gosprite64.RegisterAssetFS(assetFS)`. Without it, the game cannot load any assets.
