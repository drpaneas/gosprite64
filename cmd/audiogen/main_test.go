package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanAudioDirV1FindsSFXAndMusic(t *testing.T) {
	dir := t.TempDir()
	sfxDir := filepath.Join(dir, "assets", "audio", "sfx")
	musicDir := filepath.Join(dir, "assets", "audio", "music")
	if err := os.MkdirAll(sfxDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(musicDir, 0o755); err != nil {
		t.Fatal(err)
	}

	mustWriteFile(t, filepath.Join(sfxDir, "jump.wav"), encodePCM16WAV(t, 48000, 1, []int16{100, 200}))
	mustWriteFile(t, filepath.Join(sfxDir, "coin.wav"), encodePCM16WAV(t, 48000, 1, []int16{300}))
	mustWriteFile(t, filepath.Join(musicDir, "title.wav"), encodePCM16WAV(t, 44100, 1, []int16{500, 600, 700}))

	sfx, music, err := scanAudioDirV1(dir)
	if err != nil {
		t.Fatalf("scanAudioDirV1 error: %v", err)
	}
	if len(sfx) != 2 {
		t.Fatalf("len(sfx) = %d, want 2", len(sfx))
	}
	if len(music) != 1 {
		t.Fatalf("len(music) = %d, want 1", len(music))
	}
	if sfx[0].Name != "coin" || sfx[1].Name != "jump" {
		t.Fatalf("sfx names = [%s, %s], want [coin, jump]", sfx[0].Name, sfx[1].Name)
	}
	if music[0].Name != "title" {
		t.Fatalf("music name = %s, want title", music[0].Name)
	}
}

func TestScanAudioDirV1EmptyDir(t *testing.T) {
	dir := t.TempDir()
	sfx, music, err := scanAudioDirV1(dir)
	if err != nil {
		t.Fatalf("scanAudioDirV1 error: %v", err)
	}
	if len(sfx) != 0 || len(music) != 0 {
		t.Fatalf("expected empty, got %d sfx and %d music", len(sfx), len(music))
	}
}

func mustWriteFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("os.WriteFile(%s): %v", path, err)
	}
}
