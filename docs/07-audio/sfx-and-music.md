# Using Audio in GoSprite64

This chapter covers how to add audio to your game and how the audio system works behind the scenes.

## Part 1: Adding audio to your game

### Quick start

GoSprite64 uses a build-time pipeline. You provide `.wav` files, the `audiogen` tool compresses them into a compact format, and the engine plays them at runtime. You never deal with codecs, sample rates, or streaming in your gameplay code.

Three steps:

1. Put your `.wav` files in the right directories.
2. Run `go generate`.
3. Call `gosprite64.PlaySoundEffect` or `gosprite64.PlayMusic` from your game.

### Project layout

Organize audio files under `assets/audio/` with separate subdirectories for sound effects and music:

```text
mygame/
  main.go
  audio_embed.go        (generated)
  assets/
    audio/
      sfx/
        jump.wav
        coin.wav
        hit.wav
      music/
        overworld.wav
        boss.wav
  sfx/
    ids.go              (generated)
  music/
    ids.go              (generated)
  build/
    audio_v1.bin         (generated)
    audio_v1_aux.bin     (generated)
    audio_report.json    (generated)
  n64.env
```

The `.wav` files are your source assets. Everything in `build/`, `sfx/`, `music/`, and `audio_embed.go` is generated and should not be edited by hand.

### Setting up audiogen

Add a `go:generate` line to your `main.go`:

```go
//go:generate go run github.com/drpaneas/gosprite64/cmd/audiogen -dir .
```

Then run:

```bash
go generate ./...
```

This will:

- scan `assets/audio/sfx/` and `assets/audio/music/` for `.wav` files
- convert each WAV to mono, resample to the target rate, and compress with VADPCM
- generate typed ID constants in `sfx/ids.go` and `music/ids.go`
- generate `audio_embed.go` which registers all assets with the engine at startup
- write a size and performance report to `build/audio_report.json`

### WAV file requirements

`audiogen` accepts:

- PCM WAV files
- 16-bit samples
- mono or stereo input (stereo is downmixed to mono automatically)
- any sample rate (resampled automatically to the target rate)

If the file is not PCM 16-bit, `audiogen` rejects it with a clear error message.

### Playing sound effects

Import the generated `sfx` package and call `gosprite64.PlaySoundEffect`:

```go
import "github.com/drpaneas/gosprite64/examples/pong/sfx"

func (g *Game) Update() {
    if playerScored {
        gosprite64.PlaySoundEffect(sfx.ScorePlayer)
    }
    if ballHitWall {
        gosprite64.PlaySoundEffect(sfx.Wall)
    }
}
```

`gosprite64.PlaySoundEffect` returns `true` if the sound was accepted, `false` if it was dropped (engine not ready or command ring full). Sound effects are one-shot and can overlap. The same effect can play up to 4 times simultaneously. If you trigger a 5th instance, the oldest one is evicted.

### Playing background music

Import the generated `music` package and call `gosprite64.PlayMusic`:

```go
import "github.com/yourname/mygame/music"

func (g *Game) Init() {
    gosprite64.PlayMusic(music.Overworld)
}
```

Music always loops. If a different track is already playing, it stops and the new one starts. Calling `gosprite64.PlayMusic` with the same track that is already playing does nothing.

To stop music:

```go
gosprite64.StopMusic()
```

### Volume control

```go
gosprite64.SetSoundEffectVolume(0.5)  // SFX at half volume
gosprite64.SetMusicVolume(0.8)        // music at 80%
```

Volume is a `float32` from 0.0 (silent) to 1.0 (full). Values outside that range are clamped. Music and SFX volumes are independent.

### Complete example: Pong

Here is how the Pong example uses audio:

