package math2d

import "math"

// AABBOverlap returns true if two axis-aligned bounding boxes overlap.
func AABBOverlap(a, b Rect) bool {
	return a.Overlaps(b)
}

// AABBPenetration returns the minimum translation vector to separate a from b.
// The vector pushes a away from b along the axis of least overlap.
// Returns false if there is no overlap.
func AABBPenetration(a, b Rect) (Vec2, bool) {
	if !a.Overlaps(b) {
		return Vec2{}, false
	}

	left := a.Right() - b.X
	right := b.Right() - a.X
	top := a.Bottom() - b.Y
	bottom := b.Bottom() - a.Y

	penX := left
	signX := float32(-1)
	if right < left {
		penX = right
		signX = 1
	}

	penY := top
	signY := float32(-1)
	if bottom < top {
		penY = bottom
		signY = 1
	}

	if penX <= penY {
		return Vec2{X: signX * penX, Y: 0}, true
	}
	return Vec2{X: 0, Y: signY * penY}, true
}

// AABBResolve returns a copy of rect a moved by the minimum translation
// vector so it no longer overlaps b.
func AABBResolve(a, b Rect) Rect {
	pen, ok := AABBPenetration(a, b)
	if !ok {
		return a
	}
	return Rect{X: a.X + pen.X, Y: a.Y + pen.Y, W: a.W, H: a.H}
}

// AABBSweep performs a swept AABB test: moves rect a by velocity and checks
// for collision with static rect b.
// Returns (hit, t, normal) where t is the fraction of velocity at first contact
// and normal is the collision surface normal.
// If already overlapping, returns (true, 0, {0,0}).
func AABBSweep(a Rect, vel Vec2, b Rect) (bool, float32, Vec2) {
	if a.Overlaps(b) {
		return true, 0, Vec2{}
	}

	if vel.X == 0 && vel.Y == 0 {
		return false, 0, Vec2{}
	}

	expanded := Rect{
		X: b.X - a.W,
		Y: b.Y - a.H,
		W: b.W + a.W,
		H: b.H + a.H,
	}

	origin := Vec2{X: a.X, Y: a.Y}
	var tNearX, tFarX, tNearY, tFarY float32

	if vel.X != 0 {
		tNearX = (expanded.X - origin.X) / vel.X
		tFarX = (expanded.Right() - origin.X) / vel.X
		if tNearX > tFarX {
			tNearX, tFarX = tFarX, tNearX
		}
	} else {
		if origin.X < expanded.X || origin.X >= expanded.Right() {
			return false, 0, Vec2{}
		}
		tNearX = float32(math.Inf(-1))
		tFarX = float32(math.Inf(1))
	}

	if vel.Y != 0 {
		tNearY = (expanded.Y - origin.Y) / vel.Y
		tFarY = (expanded.Bottom() - origin.Y) / vel.Y
		if tNearY > tFarY {
			tNearY, tFarY = tFarY, tNearY
		}
	} else {
		if origin.Y < expanded.Y || origin.Y >= expanded.Bottom() {
			return false, 0, Vec2{}
		}
		tNearY = float32(math.Inf(-1))
		tFarY = float32(math.Inf(1))
	}

	if tNearX > tFarY || tNearY > tFarX {
		return false, 0, Vec2{}
	}

	tMin := tNearX
	if tNearY > tMin {
		tMin = tNearY
	}

	tMax := tFarX
	if tFarY < tMax {
		tMax = tFarY
	}

	if tMin >= 1 || tMax <= 0 {
		return false, 0, Vec2{}
	}

	if tMin < 0 {
		tMin = 0
	}

	var normal Vec2
	if tNearX > tNearY {
		if vel.X < 0 {
			normal = Vec2{X: 1, Y: 0}
		} else {
			normal = Vec2{X: -1, Y: 0}
		}
	} else {
		if vel.Y < 0 {
			normal = Vec2{X: 0, Y: 1}
		} else {
			normal = Vec2{X: 0, Y: -1}
		}
	}

	return true, tMin, normal
}
