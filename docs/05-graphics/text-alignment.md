# Text Alignment

The `TextAlign` type controls horizontal positioning when drawing text with a
custom `Font`.

## The TextAlign Type

```go
type TextAlign int

const (
    AlignLeft   TextAlign = iota // 0 - default
    AlignCenter                  // 1
    AlignRight                   // 2
)
```

You pass a `TextAlign` value to `Font.DrawTextEx`:

```go
func (f *Font) DrawTextEx(text string, x, y int, align TextAlign)
```

## How Alignment Works

The `x` coordinate acts as the alignment anchor. The text is positioned
relative to that anchor depending on the alignment mode:

### AlignLeft

Text starts at `x` and extends to the right. This is the default.

```go
font.DrawTextEx("HELLO", 10, 50, gosprite64.AlignLeft)
//  x=10
//  |
//  HELLO
```

### AlignCenter

Text is centered around `x`. More precisely, the widest line determines the
total width, and each line is centered within that width starting from `x`.

```go
font.DrawTextEx("GAME OVER", 144, 100, gosprite64.AlignCenter)
//       x=144
//         |
//    GAME OVER
```

### AlignRight

Text ends at the right edge of the measured width, starting from `x`.

```go
font.DrawTextEx("99999", 280, 4, gosprite64.AlignRight)
//              x=280
//                  |
//              99999
```

## Multiline Alignment

For multiline strings (containing `\n`), each line is aligned independently
within the bounding box of the widest line:

```go
font.DrawTextEx("GAME OVER\nPRESS START", 144, 80, gosprite64.AlignCenter)
//       x=144
//         |
//    GAME OVER
//   PRESS START
```

The total width is measured from the widest line. Shorter lines are shifted
within that width according to the alignment.

## Word Wrapping

The `Font.WrapText` method inserts newlines so no line exceeds a given pixel
width. It breaks on spaces and does not break words that are wider than the
limit.

```go
func (f *Font) WrapText(text string, maxWidth int) string
```

Combine it with alignment for paragraph-style text:

```go
wrapped := font.WrapText("This is a long message that should wrap nicely on screen.", 200)
font.DrawTextEx(wrapped, 144, 50, gosprite64.AlignCenter)
```

## Practical Examples

### Centered Title Screen

```go
func (g *Game) Draw() {
    gosprite64.ClearScreenWith(gosprite64.DarkBlue)
    g.font.DrawTextEx("MY GAME", 144, 60, gosprite64.AlignCenter)
    g.font.DrawTextEx("PRESS START", 144, 140, gosprite64.AlignCenter)
}
```

### Right-Aligned Score

```go
scoreText := gosprite64.FormatScore(g.score, 6)
g.font.DrawTextEx(scoreText, 280, 4, gosprite64.AlignRight)
```

### Left-Aligned Dialog

```go
dialog := g.font.WrapText(npcMessage, 200)
g.font.DrawTextEx(dialog, 20, 160, gosprite64.AlignLeft)
```
