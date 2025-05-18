package gosprite64

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/drpaneas/n64/drivers/display"
	n64draw "github.com/drpaneas/n64/drivers/draw"
	"github.com/drpaneas/n64/machine"
	"github.com/drpaneas/n64/rcp/video"
)

// currentScreen holds the active screen instance
var currentScreen *screen

// VideoPreset represents a predefined video configuration
type VideoPreset int

const (
	// LowRes is the most common setup: 320x240 without interlacing
	LowRes VideoPreset = iota
	// HighRes is 640x480 with interlacing
	HighRes
)

// Screen represents the display surface that can be drawn to.
// This is used internally by the package but needs to be exported
// for the API to work with the n64 drivers.
type Screen struct {
	renderer *n64draw.Rdp
}

// beginDrawing prepares for a new frame by swapping the framebuffer.
// This is an internal function used by the package.
func (s *Screen) beginDrawing() {
	fb := currentScreen.Display.Swap()
	s.renderer.SetFramebuffer(fb)
}

// endDrawing finalizes the frame by flushing the renderer.
// This is an internal function used by the package.
func (s *Screen) endDrawing() {
	s.renderer.Flush()
}

// Clear clears the screen with the specified color.
// This is a package-level function that can be called as gosprite64.Clear(color).
func Clear(c color.Color) {
	if currentScreen != nil {
		currentScreen.Renderer.Draw(currentScreen.Renderer.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)
	}
}

// screen holds all the display-related objects and state.
type screen struct {
	Display  *display.Display
	Renderer *n64draw.Rdp
}

// getPresetConfig returns the configuration for a given preset
func getPresetConfig(preset VideoPreset) (resolution image.Point, colorDepth video.ColorDepth, mode machine.VideoType, interlaced bool) {
	switch preset {
	case LowRes:
		return image.Point{X: 320, Y: 240}, video.ColorDepth(video.BPP16), machine.VideoNTSC, false
	case HighRes:
		return image.Point{X: 640, Y: 480}, video.ColorDepth(video.BPP32), machine.VideoNTSC, true
	default:
		return image.Point{X: 320, Y: 240}, video.ColorDepth(video.BPP16), machine.VideoNTSC, false
	}
}

// Init initializes display with the specified video preset.
// It sets up the framebuffer and renderer for drawing.
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
