package audioengine

import (
	"fmt"
	"slices"
)

const (
	RuntimeSampleRate     = 48000
	RuntimeChannels       = 2
	RuntimeBytesPerSample = 2
	RuntimeByteOrder      = "big-endian"
	RuntimeEncoding       = "signed 16-bit PCM"
)

type CueRef struct {
	TrackID  int
	Filename string
}

type ActiveCue struct{ Loop bool }

type MixerSourcePlan struct {
	TrackID int
	Loop    bool
}

func ResolveSFX(name string) (CueRef, bool) {
	if name == "" {
		return CueRef{}, false
	}

	return CueRef{
		TrackID:  -int(hashString(name)),
		Filename: "sfx_" + name + ".raw",
	}, true
}

func MusicFilename(id int) string {
	return fmt.Sprintf("music%d.raw", id)
}

func BuildMixerPlan(cues map[int]ActiveCue) []MixerSourcePlan {
	ids := make([]int, 0, len(cues))
	for id := range cues {
		ids = append(ids, id)
	}
	slices.Sort(ids)

	plan := make([]MixerSourcePlan, 0, len(ids))
	for _, id := range ids {
		cue := cues[id]
		plan = append(plan, MixerSourcePlan{
			TrackID: id,
			Loop:    cue.Loop,
		})
	}

	return plan
}

func hashString(s string) uint32 {
	var h uint32 = 5381
	for i := 0; i < len(s); i++ {
		h = ((h << 5) + h) + uint32(s[i])
	}
	return h
}
