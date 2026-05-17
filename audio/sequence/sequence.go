package sequence

// Player is a MIDI-like sequence player for N64 music.
// N64 games typically use a proprietary sequence format that drives
// instrument bank playback through the audio RSP microcode.
type Player struct {
	Data       []byte
	playing    bool
	paused     bool
	tempo      uint16 // BPM
	volume     uint8
	channelVol [16]uint8
	position   uint32
	loopStart  uint32
	loopCount  int
}

// NewPlayer creates a sequence player with default settings.
func NewPlayer() *Player {
	p := &Player{
		tempo:  120,
		volume: 127,
	}
	for i := range p.channelVol {
		p.channelVol[i] = 127
	}
	return p
}

// Play starts or restarts playback from the beginning.
func (p *Player) Play() {
	p.playing = true
	p.paused = false
	p.position = 0
}

// Stop halts playback.
func (p *Player) Stop() {
	p.playing = false
	p.paused = false
	p.position = 0
}

// Pause pauses playback without resetting position.
func (p *Player) Pause() {
	if p.playing {
		p.paused = true
	}
}

// Resume continues paused playback.
func (p *Player) Resume() {
	p.paused = false
}

// IsPlaying returns true if the player is actively playing.
func (p *Player) IsPlaying() bool {
	return p.playing && !p.paused
}

// SetTempo sets the playback tempo in BPM.
func (p *Player) SetTempo(bpm uint16) {
	p.tempo = bpm
}

// Tempo returns the current tempo.
func (p *Player) Tempo() uint16 {
	return p.tempo
}

// SetVolume sets the master volume (0-127).
func (p *Player) SetVolume(vol uint8) {
	if vol > 127 {
		vol = 127
	}
	p.volume = vol
}

// Volume returns the master volume.
func (p *Player) Volume() uint8 {
	return p.volume
}

// SetChannelVolume sets volume for a specific MIDI channel (0-15).
func (p *Player) SetChannelVolume(ch, vol uint8) {
	if ch < 16 {
		if vol > 127 {
			vol = 127
		}
		p.channelVol[ch] = vol
	}
}

// ChannelVolume returns the volume for a channel.
func (p *Player) ChannelVolume(ch uint8) uint8 {
	if ch < 16 {
		return p.channelVol[ch]
	}
	return 0
}

// SetLoop configures loop points. count of -1 means infinite loop.
func (p *Player) SetLoop(startPos uint32, count int) {
	p.loopStart = startPos
	p.loopCount = count
}

// Position returns the current playback position.
func (p *Player) Position() uint32 {
	return p.position
}

// NoteEvent represents a note-on or note-off event from the sequence.
type NoteEvent struct {
	Channel  uint8
	Note     uint8
	Velocity uint8
	On       bool
}

// Tick advances playback by one audio frame and returns any note events.
func (p *Player) Tick(sampleRate uint32) []NoteEvent {
	if !p.playing || p.paused || p.Data == nil {
		return nil
	}

	samplesPerBeat := sampleRate * 60 / uint32(p.tempo)
	ticksPerBeat := uint32(48)
	samplesPerTick := samplesPerBeat / ticksPerBeat

	p.position += samplesPerTick

	if p.loopCount != 0 && p.position >= uint32(len(p.Data)) {
		p.position = p.loopStart
		if p.loopCount > 0 {
			p.loopCount--
		}
		if p.loopCount == 0 {
			p.playing = false
			return nil
		}
	}

	return nil
}
