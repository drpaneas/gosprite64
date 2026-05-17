# Sprite Sheets

A sprite sheet is a single PNG image sliced into a grid of equal-sized frames.
GoSprite64 uses a compile-time tool (`mk2dsheet`) to convert the PNG into an
optimized binary `.sheet` file that can be loaded at runtime.

## Preparing Your PNG

Lay out your frames in a regular grid. Every cell must be the same width and
height. The tool reads the image left-to-right, top-to-bottom, assigning
ascending frame indices starting at 0.

```
 Frame layout for a 64x32 PNG with 16x16 tiles:
 ┌────┬────┬────┬────┐
 │  0 │  1 │  2 │  3 │   row 0
 ├────┼────┼────┼────┤
 │  4 │  5 │  6 │  7 │   row 1
 └────┴────┴────┴────┘
```

The total number of frames is `(imageWidth / tileWidth) * (imageHeight / tileHeight)`.

## Compiling with mk2dsheet

The `mk2dsheet` command converts a PNG into a `.sheet` binary:

```bash
go run github.com/drpaneas/gosprite64/cmd/mk2dsheet \
    -in character.png \
    -out character.sheet \
    -tile-width 16 \
    -tile-height 16
```

| Flag | Description |
|---|---|
| `-in` | Path to the source PNG |
| `-out` | Path for the output `.sheet` file |
| `-tile-width` | Width of each frame in pixels (default 8) |
| `-tile-height` | Height of each frame in pixels (default 8) |

A typical project runs this as a `go:generate` directive so the asset build
stays reproducible:

```go
//go:generate sh -c "go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/character.png -out assets/character.sheet -tile-width 16 -tile-height 16"
```

## Loading at Runtime

```go
func LoadSpriteSheet(path string) (*SpriteSheet, error)
```

`LoadSpriteSheet` reads a `.sheet` file from the cartridge filesystem and
returns a ready-to-use `SpriteSheet`. It returns an error if the file cannot
be read or contains zero frames.

```go
sheet, err := gosprite64.LoadSpriteSheet("assets/character.sheet")
if err != nil {
    panic(err)
}
```

## Querying the Sheet

Once loaded, you can inspect the sheet's properties:

```go
func (s *SpriteSheet) FrameCount() int
func (s *SpriteSheet) FrameWidth() int
func (s *SpriteSheet) FrameHeight() int
```

| Method | Returns |
|---|---|
| `FrameCount()` | Total number of frames in the sheet |
| `FrameWidth()` | Width of a single frame in pixels |
| `FrameHeight()` | Height of a single frame in pixels |

```go
fmt.Printf("Loaded %d frames, each %dx%d\n",
    sheet.FrameCount(),
    sheet.FrameWidth(),
    sheet.FrameHeight(),
)
// Output: Loaded 8 frames, each 16x16
```

All methods are nil-safe - calling them on a `nil` `SpriteSheet` returns 0.

## Drawing Sprites from a Sheet

After loading, use `DrawSprite` or `DrawSpriteWithOptions` to draw individual
frames. Frame indices are 0-based.

```go
// Draw frame 0 at screen position (100, 80)
gosprite64.DrawSprite(sheet, 0, 100, 80)

// Draw frame 3 flipped horizontally
gosprite64.DrawSpriteWithOptions(sheet, 3, 100, 80, gosprite64.DrawSpriteOptions{
    FlipH: true,
})
```

For world-space drawing with camera offset:

```go
gosprite64.DrawWorldSprite(sheet, frame, worldX, worldY, camera)
```

See [Sprites](sprites.md) for the full drawing API including scaling,
rotation, and alpha blending.

## Complete Example

```go
//go:generate sh -c "go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/items.png -out assets/items.sheet -tile-width 16 -tile-height 16"

package main

import "github.com/drpaneas/gosprite64"

type Game struct {
    items *gosprite64.SpriteSheet
}

func (g *Game) Init() {
    sheet, err := gosprite64.LoadSpriteSheet("assets/items.sheet")
    if err != nil {
        panic(err)
    }
    g.items = sheet
}

func (g *Game) Update() {}

func (g *Game) Draw() {
    gosprite64.ClearScreen()

    // Draw all item frames in a row
    for i := 0; i < g.items.FrameCount(); i++ {
        x := float32(10 + i*20)
        gosprite64.DrawSprite(g.items, i, x, 100)
    }
}

func main() {
    gosprite64.Run(&Game{})
}
```
