package math2d

import (
	"math"
	"testing"
)

func almostEqual(a, b, epsilon float32) bool {
	return float32(math.Abs(float64(a-b))) < epsilon
}

func TestVec2Add(t *testing.T) {
	a := Vec2{X: 1, Y: 2}
	b := Vec2{X: 3, Y: 4}
	got := a.Add(b)
	if got.X != 4 || got.Y != 6 {
		t.Fatalf("expected {4, 6}, got %v", got)
	}
}

func TestVec2Sub(t *testing.T) {
	a := Vec2{X: 5, Y: 7}
	b := Vec2{X: 2, Y: 3}
	got := a.Sub(b)
	if got.X != 3 || got.Y != 4 {
		t.Fatalf("expected {3, 4}, got %v", got)
	}
}

func TestVec2Scale(t *testing.T) {
	v := Vec2{X: 3, Y: 4}
	got := v.Scale(2)
	if got.X != 6 || got.Y != 8 {
		t.Fatalf("expected {6, 8}, got %v", got)
	}
}

func TestVec2Length(t *testing.T) {
	v := Vec2{X: 3, Y: 4}
	if !almostEqual(v.Length(), 5, 0.001) {
		t.Fatalf("expected 5, got %f", v.Length())
	}
}

func TestVec2LengthSq(t *testing.T) {
	v := Vec2{X: 3, Y: 4}
	if v.LengthSq() != 25 {
		t.Fatalf("expected 25, got %f", v.LengthSq())
	}
}

func TestVec2Normalize(t *testing.T) {
	v := Vec2{X: 3, Y: 4}
	n := v.Normalize()
	if !almostEqual(n.Length(), 1, 0.001) {
		t.Fatalf("normalized length should be 1, got %f", n.Length())
	}
	if !almostEqual(n.X, 0.6, 0.001) || !almostEqual(n.Y, 0.8, 0.001) {
		t.Fatalf("expected {0.6, 0.8}, got %v", n)
	}
}

func TestVec2NormalizeZero(t *testing.T) {
	v := Vec2{X: 0, Y: 0}
	n := v.Normalize()
	if n.X != 0 || n.Y != 0 {
		t.Fatalf("normalizing zero vector should return zero, got %v", n)
	}
}

func TestVec2Dot(t *testing.T) {
	a := Vec2{X: 1, Y: 0}
	b := Vec2{X: 0, Y: 1}
	if a.Dot(b) != 0 {
		t.Fatalf("perpendicular dot should be 0, got %f", a.Dot(b))
	}
	c := Vec2{X: 1, Y: 0}
	if a.Dot(c) != 1 {
		t.Fatalf("parallel dot should be 1, got %f", a.Dot(c))
	}
}

func TestVec2Distance(t *testing.T) {
	a := Vec2{X: 0, Y: 0}
	b := Vec2{X: 3, Y: 4}
	if !almostEqual(a.Distance(b), 5, 0.001) {
		t.Fatalf("expected 5, got %f", a.Distance(b))
	}
}

func TestVec2DistanceSq(t *testing.T) {
	a := Vec2{X: 0, Y: 0}
	b := Vec2{X: 3, Y: 4}
	if a.DistanceSq(b) != 25 {
		t.Fatalf("expected 25, got %f", a.DistanceSq(b))
	}
}

func TestVec2Lerp(t *testing.T) {
	a := Vec2{X: 0, Y: 0}
	b := Vec2{X: 10, Y: 20}
	mid := a.Lerp(b, 0.5)
	if !almostEqual(mid.X, 5, 0.001) || !almostEqual(mid.Y, 10, 0.001) {
		t.Fatalf("expected {5, 10}, got %v", mid)
	}
}

func TestVec2LerpClamp(t *testing.T) {
	a := Vec2{X: 0, Y: 0}
	b := Vec2{X: 10, Y: 20}
	over := a.Lerp(b, 1.5)
	if !almostEqual(over.X, 10, 0.001) || !almostEqual(over.Y, 20, 0.001) {
		t.Fatalf("lerp t>1 should clamp to b, got %v", over)
	}
	under := a.Lerp(b, -0.5)
	if !almostEqual(under.X, 0, 0.001) || !almostEqual(under.Y, 0, 0.001) {
		t.Fatalf("lerp t<0 should clamp to a, got %v", under)
	}
}

func TestVec2Rotate(t *testing.T) {
	v := Vec2{X: 1, Y: 0}
	r := v.Rotate(math.Pi / 2)
	if !almostEqual(r.X, 0, 0.001) || !almostEqual(r.Y, 1, 0.001) {
		t.Fatalf("expected {0, 1}, got %v", r)
	}
}

func TestVec2Angle(t *testing.T) {
	v := Vec2{X: 1, Y: 0}
	if !almostEqual(v.Angle(), 0, 0.001) {
		t.Fatalf("expected 0, got %f", v.Angle())
	}
	v2 := Vec2{X: 0, Y: 1}
	if !almostEqual(v2.Angle(), math.Pi/2, 0.001) {
		t.Fatalf("expected pi/2, got %f", v2.Angle())
	}
}

func TestVec2Negate(t *testing.T) {
	v := Vec2{X: 3, Y: -4}
	n := v.Negate()
	if n.X != -3 || n.Y != 4 {
		t.Fatalf("expected {-3, 4}, got %v", n)
	}
}

func TestVec2Abs(t *testing.T) {
	v := Vec2{X: -3, Y: -4}
	a := v.Abs()
	if a.X != 3 || a.Y != 4 {
		t.Fatalf("expected {3, 4}, got %v", a)
	}
}

func TestVec2Min(t *testing.T) {
	a := Vec2{X: 1, Y: 5}
	b := Vec2{X: 3, Y: 2}
	m := a.Min(b)
	if m.X != 1 || m.Y != 2 {
		t.Fatalf("expected {1, 2}, got %v", m)
	}
}

func TestVec2Max(t *testing.T) {
	a := Vec2{X: 1, Y: 5}
	b := Vec2{X: 3, Y: 2}
	m := a.Max(b)
	if m.X != 3 || m.Y != 5 {
		t.Fatalf("expected {3, 5}, got %v", m)
	}
}
