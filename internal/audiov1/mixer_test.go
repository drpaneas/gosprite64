package audiov1

import "testing"

const testMaxFrames = 256

func TestMixSingleVoiceFullGain(t *testing.T) {
	m := NewMixer(48000, testMaxFrames)
	src := []int16{1000, -1000, 500}
	taps := []VoiceTap{{Samples: src, SrcRate: 48000, Gain: GainFull, Active: true}}
	out := make([]int16, len(src)*2)

	m.Mix(out, taps)

	for i := 0; i < len(src); i++ {
		if out[i*2] != src[i] || out[i*2+1] != src[i] {
			t.Fatalf("out frame %d = [%d,%d], want %d", i, out[i*2], out[i*2+1], src[i])
		}
	}
	if taps[0].Consumed != len(src) {
		t.Fatalf("Consumed = %d, want %d", taps[0].Consumed, len(src))
	}
}

func TestMixClampsOverflow(t *testing.T) {
	m := NewMixer(48000, testMaxFrames)
	taps := []VoiceTap{
		{Samples: []int16{30000}, SrcRate: 48000, Gain: GainFull, Active: true},
		{Samples: []int16{30000}, SrcRate: 48000, Gain: GainFull, Active: true},
	}
	out := make([]int16, 2)

	m.Mix(out, taps)

	if out[0] != 32767 {
		t.Fatalf("out[0] = %d, want 32767", out[0])
	}
}

func TestMixHalfGain(t *testing.T) {
	m := NewMixer(48000, testMaxFrames)
	taps := []VoiceTap{{Samples: []int16{20000}, SrcRate: 48000, Gain: GainFull / 2, Active: true}}
	out := make([]int16, 2)

	m.Mix(out, taps)

	if out[0] != 10000 {
		t.Fatalf("out[0] = %d, want 10000", out[0])
	}
}

func TestMixFreshBuffersDoNotSkipWithFractionalPhase(t *testing.T) {
	m := NewMixer(48000, testMaxFrames)
	out1 := make([]int16, 20)
	taps1 := []VoiceTap{{Samples: []int16{1000, 1000, 1000, 1000, 1000, 1000}, SrcRate: 22050, Gain: GainFull, Active: true}}

	m.Mix(out1, taps1)
	if taps1[0].Consumed == 0 || taps1[0].Phase == 0 {
		t.Fatalf("first mix consumed=%d phase=%d, want progress and fractional phase", taps1[0].Consumed, taps1[0].Phase)
	}

	out2 := make([]int16, 20)
	taps2 := []VoiceTap{{Samples: []int16{2000, 2000, 2000, 2000, 2000, 2000}, SrcRate: 22050, Gain: GainFull, Active: true, Phase: taps1[0].Phase}}
	m.Mix(out2, taps2)

	if out2[0] < 1000 || out2[0] > 2000 {
		t.Fatalf("fresh-buffer first sample = %d, want interpolated boundary sample", out2[0])
	}
	if taps2[0].Consumed == 0 {
		t.Fatal("second mix consumed no source frames")
	}
}

func TestMixSteadyStateAllocations(t *testing.T) {
	m := NewMixer(48000, testMaxFrames)
	src := make([]int16, testMaxFrames)
	out := make([]int16, testMaxFrames*2)
	taps := []VoiceTap{{Samples: src, SrcRate: 48000, Gain: GainFull, Active: true}}

	allocs := testing.AllocsPerRun(100, func() {
		m.Mix(out, taps)
	})
	if allocs != 0 {
		t.Fatalf("Mix allocations = %f, want 0", allocs)
	}
}
