# The Fixed Canvas

> If you just finished the beginner journey, this page explains the concept behind the behavior you already saw on screen.

GoSprite64 gives you a single, fixed drawing surface: **288 x 216 logical pixels**. Every drawing function in the library operates in this coordinate space, and the runtime handles everything else.

## The Logical Canvas

All public drawing APIs use the same coordinate system:

- **X** ranges from `0` (left) to `287` (right)
- **Y** ranges from `0` (top) to `215` (bottom)

When you call `FillRect(10, 10, 50, 50, Red)`, those numbers are logical pixel coordinates. You never need to think about the physical framebuffer, video output format, or TV overscan.

```go
func (g *Game) Draw() {
    gosprite64.ClearScreen()

    // These coordinates are always 288x216 logical pixels
    gosprite64.FillRect(0, 0, 287, 215, gosprite64.DarkBlue)   // full screen
    gosprite64.DrawRect(10, 10, 277, 205, gosprite64.White)     // 10px inset border
    gosprite64.DrawText("288x216", 112, 100, gosprite64.Yellow) // centered-ish text
}
```

## What Happens Under the Hood

The N64's actual framebuffer is **320 x 240** pixels. GoSprite64 places your 288x216 canvas inside this framebuffer with a 16-pixel horizontal margin and 12-pixel vertical margin on each side:

```
+--- 320x240 framebuffer ------+
|  16px  +-288x216-+  16px     |
|  margin| logical |  margin   |
| 12px   | canvas  |           |
|        +---------+           |
+-------------------------------+
```

The origin offset is `(16, 12)` - every logical coordinate you pass is translated by this amount before being drawn to the framebuffer. The margin area is not accessible to your game code.

After rendering, the 320x240 framebuffer is scaled up to a 640x480 output image and then centered appropriately for the active TV mode, ensuring that pixels appear square regardless of whether the console is running in NTSC or PAL timing.

## Why 288x216?

The N64 outputs video in formats where pixels are not naturally square. NTSC displays 640x480 with a 4:3 aspect ratio, but the pixel aspect ratio is not 1:1 at 320x240. By using 288x216 as the logical canvas, GoSprite64 ensures that:

- Your game looks the same on NTSC and PAL televisions
- Circles look like circles, squares look like squares
- Sprite art can be authored at 1:1 pixel ratio in any art tool

For a deep dive into the math behind this, see the [Square Pixels](square-pixels.md) chapter.

## Functions That Use the Logical Canvas

Every drawing function in `gosprite64` uses logical coordinates:

```go
gosprite64.FillRect(x1, y1, x2, y2, color)
gosprite64.DrawRect(x1, y1, x2, y2, color)
gosprite64.DrawLine(x1, y1, x2, y2, color)
gosprite64.DrawText(str, x, y, color)
gosprite64.DrawImage(img, x, y)
gosprite64.DrawSprite(sheet, frame, x, y)
```

You can also use `DrawWorldSprite` and `DrawWorldImage` which accept world-space coordinates and a `Camera` - but the camera offset is applied before mapping into the same 288x216 space.

## Draw Regions

If you need to restrict drawing to a sub-area of the screen (for split-screen multiplayer, for example), use `SetDrawRegion`:

```go
gosprite64.SetDrawRegion(0, 0, 144, 216)   // left half
// ... draw player 1 ...
gosprite64.ResetDrawRegion()

gosprite64.SetDrawRegion(144, 0, 144, 216)  // right half
// ... draw player 2 ...
gosprite64.ResetDrawRegion()
```

The coordinates passed to `SetDrawRegion` are also in logical space. See [Draw Regions](../05-graphics/draw-regions.md) for details.
