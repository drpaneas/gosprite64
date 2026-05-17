# Easing Functions

The `math2d` package provides easing and interpolation functions. Import it with:

```go
import "github.com/drpaneas/gosprite64/math2d"
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
