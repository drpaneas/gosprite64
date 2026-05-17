# D-Pad and Buttons

Read digital button presses from the player-one controller. GoSprite64 polls controller state once per frame automatically, so these functions always reflect the current frame's input.

```go
import gs "github.com/drpaneas/gosprite64"
```

## Button constants

The N64 controller has 14 digital buttons. Each constant is a `ButtonMask` bitmask, so you can combine them with bitwise OR when needed.

| Constant | Button |
|---|---|
| `gs.ButtonA` | A (green, right thumb) |
| `gs.ButtonB` | B (blue, right thumb) |
| `gs.ButtonZ` | Z (trigger, underside) |
| `gs.ButtonStart` | Start |
| `gs.ButtonL` | Left shoulder |
| `gs.ButtonR` | Right shoulder |
| `gs.ButtonDPadUp` | D-Pad up |
| `gs.ButtonDPadDown` | D-Pad down |
| `gs.ButtonDPadLeft` | D-Pad left |
| `gs.ButtonDPadRight` | D-Pad right |
| `gs.ButtonCUp` | C-up (yellow) |
| `gs.ButtonCDown` | C-down (yellow) |
| `gs.ButtonCLeft` | C-left (yellow) |
| `gs.ButtonCRight` | C-right (yellow) |

## Checking held buttons

`IsButtonDown` returns `true` for every frame the button is physically held down:

```go
func (g *Game) Update() {
    if gs.IsButtonDown(gs.ButtonDPadRight) {
        player.X += speed
    }
    if gs.IsButtonDown(gs.ButtonDPadLeft) {
        player.X -= speed
    }
    if gs.IsButtonDown(gs.ButtonDPadUp) {
        player.Y -= speed
    }
    if gs.IsButtonDown(gs.ButtonDPadDown) {
        player.Y += speed
    }
}
```

Use `IsButtonDown` for continuous actions like movement - the callback fires every frame while the button is pressed.

## Detecting a fresh press

`IsButtonJustPressed` returns `true` only on the first frame the button transitions from released to pressed. It will not fire again until the player releases and re-presses the button.

```go
func (g *Game) Update() {
    if gs.IsButtonJustPressed(gs.ButtonA) {
        player.Jump()
    }
    if gs.IsButtonJustPressed(gs.ButtonB) {
        player.Attack()
    }
    if gs.IsButtonJustPressed(gs.ButtonStart) {
        g.TogglePause()
    }
}
```

Use `IsButtonJustPressed` for one-shot actions like jumping, attacking, or opening a menu - actions that should not repeat while the button is held.

## Combining both

A common pattern uses `IsButtonJustPressed` for the initial action and `IsButtonDown` for sustained behavior:

```go
func (g *Game) Update() {
    // Charge attack: start on press, charge while held
    if gs.IsButtonJustPressed(gs.ButtonB) {
        player.StartCharge()
    }
    if gs.IsButtonDown(gs.ButtonB) {
        player.AddCharge(1)
    }

    // D-Pad movement (continuous)
    if gs.IsButtonDown(gs.ButtonDPadRight) {
        player.X += 2
    }

    // Menu navigation (one tap per press)
    if gs.IsButtonJustPressed(gs.ButtonDPadDown) {
        menu.MoveDown()
    }
}
```

## Port-zero convenience

`IsButtonDown` and `IsButtonJustPressed` read from controller port 0 (the first controller). For multiplayer input, see [Multi-Controller Support](multi-controller.md).
