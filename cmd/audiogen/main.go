//go:build !noos

// Command audiogen converts WAV audio assets into VADPCM compressed format
// and generates the Go embed file, typed ID constants, and size reports
// for gosprite64 games. This is a host-only tool that runs on your
// development machine, not on the N64.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var verbose bool

func main() {
	var outputDir string
	var romBudgetFlag int
	var sfxBudgetFlag int
	var devRawFlag bool
	var packageFlag string
	flag.StringVar(&outputDir, "dir", ".", "Game directory containing assets/audio/")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.IntVar(&romBudgetFlag, "rom-budget", DefaultROMBudget, "ROM budget for all audio data in bytes")
	flag.IntVar(&sfxBudgetFlag, "sfx-resident-budget", DefaultSFXResidentCap, "SFX resident budget in bytes")
	flag.BoolVar(&devRawFlag, "dev-raw", false, "Emit uncompressed PCM (not implemented, reserved for future use)")
	flag.StringVar(&packageFlag, "package", "", "Override package name for generated files")
	flag.Parse()

	if devRawFlag {
		fmt.Println("-dev-raw is not implemented. Remove the flag.")
		os.Exit(1)
	}
	if err := runV1Pipeline(outputDir, packageFlag, romBudgetFlag, sfxBudgetFlag); err != nil {
		fmt.Printf("audiogen error: %v\n", err)
		os.Exit(1)
	}
}

func runV1Pipeline(dir, packageName string, romBudget, sfxBudget int) error {
	sfxSources, musicSources, err := scanAudioDirV1(dir)
	if err != nil {
		return err
	}
	if len(sfxSources) == 0 && len(musicSources) == 0 {
		fmt.Println("No audio assets found in assets/audio/")
		return nil
	}

	buildDir := filepath.Join(dir, BuildDirName)
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return err
	}

	var allData []byte
	var allAux []byte
	var manifest []manifestInput
	var sfxID uint16
	var musicID uint16

	for _, src := range sfxSources {
		result, err := encodeAsset(src)
		if err != nil {
			return fmt.Errorf("encode sfx %s: %w", src.Name, err)
		}
		if err := validateLoopInvariants(result.LoopStart, result.LoopLen, result.EncodedFrames, result.Loop); err != nil {
			return fmt.Errorf("sfx %s: %w", src.Name, err)
		}
		manifest = append(manifest, manifestInput{
			Name: src.Name, Class: "sfx", ID: sfxID, Rate: src.Rate,
			AudibleFrames: result.AudibleFrames, EncodedFrames: result.EncodedFrames,
			DataOffset: uint32(len(allData)), DataBytes: uint32(len(result.CompressedData)),
			AuxOffset: uint32(len(allAux)), AuxBytes: uint32(len(result.CodebookBytes)),
		})
		allData = append(allData, result.CompressedData...)
		allAux = append(allAux, result.CodebookBytes...)
		sfxID++
	}

	for _, src := range musicSources {
		result, err := encodeAsset(src)
		if err != nil {
			return fmt.Errorf("encode music %s: %w", src.Name, err)
		}
		if err := validateLoopInvariants(result.LoopStart, result.LoopLen, result.EncodedFrames, result.Loop); err != nil {
			return fmt.Errorf("music %s: %w", src.Name, err)
		}
		loopState := captureStateAt(&result.Codebook, result.CompressedData, int(result.LoopStart))
		loopStateBytes := stateToBytes(loopState)
		manifest = append(manifest, manifestInput{
			Name: src.Name, Class: "music", ID: musicID, Rate: src.Rate,
			AudibleFrames: result.AudibleFrames, EncodedFrames: result.EncodedFrames,
			DataOffset: uint32(len(allData)), DataBytes: uint32(len(result.CompressedData)),
			AuxOffset: uint32(len(allAux)), AuxBytes: uint32(len(result.CodebookBytes) + len(loopStateBytes)),
			Loop: true, LoopStart: result.LoopStart, LoopLen: result.LoopLen,
		})
		allData = append(allData, result.CompressedData...)
		allAux = append(allAux, result.CodebookBytes...)
		allAux = append(allAux, loopStateBytes...)
		musicID++
	}

	if err := os.WriteFile(filepath.Join(buildDir, AudioBlobName), allData, 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(buildDir, AudioAuxName), allAux, 0o644); err != nil {
		return err
	}
	if err := generateManifestAndConsts(dir, manifest, packageName); err != nil {
		return err
	}
	return computeAndWriteBudget(dir, manifest, romBudget, sfxBudget)
}
