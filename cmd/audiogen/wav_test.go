package main

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestReadWAVMonoFromMonoInput(t *testing.T) {
	wavData := encodePCM16WAV(t, 48000, 1, []int16{0x0102, -2})

	mono, rate, err := readWAVMono(bytes.NewReader(wavData))
	if err != nil {
		t.Fatalf("readWAVMono returned error: %v", err)
	}
	if rate != 48000 {
		t.Fatalf("rate = %d, want 48000", rate)
	}
	want := []int16{0x0102, -2}
	if len(mono) != len(want) {
		t.Fatalf("len(mono) = %d, want %d", len(mono), len(want))
	}
	for i := range want {
		if mono[i] != want[i] {
			t.Fatalf("mono[%d] = %d, want %d", i, mono[i], want[i])
		}
	}
}

func TestReadWAVMonoDownmixesStereo(t *testing.T) {
	wavData := encodePCM16StereoWAV(t, 48000, [][2]int16{
		{1000, 3000},
		{-2000, -4000},
	})

	mono, rate, err := readWAVMono(bytes.NewReader(wavData))
	if err != nil {
		t.Fatalf("readWAVMono returned error: %v", err)
	}
	if rate != 48000 {
		t.Fatalf("rate = %d, want 48000", rate)
	}
	if mono[0] != 2000 {
		t.Fatalf("mono[0] = %d, want 2000 ((1000+3000)/2)", mono[0])
	}
	if mono[1] != -3000 {
		t.Fatalf("mono[1] = %d, want -3000 ((-2000+-4000)/2)", mono[1])
	}
}

func TestReadWAVMonoRejectsNonPCM(t *testing.T) {
	wavData := encodeWAV(t, 3, 48000, 1, 16, []byte{0x00, 0x00})

	_, _, err := readWAVMono(bytes.NewReader(wavData))
	if err == nil {
		t.Fatal("readWAVMono error = nil, want unsupported format error")
	}
}

func TestResampleMonoIdentityRate(t *testing.T) {
	input := []int16{100, 200, 300}
	out := resampleMono(input, 16000, 16000)
	if len(out) != len(input) {
		t.Fatalf("len = %d, want %d", len(out), len(input))
	}
	for i := range input {
		if out[i] != input[i] {
			t.Fatalf("out[%d] = %d, want %d", i, out[i], input[i])
		}
	}
}

func TestResampleMonoDownsample(t *testing.T) {
	input := []int16{1000, 1000, 1000, 1000}
	out := resampleMono(input, 48000, 16000)
	if len(out) == 0 {
		t.Fatal("resampleMono returned empty")
	}
	for i, s := range out {
		if s != 1000 {
			t.Fatalf("out[%d] = %d, want 1000", i, s)
		}
	}
}

func encodePCM16WAV(t *testing.T, sampleRate, channels int, samples []int16) []byte {
	t.Helper()

	data := new(bytes.Buffer)
	for _, sample := range samples {
		if err := binary.Write(data, binary.LittleEndian, sample); err != nil {
			t.Fatalf("binary.Write sample: %v", err)
		}
	}

	return encodeWAV(t, 1, sampleRate, channels, 16, data.Bytes())
}

func encodePCM16StereoWAV(t *testing.T, sampleRate int, frames [][2]int16) []byte {
	t.Helper()

	data := new(bytes.Buffer)
	for _, frame := range frames {
		if err := binary.Write(data, binary.LittleEndian, frame[0]); err != nil {
			t.Fatalf("binary.Write left sample: %v", err)
		}
		if err := binary.Write(data, binary.LittleEndian, frame[1]); err != nil {
			t.Fatalf("binary.Write right sample: %v", err)
		}
	}

	return encodeWAV(t, 1, sampleRate, 2, 16, data.Bytes())
}

func encodeWAV(t *testing.T, audioFormat, sampleRate, channels, bitsPerSample int, data []byte) []byte {
	t.Helper()

	buf := new(bytes.Buffer)
	byteRate := sampleRate * channels * bitsPerSample / 8
	blockAlign := channels * bitsPerSample / 8
	riffSize := 4 + (8 + 16) + (8 + len(data))

	buf.WriteString("RIFF")
	mustWriteLE(t, buf, uint32(riffSize))
	buf.WriteString("WAVE")

	buf.WriteString("fmt ")
	mustWriteLE(t, buf, uint32(16))
	mustWriteLE(t, buf, uint16(audioFormat))
	mustWriteLE(t, buf, uint16(channels))
	mustWriteLE(t, buf, uint32(sampleRate))
	mustWriteLE(t, buf, uint32(byteRate))
	mustWriteLE(t, buf, uint16(blockAlign))
	mustWriteLE(t, buf, uint16(bitsPerSample))

	buf.WriteString("data")
	mustWriteLE(t, buf, uint32(len(data)))
	if _, err := buf.Write(data); err != nil {
		t.Fatalf("write wav data: %v", err)
	}

	return buf.Bytes()
}

func mustWriteLE[T ~uint16 | ~uint32](t *testing.T, buf *bytes.Buffer, value T) {
	t.Helper()
	if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
		t.Fatalf("binary.Write value: %v", err)
	}
}
