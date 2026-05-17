package scene3d

import "github.com/drpaneas/gosprite64/math3d"

// LODLevel defines a single level of detail with a maximum distance.
type LODLevel struct {
	MaxDistance float32
	Child      *Node
}

// LODData selects a child node based on distance from camera.
type LODData struct {
	Levels []LODLevel
}

// NewLODNode creates a level-of-detail node.
func NewLODNode(name string, levels []LODLevel) *Node {
	return &Node{
		Type:    NodeLOD,
		Name:    name,
		Visible: true,
		Scale:   math3d.Vec3{X: 1, Y: 1, Z: 1},
		LOD:     &LODData{Levels: levels},
	}
}

// Select returns the appropriate child for the given distance,
// or nil if no level matches.
func (l *LODData) Select(distance float32) *Node {
	for _, lv := range l.Levels {
		if distance <= lv.MaxDistance {
			return lv.Child
		}
	}
	return nil
}
