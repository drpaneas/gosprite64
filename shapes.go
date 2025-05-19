package gosprite64

import (
	"image"
	"image/draw"
	"log"
)

// Rect draws an outline rectangle using two corner points (x1, y1) and (x2, y2).
// The rectangle is drawn using the current cursor color or an optional color index.
//
// Args:
//
//	x1, y1: Coordinates of the first corner
//	x2, y2: Coordinates of the opposite corner
//	colorIndex: Optional PICO-8 color index (0-15)
func Rect(x1, y1, x2, y2 int, colorIndex ...int) {
	drawRect(x1, y1, x2, y2, false, colorIndex...)
}

// Rectfill draws a filled rectangle using two corner points (x1, y1) and (x2, y2).
// The rectangle is filled using the current cursor color or an optional color index.
//
// Args:
//
//	x1, y1: Coordinates of the first corner
//	x2, y2: Coordinates of the opposite corner
//	colorIndex: Optional PICO-8 color index (0-15)
func Rectfill(x1, y1, x2, y2 int, colorIndex ...int) {
	drawRect(x1, y1, x2, y2, true, colorIndex...)
}

// drawRect is a helper function that handles the common logic for both Rect and Rectfill.
// It draws either an outline or filled rectangle depending on the filled parameter.
func drawRect(x1, y1, x2, y2 int, filled bool, colorIndex ...int) {
	// Check if screen is ready
	if currentScreen == nil || currentScreen.Renderer == nil {
		log.Println("Warning: drawRect() called before screen was ready.")
		return
	}

	// Get the color index from arguments or use the current cursor color
	col := cursorColor
	if len(colorIndex) > 0 {
		col = colorIndex[0]
	}

	// Validate color index
	if col < 0 || col >= len(Pico8Palette) {
		log.Printf("Warning: drawRect() called with invalid color index %d. Defaulting to cursorColor (%d).", col, cursorColor)
		col = cursorColor
	}

	// Get the color from the PICO-8 palette
	c := Pico8Palette[col]

	// Ensure x1,y1 is the top-left and x2,y2 is the bottom-right
	minX := min(x1, x2)
	maxX := max(x1, x2)
	minY := min(y1, y2)
	maxY := max(y1, y2)

	// Apply screen offset for overscan
	minX += screenOffsetX
	maxX += screenOffsetX

	if filled {
		// Draw filled rectangle
		rect := image.Rect(minX, minY, maxX+1, maxY+1)
		currentScreen.Renderer.Draw(rect, &image.Uniform{c}, image.Point{}, draw.Src)
	} else {
		// Draw outline rectangle (4 lines)
		// Top
		topRect := image.Rect(minX, minY, maxX+1, minY+1)
		currentScreen.Renderer.Draw(topRect, &image.Uniform{c}, image.Point{}, draw.Src)
		// Bottom
		bottomRect := image.Rect(minX, maxY, maxX+1, maxY+1)
		currentScreen.Renderer.Draw(bottomRect, &image.Uniform{c}, image.Point{}, draw.Src)
		// Left
		leftRect := image.Rect(minX, minY+1, minX+1, maxY)
		currentScreen.Renderer.Draw(leftRect, &image.Uniform{c}, image.Point{}, draw.Src)
		// Right
		rightRect := image.Rect(maxX, minY+1, maxX+1, maxY)
		currentScreen.Renderer.Draw(rightRect, &image.Uniform{c}, image.Point{}, draw.Src)
	}

	// Update the pixel buffer
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			if y >= 0 && y < len(currentScreen.pixels) && x >= 0 && x < len(currentScreen.pixels[y]) {
				if filled || (x == x1 || x == x2 || y == y1 || y == y2) {
					currentScreen.pixels[y][x] = col
				}
			}
		}
	}
}

// min returns the smaller of x or y.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// max returns the larger of x or y.
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
