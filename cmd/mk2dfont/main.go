package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	pngPath := flag.String("png", "", "path to font atlas PNG")
	specPath := flag.String("spec", "", "path to font spec JSON")
	outSheet := flag.String("out-sheet", "", "output path for .sheet file")
	outGo := flag.String("out-go", "", "output path for generated .go file")
	fontName := flag.String("name", "font", "font name (used in generated variable names)")
	pkgName := flag.String("pkg", "main", "Go package name for generated file")
	flag.Parse()

	if *pngPath == "" || *specPath == "" || *outSheet == "" || *outGo == "" {
		fmt.Fprintln(os.Stderr, "usage: mk2dfont -png font.png -spec font.json -out-sheet font.sheet -out-go glyphs.go")
		os.Exit(1)
	}

	specData, err := os.ReadFile(*specPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read spec: %v\n", err)
		os.Exit(1)
	}

	spec, err := parseFontSpec(specData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	outDir := filepath.Dir(*outSheet)
	if outDir != "" && outDir != "." {
		os.MkdirAll(outDir, 0o755)
	}

	mk2dsheet := exec.Command("go", "run",
		"github.com/drpaneas/gosprite64/cmd/mk2dsheet",
		"-in", *pngPath,
		"-out", *outSheet,
		"-tile-width", fmt.Sprint(spec.CellWidth),
		"-tile-height", fmt.Sprint(spec.CellHeight),
	)
	mk2dsheet.Stdout = os.Stdout
	mk2dsheet.Stderr = os.Stderr
	if err := mk2dsheet.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "mk2dsheet failed: %v\n", err)
		os.Exit(1)
	}

	glyphs := spec.BuildGlyphs()
	src, err := generateGlyphSource(*fontName, *pkgName, glyphs, spec.CellHeight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "codegen failed: %v\n", err)
		os.Exit(1)
	}

	goDir := filepath.Dir(*outGo)
	if goDir != "" && goDir != "." {
		os.MkdirAll(goDir, 0o755)
	}

	if err := os.WriteFile(*outGo, []byte(src), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write go file: %v\n", err)
		os.Exit(1)
	}

	charCount := len(glyphs)
	widthType := "fixed"
	if len(spec.GlyphList) > 0 {
		widthType = "variable"
	}
	fmt.Printf("mk2dfont: %d %s-width glyphs -> %s + %s\n", charCount, widthType, *outSheet, *outGo)
}
