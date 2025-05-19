package gosprite64

import (
	"image"
	"image/draw"
	"log"
	"math"
)

// abs returns the absolute value of x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// toInterfaceSlice converts a slice of any type to []interface{}
func toInterfaceSlice(slice interface{}) []interface{} {
	s := make([]interface{}, 0)
	switch v := slice.(type) {
	case []int:
		for _, val := range v {
			s = append(s, val)
		}
	case []float64:
		for _, val := range v {
			s = append(s, val)
		}
	}
	return s
}

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

	// Ensure x1,y1 is the top-left and x2,y2 is the bottom-right
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	// Apply screen offset for overscan
	x1 += screenOffsetX
	x2 += screenOffsetX

	// Update the pixel buffer for Pget()
	if filled {
		// Draw filled rectangle
		for y := y1; y <= y2; y++ {
			for x := x1; x <= x2; x++ {
				setPixel(x, y, col)
			}
		}
	} else {
		// Draw outline rectangle
		// Top and bottom lines
		for x := x1; x <= x2; x++ {
			setPixel(x, y1, col)
			setPixel(x, y2, col)
		}
		// Left and right lines (avoiding double-drawing corners)
		for y := y1 + 1; y < y2; y++ {
			setPixel(x1, y, col)
			setPixel(x2, y, col)
		}
	}
}

// setPixel sets a pixel at (x,y) with the specified color index
func setPixel(x, y, colorIndex int) {
	if currentScreen == nil || currentScreen.Renderer == nil || currentScreen.pixels == nil {
		return
	}

	// Apply screen offset for overscan
	screenX := x + screenOffsetX

	// Check bounds
	if screenX >= 0 && screenX < len(currentScreen.pixels[0]) && y >= 0 && y < len(currentScreen.pixels) {
		// Update the pixel buffer
		currentScreen.pixels[y][screenX] = colorIndex

		// Draw the pixel using the renderer
		c := Pico8Palette[colorIndex]
		currentScreen.Renderer.Draw(
			image.Rect(screenX, y, screenX+1, y+1),
			&image.Uniform{c},
			image.Point{},
			draw.Src,
		)
	}
}

// parseLineArgs parses color index arguments for drawing functions.
// It returns the color index to use and whether parsing was successful.
func parseLineArgs(args []interface{}) (int, bool) {
	if len(args) == 0 {
		return cursorColor, true
	}

	switch v := args[0].(type) {
	case int:
		if v >= 0 && v < len(Pico8Palette) {
			return v, true
		}
		log.Printf("Warning: Invalid color index %d. Using cursor color %d.", v, cursorColor)
		return cursorColor, false
	case float64:
		intVal := int(v)
		if intVal >= 0 && intVal < len(Pico8Palette) {
			return intVal, true
		}
		log.Printf("Warning: Invalid color index %d. Using cursor color %d.", intVal, cursorColor)
		return cursorColor, false
	default:
		log.Printf("Warning: Invalid color index type %T. Using cursor color %d.", v, cursorColor)
		return cursorColor, false
	}
}

// parseCircArgs parses common arguments for Circ and Circfill.
// It returns the center coordinates (x, y), radius, the PICO-8 color index to use,
// and whether parsing was successful.
func parseCircArgs(x, y, radius float64, options []interface{}) (float64, float64, float64, int, bool) {
	// Determine drawing color
	drawColorIndex := cursorColor // Use the current cursor color
	if len(options) >= 1 {
		switch v := options[0].(type) {
		case int:
			if v >= 0 && v < len(Pico8Palette) {
				drawColorIndex = v
			} else {
				log.Printf("Warning: Circ/Circfill called with invalid color index %d. Using current color %d.", v, cursorColor)
			}
		case float64:
			intVal := int(v)
			if intVal >= 0 && intVal < len(Pico8Palette) {
				drawColorIndex = intVal
			} else {
				log.Printf("Warning: Circ/Circfill called with invalid color index %d. Using current color %d.", intVal, cursorColor)
			}
		}
	}

	return x, y, radius, drawColorIndex, true
}

// drawCirclePoints draws the 8 symmetric points of a circle
func drawCirclePoints(cx, cy, x, y int, filled bool, colorIdx int) {
	if filled {
		// Draw horizontal lines between points at the same y-level
		drawHorizontalLine(cx-x, cx+x, cy+y, colorIdx)
		if y != 0 {
			drawHorizontalLine(cx-x, cx+x, cy-y, colorIdx)
		}
	} else {
		// Draw the 8 symmetric points
		setPixel(cx+x, cy+y, colorIdx)
		setPixel(cx-x, cy+y, colorIdx)
		setPixel(cx+x, cy-y, colorIdx)
		setPixel(cx-x, cy-y, colorIdx)
		setPixel(cx+y, cy+x, colorIdx)
		setPixel(cx-y, cy+x, colorIdx)
		setPixel(cx+y, cy-x, colorIdx)
		setPixel(cx-y, cy-x, colorIdx)
	}
}

