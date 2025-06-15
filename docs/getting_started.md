# Getting Started

## Installation

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
You can find the examples at `examples/` directory, load the rom `*.z64` with your favorite emulator (e.g. `ares`).

It has installed `direnv` and created a `.envrc` file in the root directory of the repository.
You need to run `direnv allow` in the root directory of the repository.

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

2. Run `direnv allow` in the root directory of your project.

```bash
direnv allow
```

this will use the Go Nintendo everytime you use `go` in this directory.

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
go build -o mygame .
```

5. Create the ROM file:

```bash
mkrom mygame
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
