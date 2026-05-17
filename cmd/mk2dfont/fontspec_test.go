package main

import "testing"

func TestParseFontSpecFixedWidth(t *testing.T) {
	json := `{"cell_width": 8, "cell_height": 10, "chars": "ABCD"}`
	spec, err := parseFontSpec([]byte(json))
	if err != nil {
		t.Fatal(err)
	}
	if spec.CellWidth != 8 || spec.CellHeight != 10 {
		t.Fatalf("expected 8x10, got %dx%d", spec.CellWidth, spec.CellHeight)
	}
	glyphs := spec.BuildGlyphs()
	if len(glyphs) != 4 {
		t.Fatalf("expected 4 glyphs, got %d", len(glyphs))
	}
	if glyphs['A'].Frame != 0 || glyphs['A'].Advance != 8 {
		t.Fatalf("A: expected frame=0 advance=8, got %+v", glyphs['A'])
	}
	if glyphs['D'].Frame != 3 {
		t.Fatalf("D: expected frame=3, got %d", glyphs['D'].Frame)
	}
}

func TestParseFontSpecVariableWidth(t *testing.T) {
	json := `{
		"cell_width": 8,
		"cell_height": 10,
		"glyphs": [
			{"char": "A", "width": 7},
			{"char": "B", "width": 6},
			{"char": " ", "width": 4}
		]
	}`
	spec, err := parseFontSpec([]byte(json))
	if err != nil {
		t.Fatal(err)
	}
	glyphs := spec.BuildGlyphs()
	if len(glyphs) != 3 {
		t.Fatalf("expected 3 glyphs, got %d", len(glyphs))
	}
	if glyphs['A'].Width != 7 || glyphs['A'].Advance != 7 {
		t.Fatalf("A: expected width=7 advance=7, got %+v", glyphs['A'])
	}
	if glyphs[' '].Width != 4 {
		t.Fatalf("space: expected width=4, got %d", glyphs[' '].Width)
	}
}

func TestParseFontSpecInvalid(t *testing.T) {
	_, err := parseFontSpec([]byte(`{}`))
	if err == nil {
		t.Fatal("empty spec should fail validation")
	}
}

func TestParseFontSpecBothCharsAndGlyphs(t *testing.T) {
	json := `{"cell_width": 8, "cell_height": 10, "chars": "AB", "glyphs": [{"char": "C", "width": 5}]}`
	_, err := parseFontSpec([]byte(json))
	if err == nil {
		t.Fatal("specifying both chars and glyphs should fail")
	}
}

func TestBuildGlyphsSkipsMultiRuneEntry(t *testing.T) {
	json := `{
		"cell_width": 8,
		"cell_height": 10,
		"glyphs": [
			{"char": "A", "width": 7},
			{"char": "AB", "width": 5},
			{"char": "C", "width": 6}
		]
	}`
	spec, err := parseFontSpec([]byte(json))
	if err != nil {
		t.Fatal(err)
	}
	glyphs := spec.BuildGlyphs()
	if len(glyphs) != 2 {
		t.Fatalf("expected 2 glyphs (AB skipped), got %d", len(glyphs))
	}
	if _, ok := glyphs['A']; !ok {
		t.Fatal("A should be present")
	}
	if _, ok := glyphs['C']; !ok {
		t.Fatal("C should be present")
	}
}

func TestBuildGlyphsWithOffsets(t *testing.T) {
	json := `{
		"cell_width": 8,
		"cell_height": 10,
		"glyphs": [
			{"char": "g", "width": 7, "offset_x": 0, "offset_y": 2}
		]
	}`
	spec, err := parseFontSpec([]byte(json))
	if err != nil {
		t.Fatal(err)
	}
	glyphs := spec.BuildGlyphs()
	if glyphs['g'].OffsetY != 2 {
		t.Fatalf("expected OffsetY=2, got %d", glyphs['g'].OffsetY)
	}
}
