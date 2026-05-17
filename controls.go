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

// MaxControllers is the number of controller ports on the N64.
const MaxControllers = 4

var (
	states      [MaxControllers]controller.Controller
	buttons     [MaxControllers]joybus.ButtonMask
	prevButtons [MaxControllers]joybus.ButtonMask

	controllerMutex sync.Mutex
)

func updateControllerState() {
	controllerMutex.Lock()
	defer controllerMutex.Unlock()

	controller.Poll(&states)

	for i := 0; i < MaxControllers; i++ {
		prevButtons[i] = buttons[i]
		if states[i].Present() {
			buttons[i] = states[i].Down()
		} else {
			buttons[i] = 0
		}
	}
}

// --- Per-port multiplayer API ---

// PlayerButtonDown reports whether the specified button is currently pressed
// on the controller at the given port (0-3).
func PlayerButtonDown(port int, button joybus.ButtonMask) bool {
	if port < 0 || port >= MaxControllers {
		return false
	}
	controllerMutex.Lock()
	defer controllerMutex.Unlock()
	return (buttons[port] & button) != 0
}

// PlayerButtonJustPressed reports whether the button transitioned from up to
// down this frame on the controller at the given port (0-3).
func PlayerButtonJustPressed(port int, button joybus.ButtonMask) bool {
	if port < 0 || port >= MaxControllers {
		return false
	}
	controllerMutex.Lock()
	defer controllerMutex.Unlock()
	return (buttons[port]&button != 0) && (prevButtons[port]&button == 0)
}

// PlayerStickPosition returns the analog stick position in the range [-1.0, 1.0]
// for the controller at the given port (0-3).
func PlayerStickPosition(port int, deadzone float64) (float64, float64) {
	if port < 0 || port >= MaxControllers {
		return 0, 0
	}
	controllerMutex.Lock()
	defer controllerMutex.Unlock()

	if !states[port].Present() {
		return 0, 0
	}

	x := float64(states[port].X()) / 128.0
	y := -float64(states[port].Y()) / 128.0

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

// IsControllerConnected reports whether the controller at the given port (0-3)
// is currently connected.
func IsControllerConnected(port int) bool {
	if port < 0 || port >= MaxControllers {
		return false
	}
	controllerMutex.Lock()
	defer controllerMutex.Unlock()
	return states[port].Present()
}

// ConnectedControllers returns the number of controllers currently connected.
func ConnectedControllers() int {
	controllerMutex.Lock()
	defer controllerMutex.Unlock()
	count := 0
	for i := 0; i < MaxControllers; i++ {
		if states[i].Present() {
			count++
		}
	}
	return count
}

// --- Port 0 convenience wrappers (backward-compatible API) ---

// IsButtonDown reports whether the specified button is currently pressed on port 0.
func IsButtonDown(button joybus.ButtonMask) bool {
	return PlayerButtonDown(0, button)
}

// IsButtonJustPressed reports whether the button transitioned from up to down
// this frame on port 0.
func IsButtonJustPressed(button joybus.ButtonMask) bool {
	return PlayerButtonJustPressed(0, button)
}

// StickPosition returns the analog stick position in the range [-1.0, 1.0]
// for port 0.
func StickPosition(deadzone float64) (float64, float64) {
	return PlayerStickPosition(0, deadzone)
}

// SetRumble enables or disables the rumble pak on the given port (0-3).
func SetRumble(port int, enabled bool) {
	if port < 0 || port >= MaxControllers {
		return
	}
	controllerMutex.Lock()
	defer controllerMutex.Unlock()
	if !states[port].Present() {
		return
	}
	rumbleWrite(port, enabled)
}
