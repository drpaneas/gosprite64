package gosprite64

import (
	"sync"

	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/rcp/serial/joybus"
)

// Button constants for N64 controller
const (
	BtnA      = joybus.ButtonA
	BtnB      = joybus.ButtonB
	BtnZ      = joybus.ButtonZ
	BtnStart  = joybus.ButtonStart
	BtnDUp    = joybus.ButtonDUp
	BtnDDown  = joybus.ButtonDDown
	BtnDLeft  = joybus.ButtonDLeft
	BtnDRight = joybus.ButtonDRight
	BtnL      = joybus.ButtonL
	BtnR      = joybus.ButtonR
	BtnCUp    = joybus.ButtonCUp
	BtnCDown  = joybus.ButtonCDown
	BtnCLeft  = joybus.ButtonCLeft
	BtnCRight = joybus.ButtonCRight
)

// Button aliases for PICO-8 compatibility
const (
	BtnRight = BtnDRight
	BtnLeft  = BtnDLeft
	BtnUp    = BtnDUp
	BtnDown  = BtnDDown
	BtnX     = BtnCLeft  // PICO-8 X button
	BtnC     = BtnCRight // PICO-8 O button
)

var (
	// Current button state (bitmask)
	buttons joybus.ButtonMask
	// Previous button state (for detecting button presses)
	prevButtons joybus.ButtonMask
	// Mutex for thread safety
	controllerMutex sync.Mutex
)

// updateControllerState updates the current and previous button states
func updateControllerState() {
	controllerMutex.Lock()
	defer controllerMutex.Unlock()

	// Poll all controllers
	controller.Poll()

	// Save previous button state
	prevButtons = buttons

	// Get current button state from first connected controller
	if controller.States[0].Present() {
		buttons = controller.States[0].Down()
	} else {
		buttons = 0
	}
}

// Btn returns true if the specified button is currently pressed.
// Use the Btn* constants to specify which button to check.
func Btn(button joybus.ButtonMask) bool {
	controllerMutex.Lock()
	defer controllerMutex.Unlock()
	return (buttons & button) != 0
}

// Btnp returns true if the specified button was just pressed (i.e., it was not pressed in the previous frame).
// This is useful for detecting button presses rather than holds.
func Btnp(button joybus.ButtonMask) bool {
	controllerMutex.Lock()
	defer controllerMutex.Unlock()
	return (buttons&button != 0) && (prevButtons&button == 0)
}

// GetStick returns the analog stick position as x, y values in the range [-1.0, 1.0].
// The deadzone parameter specifies the minimum value to be considered as input (typically 0.1-0.3).
func GetStick(deadzone float64) (float64, float64) {
	controllerMutex.Lock()
	defer controllerMutex.Unlock()

	if !controller.States[0].Present() {
		return 0, 0
	}

	x := float64(controller.States[0].X()) / 128.0
	y := -float64(controller.States[0].Y()) / 128.0 // Invert Y axis to match screen coordinates

	// Apply deadzone
	if x < deadzone && x > -deadzone {
		x = 0
	}
	if y < deadzone && y > -deadzone {
		y = 0
	}

	// Clamp values to [-1, 1] range
	if x > 1.0 {
		x = 1.0
	} else if x < -1.0 {
		x = -1.0
	}
	if y > 1.0 {
		y = 1.0
	} else if y < -1.0 {
		y = -1.0
	}

	return x, y
}
