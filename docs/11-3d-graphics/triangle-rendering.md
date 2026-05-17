# Triangle Rendering

The `internal/rdpcpu` package provides CPU-side triangle rasterization for the N64's RDP. Instead of using the RSP's vertex pipeline (via `SPVertex` + `SP1Triangle`), these functions compute edge coefficients on the CPU and emit raw RDP triangle commands directly. This is useful for software rendering paths, debug visualization, and cases where the RSP vertex buffer is a bottleneck.

## How RDP triangles work

The N64's RDP rasterizes triangles using **edge coefficients**. For each triangle, the RDP needs:

1. **Edge coefficients** - slopes and starting positions for the three edges, sorted by Y coordinate
2. **Shade coefficients** (optional) - per-pixel color interpolation (Gouraud shading)
3. **Texture coefficients** (optional) - per-pixel texture coordinate interpolation

The `rdpcpu` functions compute all of this on the CPU and pack the result into raw 64-bit RDP command words. These words can be appended to a display list with `DPRaw`.

## Edge coefficient computation

All triangle functions share the same edge setup:

1. Vertices are sorted by Y coordinate (top to bottom)
2. The "major" edge spans from v1 (top) to v3 (bottom)
3. The "middle" vertex v2 splits the triangle into upper and lower halves
4. Inverse slopes are computed for all three edges
5. X positions are quantized to sub-pixel precision (1/4 pixel)
6. Y coordinates are clamped to the RDP's valid range

The `attrFactor` is the inverse of the cross product (nz = hx*my - hy*mx), used to interpolate attributes (shade, texture) across the triangle face.

## FillTriangle

Computes edge coefficients for a flat-colored triangle. Uses RDP opcode `0x08`.

```go
func FillTriangle(v1, v2, v3 [2]float32) []uint64
```

Vertices are in screen space: `[0]` is X, `[1]` is Y.

Returns 4 uint64 command words (edge coefficients only, no shade/texture).

```go
words := rdpcpu.FillTriangle(
    [2]float32{100, 50},   // v1 (screen X, Y)
    [2]float32{200, 150},  // v2
    [2]float32{50, 200},   // v3
)

dl := gfx.NewDisplayList(16)
dl.DPPipeSync()
// Set fill color/combine mode first...
dl.DPRaw(words...)
```

## ShadeTriangle

Computes edge + shade coefficients for a Gouraud-shaded triangle. Uses RDP opcode `0x0C`.

```go
func ShadeTriangle(v1, v2, v3 [2]float32, c1, c2, c3 [4]float32) []uint64
```

Vertices are screen-space positions. Colors are RGBA as floats in the range 0.0-1.0.

Returns 12 uint64 command words (4 edge + 8 shade coefficients).

```go
words := rdpcpu.ShadeTriangle(
    [2]float32{100, 50},   // position v1
    [2]float32{200, 150},  // position v2
    [2]float32{50, 200},   // position v3
    [4]float32{1, 0, 0, 1}, // color v1 (red)
    [4]float32{0, 1, 0, 1}, // color v2 (green)
    [4]float32{0, 0, 1, 1}, // color v3 (blue)
)

dl.DPRaw(words...)
```

The shade coefficients interpolate RGBA linearly across the triangle face using the gradients DsDx, DsDy, and DsDe (along-edge derivative).

## TexVertex

A vertex type for textured triangles that includes position, texture coordinates, and perspective correction:

```go
type TexVertex struct {
    X, Y float32   // screen-space position
    S, T float32   // texture coordinates
    InvW float32   // inverse W for perspective correction
}
```

## BuildTexturedTriangle

Computes edge + texture coefficients for a perspective-correct textured triangle. Uses RDP opcode `0x0A`.

```go
func BuildTexturedTriangle(tileIdx, mipmaps uint8, v1, v2, v3 TexVertex) []uint64
```

Parameters:
- `tileIdx` - RDP tile descriptor index (0-7)
- `mipmaps` - number of mipmap levels (0 for no mipmaps)
- `v1, v2, v3` - textured vertices with screen position, UV coords, and InvW

Returns 12 uint64 command words (4 edge + 8 texture coefficients).

```go
words := rdpcpu.BuildTexturedTriangle(
    0, 0, // tile 0, no mipmaps
    rdpcpu.TexVertex{X: 100, Y: 50,  S: 0,  T: 0,  InvW: 1},
    rdpcpu.TexVertex{X: 200, Y: 150, S: 32, T: 0,  InvW: 1},
    rdpcpu.TexVertex{X: 50,  Y: 200, S: 0,  T: 32, InvW: 1},
)

dl.DPRaw(words...)
```

The texture coefficients handle perspective correction by scaling S, T, and W by the minimum W value across all three vertices, then computing per-pixel gradients. The RDP uses these gradients to interpolate texture coordinates with perspective divide.

## Coordinate systems

All positions are in **screen space** (pixel coordinates after projection). If you are working with 3D world coordinates, you must project them through the model-view-projection matrix and viewport transform before passing them to these functions.

## When to use CPU-side triangles

| Use case | Approach |
|----------|----------|
| Standard 3D meshes | Use `SPVertex` + `SP1Triangle` (RSP pipeline) |
| Debug wireframes | Use `FillTriangle` |
| 2D effects, UI elements | Use `FillTriangle` or `ShadeTriangle` |
| Procedural geometry | Use `BuildTexturedTriangle` |
| Software rasterization | Use all three functions |

The RSP pipeline is faster for batched geometry because it handles vertex transformation in parallel. CPU-side triangles are better when you need fine control over individual triangles or when vertices are already in screen space.

## RDP triangle opcodes

| Opcode | Function | Coefficients |
|--------|----------|-------------|
| `0x08` | Fill triangle | Edge only (4 words) |
| `0x0A` | Textured triangle | Edge + texture (12 words) |
| `0x0C` | Shaded triangle | Edge + shade (12 words) |
| `0x0E` | Shaded + textured | Edge + shade + texture (20 words) |

The `triangle3d` example in the repository demonstrates CPU-side triangle rendering with perspective projection. See `examples/triangle3d/` for a working implementation.
