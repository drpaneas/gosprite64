# Sprites

This chapter covers how to load sprite sheets, draw individual sprites with transforms, animate them, and render them in world space alongside a tile scene.

## Preparing a sprite sheet PNG

A sprite sheet is a single PNG image that contains all frames of a character or object arranged in a grid. Each cell in the grid is one frame. The same `mk2dsheet` tool used for tile sheets handles sprite sheets - you just specify different frame dimensions.

For example, a character sprite sheet with 16x16 pixel frames:

```
+----+----+----+----+
| 0  | 1  | 2  | 3  |   row 0
+----+----+----+----+
| 4  | 5  | 6  | 7  |   row 1
+----+----+----+----+
```

Requirements:

- PNG format
- Image width must be evenly divisible by the frame width
- Image height must be evenly divisible by the frame height
- Pixels are stored as NRGBA internally

## Compiling the sprite sheet

Use `mk2dsheet` with the frame dimensions matching your sprite size:

```bash
go run github.com/drpaneas/gosprite64/cmd/mk2dsheet \
  -in assets-src/character.png \
  -out assets/character.sheet \
  -tile-width 16 -tile-height 16
```

This produces a `.sheet` binary that the runtime can load. For a `go:generate` line:

```go
//go:generate sh -c "mkdir -p assets && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/character.png -out assets/character.sheet -tile-width 16 -tile-height 16"
```

## Loading a sprite sheet

```go
charSheet, err := gosprite64.LoadSpriteSheet("assets/character.sheet")
if err != nil {
    panic(err)
}
```

`LoadSpriteSheet` reads the compiled `.sheet` binary and returns a `*SpriteSheet`. The sheet exposes metadata accessors:

```go
charSheet.FrameCount()  // total number of frames in the sheet
charSheet.FrameWidth()  // width of one frame in pixels (e.g. 16)
charSheet.FrameHeight() // height of one frame in pixels (e.g. 16)
```

## Drawing sprites

### Simple drawing

`DrawSprite` renders a single frame at screen-space coordinates with no transforms:

```go
gosprite64.DrawSprite(charSheet, frameIndex, x, y)
```

The `frame` argument is a zero-based index into the sheet. Out-of-range frames are silently ignored.

### Drawing with options

`DrawSpriteWithOptions` accepts a `DrawSpriteOptions` struct for transforms:

```go
gosprite64.DrawSpriteWithOptions(charSheet, frame, x, y, gosprite64.DrawSpriteOptions{
    FlipH: true,
    Blend: gosprite64.BlendAlpha,
    Alpha: 0.7,
})
```

When all options are at their defaults, this falls through to the fast `DrawSprite` path automatically.

## DrawSpriteOptions reference

| Field      | Type      | Zero value means | Description |
|------------|-----------|-----------------|-------------|
| `FlipH`    | `bool`    | no flip         | Mirror the frame horizontally |
| `FlipV`    | `bool`    | no flip         | Mirror the frame vertically |
| `ScaleX`   | `float32` | 1.0             | Horizontal scale factor. Negative values are not supported |
| `ScaleY`   | `float32` | 1.0             | Vertical scale factor. Negative values are not supported |
| `Rotation` | `float32` | no rotation     | Rotation angle in radians |
| `OriginX`  | `float32` | 0               | X component of the transform pivot in frame-local coordinates |
| `OriginY`  | `float32` | 0               | Y component of the transform pivot in frame-local coordinates |
| `Blend`    | `BlendMode` | `BlendNone`   | Blending mode (see below) |
| `Alpha`    | `float32` | 1.0             | Global alpha multiplier. Only meaningful with `BlendAlpha` |

### Blend modes

Three blend modes are available, ordered from fastest to most expensive:

- **`BlendNone`** - No blending. Every source pixel overwrites the destination. This is the fastest mode, roughly 4x faster than alpha blending.

- **`BlendMasked`** - Binary cutout. Fully transparent pixels (alpha = 0) are skipped, all other pixels are drawn at full opacity. Useful for character sprites over backgrounds.

- **`BlendAlpha`** - Full per-pixel alpha blending with an additional global `Alpha` multiplier. The most expensive mode, but required for transparency effects like shadows, ghosts, or fade-outs.

### Scale and rotation

When `ScaleX` or `ScaleY` is 0, it is treated as 1.0. This lets you use the zero-value `DrawSpriteOptions{}` without accidentally scaling to zero.

`OriginX` and `OriginY` define the pivot point for rotation in frame-local pixel coordinates. For example, to rotate a 16x16 sprite around its center, set `OriginX: 8, OriginY: 8`.

## Animation

GoSprite64 provides a tick-based `AnimationPlayer` that drives frame selection from animation clips.

### Setting up clips

