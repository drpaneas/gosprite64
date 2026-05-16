package render

import (
	"fmt"
	"image"

	"github.com/clktmr/n64/rcp/texture"
)

type Tileset struct {
	tiles []*texture.Texture
}

func NewTilesetFromAtlas(atlas image.Image, tileWidth, tileHeight int) (*Tileset, error) {
	if atlas == nil {
		return nil, fmt.Errorf("render: nil atlas")
	}
	if tileWidth <= 0 || tileHeight <= 0 {
		return nil, fmt.Errorf("render: invalid tile size %dx%d", tileWidth, tileHeight)
	}

	bounds := atlas.Bounds()
	if bounds.Dx()%tileWidth != 0 || bounds.Dy()%tileHeight != 0 {
		return nil, fmt.Errorf("render: atlas size %dx%d not divisible by tile size %dx%d", bounds.Dx(), bounds.Dy(), tileWidth, tileHeight)
	}

	cols := bounds.Dx() / tileWidth
	rows := bounds.Dy() / tileHeight
	tiles := make([]*texture.Texture, 0, cols*rows)
	for tileY := 0; tileY < rows; tileY++ {
		for tileX := 0; tileX < cols; tileX++ {
			rect := image.Rect(
				bounds.Min.X+tileX*tileWidth,
				bounds.Min.Y+tileY*tileHeight,
				bounds.Min.X+(tileX+1)*tileWidth,
				bounds.Min.Y+(tileY+1)*tileHeight,
			)
			subImager, ok := atlas.(interface{ SubImage(image.Rectangle) image.Image })
			if !ok {
				return nil, fmt.Errorf("render: atlas does not support subimages")
			}
			tiles = append(tiles, texture.NewTextureFromImage(subImager.SubImage(rect)))
		}
	}

	return &Tileset{tiles: tiles}, nil
}

func (t *Tileset) Tile(tileID uint16) *texture.Texture {
	if t == nil || tileID == 0 {
		return nil
	}
	index := int(tileID) - 1
	if index < 0 || index >= len(t.tiles) {
		return nil
	}
	return t.tiles[index]
}
