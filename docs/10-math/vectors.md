# Vectors

The `math2d` package provides 2D vectors for game math. Import it with:

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
