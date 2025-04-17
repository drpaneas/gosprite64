# 🐧 Linux Setup

If you're here, you're probably running something like Ubuntu, Arch, Fedora, or maybe even Gentoo on a thinkpad with stickers. Either way: **respects**.

You’re in a great spot to build games for the **Nintendo 64** using GoSprite64.

> Most of the setup is automated via `mage Setup`, but if you want to understand or troubleshoot the magic, read on.

---

## 🧱 Prerequisites

You’ll need just a few things to get rolling:

### ✅ Install [Go](https://go.dev/doc/install)

Your distro probably has a package for it:

```sh
# Debian/Ubuntu
sudo apt install golang

# Fedora
sudo dnf install golang

# Arch
sudo pacman -S go
```

Or download and install it manually from [go.dev](https://go.dev/doc/install).

Make sure you can use it by giving it as try, *run*: `go version`.

### ✅ Install [Mage]

```sh
go install github.com/magefile/mage@latest
```

Verify you can use it: `which mage`.

### ✅ Install [direnv]

This allows your terminal to automatically switch Go environments when entering `~/toolchains/nintendo64`.

```sh
# Debian/Ubuntu
sudo apt install direnv

# Fedora
sudo dnf install direnv

# Arch
sudo pacman -S direnv
```

After that, hook `direnv` into your shell. Edit either `~/.zshrc`or  `~/.bashrc` and add one of the following lines:

```sh
eval "$(direnv hook zsh)"   # for ZSH
eval "$(direnv hook bash)"  # for bash
```

Then run either `source ~/.bashrc` or `source ~/.zshrc` respectively.
After this, direnv will automatically load the N64 environment when you're inside `~/toolchains/nintendo64`.

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
