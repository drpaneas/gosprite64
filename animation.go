package gosprite64

import "github.com/drpaneas/gosprite64/internal/tile2d/format"

type AnimationSet struct {
	name   string
	parsed format.ParsedAnim
}

type AnimationClip struct {
	Name   string
	FPS    uint16
	Frames []uint16
}

func (a *AnimationSet) Name() string {
	if a == nil {
		return ""
	}
	return a.name
}

func (a *AnimationSet) Clips() []AnimationClip {
	if a == nil {
		return nil
	}

	clips := make([]AnimationClip, 0, len(a.parsed.Clips))
	for _, clip := range a.parsed.Clips {
		clips = append(clips, AnimationClip{
			Name:   clip.Name,
			FPS:    clip.FPS,
			Frames: append([]uint16(nil), clip.Frames...),
		})
	}
	return clips
}

func (a *AnimationSet) Clip(name string) (AnimationClip, bool) {
	if a == nil {
		return AnimationClip{}, false
	}
	for _, clip := range a.parsed.Clips {
		if clip.Name != name {
			continue
		}
		return AnimationClip{
			Name:   clip.Name,
			FPS:    clip.FPS,
			Frames: append([]uint16(nil), clip.Frames...),
		}, true
	}
	return AnimationClip{}, false
}
