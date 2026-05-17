# Colors

GoSprite64 includes a built-in 16-color palette inspired by classic fantasy consoles. You can use these named constants directly or create custom colors with Go's standard `image/color` package.

## The Built-in Palette

| Constant | R | G | B | Preview |
|----------|---|---|---|---------|
| `Black` | 0 | 0 | 0 | ![#000000](https://via.placeholder.com/16/000000/000000) |
| `DarkBlue` | 29 | 43 | 83 | ![#1D2B53](https://via.placeholder.com/16/1D2B53/1D2B53) |
| `DarkPurple` | 126 | 37 | 83 | ![#7E2553](https://via.placeholder.com/16/7E2553/7E2553) |
| `DarkGreen` | 0 | 135 | 81 | ![#008751](https://via.placeholder.com/16/008751/008751) |
| `Brown` | 171 | 82 | 54 | ![#AB5236](https://via.placeholder.com/16/AB5236/AB5236) |
| `DarkGray` | 95 | 87 | 79 | ![#5F574F](https://via.placeholder.com/16/5F574F/5F574F) |
| `LightGray` | 194 | 195 | 199 | ![#C2C3C7](https://via.placeholder.com/16/C2C3C7/C2C3C7) |
| `White` | 255 | 241 | 232 | ![#FFF1E8](https://via.placeholder.com/16/FFF1E8/FFF1E8) |
| `Red` | 255 | 0 | 77 | ![#FF004D](https://via.placeholder.com/16/FF004D/FF004D) |
| `Orange` | 255 | 163 | 0 | ![#FFA300](https://via.placeholder.com/16/FFA300/FFA300) |
| `Yellow` | 255 | 236 | 39 | ![#FFEC27](https://via.placeholder.com/16/FFEC27/FFEC27) |
| `Green` | 0 | 228 | 54 | ![#00E436](https://via.placeholder.com/16/00E436/00E436) |
| `Blue` | 41 | 173 | 255 | ![#29ADFF](https://via.placeholder.com/16/29ADFF/29ADFF) |
| `Indigo` | 131 | 118 | 156 | ![#83769C](https://via.placeholder.com/16/83769C/83769C) |
| `Pink` | 255 | 119 | 168 | ![#FF77A8](https://via.placeholder.com/16/FF77A8/FF77A8) |
| `Peach` | 255 | 204 | 170 | ![#FFCCAA](https://via.placeholder.com/16/FFCCAA/FFCCAA) |

All 16 colors are declared as `color.Color` variables in the `gosprite64` package:

```go
var (
    Black      color.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
    DarkBlue   color.Color = color.RGBA{R: 29, G: 43, B: 83, A: 255}
    DarkPurple color.Color = color.RGBA{R: 126, G: 37, B: 83, A: 255}
    DarkGreen  color.Color = color.RGBA{R: 0, G: 135, B: 81, A: 255}
    Brown      color.Color = color.RGBA{R: 171, G: 82, B: 54, A: 255}
    DarkGray   color.Color = color.RGBA{R: 95, G: 87, B: 79, A: 255}
    LightGray  color.Color = color.RGBA{R: 194, G: 195, B: 199, A: 255}
    White      color.Color = color.RGBA{R: 255, G: 241, B: 232, A: 255}
    Red        color.Color = color.RGBA{R: 255, G: 0, B: 77, A: 255}
    Orange     color.Color = color.RGBA{R: 255, G: 163, B: 0, A: 255}
    Yellow     color.Color = color.RGBA{R: 255, G: 236, B: 39, A: 255}
    Green      color.Color = color.RGBA{R: 0, G: 228, B: 54, A: 255}
    Blue       color.Color = color.RGBA{R: 41, G: 173, B: 255, A: 255}
    Indigo     color.Color = color.RGBA{R: 131, G: 118, B: 156, A: 255}
    Pink       color.Color = color.RGBA{R: 255, G: 119, B: 168, A: 255}
    Peach      color.Color = color.RGBA{R: 255, G: 204, B: 170, A: 255}
)
```

## Using Named Colors

Pass any named color directly to drawing functions:

```go
func (g *Game) Draw() {
    gosprite64.ClearScreen()                                         // fills with Black
    gosprite64.FillRect(10, 10, 100, 50, gosprite64.DarkBlue)       // filled rectangle
    gosprite64.DrawRect(10, 10, 100, 50, gosprite64.Yellow)         // outline
    gosprite64.DrawLine(0, 108, 287, 108, gosprite64.Red)           // horizontal line
    gosprite64.DrawText("HELLO N64", 100, 80, gosprite64.White)     // text
}
```

You can also clear the screen to any color:

```go
gosprite64.ClearScreenWith(gosprite64.DarkPurple)
```

## Using Custom Colors

Any value that satisfies Go's `color.Color` interface works with GoSprite64's drawing functions. The most common way to create a custom color is with `color.RGBA`:

```go
import "image/color"

// Fully opaque custom color
skyBlue := color.RGBA{R: 135, G: 206, B: 235, A: 255}
gosprite64.ClearScreenWith(skyBlue)

// Semi-transparent (for transitions or overlays)
shadow := color.RGBA{R: 0, G: 0, B: 0, A: 128}
gosprite64.FillRect(20, 20, 260, 196, shadow)
```

The `A` (alpha) field controls opacity: 255 is fully opaque, 0 is fully transparent.

## Performance Note

The 16 built-in colors are pre-cached internally as `image.Uniform` values, so using them is slightly faster than creating new `color.RGBA` values each frame. For custom colors that you use every frame, consider storing them in a variable rather than creating them inline.

```go
type Game struct {
    bgColor color.Color
}

func (g *Game) Init() {
    g.bgColor = color.RGBA{R: 20, G: 12, B: 28, A: 255}
}

func (g *Game) Draw() {
    gosprite64.ClearScreenWith(g.bgColor)
}
```
