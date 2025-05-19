package gosprite64

import (
	"image"
	"log"
	"math"

	"image/color"

	"github.com/drpaneas/n64/fonts/gomono12"
)

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

	// fontFace is the default monospace font face for text rendering
	fontFace = gomono12.NewFace(gomono12.X0000_00ff())

	// CharWidthApproximation is the approximate width of a character in pixels.
	// This is used for cursor positioning and text measurement.
	CharWidthApproximation = 8.0
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

// Print draws the given string onto the current drawing screen.
// Uses the internal `currentScreen` variable.
// It mimics the PICO-8 PRINT(str, [x, y], [color]) function, including implicit cursor tracking.
// It returns the X and Y coordinates of the pixel immediately following the printed string.
//
// Args:
//   - str: The string to print.
//   - args: Optional arguments interpreted based on PICO-8 logic:
//   - If len(args) == 0: Prints at current cursor (cursorX, cursorY) with current cursorColor.
//   - If len(args) == 1: Prints at current cursor (cursorX, cursorY) with color args[0] (overrides cursorColor).
//   - If len(args) == 2: Prints starting at (args[0], args[1]) with current cursorColor.
//   - If len(args) >= 3: Prints starting at (args[0], args[1]) with color args[2] (overrides cursorColor).
//
// Returns:
//   - int: The X coordinate after the string (drawX + stringWidth).
//   - int: The Y coordinate after the string (drawY + fontHeight).
//
// Example:
//
//	// Assume cursor starts at (0, 0), color is 7 (white)
//	Cursor(0, 0, 6) // Set current color to light gray
//	_, _ = Print("1 HELLO")         // Draws at (0,0) in light gray, cursor moves to (0, 6).
//	_, _ = Print("2 WORLD", 8)      // Draws at (0,6) in red, cursor moves to (0, 12).
//	_, _ = Print("3 AT", 20, 20)     // Draws at (20,20) in light gray, cursor moves to (20, 26).
//	endX, endY := Print("4 DONE")    // Draws at (20, 26) in light gray, cursor moves to (20, 32).
func Print(str string, args ...int) (int, int) {
	// Check if screen is ready
	if currentScreen == nil || currentScreen.Renderer == nil {
		log.Println("Warning: Print() called before screen was ready.")

		// Calculate return values based on arguments without changing cursor state
		posX, posY := cursorX, cursorY
		if len(args) >= 2 {
			posX, posY = args[0], args[1]
		}

		// Approximate measurement for return value
		advance := float64(len([]rune(str))) * CharWidthApproximation
		endX := int(math.Ceil(float64(posX) + advance))
		endY := posY + int(DefaultFontSize)

		return endX, endY
	}

	// Parse arguments
	posX, posY, col := cursorX, cursorY, cursorColor

	// If a new position is provided, override posX and posY
	if len(args) >= 2 {
		posX, posY = args[0], args[1]
	}

	// If a color is provided (in len(args)==1 for color only,
	// or len(args)==3 for position and color), use the last argument.
	if len(args) == 1 || len(args) == 3 {
		col = args[len(args)-1]
		// Update color using the Color function which handles validation
		Color(col)
	}

	// Validate the color index
	if col < 0 || col >= len(Pico8Palette) {
		log.Printf("Warning: Print() called with invalid color index %d. Defaulting to cursorColor (%d).", col, cursorColor)
		col = cursorColor // Default to current cursorColor if invalid index given
	}

	// Calculate text dimensions for cursor positioning
	// Using a more accurate width calculation based on character count and fixed width
	textWidth := len(str) * 6 // Approximate width of monospace characters
	textHeight := int(DefaultFontSize)

	// Create a font face for text rendering
	font := gomono12.NewFace(gomono12.X0000_00ff())

	// Convert PICO-8 color to RGBA
	r, g, b, _ := Pico8Palette[col].RGBA()
	textColor := color.RGBA{
		R: uint8(r >> 8), // Convert from 16-bit to 8-bit color
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: 255, // Full opacity
	}

	// Draw the text using the N64 renderer
	// Apply horizontal offset for overscan and adjust for font metrics
	currentScreen.Renderer.DrawText(
		currentScreen.Renderer.Bounds(),
		font,
		image.Point{X: posX + screenOffsetX, Y: posY + int(DefaultFontSize)},
		textColor,
		nil,
		[]byte(str),
	)

	// Calculate end position for cursor
	endX := posX + textWidth
	endY := posY + textHeight

	// Update cursor position
	// If a position was explicitly provided, use that; otherwise, keep the current cursorX.
	if len(args) >= 2 {
		cursorX = args[0]
	} else {
		cursorX = posX
	}
	cursorY = endY

	return endX, endY
}
