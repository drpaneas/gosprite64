# 🪟 Windows Setup

Welcome, traveler of the Windows realms! 💻

Yes, it’s true. You can build **real Nintendo 64 games** using **Go**, and you can do it on Windows too.  
Thanks to some automation magic, GoSprite64 makes this painless—even in a land of `.exe`, `PATH`, and `PowerShell`.

> Most of the work is handled by `mage Setup`. But if you're curious (or things go weird), read on.

---

## 🧱 Prerequisites

Here’s what you’ll need:

### ✅ Install [Go](https://go.dev/doc/install)

Download and install the official MSI package from [go.dev](https://go.dev/dl/).

During install, make sure to check the box to **add Go to your PATH**.

To verify it's working:

```powershell
go version
```

### ✅ Install Mage

In your terminal or PowerShell:

```sh
go install github.com/magefile/mage@latest
```

Make sure that Go's bin directory is in your PATH. Usually it’s at:

```sh
%USERPROFILE%\go\bin
```

You can add it permanently to your system environment variables or temporarily:

```sh
$env:Path += ";$env:USERPROFILE\go\bin"
```

### ✅ direnv (Yes, even on Windows!)

GoSprite64 uses direnv to manage the custom Go environment.

Go to the [releases] page, and download the latest `direnv.windows-amd64.exe` asset. Place it into `C:\Users\YourName\toolchains\nintendo64`.

🧙‍♂️ Run the Setup

Now clone the GoSprite64 repo and run:

```sh
git clone https://github.com/drpaneas/gosprite64
cd gosprite64
mage Setup
```

This will:

* Clone and build a custom version of Go for MIPS
* Install the emgo toolchain
* Configure your system with a new Go environment at:
  * `C:\Users\YourName\toolchains\nintendo64\go`
  * `C:\Users\YourName\toolchains\nintendo64\gopath`
* Download a Windows-compatible `direnv.exe`
* Set up .envrc to manage everything

All without touching your default Go installation. ✨

### 📦 How to use it

To develop N64 games using GoSprite64, always open a terminal in:

```sh
cd %USERPROFILE%\toolchains\nintendo64
```

[releases]: https://github.com/direnv/direnv/releases
