package scene3d

import "github.com/drpaneas/gosprite64/math3d"

// NodeType identifies the kind of scene graph node.
type NodeType int

const (
	NodeTransform NodeType = iota
	NodeCamera
	NodeMesh
	NodeBillboard
	NodeRenderFunc
	NodeLOD
	NodeGroup
)

// Node is a single element in the 3D scene graph.
type Node struct {
	Type     NodeType
	Name     string
	Children []*Node
	Visible  bool

	// Transform (local space)
	Position math3d.Vec3
	Rotation math3d.Vec3 // euler angles in degrees
	Scale    math3d.Vec3

	// Type-specific data
	Camera   *CameraData
	Mesh     *MeshData
	LOD      *LODData
	RenderFn RenderFunc
}

// RenderFunc is a callback for dynamic display list generation.
type RenderFunc func(ctx *RenderContext)

// NewNode creates a transform group node.
func NewNode(name string) *Node {
	return &Node{
		Type:    NodeGroup,
		Name:    name,
		Visible: true,
		Scale:   math3d.Vec3{X: 1, Y: 1, Z: 1},
	}
}

// AddChild appends a child node.
func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

// LocalTransform computes the local transformation matrix from position, rotation, and scale.
func (n *Node) LocalTransform() math3d.Mat4 {
	t := math3d.Translate(n.Position.X, n.Position.Y, n.Position.Z)
	rx := math3d.Rotate(n.Rotation.X, 1, 0, 0)
	ry := math3d.Rotate(n.Rotation.Y, 0, 1, 0)
	rz := math3d.Rotate(n.Rotation.Z, 0, 0, 1)
	s := math3d.Scale(n.Scale.X, n.Scale.Y, n.Scale.Z)
	return t.Mul(rz.Mul(ry.Mul(rx))).Mul(s)
}
