# Environment setup

Typically, homebrew retro-gamedev development with C/C++ is a bit messy to setup and build. Even though Go has very easy way of handling things, we're still building **actual Nintendo 64 games** here. But don’t worry — we’ve automated *everything*. There are no Makefiles, no cmake and no crying. Also, whether you're on **Windows**, **macOS**, or **Linux**, the setup is exactly the same.

All you need is actually to install some modern version of [Go] and [Mage], and you'll be ready to build games that run on real N64 hardware.

---

## 🧰 What you'll need

Before you dive in, you’ll only need two things installed:

1. ✅ **Go** (any modern version will do) — [Install Go](https://go.dev/doc/install)  
2. ✅ **Mage** – a task runner we use for setup automation:  

   ```sh
   go install github.com/magefile/mage@latest
   ```

## 🧙‍♂️ One command to rule them all

Once you have Go and Mage installed, run this magic spell:

```sh
git clone https://github.com/drpaneas/gosprite64
cd gosprite64
mage Setup
```

Sit back. Grab a coffee. Our Mage is taking over. ☕

He’ll:

* 🛠 Fetch a [custom version] of Go built for `MIPS` (what the N64 uses)
* 🔧 Build it locally
* 📦 Install [emgo], the embedded Go tool
* 📁 Set up a clean environment in `~/toolchains/nintendo64` for Mac/Linux or `%USERPROFILE%\toolchains\nintendo64` on Windows.
* ⚙ Configure [direnv] to manage your previously built Go environment automatically

No clutter in `.bashrc`, `.zshrc`, or `.profile`. It *just* works.

## 📁 Where things live

After setup, your new Go environment will be isolated in:

* `GOROOT` → `~/toolchains/nintendo64/go`
* `GOPATH` → `~/toolchains/nintendo64/gopath`

From now on, all your N64 projects should live inside `~/toolchains/nintendo64` (or open a shell inside it). Thanks to [direnv], your terminal will **automatically** switch to the embedded Go version whenever you `cd` into this folder. Leave the folder? You're back to your *normal* Go setup.
No switching. No weird aliases. No surprises.

## 🧭 Platform-specific guidance

If you really want to peek behind the curtain or you're running into edge cases, we have detailed setup pages for:

* [macOS setup](mac.md)
* [Windows setup](windows.md)
* [GNU/Linux setup](linux.md)

But honestly? *Just* run `mage Setup`. It works.

[Go]: https://go.dev/
[Mage]: https://magefile.org/
[emgo]: https://github.com/embeddedgo/tools/tree/master/emgo
[custom version]: https://github.com/clktmr/go
[direnv]: https://github.com/direnv/direnv
