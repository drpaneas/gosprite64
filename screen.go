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

type screen struct {
	Display      *display.Display
	Framebuffer  *texture.Texture
	Bounds       image.Rectangle
	uniformCache map[color.Color]*image.Uniform
}

func newScreen(disp *display.Display, framebuffer *texture.Texture) *screen {
	bounds := rendergeom.FramebufferBounds()
	defaults := []color.Color{
		Black, DarkBlue, DarkPurple, DarkGreen, Brown, DarkGray,
		LightGray, White, Red, Orange, Yellow, Green, Blue, Indigo, Pink, Peach,
	}
	cache := make(map[color.Color]*image.Uniform, len(defaults))
	for _, c := range defaults {
		cache[c] = &image.Uniform{C: c}
	}
	s := &screen{
		Display:      disp,
		Framebuffer:  framebuffer,
		Bounds:       bounds,
		uniformCache: cache,
	}
	log.Printf("Screen initialized with %d x %d pixels", bounds.Dx(), bounds.Dy())
	return s
}

var currentScreen *screen

func videoInit() {
	resolution := rendergeom.FramebufferBounds().Size()
	video.Setup(false)
	video.SetScale(squarePixelPresentationRect())
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

func beginDrawing() {
	if currentScreen == nil || currentScreen.Display == nil {
		log.Println("Warning: beginDrawing called before screen was ready.")
		return
	}
	currentScreen.Framebuffer = currentScreen.Display.Swap()
}

func endDrawing() {
	if currentScreen != nil && currentScreen.Framebuffer != nil {
		n64draw.Flush()
	}
}

func (s *screen) fill(c color.Color) {
	if s == nil || s.Framebuffer == nil {
		return
	}
	n64draw.Src.Draw(s.Framebuffer, s.Bounds, s.uniform(c), image.Point{})
}

func (s *screen) uniform(c color.Color) image.Image {
	if c == nil {
		c = Black
	}
	switch c.(type) {
	case color.RGBA, color.NRGBA, color.Gray, color.Alpha, color.Gray16, color.NRGBA64, color.RGBA64:
		if u, ok := s.uniformCache[c]; ok {
			return u
		}
	}
	return &image.Uniform{C: c}
}

// ClearScreen fills the screen with Black.
func ClearScreen() {
	if currentScreen == nil {
		log.Println("Warning: ClearScreen() called before screen was ready.")
		return
	}
	currentScreen.fill(Black)
}

// ClearScreenWith fills the screen with the given color.
func ClearScreenWith(c color.Color) {
	if currentScreen == nil {
		log.Println("Warning: ClearScreenWith() called before screen was ready.")
		return
	}
	currentScreen.fill(c)
}
