package gosprite64

import (
	"embed"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/clktmr/n64/rcp/audio"
)

// audioFS is the global embedded filesystem for audio files
// This is set by the generated audio_embed.go file in each game
var (
	audioFS         embed.FS
	audioFSInit     bool           // Flag to track if audioFS has been initialized
	initializedOnce sync.Once     // Ensures initialization happens only once
	musicFiles      map[int]string // Map track IDs to filenames
)

const (
	// numMusicTracks defines the number of music tracks to load
	numMusicTracks = 64 // Maximum number of music tracks (music0.raw to music63.raw)
)

var (
	audioPlayerInstance *audioPlayer
	audioOnce           sync.Once
)

type audioPlayer struct {
	musicData    map[int][]byte
	activeTracks map[int]*musicTrack
	mutex        sync.Mutex
}

type musicTrack struct {
	data     []byte
	position int
	loop     bool
}

func getAudioPlayer() *audioPlayer {
	audioOnce.Do(func() {
		audio.SetSampleRate(48000)
		audioPlayerInstance = &audioPlayer{
			musicData:    make(map[int][]byte),
			activeTracks: make(map[int]*musicTrack),
		}
	})
	return audioPlayerInstance
}

func (a *audioPlayer) update() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for id, track := range a.activeTracks {
		if track == nil || len(track.data) == 0 {
			delete(a.activeTracks, id)
			continue
		}

		remaining := len(track.data) - track.position
		if remaining <= 0 {
			if track.loop {
				track.position = 0
				remaining = len(track.data)
			} else {
				delete(a.activeTracks, id)
				continue
			}
		}

		chunkSize := 4096
		if remaining < chunkSize {
			chunkSize = remaining
		}

		chunk := track.data[track.position : track.position+chunkSize]
		audio.Buffer.Write(chunk)
		track.position += chunkSize
	}
}

// Music plays or stops music
// If n is -1, stops all music
func Music(n int, loop bool) {
	ap := getAudioPlayer()
	
	// Check if we already have this track loaded
	_, exists := ap.musicData[n]
	if !exists {
		// Try to load the track
		data, err := loadAudioData(n)
		if err != nil {
			log.Printf("Failed to load audio track %d: %v", n, err)
			return
		}
		// Store the loaded data
		ap.musicData[n] = data
	}

	ap.activeTracks[n] = &musicTrack{
		data: ap.musicData[n],
		loop: loop,
	}
}

// LoadAudio loads raw PCM audio data for a track
// pcmData: Raw PCM audio data (16-bit stereo, 48kHz)
func LoadAudio(id int, pcmData []byte) {
	ap := getAudioPlayer()
	ap.mutex.Lock()
	defer ap.mutex.Unlock()
	ap.musicData[id] = pcmData
}

// LoadAudioFile registers an audio file to be loaded on demand
func LoadAudioFile(id int, filename string) error {
	log.Printf("Registering audio file: %s as track ID: %d", filename, id)
	
	// Check if the audio filesystem is properly initialized
	if !audioFSInit {
		log.Printf("ERROR: audioFS has not been initialized when trying to register %s", filename)
		return fmt.Errorf("audio filesystem not initialized")
	}

	// Just store the filename, we'll load it when needed
	musicFiles[id] = filename
	log.Printf("Registered audio file %s as track ID %d", filename, id)
	return nil
}

// loadAudioData loads the actual audio data for a track
func loadAudioData(id int) ([]byte, error) {
	filename, exists := musicFiles[id]
	if !exists {
		return nil, fmt.Errorf("no audio file registered for track ID %d", id)
	}

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

	log.Printf("Loaded audio file %s as track ID %d (%d bytes)", filename, id, len(data))
	return data, nil
}

// PlaySFX plays a sound effect by name (without the "sfx_" prefix and ".raw" extension)
// Example: PlaySFX("jump") will play "sfx_jump.raw" if it exists
func PlaySFX(name string) {
	if name == "" {
		return
	}
	
	filename := "sfx_" + name + ".raw"
	id := -int(hashString(name)) // Generate the same ID as in initAudio
	
	// Check if the sound effect is already loaded
	player := getAudioPlayer()
	player.mutex.Lock()
	_, exists := player.musicData[id]
	player.mutex.Unlock()
	
	if !exists {
		// Try to load the SFX if it hasn't been loaded yet
		if err := LoadAudioFile(id, filename); err != nil {
			log.Printf("Failed to play SFX %s: %v", name, err)
			return
		}
	}
	
	// Play the SFX (with loop set to false)
	Music(id, false)
}

// UpdateAudio updates the audio system (call in game loop)
func UpdateAudio() {
	getAudioPlayer().update()
}

func initAudio() {
	initializedOnce.Do(func() {
		log.Println("Initializing audio system...")
		log.Printf("Audio filesystem type: %T, value: %+v", audioFS, audioFS)
		
		// Check if the audio filesystem is properly initialized
		if !audioFSInit {
			log.Println("ERROR: Cannot initialize audio - audioFS has not been set up")
			return
		}

		// Initialize music files map
		musicFiles = make(map[int]string)
		
		// List all files in the audio directory
		log.Println("Scanning for audio files...")
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
					log.Printf("Successfully loaded music track %d: %s", trackID, name)
					loadedCount++
				}
			}
		}

		log.Printf("Audio initialization complete. Successfully loaded %d music tracks.", loadedCount)
	})
}

// hashString generates a simple hash from a string
func hashString(s string) uint32 {
	var h uint32 = 5381
	for i := 0; i < len(s); i++ {
		h = ((h << 5) + h) + uint32(s[i])
	}
	return h
}

// SetAudioFS sets the audio filesystem instance
// This is called by the generated audio_embed.go file
func SetAudioFS(fs embed.FS) {
	log.Println("Setting up audio filesystem...")
	audioFS = fs
	audioFSInit = true
	
	// List all files in the root directory to verify the filesystem
	log.Println("Audio filesystem contents:")
	files, err := fs.ReadDir(".")
	if err != nil {
		log.Printf("ERROR: Failed to read directory: %v", err)
	} else {
		for _, file := range files {
			info, err := file.Info()
			if err != nil {
				log.Printf("  - %s (error getting info: %v)", file.Name(), err)
			} else {
				log.Printf("  - %s (dir: %v, size: %d)", file.Name(), file.IsDir(), info.Size())
			}
		}
	}
	
	// Initialize and load audio files now that we have the filesystem
	initAudio()
}
