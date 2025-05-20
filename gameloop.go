package gosprite64

import (
	_ "embed"
	"embedded/rtos"
	"fmt"
	"time"
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

const targetFPS = 60
const frameDuration = time.Second / targetFPS

// Run starts the game loop with default video settings (NTSC 320x240, no interlacing).
// It will initialize the display, call Init() once, then repeatedly call Update() and Draw().
func Run(g Gamelooper) {
	// Initialize display with default settings
	videoInit(LowRes)

	// Call Init before starting the game loop
	g.Init()

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

		// Update controller state
		updateControllerState()

		// Update game logic
		g.Update()

		// Draw game
		beginDrawing()
		g.Draw()
		PrintBitmap(fmt.Sprintf("FPS: %d", targetFPS), currentScreen.width-60, 0, 7)
		endDrawing()

		// Sleep to maintain target frame rate
		sleepDuration := frameDuration - (rtos.Nanotime() - currentTime)
		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}
	}
}
