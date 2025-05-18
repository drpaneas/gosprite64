package gosprite64

import "log"

// DefaultFontSize is the default size used for the Print function.
// PICO-8 font is typically 6px high.
const DefaultFontSize = 6.0

// Default values for cursor and display
const (
	defaultCursorX     = 0
	defaultCursorY     = 0
	defaultCursorColor = White
)

// These variables hold the internal state for the cursor used by Print.
var (
	cursorX     = defaultCursorX
	cursorY     = defaultCursorY
	cursorColor = defaultCursorColor
)

// Cursor sets the implicit print cursor position (x, y) and optionally the default draw color.
// It mimics the PICO-8 CURSOR(x, y, [color]) function.
// Calling Cursor() with no arguments resets the cursor position to (0, 0) but leaves the color unchanged.
//
// Args:
//   - args: Optional arguments interpreted as [x, y] or [x, y, colorIndex].
//   - If len(args) == 0: Resets cursor position to (0, 0).
//   - If len(args) == 2: Sets cursor position to (args[0], args[1]).
//   - If len(args) >= 3: Sets cursor position to (args[0], args[1]) and sets currentDrawColor to args[2].
//
// Example:
//
//	Cursor(10, 20)     // Set cursor to (10, 20)
//	Cursor(30, 40, 5) // Set cursor to (30, 40) and draw color to 5 (dark gray)
//	Cursor()          // Reset cursor position to (0, 0)
func Cursor(args ...int) {
	switch len(args) {
	case 0:
		cursorX = defaultCursorX
		cursorY = defaultCursorY
	case 2:
		cursorX = args[0]
		cursorY = args[1]
	case 3:
		cursorX = args[0]
		cursorY = args[1]
		// Set color using the Color function which handles validation
		Color(args[2])
	default:
		log.Printf("Warning: Cursor() called with invalid number of arguments (%d). Expected 0, 2, or 3.", len(args))
	}
}
