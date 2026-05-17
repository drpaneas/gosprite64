# Performance Notes

Tips for keeping your game running at 60 FPS on real N64 hardware.

## The 288x216 Canvas

GoSprite64 renders to a fixed **288x216 logical canvas** with 16-pixel gutters on each side. This canvas is scaled up to the full TV resolution (640x480 NTSC or 640x576 PAL) by the Video Interface. All your drawing coordinates use this logical space.

The 16-pixel gutters exist because CRT televisions overscan the edges of the picture. The visible safe area is the inner 256x184 region, but you can draw into the full 288x216 if you want edge-to-edge coverage on modern displays.

## Fixed-Step Timing

The game loop runs at a fixed **60 FPS** with a fixed-step accumulator. Each `Update()` call represents exactly one 1/60th-second tick. If the hardware falls behind, `Update()` runs multiple times to catch up before `Draw()` is called. This means your physics and input logic always see consistent time steps regardless of rendering performance.

## Blend Modes and Fill Rate

The blend mode you choose on sprites has a major impact on rendering cost:

- **`BlendNone`** is the fastest path. It performs a simple opaque blit with no per-pixel blending. Use this for backgrounds, tiles, and any sprites without transparency.
- **`BlendMasked`** treats each pixel as either fully opaque or fully transparent based on the alpha channel. It is roughly **4x slower** than `BlendNone` because the RDP must test each pixel individually.
- **`BlendAlpha`** performs full per-pixel alpha blending. It has similar cost to `BlendMasked` for a single sprite, but **overlapping blended sprites** multiply the cost because each overlapping pixel must be blended again.

Practical advice:

- Default to `BlendNone` for tiles and backgrounds.
- Use `BlendMasked` for characters and objects with hard-edged transparency.
- Reserve `BlendAlpha` for effects that genuinely need smooth transparency (particles, shadows, UI overlays).
- Avoid stacking many alpha-blended sprites in the same screen region.

## Rotation and Scaling

When you set `Rotation` to a non-zero value in `DrawSpriteOptions`, the sprite is drawn through the **transformed-quad path** instead of the fast axis-aligned blit. This path computes per-pixel source coordinates through an affine transform, which costs more per pixel than a straight copy.

Scaling alone (without rotation) is cheaper than rotation but still more expensive than an unscaled blit. If a sprite's `ScaleX`, `ScaleY`, and `Rotation` are all at their defaults, GoSprite64 automatically falls back to the fast `DrawSprite` path.

Tips:

- Pre-render rotated frames in your sprite sheet if you only need a few fixed angles (e.g., 4 or 8 directions).
- Keep the number of simultaneously rotating sprites small, especially at larger sizes.

## Scene Rendering

Tile scenes are optimized for the N64's constrained memory and fill rate:

- **Only visible cells are rendered.** The renderer uses the camera viewport to determine which tiles fall within the screen bounds and skips everything outside. Scrolling a large map costs the same as rendering a small map, as long as the visible tile count stays similar.
- **`LayerInfo()` is O(1).** Layer metadata (sheet ID, non-zero tile count) is cached at load time, so querying layer properties in your update loop is free.
- **`Stats()` is allocation-free.** The `RuntimeStats` struct is returned by value from pre-computed fields, so you can call it every frame without triggering garbage collection.

## Audio Performance

The audio engine is designed for zero allocations after initialization:

- All buffers (output, per-voice source, byte conversion) are pre-allocated once during `initAudio`.
- The mixing loop reuses fixed-size arrays and never calls `make` or `append`.
- Command dispatch uses a lock-free ring buffer, so `PlaySoundEffect` and `PlayMusic` never block the game loop.

This means audio has no GC pressure during gameplay. You can play and stop sounds freely without worrying about frame hitches.

## Math: Avoid sqrt on N64

The N64's MIPS R4300i CPU does not have a hardware square root instruction. Computing `sqrt` requires a software approximation that takes many more cycles than basic arithmetic.

The `math2d` package provides squared-distance methods on `Vec2` specifically for this:

- **`v.LengthSq() float32`** returns `x*x + y*y` without the square root.
- **`v.DistanceSq(other Vec2) float32`** returns the squared distance between two points.

Use these for distance comparisons. Instead of:

```go
if a.Distance(b) < 32 {
```

Write:

```go
if a.DistanceSq(b) < 32*32 {
```

The result is mathematically equivalent for comparisons and avoids the expensive sqrt entirely. Only call `Distance()` or `Length()` when you actually need the scalar distance value (e.g., for normalization).

## General Tips

- **Minimize draw calls for overlapping transparent sprites.** Each overlapping blended pixel is processed again. Group opaque background layers together and draw them with `BlendNone`.
- **Use smaller sprite sizes when possible.** A 16x16 sprite fills 4x fewer pixels than a 32x32 sprite. The N64's fill rate is the most common bottleneck.
- **Profile with `Stats()`.** Check `VisibleTiles` and `UploadCount` each frame to understand where rendering time is going. A spike in `UploadCount` means the renderer is loading new tile data into TMEM more often than expected.
- **Keep your map layers shallow.** Each additional layer multiplies the number of tiles the renderer must process for every visible cell. Two or three layers is typical; five or more may stress the fill rate.
- **Pre-compute what you can.** `LayerInfo` and `SheetInfo` are cached, so prefer them over re-scanning tile data manually.
