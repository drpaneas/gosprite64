package gosprite64

import (
	"strings"
	"testing"
)

func TestTextAlignValues(t *testing.T) {
	if AlignLeft != 0 {
		t.Fatal("AlignLeft should be 0")
	}
	if AlignCenter != 1 {
		t.Fatal("AlignCenter should be 1")
	}
	if AlignRight != 2 {
		t.Fatal("AlignRight should be 2")
	}
}

func TestWrapText(t *testing.T) {
	glyphs := map[rune]Glyph{}
	for _, r := range "abcdefghijklmnopqrstuvwxyz " {
		glyphs[r] = Glyph{Frame: 0, Width: 8, Advance: 8}
	}
	f := NewFont(nil, glyphs, 10)

	wrapped := f.WrapText("hello world foo bar", 50)
	lines := strings.Split(wrapped, "\n")
	for _, line := range lines {
		w, _ := f.MeasureText(line)
		if w > 50 {
			t.Fatalf("line %q is %d px wide, exceeds limit of 50", line, w)
		}
	}
}

func TestWrapTextSingleWord(t *testing.T) {
	glyphs := map[rune]Glyph{
		'a': {Frame: 0, Width: 8, Advance: 8},
	}
	f := NewFont(nil, glyphs, 10)
	wrapped := f.WrapText("aaaaaa", 20)
	if wrapped != "aaaaaa" {
		t.Fatalf("single word wider than limit should not break, got %q", wrapped)
	}
}

func TestWrapTextEmpty(t *testing.T) {
	f := NewFont(nil, map[rune]Glyph{}, 10)
	wrapped := f.WrapText("", 100)
	if wrapped != "" {
		t.Fatalf("empty input should return empty, got %q", wrapped)
	}
}
