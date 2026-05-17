# 3D Math

The `math3d` package provides vectors, matrices, and transform functions for 3D rendering on the N64. All types use `float32` as the working format, matching the N64's RSP pipeline expectations. The package includes conversion to the N64's fixed-point hardware matrix format.

## Vectors

### Vec3

A 3-component vector for positions, directions, and normals:

```go
type Vec3 struct {
    X, Y, Z float32
}
```

Operations:

```go
a := math3d.Vec3{X: 1, Y: 2, Z: 3}
b := math3d.Vec3{X: 4, Y: 5, Z: 6}

sum   := a.Add(b)        // component-wise addition
diff  := a.Sub(b)        // component-wise subtraction
scaled := a.Scale(2.0)   // multiply all components by scalar
dot   := a.Dot(b)        // dot product
cross := a.Cross(b)      // cross product
length := a.Length()      // Euclidean length
norm  := a.Normalize()   // unit vector (returns zero vector if length is 0)
```

### Vec4

A 4-component vector for homogeneous coordinates and matrix multiplication:

```go
type Vec4 struct {
    X, Y, Z, W float32
}
```

`Vec4` is used primarily with `Mat4.MulVec4` for transforming points through projection matrices.

## Matrices

### Mat4

A 4x4 matrix in row-major order:

```go
type Mat4 [4][4]float32
```

### Identity

Returns a 4x4 identity matrix:

```go
m := math3d.Identity()
```

### Matrix multiplication

Multiply two matrices or transform a vector:

```go
result := a.Mul(b)          // Mat4 * Mat4
v := m.MulVec4(math3d.Vec4{X: 1, Y: 0, Z: 0, W: 1})  // Mat4 * Vec4
```

## Transform functions

All transform functions return a new `Mat4`. Chain them with `Mul` to compose transforms.

### Translate

```go
t := math3d.Translate(10, 0, -5)
```

Builds a translation matrix. Moves objects by (x, y, z) in world space.

### Scale

```go
s := math3d.Scale(2, 2, 2)
```

Builds a uniform or non-uniform scale matrix.

### Rotate

```go
r := math3d.Rotate(45, 0, 1, 0)  // 45 degrees around the Y axis
```

Builds a rotation matrix. The angle is in **degrees**. The axis (x, y, z) is normalized internally. Matches libultra's `guRotateF`.

### Composing transforms

Apply transforms right-to-left. To scale, then rotate, then translate:

```go
model := math3d.Translate(10, 0, 0).
    Mul(math3d.Rotate(45, 0, 1, 0)).
    Mul(math3d.Scale(2, 2, 2))
```

## View matrix

### LookAt

Builds a view matrix that positions and orients the camera:

```go
view := math3d.LookAt(
    0, 5, 10,    // eye position
    0, 0, 0,     // look-at target
    0, 1, 0,     // up direction
)
```

Matches libultra's `guLookAtReflectF`. The eye looks toward the target with the given up vector defining the camera's vertical orientation.

## Projection matrices

### Perspective

Builds a perspective projection matrix for 3D rendering:

```go
proj, perspNorm := math3d.Perspective(
    45,     // fovy: vertical field of view in degrees
    1.333,  // aspect: width / height (e.g. 320/240)
    10,     // near plane distance
    1000,   // far plane distance
    1.0,    // scale factor
)
```

Returns two values:

- `Mat4` - the projection matrix
- `uint16` - the `perspNorm` value that the RSP needs for correct clipping. Pass this to the display list via `SPPerspNormalize`.

Matches libultra's `guPerspectiveF`.

### Ortho

Builds an orthographic projection matrix for 2D overlays or isometric views:

```go
ortho := math3d.Ortho(
    0, 320,     // left, right
    240, 0,     // bottom, top (flipped for screen coords)
    -1, 1,      // near, far
    1.0,        // scale factor
)
```

Matches libultra's `guOrthoF`.

## Setting up a perspective projection

A typical 3D scene setup combines all three matrices:

```go
// Projection
proj, perspNorm := math3d.Perspective(45, 320.0/240.0, 10, 1000, 1.0)

// View
view := math3d.LookAt(
    0, 100, 200,  // eye
    0, 0, 0,      // target
    0, 1, 0,      // up
)

// Model (per object)
model := math3d.Translate(0, 0, 0).
    Mul(math3d.Rotate(angle, 0, 1, 0)).
    Mul(math3d.Scale(1, 1, 1))

// Combined model-view-projection
mvp := proj.Mul(view).Mul(model)
```

## N64 fixed-point conversion

The N64 RSP uses a fixed-point matrix format (s15.16) stored as 16 uint32 words. The first 8 words hold integer portions, the last 8 hold fractional portions. Total size: 64 bytes.

### ToN64Mtx

Convert a float32 `Mat4` to the hardware format:

```go
type N64Mtx [16]uint32

hwMtx := mvp.ToN64Mtx()
```

Matches libultra's `guMtxF2L`. The resulting `N64Mtx` can be DMA'd to RDRAM and loaded by the RSP via `SPMatrix`.

### FromN64Mtx

Convert back from hardware format to float32 (useful for debugging):

```go
floatMtx := math3d.FromN64Mtx(hwMtx)
```

Matches libultra's `guMtxL2F`.

## Complete example: perspective setup with display list

```go
// Build matrices
proj, perspNorm := math3d.Perspective(45, 320.0/240.0, 10, 1000, 1.0)
view := math3d.LookAt(0, 100, 200, 0, 0, 0, 0, 1, 0)
model := math3d.Translate(0, 0, 0).Mul(math3d.Rotate(angle, 0, 1, 0))

// Convert to N64 format
projMtx := proj.ToN64Mtx()
mvMtx := view.Mul(model).ToN64Mtx()

// Load into RSP via display list
dl := gfx.NewDisplayList(32)
dl.SPPerspNormalize(perspNorm)
dl.SPMatrix(projAddr, gfx.MtxProjection|gfx.MtxLoad|gfx.MtxNoPush)
dl.SPMatrix(mvAddr, gfx.MtxModelView|gfx.MtxLoad|gfx.MtxNoPush)
```

See [Display Lists](./display-lists.md) for the full command reference.
