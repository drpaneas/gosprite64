package audiov1

const phaseFracMask = uint64(1<<32 - 1)

// VoiceTap is the per-voice input to the mixer for a single Mix call.
// Phase stores only the fractional 32.32 source position. Consumed reports
// how many whole source frames were consumed from Samples during the call.
type VoiceTap struct {
	Samples  []int16
	SrcRate  uint32
	Gain     uint16
	Active   bool
	Phase    uint64
	Consumed int
}

type Mixer struct {
	OutRate uint32
	accum   []int32
}

func NewMixer(outRate uint32, maxOutFrames int) *Mixer {
	return &Mixer{
		OutRate: outRate,
		accum:   make([]int32, maxOutFrames),
	}
}

func (m *Mixer) Mix(out []int16, taps []VoiceTap) {
	outFrames := len(out) / 2
	if outFrames == 0 {
		return
	}
	if outFrames > len(m.accum) {
		outFrames = len(m.accum)
		out = out[:outFrames*2]
	}

	acc := m.accum[:outFrames]
	clear(acc)

	for i := range taps {
		taps[i].Consumed = 0
		if !taps[i].Active || len(taps[i].Samples) == 0 || taps[i].Gain == 0 {
			continue
		}
		m.mixVoice(acc, &taps[i])
	}

	for i := 0; i < outFrames; i++ {
		out[i*2] = ClampInt16(acc[i])
		out[i*2+1] = out[i*2]
	}
}

func (m *Mixer) mixVoice(acc []int32, tap *VoiceTap) {
	gain32 := int32(tap.Gain)
	srcLen := len(tap.Samples)
	if tap.SrcRate == 0 {
		return
	}

	if tap.SrcRate == m.OutRate {
		n := srcLen
		if n > len(acc) {
			n = len(acc)
		}
		for i := 0; i < n; i++ {
			acc[i] += (int32(tap.Samples[i]) * gain32) >> 15
		}
		tap.Consumed = n
		tap.Phase = 0
		return
	}

	step := (uint64(tap.SrcRate) << 32) / uint64(m.OutRate)
	phase := tap.Phase & phaseFracMask

	for i := 0; i < len(acc); i++ {
		idx := int(phase >> 32)
		if idx >= srcLen {
			break
		}
		frac := int32((phase >> 17) & 0x7FFF)

		s0 := int32(tap.Samples[idx])
		s1 := s0
		if idx+1 < srcLen {
			s1 = int32(tap.Samples[idx+1])
		}

		sample := s0 + ((s1-s0)*frac)>>15
		acc[i] += (sample * gain32) >> 15
		phase += step
	}

	tap.Consumed = int(phase >> 32)
	if tap.Consumed > srcLen {
		tap.Consumed = srcLen
	}
	tap.Phase = phase & phaseFracMask
}

func SourceFramesNeeded(srcRate, outRate uint32, outFrames int, phase uint64) int {
	if outFrames <= 0 || srcRate == 0 || outRate == 0 {
		return 0
	}
	if srcRate == outRate {
		return outFrames
	}
	step := (uint64(srcRate) << 32) / uint64(outRate)
	endPhase := (phase & phaseFracMask) + uint64(outFrames-1)*step
	return int(endPhase>>32) + 2
}
