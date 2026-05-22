# Working With PixelLab

This example shows a PixelLab-generated sprite moving through the normal
GoSprite64 asset pipeline. PixelLab is used once to create a source PNG, then
the repository converts that checked-in image into local runtime assets that the
game loads from the cartridge filesystem.

## What the example demonstrates

The `examples/pixellab_sprite` example keeps the workflow intentionally small:

- a single checked-in source sprite generated with PixelLab
- compile-time conversion from PNG to GoSprite64 sprite-sheet data
- runtime loading of the compiled asset with no PixelLab dependency

## PixelLab source asset

The checked-in PixelLab output is:

- `examples/pixellab_sprite/assets-src/slime_right_32.png`

That PNG is the generation-time source asset for the example. The accompanying
test in `examples/pixellab_sprite/assets-src/source_asset_test.py` verifies that
the selected image exists, stays at `32x32`, and keeps transparency needed for
masked sprite compositing.

## How the source asset becomes runtime files

The example's `go:generate` directive in `examples/pixellab_sprite/generate.go`
converts the checked-in PNG into a GoSprite64 runtime sprite sheet:

```bash
go run github.com/drpaneas/gosprite64/cmd/mk2dsheet \
  -in examples/pixellab_sprite/assets-src/slime_right_32.png \
  -out examples/pixellab_sprite/assets/hero.sheet \
  -tile-width 32 \
  -tile-height 32
```

That produces:

- `examples/pixellab_sprite/assets/hero.sheet`

On N64 builds, `examples/pixellab_sprite/assets_embed.go` embeds `assets/*`
into the cartridge filesystem, and `examples/pixellab_sprite/runtime_n64.go`
loads `assets/hero.sheet` at runtime with `gosprite64.LoadSpriteSheet`.

## Rebuild the assets and ROM locally

From the repository root, regenerate the compiled sprite asset:

```bash
go generate ./examples/pixellab_sprite
```

Then rebuild the example ROM directly:

```bash
GOENV="$PWD/n64.env" go1.24.5-embedded build \
  -o examples/pixellab_sprite/game.elf \
  ./examples/pixellab_sprite

GOENV="$PWD/n64.env" n64go rom examples/pixellab_sprite/game.elf
```

If you want to rebuild every example in one pass, the repository helper also
works:

```bash
./build_examples.sh
```

The example ROM produced by the direct build lives at:

- `examples/pixellab_sprite/game.z64`

For the docs-hosted copy, see:

> **Download the ROM:** [`pixellab_sprite.z64`](../emulator/roms/pixellab_sprite.z64)

## Runtime uses local compiled assets only

PixelLab is generation-time only in this workflow. The runtime does not call
PixelLab, download assets, or depend on any external service. Once the PNG is
checked in and converted to `assets/hero.sheet`, the game runs entirely from
local compiled assets embedded into the ROM.
