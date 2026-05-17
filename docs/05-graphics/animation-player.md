# Animation Player

The `AnimationPlayer` drives sprite animations using a tick-based timing model.
You define animation clips (lists of frame indices with a playback speed), then
the player steps through them each time you call `Advance`.

## Core Types

### AnimationClip

An `AnimationClip` describes a single animation sequence:

```go
type AnimationClip struct {
    Name   string
    FPS    uint16
    Frames []uint16
}
```

| Field | Description |
|---|---|
| `Name` | Human-readable label (e.g. `"idle"`, `"walk"`) |
| `FPS` | Playback speed in frames per second |
| `Frames` | Ordered list of sprite sheet frame indices |

Clips are typically loaded from a `.anim` file via an `AnimationSet`, but you
can also build them by hand:

```go
walkClip := gosprite64.AnimationClip{
    Name:   "walk",
    FPS:    10,
    Frames: []uint16{0, 1, 2, 3},
}
```

### AnimationSet

An `AnimationSet` groups related clips loaded from a compiled `.anim` file.
You get one from a bundle:

```go
animSet, err := bundle.LoadAnimation("anims")
if err != nil {
    panic(err)
}

idle, ok := animSet.Clip("idle")
walk, ok := animSet.Clip("walk")
```

`Clips()` returns all clips in the set. `Clip(name)` returns a single clip
by name along with a boolean indicating whether it was found.

## Creating a Player

```go
func NewAnimationPlayer() *AnimationPlayer
```

Creates a stopped player with no clip loaded.

```go
player := gosprite64.NewAnimationPlayer()
```

## Playback Controls

### Play

```go
func (p *AnimationPlayer) Play(clip AnimationClip)
```

Starts playing the given clip from frame 0. If the clip has no frames, the
player stops.

```go
player.Play(walkClip)
```

### Pause / Resume

```go
func (p *AnimationPlayer) Pause()
func (p *AnimationPlayer) Resume()
```

`Pause` freezes the animation at its current frame. `Resume` continues from
where it left off. Calling `Pause` when not playing (or `Resume` when not
paused) is a no-op.

### Stop

```go
func (p *AnimationPlayer) Stop()
```

Stops playback and resets to frame 0.

### Restart

```go
func (p *AnimationPlayer) Restart()
```

Jumps back to frame 0 and starts playing. Useful for retriggering an animation
without constructing a new clip.

### SetLoop

```go
func (p *AnimationPlayer) SetLoop(loop bool)
```

When looping is enabled, the animation wraps around to frame 0 after the last
frame. When disabled (the default), the player stops on the last frame.

```go
player.SetLoop(true)
player.Play(idleClip) // will loop forever
```

## Advancing Time

```go
func (p *AnimationPlayer) Advance(ticks int)
```

Call `Advance` once per game tick (typically once per frame in your `Update`
function). The player uses an internal accumulator to convert ticks into
frame steps based on the clip's FPS.

The timing math assumes a **60-tick-per-second** base rate. A clip with
`FPS: 10` advances one animation frame every 6 ticks. A clip with `FPS: 30`
advances one frame every 2 ticks.

```go
func (g *Game) Update() {
    g.player.Advance(1) // one tick per game frame
}
```

Passing `ticks <= 0` or calling `Advance` on a stopped/paused player is a
no-op.

## Reading State

### Frame

```go
func (p *AnimationPlayer) Frame() int
```

Returns the current sprite sheet frame index. Use this with `DrawSprite`:

```go
frame := player.Frame()
gosprite64.DrawSprite(sheet, frame, x, y)
```

### Playing / Done

```go
func (p *AnimationPlayer) Playing() bool
func (p *AnimationPlayer) Done() bool
```

| Method | Returns `true` when... |
|---|---|
| `Playing()` | The player is actively advancing frames |
| `Done()` | The player is stopped (finished or never started) |

These are useful for triggering events at the end of one-shot animations:

```go
if player.Done() {
    // switch to the next game state
}
```

## Complete Example

```go
type Game struct {
    sheet  *gosprite64.SpriteSheet
    player *gosprite64.AnimationPlayer
    idle   gosprite64.AnimationClip
    walk   gosprite64.AnimationClip
    flipH  bool
    moving bool
}

func (g *Game) Init() {
    sheet, err := gosprite64.LoadSpriteSheet("assets/character.sheet")
    if err != nil {
        panic(err)
    }
    g.sheet = sheet

    g.idle = gosprite64.AnimationClip{
        Name: "idle", FPS: 6, Frames: []uint16{0, 1},
    }
    g.walk = gosprite64.AnimationClip{
        Name: "walk", FPS: 10, Frames: []uint16{2, 3, 4, 5},
    }

    g.player = gosprite64.NewAnimationPlayer()
    g.player.SetLoop(true)
    g.player.Play(g.idle)
}

func (g *Game) Update() {
    g.moving = false

    if gosprite64.IsButtonDown(gosprite64.ButtonDPadRight) {
        g.moving = true
        g.flipH = false
    }
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft) {
        g.moving = true
        g.flipH = true
    }

    if g.moving {
        g.player.Play(g.walk)
    } else {
        g.player.Play(g.idle)
    }

    g.player.Advance(1)
}

func (g *Game) Draw() {
    gosprite64.ClearScreen()

    frame := g.player.Frame()
    gosprite64.DrawSpriteWithOptions(g.sheet, frame, 136, 100, gosprite64.DrawSpriteOptions{
        FlipH: g.flipH,
    })
}
```

## Compiling Animation Data

Use `mk2danim` to compile a JSON animation definition into a binary `.anim`
file:

```bash
go run github.com/drpaneas/gosprite64/cmd/mk2danim \
    -in anims.json \
    -out anims.anim
```

The JSON format defines clips with frame lists and FPS:

```json
{
  "clips": [
    { "name": "idle", "fps": 6,  "frames": [0, 1] },
    { "name": "walk", "fps": 10, "frames": [2, 3, 4, 5] }
  ]
}
```
