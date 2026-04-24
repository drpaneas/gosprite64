package gosprite64

import (
	_ "embed"
	"embedded/rtos"
	"image"
	"image/color"
	"time"

	"github.com/clktmr/n64/rcp/rdp"
	"github.com/drpaneas/gosprite64/internal/rendergeom"
)

// Gamelooper represents a game instance that can be updated and drawn.
type Gamelooper interface {
	// Init is called once at the start of the game before the first Update.
	// Use this to initialize game state, load resources, etc.
	Init()

	// Update is called every frame to update game logic.
	Update()

	// Draw is called every frame to render the game.
	// The screen is already initialized and ready for drawing.
	Draw()
}

var TargetFPS = 60

var frameDuration = time.Second / time.Duration(TargetFPS)

// Run starts the game loop using the fixed square-pixel framebuffer path.
func Run(g Gamelooper) {
	videoInit()

	framebufferBounds := rendergeom.FramebufferBounds()
	rdp.RDP.SetScissor(
		image.Rect(0, 0, framebufferBounds.Dx()*2, framebufferBounds.Dy()*2),
		rdp.InterlaceNone,
	)

	// Call Init before starting the game loop
	g.Init()

	// Initialize audio registration as a defensive fallback in case a consumer
	// wires SetAudioFS later than the generated init path.
	initAudio()

	lastTime := rtos.Nanotime()
	accumulator := time.Duration(0)

	// Main game loop
	for {
		currentTime := rtos.Nanotime()
		elapsed := time.Duration(currentTime - lastTime)
		lastTime = currentTime
		accumulator += elapsed

		for accumulator >= frameDuration {
			updateControllerState()
			g.Update()
			accumulator -= frameDuration
		}

		// Perform lightweight audio housekeeping for completed one-shot cues.
		UpdateAudio()

		// Draw game
		beginDrawing()
		g.Draw()
		endDrawing()

		// Sleep to maintain target frame rate
		sleepDuration := frameDuration - (rtos.Nanotime() - currentTime)
		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}
	}
}

// DrawRectFill draws a rectangle on screen
// x1, y1: Top-left logical corner (0-287, 0-215)
// x2, y2: Bottom-right corner (inclusive)
// color: The color to fill the rectangle with
func DrawRectFill(x1, y1, x2, y2 int, c color.Color) {
	Rectfill(x1, y1, x2, y2, c)
}
