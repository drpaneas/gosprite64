//go:build !noos

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
)

type repeatedFlag []string

func (r *repeatedFlag) String() string {
	return strings.Join(*r, ",")
}

func (r *repeatedFlag) Set(value string) error {
	*r = append(*r, value)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "mk2dbundle: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("mk2dbundle", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var out string
	var sheets, maps, anims repeatedFlag
	fs.Var(&sheets, "sheet", "Sheet asset path (repeatable)")
	fs.Var(&maps, "map", "Map asset path (repeatable)")
	fs.Var(&anims, "anim", "Animation asset path (repeatable)")
	fs.StringVar(&out, "out", "", "Output .bundle path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if out == "" {
		return fmt.Errorf("-out is required")
	}

	entries := make([]format.BundleEntry, 0, len(sheets)+len(maps)+len(anims))
	appendEntries := func(kind uint8, paths []string) {
		for _, path := range paths {
			entries = append(entries, format.BundleEntry{
				Kind: kind,
				Name: strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
				Path: path,
			})
		}
	}

	appendEntries(format.BundleKindSheet, sheets)
	appendEntries(format.BundleKindMap, maps)
	appendEntries(format.BundleKindAnim, anims)

	raw, err := format.BuildBundle(entries)
	if err != nil {
		return err
	}

	return os.WriteFile(out, raw, 0o644)
}
