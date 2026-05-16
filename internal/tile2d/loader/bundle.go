package loader

import (
	"fmt"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
)

const supportedTileWidth = 8
const supportedTileHeight = 8

func OpenBundle(path string, l Loader) (format.ParsedBundle, error) {
	raw, err := l.ReadAsset(path)
	if err != nil {
		return format.ParsedBundle{}, err
	}
	return format.ParseBundle(raw)
}

func LoadSheet(path string, l Loader) (format.ParsedSheet, error) {
	raw, err := l.ReadAsset(path)
	if err != nil {
		return format.ParsedSheet{}, err
	}
	return format.ParseSheet(raw)
}

func LoadMap(path string, l Loader) (format.ParsedMap, error) {
	raw, err := l.ReadAsset(path)
	if err != nil {
		return format.ParsedMap{}, err
	}
	return format.ParseMap(raw)
}

func LoadAnim(path string, l Loader) (format.ParsedAnim, error) {
	raw, err := l.ReadAsset(path)
	if err != nil {
		return format.ParsedAnim{}, err
	}
	return format.ParseAnim(raw)
}

func ValidateSceneAssets(m format.ParsedMap, sheets []format.ParsedSheet) error {
	if len(m.Layers) == 0 {
		return nil
	}

	for sheetIdx, sheet := range sheets {
		if int(sheet.TileWidth) != supportedTileWidth || int(sheet.TileHeight) != supportedTileHeight {
			return fmt.Errorf(
				"scene validation: sheet %d uses unsupported tile size %dx%d, want %dx%d",
				sheetIdx+1,
				sheet.TileWidth,
				sheet.TileHeight,
				supportedTileWidth,
				supportedTileHeight,
			)
		}
	}

	for layerIdx, layer := range m.Layers {
		sheetID := layer.SheetID
		if sheetID == 0 {
			sheetID = 1
		}
		if sheetID > uint16(len(sheets)) {
			return fmt.Errorf("scene validation: layer %d references missing sheet %d", layerIdx, sheetID)
		}

		sheet := sheets[sheetID-1]
		for cellIdx, tileID := range layer.Cells {
			if tileID == 0 {
				continue
			}
			if tileID > sheet.TileCount {
				return fmt.Errorf("scene validation: layer %d cell %d references tile %d beyond sheet %d tile count %d", layerIdx, cellIdx, tileID, sheetID, sheet.TileCount)
			}
		}
	}

	return nil
}
