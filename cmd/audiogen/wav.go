package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/drpaneas/gosprite64/internal/audioengine"
)

type stereoFrame struct {
	left  int16
	right int16
}

type wavFormat struct {
	audioFormat   uint16
	channels      uint16
	sampleRate    uint32
	bitsPerSample uint16
}

func prepareAudioFiles(dir string) ([]string, error) {
	candidates := make(map[string]string)

	patterns := []string{"music*.wav", "sfx_*.wav", "music*.raw", "sfx_*.raw"}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			return nil, fmt.Errorf("error searching for %s files: %w", pattern, err)
		}

		slices.Sort(matches)
		for _, match := range matches {
			rawPath := strings.TrimSuffix(match, filepath.Ext(match)) + ".raw"
			if strings.HasSuffix(match, ".wav") || candidates[rawPath] == "" {
				candidates[rawPath] = match
			}
		}
	}

	rawFiles := make([]string, 0, len(candidates))
	for rawPath, sourcePath := range candidates {
		if strings.HasSuffix(sourcePath, ".wav") {
			if err := convertWAVFile(sourcePath, rawPath); err != nil {
				return nil, err
			}
		}
		rawFiles = append(rawFiles, rawPath)
	}

	slices.Sort(rawFiles)
	return rawFiles, nil
}

func convertWAVFile(sourcePath, rawPath string) error {
	f, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open wav %s: %w", sourcePath, err)
	}
	defer f.Close()

	rawData, err := convertWAVToRaw(f)
	if err != nil {
		return fmt.Errorf("convert wav %s: %w", sourcePath, err)
	}

	if err := os.WriteFile(rawPath, rawData, 0o644); err != nil {
		return fmt.Errorf("write raw %s: %w", rawPath, err)
	}
	return nil
}

func convertWAVToRaw(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read wav: %w", err)
	}
	if len(data) < 12 || string(data[:4]) != "RIFF" || string(data[8:12]) != "WAVE" {
		return nil, fmt.Errorf("unsupported wav header")
	}

	format, pcmData, err := parseWAVChunks(data[12:])
	if err != nil {
		return nil, err
	}

	if format.audioFormat != 1 {
		return nil, fmt.Errorf("unsupported wav audio format %d", format.audioFormat)
	}
	if format.bitsPerSample != 16 {
		return nil, fmt.Errorf("unsupported wav bit depth %d", format.bitsPerSample)
	}
	if format.channels != 1 && format.channels != 2 {
		return nil, fmt.Errorf("unsupported wav channel count %d", format.channels)
	}
	if format.sampleRate == 0 {
		return nil, fmt.Errorf("unsupported wav sample rate 0")
	}

	frames, err := decodePCM16Frames(format, pcmData)
	if err != nil {
		return nil, err
	}
	frames = resampleFrames(frames, int(format.sampleRate), audioengine.RuntimeSampleRate)

	raw := make([]byte, 0, len(frames)*audioengine.RuntimeChannels*audioengine.RuntimeBytesPerSample)
	for _, frame := range frames {
		raw = binary.BigEndian.AppendUint16(raw, uint16(frame.left))
		raw = binary.BigEndian.AppendUint16(raw, uint16(frame.right))
	}

	return raw, nil
}

func parseWAVChunks(data []byte) (wavFormat, []byte, error) {
	var format wavFormat
	var pcmData []byte

	for len(data) >= 8 {
		chunkID := string(data[:4])
		chunkSize := int(binary.LittleEndian.Uint32(data[4:8]))
		data = data[8:]
		if len(data) < chunkSize {
			return wavFormat{}, nil, fmt.Errorf("truncated %s chunk", chunkID)
		}

		chunkData := data[:chunkSize]
		switch chunkID {
		case "fmt ":
			if len(chunkData) < 16 {
				return wavFormat{}, nil, fmt.Errorf("truncated fmt chunk")
			}
			format = wavFormat{
				audioFormat:   binary.LittleEndian.Uint16(chunkData[0:2]),
				channels:      binary.LittleEndian.Uint16(chunkData[2:4]),
				sampleRate:    binary.LittleEndian.Uint32(chunkData[4:8]),
				bitsPerSample: binary.LittleEndian.Uint16(chunkData[14:16]),
			}
		case "data":
			pcmData = append([]byte(nil), chunkData...)
		}

		data = data[chunkSize:]
		if chunkSize%2 == 1 && len(data) > 0 {
			data = data[1:]
		}
	}

	if pcmData == nil {
		return wavFormat{}, nil, fmt.Errorf("wav file missing data chunk")
	}
	return format, pcmData, nil
}

func decodePCM16Frames(format wavFormat, pcmData []byte) ([]stereoFrame, error) {
	frameSize := int(format.channels) * 2
	if len(pcmData)%frameSize != 0 {
		return nil, fmt.Errorf("wav data size %d is not aligned to frame size %d", len(pcmData), frameSize)
	}

	frameCount := len(pcmData) / frameSize
	frames := make([]stereoFrame, 0, frameCount)
	for i := 0; i < len(pcmData); i += frameSize {
		left := int16(binary.LittleEndian.Uint16(pcmData[i : i+2]))
		right := left
		if format.channels == 2 {
			right = int16(binary.LittleEndian.Uint16(pcmData[i+2 : i+4]))
		}
		frames = append(frames, stereoFrame{left: left, right: right})
	}

	return frames, nil
}

func resampleFrames(frames []stereoFrame, srcRate, dstRate int) []stereoFrame {
	if len(frames) == 0 || srcRate == dstRate {
		return frames
	}

	outCount := int(math.Round(float64(len(frames)) * float64(dstRate) / float64(srcRate)))
	if outCount < 1 {
		outCount = 1
	}

	out := make([]stereoFrame, 0, outCount)
	for i := 0; i < outCount; i++ {
		srcPos := float64(i) * float64(srcRate) / float64(dstRate)
		lo := int(math.Floor(srcPos))
		if lo >= len(frames) {
			lo = len(frames) - 1
		}
		hi := lo + 1
		if hi >= len(frames) {
			hi = len(frames) - 1
		}
		frac := srcPos - float64(lo)

		out = append(out, stereoFrame{
			left:  interpolateSample(frames[lo].left, frames[hi].left, frac),
			right: interpolateSample(frames[lo].right, frames[hi].right, frac),
		})
	}

	return out
}

func interpolateSample(a, b int16, frac float64) int16 {
	if frac <= 0 {
		return a
	}
	if frac >= 1 {
		return b
	}

	value := float64(a) + (float64(b)-float64(a))*frac
	return int16(math.Round(value))
}
