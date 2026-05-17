package scene3d

import "github.com/drpaneas/gosprite64/math3d"

// RenderContext provides state during scene graph traversal.
type RenderContext struct {
	ViewMatrix       math3d.Mat4
	ProjectionMatrix math3d.Mat4
	CameraPosition   math3d.Vec3
	MatrixStack      []math3d.Mat4
}

// NewRenderContext creates a render context with default state.
func NewRenderContext() *RenderContext {
	return &RenderContext{
		MatrixStack: []math3d.Mat4{math3d.Identity()},
	}
}

// PushMatrix pushes a copy of the current matrix onto the stack.
func (rc *RenderContext) PushMatrix() {
	top := rc.MatrixStack[len(rc.MatrixStack)-1]
	rc.MatrixStack = append(rc.MatrixStack, top)
}

// PopMatrix removes the top matrix from the stack.
func (rc *RenderContext) PopMatrix() {
	if len(rc.MatrixStack) > 1 {
		rc.MatrixStack = rc.MatrixStack[:len(rc.MatrixStack)-1]
	}
}

// CurrentMatrix returns the top of the matrix stack.
func (rc *RenderContext) CurrentMatrix() math3d.Mat4 {
	return rc.MatrixStack[len(rc.MatrixStack)-1]
}

// MultiplyMatrix multiplies the top matrix by m.
func (rc *RenderContext) MultiplyMatrix(m math3d.Mat4) {
	top := len(rc.MatrixStack) - 1
	rc.MatrixStack[top] = rc.MatrixStack[top].Mul(m)
}

// Scene holds the root node and active camera for rendering.
type Scene struct {
	Root   *Node
	Camera *Node
}

// NewScene creates a scene with a root group node.
func NewScene() *Scene {
	return &Scene{
		Root: NewNode("root"),
	}
}

// Traverse walks the scene graph depth-first, calling the visitor for each visible node.
// The render context maintains the model-view matrix stack.
func (s *Scene) Traverse(ctx *RenderContext, visitor func(node *Node, ctx *RenderContext)) {
	if s.Root == nil {
		return
	}
	traverseNode(s.Root, ctx, visitor)
}

func traverseNode(node *Node, ctx *RenderContext, visitor func(*Node, *RenderContext)) {
	if !node.Visible {
		return
	}
	ctx.PushMatrix()
	ctx.MultiplyMatrix(node.LocalTransform())

	visitor(node, ctx)

	switch node.Type {
	case NodeLOD:
		if node.LOD != nil {
			pos := ctx.CameraPosition
			nodePos := node.Position
			dx := pos.X - nodePos.X
			dy := pos.Y - nodePos.Y
			dz := pos.Z - nodePos.Z
			dist := math3d.Vec3{X: dx, Y: dy, Z: dz}.Length()
			if child := node.LOD.Select(dist); child != nil {
				traverseNode(child, ctx, visitor)
			}
		}
	default:
		for _, child := range node.Children {
			traverseNode(child, ctx, visitor)
		}
	}

	ctx.PopMatrix()
}
