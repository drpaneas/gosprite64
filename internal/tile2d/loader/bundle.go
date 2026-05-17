package loader

import (
	"fmt"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
)

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

	if len(sheets) > 0 {
		refW, refH := sheets[0].TileWidth, sheets[0].TileHeight
		if refW == 0 || refH == 0 {
			return fmt.Errorf("scene validation: sheet 1 has zero tile size %dx%d", refW, refH)
		}
		for i := 1; i < len(sheets); i++ {
			if sheets[i].TileWidth == 0 || sheets[i].TileHeight == 0 {
				return fmt.Errorf("scene validation: sheet %d has zero tile size %dx%d", i+1, sheets[i].TileWidth, sheets[i].TileHeight)
			}
			if sheets[i].TileWidth != refW || sheets[i].TileHeight != refH {
				return fmt.Errorf(
					"scene validation: sheet %d tile size %dx%d differs from sheet 1 tile size %dx%d",
					i+1, sheets[i].TileWidth, sheets[i].TileHeight, refW, refH,
				)
			}
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
