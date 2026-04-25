package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type wavFormat struct {
	audioFormat   uint16
	channels      uint16
	sampleRate    uint32
	bitsPerSample uint16
}

func readWAVMono(r io.Reader) ([]int16, int, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, 0, fmt.Errorf("read wav: %w", err)
	}
	if len(data) < 12 || string(data[:4]) != "RIFF" || string(data[8:12]) != "WAVE" {
		return nil, 0, fmt.Errorf("unsupported wav header")
	}

	format, pcmData, err := parseWAVChunks(data[12:])
	if err != nil {
		return nil, 0, err
	}
	if format.audioFormat != 1 {
		return nil, 0, fmt.Errorf("unsupported wav audio format %d", format.audioFormat)
	}
	if format.bitsPerSample != 16 {
		return nil, 0, fmt.Errorf("unsupported wav bit depth %d", format.bitsPerSample)
	}
	if format.channels != 1 && format.channels != 2 {
		return nil, 0, fmt.Errorf("unsupported wav channel count %d", format.channels)
	}
	if format.sampleRate == 0 {
		return nil, 0, fmt.Errorf("unsupported wav sample rate 0")
	}

	frameSize := int(format.channels) * 2
	if len(pcmData)%frameSize != 0 {
		return nil, 0, fmt.Errorf("wav data size %d is not aligned to frame size %d", len(pcmData), frameSize)
	}

	mono := make([]int16, 0, len(pcmData)/frameSize)
	for i := 0; i < len(pcmData); i += frameSize {
		left := int32(int16(binary.LittleEndian.Uint16(pcmData[i : i+2])))
		sample := left
		if format.channels == 2 {
			right := int32(int16(binary.LittleEndian.Uint16(pcmData[i+2 : i+4])))
			sample = (left + right) / 2
		}
		mono = append(mono, int16(sample))
	}

	return mono, int(format.sampleRate), nil
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

func resampleMono(samples []int16, srcRate, dstRate int) []int16 {
	if len(samples) == 0 || srcRate == dstRate {
		return samples
	}

	outCount := int(math.Round(float64(len(samples)) * float64(dstRate) / float64(srcRate)))
	if outCount < 1 {
		outCount = 1
	}

	out := make([]int16, 0, outCount)
	for i := 0; i < outCount; i++ {
		srcPos := float64(i) * float64(srcRate) / float64(dstRate)
		lo := int(math.Floor(srcPos))
		if lo >= len(samples) {
			lo = len(samples) - 1
		}
		hi := lo + 1
		if hi >= len(samples) {
			hi = len(samples) - 1
		}
		frac := srcPos - float64(lo)
		out = append(out, interpolateSample(samples[lo], samples[hi], frac))
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