// drawHorizontalLine draws a horizontal line from x1 to x2 at y
func drawHorizontalLine(x1, x2, y, colorIdx int) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	for x := x1; x <= x2; x++ {
		setPixel(x, y, colorIdx)
	}
}

// Circ draws an outline circle.
//
// Args:
//
//	x, y: Coordinates of the center point (any Number type)
//	radius: Radius of the circle (any Number type)
//	colorIndex: Optional PICO-8 color index (0-15)
func Circ[X, Y, R Number](x X, y Y, radius R, colorIndex ...int) {
	// Check if screen is ready
	if currentScreen == nil || currentScreen.Renderer == nil {
		log.Println("Warning: Circ() called before screen was ready.")
		return
	}

	// Convert to float64 for calculations
	fx, fy, fradius := float64(x), float64(y), float64(radius)


	// Convert colorIndex to []interface{}
	var colorArgs []interface{}
	for _, c := range colorIndex {
		colorArgs = append(colorArgs, c)
	}

	// Parse optional color argument
	drawColorIndex, ok := parseLineArgs(colorArgs)
	if !ok {
		return // Error already logged
	}

	// Convert to integers for pixel-perfect drawing
	ix, iy, iradius := int(math.Round(fx)), int(math.Round(fy)), int(math.Round(fradius))

	// Apply screen offset for overscan
	ix += screenOffsetX

	// Use Midpoint Circle Algorithm to draw the circle
	x1, y1 := iradius, 0
	dp1 := 3 - 2*iradius

	for x1 >= y1 {
		drawCirclePoints(ix, iy, x1, y1, false, drawColorIndex)
		y1++
		if dp1 <= 0 {
			dp1 = dp1 + 4*y1 + 6
		} else {
			x1--
			dp1 = dp1 + 4*(y1-x1) + 10
		}
	}
}

// Circfill draws a filled circle.
//
// Args:
//
//	x, y: Coordinates of the center point (any Number type)
//	radius: Radius of the circle (any Number type)
//	colorIndex: Optional PICO-8 color index (0-15)
func Circfill[X, Y, R Number](x X, y Y, radius R, colorIndex ...int) {
	// Check if screen is ready
	if currentScreen == nil || currentScreen.Renderer == nil {
		log.Println("Warning: Circfill() called before screen was ready.")
		return
	}

	// Convert to float64 for calculations
	fx, fy, fradius := float64(x), float64(y), float64(radius)

	// Convert colorIndex to []interface{}
	var colorArgs []interface{}
	for _, c := range colorIndex {
		colorArgs = append(colorArgs, c)
	}

	// Parse optional color argument
	drawColorIndex, ok := parseLineArgs(colorArgs)
	if !ok {
		return // Error already logged
	}

	// Convert to integers for pixel-perfect drawing
	ix, iy, iradius := int(math.Round(fx)), int(math.Round(fy)), int(math.Round(fradius))

	// Apply screen offset for overscan
	ix += screenOffsetX

	// Use Midpoint Circle Algorithm to draw the filled circle
	x1, y1 := iradius, 0
	dp1 := 3 - 2*iradius

	for x1 >= y1 {
		drawCirclePoints(ix, iy, x1, y1, true, drawColorIndex)
		y1++
		if dp1 <= 0 {
			dp1 = dp1 + 4*y1 + 6
		} else {
			x1--
			dp1 = dp1 + 4*(y1-x1) + 10
		}
	}
}

// Line draws a line between two points using hardware acceleration.
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

	// Parse options
	colorIdx, ok := parseLineArgs(options)
	if !ok {
		return
	}

	// Convert to integer coordinates
	ix1, iy1 := int(math.Round(float64(x1))), int(math.Round(float64(y1)))
	ix2, iy2 := int(math.Round(float64(x2))), int(math.Round(float64(y2)))

	// Apply screen offset for overscan
	ix1 += screenOffsetX
	ix2 += screenOffsetX

	// For now, we'll use the pixel-by-pixel approach for lines
	// This can be optimized further with a proper line drawing algorithm
	// that works with the N64's RDP
	updateLineInPixelBuffer(ix1, iy1, ix2, iy2, colorIdx)
}

// updateLineInPixelBuffer updates the pixel buffer for a line using Bresenham's algorithm
func updateLineInPixelBuffer(x1, y1, x2, y2, colorIdx int) {
	// This is a simple implementation of Bresenham's line algorithm
	dx := abs(x2 - x1)
	dy := -abs(y2 - y1)
	sx, sy := 1, 1
	if x1 > x2 {
		sx = -1
	}
	if y1 > y2 {
		sy = -1
	}
	err := dx + dy // error value e_xy

	for {
		setPixel(x1, y1, colorIdx)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 >= dy { // e_xy + e_x > 0
			err += dy
			x1 += sx
		}
		if e2 <= dx { // e_xy + e_y < 0
			err += dx
			y1 += sy
		}
	}
}
