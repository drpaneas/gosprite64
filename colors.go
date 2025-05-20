package gosprite64

import "image/color"

// PICO-8 color indices
const (
	Black      = 0
	DarkBlue   = 1
	DarkPurple = 2
	DarkGreen  = 3
	Brown      = 4
	DarkGray   = 5
	LightGray  = 6
	White      = 7
	Red        = 8
	Orange     = 9
	Yellow     = 10
	Green      = 11
	Blue       = 12
	Indigo     = 13
	Pink       = 14
	Peach      = 15
)

// alpha is the default alpha value for colors.
const alpha uint8 = 255

// Pico8Palette defines the standard 16 PICO-8 colors.
var Pico8Palette = []color.Color{
	color.RGBA{R: 0, G: 0, B: 0, A: alpha},       // 0 black
	color.RGBA{R: 29, G: 43, B: 83, A: alpha},    // 1 dark-blue
	color.RGBA{R: 126, G: 37, B: 83, A: alpha},   // 2 dark-purple
	color.RGBA{R: 0, G: 135, B: 81, A: alpha},    // 3 dark-green
	color.RGBA{R: 171, G: 82, B: 54, A: alpha},   // 4 brown
	color.RGBA{R: 95, G: 87, B: 79, A: alpha},    // 5 dark-gray
	color.RGBA{R: 194, G: 195, B: 199, A: alpha}, // 6 light-gray
	color.RGBA{R: 255, G: 241, B: 232, A: alpha}, // 7 white
	color.RGBA{R: 255, G: 0, B: 77, A: alpha},    // 8 red
	color.RGBA{R: 255, G: 163, B: 0, A: alpha},   // 9 orange
	color.RGBA{R: 255, G: 236, B: 39, A: alpha},  // 10 yellow
	color.RGBA{R: 0, G: 228, B: 54, A: alpha},    // 11 green
	color.RGBA{R: 41, G: 173, B: 255, A: alpha},  // 12 blue
	color.RGBA{R: 131, G: 118, B: 156, A: alpha}, // 13 indigo
	color.RGBA{R: 255, G: 119, B: 168, A: alpha}, // 14 pink
	color.RGBA{R: 255, G: 204, B: 170, A: alpha}, // 15 peach
}

// PaletteTransparency defines which colors in the Pico8Palette should be treated as transparent.
// By default, only color 0 (black) is transparent.
var PaletteTransparency = []bool{true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false}

// Default color index for Cls
const defaultColorIndex = Black

// currentDrawColor holds the current drawing color index (0-15)
var currentDrawColor = defaultColorIndex

// Color sets the current draw color to be used by subsequent drawing operations.
// The color parameter should be a number from 0 to 15 corresponding to the PICO-8 palette.
//
// Example:
//
//	Color(8) // Set current draw color to red (color 8)
//	Pset(10, 20) // Draw a red pixel at (10, 20)
func Color(colorIndex int) {
	// Clamp color index to valid range (0-15)
	if colorIndex < 0 {
		colorIndex = 0
	} else if colorIndex >= len(Pico8Palette) {
		colorIndex = len(Pico8Palette) - 1
	}

	// Update both color variables to keep them in sync
	currentDrawColor = colorIndex
	// cursorColor = colorIndex
}
