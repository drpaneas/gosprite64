# Sprite System Design

## Goal

Add a complete sprite rendering and animation playback layer to GoSprite64 so a user can load sprite sheets, draw individual sprites with transforms (flip, scale, rotation, blending), and play frame-based animation clips with basic playback controls.

The intended outcome is not a gameplay entity framework. The intended outcome is rendering and clip playback primitives that users compose into their own game logic.

## Context

The current GoSprite64 API provides a strong tile pipeline (author, compile, bundle, load, render) and audio pipeline, but the `Sprite` type is an empty struct with only `X` and `Y` fields. There is no sprite rendering, no transforms, no animation playback. Any game that needs animated characters, enemies, items, or effects has no library support.

This is the single biggest capability gap between GoSprite64 and what a retro 2D game requires.

## Decision

GoSprite64 will treat this as one explicit feature:

**Sprite System**

The feature provides:
- sprite sheets compiled from PNG atlases
- sprite drawing with position, flip, scale, rotation, and supported blend paths
- an animation player with clip playback controls
- camera-aware world-space sprite drawing

The primary design rule is:

**The library provides rendering and clip playback primitives. Users build state machines, character controllers, and gameplay orchestration on top.**

## Product Boundary

### What the feature includes

- Sprite sheets compiled from PNG atlases with configurable frame dimensions (fixed grid, matching the tile system)
- Individual sprite drawing with position (float32), horizontal flip, vertical flip, scale, rotation, and supported masked/blended rendering paths
- An animation player type with playback controls: play, pause, resume, stop, loop, restart, set clip
- The player tracks current frame and elapsed ticks, but does not own transition logic between clips
- Camera-aware sprite drawing (world-space sprites that respect the camera viewport)
- RDP-accelerated rendering for the common sprite operations, including draw, flip, scale, rotation, and supported blend paths

### What the feature does not include

- Animation state machines or transition graphs between clips
- Character controller architecture
- Entity/component systems
- Sprite batching or automatic draw-order sorting
- Particle systems
- Sprite-to-sprite collision
- Variable-sized frames (fixed grid only)

Users build state machines, character controllers, and transition logic on top of the playback controls.

### Performance target (separate from scope)

The rendering path should sustain this representative workload without dropping below the target frame rate on the N64 at 288x216:

| Parameter | Value |
|-----------|-------|
| Sprite size | 16x16 pixels |
| Total visible sprites | 32 |
| Of which flipped (H or V) | 20 |
| Of which scaled | 8 |
| Of which rotated | 4 |
| Of which using blended rendering | 4 (non-overlapping) |
| Background | Full-screen tilemap scene (288x216 viewport, 8x8 tiles) |
| Target frame rate | 60 FPS |

This is an engineering validation target, not a public API promise. The actual sprite budget depends on sprite size, transform mix, overlap, tilemap complexity, and what else the game is doing per frame. Overlapping blended sprites carry higher per-pixel cost than non-overlapping ones.

## Architecture

The sprite system is a three-layer design that mirrors the existing tile pipeline pattern.

### Layer 1: Offline asset compilation

A sprite sheet is compiled from a PNG atlas using the existing `mk2dsheet` tool. No new authoring tool is needed for phase 1. The compiled sheet format already stores fixed-grid pixel data with frame dimensions. A 16x16 character sheet is a PNG with 16x16 frame dimensions instead of 8x8 tile dimensions.

The distinction between tile sheet and sprite sheet is a usage convention, not a format difference in phase 1.

### Layer 2: Runtime sprite types

Three public types.

**`SpriteSheet`** is a sprite-oriented wrapper over the compiled sheet data. Internally it may share the same backing implementation as `Sheet`, but publicly it is a separate type with sprite-facing vocabulary (`Frame` instead of `Tile`, frame count instead of tile count). This gives room for sprite-specific metadata to diverge from tile-specific metadata in the future without breaking the public API.

The design principle is: share implementation, not user vocabulary.

**`DrawSpriteOptions`** is a struct passed to the draw call to control transforms and blending:

```go
type DrawSpriteOptions struct {
    FlipH    bool
    FlipV    bool
    ScaleX   float32
    ScaleY   float32
    Rotation float32
    OriginX  float32
    OriginY  float32
    Blend    BlendMode
    Alpha    float32
}
```

Zero-value rule: `ScaleX=0`, `ScaleY=0`, and `Alpha=0` all mean "use default 1.0". True zero scale or zero alpha is not supported. Callers who want an invisible sprite should skip the draw call. Negative scale values are not supported; use `FlipH`/`FlipV` for mirroring.

