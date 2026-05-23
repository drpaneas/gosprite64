# Hello World

This page explains the small GoSprite64 program behind the beginner journey.

If you want the shortest route to visible progress, start with [Run Your First ROM](../02-first-journey/01-run-your-first-rom.md) instead.

## What this page is for

Use this page when you want to understand:

- the smallest standalone GoSprite64 program
- what `Init`, `Update`, and `Draw` each do
- how `n64.env` and ROM generation fit together

## Prerequisites

Complete the [Installation](./installation.md) guide first. You need:

- `go` (standard Go, for dependency resolution)
- `go1.24.5-embedded` (EmbeddedGo toolchain, for building)
- `n64go` (ROM tool)

## Create the project

```bash
mkdir -p ~/gocode/src/github.com/yourname/mygame
cd ~/gocode/src/github.com/yourname/mygame
```

## Initialize the module

```bash
go mod init github.com/yourname/mygame
```

## Write main.go

Create `main.go` with the following content:

```go
package main

import "github.com/drpaneas/gosprite64"

type Game struct{}

func (g *Game) Init()   {}
func (g *Game) Update() {}
func (g *Game) Draw()   { gosprite64.ClearScreenWith(gosprite64.Blue) }

func main() { gosprite64.Run(&Game{}) }
```

Every GoSprite64 game implements three methods on a struct:

- `Init()` runs once at startup
- `Update()` runs every frame for game logic
- `Draw()` runs every frame for rendering

`gosprite64.Run()` starts the game loop and never returns.

GoSprite64 exposes one official fixed resolution and drawing space: `288x216` logical pixels. `gosprite64.ClearScreen()` is the frame-start background clear, while drawing helpers such as `gosprite64.FillRect`, `gosprite64.DrawRect`, `gosprite64.DrawLine`, and `gosprite64.DrawText` use logical coordinates inside that fixed canvas.

## Add the n64.env file

Create `n64.env` in the project root:

```
GOTOOLCHAIN=go1.24.5-embedded
GOOS=noos
GOARCH=mips64
GOFLAGS='-tags=n64' '-trimpath' '-toolexec=n64go toolexec' '-ldflags=-M=0x00000000:8M -F=0x00000400:8M -stripfn=1'
```

This tells the Go toolchain to cross-compile for the N64 (MIPS64, no OS).

## Resolve dependencies

Run `go mod tidy` with a clean Go environment to avoid interference from any inherited N64 variables:

```bash
env -u GOENV -u GOOS -u GOARCH -u GOFLAGS -u GOTOOLCHAIN go mod tidy
```

This downloads GoSprite64 and its transitive dependencies.

## Build

Compile the project and produce the ROM:

```bash
GOENV=n64.env go1.24.5-embedded build -o game.elf .
GOENV=n64.env n64go rom game.elf
```

The first command cross-compiles your code for the N64 (MIPS64, no OS) using the settings in `n64.env`, producing `game.elf`. The second converts the ELF into an N64 ROM (`game.z64`).

## Run

Load `game.z64` in an emulator like [ares](https://ares-emu.net/) to see a blue screen. When you want to inspect the canvas boundaries and square-pixel presentation, compare it with the repository's `examples/calibration` ROM.

## Editor support

If `gopls` reports `embedded/*` packages as missing or does not recognize files guarded by `//go:build n64`, see [Editor Setup](./editor-setup.md).

## Project structure

Your project should now look like this:

```
mygame/
  main.go       # your game code
  go.mod        # module definition + dependencies
  go.sum        # dependency checksums (auto-generated)
  n64.env       # N64 build target configuration
  game.elf      # compiled binary (after build)
  game.z64      # N64 ROM (after build)
```

## Next steps

Now that your toolchain works, try changing `gosprite64.ClearScreenWith(gosprite64.Blue)` to another color like `gosprite64.Red`, `gosprite64.Green`, or `gosprite64.DarkPurple` and rebuild. See `examples/clearscreen` in the repository for this exact pattern. Then explore the [examples](https://github.com/drpaneas/gosprite64/tree/main/examples) in the GoSprite64 repository to learn about input handling, drawing shapes, text rendering, and audio.
