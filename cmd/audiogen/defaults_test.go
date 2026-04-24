//go:build !noos

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanAudioDirV1(t *testing.T) {
	dir := t.TempDir()
	sfxDir := filepath.Join(dir, "assets", "audio", "sfx")
	musicDir := filepath.Join(dir, "assets", "audio", "music")
	if err := os.MkdirAll(sfxDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(musicDir, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, filepath.Join(sfxDir, "jump.wav"), encodePCM16WAV(t, 16000, 1, make([]int16, 160)))
	mustWriteFile(t, filepath.Join(sfxDir, "hit.wav"), encodePCM16WAV(t, 16000, 1, make([]int16, 160)))
	mustWriteFile(t, filepath.Join(musicDir, "overworld.wav"), encodePCM16WAV(t, 22050, 1, make([]int16, 2205)))

	sfx, music, err := scanAudioDirV1(dir)
	if err != nil {
		t.Fatalf("scanAudioDirV1 error: %v", err)
	}
	if len(sfx) != 2 || len(music) != 1 {
		t.Fatalf("got %d sfx, %d music; want 2, 1", len(sfx), len(music))
	}
	if music[0].Name != "overworld" {
		t.Fatalf("music[0].Name = %q, want overworld", music[0].Name)
	}
}
