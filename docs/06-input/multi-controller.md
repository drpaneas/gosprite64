# Multi-Controller Support

Read input from up to four controllers for multiplayer games. The N64 has four controller ports numbered 0 through 3.

```go
import gs "github.com/drpaneas/gosprite64"
```

## Controller ports

Each port is identified by an integer from 0 to 3. Port 0 is the leftmost port on the console. The single-player convenience functions (`IsButtonDown`, `IsButtonJustPressed`, `StickPosition`) always read from port 0.

The `MaxControllers` constant is `4`.

## Detecting connected controllers

### ConnectedControllers

Returns the total number of controllers currently plugged in:

```go
count := gs.ConnectedControllers()
if count < 2 {
    showMessage("Please connect a second controller")
}
```

### IsControllerConnected

Checks whether a specific port has a controller plugged in:

```go
for port := 0; port < gs.MaxControllers; port++ {
    if gs.IsControllerConnected(port) {
        fmt.Printf("Port %d: connected\n", port)
    }
}
```

If a controller is disconnected mid-game, `IsControllerConnected` will return `false` on the next frame and all button/stick queries for that port will return zero values.

## Reading buttons per port

### PlayerButtonDown

Returns `true` every frame the button is held on the given port:

```go
func PlayerButtonDown(port int, button ButtonMask) bool
```

### PlayerButtonJustPressed

Returns `true` only on the frame the button transitions from released to pressed:

```go
func PlayerButtonJustPressed(port int, button ButtonMask) bool
```

### Example: two-player movement

```go
func (g *Game) Update() {
    // Player 1 (port 0)
    if gs.PlayerButtonDown(0, gs.ButtonDPadRight) {
        g.players[0].X += speed
    }
    if gs.PlayerButtonDown(0, gs.ButtonDPadLeft) {
        g.players[0].X -= speed
    }
    if gs.PlayerButtonJustPressed(0, gs.ButtonA) {
        g.players[0].Jump()
    }

    // Player 2 (port 1)
    if gs.PlayerButtonDown(1, gs.ButtonDPadRight) {
        g.players[1].X += speed
    }
    if gs.PlayerButtonDown(1, gs.ButtonDPadLeft) {
        g.players[1].X -= speed
    }
    if gs.PlayerButtonJustPressed(1, gs.ButtonA) {
        g.players[1].Jump()
    }
}
```

## Reading the analog stick per port

### PlayerStickPosition

Returns the stick X and Y in [-1.0, 1.0] for the given port, with deadzone filtering:

```go
func PlayerStickPosition(port int, deadzone float64) (float64, float64)
```

```go
x, y := gs.PlayerStickPosition(1, 0.15)
g.players[1].X += x * moveSpeed
g.players[1].Y += y * moveSpeed
```

If the port is out of range (not 0-3) or no controller is connected, all functions return zero values (`false` for buttons, `0, 0` for the stick). You do not need to guard every call with `IsControllerConnected`, but checking connection status is useful for UI prompts.

## Full multiplayer loop

```go
func (g *Game) Update() {
    for port := 0; port < gs.MaxControllers; port++ {
        if !gs.IsControllerConnected(port) {
            continue
        }
        p := &g.players[port]

        // Analog stick movement
        x, y := gs.PlayerStickPosition(port, 0.15)
        p.X += x * moveSpeed
        p.Y += y * moveSpeed

        // Action buttons
        if gs.PlayerButtonJustPressed(port, gs.ButtonA) {
            p.Jump()
        }
        if gs.PlayerButtonJustPressed(port, gs.ButtonB) {
            p.Attack()
        }
    }
}
```

## Relationship to single-player API

The single-player functions are thin wrappers around the per-port API:

| Single-player | Per-port equivalent |
|---|---|
| `IsButtonDown(btn)` | `PlayerButtonDown(0, btn)` |
| `IsButtonJustPressed(btn)` | `PlayerButtonJustPressed(0, btn)` |
| `StickPosition(dz)` | `PlayerStickPosition(0, dz)` |

## Reference Example

See `examples/multi_input_demo` in the GoSprite64 repository for a 4-player controller demonstration with per-port colored squares.
