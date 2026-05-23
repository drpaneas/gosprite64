# Working With Aseprite

This example shows the Aseprite-to-GoSprite64 workflow used in the repository:
draw a sprite in Aseprite, export it to PNG, convert that PNG into a `.sheet`
with `mk2dsheet`, then build and run the ROM.

## What the example includes

The example lives in `examples/aseprite_sprite` and includes these files:

- `examples/aseprite_sprite/assets-src/hero.aseprite` - the original Aseprite source
- `examples/aseprite_sprite/assets-src/hero.png` - the exported transparent PNG
- `examples/aseprite_sprite/assets/hero.sheet` - the generated sprite sheet binary
- `examples/aseprite_sprite/assets_embed.go` - embeds the generated assets into the ROM
- `examples/aseprite_sprite/main.go` - loads and draws the sprite at runtime
- `examples/aseprite_sprite/game.z64` - the locally built ROM artifact

This revised example keeps the hero art at `64x64`, but the runtime sheet is
split into four `32x32` tiles. The Aseprite source is `64x64`, the exported PNG
is `64x64`, and `mk2dsheet` slices that PNG into four runtime frames so the
example fits the N64 `RGBA32` tile limits.

## Aseprite to PNG

Start from the source file:

```text
examples/aseprite_sprite/assets-src/hero.aseprite
```

Export that sprite from Aseprite as a PNG with transparency preserved:

```text
examples/aseprite_sprite/assets-src/hero.png
```

For this example the source art is a single-frame `64x64` Aseprite file, and the
committed PNG export is the same `64x64` transparent frame.

## PNG to .sheet

GoSprite64 loads compiled `.sheet` assets at runtime, so the exported PNG is
converted with `mk2dsheet`.

The example keeps that step in a `go:generate` directive:

```go
//go:generate sh -c "mkdir -p assets && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/hero.png -out assets/hero.sheet -tile-width 32 -tile-height 32"
```

Run it with:

```bash
go generate ./examples/aseprite_sprite
```

That regenerates `examples/aseprite_sprite/assets/hero.sheet` from
`examples/aseprite_sprite/assets-src/hero.png` using `32x32` tile dimensions.
The source art still stays `64x64`; only the runtime sheet is tiled.

## Building the ROM

After generating the sheet, build the example ROM:

```bash
GOENV=n64.env go1.24.5-embedded build -o examples/aseprite_sprite/game.elf ./examples/aseprite_sprite
GOENV=n64.env n64go rom examples/aseprite_sprite/game.elf
```

This produces `examples/aseprite_sprite/game.z64`.

## Loading and drawing the sprite

At runtime the example loads the compiled sheet, centers the overall hero as a
`64x64` sprite, and draws four `32x32` frames with `BlendMasked` so transparent
pixels stay transparent:

```go
func (g *Game) Init() {
    hero, err := gosprite64.LoadSpriteSheet("assets/hero.sheet")
    if err != nil {
        panic(err)
    }

    g.hero = hero
    g.x, g.y = centeredSpritePosition(screenWidth, screenHeight, heroWidth, heroHeight)
}

func (g *Game) Draw() {
    gosprite64.ClearScreenWith(gosprite64.DarkBlue)
    gosprite64.DrawText("Aseprite sprite loaded", 56, 24, gosprite64.White)
    for _, tile := range heroCompositeTiles(g.x, g.y) {
        gosprite64.DrawSpriteWithOptions(g.hero, tile.frame, float32(tile.x), float32(tile.y), gosprite64.DrawSpriteOptions{
            Blend: gosprite64.BlendMasked,
        })
    }
}
```

That is the whole workflow:

1. Edit `examples/aseprite_sprite/assets-src/hero.aseprite`
2. Export `examples/aseprite_sprite/assets-src/hero.png`
3. Run `go generate ./examples/aseprite_sprite`
4. Build `examples/aseprite_sprite/game.elf` and `examples/aseprite_sprite/game.z64`
5. Run the ROM and confirm the centered `64x64` hero renders correctly from four `32x32` runtime tiles

## Try It

> **Download the ROM:** [`aseprite_sprite.z64`](../emulator/roms/aseprite_sprite.z64) - Open in [ares](https://ares-emu.net/) with the Expansion Pak enabled.

## Reference Example

See `examples/aseprite_sprite` for the complete source used by this workflow.
