package math2d

import "testing"

func TestAABBOverlap(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 10, H: 10}
	b := Rect{X: 5, Y: 5, W: 10, H: 10}
	if !AABBOverlap(a, b) {
		t.Fatal("overlapping rects should collide")
	}
}

func TestAABBNoOverlap(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 10, H: 10}
	b := Rect{X: 20, Y: 20, W: 10, H: 10}
	if AABBOverlap(a, b) {
		t.Fatal("non-overlapping rects should not collide")
	}
}

func TestAABBTouching(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 10, H: 10}
	b := Rect{X: 10, Y: 0, W: 10, H: 10}
	if AABBOverlap(a, b) {
		t.Fatal("edge-touching rects should NOT collide (half-open)")
	}
}

func TestAABBPenetration(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 10, H: 10}
	b := Rect{X: 7, Y: 3, W: 10, H: 10}
	pen, ok := AABBPenetration(a, b)
	if !ok {
		t.Fatal("overlapping rects should have penetration")
	}
	if !almostEqual(pen.X, -3, 0.001) || !almostEqual(pen.Y, 0, 0.001) {
		t.Fatalf("expected penetration along min axis X={-3, 0}, got %v", pen)
	}
}

func TestAABBPenetrationVertical(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 10, H: 10}
	b := Rect{X: 3, Y: 7, W: 10, H: 10}
	pen, ok := AABBPenetration(a, b)
	if !ok {
		t.Fatal("should overlap")
	}
	if pen.Y != -3 {
		t.Fatalf("should resolve along minimum penetration axis Y, got Y=%f", pen.Y)
	}
	if pen.X != 0 {
		t.Fatalf("non-min axis should be 0, got X=%f", pen.X)
	}
}

func TestAABBPenetrationNoOverlap(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 5, H: 5}
	b := Rect{X: 10, Y: 10, W: 5, H: 5}
	_, ok := AABBPenetration(a, b)
	if ok {
		t.Fatal("no overlap should return false")
	}
}

func TestAABBResolve(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 10, H: 10}
	b := Rect{X: 8, Y: 3, W: 10, H: 10}
	resolved := AABBResolve(a, b)
	if AABBOverlap(resolved, b) {
		t.Fatal("resolved rect should no longer overlap b")
	}
}

func TestAABBResolveNoOverlap(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 5, H: 5}
	b := Rect{X: 20, Y: 20, W: 5, H: 5}
	resolved := AABBResolve(a, b)
	if resolved != a {
		t.Fatal("no-overlap resolve should return original rect")
	}
}

func TestAABBSweepHit(t *testing.T) {
	a := Rect{X: 0, Y: 5, W: 4, H: 4}
	vel := Vec2{X: 10, Y: 0}
	b := Rect{X: 8, Y: 5, W: 4, H: 4}
	hit, t0, normal := AABBSweep(a, vel, b)
	if !hit {
		t.Fatal("should detect sweep collision")
	}
	if t0 < 0 || t0 > 1 {
		t.Fatalf("sweep t should be in [0,1], got %f", t0)
	}
	if !almostEqual(t0, 0.4, 0.01) {
		t.Fatalf("expected t~0.4, got %f", t0)
	}
	if normal.X != -1 {
		t.Fatalf("expected normal {-1,0}, got %v", normal)
	}
}

func TestAABBSweepMiss(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 4, H: 4}
	vel := Vec2{X: 10, Y: 0}
	b := Rect{X: 8, Y: 20, W: 4, H: 4}
	hit, _, _ := AABBSweep(a, vel, b)
	if hit {
		t.Fatal("should not detect collision - different Y lanes")
	}
}

func TestAABBSweepZeroVelocity(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 4, H: 4}
	vel := Vec2{X: 0, Y: 0}
	b := Rect{X: 8, Y: 0, W: 4, H: 4}
	hit, _, _ := AABBSweep(a, vel, b)
	if hit {
		t.Fatal("zero velocity should not hit distant object")
	}
}

func TestAABBSweepAlreadyOverlapping(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 10, H: 10}
	vel := Vec2{X: 1, Y: 0}
	b := Rect{X: 5, Y: 0, W: 10, H: 10}
	hit, t0, _ := AABBSweep(a, vel, b)
	if !hit {
		t.Fatal("already overlapping should be a hit")
	}
	if t0 != 0 {
		t.Fatalf("already overlapping should have t=0, got %f", t0)
	}
}
