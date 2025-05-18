package gosprite64

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

// Run starts the game loop with default video settings (NTSC 320x240, no interlacing).
// It will initialize the display, call Init() once, then repeatedly call Update() and Draw().
func Run(g Gamelooper) {
	// Initialize display with default settings
	qualityPreset := LowRes // TODO: make this configurable
	screen := Init(qualityPreset)

	// Call Init before starting the game loop
	g.Init()

	// Main game loop
	for {
		// Update game logic
		g.Update()

		// Draw game
		screen.beginDrawing()
		g.Draw()
		screen.endDrawing()
	}
}
