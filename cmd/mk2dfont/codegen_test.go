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

func TestGenerateGlyphSourceEscapeRune(t *testing.T) {
	glyphs := map[rune]glyphInfo{
		'\'': {Frame: 0, Width: 4, Advance: 4},
		'\\': {Frame: 1, Width: 6, Advance: 6},
	}
	src, err := generateGlyphSource("esc", "main", glyphs, 8)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(src, `'\''`) {
		t.Fatalf("should escape single quote, got:\n%s", src)
	}
	if !strings.Contains(src, `'\\'`) {
		t.Fatalf("should escape backslash, got:\n%s", src)
	}
}

func TestGenerateGlyphSourceEmpty(t *testing.T) {
	glyphs := map[rune]glyphInfo{}
	src, err := generateGlyphSource("empty", "main", glyphs, 10)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(src, "map[rune]gosprite64.Glyph{}") {
		t.Fatalf("empty glyph map should produce empty map literal, got:\n%s", src)
	}
}

func TestGenerateGlyphSourceWithOffsets(t *testing.T) {
	glyphs := map[rune]glyphInfo{
		'g': {Frame: 0, Width: 7, Advance: 8, OffsetX: 0, OffsetY: 2},
	}
	src, err := generateGlyphSource("off", "main", glyphs, 10)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(src, "OffsetY: 2") {
		t.Fatalf("should include OffsetY when non-zero, got:\n%s", src)
	}
}

func TestGenerateGlyphSourceNoOffsetsWhenZero(t *testing.T) {
	glyphs := map[rune]glyphInfo{
		'A': {Frame: 0, Width: 7, Advance: 7, OffsetX: 0, OffsetY: 0},
	}
	src, err := generateGlyphSource("noff", "main", glyphs, 10)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(src, "OffsetX") || strings.Contains(src, "OffsetY") {
		t.Fatalf("should omit OffsetX/OffsetY when both zero, got:\n%s", src)
	}
}
