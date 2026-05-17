# Display Lists

The `gfx` package provides a Go API for building N64 display lists. A display list is a sequence of 64-bit commands that drive the RSP (Reality Signal Processor) and RDP (Reality Display Processor). GoSprite64's display list API mirrors the F3DEX2 microcode command set used by most N64 games.

## F3DEX2 microcode context

The N64's RSP runs microcode that interprets display list commands. F3DEX2 is the most common variant, supporting:

- 32 vertex buffer slots (vs 16 in original Fast3D)
- Two-triangle commands for better throughput
- Matrix stack operations for hierarchical transforms
- Geometry mode flags for lighting, culling, and shading

GoSprite64's display list commands match the F3DEX2 encoding. The `gfx` package handles the bit packing so you work with typed Go calls instead of raw uint64 values.

## DisplayList type

```go
type DisplayList struct {
    cmds []Gfx
}
```

Each command is a `Gfx` struct holding two 32-bit words (matching the N64's 64-bit command format):

```go
type Gfx struct {
    W0, W1 uint32
    Raw    []uint64  // for multi-word RDP commands
}
```

### Creating and managing

```go
dl := gfx.NewDisplayList(64)    // pre-allocate capacity for 64 commands
dl.Len()                         // number of commands so far
dl.Reset()                       // clear for reuse without deallocating
cmds := dl.Commands()            // access the raw command slice
```

## RSP commands (SP prefix)

SP commands are processed by the RSP microcode.

### SPMatrix

Loads or multiplies a matrix into the RSP matrix stack:

```go
dl.SPMatrix(addr, gfx.MtxProjection|gfx.MtxLoad|gfx.MtxNoPush)
dl.SPMatrix(addr, gfx.MtxModelView|gfx.MtxLoad|gfx.MtxPush)
```

`addr` is the RDRAM address of an `N64Mtx` (64 bytes). Flags:

| Flag | Value | Description |
|------|-------|-------------|
| `MtxModelView` | 0x00 | Target the model-view stack |
| `MtxProjection` | 0x01 | Target the projection stack |
| `MtxMul` | 0x00 | Multiply with current matrix |
| `MtxLoad` | 0x02 | Replace current matrix |
| `MtxNoPush` | 0x00 | Do not push before write |
| `MtxPush` | 0x04 | Push current matrix before write |

### SPVertex

Loads vertices into the RSP vertex buffer:

```go
dl.SPVertex(addr, n, v0)
// addr: RDRAM address of vertex array
// n:    number of vertices to load (1-16)
// v0:   starting vertex buffer index
```

Each vertex is 16 bytes, matching the `Vtx` struct:

```go
type Vtx struct {
    X, Y, Z int16
    Flag    uint16
    S, T    int16      // texture coordinates
    R, G, B, A uint8   // vertex color
}
```

For lit meshes, use `VtxN` which replaces color with normals:

```go
type VtxN struct {
    X, Y, Z int16
    Flag    uint16
    S, T    int16
    NX, NY, NZ int8
    A          uint8
}
```

### SP1Triangle

Draws a single triangle from vertex buffer indices:

```go
dl.SP1Triangle(v0, v1, v2, flag)
// v0, v1, v2: vertex buffer indices (0-31)
// flag: which vertex is used for flat shading (usually 0)
```

### SP2Triangles

Draws two triangles. In the current implementation, this emits two `SP1Triangle` commands:

```go
dl.SP2Triangles(v00, v01, v02, flag0, v10, v11, v12, flag1)
```

### SPDisplayList

Calls a child display list (with return - like a function call):

```go
dl.SPDisplayList(childAddr)
```

### SPBranchList

Jumps to another display list without return:

```go
dl.SPBranchList(otherAddr)
```

### SPEndDisplayList

Terminates display list processing. Every display list must end with this:

```go
dl.SPEndDisplayList()
```

### SPSegment

Sets a segment register for address translation. The RSP uses 16 segments (0-15) to translate segmented addresses to physical RDRAM addresses:

```go
dl.SPSegment(6, baseAddr) // segment 6 = baseAddr
```

### SPViewport

Sets the viewport parameters:

```go
vp := gfx.NewViewport(320, 240)
dl.SPViewport(vpAddr)
```

The `Viewport` struct matches the N64's `Vp_t`:

```go
type Viewport struct {
    ScaleX, ScaleY, ScaleZ, ScalePad int16
    TransX, TransY, TransZ, TransPad int16
}
```

`NewViewport(width, height)` creates a viewport centered on the screen with proper scale factors.

### SPPerspNormalize

Sets the perspective normalization value for correct RSP clipping:

```go
dl.SPPerspNormalize(perspNorm) // perspNorm from math3d.Perspective()
```

### SPSetGeometryMode / SPClearGeometryMode

Enable or disable geometry mode flags:

```go
dl.SPSetGeometryMode(gfx.GeomShade | gfx.GeomZBuffer)
dl.SPClearGeometryMode(gfx.GeomCullBack)
```

## RDP commands (DP prefix)

DP commands are forwarded to the RDP for pixel processing.

### DPSetScissor

Defines the scissor rectangle. Only pixels inside this rectangle are drawn:

```go
dl.DPSetScissor(0, 0, 0, 320, 240) // mode, ulx, uly, lrx, lry
```

Coordinates are in screen pixels.

### DPSetCombineMode

Configures the RDP color combiner. The combiner controls how texture, shade, and environment colors are mixed:

```go
dl.DPSetCombineMode(w0hi, w1)
```

The two parameters are the raw 64-bit combiner encoding split into upper and lower words. Predefined modes are typically used as constants.

### DPPipeSync

Inserts a pipeline sync barrier. Required before changing RDP state (tile descriptors, combine mode, render mode) to ensure previous rendering has completed:

```go
dl.DPPipeSync()
```

### DPTileSync / DPLoadSync / DPFullSync

Other synchronization barriers:

```go
dl.DPTileSync()     // sync before tile descriptor changes
dl.DPLoadSync()     // sync before texture loads
dl.DPFullSync()     // signal RDP completion (end of frame)
```

### DPSetColorImage

Sets the RDP's render target (framebuffer):

```go
dl.DPSetColorImage(fmt, siz, width, addr)
```

### DPSetFillColor

Sets the fill color for `DPFillRect`:

```go
dl.DPSetFillColor(color) // packed 16-bit RGBA doubled, or 32-bit
```

### DPFillRect

Fills a rectangle with the current fill color. Coordinates are in 10.2 fixed-point:

```go
dl.DPFillRect(ulx, uly, lrx, lry)
```

### DPSetTextureImage / DPSetTile / DPLoadBlock / DPSetTileSize

Texture loading sequence:

```go
dl.DPSetTextureImage(fmt, siz, width, texAddr)
dl.DPSetTile(fmt, siz, line, tmem, tile, palette, cmt, maskt, shiftt, cms, masks, shifts)
dl.DPLoadBlock(tile, uls, ult, lrs, dxt)
dl.DPSetTileSize(tile, uls, ult, lrs, lrt)
```

### DPSetPrimColor / DPSetEnvColor / DPSetFogColor

Set special RDP colors:

```go
dl.DPSetPrimColor(minLevel, fracLevel, r, g, b, a)
dl.DPSetEnvColor(r, g, b, a)
dl.DPSetFogColor(r, g, b, a)
```

### DPSetRenderMode

Sets the RDP render mode via `SetOtherModeL`:

```go
dl.DPSetRenderMode(cycle0, cycle1)
```

### DPSetZImage

Sets the Z-buffer address:

```go
dl.DPSetZImage(zbufAddr)
```

### DPRaw

Appends raw 64-bit command words for advanced use cases like CPU-built triangle commands:

```go
dl.DPRaw(words...)
```

## Typical rendering sequence

A minimal 3D frame looks like this:

```go
dl := gfx.NewDisplayList(128)

// Set up viewport and projection
dl.SPSegment(0, 0)
dl.SPViewport(vpAddr)
dl.SPPerspNormalize(perspNorm)
dl.SPMatrix(projAddr, gfx.MtxProjection|gfx.MtxLoad|gfx.MtxNoPush)
dl.SPMatrix(mvAddr, gfx.MtxModelView|gfx.MtxLoad|gfx.MtxNoPush)

// Set up RDP state
dl.DPPipeSync()
dl.DPSetScissor(0, 0, 0, 320, 240)
dl.DPSetCombineMode(combinerHi, combinerLo)

// Load and draw geometry
dl.SPVertex(vertAddr, 3, 0)
dl.SP1Triangle(0, 1, 2, 0)

// Finish
dl.DPFullSync()
dl.SPEndDisplayList()
```

See [Triangle Rendering](./triangle-rendering.md) for CPU-side triangle commands that bypass the RSP vertex pipeline.
