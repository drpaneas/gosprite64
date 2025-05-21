package gosprite64

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/clktmr/n64/drivers/display"
	n64draw "github.com/clktmr/n64/drivers/draw"
	"github.com/clktmr/n64/machine"
	"github.com/clktmr/n64/rcp/video"
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
	Bounds   image.Rectangle
	// Cache of image.Uniform for each PICO-8 color
	colorUniforms [16]*image.Uniform
	width         int
	height        int
}

// newScreen creates a new screen instance with the given dimensions
func newScreen(disp *display.Display, renderer *n64draw.Rdp, width, height int) *screen {
	s := &screen{
		Display:  disp,
		Renderer: renderer,
		Bounds:   image.Rect(0, 0, width, height),
		width:    width,
		height:   height,
	}

	// Initialize color uniforms for PICO-8 palette
	for i := range s.colorUniforms {
		s.colorUniforms[i] = &image.Uniform{C: Pico8Palette[i]}
	}

	log.Printf("Screen initialized with %d x %d pixels", width, height)
	return s
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

func getPresetConfig(preset VideoPreset) (image.Point, video.ColorDepth, uint32, bool) {
	switch preset {
	case HighRes:
		return image.Point{X: highResWidth, Y: highResHeight},
			video.BPP32,
			uint32(machine.VideoNTSC), // This is a uint32 value
			true // Interlaced for high resolution
	default: // LowRes or unknown
		return image.Point{X: lowResWidth, Y: lowResHeight},
			video.BPP16,
			uint32(machine.VideoNTSC), // This is a uint32 value
			false // Non-interlaced for low resolution
	}
}

func videoInit(preset VideoPreset) {
	resolution, colorDepth, _, isInterlaced := getPresetConfig(preset)

	// Set the video mode and setup
	video.Setup(isInterlaced)

	// Create display and renderer
	disp := display.NewDisplay(resolution, colorDepth)
	renderer := n64draw.NewRdp()

	// The antialiasing is automatically set to aaResampling in SetFramebuffer
	if fb := disp.Swap(); fb != nil {
		renderer.SetFramebuffer(fb)
	}

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
	// Try to find the color in PICO-8 palette to use cached uniform
	for i, picoColor := range Pico8Palette {
		if picoColor == c {
			s.Renderer.Draw(s.Renderer.Bounds(), s.colorUniforms[i], image.Point{}, draw.Src)
			return
		}
	}
	// Fallback for colors not in PICO-8 palette (shouldn't happen with public API)
	s.Renderer.Draw(s.Renderer.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)
}

// ClearScreen clears the current drawing screen with a specified PICO-8 color index.
// If no colorIndex is provided, it defaults to Pico8Black (0).
func ClearScreen(colorIndex ...int) {
	if currentScreen == nil {
		log.Println("Warning: ClearScreen() called before screen was ready.")
		return
	}

	idx := defaultColorIndex
	if len(colorIndex) > 0 {
		idx = colorIndex[0]
	}

	if idx < 0 || idx >= len(Pico8Palette) {
		log.Printf("Warning: ClearScreen() called with invalid color index %d. Defaulting to %d.", idx, defaultColorIndex)
		idx = defaultColorIndex
	}

	// Clear the screen
	currentScreen.fill(Pico8Palette[idx])
}
