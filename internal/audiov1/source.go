package audiov1

type SourceState struct {
	asset     *AssetEntry
	data      []byte
	codebook  Codebook
	loopState State

	decState   State
	cursor     uint32
	dataOffset uint32

	decodeBuf [BlockSamples]int16
	bufStart  uint32
	bufValid  int
	rampPos   int16
	rampLen   int16
	stopping  bool
	ended     bool
}

func InitSource(s *SourceState, asset *AssetEntry, data []byte, codebook *Codebook, loopState State) {
	s.asset = asset
	s.data = data
	if codebook != nil {
		s.codebook = *codebook
	} else {
		s.codebook = Codebook{}
	}
	s.loopState = loopState
	s.Reset()
}

func (s *SourceState) Reset() {
	s.decState = State{}
	s.cursor = 0
	s.dataOffset = 0
	s.bufStart = 0
	s.bufValid = 0
	s.rampPos = 0
	s.rampLen = 0
	s.stopping = false
	s.ended = false
}

func (s *SourceState) RequestStop(rampSamples int16) {
	if !s.stopping {
		s.stopping = true
		s.rampPos = 0
		s.rampLen = rampSamples
	}
}

func (s *SourceState) Stopping() bool {
	return s.stopping
}

func (s *SourceState) Fill(dst []int16) (n int, ended bool) {
	if s.ended {
		return 0, true
	}
	if s.asset == nil {
		s.ended = true
		return 0, true
	}

	limit := s.asset.AudibleFrames
	if s.asset.Flags&FlagLoop != 0 {
		limit = s.asset.LoopStart + s.asset.LoopLen
	}

	for n < len(dst) {
		if s.cursor >= limit {
			if s.asset.Flags&FlagLoop != 0 && !s.stopping {
				s.wrapLoop()
				continue
			}
			s.ended = true
			break
		}

		avail := s.bufValid - int(s.cursor-s.bufStart)
		if avail <= 0 {
			s.decodeNextBlock()
			avail = s.bufValid - int(s.cursor-s.bufStart)
			if avail <= 0 {
				s.ended = true
				break
			}
		}

		remaining := int(limit - s.cursor)
		canWrite := len(dst) - n
		if avail > canWrite {
			avail = canWrite
		}
		if avail > remaining {
			avail = remaining
		}

		bufOff := int(s.cursor - s.bufStart)
		copy(dst[n:n+avail], s.decodeBuf[bufOff:bufOff+avail])

		if s.stopping {
			s.applyRamp(dst[n : n+avail])
		}

		n += avail
		s.cursor += uint32(avail)

		if s.stopping && s.rampPos >= s.rampLen {
			s.ended = true
			break
		}
	}

	if !s.ended && s.asset.Flags&FlagLoop == 0 && s.cursor >= limit {
		s.ended = true
	}

	return n, s.ended
}

func (s *SourceState) decodeNextBlock() {
	blockStart := s.dataOffset
	if int(blockStart)+BlockBytes > len(s.data) {
		s.bufValid = 0
		return
	}

	var block [BlockBytes]byte
	copy(block[:], s.data[blockStart:blockStart+BlockBytes])
	DecodeBlock(&s.codebook, &s.decState, block, s.decodeBuf[:])

	s.bufStart = (s.dataOffset / BlockBytes) * BlockSamples
	s.bufValid = BlockSamples
	s.dataOffset += BlockBytes
}

func (s *SourceState) wrapLoop() {
	s.decState = s.loopState
	s.cursor = s.asset.LoopStart
	s.dataOffset = (s.asset.LoopStart / BlockSamples) * BlockBytes
	s.bufValid = 0
}

func (s *SourceState) applyRamp(samples []int16) {
	if s.rampLen <= 0 {
		clear(samples)
		s.rampPos = s.rampLen
		return
	}
	for i := range samples {
		if s.rampPos >= s.rampLen {
			samples[i] = 0
			continue
		}
		fade := int32(s.rampLen-s.rampPos) * 32768 / int32(s.rampLen)
		samples[i] = int16((int32(samples[i]) * fade) >> 15)
		s.rampPos++
	}
}
