package math2d

import "math"

// Clamp restricts v to the range [lo, hi].
func Clamp(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// Lerp linearly interpolates between a and b by t (unclamped).
func Lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}

// InvLerp returns where v falls between a and b as a 0..1 ratio.
// Returns 0 if a == b.
func InvLerp(a, b, v float32) float32 {
	if a == b {
		return 0
	}
	return (v - a) / (b - a)
}

// Remap maps v from range [inMin, inMax] to [outMin, outMax].
func Remap(v, inMin, inMax, outMin, outMax float32) float32 {
	t := InvLerp(inMin, inMax, v)
	return Lerp(outMin, outMax, t)
}

// MoveToward moves current toward target by at most maxDelta.
func MoveToward(current, target, maxDelta float32) float32 {
	diff := target - current
	if float32(math.Abs(float64(diff))) <= maxDelta {
		return target
	}
	if diff > 0 {
		return current + maxDelta
	}
	return current - maxDelta
}

// EaseInQuad accelerates from zero.
func EaseInQuad(t float32) float32 { return t * t }

// EaseOutQuad decelerates to zero.
func EaseOutQuad(t float32) float32 { return t * (2 - t) }

// EaseInOutQuad accelerates then decelerates.
func EaseInOutQuad(t float32) float32 {
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}

// EaseInCubic accelerates from zero (cubic).
func EaseInCubic(t float32) float32 { return t * t * t }

// EaseOutCubic decelerates to zero (cubic).
func EaseOutCubic(t float32) float32 {
	t--
	return 1 + t*t*t
}

// EaseInOutCubic accelerates then decelerates (cubic).
func EaseInOutCubic(t float32) float32 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	t = 2*t - 2
	return 0.5*t*t*t + 1
}

// SmoothStep performs Hermite interpolation (3t^2 - 2t^3) after clamping.
func SmoothStep(edge0, edge1, x float32) float32 {
	t := Clamp(InvLerp(edge0, edge1, x), 0, 1)
	return t * t * (3 - 2*t)
}