Origin semantics: `OriginX` and `OriginY` are in frame-local pixel coordinates. `(0, 0)` is the top-left corner of the frame. The origin defines the point that sits at the draw position and the pivot for rotation and scale.

**`AnimationPlayer`** tracks playback state for one clip. It does not draw anything. It only answers what frame should be shown right now.

Playback clock model: one tick equals one fixed-step update at the game's target frame rate (60 Hz by default). A clip's `FPS` field is interpreted against this 60 Hz timebase. A clip with `FPS: 12` advances one animation frame every 5 ticks. A clip with `FPS: 30` advances every 2 ticks.

The caller controls timing:
- normal gameplay: `player.Advance(1)` per update
- dropped frame catch-up: `player.Advance(2)`
- pause: do not call `Advance`
- slow motion: `player.Advance(1)` every other update

PAL/NTSC timing differences are the caller's responsibility.

### Layer 3: Rendering bridge

Drawing goes through the existing internal RDP bridge path, extended to support sprite transforms.

Non-rotated sprites use the fast rectangle-oriented path. This covers draw, flip, scale, and masked/blended rendering. This is the common case and the cheapest path.

Rotated sprites use an internal transformed-quad path. Both paths remain hardware-accelerated. The rotated path carries higher per-sprite cost but is not a software fallback.

Blend mode affects which internal render path is valid. `BlendNone` can use the fastest copy-style path. `BlendMasked` and `BlendAlpha` require the standard blended render path, which is approximately 4x slower per pixel than copy. This is an inherent hardware trade-off, not a library limitation.

Overlap cost: overlapping blended sprites cost more per pixel because each pixel goes through the blend pipeline multiple times. This is a performance characteristic, not an API concern, but it should be documented in performance notes.

Culling: if the implementation culls off-screen sprites, it must use conservative bounds that account for scale and rotation, not just the untransformed frame origin.

## Public API Contract

### Drawing

```go
func DrawSprite(sheet *SpriteSheet, frame int, x, y float32)
func DrawSpriteWithOptions(sheet *SpriteSheet, frame int, x, y float32, opts DrawSpriteOptions)
func DrawWorldSprite(sheet *SpriteSheet, frame int, worldX, worldY float32, cam *Camera)
func DrawWorldSpriteWithOptions(sheet *SpriteSheet, frame int, worldX, worldY float32, cam *Camera, opts DrawSpriteOptions)
```

Positions are `float32` for subpixel placement. The renderer may snap to integer framebuffer coordinates internally for fast paths, but the public API supports smooth movement.

Edge case behavior:
- nil sheet: no-op
- out-of-range frame index: no-op
- nil camera in world-draw variants: draws at world coordinates as if the camera is at (0,0)

### SpriteSheet

```go
func LoadSpriteSheet(path string) (*SpriteSheet, error)

func (s *SpriteSheet) FrameCount() int
func (s *SpriteSheet) FrameWidth() int
func (s *SpriteSheet) FrameHeight() int
```

`LoadSpriteSheet` loads a compiled `.sheet` file by path through the registered asset filesystem. Frame accessors use zero-based indexing.

The sprite sheet is an opaque asset handle, not a decoded CPU-side pixel view. No `Frame(index) image.Image` method. The renderer bridge resolves frame data internally.

### AnimationPlayer

```go
func NewAnimationPlayer() *AnimationPlayer

func (p *AnimationPlayer) Play(clip AnimationClip)
func (p *AnimationPlayer) Pause()
func (p *AnimationPlayer) Resume()
func (p *AnimationPlayer) Stop()
func (p *AnimationPlayer) SetLoop(loop bool)
func (p *AnimationPlayer) Restart()
func (p *AnimationPlayer) Advance(ticks int)
func (p *AnimationPlayer) Frame() int
func (p *AnimationPlayer) Playing() bool
func (p *AnimationPlayer) Done() bool
```

The player consumes `AnimationClip` values from the existing animation system.

`Frame()` returns the current frame index. When no clip is playing or the player is stopped, `Frame()` returns 0. This is a safe numeric fallback, not a meaningful animation choice. Callers should gate on `Playing()` before passing `Frame()` to a draw call to avoid accidentally drawing frame 0 when nothing is playing.

Nil player methods do not panic and return safe zero values.

### BlendMode

```go
type BlendMode uint8

const (
    BlendNone   BlendMode = iota
    BlendMasked
    BlendAlpha
)
```

