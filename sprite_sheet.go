package gosprite64

import (
	"fmt"

	tileloader "github.com/drpaneas/gosprite64/internal/tile2d/loader"
)

type SpriteSheet struct {
	sheet *Sheet
}

func LoadSpriteSheet(path string) (*SpriteSheet, error) {
	parsed, err := tileloader.LoadSheet(path, cartLoader{})
	if err != nil {
		return nil, fmt.Errorf("load sprite sheet: %w", err)
	}
	if parsed.TileCount == 0 {
		return nil, fmt.Errorf("load sprite sheet: zero-frame sheet is invalid")
	}
	return &SpriteSheet{sheet: &Sheet{parsed: parsed}}, nil
}

func (s *SpriteSheet) FrameCount() int {
	if s == nil || s.sheet == nil || s.sheet.parsed.TileCount == 0 {
		return 0
	}
	return int(s.sheet.parsed.TileCount)
}

func (s *SpriteSheet) FrameWidth() int {
	if s == nil || s.sheet == nil || s.sheet.parsed.TileWidth == 0 {
		return 0
	}
	return int(s.sheet.parsed.TileWidth)
}

func (s *SpriteSheet) FrameHeight() int {
	if s == nil || s.sheet == nil || s.sheet.parsed.TileHeight == 0 {
		return 0
	}
	return int(s.sheet.parsed.TileHeight)
}
