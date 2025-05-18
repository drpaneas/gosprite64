// Package gosprite64 provides a simple 2D graphics library for Nintendo 64 development.
// It offers a high-level API for drawing graphics, handling the low-level details
// of the N64's Reality Display Processor (RDP) and video initialization.
package gosprite64

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/drpaneas/n64/drivers/display"
	n64draw "github.com/drpaneas/n64/drivers/draw"
	"github.com/drpaneas/n64/machine"
	"github.com/drpaneas/n64/rcp/video"
)

// VideoPreset represents a predefined video configuration.
type VideoPreset int

const (
	// LowRes is the most common setup: 320x240 without interlacing.
	LowRes VideoPreset = iota
	// HighRes is 640x480 with interlacing.
	HighRes
)

const (
	lowResWidth   = 320
	lowResHeight  = 240
	highResWidth  = 640
	highResHeight = 480
)

// screen holds all the display-related objects and state.
type screen struct {
	Display  *display.Display
	Renderer *n64draw.Rdp
	width    int
	height   int
	pixels   [][]int // 2D array to store color indices
}

// newScreen creates a new screen instance with the given dimensions
func newScreen(disp *display.Display, renderer *n64draw.Rdp, width, height int) *screen {
	// Initialize the pixel buffer with -1 (transparent/undefined)
	pixels := make([][]int, height)
	for i := range pixels {
		pixels[i] = make([]int, width)
		for j := range pixels[i] {
			pixels[i][j] = defaultColorIndex // -1 represents transparent/undefined
		}
	}

	log.Printf("Screen initialized with %d x %d pixels", width, height)

	return &screen{
		Display:  disp,
		Renderer: renderer,
		width:    width,
		height:   height,
		pixels:   pixels,
	}
}

// currentScreen holds the active screen instance.
var currentScreen *screen

// Config holds the configuration for initializing the display.
type Config struct {
	Preset VideoPreset
}

// DefaultConfig returns a default configuration.
func DefaultConfig() Config {
	return Config{
		Preset: LowRes,
	}
}

// getPresetConfig returns the configuration for a given preset.
func getPresetConfig(preset VideoPreset) (image.Point, video.ColorDepth, machine.VideoType, bool) {
	switch preset {
	case HighRes:
		return image.Point{X: highResWidth, Y: highResHeight},
			video.ColorDepth(video.BPP32),
			machine.VideoNTSC,
			true
	default: // LowRes or unknown
		return image.Point{X: lowResWidth, Y: lowResHeight},
			video.ColorDepth(video.BPP16),
			machine.VideoNTSC,
			false
	}
}

// videoInit initializes the display with the specified video preset.
// It sets up the framebuffer and renderer for drawing.
// Returns a new screen instance that can be used for drawing operations.
func videoInit(preset VideoPreset) {
	resolution, colorDepth, mode, interlaced := getPresetConfig(preset)

	// Set the video mode and setup
	machine.Video = mode
	video.Setup(interlaced)

	// Create display and renderer
	disp := display.NewDisplay(resolution, colorDepth)
	renderer := n64draw.NewRdp()
	renderer.SetFramebuffer(disp.Swap())

	// Create screen with our custom implementation
	currentScreen = newScreen(disp, renderer, resolution.X, resolution.Y)
}

// beginDrawing prepares for a new frame by swapping the framebuffer.
// This is an internal function used by the package.
func beginDrawing() {
	if currentScreen == nil || currentScreen.Display == nil {
		log.Println("Warning: beginDrawing called before screen was ready.")
		return
	}
	fb := currentScreen.Display.Swap()
	currentScreen.Renderer.SetFramebuffer(fb)
}

// endDrawing finalizes the frame by flushing the renderer.
// This is an internal function used by the package.
func endDrawing() {
	if currentScreen != nil && currentScreen.Renderer != nil {
		currentScreen.Renderer.Flush()
	}
}

// fill fills the entire screen with the specified color.
func (s *screen) fill(c color.Color) {
	s.Renderer.Draw(s.Renderer.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)
}

