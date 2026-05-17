# Bundles and Loading

A **bundle** is a manifest that packages tile sheets, maps, and animations together under named entries. Instead of loading each asset file individually, you open a single bundle and the runtime resolves all referenced assets by name. This page covers the bundle concept, how to open and load from bundles, and how `LoadScene` assembles everything into a renderable scene.

## Bundle concept

A bundle is a binary file (magic `BND2`) that acts as a table of contents. Each entry in the bundle records:

- **Kind** - what type of asset (sheet, map, or animation)
- **Name** - a human-readable lookup key
- **Path** - the file path to the compiled binary asset

The bundle itself does not embed asset data. It is a manifest that tells the runtime where to find each compiled `.sheet`, `.map`, and `.anim` file.

Bundles are created by the `mk2dbundle` tool:

```bash
go run github.com/drpaneas/gosprite64/cmd/mk2dbundle \
  -sheet assets/tiles.sheet \
  -map assets/level.map \
  -out assets/level.bundle
```

A typical `go:generate` line creates all three assets in one command chain. See [Pipeline Overview](./pipeline-overview.md) for the full build setup.

## Opening a bundle

### OpenBundle

`OpenBundle` reads a bundle manifest from the default asset filesystem (cartridge ROM on N64, embedded FS on desktop):

```go
bundle, err := gosprite64.OpenBundle("assets/level.bundle")
if err != nil {
    panic(err)
}
```

### OpenBundleWithLoader

`OpenBundleWithLoader` lets you supply a custom loader for testing or alternative storage backends:

```go
bundle, err := gosprite64.OpenBundleWithLoader("assets/level.bundle", myLoader)
if err != nil {
    panic(err)
}
```

The loader must implement the `Loader` interface, which provides raw byte access to asset files by path. The loader must not be nil.

## Loading individual assets

Once you have a `Bundle`, you can load specific assets by name:

### Bundle.LoadSheet

Loads a single tilesheet by its bundle entry name:

```go
sheet, err := bundle.LoadSheet("tiles")
if err != nil {
    panic(err)
}
info := sheet.Info()
fmt.Println(info.TileWidth, info.TileHeight, info.TileCount)
```

### Bundle.LoadMap

Loads the tile map by its bundle entry name:

```go
tileMap, err := bundle.LoadMap("level")
if err != nil {
    panic(err)
}
fmt.Println(tileMap.Width(), tileMap.Height())
```

### Bundle.LoadAnimation

Loads an animation set by its bundle entry name:

```go
anim, err := bundle.LoadAnimation("idle")
if err != nil {
    panic(err)
}
```

Each method looks up the entry by kind and name in the bundle manifest. If no matching entry is found, an error is returned.

## LoadScene - loading everything at once

For most games, you want to load all assets from a bundle into a ready-to-render scene. `LoadScene` does this in one call:

```go
bundle, err := gosprite64.OpenBundle("assets/level.bundle")
if err != nil {
    panic(err)
}

scene, err := gosprite64.LoadScene(bundle)
if err != nil {
    panic(err)
}
```

`LoadScene` iterates every entry in the bundle manifest and loads it by kind:

1. **Sheets** (kind 1) - loaded and stored in order. Multiple sheets are supported.
2. **Map** (kind 2) - exactly one map per bundle. Having zero or multiple maps is an error.
3. **Animations** (kind 3) - loaded and stored in order. Optional.

After loading, `LoadScene` performs these setup steps:

- Validates that sheet tile dimensions match the map
- Creates a default camera sized to the logical screen bounds
- Initializes the tile renderer
- Builds the internal prepared scene for fast drawing
- Computes initial statistics (sheet RAM, map RAM, layer count)

## Accessing loaded assets from a scene

```go
m := scene.Map()                  // the loaded Map
m.Width()                         // map width in tiles
m.PixelWidth()                    // map width in pixels

sheet := scene.Sheet(0)           // first sheet by index
sheet = scene.SheetByID(1)        // sheet by 1-based ID

info, sheet, ok := scene.LayerAssets(0) // layer 0's sheet + metadata
sheetInfo, ok := scene.LayerSheetInfo(0)

animSet := scene.AnimationByName("idle") // animation by name
animSet = scene.Animation(0)             // animation by index
```

## Drawing the scene

Pass a camera to render the visible portion of the map:

```go
func (g *Game) Draw() {
    gosprite64.ClearScreen()
    g.scene.Draw(g.camera)
}
```

If you pass `nil` as the camera, `Draw` uses the scene's default camera (positioned at the origin, sized to the screen). See [Camera and Scrolling](./camera-and-scrolling.md) for camera control.

## Runtime statistics

After drawing, you can query what happened:

```go
stats := scene.Stats()
fmt.Printf("visible tiles: %d\n", stats.VisibleTiles)
fmt.Printf("texture uploads: %d\n", stats.UploadCount)
fmt.Printf("sheets: %d, layers: %d\n", stats.SheetCount, stats.LayerCount)
fmt.Printf("sheet RAM: %d bytes\n", stats.SheetRAMBytes)
fmt.Printf("map RAM: %d bytes\n", stats.MapRAMBytes)
```

These stats are useful for profiling. `VisibleTiles` and `UploadCount` update every frame; the other fields are computed once at load time.

## Entry kinds

The bundle format supports three asset kinds, identified by numeric constants:

| Kind | Value | Asset type |
|------|-------|------------|
| Sheet | 1 | Compiled tilesheet (`.sheet`) |
| Map | 2 | Compiled tile map (`.map`) |
| Anim | 3 | Compiled animation set (`.anim`) |
