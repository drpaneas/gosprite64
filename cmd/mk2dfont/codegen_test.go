package main

import (
	"strings"
	"testing"
)

func TestGenerateGlyphSource(t *testing.T) {
	glyphs := map[rune]glyphInfo{
		'A': {Frame: 0, Width: 7, Advance: 7},
		'B': {Frame: 1, Width: 6, Advance: 6},
	}
	src, err := generateGlyphSource("myfont", "main", glyphs, 10)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(src, "package main") {
		t.Fatal("should contain package declaration")
	}
	if !strings.Contains(src, "gosprite64.Glyph{Frame: 0") {
		t.Fatalf("should contain glyph literal for A, got:\n%s", src)
	}
	if !strings.Contains(src, "'A':") || !strings.Contains(src, "'B':") {
		t.Fatalf("should contain character keys, got:\n%s", src)
	}
	if !strings.Contains(src, "LineHeight = 10") {
		t.Fatalf("should contain line height constant, got:\n%s", src)
	}
}

func TestGenerateGlyphSourceFormats(t *testing.T) {
	glyphs := map[rune]glyphInfo{
		'X': {Frame: 0, Width: 8, Advance: 8},
	}
	src, err := generateGlyphSource("test", "fonts", glyphs, 8)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(src, "package fonts") {
		t.Fatal("should use specified package name")
	}
}
