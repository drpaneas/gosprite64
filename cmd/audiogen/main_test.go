package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestPrepareAudioFilesPrefersWAVSourcesAndReturnsRawOutputs(t *testing.T) {
	dir := t.TempDir()

	mustWriteFile(t, filepath.Join(dir, "music0.wav"), encodePCM16WAV(t, 48000, 1, []int16{0x0102}))
	mustWriteFile(t, filepath.Join(dir, "music0.raw"), []byte("stale-raw"))
	mustWriteFile(t, filepath.Join(dir, "sfx_ping.wav"), encodePCM16WAV(t, 48000, 1, []int16{-2}))

	audioFiles, err := prepareAudioFiles(dir)
	if err != nil {
		t.Fatalf("prepareAudioFiles returned error: %v", err)
	}

	want := []string{
		filepath.Join(dir, "music0.raw"),
		filepath.Join(dir, "sfx_ping.raw"),
	}
	if !slices.Equal(audioFiles, want) {
		t.Fatalf("prepareAudioFiles = %v, want %v", audioFiles, want)
	}

	convertedMusic, err := os.ReadFile(filepath.Join(dir, "music0.raw"))
	if err != nil {
		t.Fatalf("os.ReadFile converted music: %v", err)
	}
	if string(convertedMusic) == "stale-raw" {
		t.Fatal("prepareAudioFiles did not overwrite stale raw output from the preferred wav source")
	}
}

func TestGenerateEmbedFileIncludesPreparedRawOutputs(t *testing.T) {
	dir := t.TempDir()

	audioFiles := []string{
		filepath.Join(dir, "music0.raw"),
		filepath.Join(dir, "sfx_ping.raw"),
	}
	mustWriteFile(t, audioFiles[0], []byte{0x00})
	mustWriteFile(t, audioFiles[1], []byte{0x00})

	if err := generateEmbedFile(dir, audioFiles); err != nil {
		t.Fatalf("generateEmbedFile returned error: %v", err)
	}

	embedFile, err := os.ReadFile(filepath.Join(dir, "audio_embed.go"))
	if err != nil {
		t.Fatalf("os.ReadFile audio_embed.go: %v", err)
	}

	content := string(embedFile)
	if !strings.Contains(content, "//go:embed music0.raw sfx_ping.raw") {
		t.Fatalf("audio_embed.go missing raw embed directive:\n%s", content)
	}
}

func TestPrepareAudioFilesSupportsSFXOnlyExampleLayouts(t *testing.T) {
	dir := t.TempDir()

	mustWriteFile(t, filepath.Join(dir, "sfx_start.wav"), encodePCM16WAV(t, 48000, 1, []int16{0x0102}))
	mustWriteFile(t, filepath.Join(dir, "sfx_paddle.wav"), encodePCM16WAV(t, 48000, 1, []int16{0x0203}))
	mustWriteFile(t, filepath.Join(dir, "sfx_wall.wav"), encodePCM16WAV(t, 48000, 1, []int16{0x0304}))

	audioFiles, err := prepareAudioFiles(dir)
	if err != nil {
		t.Fatalf("prepareAudioFiles returned error: %v", err)
	}

	want := []string{
		filepath.Join(dir, "sfx_paddle.raw"),
		filepath.Join(dir, "sfx_start.raw"),
		filepath.Join(dir, "sfx_wall.raw"),
	}
	if !slices.Equal(audioFiles, want) {
		t.Fatalf("prepareAudioFiles = %v, want %v", audioFiles, want)
	}

	for _, rawPath := range want {
		if _, err := os.Stat(rawPath); err != nil {
			t.Fatalf("expected generated raw file %s: %v", rawPath, err)
		}
	}
}

func mustWriteFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("os.WriteFile(%s): %v", path, err)
	}
}
