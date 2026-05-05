package gosprite64

import (
	"errors"
	"log"

	"github.com/clktmr/n64/rcp/audio"
	"github.com/drpaneas/gosprite64/audio/music"
	"github.com/drpaneas/gosprite64/audio/sfx"
	"github.com/drpaneas/gosprite64/internal/audiov1"
)

// AudioAsset mirrors internal/audiov1.AssetEntry for use in generated code,
// which cannot import internal packages. Fields must stay in sync.
type AudioAsset struct {
	ID            uint16
	Class         uint8
	Flags         uint8
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
}

type audioConfig struct {
	manifest        []audiov1.AssetEntry
	data            []byte
	aux             []byte
	sfxNameResolver func(string) (uint16, bool)
}

type audioV1Runtime struct {
	outBuf      []int16
	outByte     []byte
	taps        []audiov1.VoiceTap
	srcBufs     [audiov1.MaxVoices][]int16
	pending     [audiov1.MaxVoices]int
	startSeq    [audiov1.MaxVoices]uint32
	wasStopping [audiov1.MaxVoices]bool
	retire      [audiov1.MaxVoices]bool
}

type audioState struct {
	engine          *audiov1.Engine
	mixer           *audiov1.Mixer
	runtime         audioV1Runtime
	dacBufFrames    int
	sfxNameResolver func(string) (uint16, bool)
}

const DefaultAudioOutputRate = 48000

const defaultAudioDACBufFrames = 512

var pendingAudioConfig audioConfig

func RegisterAudioV1(assets []AudioAsset, data, aux []byte) {
	entries := make([]audiov1.AssetEntry, len(assets))
	for i, a := range assets {
		entries[i] = audiov1.AssetEntry{
			ID:            audiov1.AssetID(a.ID),
			Class:         audiov1.AssetClass(a.Class),
			Flags:         audiov1.AssetFlags(a.Flags),
			Rate:          a.Rate,
			AudibleFrames: a.AudibleFrames,
			EncodedFrames: a.EncodedFrames,
			LoopStart:     a.LoopStart,
			LoopLen:       a.LoopLen,
			DataOffset:    a.DataOffset,
			DataBytes:     a.DataBytes,
			AuxOffset:     a.AuxOffset,
			AuxBytes:      a.AuxBytes,
			MaxInstances:  a.MaxInstances,
			BlockFrames:   audiov1.BlockSamples,
		}
	}

	pendingAudioConfig.manifest = entries
	pendingAudioConfig.data = data
	pendingAudioConfig.aux = aux
}

func RegisterSFXNameResolver(fn func(string) (uint16, bool)) {
	pendingAudioConfig.sfxNameResolver = fn
}

func newAudioState(cfg audioConfig) *audioState {
	a := &audioState{
		dacBufFrames:    defaultAudioDACBufFrames,
		sfxNameResolver: cfg.sfxNameResolver,
	}
	if len(cfg.manifest) == 0 {
		return a
	}

	a.engine = &audiov1.Engine{
		Manifest:  cfg.manifest,
		Data:      cfg.data,
		Aux:       cfg.aux,
		SFXGain:   audiov1.GainFull,
		MusicGain: audiov1.GainFull,
	}
	a.mixer = audiov1.NewMixer(DefaultAudioOutputRate, a.dacBufFrames)
	a.runtime.outBuf = make([]int16, a.dacBufFrames*2)
	a.runtime.outByte = make([]byte, a.dacBufFrames*4)
	a.runtime.taps = make([]audiov1.VoiceTap, audiov1.MaxVoices)
	for i := range a.runtime.srcBufs {
		a.runtime.srcBufs[i] = make([]int16, a.dacBufFrames+4)
	}
	return a
}

func (a *audioState) ready() bool {
	return a != nil && a.engine != nil && a.engine.IsReady()
}

func PlayEffect(id sfx.ID) bool {
	audio := currentAudio()
	if !audio.ready() {
		return false
	}
	return audio.engine.Ring.Push(audiov1.Command{Kind: audiov1.CmdPlaySFX, ID: uint16(id)})
}

func PlayTrack(id music.ID) bool {
	audio := currentAudio()
	if !audio.ready() {
		return false
	}
	return audio.engine.Ring.Push(audiov1.Command{Kind: audiov1.CmdPlayMusic, ID: uint16(id)})
}

func StopTrack() {
	audio := currentAudio()
	if !audio.ready() {
		return
	}
	audio.engine.Ring.Push(audiov1.Command{Kind: audiov1.CmdStopMusic})
}

func SetEffectVolume(v float32) {
	audio := currentAudio()
	if !audio.ready() {
		return
	}
	audio.engine.Ring.Push(audiov1.Command{Kind: audiov1.CmdSetSFXGain, Gain: floatToGain(v)})
}

