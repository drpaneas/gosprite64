# Sequence Player

Play MIDI-style sequenced music using the N64's audio hardware.

```go
import "github.com/drpaneas/gosprite64/audio/sequence"
```

## How sequences work

N64 games do not stream audio like modern platforms. Instead, music is stored as a compact sequence of note events (similar to MIDI) that drive instrument samples loaded into audio RAM. The sequence player reads these events and produces sound through the RSP audio microcode.

This approach uses very little ROM space - a full song might be only a few kilobytes of sequence data plus shared instrument samples.

## Creating a player

```go
player := sequence.NewPlayer()
```

`NewPlayer` returns a player with sensible defaults:

- Tempo: 120 BPM
- Master volume: 127 (maximum)
- All 16 channel volumes: 127

## Loading sequence data

Assign raw sequence bytes to the player's `Data` field before calling `Play`:

```go
player.Data = sequenceBytes
```

The sequence data is typically loaded from a ROM asset at startup.

## Playback controls

### Play

Starts (or restarts) playback from the beginning:

```go
player.Play()
```

If the player is already playing, this resets to position 0 and starts over.

### Stop

Halts playback and resets position to the beginning:

```go
player.Stop()
```

### Pause and Resume

Pause freezes playback at the current position. Resume continues from where it left off:

```go
player.Pause()
// ...later...
player.Resume()
```

Calling `Pause` when not playing has no effect. Calling `Resume` when not paused has no effect.

### IsPlaying

Returns `true` if the player is actively producing audio (playing and not paused):

```go
if player.IsPlaying() {
    drawMusicIcon()
}
```

## Tempo

`SetTempo` controls playback speed in beats per minute:

```go
player.SetTempo(140)    // faster
player.SetTempo(80)     // slower

bpm := player.Tempo()   // read current tempo
```

The default is 120 BPM. The internal tick resolution is 48 ticks per beat.

## Volume

### Master volume

`SetVolume` sets the overall output level (0-127):

```go
player.SetVolume(100)   // slightly quieter
player.SetVolume(0)     // muted

vol := player.Volume()  // read current volume
```

### Per-channel volume

`SetChannelVolume` controls individual channels (0-15). Useful for muting or emphasizing specific instrument parts:

```go
player.SetChannelVolume(9, 0)    // mute channel 9 (often drums)
player.SetChannelVolume(0, 80)   // soften the melody channel

vol := player.ChannelVolume(0)   // read channel volume
```

## Looping

`SetLoop` configures loop behavior. The first parameter is the byte position to loop back to. The second is the loop count: use -1 for infinite looping, or a positive number for a fixed repeat count.

```go
player.SetLoop(0, -1)      // loop the whole song forever
player.SetLoop(1024, 3)    // after reaching the end, jump to byte 1024, repeat 3 times
player.SetLoop(0, 0)       // no looping (play once and stop)
```

## Advancing playback

Call `Tick` once per audio frame in your game loop, passing the audio sample rate. It advances the internal position based on the current tempo and returns any note events generated during the tick:

```go
func (g *Game) Update() {
    events := g.musicPlayer.Tick(32000) // N64 typical: 32000 Hz
    for _, ev := range events {
        // Feed events to the audio mixer
        if ev.On {
            g.mixer.NoteOn(ev.Channel, ev.Note, ev.Velocity)
        } else {
            g.mixer.NoteOff(ev.Channel, ev.Note)
        }
    }
}
```

## NoteEvent

Each event returned by `Tick` is a `NoteEvent`:

```go
type NoteEvent struct {
    Channel  uint8   // MIDI channel (0-15)
    Note     uint8   // MIDI note number (0-127)
    Velocity uint8   // Strike intensity (0-127)
    On       bool    // true = note on, false = note off
}
```

## Complete example

```go
func NewGame() *Game {
    g := &Game{}

    g.music = sequence.NewPlayer()
    g.music.Data = loadAsset("overworld.seq")
    g.music.SetLoop(0, -1)
    g.music.SetVolume(100)
    g.music.Play()

    return g
}

func (g *Game) Update() {
    events := g.music.Tick(32000)
    for _, ev := range events {
        handleNoteEvent(ev)
    }

    // Pause music when game is paused
    if gs.IsButtonJustPressed(gs.ButtonStart) {
        if g.paused {
            g.music.Resume()
        } else {
            g.music.Pause()
        }
        g.paused = !g.paused
    }
}
```
