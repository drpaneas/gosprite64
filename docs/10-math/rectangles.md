# Rectangles

The `math2d` package provides axis-aligned rectangles for 2D game math. Import it with:

```go
import "github.com/drpaneas/gosprite64/math2d"
```

## Rect

`Rect` is an axis-aligned rectangle defined by its top-left corner and dimensions. It uses half-open boundary semantics: a point at `(X+W, Y+H)` is outside the rectangle. This matches how pixels and tiles work - a 10-pixel-wide rect at X=0 occupies pixels 0 through 9.

```go
court := math2d.Rect{X: 10, Y: 10, W: 268, H: 196}

court.Right()   // 278.0 (X + W)
court.Bottom()  // 206.0 (Y + H)
court.Center()  // Vec2{144, 108}
```

### Creating from center

```go
r := math2d.RectFromCenter(math2d.Vec2{X: 144, Y: 108}, 20, 20)
// Rect{X: 134, Y: 98, W: 20, H: 20}
```

### Point containment

```go
court := math2d.Rect{X: 0, Y: 0, W: 288, H: 216}
court.ContainsPoint(math2d.Vec2{X: 100, Y: 100})  // true
court.ContainsPoint(math2d.Vec2{X: 288, Y: 216})  // false (half-open)
court.ContainsPoint(math2d.Vec2{X: -1, Y: 0})     // false
```

### Overlap and collision

```go
player := math2d.Rect{X: 50, Y: 50, W: 16, H: 16}
enemy  := math2d.Rect{X: 60, Y: 55, W: 16, H: 16}

player.Overlaps(enemy)  // true

inter, ok := player.Intersection(enemy)
// ok = true, inter = Rect{X: 60, Y: 55, W: 6, H: 11}
```

Zero-size or negative-size rects never overlap, never contain points, and never contain other rects.

### Rect containment

```go
screen := math2d.Rect{X: 0, Y: 0, W: 288, H: 216}
button := math2d.Rect{X: 100, Y: 80, W: 40, H: 20}
screen.ContainsRect(button)  // true
button.ContainsRect(screen)  // false
```

### Expand and shrink

`Expand` grows (or shrinks with a negative amount) a rect on all sides:

```go
r := math2d.Rect{X: 10, Y: 10, W: 20, H: 20}
r.Expand(5)   // {5, 5, 30, 30}
r.Expand(-2)  // {12, 12, 16, 16}
```
