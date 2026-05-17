package gosprite64

import "testing"

func TestFadeToBlackAlpha(t *testing.T) {
	tr := &Transition{Duration: 10, Style: FadeToBlack}
	tr.frame = 0
	if a := tr.alpha(); a != 0 {
		t.Fatalf("frame 0: expected 0, got %d", a)
	}
	tr.frame = 5
	if a := tr.alpha(); a < 120 || a > 135 {
		t.Fatalf("frame 5: expected ~127, got %d", a)
	}
	tr.frame = 10
	if a := tr.alpha(); a != 255 {
		t.Fatalf("frame 10: expected 255, got %d", a)
	}
}

func TestFadeFromBlackAlpha(t *testing.T) {
	tr := &Transition{Duration: 10, Style: FadeFromBlack}
	tr.frame = 0
	if a := tr.alpha(); a != 255 {
		t.Fatalf("frame 0: expected 255, got %d", a)
	}
	tr.frame = 10
	if a := tr.alpha(); a != 0 {
		t.Fatalf("frame 10: expected 0, got %d", a)
	}
}

func TestTransitionDone(t *testing.T) {
	tr := &Transition{Duration: 5, Style: FadeToBlack}
	tr.frame = 4
	if tr.Done() {
		t.Fatal("should not be done at frame 4")
	}
	tr.frame = 5
	if !tr.Done() {
		t.Fatal("should be done at frame 5")
	}
}
