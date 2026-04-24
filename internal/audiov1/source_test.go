package audiov1

import "testing"

func makeTestAsset(class AssetClass, audible, encoded uint32, loop bool, loopStart, loopLen uint32) (*AssetEntry, []byte, Codebook) {
	blocks := int(encoded / BlockSamples)
	data := make([]byte, blocks*BlockBytes)
	for b := 0; b < blocks; b++ {
		data[b*BlockBytes] = 0x10
		data[b*BlockBytes+1] = 0x11
	}
	flags := FlagResident
	if class == ClassMusic {
		flags = FlagStreamed
	}
	if loop {
		flags |= FlagLoop
	}
	return &AssetEntry{
		ID: 0, Class: class, Flags: flags, Rate: 16000,
		AudibleFrames: audible, EncodedFrames: encoded,
		LoopStart: loopStart, LoopLen: loopLen, DataBytes: uint32(len(data)),
		MaxInstances: DefaultMaxInst, BlockFrames: BlockSamples,
	}, data, Codebook{}
}

func TestFillOneShotSFXCompletes(t *testing.T) {
	entry, data, cb := makeTestAsset(ClassSFX, 32, 32, false, 0, 0)
	var src SourceState
	InitSource(&src, entry, data, &cb, State{})

	dst := make([]int16, 64)
	n, ended := src.Fill(dst)

	if n != 32 || !ended {
		t.Fatalf("Fill returned n=%d ended=%v, want 32 true", n, ended)
	}
}

func TestFillSFXRespectsAudibleFrames(t *testing.T) {
	entry, data, cb := makeTestAsset(ClassSFX, 20, 32, false, 0, 0)
	var src SourceState
	InitSource(&src, entry, data, &cb, State{})

	dst := make([]int16, 64)
	n, ended := src.Fill(dst)

	if n != 20 || !ended {
		t.Fatalf("Fill returned n=%d ended=%v, want 20 true", n, ended)
	}
}

func TestFillMusicLoopsForward(t *testing.T) {
	entry, data, cb := makeTestAsset(ClassMusic, 48, 64, true, 16, 32)
	var src SourceState
	InitSource(&src, entry, data, &cb, captureLoopState(&cb, data, 16))

	dst := make([]int16, 200)
	total := 0
	for total < len(dst) {
		n, ended := src.Fill(dst[total:])
		total += n
		if ended {
			t.Fatalf("music source ended at total=%d", total)
		}
		if n == 0 {
			t.Fatal("Fill returned 0 samples without ending")
		}
	}
}

func TestFillAntiClickRampShape(t *testing.T) {
	entry, data, cb := makeTestAsset(ClassMusic, 128, 128, true, 0, 128)
	var src SourceState
	InitSource(&src, entry, data, &cb, State{})

	src.Fill(make([]int16, 32))
	src.RequestStop(8)

	rampBuf := make([]int16, 32)
	n, _ := src.Fill(rampBuf)
	if n < 8 {
		t.Fatalf("Fill after stop returned n=%d, want >= 8", n)
	}
	for i := 1; i < 8; i++ {
		if abs16(rampBuf[i]) > abs16(rampBuf[i-1]) {
			t.Fatalf("ramp not monotonic: |%d| > |%d|", rampBuf[i], rampBuf[i-1])
		}
	}
	for i := 8; i < n; i++ {
		if rampBuf[i] != 0 {
			t.Fatalf("post-ramp sample[%d] = %d, want 0", i, rampBuf[i])
		}
	}
	_, ended := src.Fill(make([]int16, 16))
	if !ended {
		t.Fatal("source did not end after ramp")
	}
}

func captureLoopState(cb *Codebook, data []byte, loopStartFrame uint32) State {
	var state State
	for b := 0; b < int(loopStartFrame/BlockSamples); b++ {
		var block [BlockBytes]byte
		copy(block[:], data[b*BlockBytes:])
		var out [BlockSamples]int16
		DecodeBlock(cb, &state, block, out[:])
	}
	return state
}

func abs16(v int16) int16 {
	if v < 0 {
		return -v
	}
	return v
}
