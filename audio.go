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

var (
	v1engine          *audiov1.Engine
	v1mixer           *audiov1.Mixer
	v1runtime         *audioV1Runtime
	v1dacBufFrames    = 512
	v1sfxNameResolver func(string) (uint16, bool)
)

const DefaultAudioOutputRate = 48000

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

	v1engine = &audiov1.Engine{
		Manifest:  entries,
		Data:      data,
		Aux:       aux,
		SFXGain:   audiov1.GainFull,
		MusicGain: audiov1.GainFull,
	}
	v1mixer = audiov1.NewMixer(DefaultAudioOutputRate, v1dacBufFrames)
	rt := &audioV1Runtime{
		outBuf:  make([]int16, v1dacBufFrames*2),
		outByte: make([]byte, v1dacBufFrames*4),
		taps:    make([]audiov1.VoiceTap, audiov1.MaxVoices),
	}
	for i := range rt.srcBufs {
		rt.srcBufs[i] = make([]int16, v1dacBufFrames+4)
	}
	v1runtime = rt
}

func RegisterSFXNameResolver(fn func(string) (uint16, bool)) {
	v1sfxNameResolver = fn
}

func PlayEffect(id sfx.ID) bool {
	if v1engine == nil || !v1engine.IsReady() {
		return false
	}
	return v1engine.Ring.Push(audiov1.Command{Kind: audiov1.CmdPlaySFX, ID: uint16(id)})
}

func PlayTrack(id music.ID) bool {
	if v1engine == nil || !v1engine.IsReady() {
		return false
	}
	return v1engine.Ring.Push(audiov1.Command{Kind: audiov1.CmdPlayMusic, ID: uint16(id)})
}

func StopTrack() {
	if v1engine == nil || !v1engine.IsReady() {
		return
	}
	v1engine.Ring.Push(audiov1.Command{Kind: audiov1.CmdStopMusic})
}

func SetEffectVolume(v float32) {
	if v1engine == nil || !v1engine.IsReady() {
		return
	}
	v1engine.Ring.Push(audiov1.Command{Kind: audiov1.CmdSetSFXGain, Gain: floatToGain(v)})
}

func SetTrackVolume(v float32) {
	if v1engine == nil || !v1engine.IsReady() {
		return
	}
	v1engine.Ring.Push(audiov1.Command{Kind: audiov1.CmdSetMusicGain, Gain: floatToGain(v)})
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

func initAudioV1() {
	if v1engine == nil {
		return
	}
	audio.Start(DefaultAudioOutputRate)
	v1engine.SetReady(true)
	go v1AudioFeeder()
}

func v1AudioFeeder() {
	rt := v1runtime
	for {
		v1engine.DrainCommands()
		clear(rt.outBuf)
		clear(rt.retire[:])

		for i := 0; i < audiov1.MaxVoices; i++ {
			tap := &rt.taps[i]
			tap.Active = false
			tap.Samples = nil
			tap.Consumed = 0

			voice := &v1engine.Voices[i]
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
			if v1engine.Sources[i].Stopping() && !rt.wasStopping[i] {
				rt.pending[i] = 0
				rt.wasStopping[i] = true
			}
			entry := &v1engine.Manifest[voice.ManifestIndex]
			need := audiov1.SourceFramesNeeded(uint32(entry.Rate), DefaultAudioOutputRate, v1dacBufFrames, voice.Phase)
			if need > len(rt.srcBufs[i]) {
				need = len(rt.srcBufs[i])
			}
			if rt.pending[i] > need {
				need = rt.pending[i]
			}
			if rt.pending[i] < need {
				n, ended := v1engine.Sources[i].Fill(rt.srcBufs[i][rt.pending[i]:need])
				rt.pending[i] += n
				rt.retire[i] = ended
			}

			if rt.pending[i] == 0 {
				if rt.retire[i] {
					voice.State = audiov1.VoiceIdle
				}
				continue
			}

			gain := v1engine.SFXGain
			if voice.Class == audiov1.ClassMusic {
				gain = v1engine.MusicGain
			}
			tap.Samples = rt.srcBufs[i][:rt.pending[i]]
			tap.SrcRate = uint32(entry.Rate)
			tap.Gain = gain
			tap.Active = true
			tap.Phase = voice.Phase
		}

		v1mixer.Mix(rt.outBuf, rt.taps)

		for i := 0; i < audiov1.MaxVoices; i++ {
			tap := &rt.taps[i]
			voice := &v1engine.Voices[i]
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
