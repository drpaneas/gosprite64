//go:build !noos

package main

import (
	"math"

	"github.com/drpaneas/gosprite64/internal/audiov1"
)

type TrainedCodebook struct {
	Book audiov1.Codebook
}

type EncodedAsset struct {
	Codebook   audiov1.Codebook
	Data       []byte
	FrameCount int
}

func autocorrelation(samples []int16, order int) []float64 {
	r := make([]float64, order+1)
	for k := 0; k <= order; k++ {
		var sum float64
		for i := 0; i < len(samples)-k; i++ {
			sum += float64(samples[i]) * float64(samples[i+k])
		}
		if len(samples) > k {
			sum /= float64(len(samples) - k)
		}
		r[k] = sum
	}
	return r
}

func levinsonDurbin(r []float64, order int) []float64 {
	if len(r) == 0 || r[0] == 0 {
		return make([]float64, order)
	}

	a := make([]float64, order)
	e := r[0]

	for m := 0; m < order; m++ {
		var sum float64
		for j := 0; j < m; j++ {
			sum += a[j] * r[m-j]
		}
		k := (r[m+1] - sum) / e
		if k > 0.999 {
			k = 0.999
		}
		if k < -0.999 {
			k = -0.999
		}

		newA := make([]float64, order)
		copy(newA, a)
		newA[m] = k
		for j := 0; j < m; j++ {
			newA[j] = a[j] - k*a[m-1-j]
		}
		copy(a, newA)
		e *= 1 - k*k
		if e <= r[0]*0.01 {
			break
		}
	}
	return a
}

func makeVectors(lpc []float64) [audiov1.Order][audiov1.StateLen]int16 {
	var vectors [audiov1.Order][audiov1.StateLen]int16
	scale := float64(1 << 11)

	for i := 0; i < audiov1.Order; i++ {
		x1 := 0.0
		x2 := 0.0
		if i == 0 {
			x2 = scale
		} else {
			x1 = scale
		}

		for j := 0; j < audiov1.StateLen; j++ {
			x := lpc[0]*x1 + lpc[1]*x2
			v := math.Round(x)
			if v > 32767 {
				v = 32767
			}
			if v < -32768 {
				v = -32768
			}
			vectors[i][j] = int16(v)
			x2 = x1
			x1 = x
		}
	}
	return vectors
}

type autocorrFrame [audiov1.Order + 1]float64

func autocorrBlock(samples []int16) autocorrFrame {
	var r autocorrFrame
	for k := 0; k <= audiov1.Order; k++ {
		var sum float64
		for i := 0; i < len(samples)-k; i++ {
			sum += float64(samples[i]) * float64(samples[i+k])
		}
		if len(samples) > k {
			sum /= float64(len(samples) - k)
		}
		r[k] = sum
	}
	return r
}

func trainCodebook(samples []int16) TrainedCodebook {
	numBlocks := len(samples) / audiov1.BlockSamples
	if numBlocks == 0 {
		return TrainedCodebook{}
	}

	corrs := make([]autocorrFrame, numBlocks)
	for b := 0; b < numBlocks; b++ {
		frame := samples[b*audiov1.BlockSamples : (b+1)*audiov1.BlockSamples]
		corrs[b] = autocorrBlock(frame)
	}

	assignments := make([]int, numBlocks)
	assignPredictors(corrs, assignments, audiov1.PredictorCount)

	var result TrainedCodebook
	for p := 0; p < audiov1.PredictorCount; p++ {
		var avgCorr autocorrFrame
		count := 0
		for b := 0; b < numBlocks; b++ {
			if assignments[b] == p {
				for k := 0; k <= audiov1.Order; k++ {
					avgCorr[k] += corrs[b][k]
				}
				count++
			}
		}
		if count == 0 {
			continue
		}
		for k := 0; k <= audiov1.Order; k++ {
			avgCorr[k] /= float64(count)
		}

		r := []float64{avgCorr[0], avgCorr[1], avgCorr[2]}
		lpc := levinsonDurbin(r, audiov1.Order)
		result.Book[p] = makeVectors(lpc)
	}
	return result
}

