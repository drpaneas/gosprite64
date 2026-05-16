package gosprite64

type playerState uint8

const (
	playerStopped playerState = iota
	playerPlaying
	playerPaused
)

const defaultTickRate = 60

type AnimationPlayer struct {
	clip        AnimationClip
	state       playerState
	loop        bool
	frameIdx    int
	accumulator int
}

func NewAnimationPlayer() *AnimationPlayer {
	return &AnimationPlayer{}
}

func (p *AnimationPlayer) Play(clip AnimationClip) {
	if p == nil {
		return
	}
	if len(clip.Frames) == 0 {
		p.state = playerStopped
		p.frameIdx = 0
		p.accumulator = 0
		p.clip = AnimationClip{}
		return
	}
	p.clip = clip
	p.state = playerPlaying
	p.frameIdx = 0
	p.accumulator = 0
}

func (p *AnimationPlayer) Pause() {
	if p == nil || p.state != playerPlaying {
		return
	}
	p.state = playerPaused
}

func (p *AnimationPlayer) Resume() {
	if p == nil || p.state != playerPaused {
		return
	}
	p.state = playerPlaying
}

func (p *AnimationPlayer) Stop() {
	if p == nil {
		return
	}
	p.state = playerStopped
	p.frameIdx = 0
	p.accumulator = 0
}

func (p *AnimationPlayer) SetLoop(loop bool) {
	if p == nil {
		return
	}
	p.loop = loop
}

func (p *AnimationPlayer) Restart() {
	if p == nil || len(p.clip.Frames) == 0 {
		return
	}
	p.frameIdx = 0
	p.accumulator = 0
	p.state = playerPlaying
}

func (p *AnimationPlayer) Advance(ticks int) {
	if p == nil || p.state != playerPlaying || ticks <= 0 || len(p.clip.Frames) == 0 {
		return
	}

	fps := int(p.clip.FPS)
	if fps <= 0 {
		fps = defaultTickRate
	}

	p.accumulator += ticks * fps
	framesAdvanced := p.accumulator / defaultTickRate
	p.accumulator = p.accumulator % defaultTickRate

	if framesAdvanced == 0 {
		return
	}

	newIdx := p.frameIdx + framesAdvanced
	frameCount := len(p.clip.Frames)

	if p.loop {
		newIdx = newIdx % frameCount
	} else if newIdx >= frameCount {
		newIdx = frameCount - 1
		p.state = playerStopped
	}

	p.frameIdx = newIdx
}

func (p *AnimationPlayer) Frame() int {
	if p == nil || len(p.clip.Frames) == 0 {
		return 0
	}
	if p.frameIdx < 0 || p.frameIdx >= len(p.clip.Frames) {
		return 0
	}
	return int(p.clip.Frames[p.frameIdx])
}

func (p *AnimationPlayer) Playing() bool {
	return p != nil && p.state == playerPlaying
}

func (p *AnimationPlayer) Done() bool {
	return p == nil || p.state == playerStopped
}
