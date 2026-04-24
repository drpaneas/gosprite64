//go:build !noos

package main

import (
	"math"
	"testing"

	"github.com/drpaneas/gosprite64/internal/audiov1"
)

func TestAutocorrelationDCSignal(t *testing.T) {
	samples := make([]int16, 32)
	for i := range samples {
		samples[i] = 1000
	}
	r := autocorrelation(samples, audiov1.Order)

	for k := 1; k <= audiov1.Order; k++ {
		if math.Abs(r[k]/r[0]-1.0) > 1e-9 {
			t.Fatalf("R[%d]/R[0] = %f, want 1.0", k, r[k]/r[0])
		}
	}
}

func TestLevinsonDurbinDCSignal(t *testing.T) {
	lpc := levinsonDurbin([]float64{1000000, 1000000, 1000000}, audiov1.Order)
	if math.Abs(lpc[0]-0.999) > 0.01 {
		t.Fatalf("lpc[0] = %f, want ~0.999", lpc[0])
	}
	if math.Abs(lpc[1]) > 0.01 {
		t.Fatalf("lpc[1] = %f, want ~0", lpc[1])
	}
}

func TestMakeVectorsPassThrough(t *testing.T) {
	vectors := makeVectors([]float64{1.0, 0.0})
	for i := 0; i < audiov1.StateLen; i++ {
		if vectors[1][i] != 2048 {
			t.Fatalf("vectors[1][%d] = %d, want 2048", i, vectors[1][i])
		}
		if vectors[0][i] != 0 {
			t.Fatalf("vectors[0][%d] = %d, want 0", i, vectors[0][i])
		}
	}
}

func TestEncodeBlockControlByteFormat(t *testing.T) {
	sine := generateSine(440, 16000, 320)
	trained := trainCodebook(sine)
	var state audiov1.State
	block := encodeBlock(&trained.Book, &state, sine[:audiov1.BlockSamples])

	scale := int(block[0] >> 4)
	pred := int(block[0] & 0x0F)
	if pred < 0 || pred >= audiov1.PredictorCount {
		t.Fatalf("predictor index %d out of range", pred)
	}
	if scale < 0 || scale > 12 {
		t.Fatalf("scale %d out of range", scale)
	}
}

func TestEncodeVADPCMDeterminism(t *testing.T) {
	sine := generateSine(880, 16000, 320)
	enc1 := EncodeVADPCM(sine)
	enc2 := EncodeVADPCM(sine)

	if string(enc1.Data) != string(enc2.Data) {
		t.Fatal("encoded data differs between identical inputs")
	}
	if enc1.Codebook != enc2.Codebook {
		t.Fatal("codebooks differ between identical inputs")
	}
}

func TestEncodeVADPCMSilenceProducesZeroOutput(t *testing.T) {
	encoded := EncodeVADPCM(make([]int16, 160))
	decoded := decodeAll(encoded)
	for i, s := range decoded {
		if s != 0 {
			t.Fatalf("decoded silence[%d] = %d, want 0", i, s)
		}
	}
}

func TestEncodeVADPCMLoopStability(t *testing.T) {
	sine := generateSine(440, 22050, 22050)
	loopStartFrame := uint32(16 * 10)
	loopLenFrames := uint32(16 * 50)
	enc := EncodeVADPCM(sine)
	loopState := captureStateAt(&enc.Codebook, enc.Data, int(loopStartFrame))

	reference := decodeLoopBody(enc, loopState, loopStartFrame, loopLenFrames)
	for lap := 0; lap < 100; lap++ {
		got := decodeLoopBody(enc, loopState, loopStartFrame, loopLenFrames)
		for i := range reference {
			if got[i] != reference[i] {
				t.Fatalf("loop iteration %d sample %d = %d, want %d", lap, i, got[i], reference[i])
			}
		}
	}
}

func decodeLoopBody(enc EncodedAsset, state audiov1.State, start, length uint32) []int16 {
	out := make([]int16, length)
	startBlock := int(start / audiov1.BlockSamples)
	endBlock := int((start + length) / audiov1.BlockSamples)
	for b := startBlock; b < endBlock; b++ {
		var block [audiov1.BlockBytes]byte
		copy(block[:], enc.Data[b*audiov1.BlockBytes:])
		off := (b - startBlock) * audiov1.BlockSamples
		audiov1.DecodeBlock(&enc.Codebook, &state, block, out[off:off+audiov1.BlockSamples])
	}
	return out
}

func decodeAll(enc EncodedAsset) []int16 {
	decoded := make([]int16, len(enc.Data)/audiov1.BlockBytes*audiov1.BlockSamples)
	var state audiov1.State
	for b := 0; b < len(enc.Data)/audiov1.BlockBytes; b++ {
		var block [audiov1.BlockBytes]byte
		copy(block[:], enc.Data[b*audiov1.BlockBytes:])
		audiov1.DecodeBlock(&enc.Codebook, &state, block, decoded[b*audiov1.BlockSamples:])
	}
	return decoded
}

func generateSine(freq, sampleRate float64, numSamples int) []int16 {
	samples := make([]int16, numSamples)
	for i := range samples {
		samples[i] = int16(16000.0 * math.Sin(2*math.Pi*freq*float64(i)/sampleRate))
	}
	return samples
}
