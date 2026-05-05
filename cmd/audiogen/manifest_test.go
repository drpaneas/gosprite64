//go:build !noos

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateManifestAndConstPackages(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, BuildDirName), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, filepath.Join(dir, "dummy.go"), []byte("package main\n"))

	inputs := []manifestInput{
		{Name: "jump", Class: "sfx", ID: 0, Rate: DefaultSFXRate, AudibleFrames: 160, EncodedFrames: 160, DataBytes: 90, AuxBytes: 128},
		{Name: "overworld", Class: "music", ID: 0, Rate: DefaultMusicRate, AudibleFrames: 22050, EncodedFrames: 22064, DataOffset: 90, DataBytes: 12420, AuxOffset: 128, AuxBytes: 144, Loop: true, LoopLen: 22064},
	}

	if err := generateManifestAndConsts(dir, inputs, ""); err != nil {
		t.Fatalf("generateManifestAndConsts error: %v", err)
	}

	sfxContent := mustReadString(t, filepath.Join(dir, "sfx", "ids.go"))
	if !strings.Contains(sfxContent, "type ID = audiosfx.ID") || !strings.Contains(sfxContent, "Jump") {
		t.Fatalf("bad sfx ids:\n%s", sfxContent)
	}
	embedContent := mustReadString(t, filepath.Join(dir, AudioEmbedName))
	if !strings.Contains(embedContent, "AudioBundle") || !strings.Contains(embedContent, "RegisterAudioBundle") || !strings.Contains(embedContent, "package main") {
		t.Fatalf("bad embed file:\n%s", embedContent)
	}
	if strings.Contains(embedContent, "RegisterAudioV1") || strings.Contains(embedContent, "RegisterSFXNameResolver") {
		t.Fatalf("embed file still uses legacy audio names:\n%s", embedContent)
	}
	if !strings.Contains(embedContent, "ResolveSoundEffectName") || !strings.Contains(embedContent, "return 0, true") {
		t.Fatalf("embed file missing sound effect resolver:\n%s", embedContent)
	}
}

func TestValidateLoopInvariantsRejectsBadAlignment(t *testing.T) {
	tests := []struct {
		name      string
		loopStart uint32
		loopLen   uint32
		total     uint32
		wantErr   bool
	}{
		{"aligned", 16, 32, 64, false},
		{"loop start not mod 16", 15, 32, 64, true},
		{"loop len not mod 16", 16, 30, 64, true},
		{"loop exceeds total", 16, 64, 48, true},
		{"valid no loop", 0, 0, 64, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLoopInvariants(tt.loopStart, tt.loopLen, tt.total, tt.loopLen > 0)
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func mustReadString(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
