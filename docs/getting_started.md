# Getting Started

## Prerequisites

You need to have `Go` installed on your system.
And then to have `mage` installed as well.

e.g. `go install github.com/magefile/mage@latest` or `brew install mage`.

Also you need to have `$GOBIN`  to your `$PATH`, so you can run anything your install with `go install` from anywhere.

Your `$GOBIN` is usually `$GOPATH/bin`.

You can check if you have `$GOBIN` to your `$PATH` by running:

```bash
echo $PATH
```

If you don't have `$GOBIN` to your `$PATH`, you can add it by running:

```bash
export PATH=$PATH:$GOBIN
```

You can add this line to your `~/.zshrc` or `~/.bashrc` file, so it will be permanent.

### Example

1. Verify you have Go installed:

```bash
pi@raspberrypi:~ $ go version
go version go1.24.3 linux/arm64 # I have go1.24.3 installed
```

2. Now let's install `mage`:

```bash
pi@raspberrypi:~ $ go install github.com/magefile/mage@latest
```

If you try to run it, you might get an error, like this:

```bash
 $ mage --help
-bash: mage: command not found
```

This is because `$GOBIN` is not in your `$PATH`. So let's fix this.

3. Add `$GOBIN` to your `$PATH`:

```bash
pi@raspberrypi:~ $ cd $HOME/go # go where you $GOPATH is

pi@raspberrypi:~/go $ ls
bin  pkg  src # see these 3 folders, $GOBIN is the 'bin' one

pi@raspberrypi:~/go $ cd bin/
pi@raspberrypi:~/go/bin $ pwd
/home/pi/go/bin # this is your $GOBIN

pi@raspberrypi:~/go/bin $ export PATH=$PATH:$HOME/go/bin # you can put it your ~/.bashrc for permanent
```

4. Verify it works:

```bash
pi@raspberrypi:~ $ mage --help
mage [options] [target]

Mage is a make-like command runner.  See https://magefile.org for full docs.

Commands:
  -clean    clean out old generated binaries from CACHE_DIR
  -compile <string>
            output a static binary to the given path
  -h        show this help
  -init     create a starting template if no mage files exist
  -l        list mage targets in this directory
  -version  show version info for the mage binary

Options:
  -d <string> 
            directory to read magefiles from (default "." or "magefiles" if exists)
  -debug    turn on debug messages
  -f        force recreation of compiled magefile
  -goarch   sets the GOARCH for the binary created by -compile (default: current arch)
  -gocmd <string>
      use the given go binary to compile the output (default: "go")
  -goos     sets the GOOS for the binary created by -compile (default: current OS)
  -ldflags  sets the ldflags for the binary created by -compile (default: "")
  -h        show description of a target
  -keep     keep intermediate mage files around after running
  -t <string>
            timeout in duration parsable format (e.g. 5m30s)
  -v        show verbose output when running mage targets
  -w <string>
            working directory where magefiles will run (default -d value)
```

See? Now you can access `mage` and anything you install with `go install` from anywhere.

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

You can find the examples at `examples/` directory, you can load the rom `*.z64` with your favorite emulator (e.g. `ares`).

It has installed `direnv` and created a `.envrc` file in the root directory of the repository.

### Example

