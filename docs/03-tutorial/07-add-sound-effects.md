# Step 7: Add Sound Effects

Add audio to the game using the VADPCM audio pipeline.

## What you will learn

- The GoSprite64 audio pipeline (WAV to compressed VADPCM)
- Setting up `audiogen` code generation
- Calling `PlaySoundEffect` on game events
- The generated `sfx` package with typed IDs

## How the audio pipeline works

GoSprite64 compresses audio at build time using VADPCM (the same compression family the original N64 used). You provide standard `.wav` files, the `audiogen` tool compresses them, and the engine plays them at runtime. Your game code never deals with codecs or sample rates.

The pipeline:

```text
.wav files --> audiogen --> VADPCM compressed --> embedded in ROM --> pure Go mixer --> N64 DAC
```

## Adding audio to the platformer

The platformer example does not ship with audio assets, but adding sound effects follows a standard pattern. Here is how you would add a jump sound.

### 1. Create the audio directories

```bash
mkdir -p examples/platformer/assets/audio/sfx
```

### 2. Add a WAV file

Place a 16-bit PCM WAV file at `examples/platformer/assets/audio/sfx/jump.wav`. Any sample rate works - `audiogen` resamples automatically. Stereo is downmixed to mono.

### 3. Add the audiogen generate line

Add this to the top of `main.go`, alongside the existing `go:generate` lines:

```go
//go:generate go run github.com/drpaneas/gosprite64/cmd/audiogen -dir .
```

### 4. Run code generation

```bash
go generate ./examples/platformer
```

This produces:

| Generated file | Purpose |
|----------------|---------|
| `sfx/ids.go` | Typed constants like `sfx.Jump` |
| `audio_embed.go` | Registers compressed audio with the engine at startup |
| `build/audio_v1.bin` | Compressed audio data |
| `build/audio_v1_aux.bin` | VADPCM predictor coefficients |

### 5. Play sounds from game code

Import the generated `sfx` package and call `PlaySoundEffect`:

```go
import (
	"github.com/drpaneas/gosprite64"
	"github.com/yourname/mygame/sfx"
)

func (g *Game) Update() {
	// ... movement code ...

	if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) {
		gosprite64.PlaySoundEffect(sfx.Jump)
	}
}
```

`PlaySoundEffect` returns `true` if the sound was accepted, `false` if the engine was not ready or the command ring was full. Sound effects are fire-and-forget - they play once and overlap naturally. The same effect can play up to 4 simultaneous instances.

### Playing background music

Music works the same way. Place `.wav` files in `assets/audio/music/` and call `PlayMusic`:

```go
import "github.com/yourname/mygame/music"

func (g *Game) Init() {
	gosprite64.PlayMusic(music.Overworld)
}
```

Music always loops. Call `gosprite64.StopMusic()` to stop it.

### Volume control

```go
gosprite64.SetSoundEffectVolume(0.5) // SFX at half volume
gosprite64.SetMusicVolume(0.8)       // music at 80%
```

Volume is a float32 from 0.0 (silent) to 1.0 (full). Music and SFX volumes are independent.

## WAV file requirements

| Requirement | Details |
|-------------|---------|
| Format | PCM WAV, 16-bit samples |
| Channels | Mono or stereo (stereo downmixed automatically) |
| Sample rate | Any (resampled to 16 kHz for SFX, 22 kHz for music) |
| Duration | Keep SFX short (under 2 seconds) for best compression |

## How the Pong example uses audio

The `examples/pong` directory shows a complete audio setup. It has six sound effects (paddle hits, wall bounce, scoring) that compress from 663 KB of raw PCM down to 32 KB of VADPCM - a 20x reduction.

```go
func (g *Game) Update() {
	if collide(g.ball, g.player) {
		gosprite64.PlaySoundEffect(sfx.PaddlePlayer)
	}
	if g.ball.y <= courtTop || g.ball.y >= courtBottom {
		gosprite64.PlaySoundEffect(sfx.Wall)
	}
}
```

The generated `sfx` package provides typed constants. No strings, no maps, no runtime lookups.

## Build and run

If you added audio assets and the `audiogen` generate line:

```bash
go generate ./examples/platformer
GOENV=n64.env go1.24.5-embedded build -o examples/platformer/game.elf ./examples/platformer
```

The game code from Step 6 is unchanged apart from the new `go:generate` line and the `PlaySoundEffect` call. Everything else - camera following, animation, input - works exactly as before.

For the remaining tutorial steps we will continue without audio assets to keep the code simple. You can add sound effects to any step by following the pattern above.
