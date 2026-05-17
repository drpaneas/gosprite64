package gosprite64

import "strings"

// TextAlign controls horizontal text alignment.
type TextAlign int

const (
	AlignLeft   TextAlign = iota
	AlignCenter
	AlignRight
)

// WrapText inserts newlines so no line exceeds maxWidth pixels.
// Breaks on spaces. Words wider than maxWidth are not broken.
func (f *Font) WrapText(text string, maxWidth int) string {
	if f == nil || len(text) == 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	spaceGlyph, hasSpace := f.GlyphFor(' ')
	spaceW := 0
	if hasSpace {
		spaceW = spaceGlyph.Advance
	} else {
		spaceW = 8
	}

	var lines []string
	currentLine := words[0]
	currentWidth, _ := f.MeasureText(currentLine)

	for _, word := range words[1:] {
		wordW, _ := f.MeasureText(word)
		if currentWidth+spaceW+wordW <= maxWidth {
			currentLine += " " + word
			currentWidth += spaceW + wordW
		} else {
			lines = append(lines, currentLine)
			currentLine = word
			currentWidth = wordW
		}
	}
	lines = append(lines, currentLine)

	return strings.Join(lines, "\n")
}
