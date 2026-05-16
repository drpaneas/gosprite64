# Using the Tile2D Scene Pipeline

This chapter covers how to author tile assets, package them into a bundle, and render a tile scene at runtime.

## Quick start

GoSprite64 uses a build-time pipeline for tile scenes. You provide source assets (PNG tilesheets, JSON map descriptions, JSON animation clips), the `mk2d*` tools compile them into compact binary formats, and the runtime loads and draws them. You do not deal with binary layouts or renderer details in your gameplay code.

Four steps:

1. Put your source assets in the right directories.
2. Run `go generate`.
3. Call `gosprite64.OpenBundle` and `gosprite64.LoadScene` from your game.
4. Call `scene.Draw(camera)` every frame.

## Source assets

Organize source assets under `assets-src/` in your game directory:

- `assets-src/tiles.png` - PNG tilesheet atlas (must be divisible by tile size)
- `assets-src/level.json` - JSON map description
- `assets-src/idle.json` - JSON animation clip description (optional)

### PNG tilesheet requirements

The `mk2dsheet` tool accepts:

- PNG images
- Dimensions must be evenly divisible by the tile size (default 8x8)
- Pixels are stored as NRGBA internally

### JSON map format

The `mk2dmap` tool accepts JSON with this shape:

```json
{
  "width": 32,
  "height": 18,
  "layer_count": 2,
  "cell_bits": 16,
  "chunk_width": 8,
  "chunk_height": 8,
  "layers": [
    {"sheet_id": 1, "cells": [1, 2, 3]},
    {"sheet_id": 2, "cells": [0, 0, 1]}
  ]
}
```

### JSON animation format

The `mk2danim` tool accepts JSON with this shape:

```json
{
  "clips": [
    {"name": "idle", "fps": 12, "frames": [0, 1, 2, 3]}
  ]
}
```

## Setting up go generate

Add a `go:generate` line to your `main.go`:

```go
//go:generate sh -c "mkdir -p assets && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/tiles.png -out assets/tiles.sheet -tile-width 8 -tile-height 8 && go run github.com/drpaneas/gosprite64/cmd/mk2dmap -in assets-src/level.json -out assets/level.map && go run github.com/drpaneas/gosprite64/cmd/mk2dbundle -sheet assets/tiles.sheet -map assets/level.map -out assets/level.bundle"
```

Then run:

```bash
go generate ./...
```

This will:

- compile the PNG into a `.sheet` binary
- compile the JSON map into a `.map` binary
- package everything into a `.bundle` manifest

## Embedding assets for the N64

The N64 loads assets from cartridge storage using an embedded filesystem. Create `assets_embed.go` in your game directory:

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

This makes the generated `.sheet`, `.map`, and `.bundle` files available to `OpenBundle` at runtime. Without this file, `OpenBundle` cannot find the assets on the N64.

## Runtime usage

```go
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

func (g *Game) Draw() {
    gosprite64.ClearScreen()
    g.scene.Draw(g.camera)
}
```

`OpenBundle` reads the bundle manifest and makes individual assets loadable by path. `LoadScene` loads all referenced assets (sheets, map, animations) and assembles them into a runtime-ready scene. `scene.Draw(camera)` renders the visible portion of the scene into the currently active frame.

The caller still owns the outer frame lifecycle. `scene.Draw` draws into the active render pass but does not acquire or present the framebuffer.

## Phase-1 constraints

This first version supports:

- Fixed-grid tiles only (equal-sized tiles sliced from the atlas)
- No arbitrary atlas rectangles
- No mipmaps
- No per-tile transform metadata

## Inspecting loaded assets

```go
m := scene.Map()
fmt.Printf("map: %dx%d tiles, %d layers\n", m.Width(), m.Height(), m.LayerCount())

stats := scene.Stats()
fmt.Printf("visible: %d, uploads: %d\n", stats.VisibleTiles, stats.UploadCount)
```

## Reference examples

See `examples/simplegame` for a minimal working example that follows this tutorial step by step: source PNG + JSON assets, `go generate` pipeline, bundle loading, camera scrolling with D-pad input, and runtime stats overlay.

For a more advanced example with multiple layers, overlay sheets, and tile animations, see `examples/tilemap`.