func SetTrackVolume(v float32) {
	audio := currentAudio()
	if !audio.ready() {
		return
	}
	audio.engine.Ring.Push(audiov1.Command{Kind: audiov1.CmdSetMusicGain, Gain: floatToGain(v)})
}

func floatToGain(v float32) uint16 {
	if v <= 0 {
		return 0
	}
	if v >= 1 {
		return audiov1.GainFull
	}
	return uint16(v * float32(audiov1.GainFull))
}

func (rt *runtimeState) initAudio() {
	if rt == nil {
		return
	}
	rt.audio = newAudioState(pendingAudioConfig)
	rt.audio.start()
}

func (a *audioState) start() {
	if a == nil || a.engine == nil {
		return
	}
	audio.Start(DefaultAudioOutputRate)
	a.engine.SetReady(true)
	go a.feeder()
}

func (a *audioState) feeder() {
	if a == nil || a.engine == nil || a.mixer == nil {
		return
	}

	rt := &a.runtime
	for {
		a.engine.DrainCommands()
		clear(rt.outBuf)
		clear(rt.retire[:])

		for i := 0; i < audiov1.MaxVoices; i++ {
			tap := &rt.taps[i]
			tap.Active = false
			tap.Samples = nil
			tap.Consumed = 0

			voice := &a.engine.Voices[i]
			if voice.State != audiov1.VoicePlaying || voice.ManifestIndex < 0 {
				rt.pending[i] = 0
				rt.startSeq[i] = 0
				rt.wasStopping[i] = false
				continue
			}
			if rt.startSeq[i] != voice.StartSeq {
				rt.pending[i] = 0
				rt.startSeq[i] = voice.StartSeq
				rt.wasStopping[i] = false
			}
			if a.engine.Sources[i].Stopping() && !rt.wasStopping[i] {
				rt.pending[i] = 0
				rt.wasStopping[i] = true
			}
			entry := &a.engine.Manifest[voice.ManifestIndex]
			need := audiov1.SourceFramesNeeded(uint32(entry.Rate), DefaultAudioOutputRate, a.dacBufFrames, voice.Phase)
			if need > len(rt.srcBufs[i]) {
				need = len(rt.srcBufs[i])
			}
			if rt.pending[i] > need {
				need = rt.pending[i]
			}
			if rt.pending[i] < need {
				n, ended := a.engine.Sources[i].Fill(rt.srcBufs[i][rt.pending[i]:need])
				rt.pending[i] += n
				rt.retire[i] = ended
			}

			if rt.pending[i] == 0 {
				if rt.retire[i] {
					voice.State = audiov1.VoiceIdle
				}
				continue
			}

			gain := a.engine.SFXGain
			if voice.Class == audiov1.ClassMusic {
				gain = a.engine.MusicGain
			}
			tap.Samples = rt.srcBufs[i][:rt.pending[i]]
			tap.SrcRate = uint32(entry.Rate)
			tap.Gain = gain
			tap.Active = true
			tap.Phase = voice.Phase
		}

		a.mixer.Mix(rt.outBuf, rt.taps)

		for i := 0; i < audiov1.MaxVoices; i++ {
			tap := &rt.taps[i]
			voice := &a.engine.Voices[i]
			if !tap.Active {
				continue
			}
			consumed := tap.Consumed
			if consumed > rt.pending[i] {
				consumed = rt.pending[i]
			}
			if consumed > 0 {
				copy(rt.srcBufs[i], rt.srcBufs[i][consumed:rt.pending[i]])
				rt.pending[i] -= consumed
			}
			voice.Phase = tap.Phase
			if rt.retire[i] && rt.pending[i] == 0 {
				voice.State = audiov1.VoiceIdle
			}
		}

		writeV1DACOutput(rt.outBuf, rt.outByte)
	}
}

// writeV1DACOutput serializes stereo int16 samples to big-endian bytes and
// submits them to the N64 audio DMA buffer. The N64 AI DMA expects signed
// 16-bit big-endian stereo samples, which matches the byte packing here.
// audio.Buffer.Write blocks until the hardware has consumed enough prior
// data to accept the new buffer, which provides natural pacing at the DAC
// output rate.
func writeV1DACOutput(stereo []int16, buf []byte) {
	for i, s := range stereo {
		buf[i*2] = byte(s >> 8)
		buf[i*2+1] = byte(s)
	}
	_, err := audio.Buffer.Write(buf[:len(stereo)*2])
	if err != nil && !errors.Is(err, audio.ErrStop) {
		log.Printf("audio v1 write stopped: %v", err)
	}
}
