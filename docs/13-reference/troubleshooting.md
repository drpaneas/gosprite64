# Troubleshooting

Common build errors, runtime issues, and emulator quirks with solutions.

## Beginner setup and first-run problems

If you are blocked before the first ROM works:

- go back to [Installation](../02-getting-started/installation.md)
- confirm you followed [Run Your First ROM](../02-first-journey/01-run-your-first-rom.md)
- rebuild before reopening the ROM in your emulator

## Build Errors

### "ambiguous import" or multiple module errors

**Symptom:** `go build` reports an "ambiguous import" error, or complains about a package being provided by multiple modules.

**Cause:** There is a stray `go.mod` file inside one of the `examples/` subdirectories. Go sees it as a separate module that conflicts with the root module.

**Fix:** Make sure you are building from the repository root and that no example directory has its own `go.mod`. Each example should be part of the root module. If you find an extra `go.mod` in an example folder, delete it:

```bash
rm examples/mygame/go.mod
rm examples/mygame/go.sum
```

Then run `go mod tidy` from the repository root.

### "package X is not in std"

**Symptom:** The compiler reports that packages like `embedded/rtos` or `embedded/arch/r4000/systim` are not in the standard library.

**Cause:** You are using the standard `go build` path instead of the supported N64 cross-compilation workflow. GoSprite64 targets the N64 via the EmbeddedGo toolchain and the settings loaded from `n64.env`.

**Fix:** Build with `GOENV=n64.env` and the EmbeddedGo binary:

```bash
GOENV=n64.env go1.24.5-embedded build -o examples/mygame/game.elf ./examples/mygame
GOENV=n64.env n64go rom examples/mygame/game.elf
```

`GOENV=n64.env` tells the Go toolchain to load the tracked N64 settings from `n64.env`, including `GOTOOLCHAIN`, `GOOS`, `GOARCH`, and `GOFLAGS`. Without that configuration, Go tries to resolve N64-specific packages against your host toolchain.

### Build fails with embedded/* errors

**Symptom:** Compilation fails with errors referencing `embedded/rtos`, `embedded/arch/...`, or other packages under the `embedded/` prefix.

**Cause:** You need the `go1.24.5-embedded` toolchain and the `n64.env` workflow. The standard Go toolchain does not include the embedded runtime packages.

**Fix:**

1. Install the N64-capable Go toolchain. Follow the installation guide in [Installation](../02-getting-started/installation.md).
2. Make sure `n64.env` exists in your project root, then build with it explicitly:

```bash
GOENV=n64.env go1.24.5-embedded build -o examples/mygame/game.elf ./examples/mygame
GOENV=n64.env n64go rom examples/mygame/game.elf
```

3. Verify your Go version:

```bash
go1.24.5-embedded version
```

It should report `go1.24.5-embedded`. If that command is missing or fails, reinstall the EmbeddedGo toolchain before trying the build again.

### `go mod tidy` fails with network or version errors

**Symptom:** Running `go mod tidy` produces errors about missing versions or fails to fetch dependencies.

**Cause:** The N64 environment loaded through `n64.env` can confuse `go mod tidy` because it tries to resolve dependencies for the embedded target platform instead of your host environment.

**Fix:** Run `go mod tidy` with the same N64-specific environment cleared:

```bash
env -u GOENV -u GOOS -u GOARCH -u GOFLAGS -u GOTOOLCHAIN go mod tidy
```

This tells Go to use your host platform for dependency resolution while keeping the rest of your environment intact. After tidy completes, return to the normal build flow with `GOENV=n64.env`.

## Runtime Issues

### Black screen (nothing renders)

**Symptom:** The ROM boots but the screen stays black. No tiles, sprites, or text appear.

**Cause:** You forgot to call `RegisterAssetFS` before `Run()`. Without it, the cartridge filesystem is not mounted and all asset loads silently fail.

**Fix:** In your `main.go`, register the embedded filesystem before starting the game loop:

```go
package main

import (
    "embed"

    "github.com/clktmr/n64/drivers/cartfs"
    "github.com/drpaneas/gosprite64"
)

//go:embed assets/*
var embeddedAssets embed.FS

var assetFS = cartfs.Embed(embeddedAssets)

func main() {
    gosprite64.RegisterAssetFS(assetFS)
    gosprite64.Run(&MyGame{})
}
```

Make sure the `go:embed` directive matches the directory where your compiled assets live. If you renamed or moved your assets folder, update the embed path accordingly.

### No audio (sound effects and music are silent)

**Symptom:** The game runs and renders correctly, but `PlaySoundEffect` and `PlayMusic` do nothing. No sound is heard.

**Cause:** The audio bundle is not registered. GoSprite64 requires a generated `audio_embed.go` file that calls `RegisterAudioBundle` with your compiled audio assets.

**Fix:**

1. Make sure you have run the `audiogen` tool to compile your audio assets:

```bash
go run ./cmd/audiogen -manifest audio/manifest.json -out audio_embed.go
```

2. Verify that `audio_embed.go` exists in your game directory and contains a call to `gosprite64.RegisterAudioBundle(...)`.

3. The `RegisterAudioBundle` call must execute before `Run()`. The generated file typically uses an `init()` function, so simply having the file in your package is enough.

If you still hear no audio, check that your emulator supports audio output (ares, simple64, and cen64 all support it; some older emulators may not).

### Scene only fills part of the screen

**Symptom:** The tile scene renders but only covers a portion of the screen, leaving black bars or empty space.

**Cause:** The camera's `Width` and `Height` do not match the logical canvas size of 288x216.

**Fix:** Set the camera dimensions to match the canvas:

```go
cam := &gosprite64.Camera{
    Width:  288,
    Height: 216,
}
```

If you are using `Scene.Draw(nil)`, the scene creates a default camera with the correct dimensions automatically. But if you create your own camera and forget to set the size, it defaults to zero, which means no tiles fall within the viewport.

### D-pad does not respond

**Symptom:** `IsButtonDown(ButtonDPadUp)` and similar calls always return false, even though the D-pad works in other games.

**Cause:** This is usually an emulator input mapping issue. Some emulators map keyboard arrow keys to the analog stick by default rather than the D-pad.

**Fix (ares):**

1. Go to **Settings > Input**.
2. Find the D-pad mappings (DPad Up, DPad Down, DPad Left, DPad Right).
3. Bind them to your preferred keyboard keys.

**Fix (simple64):**

1. Go to **Options > Configure Controller**.
2. Verify that the D-pad entries are mapped. By default, simple64 may only map the analog stick.

If you want your game to respond to both the D-pad and analog stick, check both in your update logic:

```go
func (g *MyGame) Update() {
    sx, sy := gosprite64.StickPosition(0.2)
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft) || sx < -0.5 {
        // move left
    }
}
```

## Emulator Notes

### Which emulators work?

GoSprite64 ROMs are tested primarily with:

- **ares** - Recommended. Accurate N64 emulation with good audio and input support.
- **simple64** - Good alternative with a user-friendly interface.
- **cen64** - Cycle-accurate but slower. Useful for verifying hardware accuracy.

Project64 and Mupen64Plus may also work but are not regularly tested.

### ROM does not boot in the emulator

If the emulator shows an error or the ROM does not start:

1. Make sure you built a `.z64` ROM file (big-endian format), not just a `.elf` file. The ELF is an intermediate output; the supported conversion step is:

```bash
GOENV=n64.env n64go rom examples/mygame/game.elf
```
2. Check that the ROM header is correct. Some emulators are strict about the header format.
3. Try a different emulator to narrow down whether it is a ROM issue or an emulator compatibility issue.
