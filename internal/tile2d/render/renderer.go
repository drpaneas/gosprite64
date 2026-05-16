package render

import (
	"image"

	"github.com/clktmr/n64/rcp/texture"
	"github.com/drpaneas/gosprite64/internal/tile2d/visibility"
)

type ExecutionBridge interface {
	DrawPreparedRun(x, y, tileWidth, tileHeight int, run PreparedRun)
	DrawPreparedTile(x, y, width, height int, tile PreparedTile)
}

type RenderHooks struct {
	Executor ExecutionBridge
}

type DrawStats struct {
	VisibleTiles int
	Uploads      int
}

type Renderer struct {
	hooks RenderHooks
}

type PreparedScene struct {
	Layers []PreparedLayer
}

type Scene = PreparedScene

type PreparedLayer struct {
	Map   visibility.MapInfo
	Draws []PreparedDraw
	Runs  []PreparedRun
}

type PreparedTile struct {
	TileID  uint16
	SheetID uint16
	Source  image.Image
}

type PreparedDraw struct {
	CellX int
	CellY int
	Tile  PreparedTile
}

type PreparedRun struct {
	CellX int
	CellY int
	Count int
	Tile  PreparedTile
}

func NewRenderer(hooks RenderHooks) *Renderer {
	return &Renderer{hooks: hooks}
}

func (r *Renderer) SetHooks(hooks RenderHooks) {
	if r == nil {
		return
	}
	r.hooks = hooks
}

func (r *Renderer) DrawPreparedScene(scene PreparedScene, cam visibility.Camera) DrawStats {
	var stats DrawStats
	var lastSource *texture.Texture
	havePreparedSource := false

	for _, layer := range scene.Layers {
		bounds := visibility.VisibleCellBounds(cam, layer.Map)
		if len(layer.Runs) > 0 {
			for _, run := range layer.Runs {
				if run.CellY < bounds.MinY || run.CellY >= bounds.MaxY {
					continue
				}
				runMinX := run.CellX
				runMaxX := run.CellX + run.Count
				if runMaxX <= bounds.MinX || runMinX >= bounds.MaxX {
					continue
				}
				if run.Tile.TileID == 0 {
					continue
				}
				if src, ok := (TexturedExecutor{}).SourceTexture(run.Tile); ok {
					if !havePreparedSource || src != lastSource {
						stats.Uploads++
						lastSource = src
						havePreparedSource = true
					}
				}

				startX := max(runMinX, bounds.MinX)
				endX := min(runMaxX, bounds.MaxX)
				clippedRun := run
				clippedRun.CellX = startX
				clippedRun.Count = endX - startX
				if r.hooks.Executor != nil {
					r.hooks.Executor.DrawPreparedRun(
						(startX*layer.Map.TileWidth)-cam.X,
						(run.CellY*layer.Map.TileHeight)-cam.Y,
						layer.Map.TileWidth,
						layer.Map.TileHeight,
						clippedRun,
					)
				}
				stats.VisibleTiles += clippedRun.Count
			}
			continue
		}
		if len(layer.Draws) > 0 {
			for _, draw := range layer.Draws {
				if draw.CellX < bounds.MinX || draw.CellX >= bounds.MaxX || draw.CellY < bounds.MinY || draw.CellY >= bounds.MaxY {
					continue
				}
				if draw.Tile.TileID == 0 {
					continue
				}
				if src, ok := (TexturedExecutor{}).SourceTexture(draw.Tile); ok {
					if !havePreparedSource || src != lastSource {
						stats.Uploads++
						lastSource = src
						havePreparedSource = true
					}
				}
				if r.hooks.Executor != nil {
					r.hooks.Executor.DrawPreparedTile(
						(draw.CellX*layer.Map.TileWidth)-cam.X,
						(draw.CellY*layer.Map.TileHeight)-cam.Y,
						layer.Map.TileWidth,
						layer.Map.TileHeight,
						draw.Tile,
					)
				}
				stats.VisibleTiles++
			}
			continue
		}
	}

	return stats
}

func (r *Renderer) DrawTileLayer(scene PreparedScene, cam visibility.Camera) DrawStats {
	return r.DrawPreparedScene(scene, cam)
}
