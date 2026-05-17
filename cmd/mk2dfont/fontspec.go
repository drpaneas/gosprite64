package main

import (
	"encoding/json"
	"fmt"
)

type glyphEntry struct {
	Char  string `json:"char"`
	Width int    `json:"width"`
}

type fontSpec struct {
	CellWidth  int          `json:"cell_width"`
	CellHeight int          `json:"cell_height"`
	Chars      string       `json:"chars"`
	GlyphList  []glyphEntry `json:"glyphs"`
}

type glyphInfo struct {
	Frame   int
	Width   int
	Advance int
}

func parseFontSpec(data []byte) (*fontSpec, error) {
	var spec fontSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("parse font spec: %w", err)
	}
	if spec.CellWidth <= 0 || spec.CellHeight <= 0 {
		return nil, fmt.Errorf("parse font spec: cell_width and cell_height must be positive")
	}
	if spec.Chars != "" && len(spec.GlyphList) > 0 {
		return nil, fmt.Errorf("parse font spec: specify either chars or glyphs, not both")
	}
	if spec.Chars == "" && len(spec.GlyphList) == 0 {
		return nil, fmt.Errorf("parse font spec: must specify chars or glyphs")
	}
	return &spec, nil
}

func (s *fontSpec) BuildGlyphs() map[rune]glyphInfo {
	result := make(map[rune]glyphInfo)

	if s.Chars != "" {
		for i, ch := range s.Chars {
			result[ch] = glyphInfo{
				Frame:   i,
				Width:   s.CellWidth,
				Advance: s.CellWidth,
			}
		}
		return result
	}

	for i, entry := range s.GlyphList {
		runes := []rune(entry.Char)
		if len(runes) != 1 {
			continue
		}
		w := entry.Width
		if w <= 0 {
			w = s.CellWidth
		}
		result[runes[0]] = glyphInfo{
			Frame:   i,
			Width:   w,
			Advance: w,
		}
	}
	return result
}
