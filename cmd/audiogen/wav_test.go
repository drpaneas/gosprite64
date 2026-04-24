package main

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestConvertWAVToRawDuplicatesMonoSamplesIntoStereoBigEndianPCM(t *testing.T) {
	wavData := encodePCM16WAV(t, 48000, 1, []int16{0x0102, -2})

	rawData, err := convertWAVToRaw(bytes.NewReader(wavData))
	if err != nil {
		t.Fatalf("convertWAVToRaw returned error: %v", err)
	}

	want := []byte{
		0x01, 0x02, 0x01, 0x02,
		0xff, 0xfe, 0xff, 0xfe,
	}
	if !bytes.Equal(rawData, want) {
		t.Fatalf("convertWAVToRaw = %v, want %v", rawData, want)
	}
}

func TestConvertWAVToRawResamplesToRuntimeSampleRate(t *testing.T) {
	wavData := encodePCM16WAV(t, 24000, 1, []int16{1000, 1000})

	rawData, err := convertWAVToRaw(bytes.NewReader(wavData))
	if err != nil {
		t.Fatalf("convertWAVToRaw returned error: %v", err)
	}

	wantFrame := []byte{0x03, 0xe8, 0x03, 0xe8}
	want := bytes.Repeat(wantFrame, 4)
	if !bytes.Equal(rawData, want) {
		t.Fatalf("convertWAVToRaw resampled bytes = %v, want %v", rawData, want)
	}
}

func TestConvertWAVToRawPreservesStereoChannels(t *testing.T) {
	wavData := encodePCM16StereoWAV(t, 48000, [][2]int16{
		{0x0102, -2},
		{0x0304, -4},
	})

	rawData, err := convertWAVToRaw(bytes.NewReader(wavData))
	if err != nil {
		t.Fatalf("convertWAVToRaw returned error: %v", err)
	}

	want := []byte{
		0x01, 0x02, 0xff, 0xfe,
		0x03, 0x04, 0xff, 0xfc,
	}
	if !bytes.Equal(rawData, want) {
		t.Fatalf("convertWAVToRaw stereo bytes = %v, want %v", rawData, want)
	}
}

func TestConvertWAVToRawRejectsNonPCMInput(t *testing.T) {
	wavData := encodeWAV(t, 3, 48000, 1, 16, []byte{0x00, 0x00})

	_, err := convertWAVToRaw(bytes.NewReader(wavData))
	if err == nil {
		t.Fatal("convertWAVToRaw error = nil, want unsupported format error")
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
