package gosprite64

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
)

// Number is a constraint that permits any numeric type.
// This includes both integer and floating-point types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Rect draws an outline rectangle using two corner points.
// The rectangle is drawn using the current cursor color or an optional color index.
//
// Args:
//
//	x1, y1, x2, y2: Coordinates of two opposing corners (any Number type)
//	colorIndex: Optional PICO-8 color index (0-15)
func Rect[X1, Y1, X2, Y2 Number](x1 X1, y1 Y1, x2 X2, y2 Y2, colorIndex ...int) {
	drawRect(float64(x1), float64(y1), float64(x2), float64(y2), false, colorIndex...)
}

// Rectfill draws a filled rectangle using two corner points.
// The rectangle is filled using the current cursor color or an optional color index.
//
// Args:
//
//	x1, y1, x2, y2: Coordinates of two opposing corners (any Number type)
//	colorIndex: Optional PICO-8 color index (0-15)
func Rectfill[X1, Y1, X2, Y2 Number](x1 X1, y1 Y1, x2 X2, y2 Y2, colorIndex ...int) {
	drawRect(float64(x1), float64(y1), float64(x2), float64(y2), true, colorIndex...)
}

// drawRect is a helper function that handles the common logic for both Rect and Rectfill.
// It draws either an outline or filled rectangle depending on the filled parameter.
func drawRect(fx1, fy1, fx2, fy2 float64, filled bool, colorIndex ...int) {
	// Convert to integers for pixel-perfect drawing
	x1, y1 := int(math.Round(fx1)), int(math.Round(fy1))
	x2, y2 := int(math.Round(fx2)), int(math.Round(fy2))
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
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	// Apply screen offset for overscan in the drawing loop

	if filled {
		// Draw filled rectangle
		for y := y1; y <= y2; y++ {
			for x := x1; x <= x2; x++ {
				screenX := x + screenOffsetX
				if screenX >= 0 && screenX < len(currentScreen.pixels[0]) && y >= 0 && y < len(currentScreen.pixels) {
					currentScreen.pixels[y][screenX] = col
					currentScreen.Renderer.Draw(
						image.Rect(screenX, y, screenX+1, y+1),
						&image.Uniform{c},
						image.Point{},
						draw.Src,
					)
				}
			}
		}
	} else {
		// Draw outline rectangle
		// Top and bottom lines
		for x := x1; x <= x2; x++ {
			screenX := x + screenOffsetX
			// Top line
			if screenX >= 0 && screenX < len(currentScreen.pixels[0]) && y1 >= 0 && y1 < len(currentScreen.pixels) {
				currentScreen.pixels[y1][screenX] = col
				currentScreen.Renderer.Draw(
					image.Rect(screenX, y1, screenX+1, y1+1),
					&image.Uniform{c},
					image.Point{},
					draw.Src,
				)
			}
			// Bottom line
			if screenX >= 0 && screenX < len(currentScreen.pixels[0]) && y2 >= 0 && y2 < len(currentScreen.pixels) {
				currentScreen.pixels[y2][screenX] = col
				currentScreen.Renderer.Draw(
					image.Rect(screenX, y2, screenX+1, y2+1),
					&image.Uniform{c},
					image.Point{},
					draw.Src,
				)
			}
		}
		// Left and right lines (without corners to avoid double-drawing)
		for y := y1 + 1; y < y2; y++ {
			// Left line
			screenX := x1 + screenOffsetX
			if screenX >= 0 && screenX < len(currentScreen.pixels[0]) && y >= 0 && y < len(currentScreen.pixels) {
				currentScreen.pixels[y][screenX] = col
				currentScreen.Renderer.Draw(
					image.Rect(screenX, y, screenX+1, y+1),
					&image.Uniform{c},
					image.Point{},
					draw.Src,
				)
			}
			// Right line
			screenX = x2 + screenOffsetX
			if screenX >= 0 && screenX < len(currentScreen.pixels[0]) && y >= 0 && y < len(currentScreen.pixels) {
				currentScreen.pixels[y][screenX] = col
				currentScreen.Renderer.Draw(
					image.Rect(screenX, y, screenX+1, y+1),
					&image.Uniform{c},
					image.Point{},
					draw.Src,
				)
			}
		}
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

// parseLineArgs parses common arguments for Line function.
// It returns the PICO-8 color index to use and whether parsing was successful.
func parseLineArgs(options []interface{}) (int, bool) {
	// Use the current cursor color by default
	drawColorIndex := cursorColor

	// If color index is provided, use it
	if len(options) > 0 {
		switch v := options[0].(type) {
		case int:
			if v >= 0 && v < len(Pico8Palette) {
				drawColorIndex = v
			} else {
				log.Printf("Warning: Line() called with invalid color index %d. Using current color %d.", v, cursorColor)
			}
		case float64:
			intVal := int(v)
			if intVal >= 0 && intVal < len(Pico8Palette) {
				drawColorIndex = intVal
			} else {
				log.Printf("Warning: Line() called with invalid color index %d. Using current color %d.", intVal, cursorColor)
			}
		case float32:
			intVal := int(v)
			if intVal >= 0 && intVal < len(Pico8Palette) {
				drawColorIndex = intVal
			} else {
				log.Printf("Warning: Line() called with invalid color index %d. Using current color %d.", intVal, cursorColor)
			}
		default:
			log.Printf("Warning: Line() called with invalid color type %T. Using current color %d.", options[0], cursorColor)
		}
	}

	if len(options) > 1 {
		log.Printf("Warning: Line() called with too many arguments (%d), expected max 5.", len(options)+4)
	}

	return drawColorIndex, true
}

// Line draws a line between two points.
// The line is drawn using the current cursor color or an optional color index.
//
// Args:
//
//	x1, y1: Coordinates of the starting point (any Number type)
//	x2, y2: Coordinates of the ending point (any Number type)
//	options...:
//	  - color (int): Optional PICO-8 color index (0-15). If omitted or invalid,
//	    uses the current cursor color.
func Line[X1, Y1, X2, Y2 Number](x1 X1, y1 Y1, x2 X2, y2 Y2, options ...interface{}) {
	// Check if screen is ready
	if currentScreen == nil || currentScreen.Renderer == nil {
		log.Println("Warning: Line() called before screen was ready.")
		return
	}

	// Convert to float64 for calculations
	fx1, fy1, fx2, fy2 := float64(x1), float64(y1), float64(x2), float64(y2)

	// Parse optional color argument
	drawColorIndex, ok := parseLineArgs(options)
	if !ok {
		return // Error already logged
	}

	// Get the actual color from the palette
	var actualColor color.Color
	if drawColorIndex >= 0 && drawColorIndex < len(Pico8Palette) {
		actualColor = Pico8Palette[drawColorIndex]
	} else {
		actualColor = Pico8Palette[0] // Fallback to black
		log.Printf("Error: Invalid effective drawing color index %d for Line(). Defaulting to black.", drawColorIndex)
	}

	// Convert to integers for pixel-perfect drawing
	ix1, iy1, ix2, iy2 := int(math.Round(fx1)), int(math.Round(fy1)), int(math.Round(fx2)), int(math.Round(fy2))

	// Use Bresenham's line algorithm to draw the line
	dx := int(math.Abs(float64(ix2 - ix1)))
	dy := int(math.Abs(float64(iy2 - iy1)))
	sx, sy := 1, 1

	if ix1 > ix2 {
		sx = -1
	}
	if iy1 > iy2 {
		sy = -1
	}
	err := dx - dy

	for {
		// Draw the current pixel if it's within bounds
		if ix1 >= 0 && ix1 < len(currentScreen.pixels[0]) && iy1 >= 0 && iy1 < len(currentScreen.pixels) {
			currentScreen.pixels[iy1][ix1] = drawColorIndex
		}

		// Check if we've reached the end point
		if ix1 == ix2 && iy1 == iy2 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			ix1 += sx
		}
		if e2 < dx {
			err += dx
			iy1 += sy
		}

		// Draw the actual pixel on screen
		screenX := ix1 + screenOffsetX
		if screenX >= 0 && screenX < len(currentScreen.pixels[0]) && iy1 >= 0 && iy1 < len(currentScreen.pixels) {
			currentScreen.Renderer.Draw(
				image.Rect(screenX, iy1, screenX+1, iy1+1),
				&image.Uniform{actualColor},
				image.Point{},
				draw.Src,
			)
		}
	}
}
