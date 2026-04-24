# AudioGen

AudioGen is a tool for embedding audio files into your gosprite64 application. It generates an `audio_embed.go` file that includes all the audio files in your project.

## Supported Audio Formats

- Music source files: `music0.wav`, `music1.wav`, ..., `music63.wav`
- Music runtime files: `music0.raw`, `music1.raw`, ..., `music63.raw`
- Sound effect source files: `sfx_*.wav` (e.g., `sfx_jump.wav`, `sfx_shoot.wav`)
- Sound effect runtime files: `sfx_*.raw` (e.g., `sfx_jump.raw`, `sfx_shoot.raw`)

## Installation

AudioGen is a host-only tool that runs on your development machine (macOS, Linux, Windows), not on the N64. Build it with your regular Go toolchain, not with `go1.24.5-embedded`:

```bash
go install github.com/drpaneas/gosprite64/cmd/audiogen@latest
```

Or build it locally from the repository:

```bash
go build -o audiogen ./cmd/audiogen
```

## Usage

1. Place your audio files in your project directory
2. Run the following command in your project directory:

```bash
# From your project root
go generate ./...
```

3. Add this comment to one of your `.go` files (e.g., `main.go`):

```go
//go:generate go run github.com/drpaneas/gosprite64/cmd/audiogen -dir .
```

## Playing Audio in Your Game

### Music

Music tracks are automatically loaded with IDs 0-63 based on their filename (e.g., `music0.raw` has ID 0).

```go
// Play music track 0 in a loop
gosprite64.Music(0, true)

// Stop all music
gosprite64.Music(-1, false)
```

### Sound Effects

Sound effects can be played by their name (without the `sfx_` prefix and `.raw` extension):

```go
// Play a sound effect named "jump.raw"
gosprite64.PlaySFX("jump")
```

The normal workflow is to keep `.wav` files as source assets and let `audiogen` generate the `.raw` runtime files for embedding.

## Building for N64

When building for N64, the audio files will be embedded into the ROM. Make sure to include the generated `audio_embed.go` file in your build.

## License

MIT
