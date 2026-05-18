# Screen Transitions

Transitions provide smooth visual effects when switching between game states -
for example, fading to black before loading a new level and fading back in
once it's ready.

## Transition Styles

GoSprite64 currently supports two fade styles:

```go
type TransitionStyle int

const (
    FadeToBlack   TransitionStyle = iota // screen darkens over time
    FadeFromBlack                        // screen brightens over time
)
```

| Style | Start | End |
|---|---|---|
| `FadeToBlack` | Fully visible (alpha 0) | Fully black (alpha 255) |
| `FadeFromBlack` | Fully black (alpha 255) | Fully visible (alpha 0) |

Both styles draw a semi-transparent black overlay whose opacity changes each
frame.

## Starting a Transition

```go
func StartTransition(style TransitionStyle, durationFrames int) *Transition
```

Creates and returns an active `Transition`. The `durationFrames` parameter
controls how many frames the effect takes to complete. At 60 FPS, a duration
of 30 gives a half-second fade.

```go
fade := gosprite64.StartTransition(gosprite64.FadeToBlack, 30)
```

## Advancing and Drawing

Each frame, call `Advance` to step the transition forward, then `Draw` to
render the overlay on top of your scene:

```go
func (tr *Transition) Advance()
func (tr *Transition) Draw()
```

`Advance` increments an internal frame counter. Once the counter reaches
`Duration`, the transition is finished. `Draw` renders the black overlay at
the current alpha level. If the alpha is 0 (fully transparent), `Draw` skips
rendering entirely.

```go
func (g *Game) Update() {
    if g.fade != nil {
        g.fade.Advance()
    }
}

func (g *Game) Draw() {
    gosprite64.ClearScreen()
    // ... draw your scene ...

    if g.fade != nil {
        g.fade.Draw()
    }
}
```

## Checking State

```go
func (tr *Transition) Done() bool
func (tr *Transition) Active() bool
func (tr *Transition) Stop()
```

| Method | Returns / Does |
|---|---|
| `Done()` | `true` when the transition has reached its final frame |
| `Active()` | `true` when the transition is running and not yet done |
| `Stop()` | Immediately deactivates the transition |

`Done` returns `true` on a `nil` transition, so you can safely check without
a nil guard. `Active` is the inverse: it returns `false` on `nil` or stopped
transitions.

```go
if g.fade.Done() {
    // transition finished - safe to switch states
}
```

## Complete Fade-In / Fade-Out Example

A common pattern is to fade out, switch the game state, then fade back in:

```go
type GameState int

const (
    StatePlaying GameState = iota
    StateFadingOut
    StateLoading
    StateFadingIn
)

type Game struct {
    state   GameState
    fadeOut  *gosprite64.Transition
    fadeIn   *gosprite64.Transition
    level   int
}

func (g *Game) Update() {
    switch g.state {
    case StatePlaying:
        // Normal gameplay...
        if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) {
            g.fadeOut = gosprite64.StartTransition(gosprite64.FadeToBlack, 30)
            g.state = StateFadingOut
        }

    case StateFadingOut:
        g.fadeOut.Advance()
        if g.fadeOut.Done() {
            g.state = StateLoading
        }

    case StateLoading:
        g.level++
        g.loadLevel(g.level)
        g.fadeIn = gosprite64.StartTransition(gosprite64.FadeFromBlack, 30)
        g.state = StateFadingIn

    case StateFadingIn:
        g.fadeIn.Advance()
        if g.fadeIn.Done() {
            g.state = StatePlaying
        }
    }
}

func (g *Game) Draw() {
    gosprite64.ClearScreen()

    // Always draw the game world
    g.drawWorld()

    // Draw the active transition overlay on top
    switch g.state {
    case StateFadingOut:
        g.fadeOut.Draw()
    case StateLoading:
        // Screen is fully black during load
        gosprite64.FillRect(0, 0, 287, 215, gosprite64.Black)
    case StateFadingIn:
        g.fadeIn.Draw()
    }
}
```

## Tips

- A duration of 30 frames (0.5 seconds at 60 FPS) feels snappy. A duration
  of 60 frames (1 second) feels more cinematic.
- You can `Stop()` a transition early if the player presses a button to skip.
- Transitions draw a full-screen overlay, so call `Draw` after all your scene
  rendering.
- The alpha interpolation is linear. For eased fades, you could run a shorter
  transition and manage the alpha curve yourself.
- All methods are nil-safe. Calling `Advance`, `Draw`, `Done`, `Active`, or
  `Stop` on a `nil` `Transition` is a no-op (or returns a safe default).

## Try It

> **Download the ROM:** [`fade_demo.z64`](../emulator/roms/fade_demo.z64) - Open in [ares](https://ares-emu.net/) with the Expansion Pak enabled.
>
> **Controls:** D-Pad = movement, A = action, B = back, Start = pause, Z = trigger

## Reference Example

See `examples/fade_demo` in the GoSprite64 repository for a working fade transition example.
