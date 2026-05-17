package gosprite64

import (
	"github.com/drpaneas/gosprite64/internal/tile2d/format"
	"github.com/drpaneas/gosprite64/internal/tile2d/visibility"
)

type Map struct {
	parsed          format.ParsedMap
	cachedLayerInfo []MapLayerInfo
	tileW           int
	tileH           int
}

type MapLayerInfo struct {
	SheetID      uint16
	NonZeroTiles int
}

func newMap(parsed format.ParsedMap) *Map {
	m := &Map{parsed: parsed}
	m.cachedLayerInfo = make([]MapLayerInfo, len(parsed.Layers))
	for i, layer := range parsed.Layers {
		info := MapLayerInfo{SheetID: layer.SheetID}
		if info.SheetID == 0 {
			info.SheetID = 1
		}
		for _, tile := range layer.Cells {
			if tile != 0 {
				info.NonZeroTiles++
			}
		}
		m.cachedLayerInfo[i] = info
	}
	return m
}

func (m *Map) Width() int {
	if m == nil {
		return 0
	}
	return int(m.parsed.Width)
}

func (m *Map) Height() int {
	if m == nil {
		return 0
	}
	return int(m.parsed.Height)
}

func (m *Map) TileWidth() int {
	if m == nil || m.tileW <= 0 {
		return 8
	}
	return m.tileW
}

func (m *Map) TileHeight() int {
	if m == nil || m.tileH <= 0 {
		return 8
	}
	return m.tileH
}

func (m *Map) LayerCount() int {
	if m == nil {
		return 0
	}
	return len(m.parsed.Layers)
}

func (m *Map) LayerInfo(layer int) (MapLayerInfo, bool) {
	if m == nil || layer < 0 || layer >= len(m.cachedLayerInfo) {
		return MapLayerInfo{}, false
	}
	return m.cachedLayerInfo[layer], true
}

func (m *Map) LayerSheetID(layer int) (uint16, bool) {
	info, ok := m.LayerInfo(layer)
	if !ok {
		return 0, false
	}
	return info.SheetID, true
}

func (m *Map) TileAt(layer, x, y int) (uint16, bool) {
	if m == nil || layer < 0 || layer >= len(m.parsed.Layers) {
		return 0, false
	}
	if x < 0 || x >= m.Width() || y < 0 || y >= m.Height() {
		return 0, false
	}
	idx := y*m.Width() + x
	return m.parsed.Layers[layer].Cells[idx], true
}

func (m *Map) PixelWidth() int {
	return m.Width() * m.TileWidth()
}

func (m *Map) PixelHeight() int {
	return m.Height() * m.TileHeight()
}

func (m *Map) renderLayers(tileWidth, tileHeight int) []sceneRenderLayerSource {
	if tileWidth <= 0 {
		tileWidth = 8
	}
	if tileHeight <= 0 {
		tileHeight = 8
	}

	w, h := m.Width(), m.Height()
	mapInfo := visibility.MapInfo{
		Width:      w,
		Height:     h,
		TileWidth:  tileWidth,
		TileHeight: tileHeight,
	}
	if len(m.parsed.Layers) == 0 {
		return []sceneRenderLayerSource{{
			Map:   mapInfo,
			Tiles: make([][]uint16, h),
		}}
	}

	layers := make([]sceneRenderLayerSource, 0, len(m.parsed.Layers))
	for _, parsedLayer := range m.parsed.Layers {
		layer := sceneRenderLayerSource{
			Map:     mapInfo,
			SheetID: parsedLayer.SheetID,
			Tiles:   make([][]uint16, h),
		}
		for y := range layer.Tiles {
			start := y * w
			end := start + w
			layer.Tiles[y] = append([]uint16(nil), parsedLayer.Cells[start:end]...)
		}
		layers = append(layers, layer)
	}
	return layers
}
