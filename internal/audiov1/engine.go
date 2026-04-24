package audiov1

import "sync/atomic"

type AssetID uint16

type AssetClass uint8

const (
	ClassSFX   AssetClass = 0
	ClassMusic AssetClass = 1
)

type AssetFlags uint8

const (
	FlagResident AssetFlags = 1 << iota
	FlagStreamed
	FlagLoop
)

type AssetEntry struct {
	ID            AssetID
	Class         AssetClass
	Flags         AssetFlags
	Rate          uint16
	AudibleFrames uint32
	EncodedFrames uint32
	LoopStart     uint32
	LoopLen       uint32
	DataOffset    uint32
	DataBytes     uint32
	AuxOffset     uint32
	AuxBytes      uint32
	MaxInstances  uint8
	BlockFrames   uint8
}

const (
	MaxSFXVoices    = 8
	MaxMusicVoices  = 1
	MaxVoices       = MaxSFXVoices + MaxMusicVoices
	DefaultMaxInst  = 4
	MusicVoiceIndex = 0
	FirstSFXVoice   = 1
	GainFull        = 32768
	cmdRingSize     = 16
)

type VoiceState uint8

const (
	VoiceIdle VoiceState = iota
	VoicePlaying
)

type Voice struct {
	State         VoiceState
	AssetID       AssetID
	Class         AssetClass
	StartSeq      uint32
	ManifestIndex int16
	Phase         uint64
}

type CmdKind uint8

const (
	CmdPlaySFX CmdKind = iota
	CmdPlayMusic
	CmdStopMusic
	CmdSetSFXGain
	CmdSetMusicGain
)

type Command struct {
	Kind CmdKind
	ID   uint16
	Gain uint16
}

// CommandRing is a single-producer single-consumer (SPSC) lock-free ring buffer.
// CONCURRENCY CONTRACT:
//   - Only the gameplay goroutine may call Push (producer).
//   - Only the audio feeder goroutine may call Pop (consumer).
//   - If multiple goroutines need to call PlayEffect, external synchronization
//     is required, or all gameplay calls must be funneled through one goroutine.
type CommandRing struct {
	buf  [cmdRingSize]Command
	head atomic.Uint32
	tail atomic.Uint32
}

type Engine struct {
	Voices    [MaxVoices]Voice
	Sources   [MaxVoices]SourceState
	Ring      CommandRing
	SFXGain   uint16
	MusicGain uint16
	Seq       uint32
	Manifest  []AssetEntry
	Data      []byte
	Aux       []byte
	ready     atomic.Bool
}

func (r *CommandRing) Push(cmd Command) bool {
	head := r.head.Load()
	next := (head + 1) % cmdRingSize
	if next == r.tail.Load() {
		return false
	}
	r.buf[head] = cmd
	r.head.Store(next)
	return true
}

func (r *CommandRing) Pop() (Command, bool) {
	tail := r.tail.Load()
	if tail == r.head.Load() {
		return Command{}, false
	}
	cmd := r.buf[tail]
	r.tail.Store((tail + 1) % cmdRingSize)
	return cmd, true
}

func (e *Engine) SetReady(r bool) {
	e.ready.Store(r)
}

func (e *Engine) IsReady() bool {
	return e.ready.Load()
}

func (e *Engine) FindManifestEntry(class AssetClass, id uint16) (*AssetEntry, int) {
	for i := range e.Manifest {
		if e.Manifest[i].Class == class && uint16(e.Manifest[i].ID) == id {
			return &e.Manifest[i], i
		}
	}
	return nil, -1
}

// sfxVoiceEnd is the exclusive upper bound for SFX voice indices.
const sfxVoiceEnd = FirstSFXVoice + MaxSFXVoices

