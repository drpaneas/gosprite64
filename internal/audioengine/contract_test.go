package audioengine

import "testing"

func TestRuntimePCMContract(t *testing.T) {
	if RuntimeSampleRate != 48000 {
		t.Fatalf("RuntimeSampleRate = %d, want 48000", RuntimeSampleRate)
	}
	if RuntimeChannels != 2 {
		t.Fatalf("RuntimeChannels = %d, want 2", RuntimeChannels)
	}
	if RuntimeBytesPerSample != 2 {
		t.Fatalf("RuntimeBytesPerSample = %d, want 2", RuntimeBytesPerSample)
	}
	if RuntimeByteOrder != "big-endian" {
		t.Fatalf("RuntimeByteOrder = %q, want %q", RuntimeByteOrder, "big-endian")
	}
	if RuntimeEncoding != "signed 16-bit PCM" {
		t.Fatalf("RuntimeEncoding = %q, want %q", RuntimeEncoding, "signed 16-bit PCM")
	}
}

func TestResolveSFXReturnsStableFilenameAndTrackID(t *testing.T) {
	cue, ok := ResolveSFX("jump")
	if !ok {
		t.Fatal("ResolveSFX returned ok=false for non-empty name")
	}

	if cue.Filename != "sfx_jump.raw" {
		t.Fatalf("ResolveSFX filename = %q, want %q", cue.Filename, "sfx_jump.raw")
	}

	const wantID = -2090414049
	if cue.TrackID != wantID {
		t.Fatalf("ResolveSFX track id = %d, want %d", cue.TrackID, wantID)
	}
}

func TestResolveSFXRejectsEmptyName(t *testing.T) {
	if _, ok := ResolveSFX(""); ok {
		t.Fatal("ResolveSFX returned ok=true for empty name")
	}
}

func TestBuildMixerPlanCreatesSeparateSourcesForEachCue(t *testing.T) {
	plan := BuildMixerPlan(map[int]ActiveCue{
		11: {Loop: true},
		3:  {Loop: false},
	})

	if len(plan) != 2 {
		t.Fatalf("BuildMixerPlan returned %d sources, want 2", len(plan))
	}

	if plan[0].TrackID != 3 || plan[0].Loop {
		t.Fatalf("first source = %+v, want track 3 loop=false", plan[0])
	}
	if plan[1].TrackID != 11 || !plan[1].Loop {
		t.Fatalf("second source = %+v, want track 11 loop=true", plan[1])
	}
}

func TestBuildMixerPlanHandlesEmptyAndSingleCueSets(t *testing.T) {
	if plan := BuildMixerPlan(nil); len(plan) != 0 {
		t.Fatalf("BuildMixerPlan(nil) returned %d sources, want 0", len(plan))
	}

	plan := BuildMixerPlan(map[int]ActiveCue{
		9: {Loop: true},
	})
	if len(plan) != 1 {
		t.Fatalf("BuildMixerPlan(single) returned %d sources, want 1", len(plan))
	}
	if plan[0].TrackID != 9 || !plan[0].Loop {
		t.Fatalf("single source = %+v, want track 9 loop=true", plan[0])
	}
}
