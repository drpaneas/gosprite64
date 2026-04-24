package gosprite64

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"sync/atomic"

	"github.com/clktmr/n64/drivers/rspq"
	"github.com/clktmr/n64/drivers/rspq/mixer"
	"github.com/clktmr/n64/rcp/audio"
	"github.com/drpaneas/gosprite64/audio/music"
	"github.com/drpaneas/gosprite64/audio/sfx"
	"github.com/drpaneas/gosprite64/internal/audioengine"
	"github.com/drpaneas/gosprite64/internal/audiov1"
)

// audioFS is the global embedded filesystem for runtime audio assets.
// Runtime audio data is signed 16-bit stereo PCM at 48 kHz, stored as
// big-endian interleaved left/right samples.
// This is set by the generated audio_embed.go file in each game.
var (
	audioFS         embed.FS
	audioFSInit     bool      // Flag to track if audioFS has been initialized
	initializedOnce sync.Once // Ensures initialization happens only once
)

var (
	audioPlayerInstance *audioPlayer
	audioOnce           sync.Once
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

type audioPlayer struct {
	registry *audioengine.Registry
	playback *audioengine.PlaybackState
	sources  map[int]*playbackSource
	mutex    sync.Mutex
}

type playbackSource struct {
	channel  int
	loop     bool
	finished atomic.Bool
}

type trackedReadSeeker struct {
	io.ReadSeeker
	finished *atomic.Bool
}

func (t *trackedReadSeeker) Read(p []byte) (int, error) {
	n, err := t.ReadSeeker.Read(p)
	if errors.Is(err, io.EOF) {
		t.finished.Store(true)
	}
	return n, err
}

func getAudioPlayer() *audioPlayer {
	audioOnce.Do(func() {
		audioPlayerInstance = &audioPlayer{
			registry: audioengine.NewRegistry(),
			playback: audioengine.NewPlaybackState(),
			sources:  make(map[int]*playbackSource),
		}

		audioengine.StartMixerRuntime(audioengine.MixerRuntimeHooks{
			ResetQueue: rspq.Reset,
			InitMixer:  mixer.Init,
			StartDAC:   audio.Start,
			SetMixerRate: func(rate uint) {
				mixer.SetSampleRate(rate)
			},
			StartFeeder: func() {
				go func() {
					if _, err := audio.Buffer.ReadFrom(mixer.Output); err != nil && !errors.Is(err, audio.ErrStop) {
						log.Printf("audio feeder stopped: %v", err)
					}
				}()
			},
		})
	})
	return audioPlayerInstance
}

func (a *audioPlayer) update() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for trackID, source := range a.sources {
		if source.loop || !source.finished.Load() {
			continue
		}

		channel, ok := a.playback.Release(trackID)
		if ok {
			mixer.SetSource(channel, nil)
		}
		delete(a.sources, trackID)
	}
}

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

const DefaultAudioOutputRate = 48000

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

// Music plays or stops music
// If n is -1, stops all music
// Deprecated: Use PlayTrack with a generated music.ID constant instead.
func Music(n int, loop bool) {
	if v1engine != nil {
		if n == -1 {
			StopTrack()
			return
		}
		if n >= 0 {
			PlayTrack(music.ID(n))
		}
		return
	}

	ap := getAudioPlayer()

	if n == -1 {
		ap.mutex.Lock()
		defer ap.mutex.Unlock()

		channels := ap.playback.StopAll()
		clear(ap.sources)
		for _, channel := range channels {
			mixer.SetSource(channel, nil)
		}
		return
	}

	ap.mutex.Lock()
	defer ap.mutex.Unlock()

	data, err := ap.registry.Load(n, loadAudioDataByFilename)
	if err != nil {
		log.Printf("Failed to load audio track %d: %v", n, err)
		return
	}

	playback, err := ap.playback.Activate(n, loop)
	if err != nil {
		log.Printf("Failed to activate audio track %d: %v", n, err)
		return
	}

	playSource := &playbackSource{
		channel: playback.Channel,
		loop:    loop,
	}
	var reader io.ReadSeeker = bytes.NewReader(data)
	if loop {
		reader = mixer.Loop(reader)
	} else {
		reader = &trackedReadSeeker{ReadSeeker: reader, finished: &playSource.finished}
	}

	source := mixer.NewSource(reader, uint(audioengine.RuntimeSampleRate))
	ap.sources[n] = playSource
	mixer.SetSource(playback.Channel, source)
}

