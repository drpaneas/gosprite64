package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
)

func TestMk2DAnimProducesClipFrames(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "idle.json")
	out := filepath.Join(dir, "idle.anim")

	input := []byte(`{"clips":[{"name":"idle","fps":12,"frames":[0,1]}]}`)
	if err := os.WriteFile(in, input, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := run([]string{"-in", in, "-out", out}); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	parsed, err := format.ParseAnim(raw)
	if err != nil {
		t.Fatalf("ParseAnim() error = %v", err)
	}
	if len(parsed.Clips) != 1 || parsed.Clips[0].Name != "idle" {
		t.Fatalf("clips = %+v, want single idle clip", parsed.Clips)
	}
}
