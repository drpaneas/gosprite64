# Drawing Primitives

GoSprite64 provides a small set of drawing functions for shapes, lines, text,
and images. All coordinates use the **288x216 logical canvas** - the library
maps them to the actual framebuffer resolution automatically.

## Clearing the Screen

Every frame typically starts by wiping the previous contents:

```go
func ClearScreen()
func ClearScreenWith(c color.Color)
```

`ClearScreen` fills the entire screen with `gosprite64.Black`.
`ClearScreenWith` lets you pick the color:

```go
func (g *Game) Draw() {
    gosprite64.ClearScreen()            // solid black
    // or
    gosprite64.ClearScreenWith(gosprite64.DarkBlue) // night-sky blue
}
```

## Filled Rectangles

```go
func FillRect(x1, y1, x2, y2 int, c color.Color)
```

Draws a solid rectangle from the top-left corner `(x1, y1)` to the
bottom-right corner `(x2, y2)`, inclusive. The coordinates are automatically
swapped if you pass them in the wrong order.

```go
// A 32x16 red rectangle starting at (10, 10)
gosprite64.FillRect(10, 10, 41, 25, gosprite64.Red)
```

Coordinates are clipped to the screen bounds, so you don't need to worry about
drawing outside the 288x216 area.

## Outlined Rectangles

```go
func DrawRect(x1, y1, x2, y2 int, c color.Color)
```

Draws a 1-pixel outline around the given rectangle. Corners are drawn exactly
once (no overdraw).

```go
// A green border around a 50x50 area
gosprite64.DrawRect(20, 20, 69, 69, gosprite64.Green)
```

## Lines

```go
func DrawLine(x1, y1, x2, y2 int, c color.Color)
```

Draws a 1-pixel line between two points. Horizontal and vertical lines are
optimized into a single `FillRect` call. Diagonal lines use Bresenham's
algorithm.

```go
// Horizontal line across the top
gosprite64.DrawLine(0, 0, 287, 0, gosprite64.White)

// Diagonal line
gosprite64.DrawLine(0, 0, 100, 80, gosprite64.Yellow)
```

## Text

```go
func DrawText(str string, x, y int, c color.Color)
```

Renders a string using the built-in 8x8 monospace bitmap font. Each character
occupies exactly 8 pixels wide. Characters outside the printable ASCII range
(32-127) are skipped but still advance the cursor by 8 pixels.

```go
gosprite64.DrawText("HELLO WORLD", 10, 10, gosprite64.White)

// Display coordinates
gosprite64.DrawText(
    fmt.Sprintf("x:%d y:%d", playerX, playerY),
    2, 2,
    gosprite64.White,
)
```

The built-in font is good for debug overlays and quick prototyping. For
styled text with variable-width characters, see
[Custom Fonts](custom-fonts.md).

## Drawing Images

```go
func DrawImage(src image.Image, x, y int)
```

Blits a Go `image.Image` at the given logical position. The image is clipped
to the 288x216 canvas automatically.

```go
gosprite64.DrawImage(myImage, 50, 30)
```

### World-Space Images

```go
func DrawWorldImage(src image.Image, worldX, worldY int, cam *Camera)
```

Works like `DrawImage` but offsets the position by the camera's scroll. If
`cam` is `nil`, it behaves identically to `DrawImage`.

```go
// Draw a pickup item at world coordinates, offset by the camera
gosprite64.DrawWorldImage(coinImg, coinWorldX, coinWorldY, g.camera)
```

## Coordinate Quick Reference

| Concept | Range |
|---|---|
| Logical canvas width | 0 - 287 |
| Logical canvas height | 0 - 215 |
| Origin | Top-left corner (0, 0) |
| Inclusive corners | Both `(x1,y1)` and `(x2,y2)` are drawn |

All drawing functions accept coordinates in this logical space. The library
handles scaling to the physical framebuffer and clipping to screen bounds.

## Complete Example

```go
func (g *Game) Draw() {
    gosprite64.ClearScreenWith(gosprite64.DarkBlue)

    // Sky gradient (stacked horizontal bars)
    for y := 0; y < 80; y++ {
        gosprite64.DrawLine(0, y, 287, y, gosprite64.DarkPurple)
    }

    // Ground
    gosprite64.FillRect(0, 160, 287, 215, gosprite64.DarkGreen)

    // House outline
    gosprite64.DrawRect(100, 120, 180, 160, gosprite64.White)

    // Door
    gosprite64.FillRect(130, 140, 150, 160, gosprite64.Brown)

    // HUD text
    gosprite64.DrawText("SCORE: 0042", 2, 2, gosprite64.Yellow)
}
```
