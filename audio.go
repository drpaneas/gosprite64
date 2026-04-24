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
	"github.com/drpaneas/gosprite64/internal/audioengine"
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

// Music plays or stops music
// If n is -1, stops all music
func Music(n int, loop bool) {
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
func PlaySFX(name string) {
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
