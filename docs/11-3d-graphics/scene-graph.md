# Scene Graph

The `scene3d` package provides a hierarchical scene graph for organizing 3D objects. Nodes form a tree where each child inherits its parent's transform. The scene graph handles transform composition, camera setup, level-of-detail selection, and depth-first traversal for rendering.

## Node types

Every node has a `Type` field that determines its behavior during traversal:

| Type | Value | Purpose |
|------|-------|---------|
| `NodeTransform` | 0 | Pure transform (position/rotation/scale) |
| `NodeCamera` | 1 | Camera with projection parameters |
| `NodeMesh` | 2 | Renderable mesh with a display list |
| `NodeBillboard` | 3 | Billboard that always faces the camera |
| `NodeRenderFunc` | 4 | Custom rendering callback |
| `NodeLOD` | 5 | Level-of-detail selector |
| `NodeGroup` | 6 | Generic group for organizing children |

## Creating nodes

### NewNode

Creates a group node for organizing children. Visibility defaults to `true` and scale defaults to (1, 1, 1):

```go
root := scene3d.NewNode("world")
```

### NewMeshNode

Creates a mesh node that renders a pre-built display list:

```go
dl := gfx.NewDisplayList(64)
// ... build display list commands ...
dl.SPEndDisplayList()

cube := scene3d.NewMeshNode("cube", dl)
```

The `MeshData` also supports a bounding sphere for frustum culling:

```go
cube.Mesh.BoundsCenter = math3d.Vec3{X: 0, Y: 0, Z: 0}
cube.Mesh.BoundsRadius = 5.0
```

### NewPerspectiveCamera

Creates a camera node with perspective projection:

```go
cam := scene3d.NewPerspectiveCamera("main-cam", 45, 1.333, 10, 1000)
// fov=45 degrees, aspect=1.333, near=10, far=1000
```

### NewOrthoCamera

Creates a camera node with orthographic projection:

```go
cam := scene3d.NewOrthoCamera("ui-cam", 0, 320, 240, 0, -1, 1)
```

### NewLODNode

Creates a level-of-detail node that selects a child based on camera distance:

```go
lod := scene3d.NewLODNode("tree-lod", []scene3d.LODLevel{
    {MaxDistance: 100, Child: highDetailTree},
    {MaxDistance: 500, Child: lowDetailTree},
})
```

The `LODData.Select(distance)` method returns the first child whose `MaxDistance` is >= the given distance, or `nil` if no level matches.

## Building a node hierarchy

Use `AddChild` to attach nodes to parents. Transforms compose through the hierarchy:

```go
root := scene3d.NewNode("root")

// Position a group
platform := scene3d.NewNode("platform")
platform.Position = math3d.Vec3{X: 10, Y: 0, Z: 0}
root.AddChild(platform)

// Add a mesh to the group - it inherits the platform's position
crate := scene3d.NewMeshNode("crate", crateDL)
crate.Position = math3d.Vec3{X: 0, Y: 2, Z: 0} // 2 units above the platform
platform.AddChild(crate)
```

## Node transform

Each node has `Position`, `Rotation`, and `Scale` fields:

```go
type Node struct {
    Position math3d.Vec3   // translation
    Rotation math3d.Vec3   // euler angles in degrees (X, Y, Z)
    Scale    math3d.Vec3   // scale factors
    Visible  bool          // if false, node and children are skipped
    // ...
}
```

### LocalTransform

Computes the local transformation matrix from position, rotation, and scale. The transform order is: translate, then rotate (Z * Y * X), then scale:

```go
localMat := node.LocalTransform()
```

This returns a `math3d.Mat4` that can be used directly or composed with parent transforms.

## CameraData

The `CameraData` struct holds projection parameters:

```go
type CameraData struct {
    FOV    float32  // field of view in degrees (perspective only)
    Aspect float32  // width/height ratio
    Near   float32
    Far    float32
    Ortho  bool     // if true, use orthographic projection

    // Orthographic extents (only used when Ortho is true)
    Left, Right, Bottom, Top float32
}
```

Get the projection matrix:

```go
projMat, perspNorm := camData.ProjectionMatrix()
```

For perspective cameras, `perspNorm` is the RSP normalization value. For orthographic cameras, it is 0.

