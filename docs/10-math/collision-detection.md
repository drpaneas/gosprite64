# Collision Detection

The `math2d` package provides AABB collision detection and resolution for 2D games.

```go
import "github.com/drpaneas/gosprite64/math2d"
```

## AABB collision functions

All functions operate on `math2d.Rect` values - axis-aligned bounding boxes defined by position and size. See [Rectangles](rectangles.md) for the full `Rect` API.

### AABBOverlap

Returns `true` if two rectangles overlap:

```go
player := math2d.Rect{X: 50, Y: 50, W: 16, H: 24}
coin   := math2d.Rect{X: 55, Y: 48, W: 8, H: 8}

if math2d.AABBOverlap(player, coin) {
    collectCoin()
}
```

### AABBPenetration

Returns the minimum translation vector (MTV) needed to push rect `a` out of rect `b`. The MTV points along the axis of least overlap. Returns `false` if there is no overlap.

```go
pen, ok := math2d.AABBPenetration(player, wall)
if ok {
    // pen is a Vec2 that moves 'player' out of 'wall'
    player.X += pen.X
    player.Y += pen.Y
}
```

The MTV always pushes along a single axis (the shorter overlap), so either `pen.X` or `pen.Y` will be zero.

### AABBResolve

A convenience function that returns a copy of rect `a` moved so it no longer overlaps `b`:

```go
playerRect = math2d.AABBResolve(playerRect, wallRect)
```

This is equivalent to computing the penetration vector and adding it to `a`'s position.

### AABBSweep

Performs a swept (continuous) collision test. Moves rect `a` by a velocity vector and checks for collision with static rect `b`. Returns three values:

- `hit` - whether a collision occurred during the sweep
- `t` - the fraction of velocity at first contact (0.0 to 1.0)
- `normal` - the surface normal at the collision point

```go
func AABBSweep(a Rect, vel Vec2, b Rect) (hit bool, t float32, normal Vec2)
```

```go
vel := math2d.Vec2{X: 5, Y: 0}
hit, t, normal := math2d.AABBSweep(player, vel, wall)
if hit {
    // Move only up to the contact point
    player.X += vel.X * t
    player.Y += vel.Y * t

    // Slide along the wall using the normal
    if normal.X != 0 {
        vel.X = 0  // hit a vertical wall, stop horizontal movement
    }
    if normal.Y != 0 {
        vel.Y = 0  // hit a horizontal wall, stop vertical movement
    }
}
```

If the two rects already overlap before the sweep, it returns `(true, 0, {0,0})`.

## Layers and colliders

For games with many object types (player, enemies, projectiles, pickups), use layers to control which objects can collide with each other.

### Layer

`Layer` is a `uint32` bitmask. Define your game's layers as powers of two:

```go
const (
    LayerPlayer     math2d.Layer = 1 << 0  // 0x01
    LayerEnemy      math2d.Layer = 1 << 1  // 0x02
    LayerProjectile math2d.Layer = 1 << 2  // 0x04
    LayerPickup     math2d.Layer = 1 << 3  // 0x08
)
```

Two built-in constants:

- `math2d.LayerNone` (0x00000000) - matches nothing
- `math2d.LayerAll` (0xFFFFFFFF) - matches everything

Use `Matches` to test whether two layers share any bits:

```go
a := LayerPlayer | LayerEnemy
b := LayerEnemy
a.Matches(b) // true - both include LayerEnemy
```

### Collider

A `Collider` combines a bounding box with layer membership:

```go
type Collider struct {
    Bounds Rect    // the AABB
    Layer  Layer   // what this object IS
    Mask   Layer   // what this object COLLIDES WITH
}
```

- `Layer` identifies the object's type.
- `Mask` defines which layers this object reacts to.

```go
playerCol := math2d.Collider{
    Bounds: math2d.Rect{X: 50, Y: 50, W: 16, H: 24},
    Layer:  LayerPlayer,
    Mask:   LayerEnemy | LayerPickup, // collide with enemies and pickups
}

enemyCol := math2d.Collider{
    Bounds: math2d.Rect{X: 80, Y: 50, W: 16, H: 16},
    Layer:  LayerEnemy,
    Mask:   LayerPlayer | LayerProjectile, // collide with player and bullets
}
```

### ColliderOverlap

Tests both spatial overlap and layer compatibility. The check is bidirectional - either collider's mask matching the other's layer is sufficient:

```go
if math2d.ColliderOverlap(playerCol, enemyCol) {
    player.TakeDamage()
}
```

Two objects only collide if their bounding boxes overlap AND at least one of them has a mask that includes the other's layer.

## Platformer collision example

```go
const (
    LayerPlayer math2d.Layer = 1 << 0
    LayerSolid  math2d.Layer = 1 << 1
    LayerCoin   math2d.Layer = 1 << 2
)

type Entity struct {
    Pos      math2d.Vec2
    Size     math2d.Vec2
    Vel      math2d.Vec2
    Collider math2d.Collider
}

func (e *Entity) Bounds() math2d.Rect {
    return math2d.Rect{X: e.Pos.X, Y: e.Pos.Y, W: e.Size.X, H: e.Size.Y}
}

func (g *Game) Update() {
    p := &g.player

    // Apply gravity
    p.Vel.Y += 0.5

    // Sweep against all solid tiles
    bounds := p.Bounds()
    for _, wall := range g.walls {
        hit, t, normal := math2d.AABBSweep(bounds, p.Vel, wall)
        if hit {
            // Move up to contact
            p.Vel = math2d.Vec2{
                X: p.Vel.X * t,
                Y: p.Vel.Y * t,
            }
            // Cancel velocity into the wall
            if normal.Y < 0 {
                p.Vel.Y = 0
                p.OnGround = true
            }
            if normal.Y > 0 {
                p.Vel.Y = 0 // hit ceiling
            }
            if normal.X != 0 {
                p.Vel.X = 0 // hit side wall
            }
        }
    }

    // Apply velocity
    p.Pos.X += p.Vel.X
    p.Pos.Y += p.Vel.Y

    // Check coin pickups (simple overlap, no physics)
    playerCol := math2d.Collider{
        Bounds: p.Bounds(),
        Layer:  LayerPlayer,
        Mask:   LayerCoin,
    }
    for i, coin := range g.coins {
        coinCol := math2d.Collider{
            Bounds: coin,
            Layer:  LayerCoin,
            Mask:   LayerPlayer,
        }
        if math2d.ColliderOverlap(playerCol, coinCol) {
            g.score++
            g.coins = append(g.coins[:i], g.coins[i+1:]...)
            break
        }
    }
}
```
