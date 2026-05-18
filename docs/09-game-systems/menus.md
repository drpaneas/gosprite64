# Menus

The `Menu` type handles the boilerplate of D-pad-navigated option lists: cursor tracking, wrapping, disabled items, and confirmation callbacks. Use it inside any `GameState` for title screens, option screens, pause menus, and game-over screens.

## Quick start

```go
type TitleState struct {
    sm   *gosprite64.StateMachine
    menu *gosprite64.Menu
}

func (s *TitleState) Enter() {
    s.menu = gosprite64.NewMenu([]gosprite64.MenuItem{
        {Label: "Start Game", OnConfirm: func() {
            s.sm.Switch(&GameplayState{sm: s.sm})
        }},
        {Label: "Options", OnConfirm: func() {
            s.sm.Push(&OptionsState{sm: s.sm})
        }},
        {Label: "Quit", Disabled: true},
    })
    s.menu.X = 100
    s.menu.Y = 100
    s.menu.Wrap = true
}

func (s *TitleState) Update() {
    s.menu.HandleInput()
}

func (s *TitleState) Draw() {
    gosprite64.ClearScreen()
    gosprite64.DrawText("MY GAME", 112, 40, gosprite64.White)
    s.menu.Draw()
}
```

`HandleInput` reads D-pad Up/Down for navigation and A for confirmation. `Draw` renders the menu with a cursor indicator. That's the entire integration.

## MenuItem

Each item has a label, an optional callback, and a disabled flag:

```go
type MenuItem struct {
    Label     string      // displayed text
    Disabled  bool        // grayed out, skipped by cursor, can't confirm
    OnConfirm func()      // called when A is pressed on this item
}
```

`OnConfirm` can do anything: switch states, change settings, start gameplay. If `OnConfirm` is nil or the item is disabled, pressing A does nothing.

## Navigation

The cursor moves with D-pad Up and Down. Disabled items are automatically skipped:

```go
menu := gosprite64.NewMenu([]gosprite64.MenuItem{
    {Label: "Play"},
    {Label: "Locked", Disabled: true},
    {Label: "Settings"},
})
// Pressing Down from "Play" skips "Locked" and lands on "Settings"
```

If all items are disabled, the cursor stays where it is.

### Wrapping

By default, the cursor stops at the first and last item. Set `Wrap = true` to wrap around:

```go
menu.Wrap = true
// Pressing Down on the last item goes to the first
// Pressing Up on the first item goes to the last
```

## Manual control

If `HandleInput` doesn't fit your needs (e.g. you want analog stick support or different buttons), use the navigation methods directly:

```go
func (s *MyState) Update() {
    if gosprite64.IsButtonJustPressed(gosprite64.ButtonDPadDown) {
        s.menu.MoveDown()
    }
    if gosprite64.IsButtonJustPressed(gosprite64.ButtonDPadUp) {
        s.menu.MoveUp()
    }
    if gosprite64.IsButtonJustPressed(gosprite64.ButtonStart) {
        s.menu.Confirm()
    }
}
```

Other useful methods:

```go
menu.Cursor()       // current index (0-based)
menu.SetCursor(2)   // jump to index (clamped to valid range)
menu.Selected()     // get the highlighted MenuItem
menu.Count()        // number of items
```

## Customizing appearance

The built-in `Draw` uses the 8x8 `DrawText` font. Customize position, spacing, colors, and cursor character:

```go
menu.X = 80           // left edge
menu.Y = 60           // top edge
menu.LineHeight = 14   // pixels between items (default 12)
menu.Color = gosprite64.Yellow          // text color
menu.CursorChar = "-> "                // cursor prefix (default "> ")
```

Disabled items are drawn in `DarkGray` regardless of `Color`.

### Custom rendering

For full control (custom fonts, icons, backgrounds), skip `Draw` and render manually using `Cursor()`:

```go
func (s *MyState) Draw() {
    for i, item := range items {
        y := 60 + i*16
        if i == s.menu.Cursor() {
            gosprite64.FillRect(78, y-2, 220, y+10, gosprite64.DarkBlue)
        }
        c := gosprite64.White
        if item.Disabled {
            c = gosprite64.DarkGray
        }
        myFont.DrawTextEx(item.Label, 84, y, gosprite64.AlignLeft)
    }
}
```

## Multiple menus

A state can have multiple menus. Track which is active:

```go
type OptionsState struct {
    difficultyMenu *gosprite64.Menu
    speedMenu      *gosprite64.Menu
    activeMenu     int  // 0 = difficulty, 1 = speed
}

func (s *OptionsState) Update() {
    switch s.activeMenu {
    case 0:
        s.difficultyMenu.HandleInput()
    case 1:
        s.speedMenu.HandleInput()
    }
    if gosprite64.IsButtonJustPressed(gosprite64.ButtonR) {
        s.activeMenu = (s.activeMenu + 1) % 2
    }
}
```

## Try It

<iframe src="../emulator/play.html?rom=menu_demo.z64" width="640" height="480" frameborder="0" allow="autoplay" style="display:block;margin:0 auto;max-width:100%;"></iframe>

> **Controls:** Arrow keys = D-Pad, X = A button, C = B button, Enter = Start, Z = Z trigger

## Complete example

See `examples/menu_demo` for a focused menu demonstration, or `examples/splitscreen_demo` for menus used alongside other systems.
