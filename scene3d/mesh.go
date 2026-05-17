package scene3d

import (
	"github.com/drpaneas/gosprite64/gfx"
	"github.com/drpaneas/gosprite64/math3d"
)

// MeshData references a pre-built display list for a 3D model.
type MeshData struct {
	DisplayList *gfx.DisplayList
	// BoundingSphere for frustum culling: center (local) + radius.
	BoundsCenter math3d.Vec3
	BoundsRadius float32
}

// NewMeshNode creates a node that renders a display list.
func NewMeshNode(name string, dl *gfx.DisplayList) *Node {
	return &Node{
		Type:    NodeMesh,
		Name:    name,
		Visible: true,
		Scale:   math3d.Vec3{X: 1, Y: 1, Z: 1},
		Mesh:    &MeshData{DisplayList: dl},
	}
}
