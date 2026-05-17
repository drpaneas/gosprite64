package gosprite64

import "testing"

func TestParallaxHalfSpeed(t *testing.T) {
	p := ParallaxLayer{SpeedX: 0.5, SpeedY: 0.5}
	ox, oy := p.Offset(100, 200)
	if ox != 50 {
		t.Fatalf("expected X=50, got %d", ox)
	}
	if oy != 100 {
		t.Fatalf("expected Y=100, got %d", oy)
	}
}

func TestParallaxNoScroll(t *testing.T) {
	p := ParallaxLayer{SpeedX: 0, SpeedY: 0}
	ox, oy := p.Offset(100, 200)
	if ox != 0 || oy != 0 {
		t.Fatalf("expected 0,0 got %d,%d", ox, oy)
	}
}

func TestParallaxFullSpeed(t *testing.T) {
	p := ParallaxLayer{SpeedX: 1.0, SpeedY: 1.0}
	ox, oy := p.Offset(50, 80)
	if ox != 50 || oy != 80 {
		t.Fatalf("expected 50,80 got %d,%d", ox, oy)
	}
}

func TestParallaxConfigDefault(t *testing.T) {
	pc := ParallaxConfig{}
	ox, oy := pc.LayerOffset(0, 100, 200)
	if ox != 100 || oy != 200 {
		t.Fatalf("unconfigured layer should scroll at full speed, got %d,%d", ox, oy)
	}
}
