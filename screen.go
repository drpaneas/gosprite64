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
}

// currentScreen holds the active screen instance.
var currentScreen *screen

// Screen represents the display surface that can be drawn to.
// This is a thin wrapper around the N64's display and renderer.
type Screen struct {
	renderer *n64draw.Rdp
}

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

// Init initializes the display with the specified video preset.
// It sets up the framebuffer and renderer for drawing.
// Returns a new Screen instance that can be used for drawing operations.
func Init(preset VideoPreset) *Screen {
	resolution, colorDepth, mode, interlaced := getPresetConfig(preset)

	// Set the video mode and setup
	machine.Video = mode
	video.Setup(interlaced)

	// Create display and renderer
	disp := display.NewDisplay(resolution, colorDepth)
	renderer := n64draw.NewRdp()
	renderer.SetFramebuffer(disp.Swap())

	currentScreen = &screen{
		Display:  disp,
		Renderer: renderer,
	}

	return &Screen{
		renderer: renderer,
	}
}

// beginDrawing prepares for a new frame by swapping the framebuffer.
// This is an internal function used by the package.
func (s *Screen) beginDrawing() {
	if currentScreen == nil || currentScreen.Display == nil {
		log.Println("Warning: beginDrawing called before screen was ready.")
		return
	}
	fb := currentScreen.Display.Swap()
	s.renderer.SetFramebuffer(fb)
}

// endDrawing finalizes the frame by flushing the renderer.
// This is an internal function used by the package.
func (s *Screen) endDrawing() {
	if s.renderer != nil {
		s.renderer.Flush()
	}
}

// Fill fills the entire screen with the specified color.
func (s *screen) Fill(c color.Color) {
	s.Renderer.Draw(s.Renderer.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)
}

// Default color index for Cls
const defaultColorIndex = Black

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

	currentScreen.Fill(Pico8Palette[idx])

	// Reset the global print cursor position
	cursorX = defaultCursorX
	cursorY = defaultCursorY
}
