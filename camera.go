package gosprite64

import (
	"github.com/drpaneas/gosprite64/internal/rendergeom"
	"github.com/drpaneas/gosprite64/math2d"
)

// Camera defines the visible region of the game world.
type Camera struct {
	X, Y          int
	Width, Height int

	Zoom float32

	FollowTarget *math2d.Vec2
	FollowSpeed  float32

	Bounds *math2d.Rect

	trauma  float32
	shakeRng *math2d.Rand
}

func newDefaultCamera() *Camera {
	bounds := rendergeom.LogicalBounds()
	return &Camera{
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
	}
}

// EffectiveZoom returns the zoom level, defaulting to 1 if unset.
func (c *Camera) EffectiveZoom() float32 {
	if c == nil || c.Zoom == 0 {
		return 1
	}
	return c.Zoom
}

// WorldToScreen converts world coordinates to screen coordinates,
// accounting for camera position and zoom.
func (c *Camera) WorldToScreen(worldX, worldY float32) (float32, float32) {
	if c == nil {
		return worldX, worldY
	}
	z := c.EffectiveZoom()
	return (worldX - float32(c.X)) * z, (worldY - float32(c.Y)) * z
}

// UpdateFollow moves the camera toward FollowTarget by FollowSpeed (0..1).
// Speed of 1.0 snaps instantly. Speed of 0.1 gives smooth lerp.
// Centers the target in the viewport.
func (c *Camera) UpdateFollow() {
	if c == nil || c.FollowTarget == nil {
		return
	}
	targetX := c.FollowTarget.X - float32(c.Width)/2
	targetY := c.FollowTarget.Y - float32(c.Height)/2

	speed := c.FollowSpeed
	if speed <= 0 {
		speed = 1
	}
	if speed > 1 {
		speed = 1
	}

	newX := math2d.Lerp(float32(c.X), targetX, speed)
	newY := math2d.Lerp(float32(c.Y), targetY, speed)
	c.X = int(newX)
	c.Y = int(newY)
}

// ClampToBounds restricts the camera position to stay within Bounds.
// No-op if Bounds is nil.
func (c *Camera) ClampToBounds() {
	if c == nil || c.Bounds == nil {
		return
	}
	minX := int(c.Bounds.X)
	minY := int(c.Bounds.Y)
	maxX := int(c.Bounds.X+c.Bounds.W) - c.Width
	maxY := int(c.Bounds.Y+c.Bounds.H) - c.Height

	if maxX < minX {
		maxX = minX
	}
	if maxY < minY {
		maxY = minY
	}

	if c.X < minX {
		c.X = minX
	}
	if c.X > maxX {
		c.X = maxX
	}
	if c.Y < minY {
		c.Y = minY
	}
	if c.Y > maxY {
		c.Y = maxY
	}
}

// AddTrauma adds screen shake intensity (0..1). Multiple hits accumulate
// up to 1.0. Shake magnitude is trauma squared for a more natural feel.
func (c *Camera) AddTrauma(amount float32) {
	if c == nil {
		return
	}
	c.trauma += amount
	if c.trauma > 1.0 {
		c.trauma = 1.0
	}
	if c.shakeRng == nil {
		c.shakeRng = math2d.NewRand(12345)
	}
}

// UpdateShake decays trauma each frame. Call once per Update().
func (c *Camera) UpdateShake() {
	if c == nil {
		return
	}
	c.trauma -= 1.0 / 60.0
	if c.trauma < 0 {
		c.trauma = 0
	}
}

// ShakeOffset returns the current frame's shake displacement.
// Apply this to your draw offset: drawX = cam.X + shakeX.
func (c *Camera) ShakeOffset() (int, int) {
	if c == nil || c.trauma <= 0 {
		return 0, 0
	}
	if c.shakeRng == nil {
		c.shakeRng = math2d.NewRand(12345)
	}
	magnitude := c.trauma * c.trauma
	maxOffset := float32(8)
	ox := (c.shakeRng.Float32()*2 - 1) * maxOffset * magnitude
	oy := (c.shakeRng.Float32()*2 - 1) * maxOffset * magnitude
	return int(ox), int(oy)
}
