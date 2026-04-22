# Hello World

This guide walks you through creating a standalone N64 game project from scratch using GoSprite64. By the end you will have a blue screen ROM that proves your toolchain is set up correctly.

## Prerequisites

Complete the [Getting Started](./getting_started.md) guide first. You need:

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

import . "github.com/drpaneas/gosprite64"

type Game struct{}

func (g *Game) Init()   {}
func (g *Game) Update() {}
func (g *Game) Draw()   { ClearScreen(Blue) }

func main() { Run(&Game{}) }
```

Every GoSprite64 game implements three methods on a struct:

- `Init()` runs once at startup
- `Update()` runs every frame for game logic
- `Draw()` runs every frame for rendering

`Run()` starts the game loop and never returns.

## Add the n64.env file

Create `n64.env` in the project root:

```
GOTOOLCHAIN=go1.24.5-embedded
GOOS=noos
GOARCH=mips64
GOFLAGS='-tags=n64' '-trimpath' '-ldflags=-M=0x00000000:8M -F=0x00000400:8M -stripfn=1'
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

Load `game.z64` in an emulator like [ares](https://ares-emu.net/) to see a blue screen.

## Editor support

No special editor configuration is needed. The GoSprite64 source files have no build tags, so gopls and the VS Code Go extension work with standard Go for code navigation and autocompletion. The EmbeddedGo toolchain is only needed at build time.

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

Now that your toolchain works, try changing `ClearScreen(Blue)` to another color like `Red`, `Green`, or `DarkPurple` and rebuild. Then explore the [examples](https://github.com/drpaneas/gosprite64/tree/main/examples) in the GoSprite64 repository to learn about input handling, drawing shapes, text rendering, and audio.