## MeshData

References a pre-built display list for rendering:

```go
type MeshData struct {
    DisplayList  *gfx.DisplayList
    BoundsCenter math3d.Vec3
    BoundsRadius float32
}
```

## LODData

Selects children based on distance from the camera:

```go
type LODData struct {
    Levels []LODLevel
}

type LODLevel struct {
    MaxDistance float32
    Child      *Node
}
```

During traversal, LOD nodes do not traverse their `Children` list. Instead, they call `LOD.Select(distance)` and traverse only the selected child.

## Scene

A `Scene` holds the root node and the active camera:

```go
type Scene struct {
    Root   *Node
    Camera *Node
}
```

### NewScene

Creates a scene with an empty root group node named "root":

```go
scene := scene3d.NewScene()
```

### Setting up the camera

Assign a camera node to the scene:

```go
cam := scene3d.NewPerspectiveCamera("main", 45, 320.0/240.0, 10, 1000)
cam.Position = math3d.Vec3{X: 0, Y: 100, Z: 200}
scene.Root.AddChild(cam)
scene.Camera = cam
```

## Traversal

### Traverse

`Traverse` walks the scene graph depth-first, calling a visitor function for each visible node. The `RenderContext` maintains a matrix stack:

```go
ctx := scene3d.NewRenderContext()
scene.Traverse(ctx, func(node *scene3d.Node, rc *scene3d.RenderContext) {
    if node.Type == scene3d.NodeMesh {
        modelView := rc.CurrentMatrix()
        // render the mesh with this transform
    }
})
```

The traversal automatically:
- Skips nodes where `Visible` is `false`
- Pushes/pops the matrix stack at each level
- Multiplies each node's `LocalTransform` into the stack
- For LOD nodes, selects the appropriate child based on camera distance

### RenderContext

The `RenderContext` tracks state during traversal:

```go
type RenderContext struct {
    ViewMatrix       math3d.Mat4
    ProjectionMatrix math3d.Mat4
    CameraPosition   math3d.Vec3
    MatrixStack      []math3d.Mat4
}
```

Matrix stack operations:

```go
ctx.PushMatrix()                  // duplicate top of stack
ctx.PopMatrix()                   // remove top of stack
ctx.MultiplyMatrix(localMat)      // multiply top by given matrix
current := ctx.CurrentMatrix()    // read top of stack
```

## DrawScene

`DrawScene` is the high-level rendering entry point. It sets up the camera, traverses the scene graph, and executes display lists for all visible mesh nodes:

```go
scene3d.DrawScene(scene)
```

Internally, `DrawScene`:

1. Creates a `RenderContext`
2. If a camera node is assigned, computes the projection and view matrices
3. Traverses the graph, executing the display list for each `NodeMesh`
4. Calls `gfx.Flush()` to submit all commands to the RDP

On non-N64 builds, `DrawScene` is a no-op.

## Complete example

```go
// Create scene
scene := scene3d.NewScene()

// Add camera
cam := scene3d.NewPerspectiveCamera("cam", 45, 320.0/240.0, 10, 1000)
cam.Position = math3d.Vec3{X: 0, Y: 50, Z: 100}
scene.Root.AddChild(cam)
scene.Camera = cam

// Add a mesh
cubeDL := gfx.NewDisplayList(64)
// ... build cube display list ...
cubeDL.SPEndDisplayList()
cube := scene3d.NewMeshNode("cube", cubeDL)
cube.Position = math3d.Vec3{X: 0, Y: 0, Z: 0}
cube.Rotation = math3d.Vec3{Y: 45} // rotated 45 degrees around Y
scene.Root.AddChild(cube)

// Add LOD object
highDetail := scene3d.NewMeshNode("tree-hi", highDL)
lowDetail := scene3d.NewMeshNode("tree-lo", lowDL)
treeLOD := scene3d.NewLODNode("tree", []scene3d.LODLevel{
    {MaxDistance: 200, Child: highDetail},
    {MaxDistance: 1000, Child: lowDetail},
})
treeLOD.Position = math3d.Vec3{X: 50, Y: 0, Z: -30}
scene.Root.AddChild(treeLOD)

// Render
scene3d.DrawScene(scene)
```
