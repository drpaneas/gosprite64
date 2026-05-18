# Analog Stick

Read the N64 analog stick for smooth directional movement and aiming.

```go
import gs "github.com/drpaneas/gosprite64"
```

## StickPosition

`StickPosition` returns the analog stick's X and Y position as two `float64` values, each in the range -1.0 to 1.0.

```go
func StickPosition(deadzone float64) (float64, float64)
```

- **X axis:** -1.0 is full left, 1.0 is full right.
- **Y axis:** -1.0 is full up, 1.0 is full down. The Y axis is flipped compared to raw hardware values so that "down on the screen" is positive, matching screen coordinates.
- Values are clamped to [-1.0, 1.0] even if the hardware reports values slightly outside that range.

### The deadzone parameter

N64 analog sticks rarely rest at a perfect zero position. A worn stick might report small non-zero values even when the player is not touching it. The `deadzone` parameter defines a threshold: any axis value between `-deadzone` and `+deadzone` is snapped to zero.

A deadzone of `0.15` is a good starting point for most games:

```go
x, y := gs.StickPosition(0.15)
```

- `0.0` - No deadzone. The raw (clamped) value is returned. Fine for menus or testing, but can cause drift on worn controllers.
- `0.10 - 0.20` - Typical range for action games. Filters out stick drift while preserving sensitivity.
- `0.25+` - Large deadzone. The player must push the stick further before movement registers. Useful for games where accidental input is costly.

## Character movement example

```go
const moveSpeed = 3.0

func (g *Game) Update() {
    x, y := gs.StickPosition(0.15)

    player.X += x * moveSpeed
    player.Y += y * moveSpeed
}
```

Because `StickPosition` returns values between -1.0 and 1.0, multiplying by a speed constant gives smooth, proportional movement. Gentle stick tilts produce slow movement; full tilts produce maximum speed.

## Eight-directional movement with normalization

When the stick is pushed diagonally, both X and Y are non-zero. Naively adding both axes makes diagonal movement ~41% faster than cardinal movement. Normalize the direction vector to fix this:

```go
import "math"

func (g *Game) Update() {
    x, y := gs.StickPosition(0.15)

    length := math.Sqrt(x*x + y*y)
    if length > 0 {
        if length > 1.0 {
            length = 1.0
        }
        player.X += (x / length) * moveSpeed * length
        player.Y += (y / length) * moveSpeed * length
    }
}
```

## Aiming and cursor control

The analog stick works well for aiming a cursor or rotating a character:

```go
func (g *Game) Update() {
    x, y := gs.StickPosition(0.2)

    if x != 0 || y != 0 {
        player.AimAngle = math.Atan2(y, x)
    }
}
```

## Port-zero convenience

`StickPosition` reads from controller port 0. For multiplayer input, use `PlayerStickPosition(port, deadzone)` - see [Multi-Controller Support](multi-controller.md).

## Try It

> **Download the ROM:** [`analog_demo.z64`](../emulator/roms/analog_demo.z64) - Open in [ares](https://ares-emu.net/) with the Expansion Pak enabled.
>
> **Controls:** D-Pad = movement, A = action, B = back, Start = pause, Z = trigger

## Reference Example

See `examples/analog_demo` in the GoSprite64 repository for a visual demonstration of analog stick input with a crosshair following the stick position.
