# Game State Machine

Most games have multiple screens: a title screen, gameplay, pause menu, game over. The `StateMachine` manages transitions between them so you don't have to write a giant switch statement in your `Update` and `Draw` methods.

## GameState interface

Each screen implements `GameState`:

```go
type GameState interface {
    Enter()   // called when this state becomes active
    Update()  // called every frame while active
    Draw()    // called every frame while active
    Exit()    // called when this state is removed or replaced
}
```

`Enter` and `Exit` are your setup and teardown hooks. Load resources in `Enter`, release them in `Exit`. The state machine guarantees that `Enter` is called before any `Update`/`Draw`, and `Exit` is called before another state takes over.

## Creating a StateMachine

Create the machine with an initial state, then call `Init` to trigger the first `Enter`:

```go
type Game struct {
    sm *gosprite64.StateMachine
}

func (g *Game) Init() {
    title := &TitleState{sm: g.sm}
    g.sm = gosprite64.NewStateMachine(title)
    title.sm = g.sm
    g.sm.Init()
}

func (g *Game) Update() { g.sm.Update() }
func (g *Game) Draw()   { g.sm.Draw() }
```

The `StateMachine` delegates `Update` and `Draw` to whichever state is on top of the stack.

## Switching screens

`Switch` replaces the current state. It calls `Exit` on the old state and `Enter` on the new one:

```go
// In your title screen:
func (s *TitleState) Update() {
    if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) {
        s.sm.Switch(&GameplayState{sm: s.sm})
    }
}
```

The lifecycle is: old.Exit() -> new.Enter() -> new.Update() on the next frame.

## Overlays with Push and Pop

`Push` adds a state on top without removing the one below. This is how you implement pause menus, dialog boxes, or inventory screens that overlay gameplay:

```go
// In gameplay, pressing Start opens the pause menu:
func (s *GameplayState) Update() {
    if gosprite64.IsButtonJustPressed(gosprite64.ButtonStart) {
        s.sm.Push(&PauseState{sm: s.sm})
        return
    }
    // normal gameplay logic...
}
```

The pause state draws on top. When the player unpauses, `Pop` removes it and returns to gameplay:

```go
func (s *PauseState) Update() {
    if gosprite64.IsButtonJustPressed(gosprite64.ButtonStart) {
        s.sm.Pop()
    }
}
```

`Pop` calls `Exit` on the removed state. The state below resumes receiving `Update` and `Draw` calls. If only one state remains, `Pop` is a no-op - the game always has at least one active state.

## Drawing overlays

When a state is pushed, only the top state receives `Update` and `Draw`. If you want the underlying state to remain visible (e.g. gameplay behind a semi-transparent pause menu), the overlay state should draw the background explicitly or you can draw all states in the stack manually.

A common pattern for pause overlays is to draw a semi-transparent box over whatever was previously rendered:

```go
func (s *PauseState) Draw() {
    // The previous frame's gameplay is still in the framebuffer.
    // Draw a darkened overlay and menu text on top.
    gosprite64.FillRect(60, 80, 228, 136, gosprite64.DarkPurple)
    gosprite64.DrawRect(60, 80, 228, 136, gosprite64.White)
    gosprite64.DrawText("PAUSED", 120, 96, gosprite64.White)
    gosprite64.DrawText("PRESS START TO RESUME", 72, 116, gosprite64.LightGray)
}
```

## Inspecting the stack

`Current` returns the active (top) state. `Depth` returns how many states are on the stack:

```go
sm.Current()  // the active GameState
sm.Depth()    // 1 = just gameplay, 2 = gameplay + pause, etc.
```

## Nil safety

Passing `nil` to `Switch` or `Push` is a no-op. This prevents crashes from conditional state transitions:

```go
var nextState gosprite64.GameState
if condition {
    nextState = &SomeState{}
}
sm.Switch(nextState)  // safe even if nextState is nil
```

## Typical game structure

A game with title, gameplay, pause, and game over uses four state types:

```go
type TitleState struct {
    sm *gosprite64.StateMachine
}

type GameplayState struct {
    sm      *gosprite64.StateMachine
    playerX float32
    score   int
}

type PauseState struct {
    sm *gosprite64.StateMachine
}

type GameOverState struct {
    sm         *gosprite64.StateMachine
    finalScore int
}
```

The flow:
1. Title -- (A pressed) --> Switch to Gameplay
2. Gameplay -- (Start pressed) --> Push Pause
3. Pause -- (Start pressed) --> Pop (back to Gameplay)
4. Gameplay -- (player dies) --> Switch to GameOver
5. GameOver -- (A pressed) --> Switch to Title

Each state has its own `Enter`, `Update`, `Draw`, `Exit`. No state knows about the internals of other states. The `StateMachine` handles the transitions.

## Try It

> **Download the ROM:** [`state_demo.z64`](../emulator/roms/state_demo.z64) - Open in [ares](https://ares-emu.net/) with the Expansion Pak enabled.
>
> **Controls:** D-Pad = movement, A = action, B = back, Start = pause, Z = trigger

## Complete example

See `examples/state_demo` for a working demo with all four states. Build it with:

```bash
GOENV=n64.env go1.24.5-embedded build -o state_demo.elf ./examples/state_demo
n64go rom state_demo.elf
```