// LoadAudio loads raw runtime PCM audio data for a track.
// pcmData must be signed 16-bit stereo PCM at 48 kHz, stored as big-endian
// interleaved left/right samples.
func LoadAudio(id int, pcmData []byte) {
	ap := getAudioPlayer()
	ap.mutex.Lock()
	defer ap.mutex.Unlock()
	ap.registry.StorePCM(id, pcmData)
}

// LoadAudioFile registers an audio file to be loaded on demand
func LoadAudioFile(id int, filename string) error {
	// Check if the audio filesystem is properly initialized
	if !audioFSInit {
		log.Printf("ERROR: audioFS has not been initialized when trying to register %s", filename)
		return fmt.Errorf("audio filesystem not initialized")
	}

	ap := getAudioPlayer()
	ap.mutex.Lock()
	defer ap.mutex.Unlock()
	ap.registry.RegisterFile(id, filename)
	return nil
}

// loadAudioDataByFilename loads the actual audio data for a track.
func loadAudioDataByFilename(filename string) ([]byte, error) {
	f, err := audioFS.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file %s: %w", filename, err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio file %s: %w", filename, err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("audio file %s is empty", filename)
	}

	return data, nil
}

// PlaySFX plays a sound effect by name (without the "sfx_" prefix and ".raw" extension)
// Example: PlaySFX("jump") will play "sfx_jump.raw" if it exists
// Deprecated: Use PlayEffect with a generated sfx.ID constant instead.
func PlaySFX(name string) {
	if v1engine != nil {
		if v1sfxNameResolver != nil {
			if id, ok := v1sfxNameResolver(name); ok {
				PlayEffect(sfx.ID(id))
			}
		}
		return
	}

	cue, ok := audioengine.ResolveSFX(name)
	if !ok {
		return
	}

	filename := cue.Filename
	id := cue.TrackID

	player := getAudioPlayer()
	player.mutex.Lock()
	registered := player.registry.Has(id)
	player.mutex.Unlock()

	if !registered {
		if err := LoadAudioFile(id, filename); err != nil {
			log.Printf("Failed to play SFX %s: %v", name, err)
			return
		}
	}

	// Play the SFX (with loop set to false)
	Music(id, false)
}

// UpdateAudio performs lightweight audio housekeeping.
// Playback itself is driven by the mixer feeder goroutine once the runtime
// audio engine has been initialized.
func UpdateAudio() {
	ap := getAudioPlayer()
	ap.update()
}

func initAudio() {
	initializedOnce.Do(func() {
		log.Println("Initializing audio system...")
		// log.Printf("Audio filesystem type: %T, value: %+v", audioFS, audioFS)

		// Check if the audio filesystem is properly initialized
		if !audioFSInit {
			log.Println("ERROR: Cannot initialize audio - audioFS has not been set up")
			return
		}

		// List all files in the audio directory
		// log.Println("Scanning for audio files...")
		files, err := audioFS.ReadDir(".")
		if err != nil {
			log.Printf("ERROR: Failed to list audio files: %v", err)
			return
		}

		// Track loaded files
		loadedCount := 0

		// Only try to load files that actually exist
		for _, file := range files {
			name := file.Name()
			// Check if it's a music file with pattern musicN.raw
			var trackID int
			_, err := fmt.Sscanf(name, "music%d.raw", &trackID)
			if err == nil {
				// Found a valid music file, try to load it
				if err := LoadAudioFile(trackID, name); err != nil {
					log.Printf("Could not load music track %d (%s): %v", trackID, name, err)
				} else {
					// log.Printf("Successfully loaded music track %d: %s", trackID, name)
					loadedCount++
				}
			}
		}

		log.Printf("Audio initialization complete. Successfully loaded %d music tracks.", loadedCount)
	})
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
// output rate. This is the same self-pacing mechanism used by the legacy
// ReadFrom-based feeder.
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

// SetAudioFS sets the audio filesystem instance
// This is called by the generated audio_embed.go file
func SetAudioFS(fs embed.FS) {
	// log.Println("Setting up audio filesystem...")
	audioFS = fs
	audioFSInit = true

	// List all files in the root directory to verify the filesystem
	// log.Println("Audio filesystem contents:")
	files, err := fs.ReadDir(".")
	if err != nil {
		log.Printf("ERROR: Failed to read directory: %v", err)
	} else {
		for _, file := range files {
			_, err := file.Info()
			if err != nil {
				log.Printf("  - %s (error getting info: %v)", file.Name(), err)
			}
		}
	}

	// Initialize and load audio files now that we have the filesystem
	initAudio()
}