Animation clips are defined in JSON and compiled with `mk2danim` (see the tile2d chapter). Each clip has a name, an FPS rate, and a list of frame indices:

```json
{
  "clips": [
    {"name": "idle", "fps": 12, "frames": [0, 1, 2, 3]},
    {"name": "walk", "fps": 8,  "frames": [4, 5, 6, 7]}
  ]
}
```

At runtime, retrieve clips from a loaded animation set:

```go
animSet, err := bundle.LoadAnimation("anims")
if err != nil {
    panic(err)
}

idleClip, ok := animSet.Clip("idle")
if !ok {
    panic("idle clip not found")
}
```

### Using AnimationPlayer

```go
player := gosprite64.NewAnimationPlayer()
player.SetLoop(true)
player.Play(idleClip)
```

Each frame of your game loop, advance the player by one tick:

```go
func (g *Game) Update() {
    g.player.Advance(1)
}
```

Then use `Frame()` to get the current sprite sheet frame index:

```go
func (g *Game) Draw() {
    frame := g.player.Frame()
    gosprite64.DrawSprite(g.charSS, frame, x, y)
}
```

### The tick model

The player runs at a base rate of 60 ticks per second (matching the N64's 60Hz refresh). `Advance(1)` means "one tick has passed." The clip's FPS is interpreted against this 60-tick base:

- A clip at 60 FPS advances one animation frame per tick
- A clip at 12 FPS advances one animation frame every 5 ticks
- A clip at 30 FPS advances one animation frame every 2 ticks

### Switching clips

Call `Play` with a different clip to switch animations. The player resets to frame 0:

```go
if moving {
    if currentClip != "walk" {
        player.Play(walkClip)
        currentClip = "walk"
    }
} else {
    if currentClip != "idle" {
        player.Play(idleClip)
        currentClip = "idle"
    }
}
```

Guard clip switches with a state check to avoid resetting the animation every frame.

### Frame() returns 0 when stopped

`Frame()` returns 0 when the player is not playing (stopped or no clip loaded). If frame 0 in your sheet is a valid sprite, this can cause unwanted drawing. Gate your draw calls on `Playing()`:

```go
if g.player.Playing() {
    frame := g.player.Frame()
    gosprite64.DrawSprite(g.charSS, frame, x, y)
}
```

### Other controls

```go
player.Pause()          // freeze at current frame
player.Resume()         // continue from paused frame
player.Stop()           // stop and reset to frame 0
player.Restart()        // restart the current clip from the beginning
player.Playing()        // true if currently advancing
player.Done()           // true if stopped (finished or never started)
```

## World-space drawing

When drawing sprites alongside a tile scene that uses a camera, use the world-space variants to automatically offset by the camera position:

```go
gosprite64.DrawWorldSprite(charSheet, frame, worldX, worldY, camera)

gosprite64.DrawWorldSpriteWithOptions(charSheet, frame, worldX, worldY, camera, gosprite64.DrawSpriteOptions{
    FlipH: true,
})
```

These subtract `camera.X` and `camera.Y` from the world coordinates before drawing. If the camera is nil, they fall through to the screen-space versions.

A common pattern is drawing a character sprite over the tile scene with a shadow underneath:

```go
// Draw shadow first (stretched, semi-transparent)
gosprite64.DrawWorldSpriteWithOptions(charSheet, frame, playerX, playerY+12, camera, gosprite64.DrawSpriteOptions{
    ScaleX: 1.5,
    ScaleY: 0.3,
    Blend:  gosprite64.BlendAlpha,
    Alpha:  0.3,
})

// Draw character on top
gosprite64.DrawWorldSpriteWithOptions(charSheet, frame, playerX, playerY, camera, gosprite64.DrawSpriteOptions{
    FlipH: facingLeft,
})
```

## Performance notes

- `BlendNone` is roughly 4x faster than `BlendAlpha`. Use it for opaque sprites that fully cover their footprint.
- `BlendMasked` is between the two - it skips transparent pixels but avoids the per-pixel alpha math.
- Rotation adds per-pixel coordinate transforms. Prefer axis-aligned sprites when performance is tight.
- Overlapping blended sprites compound cost: each overlapping pixel runs the blend math again. Minimize large transparent overlaps in hot scenes.
- The `DrawSpriteWithOptions` fast path kicks in when all options are at defaults, falling through to the plain `DrawSprite` code.

## Reference example

See `examples/sprite_demo` for a complete working example that demonstrates:

- Compiling a 16x16 character sprite sheet alongside an 8x8 tile sheet
- Loading both the tile scene and a standalone sprite sheet
- Animation clips (idle and walk) driven by `AnimationPlayer`
- World-space rendering with camera tracking
- Shadow effect using scaled, alpha-blended sprites
- HUD overlay drawing in screen space
