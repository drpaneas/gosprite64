package gosprite64

import (
	tilerender "github.com/drpaneas/gosprite64/internal/tile2d/render"
	"github.com/drpaneas/gosprite64/internal/tile2d/visibility"
)

type sceneRendererAdapter struct {
	renderer *tilerender.Renderer
	bridge   *sceneRenderBridge
}

func newSceneRendererAdapter(renderer *tilerender.Renderer, bridge *sceneRenderBridge) *sceneRendererAdapter {
	return &sceneRendererAdapter{renderer: renderer, bridge: bridge}
}

func (a *sceneRendererAdapter) configure() {
	if a == nil || a.renderer == nil || a.bridge == nil {
		return
	}
	a.renderer.SetHooks(tilerender.RenderHooks{
		Executor: a.bridge,
	})
}

func (a *sceneRendererAdapter) drawPreparedScene(scene tilerender.PreparedScene, cam Camera) tilerender.DrawStats {
	if a == nil || a.renderer == nil {
		return tilerender.DrawStats{}
	}
	return a.renderer.DrawPreparedScene(
		scene,
		visibility.Camera{
			X:      cam.X,
			Y:      cam.Y,
			Width:  cam.Width,
			Height: cam.Height,
		},
	)
}
