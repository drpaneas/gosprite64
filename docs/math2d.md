# 2D Math

The `math2d` package provides the building blocks for 2D game math: vectors, rectangles, random numbers, and interpolation. Everything is portable Go with no N64 build tags, so you can use and test it on your host machine.

Import it alongside the main library:

```go
import "github.com/drpaneas/gosprite64/math2d"
```

## Vec2

`Vec2` is a 2D vector with `float32` components. It is a value type - all methods return a new Vec2 and never mutate the receiver.

```go
pos := math2d.Vec2{X: 100, Y: 50}
vel := math2d.Vec2{X: 2, Y: -1}

pos = pos.Add(vel)          // {102, 49}
pos = pos.Sub(vel)          // {100, 50}
doubled := vel.Scale(2)     // {4, -2}
flipped := vel.Negate()     // {-2, 1}
```

### Length and distance

```go
v := math2d.Vec2{X: 3, Y: 4}
v.Length()      // 5.0
v.LengthSq()   // 25.0 (no sqrt - faster for comparisons)

a := math2d.Vec2{X: 0, Y: 0}
b := math2d.Vec2{X: 3, Y: 4}
a.Distance(b)    // 5.0
a.DistanceSq(b)  // 25.0
```

Use `LengthSq` and `DistanceSq` when you only need to compare distances. They skip the square root, which matters on N64 hardware.

### Normalize

Returns a unit-length vector pointing in the same direction. Returns `{0, 0}` for zero or near-zero vectors (length squared < 1e-12).

```go
dir := math2d.Vec2{X: 3, Y: 4}.Normalize()  // {0.6, 0.8}
zero := math2d.Vec2{}.Normalize()             // {0, 0}
```

### Dot product

```go
right := math2d.Vec2{X: 1, Y: 0}
up    := math2d.Vec2{X: 0, Y: 1}
right.Dot(up)     // 0.0 (perpendicular)
right.Dot(right)  // 1.0 (parallel)
```

### Interpolation

`Lerp` linearly interpolates between two vectors. The `t` parameter is unclamped, so values outside 0..1 extrapolate beyond the endpoints. Use `math2d.Clamp` on `t` if you need clamped behavior.

```go
a := math2d.Vec2{X: 0, Y: 0}
b := math2d.Vec2{X: 100, Y: 200}

a.Lerp(b, 0.5)   // {50, 100}
a.Lerp(b, 0.0)   // {0, 0}
a.Lerp(b, 1.0)   // {100, 200}
a.Lerp(b, 1.5)   // {150, 300} - extrapolation
```

### Rotation and angle

`Rotate` rotates a vector by the given angle in radians (counterclockwise). `Angle` returns the angle of a vector in radians.

```go
right := math2d.Vec2{X: 1, Y: 0}
up := right.Rotate(math.Pi / 2)  // {0, 1}

up.Angle()    // ~1.5708 (pi/2)
right.Angle() // 0.0
```

### Min, Max, Abs

Per-component operations useful for bounding calculations:

```go
a := math2d.Vec2{X: 1, Y: 5}
b := math2d.Vec2{X: 3, Y: 2}

a.Min(b)  // {1, 2}
a.Max(b)  // {3, 5}

math2d.Vec2{X: -3, Y: -4}.Abs()  // {3, 4}
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

## Rand

`Rand` is a seedable pseudo-random number generator using the xoshiro128** algorithm. It is deterministic: the same seed always produces the same sequence. This is important for gameplay that needs to be reproducible - replays, netplay, or consistent procedural generation.

Unlike Go's `math/rand`, each `Rand` instance has its own state. There is no global generator to worry about.

```go
rng := math2d.NewRand(42)  // seed with 42
```

### Integers

```go
rng.Uint32()       // raw 32-bit value
rng.Intn(6)        // [0, 6) - like a die roll 0..5
rng.RangeInt(1, 7) // [1, 7) - die roll 1..6
```

### Floats

```go
rng.Float32()                  // [0.0, 1.0)
rng.RangeFloat32(0.5, 2.0)    // [0.5, 2.0)
```

### Bool

```go
if rng.Bool() {
    // roughly 50/50
}
```

### Re-seeding

Call `Seed` to reset the generator to a known state:

```go
rng := math2d.NewRand(42)
first := rng.Uint32()

