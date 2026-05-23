# The Game Loop

> If you just finished the beginner journey, this page explains the concept behind the behavior you already saw on screen.

GoSprite64 runs your game with a fixed-timestep loop at 60 FPS. You provide the logic; the engine handles timing, input polling, and frame presentation.

## The Game Interface

Every GoSprite64 game implements the `Game` interface:

```go
type Game interface {
    Init()
    Update()
    Draw()
}
```

Each method has a specific role:

- **`Init()`** - Called once before the game loop starts. Load your resources, set up your initial state, and create your game objects here.
- **`Update()`** - Called once per logic tick at a fixed 60 Hz rate. Read input, move objects, check collisions, and update game state here. Do not draw anything in Update.
- **`Draw()`** - Called once per rendered frame. Clear the screen, draw your sprites, text, and UI here. Do not update game state in Draw.

## Starting the Game

Call `Run` with a pointer to your Game implementation:

```go
func main() {
    gosprite64.Run(&Game{})
}
```

`Run` never returns. It initializes the N64 video hardware, calls your `Init()` once, initializes audio, and then enters the main loop.

## How the Loop Works

The game loop uses a fixed-timestep accumulator pattern:

1. Measure the elapsed time since the last frame
2. Add elapsed time to an accumulator
3. While the accumulator has enough time for a tick (1/60th of a second), call `Update()` and subtract the tick duration
4. Call `Draw()` once per frame
5. Sleep for the remaining time to hit the target frame rate

This means `Update()` always runs at a consistent rate regardless of how long drawing takes. If the system falls behind, multiple `Update()` calls run before the next `Draw()` to catch up.

The relevant source code in `gameloop.go`:

```go
var TargetFPS = 60
var frameDuration = time.Second / time.Duration(TargetFPS)

func Run(g Game) {
    // ... hardware initialization ...

    g.Init()

    lastTime := rtos.Nanotime()
    accumulator := time.Duration(0)

    for {
        currentTime := rtos.Nanotime()
        elapsed := time.Duration(currentTime - lastTime)
        lastTime = currentTime
        accumulator += elapsed

        for accumulator >= frameDuration {
            updateControllerState()
            g.Update()
            accumulator -= frameDuration
        }

        beginDrawing()
        g.Draw()
        endDrawing()

        sleepDuration := frameDuration - (rtos.Nanotime() - currentTime)
        if sleepDuration > 0 {
            time.Sleep(sleepDuration)
        }
    }
}
```

Key details:

- Controller input is polled automatically before each `Update()` call
- Audio is initialized after `Init()` returns, so audio calls in `Init()` are silent no-ops
- The loop runs forever - there is no quit mechanism (the N64 has no OS to return to)

## Minimal Example

Here is the simplest possible GoSprite64 game - a solid red screen:

```go
package main

import "github.com/drpaneas/gosprite64"

type Game struct{}

func (g *Game) Init()   {}
func (g *Game) Update() {}

func (g *Game) Draw() {
    gosprite64.ClearScreenWith(gosprite64.Red)
}

func main() {
    gosprite64.Run(&Game{})
}
```

A slightly more interesting example that responds to input:

```go
package main

import "github.com/drpaneas/gosprite64"

type Game struct {
    x, y int
}

func (g *Game) Init() {
    g.x = 144
    g.y = 108
}

func (g *Game) Update() {
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadUp) {
        g.y--
    }
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadDown) {
        g.y++
    }
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft) {
        g.x--
    }
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadRight) {
        g.x++
    }
}

func (g *Game) Draw() {
    gosprite64.ClearScreen()
    gosprite64.FillRect(g.x-4, g.y-4, g.x+4, g.y+4, gosprite64.Green)
    gosprite64.DrawText("MOVE WITH D-PAD", 80, 4, gosprite64.White)
}

func main() {
    gosprite64.Run(&Game{})
}
```

## Update vs Draw

Keep these two methods cleanly separated:

| `Update()` | `Draw()` |
|-------------|----------|
| Read input | Clear the screen |
| Move objects | Draw sprites and shapes |
| Check collisions | Draw text and UI |
| Update timers | Render transitions |
| Change game state | Read-only access to state |

This separation matters because `Update()` and `Draw()` can run a different number of times per frame. If the game falls behind, `Update()` catches up with multiple calls while `Draw()` only runs once.
