package gosprite64

import "image/color"

// MenuItem is a single entry in a menu.
type MenuItem struct {
	Label     string
	Disabled  bool
	OnConfirm func()
}

// Menu manages a D-pad-navigated list of options with cursor tracking.
// Use it inside a GameState's Update/Draw to build menu screens.
type Menu struct {
	items  []MenuItem
	cursor int
	Wrap   bool

	X, Y       int
	LineHeight int
	Color      color.Color
	CursorChar string
}

// NewMenu creates a menu with the given items.
func NewMenu(items []MenuItem) *Menu {
	return &Menu{
		items:      items,
		LineHeight: 12,
		Color:      White,
		CursorChar: "> ",
	}
}

// Count returns the number of menu items.
func (m *Menu) Count() int {
	if m == nil {
		return 0
	}
	return len(m.items)
}

// Cursor returns the current cursor position.
func (m *Menu) Cursor() int {
	if m == nil {
		return 0
	}
	return m.cursor
}

// SetCursor sets the cursor position, clamped to valid range.
func (m *Menu) SetCursor(index int) {
	if m == nil || len(m.items) == 0 {
		return
	}
	if index < 0 {
		index = 0
	}
	if index >= len(m.items) {
		index = len(m.items) - 1
	}
	m.cursor = index
}

// Selected returns the currently highlighted item.
func (m *Menu) Selected() MenuItem {
	if m == nil || len(m.items) == 0 {
		return MenuItem{}
	}
	return m.items[m.cursor]
}

// MoveDown moves the cursor down, skipping disabled items.
// If all items are disabled, the cursor stays where it is.
func (m *Menu) MoveDown() {
	if m == nil || len(m.items) == 0 {
		return
	}
	start := m.cursor
	for {
		m.cursor++
		if m.cursor >= len(m.items) {
			if m.Wrap {
				m.cursor = 0
			} else {
				m.cursor = start
				return
			}
		}
		if !m.items[m.cursor].Disabled {
			return
		}
		if m.cursor == start {
			return
		}
	}
}

// MoveUp moves the cursor up, skipping disabled items.
// If all items are disabled, the cursor stays where it is.
func (m *Menu) MoveUp() {
	if m == nil || len(m.items) == 0 {
		return
	}
	start := m.cursor
	for {
		m.cursor--
		if m.cursor < 0 {
			if m.Wrap {
				m.cursor = len(m.items) - 1
			} else {
				m.cursor = start
				return
			}
		}
		if !m.items[m.cursor].Disabled {
			return
		}
		if m.cursor == start {
			return
		}
	}
}

// Confirm triggers the OnConfirm callback of the selected item.
func (m *Menu) Confirm() {
	if m == nil || len(m.items) == 0 {
		return
	}
	item := m.items[m.cursor]
	if item.Disabled || item.OnConfirm == nil {
		return
	}
	item.OnConfirm()
}

// HandleInput reads the controller and updates the cursor.
// Call this in your GameState's Update().
// Returns true if A was pressed (confirmation).
func (m *Menu) HandleInput() bool {
	if m == nil || len(m.items) == 0 {
		return false
	}
	if IsButtonJustPressed(ButtonDPadDown) {
		m.MoveDown()
	}
	if IsButtonJustPressed(ButtonDPadUp) {
		m.MoveUp()
	}
	if IsButtonJustPressed(ButtonA) {
		m.Confirm()
		return true
	}
	return false
}

// Draw renders the menu using DrawText. Call this in your GameState's Draw().
func (m *Menu) Draw() {
	if m == nil || len(m.items) == 0 {
		return
	}
	c := m.Color
	if c == nil {
		c = White
	}
	for i, item := range m.items {
		y := m.Y + i*m.LineHeight
		prefix := "  "
		if i == m.cursor {
			prefix = m.CursorChar
		}
		itemColor := c
		if item.Disabled {
			itemColor = DarkGray
		}
		DrawText(prefix+item.Label, m.X, y, itemColor)
	}
}
