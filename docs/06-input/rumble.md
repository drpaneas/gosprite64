# Rumble

Trigger controller vibration for haptic feedback using the Rumble Pak.

```go
import gs "github.com/drpaneas/gosprite64"
```

## What is the Rumble Pak?

The Rumble Pak is an accessory that plugs into the back of an N64 controller. When activated, a small motor vibrates the controller, giving the player physical feedback during gameplay. The N64 was one of the first consoles to support force feedback.

The Rumble Pak occupies the same slot as the Controller Pak (memory card). A controller can have one or the other plugged in at a time, but not both.

## SetRumble

```go
func SetRumble(port int, enabled bool)
```

`SetRumble` turns the rumble motor on or off for the controller at the given port (0-3).

- `port` - Controller port, 0 through 3.
- `enabled` - `true` to start vibrating, `false` to stop.

The function is a no-op if the port is out of range or the controller is not connected. If no Rumble Pak is inserted, the call has no effect.

Rumble stays on until you explicitly turn it off. Always pair an "on" call with a later "off" call, or the controller will vibrate indefinitely.

## Basic usage

```go
// Start rumble
gs.SetRumble(0, true)

// Stop rumble
gs.SetRumble(0, false)
```

## Collision feedback example

A common pattern is to rumble briefly when the player takes damage. Use a frame counter to control the duration:

```go
type Game struct {
    rumbleTimer int
}

func (g *Game) Update() {
    // When the player collides with an enemy
    if playerHit {
        gs.SetRumble(0, true)
        g.rumbleTimer = 10  // rumble for 10 frames
    }

    // Count down and stop
    if g.rumbleTimer > 0 {
        g.rumbleTimer--
        if g.rumbleTimer == 0 {
            gs.SetRumble(0, false)
        }
    }
}
```

## Multiplayer rumble

In a multiplayer game, rumble the controller of the player who was hit:

```go
func (g *Game) OnPlayerHit(port int) {
    gs.SetRumble(port, true)
    g.rumbleTimers[port] = 15
}

func (g *Game) Update() {
    for port := 0; port < gs.MaxControllers; port++ {
        if g.rumbleTimers[port] > 0 {
            g.rumbleTimers[port]--
            if g.rumbleTimers[port] == 0 {
                gs.SetRumble(port, false)
            }
        }
    }
}
```

## Tips

- Keep rumble bursts short (5-15 frames). Constant vibration is annoying and drains batteries.
- Vary the duration to convey intensity: a light bump might rumble for 3 frames, a heavy hit for 12.
- Always stop rumble when pausing or transitioning screens. A vibrating controller during a pause menu is distracting.
- The Rumble Pak runs on two AAA batteries. Excessive use drains them faster.
