package gosprite64

import (
	"image"
	"image/color"
	"log"

	"github.com/clktmr/n64/drivers/display"
	n64draw "github.com/clktmr/n64/drivers/draw"
	"github.com/clktmr/n64/machine"
	"github.com/clktmr/n64/rcp/texture"
	"github.com/clktmr/n64/rcp/video"
	"github.com/drpaneas/gosprite64/internal/rendergeom"
)

// screen holds all the display-related objects and state.
type screen struct {
	Display     *display.Display
	Framebuffer *texture.Texture
	Bounds      image.Rectangle
	// Cache of image.Uniform for each PICO-8 color
	colorUniforms [16]*image.Uniform
}

// newScreen creates a new screen instance for the fixed framebuffer profile.
func newScreen(disp *display.Display, framebuffer *texture.Texture) *screen {
	bounds := rendergeom.FramebufferBounds()
	s := &screen{
		Display:     disp,
		Framebuffer: framebuffer,
		Bounds:      bounds,
	}

	// Initialize color uniforms for PICO-8 palette
	for i := range s.colorUniforms {
		s.colorUniforms[i] = &image.Uniform{C: Pico8Palette[i]}
	}

	log.Printf("Screen initialized with %d x %d pixels", bounds.Dx(), bounds.Dy())
	return s
}

// currentScreen holds the active screen instance.
var currentScreen *screen

func videoInit() {
	resolution := rendergeom.FramebufferBounds().Size()

	// The square-pixel reset uses one 320x240 progressive framebuffer path.
	video.Setup(false)
	video.SetScale(squarePixelPresentationRect())

	// Create display and seed the first framebuffer.
	disp := display.NewDisplay(resolution, video.BPP16)
	fb := disp.Swap()

	currentScreen = newScreen(disp, fb)
}

func squarePixelPresentationRect() image.Rectangle {
	outputSize := rendergeom.FramebufferBounds().Size().Mul(2)

	switch machine.VideoType {
	case machine.VideoPAL:
		return rendergeom.CenteredRect(image.Rect(128, 45, 128+640, 45+576), outputSize)
	case machine.VideoMPAL, machine.VideoNTSC:
		return rendergeom.CenteredRect(image.Rect(108, 35, 108+640, 35+480), outputSize)
	default:
		return rendergeom.CenteredRect(video.Scale(), outputSize)
	}
}

// beginDrawing prepares for a new frame by swapping the framebuffer.
// This is an internal function used by the package.
func beginDrawing() {
	if currentScreen == nil || currentScreen.Display == nil {
		log.Println("Warning: beginDrawing called before screen was ready.")
		return
	}
	currentScreen.Framebuffer = currentScreen.Display.Swap()
}

// endDrawing finalizes the frame by flushing the renderer.
// This is an internal function used by the package.
func endDrawing() {
	if currentScreen != nil && currentScreen.Framebuffer != nil {
		n64draw.Flush()
	}
}

// fill fills the entire screen with the specified color.
func (s *screen) fill(c color.Color) {
	if s == nil || s.Framebuffer == nil {
		return
	}

	n64draw.Src.Draw(s.Framebuffer, s.Bounds, s.uniform(c), image.Point{})
}

func (s *screen) uniform(c color.Color) image.Image {
	for i, picoColor := range Pico8Palette {
		if picoColor == c {
			return s.colorUniforms[i]
		}
	}
	return &image.Uniform{c}
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