- `BlendNone`: opaque fast path (copy-style, fastest)
- `BlendMasked`: 1-bit alpha cutout, zero-alpha pixels discarded
- `BlendAlpha`: per-pixel alpha blending from the source texture, with `Alpha` acting as a global opacity multiplier applied on top of the texture's own alpha

## Error Handling

### Load-time errors

`LoadSpriteSheet` returns an explicit error when:
- the path cannot be opened through the registered asset filesystem
- the compiled sheet data is malformed
- the sheet has zero frames

These follow the same strict-and-early pattern as `OpenBundle` and `LoadScene`.

### Frame-time behavior

Drawing functions are defensive and never panic:
- nil sheet: no-op
- out-of-range frame index: no-op
- nil camera in world-draw variants: treat as camera at origin

`AnimationPlayer` is safe on nil receiver and with unset clips. `Frame()` returns 0 if no clip is playing. `Done()` returns true if stopped or no clip is set.

## Testing Strategy

### Format and loading tests

- `LoadSpriteSheet` success with a valid compiled sheet
- `LoadSpriteSheet` failure with malformed data
- `FrameCount`, `FrameWidth`, `FrameHeight` return correct metadata

### AnimationPlayer tests (host-testable)

- `Play` + `Advance` + `Frame` produces correct frame sequence
- FPS-to-tick conversion: clip at FPS 12 advances frame every 5 ticks at 60 Hz base
- `Pause` / `Resume` freezes and unfreezes advancement
- `Stop` resets to frame 0 and `Done()` returns true
- `SetLoop(true)` wraps around; `SetLoop(false)` stops at last frame
- `Restart` replays from frame 0
- `Advance(0)` is a no-op
- `Advance(n)` with large n skips frames correctly
- nil player methods do not panic

### Structural guardrails

- `SpriteSheet` does not expose `image.Image` in its public API
- Draw functions exist with the documented signatures
- `DrawSpriteOptions` zero-value produces default behavior

### Vertical-slice proof

One example that demonstrates loading a sprite sheet, drawing sprites with flip/scale/blend, playing animation clips, and camera-aware world-space drawing. The example must generate, build with the embedded toolchain, and run.

### Performance validation (separate from functional tests)

A benchmark workload matching the performance target measured on N64 hardware or the most representative emulator available.

## Example Story

`examples/sprite_demo` should serve as the canonical proof of the sprite system.

The example should demonstrate:
- a character sprite sheet with at least 4 animation frames
- idle and walk animation clips
- the animation player advancing each frame
- D-pad movement with horizontal flip when changing direction
- camera following the character across a tilemap background
- at least one scaled sprite
- at least one blended sprite
- draw order controlled by call sequence

The example exists to prove the full sprite pipeline works end to end, not to be a complete game.

## Documentation

Documentation should teach:
- how to prepare a sprite sheet PNG
- how to compile it with `mk2dsheet`
- how to load it at runtime with `LoadSpriteSheet`
- how to draw sprites with and without options
- how to use the animation player for clip playback
- the cost model for different blend modes and transform types
- how `examples/sprite_demo` proves the feature

Documentation should not teach:
- animation state machine design
- character controller patterns
- entity management strategies

## Rollout Order

1. Freeze the public API contract and document the zero-value and edge-case rules
2. Implement `SpriteSheet` loading and metadata accessors, with tests
3. Implement `DrawSprite` and `DrawSpriteWithOptions` for the non-rotated fast path, with tests
4. Get `examples/sprite_demo` generating, loading, and drawing static sprites
5. Implement `AnimationPlayer` with full playback controls, with host-testable tests
6. Implement rotation support through the transformed-quad path
7. Implement blend mode support
8. Extend the example to cover animation, flip, scale, rotation, and blending
9. Add structural guardrails and performance validation
10. Write documentation

Each step should include its own tests. Do not defer testing to a final step.

## Non-Goals For This Phase

Do not expand this feature into:
- an entity/component system
- automatic sprite batching or draw-order sorting
- a particle system
- sprite-to-sprite collision detection
- animation state machines or clip transition logic
- variable-sized sprite frames

The phase succeeds when one sprite pipeline is correct and teachable, not when every future gameplay system is anticipated.

## Final Position

The right way to treat the sprite system is as a rendering and clip playback primitive layer: load compiled sprite sheets, draw sprites with transforms through hardware-accelerated paths, and play animation clips with explicit caller-controlled timing. The library provides the drawing and playback tools. Users provide the gameplay logic.
