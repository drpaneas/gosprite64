# Input Recording and Replay

The replay system records controller input frame-by-frame and plays it back later. Use it for attract-mode demos (the game plays itself on the title screen), debugging (reproduce a bug by replaying the exact inputs), or competitive replay viewing.

## How it works

1. An `InputRecorder` captures a `FrameInput` (buttons + stick) for each player every frame.
2. When you're done, `Finish` returns a `ReplayData` containing the complete recording.
3. An `InputPlayer` feeds those frames back one at a time.

The recording is deterministic: the same sequence of `FrameInput` values always produces the same replay. If your game logic is also deterministic (same inputs = same outcome), the replay will reproduce the gameplay exactly.

## Recording

Create a recorder with the number of players, then call `CaptureFrame` once per player per frame:

```go
recorder := gosprite64.NewInputRecorder(1)  // 1 player

// In your Update():
func (g *Game) Update() {
    var buttons gosprite64.ButtonMask
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadUp) {
        buttons |= gosprite64.ButtonDPadUp
    }
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadDown) {
        buttons |= gosprite64.ButtonDPadDown
    }
    // ... other buttons

    recorder.CaptureFrame(0, gosprite64.FrameInput{
        Buttons: buttons,
        StickX:  stickX,
        StickY:  stickY,
    })
}
```

## Finishing a recording

Call `Finish` to get the replay data:

```go
replay := recorder.Finish()
// replay.FrameCount - total frames recorded
// replay.PlayerCount - number of players
```

## Playback

Create an `InputPlayer` from the replay data and read frames one at a time:

```go
player := gosprite64.NewInputPlayer(replay)

// In your Update():
input, ok := player.NextFrame(0)  // player 0
if !ok {
    // all frames consumed
    return
}

// Use input.Buttons, input.StickX, input.StickY
// to drive game logic instead of reading the controller
if input.Buttons&gosprite64.ButtonDPadUp != 0 {
    moveUp()
}
```

## Checking completion

`Done` returns true when all players have consumed all their frames:

```go
if player.Done() {
    // replay finished
}
```

## Looping playback

Call `Reset` to restart from the beginning. This is how you make an attract-mode demo that loops:

```go
if player.Done() {
    player.Reset()
    // also reset your game state to the starting position
}
```

## Multiplayer recording

Pass the number of players to `NewInputRecorder` and capture frames for each port:

```go
recorder := gosprite64.NewInputRecorder(2)

// Each frame, capture both players:
recorder.CaptureFrame(0, gosprite64.FrameInput{Buttons: p1Buttons, StickX: p1X, StickY: p1Y})
recorder.CaptureFrame(1, gosprite64.FrameInput{Buttons: p2Buttons, StickX: p2X, StickY: p2Y})
```

During playback, read each player's input separately:

```go
p1Input, _ := player.NextFrame(0)
p2Input, _ := player.NextFrame(1)
```

## FrameInput fields

| Field | Type | Description |
|-------|------|-------------|
| `Buttons` | `ButtonMask` | Bitmask of pressed buttons (same type as `ButtonA`, `ButtonDPadUp`, etc.) |
| `StickX` | `int8` | Analog stick horizontal: -128 (full left) to 127 (full right) |
| `StickY` | `int8` | Analog stick vertical: -128 (full down) to 127 (full up) |

## Typical attract-mode pattern

```go
type TitleState struct {
    sm       *gosprite64.StateMachine
    player   *gosprite64.InputPlayer
    ghostX   float32
    ghostY   float32
    demoData *gosprite64.ReplayData
}

func (s *TitleState) Enter() {
    // Pre-recorded demo data (could be loaded from cartridge FS)
    s.player = gosprite64.NewInputPlayer(s.demoData)
    s.ghostX = 144
    s.ghostY = 108
}

func (s *TitleState) Update() {
    // Play the demo in the background
    input, ok := s.player.NextFrame(0)
    if !ok {
        s.player.Reset()
        s.ghostX = 144
        s.ghostY = 108
        return
    }
    s.ghostX += float32(input.StickX) * 2
    s.ghostY += float32(input.StickY) * 2

    // Player presses Start to actually begin
    if gosprite64.IsButtonJustPressed(gosprite64.ButtonStart) {
        s.sm.Switch(&GameplayState{sm: s.sm})
    }
}

func (s *TitleState) Draw() {
    gosprite64.ClearScreen()
    // Draw the demo ghost
    gosprite64.FillRect(int(s.ghostX)-4, int(s.ghostY)-4, int(s.ghostX)+4, int(s.ghostY)+4, gosprite64.DarkGray)
    // Draw title text on top
    gosprite64.DrawText("MY GAME", 112, 40, gosprite64.White)
    gosprite64.DrawText("PRESS START", 100, 160, gosprite64.Yellow)
}
```

## Try It

<iframe src="../emulator/play.html?rom=replay_demo.z64" width="640" height="480" frameborder="0" allow="autoplay" style="display:block;margin:0 auto;max-width:100%;"></iframe>

> **Controls:** Arrow keys = D-Pad, X = A button, C = B button, Enter = Start, Z = Z trigger

## Complete example

See `examples/replay_demo` for a working demo that lets you record input and then watch it play back as a ghost trail. Build it with:

```bash
GOENV=n64.env go1.24.5-embedded build -o replay_demo.elf ./examples/replay_demo
n64go rom replay_demo.elf
```