```bash
pi@raspberrypi:~/go/src/gosprite64 $ mage Setup
Using default GOPATH=/home/pi/go, GOBIN=/home/pi/go/bin
Running: go install github.com/clktmr/dl/gotip-embedded@latest
go: downloading github.com/clktmr/dl v0.0.0-20250603124022-78d0cf544a51
Running: gotip-embedded download
Cloning into '/home/pi/sdk/gotip-embedded'...
warning: redirecting to https://github.com/clktmr/go.git/
remote: Enumerating objects: 14266, done.
remote: Counting objects: 100% (14266/14266), done.
remote: Compressing objects: 100% (11372/11372), done.
remote: Total 14266 (delta 2485), reused 12563 (delta 2362), pack-reused 0 (from 0)
Receiving objects: 100% (14266/14266), 28.63 MiB | 10.94 MiB/s, done.
Resolving deltas: 100% (2485/2485), done.
Updating files: 100% (13112/13112), done.
Updating the go development tree...
warning: redirecting to https://github.com/clktmr/go.git/
From https://www.github.com/clktmr/go
 * branch            master-embedded -> FETCH_HEAD
HEAD is now at 83320e9 noos: Make use of free memory before the arena
Building Go cmd/dist using /usr/local/go. (go1.24.3 linux/arm64)
Building Go toolchain1 using /usr/local/go.
Building Go bootstrap cmd/go (go_bootstrap) using Go toolchain1.
Building Go toolchain2 using go_bootstrap and Go toolchain1.
Building Go toolchain3 using go_bootstrap and Go toolchain2.
Building packages and commands for linux/arm64.
---
Installed Go for linux/arm64 in /home/pi/sdk/gotip-embedded
Installed commands in /home/pi/sdk/gotip-embedded/bin
Success. You may now run 'gotip-embedded'!
Created symlink from /home/pi/go/bin/go1.22 to /home/pi/go/bin/gotip-embedded
Running: go install github.com/clktmr/n64/tools/mkrom
go: downloading github.com/clktmr/n64 v0.0.0-20250514144408-7410fc26cec0
Created .envrc file at /home/pi/go/src/gosprite64/.envrc (GOTOOLCHAIN=go1.22)
Installing direnv via apt-get
Running: sudo apt-get update
Get:1 http://security.debian.org/debian-security bullseye-security InRelease [27.2 kB]
Hit:2 http://deb.debian.org/debian bullseye InRelease                                                               
Get:3 http://deb.debian.org/debian bullseye-updates InRelease [44.1 kB]           
Get:4 http://deb.debian.org/debian bullseye-backports InRelease [49.0 kB]         
Get:5 http://archive.raspberrypi.org/debian bullseye InRelease [39.0 kB]          
Get:6 http://security.debian.org/debian-security bullseye-security/main arm64 Packages [372 kB]
Get:7 http://security.debian.org/debian-security bullseye-security/main armhf Packages [369 kB]
Get:8 http://security.debian.org/debian-security bullseye-security/main Translation-en [249 kB]
Get:9 http://archive.raspberrypi.org/debian bullseye/main arm64 Packages [323 kB]
Get:10 http://archive.raspberrypi.org/debian bullseye/main armhf Packages [330 kB]
Fetched 1802 kB in 2s (889 kB/s)                          
Reading package lists... Done
Running: sudo apt-get install -y direnv
Reading package lists... Done
Building dependency tree... Done
Reading state information... Done
The following package was automatically installed and is no longer required:
  raspinfo
Use 'sudo apt autoremove' to remove it.
The following NEW packages will be installed:
  direnv
0 upgraded, 1 newly installed, 0 to remove and 4 not upgraded.
Need to get 1816 kB of archives.
After this operation, 6357 kB of additional disk space will be used.
Get:1 http://deb.debian.org/debian bullseye/main arm64 direnv arm64 2.25.2-2 [1816 kB]
Fetched 1816 kB in 0s (16.9 MB/s)
perl: warning: Setting locale failed.
perl: warning: Please check that your locale settings:
 LANGUAGE = (unset),
 LC_ALL = (unset),
 LC_CTYPE = "UTF-8",
 LANG = "en_GB.UTF-8"
    are supported and installed on your system.
perl: warning: Falling back to a fallback locale ("en_GB.UTF-8").
locale: Cannot set LC_CTYPE to default locale: No such file or directory
locale: Cannot set LC_ALL to default locale: No such file or directory
Selecting previously unselected package direnv.
(Reading database ... 85673 files and directories currently installed.)
Preparing to unpack .../direnv_2.25.2-2_arm64.deb ...
Unpacking direnv (2.25.2-2) ...
Setting up direnv (2.25.2-2) ...
Processing triggers for man-db (2.9.4-2) ...
Added direnv hook to /home/pi/.bashrc
Running in /home/pi/go/src/gosprite64: direnv allow
Successfully ran 'direnv allow' for the .envrc file.

========== Setup Complete ==========

To build:      go build -o test.elf .
To create rom: mkrom test.elf


run this command: go env
AR='ar'
CC='gcc'
CGO_CFLAGS='-O2 -g'
CGO_CPPFLAGS=''
CGO_CXXFLAGS='-O2 -g'
CGO_ENABLED='1'
CGO_FFLAGS='-O2 -g'
CGO_LDFLAGS='-O2 -g'
CXX='g++'
GCCGO='gccgo'
GO111MODULE=''
GOARCH='arm64'
GOARM64='v8.0'
GOAUTH='netrc'
GOBIN=''
GOCACHE='/home/pi/.cache/go-build'
GOCACHEPROG=''
GODEBUG=''
GOENV='/home/pi/.config/go/env'
GOEXE=''
GOEXPERIMENT=''
GOFIPS140='off'
GOFLAGS=''
GOGCCFLAGS='-fPIC -pthread -Wl,--no-gc-sections -fmessage-length=0 -ffile-prefix-map=/tmp/go-build752222880=/tmp/go-build -gno-record-gcc-switches'
GOHOSTARCH='arm64'
GOHOSTOS='linux'
GOINSECURE=''
GOMOD='/home/pi/go/src/gosprite64/go.mod'
GOMODCACHE='/home/pi/go/pkg/mod'
GONOPROXY=''
GONOSUMDB=''
GOOS='linux'
GOPATH='/home/pi/go'
GOPRIVATE=''
GOPROXY='https://proxy.golang.org,direct'
GOROOT='/usr/local/go'
GOSUMDB='sum.golang.org'
GOTELEMETRY='local'
GOTELEMETRYDIR='/home/pi/.config/go/telemetry'
GOTMPDIR=''
GOTOOLCHAIN='auto'
GOTOOLDIR='/usr/local/go/pkg/tool/linux_arm64'
GOVCS=''
GOVERSION='go1.24.3'
GOWORK=''
PKG_CONFIG='pkg-config'
```