func assignPredictors(corrs []autocorrFrame, assignments []int, k int) {
	n := len(corrs)
	if n == 0 || k == 0 {
		return
	}

	centers := make([]autocorrFrame, k)
	used := 0
	for _, c := range corrs {
		if used >= k {
			break
		}
		dup := false
		for j := 0; j < used; j++ {
			if centers[j] == c {
				dup = true
				break
			}
		}
		if !dup {
			centers[used] = c
			used++
		}
	}
	for i := used; i < k; i++ {
		centers[i] = centers[0]
	}

	for iter := 0; iter < 50; iter++ {
		changed := false
		for i, c := range corrs {
			best := 0
			bestDist := math.Inf(1)
			for p := 0; p < k; p++ {
				d := corrDist(c, centers[p])
				if d < bestDist {
					bestDist = d
					best = p
				}
			}
			if assignments[i] != best {
				assignments[i] = best
				changed = true
			}
		}
		if !changed {
			break
		}

		counts := make([]int, k)
		newCenters := make([]autocorrFrame, k)
		for i, c := range corrs {
			p := assignments[i]
			counts[p]++
			for j := 0; j <= audiov1.Order; j++ {
				newCenters[p][j] += c[j]
			}
		}
		for p := 0; p < k; p++ {
			if counts[p] == 0 {
				newCenters[p] = centers[p]
				continue
			}
			for j := 0; j <= audiov1.Order; j++ {
				newCenters[p][j] /= float64(counts[p])
			}
		}
		centers = newCenters
	}
}

func corrDist(a, b autocorrFrame) float64 {
	var sum float64
	for i := 0; i <= audiov1.Order; i++ {
		d := a[i] - b[i]
		sum += d * d
	}
	return sum
}

func encodeBlock(cb *audiov1.Codebook, encState *audiov1.State, input []int16) [audiov1.BlockBytes]byte {
	bestPred := 0
	bestShift := 0
	bestError := math.MaxFloat64
	var bestNibbles [audiov1.BlockSamples]int32

	for pred := 0; pred < audiov1.PredictorCount; pred++ {
		pvec := &cb[pred]
		shift := computeShiftForPred(pvec, encState, input)

		minShift := shift - 1
		if minShift < 0 {
			minShift = 0
		}
		maxShift := shift + 1
		if maxShift > 12 {
			maxShift = 12
		}

		for tryShift := minShift; tryShift <= maxShift; tryShift++ {
			nibbles, totalError, _ := tryEncode(pvec, encState, input, tryShift)
			if totalError < bestError {
				bestError = totalError
				bestPred = pred
				bestShift = tryShift
				bestNibbles = nibbles
			}
		}
	}

	block := packBlock(bestShift, bestPred, bestNibbles)

	// Run the real decoder on the chosen block to get the canonical state.
	// This ensures encoder and decoder state stay perfectly synchronized,
	// including all 8 positions of the State vector, not just the 2 history
	// samples tracked by tryEncode.
	var decoded [audiov1.BlockSamples]int16
	audiov1.DecodeBlock(cb, encState, block, decoded[:])

	return block
}

func computeShiftForPred(pvec *[audiov1.Order][audiov1.StateLen]int16, encState *audiov1.State, input []int16) int {
	s0 := int32(encState[audiov1.StateLen-2])
	s1 := int32(encState[audiov1.StateLen-1])

	minVal, maxVal := int32(0), int32(0)

	for vector := 0; vector < 2; vector++ {
		var accumulator [audiov1.StateLen]int32
		for i := 0; i < audiov1.StateLen; i++ {
			accumulator[i] = s0*int32(pvec[0][i]) + s1*int32(pvec[1][i])
		}

		for i := 0; i < audiov1.StateLen; i++ {
			src := int32(input[vector*audiov1.StateLen+i])
			res := src*(1<<11) - accumulator[i]
			s := res >> 11
			if s < minVal {
				minVal = s
			}
			if s > maxVal {
				maxVal = s
			}
			for j := 0; j < 7-i; j++ {
				accumulator[i+1+j] -= s * int32(pvec[audiov1.Order-1][j])
			}
			s0 = s1
			s1 = src
		}
	}

	shift := 0
	for shift < 12 && (minVal < -8 || maxVal > 7) {
		minVal >>= 1
		maxVal >>= 1
		shift++
	}
	return shift
}

