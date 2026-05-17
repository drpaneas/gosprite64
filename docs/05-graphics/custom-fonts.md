# Custom Fonts

The built-in 8x8 font (`DrawText`) is fine for debug output, but most games
want a distinctive look. GoSprite64 lets you build bitmap fonts from sprite
sheets where each frame is a glyph.

## How It Works

1. Create a PNG atlas with all your glyphs laid out in a grid.
2. Write a font spec JSON that maps characters to grid cells.
3. Run `mk2dfont` to produce a `.sheet` file and a generated `.go` file with
   the glyph map.
4. At runtime, load the sheet and create a `Font` with the generated map.

## The Glyph Type

Each character in a font is described by a `Glyph`:

```go
type Glyph struct {
    Frame   int // sprite sheet frame index
    Width   int // visible width in pixels
    Advance int // cursor advance after this glyph
    OffsetX int // horizontal draw offset
    OffsetY int // vertical draw offset
}
```

| Field | Purpose |
|---|---|
| `Frame` | Which frame in the sprite sheet contains this glyph |
| `Width` | The visible pixel width of the character |
| `Advance` | How many pixels the cursor moves right after drawing |
| `OffsetX` | Shifts the glyph left/right relative to the cursor |
| `OffsetY` | Shifts the glyph up/down relative to the baseline |

For fixed-width fonts, `Width` and `Advance` are typically the same and offsets
are zero.

## Creating a Font

```go
func NewFont(sheet *SpriteSheet, glyphs map[rune]Glyph, lineHeight int) *Font
```

| Parameter | Description |
|---|---|
| `sheet` | A loaded `SpriteSheet` containing the glyph atlas |
| `glyphs` | A map from rune to `Glyph` describing each character |
| `lineHeight` | Pixel distance between lines for multiline text |

The returned `Font` has a default `Spacing` of 2 pixels between lines and no
fallback character.

```go
sheet, _ := gosprite64.LoadSpriteSheet("assets/myfont.sheet")

font := gosprite64.NewFont(sheet, MyFontGlyphs, MyFontLineHeight)
font.Fallback = '?' // show '?' for unknown characters
```

## Font Spec JSON

The `mk2dfont` tool reads a JSON spec that defines your glyph layout. There
are two modes:

### Fixed-Width (simple)

List all characters in order with the `chars` field. Every glyph uses the
full cell dimensions:

```json
{
  "cell_width": 8,
  "cell_height": 10,
  "chars": "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 .!?"
}
```

### Variable-Width (per-glyph)

Use the `glyphs` array to specify per-character widths and offsets:

```json
{
  "cell_width": 16,
  "cell_height": 16,
  "glyphs": [
    { "char": "A", "width": 10 },
    { "char": "B", "width": 9 },
    { "char": "i", "width": 4, "offset_x": 2 },
    { "char": " ", "width": 6 }
  ]
}
```

You must use either `chars` or `glyphs`, not both.

## Building with mk2dfont

```bash
go run github.com/drpaneas/gosprite64/cmd/mk2dfont \
    -png font.png \
    -spec font.json \
    -out-sheet assets/font.sheet \
    -out-go glyphs.go \
    -name myFont \
    -pkg main
```

| Flag | Description |
|---|---|
| `-png` | Path to the font atlas PNG |
| `-spec` | Path to the font spec JSON |
| `-out-sheet` | Output path for the compiled `.sheet` file |
| `-out-go` | Output path for the generated Go source file |
| `-name` | Font name used in generated variable names |
| `-pkg` | Package name for the generated file |

This produces two files:
- A `.sheet` binary (same format as `mk2dsheet`)
- A `.go` file with a glyph map and line height constant

The generated Go file looks like:

```go
package main

import "github.com/drpaneas/gosprite64"

const MyFontLineHeight = 16

var MyFontGlyphs = map[rune]gosprite64.Glyph{
    'A': gosprite64.Glyph{Frame: 0, Width: 10, Advance: 10},
    'B': gosprite64.Glyph{Frame: 1, Width: 9, Advance: 9},
    'i': gosprite64.Glyph{Frame: 2, Width: 4, Advance: 4, OffsetX: 2},
    ' ': gosprite64.Glyph{Frame: 3, Width: 6, Advance: 6},
}
```

## Drawing Text

### DrawTextEx

```go
func (f *Font) DrawTextEx(text string, x, y int, align TextAlign)
```

Renders text at `(x, y)` using the font's sprite sheet. Supports multiline
strings (split on `\n`) and alignment via the `TextAlign` type. See
[Text Alignment](text-alignment.md) for details.

```go
font.DrawTextEx("GAME OVER", 144, 100, gosprite64.AlignCenter)
```

### GlyphFor

```go
func (f *Font) GlyphFor(r rune) (Glyph, bool)
```

Looks up the glyph for a rune. If the rune is not in the font and `Fallback`
is set, returns the fallback glyph instead.

## Measuring Text

```go
func (f *Font) MeasureText(text string) (width int, height int)
```

Returns the pixel dimensions of the rendered text without drawing anything.
Supports multiline strings. Useful for centering or positioning UI elements.

```go
w, h := font.MeasureText("SCORE: 0042")
// Position text so it's centered on screen
x := (288 - w) / 2
y := (216 - h) / 2
font.DrawTextEx("SCORE: 0042", x, y, gosprite64.AlignLeft)
```

## Formatting Scores

```go
func FormatScore(score int, width int) string
```

Formats an integer with leading zeros to a fixed width. Negative values are
clamped to 0. If the number has more digits than `width`, the full number is
returned.

```go
gosprite64.FormatScore(42, 6)   // "000042"
gosprite64.FormatScore(0, 4)    // "0000"
gosprite64.FormatScore(99999, 3) // "99999" (not truncated)
```

This pairs well with custom fonts for score displays:

```go
text := gosprite64.FormatScore(g.score, 6)
font.DrawTextEx(text, 240, 4, gosprite64.AlignRight)
```

## Complete Example

```go
type Game struct {
    font  *gosprite64.Font
    score int
}

func (g *Game) Init() {
    sheet, err := gosprite64.LoadSpriteSheet("assets/font.sheet")
    if err != nil {
        panic(err)
    }
    g.font = gosprite64.NewFont(sheet, MyFontGlyphs, MyFontLineHeight)
    g.font.Fallback = '?'
}

func (g *Game) Update() {
    g.score++
}

func (g *Game) Draw() {
    gosprite64.ClearScreen()

    // Title centered at top
    g.font.DrawTextEx("MY COOL GAME", 144, 10, gosprite64.AlignCenter)

    // Score right-aligned
    scoreText := gosprite64.FormatScore(g.score, 6)
    g.font.DrawTextEx(scoreText, 280, 4, gosprite64.AlignRight)

    // Multiline message
    g.font.DrawTextEx("PRESS START\nTO BEGIN", 144, 100, gosprite64.AlignCenter)
}
```