See it created this `.envrc` file:

```bash
pi@raspberrypi:~/go/src/gosprite64 $ cat .envrc 
export GOOS="noos"
export GOARCH="mips64"
export GOPATH="/home/pi/go"
export GOBIN="/home/pi/go/bin"
export GOFLAGS="-tags=n64 '-ldflags=-M=0x00000000:8M -F=0x00000400:8M -stripfn=1'"
export GOTOOLCHAIN="go1.22"
```

However, for some reason, it might be that `direnv` is not picking up the `.envrc` file.

```bash
pi@raspberrypi:~/go/src/gosprite64 $ direnv status
direnv exec path /usr/bin/direnv
DIRENV_CONFIG /home/pi/.config/direnv
bash_path /usr/bin/bash
disable_stdin false
warn_timeout 5s
whitelist.prefix []
whitelist.exact map[]
No .envrc loaded
Found RC path /home/pi/go/src/gosprite64/.envrc
Found watch: ".envrc" - 2025-06-15T21:41:01+02:00
Found watch: "../../../.local/share/direnv/allow/a1f181f044bc3e2d5b98e6257480b8d4742f7bbd0f2fc95a5b086205ab65acef" - 2025-06-15T21:35:12+02:00
Found RC allowed true
Found RC allowPath /home/pi/.local/share/direnv/allow/a1f181f044bc3e2d5b98e6257480b8d4742f7bbd0f2fc95a5b086205ab65acef
```

See: `No .envrc loaded`.

To fix that, do:

```bash
pi@raspberrypi:~/go/src/gosprite64 $ echo 'eval "$(direnv hook bash)"' >> ~/.bashrc
pi@raspberrypi:~/go/src/gosprite64 $ source ~/.bashrc
direnv: loading ~/go/src/gosprite64/.envrc
direnv: export +GOARCH +GOBIN +GOFLAGS +GOOS +GOPATH +GOTOOLCHAIN
```

So now it should work:

```bash
pi@raspberrypi:~/go/src/gosprite64 $ go version
go version devel go1.22-83320e9 Thu May 22 23:46:00 2025 +0200 linux/arm64
```

If you do not have bash, then you can use the one for your shell:

```bash
# for ZSH
echo 'eval "$(direnv hook zsh)"' >> ~/.zshrc 

# for Fish
echo 'eval (direnv hook fish)' >> ~/.config/fish/config.fish

# for Windows
Invoke-Expression "$(direnv hook powershell)"
Add-Content $PROFILE 'Invoke-Expression "$(direnv hook powershell)"'
```

Great, so now you can build the examples:

```bash
pi@raspberrypi:~/go/src/gosprite64 $ ./build_examples.sh 
Finding example directories...
Found example: clearscreen
Building example in clearscreen
  Running go build -o game.elf .
go: downloading github.com/embeddedgo/fs v0.1.3
go: downloading github.com/embeddedgo/display v1.1.0
go: downloading github.com/sigurn/crc8 v0.0.0-20220107193325-2243fe600f9f
go: downloading golang.org/x/text v0.13.0
  Running mkrom game.elf
  Successfully built clearscreen
Found example: pong
Building example in pong
  Running go build -o game.elf .
  Running mkrom game.elf
  Successfully built pong
Found example: space_invaders
Building example in space_invaders
  Running go build -o game.elf .
  Running mkrom game.elf
  Successfully built space_invaders
All examples built successfully!
```

And indeed, here are your N64 roms:

```bash
pi@raspberrypi:~/go/src/gosprite64 $ find . -iname "*.z64"
./examples/clearscreen/game.z64
./examples/space_invaders/game.z64
./examples/pong/game.z64
```

## Create your own project

1. Create a new directory for your project:

```bash

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