func (e *Engine) AllocateSFXVoice(id AssetID, maxInstances uint8) int {
	if maxInstances == 0 {
		maxInstances = DefaultMaxInst
	}

	var sameCount int
	oldestSameIdx := -1
	var oldestSameSeq uint32 = ^uint32(0)

	for i := FirstSFXVoice; i < sfxVoiceEnd; i++ {
		v := &e.Voices[i]
		if v.State == VoicePlaying && v.AssetID == id {
			sameCount++
			if v.StartSeq < oldestSameSeq {
				oldestSameSeq = v.StartSeq
				oldestSameIdx = i
			}
		}
	}

	if sameCount >= int(maxInstances) && oldestSameIdx >= 0 {
		return oldestSameIdx
	}

	for i := FirstSFXVoice; i < sfxVoiceEnd; i++ {
		if e.Voices[i].State == VoiceIdle {
			return i
		}
	}

	oldestIdx := FirstSFXVoice
	var oldestSeq uint32 = ^uint32(0)
	for i := FirstSFXVoice; i < sfxVoiceEnd; i++ {
		if e.Voices[i].StartSeq < oldestSeq {
			oldestSeq = e.Voices[i].StartSeq
			oldestIdx = i
		}
	}
	return oldestIdx
}

func (e *Engine) DrainCommands() {
	for {
		cmd, ok := e.Ring.Pop()
		if !ok {
			break
		}
		switch cmd.Kind {
		case CmdPlaySFX:
			entry, manifestIndex := e.FindManifestEntry(ClassSFX, cmd.ID)
			if entry == nil {
				continue
			}
			e.Seq++
			idx := e.AllocateSFXVoice(AssetID(cmd.ID), entry.MaxInstances)
			e.startVoice(idx, entry, manifestIndex)

		case CmdPlayMusic:
			if e.Voices[MusicVoiceIndex].State == VoicePlaying &&
				uint16(e.Voices[MusicVoiceIndex].AssetID) == cmd.ID {
				continue
			}
			entry, manifestIndex := e.FindManifestEntry(ClassMusic, cmd.ID)
			if entry == nil {
				continue
			}
			e.Seq++
			e.startVoice(MusicVoiceIndex, entry, manifestIndex)

		case CmdStopMusic:
			if e.Voices[MusicVoiceIndex].State == VoicePlaying {
				e.Sources[MusicVoiceIndex].RequestStop(64)
			}

		case CmdSetSFXGain:
			e.SFXGain = cmd.Gain

		case CmdSetMusicGain:
			e.MusicGain = cmd.Gain
		}
	}
}

func (e *Engine) startVoice(idx int, entry *AssetEntry, manifestIndex int) {
	data := e.assetData(entry)
	codebook, loopState := e.assetAux(entry)
	InitSource(&e.Sources[idx], entry, data, &codebook, loopState)
	e.Voices[idx] = Voice{
		State:         VoicePlaying,
		AssetID:       entry.ID,
		Class:         entry.Class,
		StartSeq:      e.Seq,
		ManifestIndex: int16(manifestIndex),
		Phase:         0,
	}
}

func (e *Engine) assetData(entry *AssetEntry) []byte {
	start := int(entry.DataOffset)
	end := start + int(entry.DataBytes)
	if start < 0 || end > len(e.Data) || start > end {
		return nil
	}
	return e.Data[start:end]
}

func (e *Engine) assetAux(entry *AssetEntry) (Codebook, State) {
	var cb Codebook
	var loop State
	start := int(entry.AuxOffset)
	end := start + int(entry.AuxBytes)
	if start < 0 || end > len(e.Aux) || start > end {
		return cb, loop
	}
	aux := e.Aux[start:end]
	codebookFromBytes(aux, &cb)
	if len(aux) >= CodebookInts*2+StateLen*2 {
		stateFromBytes(aux[CodebookInts*2:], &loop)
	}
	return cb, loop
}

func codebookFromBytes(data []byte, cb *Codebook) {
	if len(data) < CodebookInts*2 {
		return
	}
	idx := 0
	for p := 0; p < PredictorCount; p++ {
		for o := 0; o < Order; o++ {
			for s := 0; s < StateLen; s++ {
				cb[p][o][s] = int16(uint16(data[idx])<<8 | uint16(data[idx+1]))
				idx += 2
			}
		}
	}
}

func stateFromBytes(data []byte, state *State) {
	if len(data) < StateLen*2 {
		return
	}
	for i := 0; i < StateLen; i++ {
		state[i] = int16(uint16(data[i*2])<<8 | uint16(data[i*2+1]))
	}
}
