# Step 8: Add a Title Screen

Use the StateMachine to add a title screen before gameplay begins.

## What you will learn

- The `GameState` interface (Enter, Update, Draw, Exit)
- Creating a `StateMachine` to manage game states
- Switching between a title screen and gameplay
- Drawing text with `DrawText`

## The code

The single `Game` struct from Step 6 has been split into three types. All gameplay fields and asset loading moved into `PlayState`. A new `TitleState` draws the start screen. The top-level `Game` just owns the state machine.

### TitleState (new)

```go
type TitleState struct {
	sm    *gosprite64.StateMachine
	blink int
	show  bool
}

func (s *TitleState) Enter() {
	s.show = true
}

func (s *TitleState) Update() {
	s.blink++
	if s.blink%30 == 0 {
		s.show = !s.show
	}
	if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) || gosprite64.IsButtonJustPressed(gosprite64.ButtonStart) {
		s.sm.Switch(&PlayState{sm: s.sm})
	}
}

func (s *TitleState) Draw() {
	gosprite64.ClearScreenWith(gosprite64.DarkBlue)
	gosprite64.DrawText("PLATFORMER", 104, 60, gosprite64.White)
	gosprite64.DrawText("A GoSprite64 Tutorial", 60, 80, gosprite64.LightGray)
	if s.show {
		gosprite64.DrawText("PRESS START", 96, 140, gosprite64.Yellow)
	}
}

func (s *TitleState) Exit() {}
```

### PlayState (moved from Game)

`PlayState` is the same code as Step 6's `Game` struct. The only difference is that it now implements `GameState` instead of the top-level Game interface, and it holds a reference to the state machine (`sm`). The `Enter` method replaces `Init` - asset loading happens there. `Update` and `Draw` are unchanged. `Exit` is empty.

### Game (new top-level wrapper)

```go
type Game struct {
	sm *gosprite64.StateMachine
}

func (g *Game) Init() {
	title := &TitleState{}
	g.sm = gosprite64.NewStateMachine(title)
	title.sm = g.sm
	g.sm.Init()
}

func (g *Game) Update() { g.sm.Update() }
func (g *Game) Draw()   { g.sm.Draw() }
```

## How it works

### The GameState interface

Every state implements four methods:

| Method | When called |
|--------|-----------|
| `Enter()` | State becomes active |
| `Update()` | Every frame |
| `Draw()` | Every frame after Update |
| `Exit()` | State is replaced |

### The StateMachine

`NewStateMachine(title)` creates the machine with an initial state. `Init()` calls `Enter` on it. From then on `Game` just forwards `Update` and `Draw` to the machine. `Switch` calls `Exit` on the current state, replaces it, and calls `Enter` on the new one.

### The blinking text

A frame counter toggles "PRESS START" visibility every 30 frames (half a second at 60 FPS). `DrawText` renders a string at a screen pixel position in the given color using the engine's built-in font.

## Build and run

```bash
go generate ./examples/platformer
GOENV=n64.env go1.24.5-embedded build -o examples/platformer/game.elf ./examples/platformer
```

The game now starts on a dark blue title screen with "PLATFORMER" and blinking "PRESS START" text. Press A or Start to switch to gameplay. All the movement, animation, and camera following from Step 6 works exactly as before.
