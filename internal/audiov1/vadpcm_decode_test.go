package audiov1

import "testing"

func TestDecodeBlockAllZero(t *testing.T) {
	var cb Codebook
	var state State
	block := [BlockBytes]byte{}
	var out [BlockSamples]int16

	DecodeBlock(&cb, &state, block, out[:])

	for i, s := range out {
		if s != 0 {
			t.Fatalf("out[%d] = %d, want 0", i, s)
		}
	}
}

func TestDecodeBlockWithResiduals(t *testing.T) {
	var cb Codebook
	var state State
	block := [BlockBytes]byte{0x30, 0x72}
	var out [BlockSamples]int16

	DecodeBlock(&cb, &state, block, out[:])

	want := [BlockSamples]int16{56, 16}
	for i := range out {
		if out[i] != want[i] {
			t.Fatalf("out[%d] = %d, want %d", i, out[i], want[i])
		}
	}
}

func TestDecodeBlockNegativeNibbles(t *testing.T) {
	var cb Codebook
	var state State
	block := [BlockBytes]byte{0x20, 0xF0}
	var out [BlockSamples]int16

	DecodeBlock(&cb, &state, block, out[:])

	if out[0] != -4 {
		t.Fatalf("out[0] = %d, want -4", out[0])
	}
	if out[1] != 0 {
		t.Fatalf("out[1] = %d, want 0", out[1])
	}
}

func TestDecodeBlockForwardFeedingWithinSubvector(t *testing.T) {
	var cb Codebook
	cb[0][1][0] = 2048
	var state State
	block := [BlockBytes]byte{0x00, 0x40}
	var out [BlockSamples]int16

	DecodeBlock(&cb, &state, block, out[:])

	if out[0] != 4 {
		t.Fatalf("out[0] = %d, want 4", out[0])
	}
	if out[1] != 4 {
		t.Fatalf("out[1] = %d, want 4", out[1])
	}
	if out[2] != 0 {
		t.Fatalf("out[2] = %d, want 0", out[2])
	}
}

func TestDecodeBlockPreservesStateBetweenCalls(t *testing.T) {
	var cb Codebook
	for i := 0; i < StateLen; i++ {
		cb[0][1][i] = 2048
	}
	state := State{0, 0, 0, 0, 0, 0, 0, 300}
	block := [BlockBytes]byte{}
	var out [BlockSamples]int16

	DecodeBlock(&cb, &state, block, out[:])
	DecodeBlock(&cb, &state, block, out[:])

	for i, s := range out {
		if s != 300 {
			t.Fatalf("second block out[%d] = %d, want 300", i, s)
		}
	}
}
