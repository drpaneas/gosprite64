package gosprite64

import (
	"image"
	"image/color"
	"image/draw"
	"log"
)

// DrawRect draws the outline of a rectangle using DrawLine.
// x1, y1: Top-left corner (0-319, 0-239)
// x2, y2: Bottom-right corner (inclusive)
// c: The color of the rectangle outline
func DrawRect(x1, y1, x2, y2 int, c color.Color) {
	// Ensure x1 <= x2 and y1 <= y2
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	// Draw top line
	Line(x1, y1, x2, y1, c)
	// Draw bottom line
	Line(x1, y2, x2, y2, c)
	// Draw left line (offset by 1 to avoid double-drawing corners)
	if y2-y1 > 1 {
		Line(x1, y1+1, x1, y2-1, c)
	}
	// Draw right line (offset by 1 to avoid double-drawing corners)
	if y2-y1 > 1 {
		Line(x2, y1+1, x2, y2-1, c)
	}
}

// Line draws a line from (x1,y1) to (x2,y2) using DrawRectFill.
// The line will be 1 pixel thick.
// Coordinates must be within the screen bounds (0-319, 0-239).
func Line(x1, y1, x2, y2 int, c color.Color) {
	// Check bounds
	if x1 < 0 || x1 > 319 || y1 < 0 || y1 > 239 || x2 < 0 || x2 > 319 || y2 < 0 || y2 > 239 {
		log.Printf("DrawLine: (%d,%d) to (%d,%d) out of bounds", x1, y1, x2, y2)
		return
	}

	// For horizontal lines
	if y1 == y2 {
		if x1 > x2 {
			x1, x2 = x2, x1 // Ensure x1 <= x2
		}
		DrawRectFill(x1, y1, x2, y1, c)
		return
	}

	// For vertical lines
	if x1 == x2 {
		if y1 > y2 {
			y1, y2 = y2, y1 // Ensure y1 <= y2
		}
		DrawRectFill(x1, y1, x1, y2, c)
		return
	}

	// For diagonal lines, we'll draw a series of 1x1 rectangles
	// This is a simple implementation using Bresenham's algorithm
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := 1
	if x1 > x2 {
		sx = -1
	}
	sy := 1
	if y1 > y2 {
		sy = -1
	}
	err := dx - dy

	for {
		DrawRectFill(x1, y1, x1, y1, c)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// abs returns the absolute value of x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Rectfill draws a filled rectangle on screen
// x1, y1: Top-left corner (0-319, 0-239)
// x2, y2: Bottom-right corner (inclusive)
// color: The color to fill the rectangle with
func Rectfill(x1, y1, x2, y2 int, color color.Color) {
	if currentScreen == nil || currentScreen.Renderer == nil {
		return
	}

	// Ensure x1 <= x2 and y1 <= y2
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	// For NTSC, the display is typically 640x480 but the visible area is smaller
	// We'll work in the visible area (approximately 320x240)
	const (
		screenWidth  = 320
		screenHeight = 240
	)

	// Clip to screen bounds
	if x1 < 0 {
		x1 = 0
	}
	if y1 < 0 {
		y1 = 0
	}
	if x2 >= screenWidth {
		x2 = screenWidth - 1
	}
	if y2 >= screenHeight {
		y2 = screenHeight - 1
	}

	// Create and draw the rectangle
	// Note: image.Rect is half-open interval [Min, Max)
	rect := image.Rect(x1, y1, x2+1, y2+1)
	img := image.NewUniform(color)
	currentScreen.Renderer.Draw(rect, img, image.Point{}, draw.Src)
}

// // Rect draws an outlined rectangle with the specified color index.
// // It works with any numeric type for coordinates.
// func Rect[T Number](x1, y1, x2, y2 T, colorIdx int) {
// 	if currentScreen == nil {
// 		return
// 	}

// 	// Convert coordinates to int for comparison
// 	ix1, iy1, ix2, iy2 := int(x1), int(y1), int(x2), int(y2)

// 	// Ensure coordinates are in the right order
// 	if ix1 > ix2 {
// 		ix1, ix2 = ix2, ix1
// 	}
// 	if iy1 > iy2 {
// 		iy1, iy2 = iy2, iy1
// 	}

// 	// Convert back to type T for the function call
// 	x1, y1, x2, y2 = T(ix1), T(iy1), T(ix2), T(iy2)

// 	// Draw the outline using Rectfill with type parameters
// 	Rectfill(x1, y1, x2, y1, colorIdx)         // Top
// 	Rectfill(x1, y2, x2, y2, colorIdx)         // Bottom
// 	Rectfill(x1+1, y1+1, x1+1, y2-1, colorIdx) // Left
// 	Rectfill(x2-1, y1+1, x2-1, y2-1, colorIdx) // Right
// }

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	SafeBorder   = 80                          // 80 pixels on each side
	SafeWidth    = ScreenWidth - 2*SafeBorder  // 480
	SafeHeight   = ScreenHeight - 2*SafeBorder // 320
)

// DrawBorder draws a border around the safe area for debugging
func (s *screen) DrawBorder() {
	// Outer border (red)
	Rectfill(0, 0, ScreenWidth-1, ScreenHeight-1, color.RGBA{255, 0, 0, 255})
	// Inner safe area (black)
	Rectfill(
		SafeBorder,
		SafeBorder,
		ScreenWidth-SafeBorder-1,
		ScreenHeight-SafeBorder-1,
		color.Black,
	)
}

// // InSafeArea checks if coordinates are within the safe area
// func (s *screen) InSafeArea(x, y int) bool {
// 	return x >= SafeBorder &&
// 		x < ScreenWidth-SafeBorder &&
// 		y >= SafeBorder &&
// 		y < ScreenHeight-SafeBorder
// }
