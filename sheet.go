package gosprite64

import (
	"image"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
	tilerender "github.com/drpaneas/gosprite64/internal/tile2d/render"
)

type Sheet struct {
	parsed  format.ParsedSheet
	tileset *tilerender.Tileset
}

type SheetInfo struct {
	TileWidth   int
	TileHeight  int
	TileCount   int
	AtlasWidth  int
	AtlasHeight int
}

func (s *Sheet) Info() SheetInfo {
	if s == nil {
		return SheetInfo{}
	}
	return SheetInfo{
		TileWidth:   int(s.parsed.TileWidth),
		TileHeight:  int(s.parsed.TileHeight),
		TileCount:   int(s.parsed.TileCount),
		AtlasWidth:  int(s.parsed.AtlasWidth),
		AtlasHeight: int(s.parsed.AtlasHeight),
	}
}

func (s *Sheet) Tile(tileID uint16) image.Image {
	return s.tileImage(tileID)
}

func (s *Sheet) atlasImage() image.Image {
	if s == nil {
		return nil
	}
	if len(s.parsed.Pixels) == 0 {
		return nil
	}

	nrgba := &image.NRGBA{
		Pix:    append([]byte(nil), s.parsed.Pixels...),
		Stride: int(s.parsed.AtlasWidth) * 4,
		Rect:   image.Rect(0, 0, int(s.parsed.AtlasWidth), int(s.parsed.AtlasHeight)),
	}
	if allOpaque(nrgba.Pix) {
		rgba := image.NewRGBA(nrgba.Rect)
		for y := 0; y < nrgba.Rect.Dy(); y++ {
			for x := 0; x < nrgba.Rect.Dx(); x++ {
				rgba.Set(x, y, nrgba.At(x, y))
			}
		}
		return rgba
	}
	return nrgba
}

func (s *Sheet) tileImage(tileID uint16) image.Image {
	if s == nil || tileID == 0 {
		return nil
	}
	if s.tileset == nil {
		atlas := s.atlasImage()
		if atlas == nil || s.parsed.TileWidth == 0 || s.parsed.TileHeight == 0 {
			return nil
		}
		tileset, err := tilerender.NewTilesetFromAtlas(atlas, int(s.parsed.TileWidth), int(s.parsed.TileHeight))
		if err != nil {
			return nil
		}
		s.tileset = tileset
	}
	return s.tileset.Tile(tileID)
}

func allOpaque(pix []byte) bool {
	for i := 3; i < len(pix); i += 4 {
		if pix[i] != 0xFF {
			return false
		}
	}
	return true
}
