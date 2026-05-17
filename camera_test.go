package gosprite64

import (
	"testing"

	"github.com/drpaneas/gosprite64/math2d"
)

func almostEq(a, b, eps float32) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d < eps
}

func TestCameraZoomDefault(t *testing.T) {
	c := &Camera{}
	if c.EffectiveZoom() != 1 {
		t.Fatalf("zero Zoom should default to 1, got %f", c.EffectiveZoom())
	}
}

func TestCameraZoomExplicit(t *testing.T) {
	c := &Camera{Zoom: 2.0}
	if c.EffectiveZoom() != 2.0 {
		t.Fatalf("expected 2.0, got %f", c.EffectiveZoom())
	}
}

func TestCameraWorldToScreen(t *testing.T) {
	c := &Camera{X: 100, Y: 50, Width: 288, Height: 216}
	sx, sy := c.WorldToScreen(150, 75)
	if sx != 50 || sy != 25 {
		t.Fatalf("expected (50, 25), got (%f, %f)", sx, sy)
	}
}

func TestCameraWorldToScreenWithZoom(t *testing.T) {
	c := &Camera{X: 100, Y: 50, Width: 288, Height: 216, Zoom: 2.0}
	sx, sy := c.WorldToScreen(150, 75)
	if !almostEq(sx, 100, 0.1) || !almostEq(sy, 50, 0.1) {
		t.Fatalf("expected (100, 50), got (%f, %f)", sx, sy)
	}
}

func TestCameraFollowInstant(t *testing.T) {
	c := &Camera{Width: 288, Height: 216}
	c.FollowSpeed = 1.0
	c.FollowTarget = &math2d.Vec2{X: 200, Y: 150}
	c.UpdateFollow()
	expectedX := float32(200 - 288/2)
	expectedY := float32(150 - 216/2)
	if !almostEq(float32(c.X), expectedX, 0.1) {
		t.Fatalf("expected X=%f, got %d", expectedX, c.X)
	}
	if !almostEq(float32(c.Y), expectedY, 0.1) {
		t.Fatalf("expected Y=%f, got %d", expectedY, c.Y)
	}
}

func TestCameraFollowLerp(t *testing.T) {
	c := &Camera{X: 0, Y: 0, Width: 288, Height: 216}
	c.FollowSpeed = 0.1
	c.FollowTarget = &math2d.Vec2{X: 200, Y: 150}

	c.UpdateFollow()
	if c.X == 0 && c.Y == 0 {
		t.Fatal("camera should have moved toward target")
	}
	targetX := 200 - 288/2
	if c.X >= targetX {
		t.Fatal("camera should not have reached target in one step at speed 0.1")
	}
}

func TestCameraFollowNilTarget(t *testing.T) {
	c := &Camera{X: 50, Y: 50}
	c.FollowSpeed = 1.0
	c.FollowTarget = nil
	c.UpdateFollow()
	if c.X != 50 || c.Y != 50 {
		t.Fatal("nil target should be a no-op")
	}
}

func TestCameraBoundsClamp(t *testing.T) {
	c := &Camera{X: -10, Y: -20, Width: 288, Height: 216}
	c.Bounds = &math2d.Rect{X: 0, Y: 0, W: 500, H: 400}
	c.ClampToBounds()
	if c.X < 0 {
		t.Fatalf("X should be clamped to >= 0, got %d", c.X)
	}
	if c.Y < 0 {
		t.Fatalf("Y should be clamped to >= 0, got %d", c.Y)
	}
}

func TestCameraBoundsClampRight(t *testing.T) {
	c := &Camera{X: 300, Y: 0, Width: 288, Height: 216}
	c.Bounds = &math2d.Rect{X: 0, Y: 0, W: 500, H: 400}
	c.ClampToBounds()
	maxX := 500 - 288
	if c.X > maxX {
		t.Fatalf("X should be clamped to <= %d, got %d", maxX, c.X)
	}
}

func TestCameraBoundsNilIsNoop(t *testing.T) {
	c := &Camera{X: -100, Y: -100}
	c.Bounds = nil
	c.ClampToBounds()
	if c.X != -100 {
		t.Fatal("nil bounds should not clamp")
	}
}

func TestCameraShake(t *testing.T) {
	c := &Camera{X: 100, Y: 100}
	c.AddTrauma(0.5)
	if c.trauma == 0 {
		t.Fatal("trauma should be set")
	}
	ox, oy := c.ShakeOffset()
	_ = ox
	_ = oy

	for i := 0; i < 60; i++ {
		c.UpdateShake()
	}
	if c.trauma != 0 {
		t.Fatalf("trauma should decay to 0 after 60 frames, got %f", c.trauma)
	}
}

func TestCameraShakeTraumaCaps(t *testing.T) {
	c := &Camera{}
	c.AddTrauma(0.8)
	c.AddTrauma(0.5)
	if c.trauma > 1.0 {
		t.Fatalf("trauma should cap at 1.0, got %f", c.trauma)
	}
}

func TestCameraShakeZeroTrauma(t *testing.T) {
	c := &Camera{}
	ox, oy := c.ShakeOffset()
	if ox != 0 || oy != 0 {
		t.Fatal("zero trauma should produce zero offset")
	}
}
