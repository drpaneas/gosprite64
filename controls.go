package gosprite64

import (
	"sync"

	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/rcp/serial/joybus"
)

// Button constants for N64 controller.
const (
	ButtonA         = joybus.ButtonA
	ButtonB         = joybus.ButtonB
	ButtonZ         = joybus.ButtonZ
	ButtonStart     = joybus.ButtonStart
	ButtonDPadUp    = joybus.ButtonDUp
	ButtonDPadDown  = joybus.ButtonDDown
	ButtonDPadLeft  = joybus.ButtonDLeft
	ButtonDPadRight = joybus.ButtonDRight
	ButtonL         = joybus.ButtonL
	ButtonR         = joybus.ButtonR
	ButtonCUp       = joybus.ButtonCUp
	ButtonCDown     = joybus.ButtonCDown
	ButtonCLeft     = joybus.ButtonCLeft
	ButtonCRight    = joybus.ButtonCRight
)

var (
	// Controller state buffer owned by this package.
	states [4]controller.Controller
	// Current button state (bitmask)
	buttons joybus.ButtonMask
	// Previous button state (for detecting button presses)
	prevButtons joybus.ButtonMask
	// Mutex for thread safety
	controllerMutex sync.Mutex
)

// updateControllerState updates the current and previous button states.
func updateControllerState() {
	controllerMutex.Lock()
	defer controllerMutex.Unlock()

	controller.Poll(&states)

	prevButtons = buttons

	if states[0].Present() {
		buttons = states[0].Down()
	} else {
		buttons = 0
	}
}

// IsButtonDown reports whether the specified button is currently pressed.
func IsButtonDown(button joybus.ButtonMask) bool {
	controllerMutex.Lock()
	defer controllerMutex.Unlock()
	return (buttons & button) != 0
}

// IsButtonJustPressed reports whether the button transitioned from up to down this frame.
func IsButtonJustPressed(button joybus.ButtonMask) bool {
	controllerMutex.Lock()
	defer controllerMutex.Unlock()
	return (buttons&button != 0) && (prevButtons&button == 0)
}

// StickPosition returns the analog stick position in the range [-1.0, 1.0].
func StickPosition(deadzone float64) (float64, float64) {
	controllerMutex.Lock()
	defer controllerMutex.Unlock()

	if !states[0].Present() {
		return 0, 0
	}

	x := float64(states[0].X()) / 128.0
	y := -float64(states[0].Y()) / 128.0

	if x < deadzone && x > -deadzone {
		x = 0
	}
	if y < deadzone && y > -deadzone {
		y = 0
	}

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