func tryEncode(pvec *[audiov1.Order][audiov1.StateLen]int16, encState *audiov1.State, input []int16, shift int) ([audiov1.BlockSamples]int32, float64, audiov1.State) {
	var nibbles [audiov1.BlockSamples]int32
	var totalError float64
	var finalState audiov1.State

	s0 := int32(encState[audiov1.StateLen-2])
	s1 := int32(encState[audiov1.StateLen-1])

	for vector := 0; vector < 2; vector++ {
		var accumulator [audiov1.StateLen]int32
		for i := 0; i < audiov1.StateLen; i++ {
			accumulator[i] = s0*int32(pvec[0][i]) + s1*int32(pvec[1][i])
		}

		for i := 0; i < audiov1.StateLen; i++ {
			idx := vector*audiov1.StateLen + i
			s := int32(input[idx])
			a := accumulator[i] >> 11

			r := (s - a) >> shift
			if r > 7 {
				r = 7
			}
			if r < -8 {
				r = -8
			}
			nibbles[idx] = r

			sout := r * (1 << shift)
			for j := 0; j < 7-i; j++ {
				accumulator[i+1+j] += sout * int32(pvec[audiov1.Order-1][j])
			}
			sout += a
			if sout > 32767 {
				sout = 32767
			}
			if sout < -32768 {
				sout = -32768
			}

			s0 = s1
			s1 = sout

			diff := float64(s - sout)
			totalError += diff * diff

			finalState[i] = int16(sout)
		}
	}

	return nibbles, totalError, finalState
}

func packBlock(scale, pred int, nibbles [audiov1.BlockSamples]int32) [audiov1.BlockBytes]byte {
	var block [audiov1.BlockBytes]byte
	block[0] = byte(scale<<4) | byte(pred&0x0F)
	for i := 0; i < audiov1.BlockSamples; i += 2 {
		hi := byte(nibbles[i]) & 0x0F
		lo := byte(nibbles[i+1]) & 0x0F
		block[1+i/2] = (hi << 4) | lo
	}
	return block
}

func EncodeVADPCM(samples []int16) EncodedAsset {
	numBlocks := (len(samples) + audiov1.BlockSamples - 1) / audiov1.BlockSamples
	padded := make([]int16, numBlocks*audiov1.BlockSamples)
	copy(padded, samples)

	trained := trainCodebook(padded)
	data := make([]byte, 0, numBlocks*audiov1.BlockBytes)
	var state audiov1.State

	for b := 0; b < numBlocks; b++ {
		frame := padded[b*audiov1.BlockSamples : (b+1)*audiov1.BlockSamples]
		block := encodeBlock(&trained.Book, &state, frame)
		data = append(data, block[:]...)
	}

	return EncodedAsset{
		Codebook:   trained.Book,
		Data:       data,
		FrameCount: len(padded),
	}
}

func captureStateAt(cb *audiov1.Codebook, data []byte, frameIndex int) audiov1.State {
	blockIndex := frameIndex / audiov1.BlockSamples
	var state audiov1.State
	for b := 0; b < blockIndex; b++ {
		var block [audiov1.BlockBytes]byte
		copy(block[:], data[b*audiov1.BlockBytes:])
		var out [audiov1.BlockSamples]int16
		audiov1.DecodeBlock(cb, &state, block, out[:])
	}
	return state
}
