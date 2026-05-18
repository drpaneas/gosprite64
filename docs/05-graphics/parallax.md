# Parallax Scrolling

Parallax scrolling creates an illusion of depth by moving background layers at
different speeds relative to the camera. Layers closer to the "viewer" scroll
faster, while distant layers scroll slower.

GoSprite64 provides a lightweight parallax system that computes layer offsets
from camera position. You handle the actual drawing - the library just does
the math.

## Core Types

### ParallaxLayer

A single layer with independent horizontal and vertical scroll speeds:

```go
type ParallaxLayer struct {
    SpeedX float32
    SpeedY float32
}
```

Speed values are multipliers applied to the camera position:

| Speed | Effect |
|---|---|
| `0.0` | Layer is fixed (does not scroll) |
| `0.5` | Scrolls at half the camera speed (distant background) |
| `1.0` | Scrolls at camera speed (same as the main game layer) |
| `> 1.0` | Scrolls faster than the camera (foreground elements) |

Each layer has an `Offset` method that converts a camera position into the
layer's scroll offset:

```go
func (p ParallaxLayer) Offset(cameraX, cameraY int) (int, int)
```

### ParallaxConfig

A collection of layers:

```go
type ParallaxConfig struct {
    Layers []ParallaxLayer
}
```

## Creating a Parallax Config

```go
func NewParallaxConfig(speeds ...ParallaxLayer) ParallaxConfig
```

Pass layers from back (slowest) to front (fastest):

```go
parallax := gosprite64.NewParallaxConfig(
    gosprite64.ParallaxLayer{SpeedX: 0.0, SpeedY: 0.0},  // sky (fixed)
    gosprite64.ParallaxLayer{SpeedX: 0.2, SpeedY: 0.0},  // distant mountains
    gosprite64.ParallaxLayer{SpeedX: 0.5, SpeedY: 0.0},  // hills
    gosprite64.ParallaxLayer{SpeedX: 1.0, SpeedY: 1.0},  // main game layer
)
```

## Getting Layer Offsets

```go
func (pc ParallaxConfig) LayerOffset(layer, cameraX, cameraY int) (int, int)
```

Returns the scroll offset for a specific layer given the current camera
position. If the layer index is out of range, it returns the raw camera
position as a safe fallback.

```go
offsetX, offsetY := parallax.LayerOffset(1, camera.X, camera.Y)
```

You can also call `Offset` directly on a layer:

```go
ox, oy := parallax.Layers[1].Offset(camera.X, camera.Y)
```

## Drawing with Parallax

The parallax system only computes offsets - you choose how to draw each layer.
A common pattern is to draw tiled background images offset by the parallax
values:

```go
func (g *Game) Draw() {
    gosprite64.ClearScreen()

    // Layer 0: fixed sky background
    ox0, oy0 := g.parallax.LayerOffset(0, g.camera.X, g.camera.Y)
    gosprite64.DrawImage(g.skyImage, -ox0, -oy0)

    // Layer 1: distant mountains (slow scroll)
    ox1, oy1 := g.parallax.LayerOffset(1, g.camera.X, g.camera.Y)
    gosprite64.DrawImage(g.mountainImage, -ox1, -oy1)

    // Layer 2: hills (medium scroll)
    ox2, oy2 := g.parallax.LayerOffset(2, g.camera.X, g.camera.Y)
    gosprite64.DrawImage(g.hillImage, -ox2, -oy2)

    // Layer 3: main game world (1:1 with camera)
    g.scene.Draw(g.camera)
}
```

The offset is subtracted from the draw position because the camera position
represents how far the view has scrolled into the world - the background needs
to move in the opposite direction.

## Complete Example

```go
type Game struct {
    camera   *gosprite64.Camera
    parallax gosprite64.ParallaxConfig
    scene    *gosprite64.Scene
    bgFar    image.Image
    bgNear   image.Image
}

func (g *Game) Init() {
    g.camera = &gosprite64.Camera{Width: 288, Height: 216}

    g.parallax = gosprite64.NewParallaxConfig(
        gosprite64.ParallaxLayer{SpeedX: 0.1, SpeedY: 0.0},  // far clouds
        gosprite64.ParallaxLayer{SpeedX: 0.4, SpeedY: 0.0},  // near trees
        gosprite64.ParallaxLayer{SpeedX: 1.0, SpeedY: 1.0},  // game world
    )

    // ... load scene, images ...
}

func (g *Game) Update() {
    // scroll the camera
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadRight) {
        g.camera.X += 2
    }
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft) {
        g.camera.X -= 2
    }
}

func (g *Game) Draw() {
    gosprite64.ClearScreenWith(gosprite64.DarkBlue)

    // Draw background layers with parallax offsets
    ox, _ := g.parallax.LayerOffset(0, g.camera.X, g.camera.Y)
    gosprite64.DrawImage(g.bgFar, -ox, 0)

    ox, _ = g.parallax.LayerOffset(1, g.camera.X, g.camera.Y)
    gosprite64.DrawImage(g.bgNear, -ox, 40)

    // Draw the main game layer
    g.scene.Draw(g.camera)
}
```

## Tips

- Use `SpeedY: 0.0` for side-scrolling games where backgrounds only move
  horizontally.
- For seamless scrolling, make your background images wider than 288 pixels
  and tile them.
- The layer ordering is up to you. Draw the slowest (most distant) layers
  first and the fastest (closest) layers last.
- You can change layer speeds at runtime by modifying the `Layers` slice
  directly, for example to speed up scrolling during a boost effect.

## Try It

> **Download the ROM:** [`parallax_demo.z64`](../emulator/roms/parallax_demo.z64) - Open in [ares](https://ares-emu.net/) with the Expansion Pak enabled.
>
> **Controls:** D-Pad = movement, A = action, B = back, Start = pause, Z = trigger

## Reference Example

See `examples/parallax_demo` in the GoSprite64 repository for a working multi-layer parallax scrolling example.
