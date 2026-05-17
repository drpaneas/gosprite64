# Step 9: Screen Transitions

Add fade-in and fade-out transitions between game states.

## What you will learn

- Creating transitions with `StartTransition`
- The `FadeToBlack` and `FadeFromBlack` transition styles
- Advancing and drawing transitions each frame
- Checking `Done()` to clean up finished transitions

## What changed from Step 8

Both states now have a `fade *gosprite64.Transition` field. The title screen fades in when it appears and fades out when the player presses Start. The gameplay state fades in when it starts. The changes are small but the visual difference is significant.

### TitleState changes

```go
type TitleState struct {
	sm    *gosprite64.StateMachine
	blink int
	show  bool
	fade  *gosprite64.Transition
}

func (s *TitleState) Enter() {
	s.show = true
	s.fade = gosprite64.StartTransition(gosprite64.FadeFromBlack, 30)
}
```

On Enter, a 30-frame fade-from-black plays. The screen starts fully black and gradually reveals the title over half a second.

When the player presses Start, a fade-to-black begins before switching states:

```go
if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) || gosprite64.IsButtonJustPressed(gosprite64.ButtonStart) {
	s.fade = gosprite64.StartTransition(gosprite64.FadeToBlack, 20)
	s.sm.Switch(&PlayState{sm: s.sm})
}
```

### PlayState changes

```go
func (s *PlayState) Enter() {
	// ... asset loading (unchanged) ...
	s.fade = gosprite64.StartTransition(gosprite64.FadeFromBlack, 30)
}
```

The gameplay state fades in from black when it starts. This creates a smooth visual bridge: title fades to black, gameplay fades from black.

### The transition loop

Both states share the same pattern in Update and Draw:

```go
func (s *PlayState) Update() {
	if s.fade != nil {
		s.fade.Advance()
		if s.fade.Done() {
			s.fade.Stop()
			s.fade = nil
		}
	}
	// ... rest of update ...
}

func (s *PlayState) Draw() {
	// ... draw everything ...
	if s.fade != nil {
		s.fade.Draw()
	}
}
```

`Advance()` ticks the transition forward by one frame. `Done()` returns true when all frames have played. `Stop()` deactivates it and `Draw()` renders a semi-transparent black overlay whose opacity changes each frame.

The transition must be drawn last so it appears on top of everything else.

## The complete code

The full listing is the same as Step 8 with the `fade` field added to both states. Here are the key differences:

**TitleState** - Added `fade` field. Enter starts `FadeFromBlack`. Update advances and cleans up the fade. Draw renders the fade overlay last. The button press starts `FadeToBlack`.

**PlayState** - Added `fade` field. Enter starts `FadeFromBlack`. Update advances and cleans up the fade. Draw renders the fade overlay last.

## Transition API reference

| Function/Method | What it does |
|----------------|-------------|
| `StartTransition(style, frames)` | Create and start a new transition |
| `Advance()` | Tick forward one frame |
| `Done()` | True when all frames have played |
| `Stop()` | Deactivate the transition |
| `Draw()` | Render the overlay (call last in Draw) |

| Style | Effect |
|-------|--------|
| `FadeToBlack` | Screen gradually goes black (alpha 0 to 255) |
| `FadeFromBlack` | Screen gradually reveals (alpha 255 to 0) |

The duration is in frames. At 60 FPS, 30 frames is half a second and 60 frames is one full second.

## Build and run

```bash
go generate ./examples/platformer
GOENV=n64.env go1.24.5-embedded build -o examples/platformer/game.elf ./examples/platformer
```

The title screen now fades in from black when the game starts. Press Start and the screen fades to black, then the gameplay fades in. The transitions are short (20-30 frames) to feel snappy rather than sluggish.
