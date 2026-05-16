//go:build !noos

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "mk2danim: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("mk2danim", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var in, out string
	fs.StringVar(&in, "in", "", "Input JSON path")
	fs.StringVar(&out, "out", "", "Output .anim path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if in == "" || out == "" {
		return fmt.Errorf("both -in and -out are required")
	}

	raw, err := os.ReadFile(in)
	if err != nil {
		return err
	}

	var cfg format.AnimConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return err
	}

	built, err := format.BuildAnim(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(out, built, 0o644)
}
