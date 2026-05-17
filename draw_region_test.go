package gosprite64

import "testing"

func TestDrawRegionOffset(t *testing.T) {
	r := DrawRegion{X: 50, Y: 30, W: 100, H: 80}
	ox, oy := r.Offset(10, 20)
	if ox != 60 || oy != 50 {
		t.Fatalf("expected (60, 50), got (%d, %d)", ox, oy)
	}
}

func TestDrawRegionClipInside(t *testing.T) {
	r := DrawRegion{X: 50, Y: 30, W: 100, H: 80}
	x1, y1, x2, y2, ok := r.Clip(10, 10, 60, 40)
	if !ok {
		t.Fatal("rect inside region should not be clipped away")
	}
	if x1 != 60 || y1 != 40 || x2 != 110 || y2 != 70 {
		t.Fatalf("expected (60,40,110,70), got (%d,%d,%d,%d)", x1, y1, x2, y2)
	}
}

func TestDrawRegionClipPartial(t *testing.T) {
	r := DrawRegion{X: 50, Y: 30, W: 100, H: 80}
	x1, y1, x2, y2, ok := r.Clip(-10, -10, 200, 200)
	if !ok {
		t.Fatal("overlapping rect should not be clipped away")
	}
	if x1 != 50 || y1 != 30 {
		t.Fatalf("should clamp to region min, got (%d,%d)", x1, y1)
	}
	if x2 != 150 || y2 != 110 {
		t.Fatalf("should clamp to region max, got (%d,%d)", x2, y2)
	}
}

func TestDrawRegionClipOutside(t *testing.T) {
	r := DrawRegion{X: 50, Y: 30, W: 100, H: 80}
	_, _, _, _, ok := r.Clip(200, 200, 250, 250)
	if ok {
		t.Fatal("rect fully outside region should be clipped away")
	}
}

func TestDrawRegionZeroIsFullScreen(t *testing.T) {
	r := DrawRegion{}
	if r.Active() {
		t.Fatal("zero DrawRegion should not be active")
	}
}

func TestDrawRegionContainsPoint(t *testing.T) {
	r := DrawRegion{X: 50, Y: 30, W: 100, H: 80}
	if !r.ContainsPoint(70, 50) {
		t.Fatal("point inside should be contained")
	}
	if r.ContainsPoint(200, 200) {
		t.Fatal("point outside should not be contained")
	}
}

func TestCurrentDrawRegionDefault(t *testing.T) {
	r := currentDrawRegion()
	if r.Active() {
		t.Fatal("default draw region should not be active")
	}
}

func TestSetResetDrawRegionNesting(t *testing.T) {
	drawRegionStack = nil

	SetDrawRegion(10, 10, 100, 100)
	if currentDrawRegion().X != 10 {
		t.Fatal("first region should be active")
	}

	SetDrawRegion(20, 20, 50, 50)
	if currentDrawRegion().X != 20 {
		t.Fatal("nested region should be active")
	}
	if len(drawRegionStack) != 2 {
		t.Fatalf("expected stack depth 2, got %d", len(drawRegionStack))
	}

	ResetDrawRegion()
	if currentDrawRegion().X != 10 {
		t.Fatal("after pop, first region should be active")
	}

	ResetDrawRegion()
	if currentDrawRegion().Active() {
		t.Fatal("after second pop, should be back to full screen")
	}

	drawRegionStack = nil
}

func TestResetDrawRegionEmptyStack(t *testing.T) {
	drawRegionStack = nil
	ResetDrawRegion()
	if len(drawRegionStack) != 0 {
		t.Fatal("resetting empty stack should not panic or grow")
	}
}
