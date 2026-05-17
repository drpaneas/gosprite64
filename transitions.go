package gosprite64

import "image/color"

type TransitionStyle int

const (
	FadeToBlack TransitionStyle = iota
	FadeFromBlack
)

type Transition struct {
	Style    TransitionStyle
	Duration int
	frame    int
	active   bool
}

func StartTransition(style TransitionStyle, durationFrames int) *Transition {
	return &Transition{Style: style, Duration: durationFrames, active: true}
}

func (tr *Transition) Advance() {
	if tr == nil || !tr.active {
		return
	}
	if tr.frame < tr.Duration {
		tr.frame++
	}
}

func (tr *Transition) Done() bool {
	if tr == nil {
		return true
	}
	return tr.frame >= tr.Duration
}

func (tr *Transition) Active() bool {
	return tr != nil && tr.active && !tr.Done()
}

func (tr *Transition) Stop() {
	if tr != nil {
		tr.active = false
	}
}

func (tr *Transition) alpha() uint8 {
	if tr.Duration <= 0 {
		return 255
	}
	t := float32(tr.frame) / float32(tr.Duration)
	if t > 1 {
		t = 1
	}
	switch tr.Style {
	case FadeToBlack:
		return uint8(t * 255)
	case FadeFromBlack:
		return uint8((1 - t) * 255)
	}
	return 0
}

func (tr *Transition) Draw() {
	if tr == nil || !tr.active {
		return
	}
	a := tr.alpha()
	if a == 0 {
		return
	}
	drawTransitionOverlay(color.RGBA{R: 0, G: 0, B: 0, A: a})
}
