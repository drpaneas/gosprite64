package audioengine

import (
	"bytes"
	"errors"
	"slices"
	"testing"
)

func TestRegistryLoadCachesAudioLoadedFromFile(t *testing.T) {
	registry := NewRegistry()
	registry.RegisterFile(7, MusicFilename(7))

	loadCalls := 0
	want := []byte{0x01, 0x02}
	load := func(filename string) ([]byte, error) {
		loadCalls++
		if filename != "music7.raw" {
			t.Fatalf("loader filename = %q, want %q", filename, "music7.raw")
		}
		return append([]byte(nil), want...), nil
	}

	got, err := registry.Load(7, load)
	if err != nil {
		t.Fatalf("registry.Load returned error: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("registry.Load bytes = %v, want %v", got, want)
	}

	got, err = registry.Load(7, load)
	if err != nil {
		t.Fatalf("second registry.Load returned error: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("second registry.Load bytes = %v, want %v", got, want)
	}
	if loadCalls != 1 {
		t.Fatalf("loader call count = %d, want 1", loadCalls)
	}
}

func TestRegistryLoadUsesStoredPCMWithoutLoader(t *testing.T) {
	registry := NewRegistry()
	registry.StorePCM(3, []byte{0xaa, 0xbb})

	got, err := registry.Load(3, func(string) ([]byte, error) {
		t.Fatal("loader should not be called for stored PCM")
		return nil, nil
	})
	if err != nil {
		t.Fatalf("registry.Load returned error: %v", err)
	}
	if !bytes.Equal(got, []byte{0xaa, 0xbb}) {
		t.Fatalf("registry.Load bytes = %v, want %v", got, []byte{0xaa, 0xbb})
	}
}

func TestRegistryLoadReturnsErrorForUnknownTrack(t *testing.T) {
	registry := NewRegistry()

	_, err := registry.Load(99, func(string) ([]byte, error) {
		return nil, nil
	})
	if err == nil {
		t.Fatal("registry.Load error = nil, want unknown track error")
	}
}

func TestRegistryLoadPropagatesLoaderErrors(t *testing.T) {
	registry := NewRegistry()
	registry.RegisterFile(4, MusicFilename(4))

	wantErr := errors.New("boom")
	_, err := registry.Load(4, func(string) ([]byte, error) {
		return nil, wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("registry.Load error = %v, want %v", err, wantErr)
	}
}

func TestRegistryHasTracksRegisteredOrStoredAudio(t *testing.T) {
	registry := NewRegistry()
	if registry.Has(1) {
		t.Fatal("Registry.Has returned true for missing track")
	}

	registry.RegisterFile(1, MusicFilename(1))
	if !registry.Has(1) {
		t.Fatal("Registry.Has returned false for registered file")
	}

	registry.StorePCM(2, []byte{0x01})
	if !registry.Has(2) {
		t.Fatal("Registry.Has returned false for stored PCM")
	}
}

func TestPlaybackStateAllocatesSeparateChannelsPerCue(t *testing.T) {
	state := NewPlaybackState()

	first, err := state.Activate(10, false)
	if err != nil {
		t.Fatalf("Activate first cue returned error: %v", err)
	}
	second, err := state.Activate(20, true)
	if err != nil {
		t.Fatalf("Activate second cue returned error: %v", err)
	}

	if first.Channel != 0 || first.TrackID != 10 || first.Loop {
		t.Fatalf("first playback = %+v, want channel 0 track 10 loop=false", first)
	}
	if second.Channel != 1 || second.TrackID != 20 || !second.Loop {
		t.Fatalf("second playback = %+v, want channel 1 track 20 loop=true", second)
	}
}

func TestPlaybackStateReusesExistingChannelForSameCue(t *testing.T) {
	state := NewPlaybackState()

	first, err := state.Activate(5, false)
	if err != nil {
		t.Fatalf("Activate first cue returned error: %v", err)
	}
	second, err := state.Activate(5, true)
	if err != nil {
		t.Fatalf("Activate second cue returned error: %v", err)
	}

	if second.Channel != first.Channel {
		t.Fatalf("second activation channel = %d, want %d", second.Channel, first.Channel)
	}

	snapshot := state.Snapshot()
	if len(snapshot) != 1 {
		t.Fatalf("Snapshot returned %d active cues, want 1", len(snapshot))
	}
	if !snapshot[0].Loop {
		t.Fatalf("Snapshot loop = %v, want true", snapshot[0].Loop)
	}
}

func TestPlaybackStateStopAllReturnsOccupiedChannels(t *testing.T) {
	state := NewPlaybackState()
	if _, err := state.Activate(1, false); err != nil {
		t.Fatalf("Activate first cue returned error: %v", err)
	}
	if _, err := state.Activate(2, true); err != nil {
		t.Fatalf("Activate second cue returned error: %v", err)
	}

	channels := state.StopAll()
	if len(channels) != 2 || channels[0] != 0 || channels[1] != 1 {
		t.Fatalf("StopAll channels = %v, want [0 1]", channels)
	}
	if len(state.Snapshot()) != 0 {
		t.Fatal("Snapshot not empty after StopAll")
	}
}

func TestPlaybackStateReleaseFreesChannelForReuse(t *testing.T) {
	state := NewPlaybackState()

	first, err := state.Activate(1, false)
	if err != nil {
		t.Fatalf("Activate first cue returned error: %v", err)
	}
	second, err := state.Activate(2, false)
	if err != nil {
		t.Fatalf("Activate second cue returned error: %v", err)
	}

	channel, ok := state.Release(1)
	if !ok {
		t.Fatal("Release returned ok=false, want true")
	}
	if channel != first.Channel {
		t.Fatalf("Release channel = %d, want %d", channel, first.Channel)
	}

	third, err := state.Activate(3, true)
	if err != nil {
		t.Fatalf("Activate third cue returned error: %v", err)
	}
	if third.Channel != first.Channel {
		t.Fatalf("reused channel = %d, want %d", third.Channel, first.Channel)
	}
	if third.Channel == second.Channel {
		t.Fatalf("third channel = %d, want different from second channel %d", third.Channel, second.Channel)
	}
}

func TestPlaybackStateReturnsErrorWhenChannelsExhausted(t *testing.T) {
	state := NewPlaybackState()

	for i := 0; i < MaxMixerChannels; i++ {
		if _, err := state.Activate(i, false); err != nil {
			t.Fatalf("Activate(%d) returned error: %v", i, err)
		}
	}

	if _, err := state.Activate(999, false); err == nil {
		t.Fatal("Activate returned nil error when all mixer channels were occupied")
	}
}

func TestStartMixerRuntimeCallsHooksInUpstreamOrder(t *testing.T) {
	var calls []string

	StartMixerRuntime(MixerRuntimeHooks{
		ResetQueue: func() {
			calls = append(calls, "reset")
		},
		InitMixer: func() {
			calls = append(calls, "init")
		},
		StartDAC: func(rate int) {
			if rate != RuntimeSampleRate {
				t.Fatalf("StartDAC rate = %d, want %d", rate, RuntimeSampleRate)
			}
			calls = append(calls, "start-dac")
		},
		SetMixerRate: func(rate uint) {
			if rate != uint(RuntimeSampleRate) {
				t.Fatalf("SetMixerRate rate = %d, want %d", rate, RuntimeSampleRate)
			}
			calls = append(calls, "set-rate")
		},
		StartFeeder: func() {
			calls = append(calls, "start-feeder")
		},
	})

	want := []string{"reset", "init", "start-dac", "set-rate", "start-feeder"}
	if !slices.Equal(calls, want) {
		t.Fatalf("StartMixerRuntime calls = %v, want %v", calls, want)
	}
}
