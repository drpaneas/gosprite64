//go:build !noos

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestComputeBudgetWritesReport(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, BuildDirName), 0o755); err != nil {
		t.Fatal(err)
	}
	inputs := []manifestInput{
		{Name: "jump", Class: "sfx", DataBytes: 100, AuxBytes: 128},
		{Name: "overworld", Class: "music", DataBytes: 12000, AuxBytes: 144, Rate: DefaultMusicRate},
	}

	if err := computeAndWriteBudget(dir, inputs, DefaultROMBudget, DefaultSFXResidentCap); err != nil {
		t.Fatalf("computeAndWriteBudget error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, BuildDirName, AudioReportName))
	if err != nil {
		t.Fatalf("read report: %v", err)
	}
	var report map[string]any
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("parse report: %v", err)
	}
	if report["romTotal"].(float64) == 0 || report["fixedRuntimeRAMBytes"].(float64) == 0 {
		t.Fatalf("bad budget report: %v", report)
	}
	if report["sfxResident"].(float64) != 100 {
		t.Fatalf("sfxResident = %v, want 100", report["sfxResident"])
	}
}

func TestComputeBudgetFailsOnOversize(t *testing.T) {
	dir := t.TempDir()
	inputs := []manifestInput{{Name: "big", Class: "sfx", DataBytes: 50000, AuxBytes: 128}}
	if err := computeAndWriteBudget(dir, inputs, DefaultROMBudget, DefaultSFXResidentCap); err == nil {
		t.Fatal("expected budget error")
	}
}
