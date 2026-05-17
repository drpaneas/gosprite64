# Tile Sheets and Maps

GoSprite64 uses two types of source assets for tile scenes: **PNG tilesheets** (the graphics) and **JSON map files** (the layout). Build-time tools compile these into compact binary formats that the runtime loads efficiently. This page covers the asset formats, the compilation tools, and the key configuration fields. For the full build-to-render pipeline, see [Pipeline Overview](./pipeline-overview.md).

## Source asset overview

| Asset | Source format | Tool | Output | Binary magic |
|-------|--------------|------|--------|-------------|
| Tilesheet | PNG image | `mk2dsheet` | `.sheet` | `SHT2` |
| Tile map | JSON file | `mk2dmap` | `.map` | `MAP2` |

## PNG tilesheets

A tilesheet is a PNG atlas where every tile has the same dimensions. The image is sliced into a grid of equal-sized cells:

```
 ┌────┬────┬────┬────┐
 │ 0  │ 1  │ 2  │ 3  │    tiles.png (32x16, tile size 8x8)
 ├────┼────┼────┼────┤    -> 8 tiles total
 │ 4  │ 5  │ 6  │ 7  │
 └────┴────┴────┴────┘
```

Requirements:

- Image width must be evenly divisible by `tile-width`
- Image height must be evenly divisible by `tile-height`
- Maximum tile count is 65535 (uint16)
- Pixels are stored internally as NRGBA (pre-multiplied alpha is not used)

### Compiling with mk2dsheet

The `mk2dsheet` tool converts a PNG tilesheet into the binary `.sheet` format:

```bash
go run github.com/drpaneas/gosprite64/cmd/mk2dsheet \
  -in assets-src/tiles.png \
  -out assets/tiles.sheet \
  -tile-width 8 \
  -tile-height 8
```

Flags:

| Flag | Default | Description |
|------|---------|-------------|
| `-in` | (required) | Input PNG path |
| `-out` | (required) | Output `.sheet` path |
| `-tile-width` | `8` | Tile width in pixels |
| `-tile-height` | `8` | Tile height in pixels |

The output binary encodes tile dimensions, tile count, palette entry count, image dimensions, and the raw NRGBA pixel data.

## JSON map files

A map file describes the tile layout as a grid of cell indices that reference tiles in one or more tilesheets. Each map has one or more **layers**, and each layer references a tilesheet by `sheet_id`.

### Map JSON structure

```json
{
  "width": 32,
  "height": 18,
  "layer_count": 2,
  "cell_bits": 16,
  "chunk_width": 8,
  "chunk_height": 8,
  "layers": [
    {"sheet_id": 1, "cells": [1, 2, 3, 0, 0, ...]},
    {"sheet_id": 2, "cells": [0, 0, 1, 4, 0, ...]}
  ]
}
```

### Field reference

#### Top-level fields

| Field | Type | Description |
|-------|------|-------------|
| `width` | uint16 | Map width in tiles (must be > 0) |
| `height` | uint16 | Map height in tiles (must be > 0) |
| `layer_count` | uint16 | Number of tile layers (must be > 0, must match length of `layers` array) |
| `cell_bits` | uint8 | Bits per cell index: `8` (max 255 tiles) or `16` (max 65535 tiles) |
| `chunk_width` | uint16 | Chunk width in tiles for streaming/culling |
| `chunk_height` | uint16 | Chunk height in tiles for streaming/culling |
| `layers` | array | Per-layer tile data |

#### Layer fields

| Field | Type | Description |
|-------|------|-------------|
| `sheet_id` | uint16 | Which tilesheet this layer uses (1-based; defaults to 1 if omitted) |
| `cells` | []uint16 | Flat array of tile indices, length must equal `width * height` |

### cell_bits

The `cell_bits` field controls how many bits each tile index occupies in the compiled binary:

- **8** - Each cell is a single byte. Supports up to 255 unique tiles per layer. Use this for small tilesets to save memory.
- **16** - Each cell is two bytes. Supports up to 65535 unique tiles per layer. Use this for larger tilesets or when you need more than 255 distinct tiles.

A cell value of `0` means "empty" (no tile drawn at this position).

### Chunk dimensions

`chunk_width` and `chunk_height` define the size of streaming/culling chunks. The renderer uses chunks to skip drawing regions of the map that are off-screen. Typical values are 8x8 or 16x16 tiles. Smaller chunks give finer culling granularity at the cost of more chunk metadata.

### sheet_id

Each layer references a tilesheet by its 1-based `sheet_id`. When a bundle contains multiple tilesheets, different layers can draw from different sheets. For example, a background layer might use a terrain sheet (sheet_id 1) while a foreground layer uses a decoration sheet (sheet_id 2).

If `sheet_id` is omitted or set to 0, it defaults to 1.

### Compiling with mk2dmap

The `mk2dmap` tool converts a JSON map file into the binary `.map` format:

```bash
go run github.com/drpaneas/gosprite64/cmd/mk2dmap \
  -in assets-src/level.json \
  -out assets/level.map
```

Flags:

| Flag | Default | Description |
|------|---------|-------------|
| `-in` | (required) | Input JSON path |
| `-out` | (required) | Output `.map` path |

The tool validates that:
- Map dimensions are non-zero
- `layer_count` is non-zero and matches the actual number of layers
- `cell_bits` is 8 or 16
- Chunk dimensions are non-zero
- Each layer's `cells` array length equals `width * height`
- 8-bit cells do not exceed 255

## Putting it together with go generate

The typical workflow uses `go:generate` to run both tools as part of your build:

```go
//go:generate sh -c "mkdir -p assets && go run github.com/drpaneas/gosprite64/cmd/mk2dsheet -in assets-src/tiles.png -out assets/tiles.sheet -tile-width 8 -tile-height 8 && go run github.com/drpaneas/gosprite64/cmd/mk2dmap -in assets-src/level.json -out assets/level.map"
```

Then combine the compiled assets into a bundle using `mk2dbundle`. See [Bundles and Loading](./bundles-and-loading.md) for details.

## Runtime map access

Once loaded through a bundle, you can query map properties:

```go
scene, _ := gosprite64.LoadScene(bundle)
m := scene.Map()

fmt.Println(m.Width(), m.Height())         // map size in tiles
fmt.Println(m.TileWidth(), m.TileHeight()) // tile size in pixels
fmt.Println(m.PixelWidth(), m.PixelHeight()) // total size in pixels
fmt.Println(m.LayerCount())                // number of layers

tile, ok := m.TileAt(0, 5, 3) // layer 0, column 5, row 3
info, ok := m.LayerInfo(0)    // SheetID and NonZeroTiles for layer 0
```
