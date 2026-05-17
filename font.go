package gosprite64

import (
	"fmt"
	"strings"
)

// Glyph describes a single character in a font atlas.
type Glyph struct {
	Frame   int
	Width   int
	Advance int
	OffsetX int
	OffsetY int
}

// Font maps runes to sprite sheet frames with per-glyph metrics.
type Font struct {
	sheet      *SpriteSheet
	glyphs     map[rune]Glyph
	lineHeight int
	Fallback   rune
	Spacing    int
}

// NewFont creates a font from a sprite sheet and glyph map.
// lineHeight is the pixel distance between lines.
func NewFont(sheet *SpriteSheet, glyphs map[rune]Glyph, lineHeight int) *Font {
	return &Font{
		sheet:      sheet,
		glyphs:     glyphs,
		lineHeight: lineHeight,
		Spacing:    2,
	}
}

// LineHeight returns the vertical distance between baselines.
func (f *Font) LineHeight() int {
	if f == nil {
		return 0
	}
	return f.lineHeight
}

// GlyphFor returns the glyph for a rune. If the rune is not in the font
// and Fallback is set, returns the fallback glyph.
func (f *Font) GlyphFor(r rune) (Glyph, bool) {
	if f == nil || f.glyphs == nil {
		return Glyph{}, false
	}
	g, ok := f.glyphs[r]
	if ok {
		return g, true
	}
	if f.Fallback != 0 {
		g, ok = f.glyphs[f.Fallback]
		return g, ok
	}
	return Glyph{}, false
}

// MeasureText returns the pixel dimensions of the rendered text.
// Supports newlines for multiline measurement.
func (f *Font) MeasureText(text string) (width int, height int) {
	if f == nil || len(text) == 0 {
		return 0, 0
	}

	lines := strings.Split(text, "\n")
	maxWidth := 0

	for _, line := range lines {
		lineWidth := 0
		for _, r := range line {
			g, ok := f.GlyphFor(r)
			if !ok {
				continue
			}
			lineWidth += g.Advance
		}
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	totalHeight := len(lines)*f.lineHeight + (len(lines)-1)*f.Spacing
	return maxWidth, totalHeight
}

// DrawTextEx draws text at (x, y) using this font.
// Each character is drawn as a sprite frame from the font's sheet.
func (f *Font) DrawTextEx(text string, x, y int, align TextAlign) {
	if f == nil || f.sheet == nil || len(text) == 0 {
		return
	}

	lines := strings.Split(text, "\n")
	totalW, _ := f.MeasureText(text)

	curY := y
	for _, line := range lines {
		lineW := 0
		for _, r := range line {
			g, ok := f.GlyphFor(r)
			if ok {
				lineW += g.Advance
			}
		}

		var curX int
		switch align {
		case AlignCenter:
			curX = x + (totalW-lineW)/2
		case AlignRight:
			curX = x + totalW - lineW
		default:
			curX = x
		}

		for _, r := range line {
			g, ok := f.GlyphFor(r)
			if !ok {
				continue
			}
			DrawSprite(f.sheet, g.Frame, float32(curX+g.OffsetX), float32(curY+g.OffsetY))
			curX += g.Advance
		}

		curY += f.lineHeight + f.Spacing
	}
}

// FormatScore formats an integer score with leading zeros to the given width.
// If the number has more digits than width, the full number is returned.
func FormatScore(score int, width int) string {
	if score < 0 {
		score = 0
	}
	s := fmt.Sprintf("%d", score)
	if len(s) >= width {
		return s
	}
	return strings.Repeat("0", width-len(s)) + s
}
