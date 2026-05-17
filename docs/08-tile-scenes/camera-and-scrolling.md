# Camera and Scrolling

The `Camera` struct defines the visible region of your tile world. By moving the camera, you scroll through maps that are larger than the screen. This page covers camera creation, manual scrolling, smooth follow, bounds clamping, screen shake, and the `scene.Draw(camera)` rendering call.

## Camera struct

```go
type Camera struct {
    X, Y          int       // top-left corner of the viewport in world pixels
    Width, Height int       // viewport size in pixels

    Zoom         float32    // zoom level (0 or unset defaults to 1.0)

    FollowTarget *math2d.Vec2 // world position to follow
    FollowSpeed  float32      // lerp speed: 0.0-1.0 (1.0 = instant snap)

    Bounds       *math2d.Rect // optional clamping rectangle
}
```

Create a camera by specifying the viewport dimensions. On the N64, the standard logical resolution is 288x216:

```go
camera := &gosprite64.Camera{Width: 288, Height: 216}
```

## Drawing with a camera

Pass the camera to `scene.Draw` each frame:

```go
func (g *Game) Draw() {
    gosprite64.ClearScreen()
    g.scene.Draw(g.camera)
}
```

`scene.Draw` renders only the tiles visible within the camera's viewport. The renderer uses chunk-based culling to skip off-screen regions efficiently. If the camera is `nil`, the scene uses its default camera (positioned at the origin).

## Manual scrolling

Move the camera by updating `X` and `Y` directly. The `simplegame` example scrolls with the D-pad:

```go
func (g *Game) Update() {
    speed := 1

    if gosprite64.IsButtonDown(gosprite64.ButtonDPadUp) {
        g.camera.Y -= speed
    }
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadDown) {
        g.camera.Y += speed
    }
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft) {
        g.camera.X -= speed
    }
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadRight) {
        g.camera.X += speed
    }
}
```

## Clamping to map bounds

Without clamping, the camera can scroll past the map edges, showing empty space. There are two approaches to restrict the camera.

### Manual clamping

Compute the maximum scroll position from the map's pixel dimensions:

```go
m := g.scene.Map()
maxX := m.PixelWidth() - g.camera.Width
maxY := m.PixelHeight() - g.camera.Height

if g.camera.X < 0 {
    g.camera.X = 0
}
if g.camera.Y < 0 {
    g.camera.Y = 0
}
if g.camera.X > maxX {
    g.camera.X = maxX
}
if g.camera.Y > maxY {
    g.camera.Y = maxY
}
```

### Using Camera.Bounds

Set the `Bounds` field and call `ClampToBounds` for automatic clamping:

```go
m := g.scene.Map()
g.camera.Bounds = &math2d.Rect{
    X: 0, Y: 0,
    W: float32(m.PixelWidth()),
    H: float32(m.PixelHeight()),
}

// In Update(), after moving the camera:
g.camera.ClampToBounds()
```

`ClampToBounds` ensures the viewport stays within the bounds rectangle. It accounts for the viewport size so the right/bottom edges do not exceed the map. If `Bounds` is nil, the call is a no-op.

## Smooth camera follow

To make the camera track a player or other target smoothly, set `FollowTarget` and `FollowSpeed`, then call `UpdateFollow` each frame:

```go
g.camera.FollowTarget = &math2d.Vec2{X: playerX, Y: playerY}
g.camera.FollowSpeed = 0.1 // smooth lerp (1.0 = instant snap)

func (g *Game) Update() {
    g.camera.FollowTarget.X = g.playerX
    g.camera.FollowTarget.Y = g.playerY
    g.camera.UpdateFollow()
    g.camera.ClampToBounds()
}
```

`UpdateFollow` lerps the camera position toward the target, centering it in the viewport. A speed of `0.1` gives a smooth trailing feel; `1.0` snaps instantly. When the camera is within 1 pixel of the target, it snaps to avoid sub-pixel jitter.

## Coordinate conversion

`WorldToScreen` converts a world position to screen coordinates, accounting for camera position and zoom:

```go
screenX, screenY := g.camera.WorldToScreen(worldX, worldY)
```

This is useful for placing UI elements or debug overlays relative to world objects.

## Zoom

Set the `Zoom` field to scale the viewport. A zoom of `2.0` means each world pixel occupies 2 screen pixels (zoomed in). The default zoom is `1.0`.

```go
g.camera.Zoom = 2.0
```

`EffectiveZoom` returns the active zoom level, defaulting to 1.0 when `Zoom` is zero:

```go
z := g.camera.EffectiveZoom()
```

## Screen shake

Camera shake adds visual impact to events like explosions or hits. The system uses a trauma model where shake magnitude is the square of the trauma value, producing a natural decay.

```go
// On impact:
g.camera.AddTrauma(0.5) // 0.0-1.0, multiple hits accumulate up to 1.0

// Every frame in Update():
g.camera.UpdateShake() // decays trauma over time

// When drawing, apply the shake offset:
shakeX, shakeY := g.camera.ShakeOffset()
// Use shakeX/shakeY as an additional draw offset
```

- `AddTrauma(amount)` adds to the trauma value, capping at 1.0
- `UpdateShake()` decays trauma by 1/60 per frame
- `ShakeOffset()` returns pixel displacement for the current frame, with a maximum offset of 8 pixels in each direction

## Complete example

This is a condensed version of `examples/simplegame`:

```go
type Game struct {
    scene  *gosprite64.Scene
    camera *gosprite64.Camera
}

func (g *Game) Init() {
    bundle, err := gosprite64.OpenBundle("assets/level.bundle")
    if err != nil {
        panic(err)
    }
    scene, err := gosprite64.LoadScene(bundle)
    if err != nil {
        panic(err)
    }
    g.scene = scene
    g.camera = &gosprite64.Camera{Width: 288, Height: 216}
}

func (g *Game) Update() {
    speed := 1
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadUp)    { g.camera.Y -= speed }
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadDown)  { g.camera.Y += speed }
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft)  { g.camera.X -= speed }
    if gosprite64.IsButtonDown(gosprite64.ButtonDPadRight) { g.camera.X += speed }

    m := g.scene.Map()
    maxX := m.PixelWidth() - g.camera.Width
    maxY := m.PixelHeight() - g.camera.Height
    if g.camera.X < 0    { g.camera.X = 0 }
    if g.camera.Y < 0    { g.camera.Y = 0 }
    if g.camera.X > maxX { g.camera.X = maxX }
    if g.camera.Y > maxY { g.camera.Y = maxY }
}

func (g *Game) Draw() {
    gosprite64.ClearScreen()
    g.scene.Draw(g.camera)

    stats := g.scene.Stats()
    gosprite64.DrawText(fmt.Sprintf("vis:%d", stats.VisibleTiles), 2, 2, gosprite64.White)
}
```
