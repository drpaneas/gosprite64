package math2d

// Layer is a bitmask identifying which collision group an object belongs to.
// Use powers of 2: const LayerPlayer Layer = 1 << 0, const LayerEnemy Layer = 1 << 1, etc.
type Layer uint32

const (
	LayerNone Layer = 0
	LayerAll  Layer = 0xFFFFFFFF
)

// Matches returns true if any bit in this layer overlaps with other.
func (l Layer) Matches(other Layer) bool {
	return l&other != 0
}

// Collider combines a bounding box with layer membership and a collision mask.
// Layer is what this object IS. Mask is what this object COLLIDES WITH.
type Collider struct {
	Bounds Rect
	Layer  Layer
	Mask   Layer
}

// ColliderOverlap returns true if two colliders overlap both spatially
// and in terms of layer filtering. The check is bidirectional: either
// collider's mask matching the other's layer is sufficient.
func ColliderOverlap(a, b Collider) bool {
	if !a.Mask.Matches(b.Layer) && !b.Mask.Matches(a.Layer) {
		return false
	}
	return a.Bounds.Overlaps(b.Bounds)
}
