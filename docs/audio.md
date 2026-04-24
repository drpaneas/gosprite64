# Using Audio in GoSprite64

This chapter explains how audio works in GoSprite64 today, what files you should keep in your game project, and which APIs you should call from gameplay code.

## The mental model

GoSprite64 uses a two-step audio pipeline:

1. You author or keep your source assets as `.wav` files.
2. The build-time `audiogen` tool converts those assets into runtime `.raw` files and generates `audio_embed.go`.

At runtime, GoSprite64 does **not** decode WAV containers. It plays raw PCM audio that has already been converted into the one format the engine expects.

## The runtime audio format

Runtime audio assets must be:

- signed 16-bit PCM
- stereo
- 48 kHz
- big-endian
- interleaved left/right samples

That is the format stored in `music*.raw` and `sfx_*.raw` after conversion.

If you load raw audio yourself with `LoadAudio`, the byte slice you provide must already follow that exact format.

## The normal workflow

The intended workflow looks like this:

1. Add source audio files to your game directory using names like `music0.wav` for numbered music tracks or `sfx_jump.wav` for named one-shot effects.
2. Add a `go:generate` line to your game:

```go
//go:generate go run github.com/drpaneas/gosprite64/cmd/audiogen -dir .
```

3. Run:

```bash
go generate ./...
```

4. Build your game as usual.

`audiogen` will:

- discover your audio files
- convert supported WAV sources into runtime `.raw`
- generate `audio_embed.go`
- call `gosprite64.SetAudioFS(...)` from the generated file so the runtime can find the embedded assets

## Which files matter

GoSprite64 currently recognizes these filename patterns in the selected directory:

- `music*.wav`
- `music*.raw`
- `sfx_*.wav`
- `sfx_*.raw`

The current tool scans the directory you pass with `-dir`. It does not recursively walk subdirectories, so keep your audio files directly in that directory if you want them included.

When both a `.wav` and a `.raw` exist for the same logical asset, the `.wav` is treated as the source and the `.raw` is regenerated from it.

In practice, this means:

- `.wav` is the editable source asset
- `.raw` is the runtime asset the N64 code actually embeds and plays

## What `audiogen` converts

The current WAV conversion path supports:

- PCM WAV
- 16-bit samples
- mono or stereo input

During conversion:

- mono input is duplicated to stereo
- non-48 kHz input is resampled to 48 kHz
- samples are written out as big-endian runtime PCM

If the WAV file is not PCM 16-bit, `audiogen` rejects it instead of guessing.

## Music tracks

Background music and numbered cues use the `Music` API:

```go
gosprite64.Music(0, true)  // play music0.raw in a loop
gosprite64.Music(3, false) // play music3.raw once
gosprite64.Music(-1, false) // stop all active audio
```

Track IDs come from filenames:

- `music0.raw` -> `Music(0, ...)`
- `music1.raw` -> `Music(1, ...)`
- and so on

Embedded `music*.raw` files are registered during audio initialization and loaded on first use.

## Sound effects

Named sound effects use `PlaySFX`:

```go
gosprite64.PlaySFX("jump")
```

That looks for:

```text
sfx_jump.raw
```

and loads it on demand the first time it is used.

If you are choosing between the two styles:

- use `Music(id, loop)` for numbered music tracks and fixed cue slots
- use `PlaySFX(name)` for named one-shot effects

The Pong example uses named `PlaySFX(...)` calls for its one-shot cues. That is the recommended style when you are thinking in terms of gameplay events rather than numbered track slots.

## What happens inside `Run()`

If your game uses `Run(&Game{})`, you do not need to manually initialize the audio engine.

At startup, GoSprite64:

- sets up the embedded audio filesystem via generated `audio_embed.go`
- initializes the mixer runtime once
- starts a background feeder that streams mixed audio to the N64 audio hardware

During the game loop, `UpdateAudio()` is still called, but it is now only lightweight housekeeping. It cleans up finished one-shot sounds and frees mixer channels. It does not manually push PCM chunks every frame anymore.

## Build requirements

Audio playback now depends on upstream mixer assets that are packaged through `n64go toolexec`, so your N64 build environment must include the current `GOFLAGS` setup.

The repository's `n64.env` uses:

```text
GOTOOLCHAIN=go1.24.5-embedded
GOOS=noos
GOARCH=mips64
GOFLAGS='-tags=n64' '-trimpath' '-toolexec=n64go toolexec' '-ldflags=-M=0x00000000:8M -F=0x00000400:8M -stripfn=1'
```

If that `-toolexec=n64go toolexec` part is missing, ROMs that depend on mixer cart assets can fail at startup.

## Recommended project layout

For a small game package, a practical layout is:

```text
mygame/
  main.go
  audio_embed.go
  music0.wav
  music0.raw
  sfx_jump.wav
  sfx_jump.raw
  n64.env
```

The `.wav` files are your source-of-truth assets. The `.raw` files are generated runtime artifacts.

## Troubleshooting

If audio does not work, check these first:

- `audio_embed.go` exists and is generated from your current assets
- your filenames match the expected patterns exactly
- your audio files are in the directory passed to `audiogen`
- your generated `.raw` files are non-empty
- your build uses the current `n64.env` with `n64go toolexec`

If you are loading raw PCM manually, make sure it is already:

- stereo
- 48 kHz
- signed 16-bit
- big-endian

If any of those are wrong, the result will usually sound distorted rather than quietly failing.
