# Getting Started

## Prerequisites

You need to have `Go` installed on your system.

## Installation

0. Install `mage`: `go install github.com/magefile/mage@latest` or `brew install mage`.

1. Clone the repository:

```bash
cd $GOPATH/src
git clone https://github.com/drpaneas/gosprite64.git
cd gosprite64
```

2. Run the setup script:

```bash
mage Setup
./build_examples.sh
```

This is all you need.

You can find the examples at `examples/` directory, you can load the rom `*.z64` with your favorite emulator (e.g. `ares`).

It has installed `direnv` and created a `.envrc` file in the root directory of the repository.

## Create your own project

1. Create a new directory for your project:

```bash
cd $GOPATH/src
mkdir myproject
cd myproject
```

Create a `.envrc` file in the root directory of your project.

```bash
cp $GOPATH/src/gosprite64/.envrc .
```

e.g.:

```bash
drpaneas@m2:~/gocode/src/github.com/drpaneas/gosprite64 (main)% cat .envrc 
export GOOS="noos"
export GOARCH="mips64"
export GOFLAGS="-tags=n64 '-ldflags=-M=0x00000000:8M -F=0x00000400:8M -stripfn=1'"
export GOTOOLCHAIN="go1.22"
```

2. Run `direnv allow` in the root directory of your project only the very first time.

```bash
direnv allow
```

this will use the Go64 (from clktmr) everytime you use `go` in this directory.
As soon as you leave this directory, your system will use your default Go version, instead of Go64.
As soon as you `cd` back in the directory, it will use Go64 again.

3. Create a new main.go file:

```go
package main

import (
 . "github.com/drpaneas/gosprite64"
)

// Game instances to store game state
type Game struct{}

// Init is called once at the start of the game
func (g *Game) Init() {}

// Update game logic here
func (g *Game) Update() {}

// Draw game here
func (g *Game) Draw() {
 ClearScreen(Red)
}

func main() {
 Run(&Game{})
}
```

4. Create the ELF file:

```bash
go build -o mygame.elf .
```

This will create a `mygame.elf` file in the current directory.

5. Create the ROM file:

```bash
mkrom mygame.elf
```

6. Load the ROM file with your favorite emulator (e.g. `ares`).

## Game with audio

If you want to use audio, your files must be of type `raw` and have the following naming convention:

- Music files: `music0.raw`, `music1.raw`, ..., `music63.raw`

Then, at your `main.go` file, at the very top, add this line:

```bash
//go:generate go run github.com/drpaneas/gosprite64/cmd/audiogen -dir .
```

and then run `go generate ./...`.
This will download and install `audiogen` tool, which will scan for these audio raw files
and create an `audio_embed.go` file to embed them in your game automatically.

You can then use the audio files in your game using `Music(0, false)` or `Music(0, true)`, instead of `0` put your `music$0.raw` number
and selecte `true` or `false` depending if you want loop or not.

So far you can not draw sprites yet, only some basic primitives, such as `DrawRect`, `DrawRectFill`, `Line`, `Print` and `ClearScreen`.
