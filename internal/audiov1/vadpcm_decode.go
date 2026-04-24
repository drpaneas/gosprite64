package audiov1

const (
	PredictorCount = 4
	Order          = 2
	BlockSamples   = 16
	BlockBytes     = 9
	StateLen       = 8
	CodebookInts   = PredictorCount * Order * StateLen
)

// Codebook stores predictor coefficients in Q11 fixed-point.
// Layout: [predictor_index][order_index][vector_position].
// order 0 = response to state[-2] (older history sample).
// order 1 = response to state[-1] (more recent history sample).
type Codebook [PredictorCount][Order][StateLen]int16

// State holds the last 8 decoded samples from the previous block.
// state[StateLen-Order] .. state[StateLen-1] are the history samples
// used for prediction (state[-2] and state[-1] respectively).
type State [StateLen]int16

func ClampInt16(v int32) int16 {
	if v > 32767 {
		return 32767
	}
	if v < -32768 {
		return -32768
	}
	return int16(v)
}

func signExtend4(x int32) int32 {
	if x > 7 {
		return x - 16
	}
	return x
}

// DecodeBlock decodes one 9-byte VADPCM block into 16 mono int16 samples.
// Control byte layout: high nibble = scale, low nibble = predictor index.
func DecodeBlock(cb *Codebook, state *State, block [BlockBytes]byte, out []int16) {
	control := block[0]
	scale := uint(control >> 4)
	pred := int(control & 0x0F)
	if pred >= PredictorCount {
		pred = 0
	}

	for vector := 0; vector < 2; vector++ {
		var accumulator [StateLen]int32

		for k := 0; k < Order; k++ {
			sample := int32(state[StateLen-Order+k])
			for i := 0; i < StateLen; i++ {
				accumulator[i] += sample * int32(cb[pred][k][i])
			}
		}

		var residuals [StateLen]int32
		for i := 0; i < 4; i++ {
			b := block[1+4*vector+i]
			residuals[2*i] = signExtend4(int32(b >> 4))
			residuals[2*i+1] = signExtend4(int32(b & 0x0F))
		}

		for k := 0; k < StateLen; k++ {
			residual := residuals[k] << scale
			accumulator[k] += residual << 11
			for i := 0; i < 7-k; i++ {
				accumulator[k+1+i] += residual * int32(cb[pred][Order-1][i])
			}
		}

		for i := 0; i < StateLen; i++ {
			sample := ClampInt16(accumulator[i] >> 11)
			out[vector*StateLen+i] = sample
			state[i] = sample
		}
	}
}
