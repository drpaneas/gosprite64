//go:build !noos

package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "mk2dsheet: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("mk2dsheet", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var in, out string
	var tileWidth, tileHeight int
	fs.StringVar(&in, "in", "", "Input PNG path")
	fs.StringVar(&out, "out", "", "Output .sheet path")
	fs.IntVar(&tileWidth, "tile-width", 8, "Tile width in pixels")
	fs.IntVar(&tileHeight, "tile-height", 8, "Tile height in pixels")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if in == "" || out == "" {
		return fmt.Errorf("both -in and -out are required")
	}

	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return err
	}

	raw, err := format.BuildSheet(img, tileWidth, tileHeight)
	if err != nil {
		return err
	}

	return os.WriteFile(out, raw, 0o644)
}
