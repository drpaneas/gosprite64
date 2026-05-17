package main

import (
	"fmt"
	"go/format"
	"sort"
	"strings"
	"unicode"
)

func titleCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func generateGlyphSource(fontName, pkgName string, glyphs map[rune]glyphInfo, lineHeight int) (string, error) {
	var b strings.Builder

	fmt.Fprintf(&b, "package %s\n\n", pkgName)
	fmt.Fprintf(&b, "import \"github.com/drpaneas/gosprite64\"\n\n")
	fmt.Fprintf(&b, "const %sLineHeight = %d\n\n", titleCase(fontName), lineHeight)
	fmt.Fprintf(&b, "var %sGlyphs = map[rune]gosprite64.Glyph{\n", titleCase(fontName))

	runes := make([]rune, 0, len(glyphs))
	for r := range glyphs {
		runes = append(runes, r)
	}
	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })

	for _, r := range runes {
		g := glyphs[r]
		fmt.Fprintf(&b, "\t'%s': gosprite64.Glyph{Frame: %d, Width: %d, Advance: %d},\n",
			escapeRune(r), g.Frame, g.Width, g.Advance)
	}

	fmt.Fprintf(&b, "}\n")

	formatted, err := format.Source([]byte(b.String()))
	if err != nil {
		return b.String(), nil
	}
	return string(formatted), nil
}

func escapeRune(r rune) string {
	switch r {
	case '\'':
		return "\\'"
	case '\\':
		return "\\\\"
	default:
		return string(r)
	}
}
