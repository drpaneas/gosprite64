package gosprite64

import (
	_ "embed"
	"embedded/rtos"
	"image"
	"time"

	"github.com/clktmr/n64/rcp/rdp"
	"github.com/drpaneas/gosprite64/internal/rendergeom"
)

// Game represents a game instance that can be initialized, updated, and drawn.
type Game interface {
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
func Run(g Game) {
	setupConsole()
	rt := newRuntimeState()
	rt.initVideo()
	activateRuntime(rt)

	framebufferBounds := rendergeom.FramebufferBounds()
	rdp.RDP.SetScissor(
		image.Rect(0, 0, framebufferBounds.Dx()*2, framebufferBounds.Dy()*2),
		rdp.InterlaceNone,
	)

	// Call Init before starting the game loop
	g.Init()

	// Audio init runs after g.Init() so that pre-init PlayEffect/PlayTrack
	// calls from g.Init() are silent no-ops, matching the spec (section 3.3).
	rt.initAudio()

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
