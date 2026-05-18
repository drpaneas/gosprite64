# Timers

Games are full of "wait N frames, then do something" patterns: countdown before a round starts, cooldown between attacks, flashing a hit sprite for 10 frames, blinking a cursor every half second. The `Timer` and `RepeatingTimer` types handle these without manual frame counters.

## Timer

A `Timer` counts a fixed number of frames, then stops.

```go
// Wait 60 frames (1 second at 60 FPS) before dropping the next pill
dropDelay := gosprite64.NewTimer(60)

// In Update():
if dropDelay.Tick() {
    // This runs on exactly the frame the timer finishes
    dropPill()
    dropDelay.Reset()  // start the next delay
}
```

### Checking state

```go
timer.Done()       // true after all frames elapsed
timer.Elapsed()    // frames so far
timer.Remaining()  // frames left
timer.Duration()   // total frame count
timer.Progress()   // 0.0 to 1.0 ratio
```

### Progress for animation

`Progress` returns a 0..1 value that pairs naturally with easing functions:

```go
// Fade out over 30 frames
fadeTimer := gosprite64.NewTimer(30)

// In Update():
fadeTimer.Tick()

// In Draw():
alpha := 1.0 - fadeTimer.Progress()
// or with easing:
alpha := 1.0 - math2d.EaseOutQuad(fadeTimer.Progress())
```

### Reset

`Reset` restarts with the same duration. `ResetWith` changes the duration:

```go
timer.Reset()          // restart same countdown
timer.ResetWith(120)   // restart with 2 seconds
```

### Zero and negative durations

A zero-duration timer is immediately `Done`. Negative durations are clamped to 0:

```go
gosprite64.NewTimer(0).Done()   // true
gosprite64.NewTimer(-5).Done()  // true (clamped to 0)
```

## RepeatingTimer

A `RepeatingTimer` fires at a fixed interval and counts how many times it has triggered.

```go
// Blink cursor every 30 frames (twice per second)
blink := gosprite64.NewRepeatingTimer(30)

// In Update():
blink.Tick()

// In Draw():
if blink.Count()%2 == 0 {
    drawCursor()
}
```

### Spawn waves

```go
// Spawn an enemy every 90 frames
spawner := gosprite64.NewRepeatingTimer(90)

// In Update():
if spawner.Tick() {
    spawnEnemy()
}
```

### Count and Reset

```go
spawner.Count()  // how many times it has triggered
spawner.Reset()  // clear count and elapsed
```

## Typical patterns

### Countdown before round start

```go
type PlayState struct {
    countdown *gosprite64.Timer
    started   bool
}

func (s *PlayState) Enter() {
    s.countdown = gosprite64.NewTimer(180) // 3 seconds
}

func (s *PlayState) Update() {
    if !s.countdown.Done() {
        s.countdown.Tick()
        return  // skip gameplay logic during countdown
    }
    s.started = true
    // normal gameplay...
}

func (s *PlayState) Draw() {
    if !s.countdown.Done() {
        remaining := 3 - s.countdown.Elapsed()/60
        gosprite64.DrawText(fmt.Sprintf("%d", remaining), 140, 100, gosprite64.Yellow)
    }
}
```

### Flash on hit

```go
hitFlash := gosprite64.NewTimer(10)

// When hit:
hitFlash.Reset()

// In Draw():
if !hitFlash.Done() {
    hitFlash.Tick()
    if hitFlash.Elapsed()%2 == 0 {
        drawPlayer()  // skip every other frame for flash effect
    }
} else {
    drawPlayer()
}
```

## Try It

> **Download the ROM:** [`timer_demo.z64`](../emulator/roms/timer_demo.z64) - Open in [ares](https://ares-emu.net/) with the Expansion Pak enabled.
>
> **Controls:** D-Pad = movement, A = action, B = back, Start = pause, Z = trigger

## Complete example

See `examples/timer_demo` for a focused timer demonstration, or `examples/splitscreen_demo` for timers used alongside other systems. Build with:

```bash
GOENV=n64.env go1.24.5-embedded build -o timer.elf ./examples/timer_demo
n64go rom timer.elf
```
