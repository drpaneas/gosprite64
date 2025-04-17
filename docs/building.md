# 🎮 Building Your First GoSprite64 Example

So you’ve set up your environment? Sweet.  
Now it’s time to **build your first actual Nintendo 64 ROM** using GoSprite64 and watch Bitron come to life. ⚡

And yes; it’s as easy as running **one command**:

```sh
mage Test
```

🧙 What does mage Test do?

Here's what happens under the hood:

1. Clones the GoSprite64 source repo (if it's not already cloned)
2. Navigates to examples/clearscreen
3. Initializes a Go module with emgo mod init
4. Replaces the Go module import path with your local GoSprite64 path
5. Runs emgo mod tidy to fetch dependencies
6. Builds the example with emgo build
7. Confirms a .z64 ROM file was created

If everything works, you’ll see something like:

```sh
Build succeeded! Generated .z64 file(s): [/Users/drpaneas/toolchains/nintendo64/gopath/src/gosprite64/examples/clearscreen/clearscreen.z64]

Please load the .z64 file into the Ares emulator.
```

Boom 💥! You've got yourself a playable N64 ROM!

So if you go to `cd ~/toolchains/nintendo64/gopath/src/gosprite64` you will find our repo again, cloned over there. Actually we do not need the whole repo, just the `examples` folder would do, but this makes our life easier in case you want to contribute to the project later.

Why? Because every example is using a `replace` directory, where instead of using the GoSprite repo from GitHub, it uses the one locally at your toolchain. In other words, the example Clearscreen: `~/toolchains/nintendo64/gopath/src/gosprite64/examples/clearscreen` uses `~/toolchains/nintendo64/gopath/src/gosprite64`. That means you can easily make any changes to GoSprite64 and test them locally.

### 🛠 Want to tinker?

Crack open the `main.go` inside the example and start playing around:

```sh
cd ~/toolchains/nintendo64/gopath/src/gosprite64/examples/clearscreen
nvim main.go  # or your favorite editor
```

```go
package main

import (
        "image/color"

        gospr64 "github.com/drpaneas/gosprite64"
)

var Azure = color.RGBA{0xf0, 0xff, 0xff, 0xff} // rgb(240, 255, 255)

type Game struct {
}

func (g *Game) Update() error {
        return nil
}

func (g *Game) Draw(screen *gospr64.Screen) {
        screen.Clear(Azure)
}

func main() {
        // Initialize the game
        game := &Game{}

        // Run the game
        if err := gospr64.Run(game); err != nil {
                panic(err)
        }

}
```

Then make sure you have a `build.cfg` file next to it:

### ⚙️ `build.cfg` — The Memory Blueprint

In every **GoSprite64** project, you’ll need a file named **`build.cfg`** sitting next to your `main.go`.

Why? Because this file tells the toolchain **how to lay out memory** for your N64 ROM.  
Think of it as the **blueprint** for where code and data live in the final game binary.

Without it, your build may fail—or worse, the ROM may not boot at all on real hardware.

#### 🧠 What’s inside `build.cfg`?

Here’s an example:

```ini
GOTARGET = n64
GOMEM = 0x00000000:4M
GOTEXT = 0x00000400:4M
GOSTRIPFN = 0
GOOUT = z64
```

> Note: Most emulators emulate the Expansion Pak by default, but real hardware will need the actual RAM module installed.

If you do have an Expansion Pak with your Nintendo64, you can make sure of aditional 8M memory, required for larger assets. To use 8MB:

```ini
GOMEM = 0x00000000:8M
GOTEXT = 0x00000400:8M
```

#### 📁 Where does build.cfg go?

It must be in the same folder as your main.go and where you run emgo build.

```sh
your-project/
├── build.cfg
├── main.go
└── ...
```

Then rebuild:

```sh
emgo build
```

Your updated .z64 file will be ready for action in seconds.
