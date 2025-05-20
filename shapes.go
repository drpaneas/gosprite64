package gosprite64

import (
	"image"
	"image/draw"
)

// Rectfill draws a filled rectangle with the specified color index.
// It works with any numeric type for coordinates.
func Rectfill[T Number](x1, y1, x2, y2 T, colorIdx int) {
	if currentScreen == nil {
		return
	}

	// Convert coordinates to int
	ix1, iy1, ix2, iy2 := int(x1), int(y1), int(x2), int(y2)

	// Ensure coordinates are in the right order
	if ix1 > ix2 {
		ix1, ix2 = ix2, ix1
	}
	if iy1 > iy2 {
		iy1, iy2 = iy2, iy1
	}

	dstRect := image.Rect(ix1+widthOffset, iy1, ix2+1+widthOffset, iy2+1)
	col := Pico8Palette[colorIdx]
	currentScreen.Renderer.Draw(dstRect, &image.Uniform{C: col}, image.Point{}, draw.Src)
}

// Rect draws an outlined rectangle with the specified color index.
// It works with any numeric type for coordinates.
func Rect[T Number](x1, y1, x2, y2 T, colorIdx int) {
	if currentScreen == nil {
		return
	}

	// Convert coordinates to int for comparison
	ix1, iy1, ix2, iy2 := int(x1), int(y1), int(x2), int(y2)

	// Ensure coordinates are in the right order
	if ix1 > ix2 {
		ix1, ix2 = ix2, ix1
	}
	if iy1 > iy2 {
		iy1, iy2 = iy2, iy1
	}

	// Convert back to type T for the function call
	x1, y1, x2, y2 = T(ix1), T(iy1), T(ix2), T(iy2)

	// Draw the outline using Rectfill with type parameters
	Rectfill(x1, y1, x2, y1, colorIdx)     // Top
	Rectfill(x1, y2, x2, y2, colorIdx)     // Bottom
	Rectfill(x1, y1+1, x1, y2-1, colorIdx) // Left
	Rectfill(x2, y1+1, x2, y2-1, colorIdx) // Right
}
