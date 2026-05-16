//go:build !noos

package main

import (
	"fmt"
	"os"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "dump2dasset: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: dump2dasset <asset-file>")
	}

	raw, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	if len(raw) < 4 {
		return fmt.Errorf("asset too short")
	}

	switch string(raw[:4]) {
	case "SHT2":
		sheet, err := format.ParseSheet(raw)
		if err != nil {
			return err
		}
		fmt.Printf("sheet %dx%d tiles=%d palette=%d\n", sheet.TileWidth, sheet.TileHeight, sheet.TileCount, sheet.PaletteEntries)
	case "MAP2":
		m, err := format.ParseMap(raw)
		if err != nil {
			return err
		}
		fmt.Printf("map %dx%d layers=%d cell_bits=%d\n", m.Width, m.Height, m.LayerCount, m.CellBits)
	case "ANM2":
		anim, err := format.ParseAnim(raw)
		if err != nil {
			return err
		}
		fmt.Printf("anim clips=%d\n", len(anim.Clips))
	case "BND2":
		bundle, err := format.ParseBundle(raw)
		if err != nil {
			return err
		}
		fmt.Printf("bundle entries=%d\n", len(bundle.Entries))
	default:
		return fmt.Errorf("unknown magic %q", string(raw[:4]))
	}

	return nil
}
