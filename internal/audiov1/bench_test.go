package audiov1

import (
	"math"
	"testing"
)

func BenchmarkDecodeBlock(b *testing.B) {
	var cb Codebook
	cb[0][0][0] = 1024
	cb[0][1][0] = 2048
	cb[1][0][3] = 512
	cb[1][1][5] = 1536
	var state State
	block := [BlockBytes]byte{0x31, 0x72, 0x45, 0x89, 0xAB, 0xCD, 0xEF, 0x12, 0x34}
	var out [BlockSamples]int16

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		DecodeBlock(&cb, &state, block, out[:])
	}
}

func BenchmarkMix9Voices(b *testing.B) {
	m := NewMixer(48000, 512)
	out := make([]int16, 512*2)

	srcData := make([]int16, 300)
	for i := range srcData {
		srcData[i] = int16(1000 * math.Sin(float64(i)*0.1))
	}

	taps := make([]VoiceTap, MaxVoices)
	for i := range taps {
		taps[i] = VoiceTap{
			Samples: srcData,
			SrcRate: 16000,
			Gain:    GainFull,
			Active:  true,
		}
	}
	taps[0].SrcRate = 22050

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := range taps {
			taps[j].Consumed = 0
			taps[j].Phase = 0
		}
		m.Mix(out, taps)
	}
}

func BenchmarkMix1Voice48kPassthrough(b *testing.B) {
	m := NewMixer(48000, 512)
	out := make([]int16, 512*2)
	src := make([]int16, 512)
	for i := range src {
		src[i] = int16(i)
	}
	taps := []VoiceTap{{Samples: src, SrcRate: 48000, Gain: GainFull, Active: true}}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		taps[0].Consumed = 0
		taps[0].Phase = 0
		m.Mix(out, taps)
	}
}

func BenchmarkSourceFill256Samples(b *testing.B) {
	blocks := 20
	data := make([]byte, blocks*BlockBytes)
	for i := range data {
		data[i] = byte(i % 256)
	}
	for j := 0; j < blocks; j++ {
		data[j*BlockBytes] = 0x21
	}

	asset := &AssetEntry{
		AudibleFrames: uint32(blocks * BlockSamples),
		EncodedFrames: uint32(blocks * BlockSamples),
		Rate:          16000,
		BlockFrames:   BlockSamples,
	}

	var s SourceState
	var cb Codebook
	InitSource(&s, asset, data, &cb, State{})
	dst := make([]int16, 256)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.Reset()
		s.Fill(dst)
	}
}

func BenchmarkCommandRingPushPop(b *testing.B) {
	var ring CommandRing
	cmd := Command{Kind: CmdPlaySFX, ID: 1}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ring.Push(cmd)
		ring.Pop()
	}
}
