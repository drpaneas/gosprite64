package math2d

import "testing"

func TestClamp(t *testing.T) {
	if Clamp(5, 0, 10) != 5 {
		t.Fatal("in-range value should pass through")
	}
	if Clamp(-1, 0, 10) != 0 {
		t.Fatal("below-min should clamp to min")
	}
	if Clamp(15, 0, 10) != 10 {
		t.Fatal("above-max should clamp to max")
	}
}

func TestLerp(t *testing.T) {
	if !almostEqual(Lerp(0, 10, 0.5), 5, 0.001) {
		t.Fatalf("Lerp(0, 10, 0.5) = %f, want 5", Lerp(0, 10, 0.5))
	}
	if !almostEqual(Lerp(0, 10, 0), 0, 0.001) {
		t.Fatal("Lerp at 0 should be a")
	}
	if !almostEqual(Lerp(0, 10, 1), 10, 0.001) {
		t.Fatal("Lerp at 1 should be b")
	}
}

func TestInvLerp(t *testing.T) {
	if !almostEqual(InvLerp(0, 10, 5), 0.5, 0.001) {
		t.Fatalf("InvLerp(0, 10, 5) = %f, want 0.5", InvLerp(0, 10, 5))
	}
	if !almostEqual(InvLerp(0, 10, 0), 0, 0.001) {
		t.Fatal("InvLerp at a should be 0")
	}
	if !almostEqual(InvLerp(0, 10, 10), 1, 0.001) {
		t.Fatal("InvLerp at b should be 1")
	}
}

func TestInvLerpSameEndpoints(t *testing.T) {
	if InvLerp(5, 5, 5) != 0 {
		t.Fatal("InvLerp with same endpoints should return 0")
	}
}

func TestEaseInQuad(t *testing.T) {
	if !almostEqual(EaseInQuad(0), 0, 0.001) {
		t.Fatal("EaseInQuad(0) should be 0")
	}
	if !almostEqual(EaseInQuad(1), 1, 0.001) {
		t.Fatal("EaseInQuad(1) should be 1")
	}
	if !almostEqual(EaseInQuad(0.5), 0.25, 0.001) {
		t.Fatalf("EaseInQuad(0.5) = %f, want 0.25", EaseInQuad(0.5))
	}
}

func TestEaseOutQuad(t *testing.T) {
	if !almostEqual(EaseOutQuad(0), 0, 0.001) {
		t.Fatal("EaseOutQuad(0) should be 0")
	}
	if !almostEqual(EaseOutQuad(1), 1, 0.001) {
		t.Fatal("EaseOutQuad(1) should be 1")
	}
	if !almostEqual(EaseOutQuad(0.5), 0.75, 0.001) {
		t.Fatalf("EaseOutQuad(0.5) = %f, want 0.75", EaseOutQuad(0.5))
	}
}

func TestEaseInOutQuad(t *testing.T) {
	if !almostEqual(EaseInOutQuad(0), 0, 0.001) {
		t.Fatal("EaseInOutQuad(0) should be 0")
	}
	if !almostEqual(EaseInOutQuad(1), 1, 0.001) {
		t.Fatal("EaseInOutQuad(1) should be 1")
	}
	if !almostEqual(EaseInOutQuad(0.5), 0.5, 0.001) {
		t.Fatalf("EaseInOutQuad(0.5) = %f, want 0.5", EaseInOutQuad(0.5))
	}
}

func TestMoveToward(t *testing.T) {
	if !almostEqual(MoveToward(0, 10, 3), 3, 0.001) {
		t.Fatal("should move 3 toward 10")
	}
	if !almostEqual(MoveToward(8, 10, 5), 10, 0.001) {
		t.Fatal("should not overshoot target")
	}
	if !almostEqual(MoveToward(10, 5, 3), 7, 0.001) {
		t.Fatal("should move backward toward target")
	}
}

func TestRemap(t *testing.T) {
	v := Remap(5, 0, 10, 100, 200)
	if !almostEqual(v, 150, 0.001) {
		t.Fatalf("Remap(5, 0,10, 100,200) = %f, want 150", v)
	}
}

func TestMoveTowardNegativeDelta(t *testing.T) {
	v := MoveToward(5, 10, -1)
	if v != 5 {
		t.Fatalf("negative maxDelta should return current, got %f", v)
	}
}

func TestMoveTowardZeroDelta(t *testing.T) {
	v := MoveToward(5, 10, 0)
	if v != 5 {
		t.Fatalf("zero maxDelta should return current, got %f", v)
	}
}

func TestEaseInCubic(t *testing.T) {
	if !almostEqual(EaseInCubic(0), 0, 0.001) {
		t.Fatal("EaseInCubic(0) should be 0")
	}
	if !almostEqual(EaseInCubic(1), 1, 0.001) {
		t.Fatal("EaseInCubic(1) should be 1")
	}
	if !almostEqual(EaseInCubic(0.5), 0.125, 0.001) {
		t.Fatalf("EaseInCubic(0.5) = %f, want 0.125", EaseInCubic(0.5))
	}
}

func TestEaseOutCubic(t *testing.T) {
	if !almostEqual(EaseOutCubic(0), 0, 0.001) {
		t.Fatalf("EaseOutCubic(0) = %f, want 0", EaseOutCubic(0))
	}
	if !almostEqual(EaseOutCubic(1), 1, 0.001) {
		t.Fatal("EaseOutCubic(1) should be 1")
	}
	if !almostEqual(EaseOutCubic(0.5), 0.875, 0.001) {
		t.Fatalf("EaseOutCubic(0.5) = %f, want 0.875", EaseOutCubic(0.5))
	}
}

func TestEaseInOutCubic(t *testing.T) {
	if !almostEqual(EaseInOutCubic(0), 0, 0.001) {
		t.Fatal("EaseInOutCubic(0) should be 0")
	}
	if !almostEqual(EaseInOutCubic(1), 1, 0.001) {
		t.Fatal("EaseInOutCubic(1) should be 1")
	}
	if !almostEqual(EaseInOutCubic(0.5), 0.5, 0.001) {
		t.Fatalf("EaseInOutCubic(0.5) = %f, want 0.5", EaseInOutCubic(0.5))
	}
}

func TestSmoothStep(t *testing.T) {
	if !almostEqual(SmoothStep(0, 1, -1), 0, 0.001) {
		t.Fatal("SmoothStep below edge0 should be 0")
	}
	if !almostEqual(SmoothStep(0, 1, 2), 1, 0.001) {
		t.Fatal("SmoothStep above edge1 should be 1")
	}
	if !almostEqual(SmoothStep(0, 1, 0.5), 0.5, 0.001) {
		t.Fatalf("SmoothStep(0,1,0.5) = %f, want 0.5", SmoothStep(0, 1, 0.5))
	}
	if !almostEqual(SmoothStep(0, 1, 0), 0, 0.001) {
		t.Fatal("SmoothStep at edge0 should be 0")
	}
	if !almostEqual(SmoothStep(0, 1, 1), 1, 0.001) {
		t.Fatal("SmoothStep at edge1 should be 1")
	}
}

func TestClampLoGreaterThanHi(t *testing.T) {
	v := Clamp(5, 10, 0)
	if v != 10 {
		t.Fatalf("Clamp(5, lo=10, hi=0) with lo>hi returns lo, got %f", v)
	}
}
