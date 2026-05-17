//go:build n64

package scene3d

import (
	"github.com/drpaneas/gosprite64/gfx"
	"github.com/drpaneas/gosprite64/math3d"
)

func DrawScene(scene *Scene) {
	if scene == nil || scene.Root == nil {
		return
	}

	ctx := NewRenderContext()

	if scene.Camera != nil && scene.Camera.Camera != nil {
		cam := scene.Camera.Camera
		proj, _ := cam.ProjectionMatrix()
		ctx.ProjectionMatrix = proj

		pos := scene.Camera.Position
		ctx.CameraPosition = pos
		ctx.ViewMatrix = math3d.LookAt(
			pos.X, pos.Y, pos.Z,
			pos.X+scene.Camera.Rotation.X,
			pos.Y+scene.Camera.Rotation.Y,
			pos.Z+scene.Camera.Rotation.Z,
			0, 1, 0,
		)
		ctx.MatrixStack[0] = ctx.ViewMatrix
	}

	scene.Traverse(ctx, func(node *Node, rc *RenderContext) {
		if node.Type != NodeMesh {
			return
		}
		if node.Mesh == nil || node.Mesh.DisplayList == nil {
			return
		}
		if node.RenderFn != nil {
			node.RenderFn(rc)
			return
		}
		gfx.Execute(node.Mesh.DisplayList)
	})

	gfx.Flush()
}
