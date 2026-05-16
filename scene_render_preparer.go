package gosprite64

import (
	tilerender "github.com/drpaneas/gosprite64/internal/tile2d/render"
	"github.com/drpaneas/gosprite64/internal/tile2d/visibility"
)

type sceneRenderPreparer struct {
	scene *Scene
}

type sceneRenderLayerSource struct {
	Map     visibility.MapInfo
	SheetID uint16
	Tiles   [][]uint16
}

func newSceneRenderPreparer(scene *Scene) *sceneRenderPreparer {
	return &sceneRenderPreparer{scene: scene}
}

func (p *sceneRenderPreparer) buildScene() tilerender.PreparedScene {
	if p == nil || p.scene == nil || p.scene.gameMap == nil {
		return tilerender.PreparedScene{}
	}

	layers := p.scene.gameMap.renderLayers(p.primaryTileWidth(), p.primaryTileHeight())
	preparedLayers := make([]tilerender.PreparedLayer, 0, len(layers))
	for _, layer := range layers {
		prepared := p.prepareLayerTiles(layer)
		preparedLayers = append(preparedLayers, tilerender.PreparedLayer{
			Map:   layer.Map,
			Draws: p.prepareLayerDraws(prepared),
			Runs:  p.prepareLayerRuns(prepared),
		})
	}

	return tilerender.PreparedScene{
		Layers: preparedLayers,
	}
}

func (p *sceneRenderPreparer) primaryTileWidth() int {
	if p != nil && p.scene != nil && len(p.scene.sheets) > 0 && p.scene.sheets[0] != nil && p.scene.sheets[0].parsed.TileWidth > 0 {
		return int(p.scene.sheets[0].parsed.TileWidth)
	}
	return 8
}

func (p *sceneRenderPreparer) primaryTileHeight() int {
	if p != nil && p.scene != nil && len(p.scene.sheets) > 0 && p.scene.sheets[0] != nil && p.scene.sheets[0].parsed.TileHeight > 0 {
		return int(p.scene.sheets[0].parsed.TileHeight)
	}
	return 8
}

func nonZeroCount(prepared [][]tilerender.PreparedTile) int {
	n := 0
	for _, row := range prepared {
		for _, tile := range row {
			if tile.TileID != 0 {
				n++
			}
		}
	}
	return n
}

func (p *sceneRenderPreparer) prepareLayerTiles(layer sceneRenderLayerSource) [][]tilerender.PreparedTile {
	prepared := make([][]tilerender.PreparedTile, len(layer.Tiles))
	for y := range layer.Tiles {
		prepared[y] = make([]tilerender.PreparedTile, len(layer.Tiles[y]))
		for x, tileID := range layer.Tiles[y] {
			if tileID == 0 {
				continue
			}

			sheetID := layer.SheetID
			if sheetID == 0 {
				sheetID = 1
			}

			entry := tilerender.PreparedTile{
				TileID:  tileID,
				SheetID: sheetID,
			}
			if sheet := p.scene.SheetByID(sheetID); sheet != nil {
				entry.Source = sheet.tileImage(tileID)
			}
			prepared[y][x] = entry
		}
	}
	return prepared
}

func (p *sceneRenderPreparer) prepareLayerDraws(prepared [][]tilerender.PreparedTile) []tilerender.PreparedDraw {
	draws := make([]tilerender.PreparedDraw, 0, nonZeroCount(prepared))
	for y := range prepared {
		for x, tile := range prepared[y] {
			if tile.TileID == 0 {
				continue
			}
			draws = append(draws, tilerender.PreparedDraw{
				CellX: x,
				CellY: y,
				Tile:  tile,
			})
		}
	}
	return draws
}

func (p *sceneRenderPreparer) prepareLayerRuns(prepared [][]tilerender.PreparedTile) []tilerender.PreparedRun {
	runs := make([]tilerender.PreparedRun, 0, nonZeroCount(prepared))
	for y := range prepared {
		row := prepared[y]
		for x := 0; x < len(row); {
			tile := row[x]
			if tile.TileID == 0 {
				x++
				continue
			}

			run := tilerender.PreparedRun{
				CellX: x,
				CellY: y,
				Count: 1,
				Tile:  tile,
			}
			for next := x + 1; next < len(row); next++ {
				candidate := row[next]
				if candidate.TileID != tile.TileID || candidate.SheetID != tile.SheetID || candidate.Source != tile.Source {
					break
				}
				run.Count++
				x = next
			}
			runs = append(runs, run)
			x++
		}
	}
	return runs
}
