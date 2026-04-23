package rendergeom

import (
	"image"
	"testing"
)

func TestProfileAccessors(t *testing.T) {
	if got := LogicalBounds(); got != image.Rect(0, 0, 288, 216) {
		t.Fatalf("LogicalBounds() = %v, want %v", got, image.Rect(0, 0, 288, 216))
	}

	if got := FramebufferBounds(); got != image.Rect(0, 0, 320, 240) {
		t.Fatalf("FramebufferBounds() = %v, want %v", got, image.Rect(0, 0, 320, 240))
	}

	if got := Origin(); got != image.Pt(16, 12) {
		t.Fatalf("Origin() = %v, want %v", got, image.Pt(16, 12))
	}
}

func TestMapPoint(t *testing.T) {
	tests := []struct {
		name    string
		logical image.Point
		want    image.Point
	}{
		{
			name:    "top left logical pixel",
			logical: image.Pt(0, 0),
			want:    image.Pt(16, 12),
		},
		{
			name:    "bottom right logical pixel",
			logical: image.Pt(287, 215),
			want:    image.Pt(303, 227),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := MapPoint(tt.logical)
			if !ok {
				t.Fatalf("MapPoint(%v) reported out of bounds", tt.logical)
			}
			if got != tt.want {
				t.Fatalf("MapPoint(%v) = %v, want %v", tt.logical, got, tt.want)
			}
		})
	}
}

func TestMapPointRejectsOutOfBounds(t *testing.T) {
	tests := []image.Point{
		image.Pt(-1, 0),
		image.Pt(0, -1),
		image.Pt(288, 0),
		image.Pt(0, 216),
	}

	for _, logical := range tests {
		if _, ok := MapPoint(logical); ok {
			t.Fatalf("MapPoint(%v) reported in bounds", logical)
		}
	}
}

func TestMapRectInclusive(t *testing.T) {
	got, ok := MapRectInclusive(image.Rectangle{
		Min: image.Pt(0, 0),
		Max: image.Pt(287, 215),
	})
	if !ok {
		t.Fatal("MapRectInclusive(full logical bounds) reported out of bounds")
	}

	want := image.Rectangle{
		Min: image.Pt(16, 12),
		Max: image.Pt(303, 227),
	}
	if got != want {
		t.Fatalf("MapRectInclusive(full logical bounds) = %v, want %v", got, want)
	}
}

func TestMapRectInclusiveSinglePixel(t *testing.T) {
	got, ok := MapRectInclusive(image.Rectangle{
		Min: image.Pt(287, 215),
		Max: image.Pt(287, 215),
	})
	if !ok {
		t.Fatal("MapRectInclusive(single pixel) reported out of bounds")
	}

	want := image.Rectangle{
		Min: image.Pt(303, 227),
		Max: image.Pt(303, 227),
	}
	if got != want {
		t.Fatalf("MapRectInclusive(single pixel) = %v, want %v", got, want)
	}
}

func TestMapRectInclusiveClipsNegativeInput(t *testing.T) {
	got, ok := MapRectInclusive(image.Rectangle{
		Min: image.Pt(-8, -4),
		Max: image.Pt(10, 20),
	})
	if !ok {
		t.Fatal("MapRectInclusive(negative overlap) reported out of bounds")
	}

	want := image.Rectangle{
		Min: image.Pt(16, 12),
		Max: image.Pt(26, 32),
	}
	if got != want {
		t.Fatalf("MapRectInclusive(negative overlap) = %v, want %v", got, want)
	}
}

func TestMapRectInclusiveClipsOverflowingInput(t *testing.T) {
	got, ok := MapRectInclusive(image.Rectangle{
		Min: image.Pt(280, 210),
		Max: image.Pt(400, 300),
	})
	if !ok {
		t.Fatal("MapRectInclusive(overflow overlap) reported out of bounds")
	}

	want := image.Rectangle{
		Min: image.Pt(296, 222),
		Max: image.Pt(303, 227),
	}
	if got != want {
		t.Fatalf("MapRectInclusive(overflow overlap) = %v, want %v", got, want)
	}
}

func TestMapRectInclusiveRejectsFullyOutOfBounds(t *testing.T) {
	tests := []image.Rectangle{
		{
			Min: image.Pt(-10, 0),
			Max: image.Pt(-1, 5),
		},
		{
			Min: image.Pt(0, -10),
			Max: image.Pt(5, -1),
		},
		{
			Min: image.Pt(288, 0),
			Max: image.Pt(290, 5),
		},
		{
			Min: image.Pt(0, 216),
			Max: image.Pt(10, 220),
		},
		{
			Min: image.Pt(10, 10),
			Max: image.Pt(5, 5),
		},
	}

	for _, rect := range tests {
		if _, ok := MapRectInclusive(rect); ok {
			t.Fatalf("MapRectInclusive(%v) reported in bounds", rect)
		}
	}
}

func TestCenteredRect(t *testing.T) {
	tests := []struct {
		name   string
		bounds image.Rectangle
		size   image.Point
		want   image.Rectangle
	}{
		{
			name:   "ntsc sized output stays aligned",
			bounds: image.Rect(108, 35, 108+640, 35+480),
			size:   image.Pt(640, 480),
			want:   image.Rect(108, 35, 108+640, 35+480),
		},
		{
			name:   "pal centers shorter output vertically",
			bounds: image.Rect(128, 45, 128+640, 45+576),
			size:   image.Pt(640, 480),
			want:   image.Rect(128, 93, 128+640, 93+480),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CenteredRect(tt.bounds, tt.size); got != tt.want {
				t.Fatalf("CenteredRect(%v, %v) = %v, want %v", tt.bounds, tt.size, got, tt.want)
			}
		})
	}
}
