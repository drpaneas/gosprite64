package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
)

func TestMk2DBundleEmitsManifestReferences(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "level.bundle")

	if err := run([]string{
		"-sheet", "tiles.sheet",
		"-map", "level.map",
		"-anim", "idle.anim",
		"-out", out,
	}); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if got := string(raw[:4]); got != "BND2" {
		t.Fatalf("magic = %q, want BND2", got)
	}

	parsed, err := format.ParseBundle(raw)
	if err != nil {
		t.Fatalf("ParseBundle() error = %v", err)
	}
	if len(parsed.Entries) != 3 {
		t.Fatalf("len(Entries) = %d, want 3", len(parsed.Entries))
	}
	if parsed.Entries[0].Path != "tiles.sheet" {
		t.Fatalf("first path = %q, want tiles.sheet", parsed.Entries[0].Path)
	}
}

func TestMk2DBundlePreservesSheetFlagOrder(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "level.bundle")

	if err := run([]string{
		"-sheet", "z.sheet",
		"-sheet", "a.sheet",
		"-map", "level.map",
		"-out", out,
	}); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	parsed, err := format.ParseBundle(raw)
	if err != nil {
		t.Fatalf("ParseBundle() error = %v", err)
	}

	if len(parsed.Entries) < 2 {
		t.Fatalf("len(Entries) = %d, want at least 2", len(parsed.Entries))
	}
	if parsed.Entries[0].Path != "z.sheet" || parsed.Entries[1].Path != "a.sheet" {
		t.Fatalf("sheet order = [%q %q], want [z.sheet a.sheet]", parsed.Entries[0].Path, parsed.Entries[1].Path)
	}
}
