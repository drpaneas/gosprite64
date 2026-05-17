# Step 6: Camera Following

Make the camera follow the player so the world scrolls as they move.

## What you will learn

- Configuring the Camera with `FollowTarget`, `FollowSpeed`, and `Bounds`
- Smooth camera lerp with `UpdateFollow`
- Clamping the camera to the map edges with `ClampToBounds`

## What changed from Step 5

The code is the same as Step 5 with three changes: a new `math2d` import, a new camera setup in `Init`, and two new calls at the end of `Update`. The movement, animation, and draw code are unchanged.

### New import

```go
import (
	"github.com/drpaneas/gosprite64"
	"github.com/drpaneas/gosprite64/math2d"
)
```

### New camera setup in Init

Replace the old `g.camera = &gosprite64.Camera{Width: 288, Height: 216}` with:

```go
g.playerX = 80
g.playerY = 180

g.camera = &gosprite64.Camera{
	Width:       288,
	Height:      216,
	FollowSpeed: 0.1,
}
g.camera.FollowTarget = &math2d.Vec2{X: g.playerX, Y: g.playerY}
g.camera.Bounds = &math2d.Rect{
	X: 0, Y: 0,
	W: float32(scene.Map().PixelWidth()),
	H: float32(scene.Map().PixelHeight()),
}
```

The starting position moves to 80,180 (near the bottom-left of the map) so there is room to scroll in every direction.

### New lines at the end of Update

After the animation code, add:

```go
g.camera.FollowTarget.X = g.playerX
g.camera.FollowTarget.Y = g.playerY
g.camera.UpdateFollow()
g.camera.ClampToBounds()
```

Each frame we update the target position to match the player, then tell the camera to move toward it and stay inside the map.

## How it works

### The follow system

`UpdateFollow` calculates where the camera wants to be (player position minus half the screen size, centering the player) and lerps toward it. `FollowSpeed` of 0.1 means the camera covers 10% of the remaining distance per frame - fast enough to keep up, slow enough to feel smooth. Setting it to 1.0 snaps instantly.

| Field | Type | Purpose |
|-------|------|---------|
| `FollowTarget` | `*math2d.Vec2` | World position the camera tracks |
| `FollowSpeed` | `float32` | Lerp factor, 0.0 to 1.0 |
| `Bounds` | `*math2d.Rect` | Rectangle the camera stays inside |

### Clamping to bounds

Without bounds the camera would show empty space beyond the map edges. `ClampToBounds` restricts the camera to the `Bounds` rectangle. The valid range is 0,0 to `mapWidth - cameraWidth` on each axis. If the map is smaller than the camera viewport, the camera pins to 0.

### math2d types

`math2d.Vec2` is `{X, Y float32}` for positions. `math2d.Rect` is `{X, Y, W, H float32}` for axis-aligned rectangles. Both are value types in the `github.com/drpaneas/gosprite64/math2d` package.

## Build and run

```bash
go generate ./examples/platformer
GOENV=n64.env go1.24.5-embedded build -o examples/platformer/game.elf ./examples/platformer
```

Walk around with the D-pad. The camera now smoothly follows the player, and the tile world scrolls underneath. Walk to any edge of the map and the camera stops scrolling while the player can still move within the visible area.
