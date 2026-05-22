package main

import "testing"

func TestCenteredHeroPosition(t *testing.T) {
	x, y := centeredHeroPosition(32, 32)
	if x != 128 {
		t.Fatalf("centeredHeroPosition() x = %v, want 128", x)
	}
	if y != 92 {
		t.Fatalf("centeredHeroPosition() y = %v, want 92", y)
	}
}
