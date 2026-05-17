package gosprite64

import "testing"

func TestNewFontFromGlyphs(t *testing.T) {
	glyphs := map[rune]Glyph{
		'A': {Frame: 0, Width: 8, Advance: 9},
		'B': {Frame: 1, Width: 7, Advance: 8},
	}
	f := NewFont(nil, glyphs, 10)
	if f == nil {
		t.Fatal("NewFont should not return nil")
	}
	if f.LineHeight() != 10 {
		t.Fatalf("expected line height 10, got %d", f.LineHeight())
	}
}

func TestFontGlyphLookup(t *testing.T) {
	glyphs := map[rune]Glyph{
		'A': {Frame: 0, Width: 8, Advance: 9},
		'B': {Frame: 1, Width: 7, Advance: 8},
	}
	f := NewFont(nil, glyphs, 10)
	g, ok := f.GlyphFor('A')
	if !ok {
		t.Fatal("should find glyph for 'A'")
	}
	if g.Width != 8 || g.Advance != 9 {
		t.Fatalf("expected width=8 advance=9, got width=%d advance=%d", g.Width, g.Advance)
	}
}

func TestFontGlyphMissing(t *testing.T) {
	glyphs := map[rune]Glyph{
		'A': {Frame: 0, Width: 8, Advance: 9},
	}
	f := NewFont(nil, glyphs, 10)
	_, ok := f.GlyphFor('Z')
	if ok {
		t.Fatal("should not find glyph for 'Z' when not registered")
	}
}

func TestFontGlyphWithFallback(t *testing.T) {
	glyphs := map[rune]Glyph{
		'A': {Frame: 0, Width: 8, Advance: 9},
		'?': {Frame: 2, Width: 6, Advance: 7},
	}
	f := NewFont(nil, glyphs, 10)
	f.Fallback = '?'
	g, ok := f.GlyphFor('Z')
	if !ok {
		t.Fatal("should find fallback glyph")
	}
	if g.Frame != 2 {
		t.Fatalf("should use '?' glyph (frame 2), got frame %d", g.Frame)
	}
}

func TestMeasureText(t *testing.T) {
	glyphs := map[rune]Glyph{
		'H': {Frame: 0, Width: 8, Advance: 9},
		'i': {Frame: 1, Width: 4, Advance: 5},
	}
	f := NewFont(nil, glyphs, 10)
	w, h := f.MeasureText("Hi")
	expectedW := 9 + 5
	if w != expectedW {
		t.Fatalf("expected width %d, got %d", expectedW, w)
	}
	if h != 10 {
		t.Fatalf("expected height 10, got %d", h)
	}
}

func TestMeasureTextMultiline(t *testing.T) {
	glyphs := map[rune]Glyph{
		'A': {Frame: 0, Width: 8, Advance: 9},
		'B': {Frame: 1, Width: 8, Advance: 9},
	}
	f := NewFont(nil, glyphs, 10)
	_, h := f.MeasureText("A\nB")
	if h != 22 {
		t.Fatalf("expected height 22 (10 + 2 spacing + 10), got %d", h)
	}
}

func TestMeasureTextEmpty(t *testing.T) {
	f := NewFont(nil, map[rune]Glyph{}, 10)
	w, h := f.MeasureText("")
	if w != 0 || h != 0 {
		t.Fatalf("empty text should measure 0x0, got %dx%d", w, h)
	}
}

func TestMeasureNumber(t *testing.T) {
	glyphs := map[rune]Glyph{}
	for i := 0; i < 10; i++ {
		glyphs[rune('0'+i)] = Glyph{Frame: i, Width: 8, Advance: 9}
	}
	f := NewFont(nil, glyphs, 10)
	w, _ := f.MeasureText("12345")
	expectedW := 5 * 9
	if w != expectedW {
		t.Fatalf("expected width %d, got %d", expectedW, w)
	}
}

func TestFormatScore(t *testing.T) {
	s := FormatScore(1234, 8)
	if s != "00001234" {
		t.Fatalf("expected '00001234', got '%s'", s)
	}
}

func TestFormatScoreNoOverflow(t *testing.T) {
	s := FormatScore(123456789, 6)
	if s != "123456789" {
		t.Fatalf("expected '123456789' (no truncation), got '%s'", s)
	}
}

func TestFormatScoreZero(t *testing.T) {
	s := FormatScore(0, 4)
	if s != "0000" {
		t.Fatalf("expected '0000', got '%s'", s)
	}
}
