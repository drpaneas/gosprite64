# 🍎 macOS Setup

If you’ve landed on this page, it means you’re using macOS 🖥️ to build legendary Nintendo 64 games with **GoSprite64**.   Good news: macOS is *great* for this. Even better news? You probably don’t need to do much here.

> 👉 All of the setup is fully 100% automated via `mage Setup`. But if you want to know what's happening under the hood, read on.

---

## 🧱 Prerequisites

Here’s what you’ll need before you can start hacking around:

### Brew

First, make sure you install Apple's developer resources, that is Xcode, including the Command Line tools.

```sh
xcode-select --install
xcode-select -p # to verify they are installed
```

```sh
/Library/Developer/CommandLineTools
```

You will also need to install several other packages for which we'll include instructions assuming you have installed the Homebrew package manager on your system:

First install the package manager, that is [Homebrew](https://brew.sh/):

```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

Output:

```bash
/Library/Developer/CommandLineTools
```

### ✅ Install [Go](https://go.dev/doc/install)

If you don’t already have Go installed:

```sh
brew install go
```

Make sure you can use it by giving it as try, *run*: `go version`.

### ✅ Install [Mage]

```sh
go install github.com/magefile/mage@latest
```

Make sure you can use it `which mage`.

## ✅ Install [direnv]

```sh
brew install direnv
```

After that, hook `direnv` into your shell. Edit either `~/.zshrc`or  `~/.bashrc` and add one of the following lines:

```sh
eval "$(direnv hook zsh)"   # for ZSH
eval "$(direnv hook bash)"  # for bash
```

Then run either `source ~/.bashrc` or `source ~/.zshrc` respectively.

## 🧙‍♂️ Run the Setup

Now clone the repo and run:

```sh
git clone https://github.com/drpaneas/gosprite64
cd gosprite64
mage Setup
```

What happens now:

* 🧬 A [custom Go toolchain] for MIPS gets cloned and compiled
* 🛠 A `GOPATH` and `GOROOT` get set up at `~/toolchains/nintendo64/gopath` and `~/toolchains/nintendo64/go` respectively.
* 🧪 [emgo] tool gets installed
* ⚙ An `.envrc` file is created at `~/toolchains/nintendo64` so that [direnv] can auto-activate your dev environment

You can verify all those by going to `cd ~/toolchains/nintendo64/` and witness your Go env setup changing.

Outside of the toolchain's directory:

```sh
$ cd; go env GOROOT

# Output:
/opt/homebrew/Cellar/go/1.23.3/libexec # my system uses Go 1.23
```

```sh
$ cd; emgo

# Output:
zsh: command not found: emgo # expected, not from my std GOBIN
```

Inside of the toolchain's directory:

```sh
cd ~/toolchains/nintendo64 
```

```sh
# Output:
direnv: loading ~/toolchains/nintendo64/.envrc
direnv: export +GOBIN +GOROOT ~GOPATH ~PATH 


go env GOROOT

# Output
/Users/pgeorgia/toolchains/nintendo64/go

emgo

# Output
Go is a tool for managing Go source code.
```

[direnv]: https://github.com/direnv/direnv
[custom Go toolchain]: https://github.com/clktmr/go
[emgo]: https://github.com/embeddedgo/tools/tree/master/emgo
[Mage]: https://magefile.org/