// ClearScreen clears the current drawing screen with a specified PICO-8 color index.
// If no colorIndex is provided, it defaults to Pico8Black (0).
func ClearScreen(colorIndex ...int) {
	if currentScreen == nil {
		log.Println("Warning: Cls() called before screen was ready.")
		return
	}

	idx := defaultColorIndex
	if len(colorIndex) > 0 {
		idx = colorIndex[0]
	}

	if idx < 0 || idx >= len(Pico8Palette) {
		log.Printf("Warning: Cls() called with invalid color index %d. Defaulting to %d.", idx, defaultColorIndex)
		idx = defaultColorIndex
	}

	// Clear our pixel buffer
	for y := range currentScreen.pixels {
		for x := range currentScreen.pixels[y] {
			currentScreen.pixels[y][x] = idx
		}
	}

	// Clear the screen
	currentScreen.fill(Pico8Palette[idx])

	// Reset the global print cursor position
	cursorX = defaultCursorX
	cursorY = defaultCursorY
}

// Pget returns the PICO-8 color index (0-15) of the pixel at coordinates (x, y)
// on the current drawing screen.
//
// If the coordinates are outside the screen bounds, it returns 0 (black).
// If the pixel has not been set (is transparent), it returns 0.
//
// Example:
//
//	// Set pixel at (10, 20) to red (index 8)
//	Pset(10, 20, 8)
//
//	// Get the color index of the pixel we just set
//	pixelColorIndex := Pget(10, 20) // Returns 8 (red)
func Pget(x, y int) int {
	if currentScreen == nil {
		log.Println("Warning: Pget() called before screen was ready.")
		return 0
	}

	// Check bounds
	if x < 0 || x >= currentScreen.width || y < 0 || y >= currentScreen.height {
		return 0 // PICO-8 pget returns 0 for out-of-bounds
	}

	// Return the color index from our buffer
	colorIdx := currentScreen.pixels[y][x]
	if colorIdx == -1 {
		return 0 // Return black for undefined/transparent pixels
	}
	return colorIdx
}

// Pset draws a single pixel at coordinates (x, y) on the current drawing screen
// using the specified PICO-8 color index or the current cursorColor.
//
// The color is specified by its index (0-15) in the standard Pico8Palette.
// If no colorIndex is provided, the current cursorColor is used.
//
// If the coordinates (x, y) are outside the screen bounds, the function does nothing.
// If an invalid colorIndex is provided (e.g., < 0 or > 15), a warning is logged,
// and the function does nothing.
//
// Example:
//
//	Cursor(0, 0, 8) // Set current color to red
//	Pset(10, 20) // Draws a red pixel at (10, 20)
//	Pset(50, 50, 12) // Draws a blue pixel at (50, 50), color overrides cursorColor
func Pset(x, y int, colorIndex ...int) {
	// Check if screen is ready
	if currentScreen == nil || currentScreen.Renderer == nil {
		log.Println("Warning: Pset() called before screen was ready.")
		return
	}

	// Check bounds
	if x < 0 || x >= currentScreen.width || y < 0 || y >= currentScreen.height {
		return
	}

	// Determine color to use
	color := cursorColor // Default to current cursor color
	if len(colorIndex) > 0 {
		color = colorIndex[0]
		if color < 0 || color >= len(Pico8Palette) {
			log.Printf("Warning: Pset() called with invalid color index %d. Palette has %d colors. Ignoring.", color, len(Pico8Palette))
			return
		}
	}

	// Check if this is a transparent color (binary transparency from PaletteTransparency)
	if color < len(PaletteTransparency) && PaletteTransparency[color] {
		// Don't draw transparent pixels (binary transparency)
		return
	}

	// Update our pixel buffer
	currentScreen.pixels[y][x] = color

	// Draw the pixel using a 1x1 rectangle
	rect := image.Rect(x, y, x+1, y+1)
	currentScreen.Renderer.Draw(rect, &image.Uniform{Pico8Palette[color]}, image.Point{}, draw.Over)
}