```go
//go:generate go run github.com/drpaneas/gosprite64/cmd/audiogen -dir .

package main

import (
    "github.com/drpaneas/gosprite64"
    "github.com/drpaneas/gosprite64/examples/pong/sfx"
)

func (g *Game) Init() {
    switch g.Scored {
    case "Player":
        gosprite64.PlaySoundEffect(sfx.ScorePlayer)
    case "Computer":
        gosprite64.PlaySoundEffect(sfx.ScoreComputer)
    default:
        gosprite64.PlaySoundEffect(sfx.Start)
    }
}

func (g *Game) Update() {
    if collide(g.ball, g.computer) {
        gosprite64.PlaySoundEffect(sfx.PaddleComputer)
    }
    if collide(g.ball, g.player) {
        gosprite64.PlaySoundEffect(sfx.PaddlePlayer)
    }
    if g.ball.y <= courtTop || g.ball.y >= courtBottom {
        gosprite64.PlaySoundEffect(sfx.Wall)
    }
}
```

The generated `sfx/ids.go` provides typed constants like `sfx.PaddleComputer`, `sfx.Wall`, etc. No strings, no maps, no runtime lookups.

### Build budget flags

`audiogen` enforces size budgets by default. If your audio exceeds the limits, the build fails with a clear message. You can override the defaults:

```go
//go:generate go run github.com/drpaneas/gosprite64/cmd/audiogen -dir . -rom-budget=1048576 -sfx-resident-budget=65536
```

| Flag | Default | Description |
|------|---------|-------------|
| `-rom-budget` | 524,288 bytes (512 KB) | Maximum total size for all audio data in ROM |
| `-sfx-resident-budget` | 32,768 bytes (32 KB) | Maximum compressed SFX data resident in memory |

The budget report at `build/audio_report.json` shows exactly where your audio bytes are going.

## Part 2: How it works behind the scenes

### The pipeline at a glance

```text
                   BUILD TIME                          RUNTIME
  .wav files --> audiogen --> VADPCM compressed --> embedded in ROM
                                                        |
                                                   pure Go mixer
                                                        |
                                                   48 kHz stereo
                                                        |
                                                    N64 DAC out
```

At build time, `audiogen` converts WAV files into 4-bit VADPCM compressed mono audio. At runtime, a pure Go software mixer decodes, resamples, mixes, and writes stereo output to the N64 DAC at 48 kHz.

### VADPCM compression

GoSprite64 uses 4-bit VADPCM (Vector Adaptive Differential Pulse Code Modulation), the same compression family used by Nintendo for N64 audio. Each 9-byte compressed block decodes to 16 mono samples.

The codec achieves roughly 3.5:1 compression on the raw PCM data. Combined with mono downmix (2x) and lower sample rates (up to 3x), the total size reduction compared to the old 48 kHz stereo raw format is typically around 20x.

For the Pong example with 6 sound effects:

| | Old system | New system |
|---|-----------|-----------|
| Format | 48 kHz stereo raw PCM | 16 kHz mono VADPCM |
| Total audio ROM | 662,896 bytes | 31,863 bytes |
| Reduction | - | 20.8x smaller (95.2%) |

### Sample rates

Assets are resampled at build time. Gameplay code never sees sample rates.

| Asset class | Native rate | Rationale |
|-------------|------------|-----------|
| SFX | 16,000 Hz | Short effects do not need high fidelity. 16 kHz is clean 3x resampling to 48 kHz. |
| Music | 22,050 Hz | Longer tracks benefit from slightly higher quality. |
| DAC output | 48,000 Hz | Hardware output rate, unchanged from previous versions. |

### Voice model

The engine pre-allocates a fixed set of voices at startup:

- **1 music voice** - reserved, never stolen by sound effects
- **8 SFX voices** - shared pool with priority-based stealing

When all 8 SFX voices are busy and a new effect is triggered, the oldest playing SFX is evicted. Music is never a victim.

The same SFX can overlap up to 4 times. If a 5th instance is triggered, the oldest instance of that specific effect is evicted first.

### Memory usage

The engine uses a fixed memory budget that does not grow after initialization:

