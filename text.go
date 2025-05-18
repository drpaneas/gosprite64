package gosprite64

// DefaultFontSize is the default size used for the Print function.
// PICO-8 font is typically 6px high.
const DefaultFontSize = 6.0

// Default values for cursor and display
const (
	defaultCursorX     = 0
	defaultCursorY     = 0
	defaultCursorColor = White
)

// These variables hold the internal state for the cursor used by Print.
var (
	cursorX     = defaultCursorX
	cursorY     = defaultCursorY
	cursorColor = defaultCursorColor
)
