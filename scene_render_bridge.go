package gosprite64

import (
	tilerender "github.com/drpaneas/gosprite64/internal/tile2d/render"
)

type sceneRenderBridge struct{}

func newSceneRenderBridge() *sceneRenderBridge {
	return &sceneRenderBridge{}
}

func (b *sceneRenderBridge) DrawPreparedRun(x, y, tileWidth, tileHeight int, run tilerender.PreparedRun) {
	if run.Count <= 0 {
		return
	}
	exec, ok := b.currentTexturedExecutor()
	if !ok || !exec.EnsurePrepared(run.Tile) {
		for i := 0; i < run.Count; i++ {
			drawLogicalImage(run.Tile.Source, x+(i*tileWidth), y)
		}
		return
	}
	exec.BlitRun(x, y, tileWidth, run.Count)
}

func (b *sceneRenderBridge) DrawPreparedTile(x, y, width, height int, tile tilerender.PreparedTile) {
	exec, ok := b.currentTexturedExecutor()
	if !ok || !exec.EnsurePrepared(tile) {
		drawLogicalImage(tile.Source, x, y)
		return
	}
	exec.BlitTile(x, y)
}

func (b *sceneRenderBridge) currentTexturedExecutor() (tilerender.TexturedExecutor, bool) {
	video := currentVideo()
	rt := currentTile()
	if video == nil || video.Framebuffer == nil || rt == nil {
		return tilerender.TexturedExecutor{}, false
	}
	return tilerender.TexturedExecutor{
		Framebuffer: video.Framebuffer,
		State:       &rt.textured,
	}, true
}