| Component | Size |
|-----------|------|
| Voice states | 144 bytes |
| Source decode buffers | 288 bytes |
| Source structs | 2,304 bytes |
| Mixer accumulator | 2,048 bytes |
| Output buffer | 2,048 bytes |
| Command ring | 96 bytes |
| **Fixed runtime total** | **6,928 bytes** |

SFX data stays resident in compressed form. For the Pong example, that is 31 KB of compressed audio in ROM versus 663 KB of raw PCM that the old system loaded into RDRAM.

### Zero allocations after init

Every hot-path operation runs with zero heap allocations:

| Operation | Time (Apple M2 Pro) | Allocations |
|-----------|-------------------|-------------|
| Decode one VADPCM block (16 samples) | 135 ns | 0 |
| Mix 9 voices, 512 output frames | 9.0 us | 0 |
| Fill 256 decoded samples from source | 2.4 us | 0 |
| Command ring push + pop | 15 ns | 0 |

The old system allocated memory on the first play of every SFX (about 114 KB per 50 KB asset for the read + cache copy).

### Concurrency

Gameplay code and the audio feeder run in separate goroutines. They communicate through a lock-free single-producer single-consumer command ring. Gameplay calls like `gosprite64.PlaySoundEffect` push a small command struct into the ring. The feeder drains all pending commands at the top of each fill cycle before touching any voice state. No mutexes are used in the audio hot path.

### The feeder loop

A background goroutine runs continuously:

1. Drain all pending commands from the ring (play, stop, volume changes).
2. For each active voice, decode enough VADPCM blocks to fill the source buffer.
3. The mixer resamples each voice from its native rate to 48 kHz using linear interpolation, applies bus gain, and mixes into a stereo output buffer.
4. Write the output buffer to the N64 DAC. The hardware write blocks until the DAC is ready for more data, which naturally paces the loop at the output sample rate.

### Anti-click ramp

When music is stopped via `gosprite64.StopMusic()`, the engine applies a short linear ramp to zero over 1-2 ms (about 22-44 samples) before releasing the voice. This prevents the audible click that would otherwise occur from abruptly zeroing a playing waveform.

### Loop handling

Music assets always loop. Loop boundaries are enforced at build time by `audiogen`:

- Only forward loops are supported
- Loop start and length must be aligned to 16 decoded samples (one VADPCM block)
- The decoder state at the loop start is captured and saved during encoding
- At runtime, when the source reaches the loop end, it restores the saved decoder state and jumps to the loop start

This avoids any "decode from the beginning" logic or runtime loop repair.

### Build-time budget report

Every `audiogen` run produces `build/audio_report.json` with metrics including:

- Total ROM bytes used by audio
- SFX resident vs music streamed breakdown
- Estimated streaming bandwidth for active music (about 12.4 KB/sec at 22,050 Hz VADPCM)
- Estimated decode CPU cost per 10 ms of audio
- Fixed runtime RAM breakdown by component

If any hard budget limit is exceeded, `audiogen` exits with a non-zero status so CI catches the problem.

### Troubleshooting

If audio does not work, check these first:

- `audio_embed.go` exists and was generated from your current assets
- your `.wav` files are in `assets/audio/sfx/` or `assets/audio/music/`
- your `.wav` files are 16-bit PCM (not 24-bit, not float, not compressed)
- the `build/` directory contains `audio_v1.bin` and `audio_v1_aux.bin`
- your build uses the current `n64.env` with `n64go toolexec`
- check `build/audio_report.json` for budget violations

If you hear distortion or clicking, verify that `go generate` ran successfully after your last asset change.

## Try It

<iframe src="../emulator/play.html?rom=pong.z64" width="640" height="480" frameborder="0" allow="autoplay" style="display:block;margin:0 auto;max-width:100%;"></iframe>

> **Controls:** Arrow keys = D-Pad, X = A button, C = B button, Enter = Start, Z = Z trigger