rng.Seed(42)
second := rng.Uint32()
// first == second
```

### Typical game usage

```go
func (g *Game) Init() {
    g.rng = math2d.NewRand(12345)

    // Spawn enemies at random positions
    for i := 0; i < 10; i++ {
        x := g.rng.RangeFloat32(20, 268)
        y := g.rng.RangeFloat32(20, 196)
        spawnEnemy(x, y)
    }

    // Pick a random color
    colors := []color.Color{gosprite64.Red, gosprite64.Blue, gosprite64.Green}
    c := colors[g.rng.Intn(len(colors))]
}
```

## Easing and interpolation

### Lerp

Linearly interpolate between two scalar values. Unclamped - values of `t` outside 0..1 extrapolate.

```go
math2d.Lerp(0, 100, 0.5)   // 50
math2d.Lerp(0, 100, 0.0)   // 0
math2d.Lerp(0, 100, 1.0)   // 100
math2d.Lerp(0, 100, 1.5)   // 150 (extrapolation)
```

### InvLerp

The inverse of Lerp - given a value, find where it falls in a range as a 0..1 ratio:

```go
math2d.InvLerp(0, 100, 50)   // 0.5
math2d.InvLerp(0, 100, 0)    // 0.0
math2d.InvLerp(0, 100, 100)  // 1.0
```

### Remap

Map a value from one range to another. This is `InvLerp` followed by `Lerp`:

```go
// Map health (0..100) to a bar width (0..64 pixels)
barWidth := math2d.Remap(health, 0, 100, 0, 64)
```

### Clamp

Restrict a value to a range:

```go
math2d.Clamp(150, 0, 100)  // 100
math2d.Clamp(-5, 0, 100)   // 0
math2d.Clamp(50, 0, 100)   // 50
```

### MoveToward

Move a value toward a target by at most a fixed step. Useful for smooth following that arrives at the target in finite time (unlike lerp which asymptotically approaches).

```go
// In Update(), smoothly move camera X toward player
cam.X = math2d.MoveToward(cam.X, playerX, 3)
```

If `maxDelta` is zero or negative, the value does not change. The function never overshoots the target.

### Easing curves

Easing functions take a `t` in 0..1 and return a shaped 0..1 value. Combine them with `Lerp` to animate anything:

```go
t := float32(frame) / float32(totalFrames)  // 0..1 over time

// Ease in (slow start, fast end)
x := math2d.Lerp(startX, endX, math2d.EaseInQuad(t))

// Ease out (fast start, slow end)
x := math2d.Lerp(startX, endX, math2d.EaseOutQuad(t))

// Ease in-out (slow start, fast middle, slow end)
x := math2d.Lerp(startX, endX, math2d.EaseInOutQuad(t))
```

Available curves:

| Function | Shape | Use for |
|----------|-------|---------|
| `EaseInQuad` | Accelerate | Objects starting from rest |
| `EaseOutQuad` | Decelerate | Objects coming to a stop |
| `EaseInOutQuad` | Accelerate then decelerate | Smooth menu transitions |
| `EaseInCubic` | Stronger accelerate | More dramatic start |
| `EaseOutCubic` | Stronger decelerate | Snappy UI animations |
| `EaseInOutCubic` | Stronger both | Emphasis on endpoints |
| `SmoothStep` | Hermite S-curve | Thresholds, fog edges |

### SmoothStep

`SmoothStep` clamps the input to a range and applies Hermite interpolation. Useful for smooth thresholds:

```go
// Fade alpha from 0 to 1 as distance goes from 50 to 100
alpha := math2d.SmoothStep(50, 100, distance)
```

## Complete example

The `examples/math2d_demo` shows all of these working together: bouncing particles with Vec2 motion, Rect collision, seeded randomness, and easing-animated crosshair. Build it with:

```bash
GOENV=n64.env go1.24.5-embedded build -o math2d_demo.elf ./examples/math2d_demo
```
