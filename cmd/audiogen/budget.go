//go:build !noos

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/drpaneas/gosprite64/internal/audiov1"
)

func computeAndWriteBudget(dir string, inputs []manifestInput, romCap, sfxCap int) error {
	var sfxResident, musicStreamed, manifestBytes, codebookBytes uint64

	for _, in := range inputs {
		switch in.Class {
		case "sfx":
			sfxResident += uint64(in.DataBytes)
		case "music":
			musicStreamed += uint64(in.DataBytes)
		}
		codebookBytes += uint64(in.AuxBytes)
	}

	manifestBytes = uint64(len(inputs)) * 48
	romTotal := sfxResident + musicStreamed + manifestBytes + codebookBytes

	if sfxResident > uint64(sfxCap) {
		return fmt.Errorf("SFX resident size %d exceeds budget %d bytes", sfxResident, sfxCap)
	}
	if romTotal > uint64(romCap) {
		return fmt.Errorf("total audio ROM %d exceeds budget %d bytes", romTotal, romCap)
	}

	var estBandwidth float64
	for _, in := range inputs {
		if in.Class == "music" && in.Rate > 0 {
			blocksPerSec := math.Ceil(float64(in.Rate) / float64(audiov1.BlockSamples))
			estBandwidth = blocksPerSec * float64(audiov1.BlockBytes)
		}
	}

	blocksPer10ms := math.Ceil(float64(DefaultMusicRate) / float64(audiov1.BlockSamples) / 100.0)
	estDecodeCost := blocksPer10ms * float64(DefaultDecodeCostUsec)

	dacBufferFrames := 512
	voiceStateBytes := audiov1.MaxVoices * 16
	sourceDecBufBytes := audiov1.MaxVoices * audiov1.BlockSamples * 2
	sourceStructBytes := audiov1.MaxVoices * 256
	accumBytes := dacBufferFrames * 4
	outputBufBytes := dacBufferFrames * 2 * 2
	cmdRingBytes := 16 * 6
	fixedRuntimeRAM := voiceStateBytes + sourceDecBufBytes + sourceStructBytes + accumBytes + outputBufBytes + cmdRingBytes

	report := map[string]any{
		"romTotal":                      romTotal,
		"sfxResident":                   sfxResident,
		"musicStreamed":                 musicStreamed,
		"manifestBytes":                 manifestBytes,
		"codebookBytes":                 codebookBytes,
		"worstCaseVoicesSfx":            audiov1.MaxSFXVoices,
		"musicReserves":                 audiov1.MaxMusicVoices,
		"estStreamBandwidthBytesPerSec": estBandwidth,
		"estDecodeUsecPer10msAudio":     estDecodeCost,
		"fixedRuntimeRAMBytes":          fixedRuntimeRAM,
		"fixedRuntimeRAMBreakdown": map[string]int{
			"voiceStates":     voiceStateBytes,
			"sourceDecodeBuf": sourceDecBufBytes,
			"sourceStructs":   sourceStructBytes,
			"mixerAccum":      accumBytes,
			"outputBuffer":    outputBufBytes,
			"commandRing":     cmdRingBytes,
		},
		"embedResidencyNote": "//go:embed data is assumed ROM-resident. If the toolchain copies embedded bytes into RDRAM, music is not truly streamed and resident SFX budgets must account for the full copy.",
	}

	buildDir := filepath.Join(dir, BuildDirName)
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(buildDir, AudioReportName), data, 0o644)
}
