package gosprite64

import (
	_ "embed"
	"embedded/rtos"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/clktmr/n64/rcp/rdp"
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

// Run starts the game loop. It accepts a Gamelooper and an optional Config (from screen.go).
// It will initialize the display, call Init() once, then repeatedly call Update() and Draw().
func Run(g Gamelooper, configs ...Config) {
	chosenPreset := LowRes // Default preset

	if len(configs) > 0 {
		// If a Config object is provided, use its Preset.
		// Note: VideoPreset is an int. If configs[0].Preset is not explicitly set by the user
		// (e.g. pigo8.Run(game, pigo8.Config{})), it will have its zero value.
		// The zero value for VideoPreset (LowRes = iota = 0) aligns with our default.
		// So, we can directly use configs[0].Preset.
		chosenPreset = configs[0].Preset
		log.Printf("Config provided. Using VideoPreset: %v", chosenPreset)
	} else {
		log.Println("No Config provided, defaulting to LowRes video preset.")
	}

	// Initialize display with the determined preset
	videoInit(chosenPreset)

	// Initialize RDP viewport - ADD THESE LINES
	rdp.RDP.SetScissor(image.Rect(0, 0, 640, 480), rdp.InterlaceNone)

	// Call Init before starting the game loop
	g.Init()

	// Initialize audio
	initAudio()

	lastTime := rtos.Nanotime()
	accumulator := time.Duration(0)

	// counter for looping pixel per pixel the bounds of the screen
	// x := 0
	// y := 0

	// Update audio
	UpdateAudio()

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

		// Update audio
		UpdateAudio()

		// Draw game
		beginDrawing()
		// Draw debug border
		currentScreen.DrawBorder()
		g.Draw()

		// In your draw function:
		// currentScreen.DrawDebugOverlay()

		// In your draw function, after DrawDebugOverlay():

		// // Draw test squares at key positions
		// drawTestSquare(0, 0, Pico8Palette[Red], "TL")                                 // Top-left
		// drawTestSquare(VisibleWidth-20, 0, Pico8Palette[Blue], "TR")                  // Top-right
		// drawTestSquare(0, VisibleHeight-20, Pico8Palette[Green], "BL")                // Bottom-left
		// drawTestSquare(VisibleWidth-20, VisibleHeight-20, Pico8Palette[Yellow], "BR") // Bottom-right

		// // Draw a square at (100,100)
		// drawTestSquare(100, 100, Pico8Palette[White], "100,100")

		// // Draw a square at the center
		// centerX := VisibleWidth/2 - 10
		// centerY := VisibleHeight/2 - 10
		// drawTestSquare(centerX, centerY, Pico8Palette[DarkPurple], "CENTER")
		// ClearScreen(Blue)

		// // Upper left
		// drawPixel(0, 0, Pico8Palette[Red])
		// DrawRect(0, 0, 10, 10, Pico8Palette[Red])

		// // Upper right
		// drawPixel(312, 0, Pico8Palette[Red])
		// DrawRect(312-10, 0, 312, 10, Pico8Palette[Red])

		// // // Bottom left
		// // drawPixel(0, 239, Pico8Palette[Red])
		// DrawRect(0, 239-10, 10, 239, Pico8Palette[Red])

		// // // Bottom right
		// // drawPixel(312, 239, Pico8Palette[Red])
		// DrawRect(312-10, 239-10, 312, 239, Pico8Palette[Red])

		// PrintBitmap(fmt.Sprintf("FPS: %d", targetFPS), currentScreen.width-80, 20, 9)
		endDrawing()

		// Sleep to maintain target frame rate
		sleepDuration := frameDuration - (rtos.Nanotime() - currentTime)
		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}
	}
}

const (
	// Calculate borders based on your screenshot
	BorderLeft = 2
	BorderTop  = 0
)

// ToScreen converts game coordinates to actual screen coordinates
func (s *screen) ToScreen(x, y int) (int, int) {
	return 2*BorderLeft + x, BorderTop + y
}

// func drawPixel(x, y int, c color.Color) {
// 	// x, y out of bounds check
// 	// top left: [0,0]
// 	// top right: [319,0]
// 	// bottom left: [0,239]
// 	// bottom right: [319,239]
// 	if x < 0 || x > 319 || y < 0 || y > 239 {
// 		log.Printf("drawPixel: x=%d, y=%d out of bounds", x, y)
// 		return
// 	}

// 	screenX, screenY := currentScreen.ToScreen(x, y)
// 	Rectfill(screenX, screenY, screenX, screenY, c)
// }

// DrawRectFill draws a rectangle on screen
// x1, y1: Top-left corner (0-319, 0-239)
// x2, y2: Bottom-right corner (inclusive)
// color: The color to fill the rectangle with
func DrawRectFill(x1, y1, x2, y2 int, c color.Color) {
	if x1 < 0 || x1 > 319 || y1 < 0 || y1 > 239 || x2 < 0 || x2 > 319 || y2 < 0 || y2 > 239 {
		log.Printf("drawRect: x1=%d, y1=%d, x2=%d, y2=%d out of bounds", x1, y1, x2, y2)
		return
	}

	tmpX1, tmpY1 := currentScreen.ToScreen(x1, y1)
	tmpX2, tmpY2 := currentScreen.ToScreen(x2, y2)
	Rectfill(tmpX1, tmpY1, tmpX2, tmpY2, c)
}
