package math2d

import "testing"

func TestRectContainsPoint(t *testing.T) {
	r := Rect{X: 10, Y: 10, W: 100, H: 50}
	if !r.ContainsPoint(Vec2{X: 50, Y: 30}) {
		t.Fatal("center point should be contained")
	}
	if !r.ContainsPoint(Vec2{X: 10, Y: 10}) {
		t.Fatal("top-left corner should be contained")
	}
	if r.ContainsPoint(Vec2{X: 110, Y: 60}) {
		t.Fatal("bottom-right edge should NOT be contained (half-open)")
	}
	if r.ContainsPoint(Vec2{X: 5, Y: 30}) {
		t.Fatal("point left of rect should not be contained")
	}
}

func TestRectOverlaps(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 10, H: 10}
	b := Rect{X: 5, Y: 5, W: 10, H: 10}
	if !a.Overlaps(b) {
		t.Fatal("overlapping rects should overlap")
	}
	c := Rect{X: 10, Y: 0, W: 10, H: 10}
	if a.Overlaps(c) {
		t.Fatal("touching-edge rects should NOT overlap (half-open)")
	}
	d := Rect{X: 20, Y: 20, W: 5, H: 5}
	if a.Overlaps(d) {
		t.Fatal("non-overlapping rects should not overlap")
	}
}

func TestRectIntersection(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 10, H: 10}
	b := Rect{X: 5, Y: 5, W: 10, H: 10}
	inter, ok := a.Intersection(b)
	if !ok {
		t.Fatal("overlapping rects should have intersection")
	}
	if inter.X != 5 || inter.Y != 5 || inter.W != 5 || inter.H != 5 {
		t.Fatalf("expected {5,5,5,5}, got %v", inter)
	}
}

func TestRectIntersectionDisjoint(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 5, H: 5}
	b := Rect{X: 10, Y: 10, W: 5, H: 5}
	_, ok := a.Intersection(b)
	if ok {
		t.Fatal("disjoint rects should have no intersection")
	}
}

func TestRectCenter(t *testing.T) {
	r := Rect{X: 10, Y: 20, W: 100, H: 50}
	c := r.Center()
	if !almostEqual(c.X, 60, 0.001) || !almostEqual(c.Y, 45, 0.001) {
		t.Fatalf("expected {60, 45}, got %v", c)
	}
}

func TestRectExpand(t *testing.T) {
	r := Rect{X: 10, Y: 10, W: 20, H: 20}
	e := r.Expand(5)
	if e.X != 5 || e.Y != 5 || e.W != 30 || e.H != 30 {
		t.Fatalf("expected {5,5,30,30}, got %v", e)
	}
}

func TestRectFromCenter(t *testing.T) {
	r := RectFromCenter(Vec2{X: 50, Y: 50}, 20, 10)
	if r.X != 40 || r.Y != 45 || r.W != 20 || r.H != 10 {
		t.Fatalf("expected {40,45,20,10}, got %v", r)
	}
}

func TestRectContainsRect(t *testing.T) {
	outer := Rect{X: 0, Y: 0, W: 100, H: 100}
	inner := Rect{X: 10, Y: 10, W: 20, H: 20}
	if !outer.ContainsRect(inner) {
		t.Fatal("inner should be contained")
	}
	if inner.ContainsRect(outer) {
		t.Fatal("outer should not be contained in inner")
	}
}

func TestRectRight(t *testing.T) {
	r := Rect{X: 10, Y: 20, W: 30, H: 40}
	if r.Right() != 40 {
		t.Fatalf("expected 40, got %f", r.Right())
	}
}

func TestRectBottom(t *testing.T) {
	r := Rect{X: 10, Y: 20, W: 30, H: 40}
	if r.Bottom() != 60 {
		t.Fatalf("expected 60, got %f", r.Bottom())
	}
}

func TestRectZeroSize(t *testing.T) {
	r := Rect{X: 5, Y: 5, W: 0, H: 0}
	if r.ContainsPoint(Vec2{X: 5, Y: 5}) {
		t.Fatal("zero-size rect should contain nothing")
	}
	if r.Overlaps(Rect{X: 0, Y: 0, W: 10, H: 10}) {
		t.Fatal("zero-size rect should not overlap anything")
	}
}

func TestRectZeroSizeContainsRect(t *testing.T) {
	r := Rect{X: 5, Y: 5, W: 0, H: 0}
	inner := Rect{X: 5, Y: 5, W: 1, H: 1}
	if r.ContainsRect(inner) {
		t.Fatal("zero-size rect should not contain anything")
	}
}

func TestRectExpandNegative(t *testing.T) {
	r := Rect{X: 10, Y: 10, W: 20, H: 20}
	e := r.Expand(-5)
	if e.X != 15 || e.Y != 15 || e.W != 10 || e.H != 10 {
		t.Fatalf("expected {15,15,10,10}, got %v", e)
	}
}

func TestRectExpandNegativeCollapse(t *testing.T) {
	r := Rect{X: 10, Y: 10, W: 4, H: 4}
	e := r.Expand(-5)
	if e.W >= 0 || e.H >= 0 {
		t.Logf("shrinking past zero produces negative W/H: %v (expected)", e)
	}
	if e.Overlaps(Rect{X: 0, Y: 0, W: 100, H: 100}) {
		t.Fatal("collapsed rect should not overlap anything")
	}
}
