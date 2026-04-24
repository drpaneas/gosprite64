package audioengine

import (
	"fmt"
	"slices"
)

const MaxMixerChannels = 32

type MixerRuntimeHooks struct {
	ResetQueue   func()
	InitMixer    func()
	StartDAC     func(rate int)
	SetMixerRate func(rate uint)
	StartFeeder  func()
}

func StartMixerRuntime(hooks MixerRuntimeHooks) {
	hooks.ResetQueue()
	hooks.InitMixer()
	hooks.StartDAC(RuntimeSampleRate)
	hooks.SetMixerRate(uint(RuntimeSampleRate))
	hooks.StartFeeder()
}

type Registry struct {
	files map[int]string
	pcm   map[int][]byte
}

func NewRegistry() *Registry {
	return &Registry{
		files: make(map[int]string),
		pcm:   make(map[int][]byte),
	}
}

func (r *Registry) RegisterFile(id int, filename string) {
	r.files[id] = filename
}

func (r *Registry) Has(id int) bool {
	if _, ok := r.pcm[id]; ok {
		return true
	}
	_, ok := r.files[id]
	return ok
}

func (r *Registry) StorePCM(id int, data []byte) {
	r.pcm[id] = append([]byte(nil), data...)
}

func (r *Registry) Load(id int, loader func(filename string) ([]byte, error)) ([]byte, error) {
	if data, ok := r.pcm[id]; ok {
		return data, nil
	}

	filename, ok := r.files[id]
	if !ok {
		return nil, fmt.Errorf("no audio file registered for track ID %d", id)
	}

	data, err := loader(filename)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("audio file %s is empty", filename)
	}

	r.StorePCM(id, data)
	return r.pcm[id], nil
}

type CuePlayback struct {
	Channel int
	TrackID int
	Loop    bool
}

type PlaybackState struct {
	channelByTrack map[int]int
	playbackByChan map[int]CuePlayback
}

func NewPlaybackState() *PlaybackState {
	return &PlaybackState{
		channelByTrack: make(map[int]int),
		playbackByChan: make(map[int]CuePlayback),
	}
}

func (p *PlaybackState) Activate(trackID int, loop bool) (CuePlayback, error) {
	if channel, ok := p.channelByTrack[trackID]; ok {
		playback := CuePlayback{Channel: channel, TrackID: trackID, Loop: loop}
		p.playbackByChan[channel] = playback
		return playback, nil
	}

	channel := p.firstFreeChannel()
	if channel == -1 {
		return CuePlayback{}, fmt.Errorf("no free mixer channels for track ID %d", trackID)
	}

	playback := CuePlayback{Channel: channel, TrackID: trackID, Loop: loop}
	p.channelByTrack[trackID] = channel
	p.playbackByChan[channel] = playback
	return playback, nil
}

func (p *PlaybackState) Snapshot() []CuePlayback {
	channels := make([]int, 0, len(p.playbackByChan))
	for channel := range p.playbackByChan {
		channels = append(channels, channel)
	}
	slices.Sort(channels)

	snapshot := make([]CuePlayback, 0, len(channels))
	for _, channel := range channels {
		snapshot = append(snapshot, p.playbackByChan[channel])
	}
	return snapshot
}

func (p *PlaybackState) StopAll() []int {
	channels := make([]int, 0, len(p.playbackByChan))
	for channel := range p.playbackByChan {
		channels = append(channels, channel)
	}
	slices.Sort(channels)

	clear(p.channelByTrack)
	clear(p.playbackByChan)
	return channels
}

func (p *PlaybackState) Release(trackID int) (int, bool) {
	channel, ok := p.channelByTrack[trackID]
	if !ok {
		return 0, false
	}

	delete(p.channelByTrack, trackID)
	delete(p.playbackByChan, channel)
	return channel, true
}

func (p *PlaybackState) firstFreeChannel() int {
	for channel := 0; channel < MaxMixerChannels; channel++ {
		if _, ok := p.playbackByChan[channel]; !ok {
			return channel
		}
	}
	return -1
}
