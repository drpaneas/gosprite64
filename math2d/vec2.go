package math2d

import "math"

// Vec2 is a 2D vector with float32 components.
type Vec2 struct {
	X, Y float32
}

func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{X: v.X + other.X, Y: v.Y + other.Y}
}

func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{X: v.X - other.X, Y: v.Y - other.Y}
}

func (v Vec2) Scale(s float32) Vec2 {
	return Vec2{X: v.X * s, Y: v.Y * s}
}

func (v Vec2) Negate() Vec2 {
	return Vec2{X: -v.X, Y: -v.Y}
}

func (v Vec2) Abs() Vec2 {
	return Vec2{
		X: float32(math.Abs(float64(v.X))),
		Y: float32(math.Abs(float64(v.Y))),
	}
}

func (v Vec2) LengthSq() float32 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vec2) Length() float32 {
	return float32(math.Sqrt(float64(v.LengthSq())))
}

func (v Vec2) Normalize() Vec2 {
	lsq := v.LengthSq()
	if lsq < 1e-12 {
		return Vec2{}
	}
	l := float32(math.Sqrt(float64(lsq)))
	return Vec2{X: v.X / l, Y: v.Y / l}
}

func (v Vec2) Dot(other Vec2) float32 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vec2) Distance(other Vec2) float32 {
	return v.Sub(other).Length()
}

func (v Vec2) DistanceSq(other Vec2) float32 {
	return v.Sub(other).LengthSq()
}

// Lerp linearly interpolates between v and other by t.
// t is unclamped, matching the scalar Lerp in easing.go.
// Use Clamp on t beforehand if you need clamped behavior.
func (v Vec2) Lerp(other Vec2, t float32) Vec2 {
	return Vec2{
		X: v.X + (other.X-v.X)*t,
		Y: v.Y + (other.Y-v.Y)*t,
	}
}

func (v Vec2) Rotate(radians float64) Vec2 {
	cos := float32(math.Cos(radians))
	sin := float32(math.Sin(radians))
	return Vec2{
		X: v.X*cos - v.Y*sin,
		Y: v.X*sin + v.Y*cos,
	}
}

func (v Vec2) Angle() float32 {
	return float32(math.Atan2(float64(v.Y), float64(v.X)))
}

func (v Vec2) Min(other Vec2) Vec2 {
	return Vec2{
		X: float32(math.Min(float64(v.X), float64(other.X))),
		Y: float32(math.Min(float64(v.Y), float64(other.Y))),
	}
}

func (v Vec2) Max(other Vec2) Vec2 {
	return Vec2{
		X: float32(math.Max(float64(v.X), float64(other.X))),
		Y: float32(math.Max(float64(v.Y), float64(other.Y))),
	}
}
