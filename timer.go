package gosprite64

// Timer counts down a fixed number of frames. Use it for delays,
// cooldowns, animation timing, or any "wait N frames" pattern.
type Timer struct {
	duration int
	elapsed  int
}

// NewTimer creates a timer that runs for the given number of frames.
func NewTimer(durationFrames int) *Timer {
	return &Timer{duration: durationFrames}
}

// Tick advances the timer by one frame. Returns true on the frame it finishes.
func (t *Timer) Tick() bool {
	if t == nil || t.elapsed >= t.duration {
		return false
	}
	t.elapsed++
	return t.elapsed == t.duration
}

// Done reports whether the timer has finished.
func (t *Timer) Done() bool {
	if t == nil {
		return true
	}
	return t.elapsed >= t.duration
}

// Progress returns a 0..1 ratio of elapsed/duration.
func (t *Timer) Progress() float32 {
	if t == nil || t.duration <= 0 {
		return 1
	}
	p := float32(t.elapsed) / float32(t.duration)
	if p > 1 {
		p = 1
	}
	return p
}

// Elapsed returns the number of frames that have passed.
func (t *Timer) Elapsed() int {
	if t == nil {
		return 0
	}
	return t.elapsed
}

// Remaining returns the number of frames left.
func (t *Timer) Remaining() int {
	if t == nil {
		return 0
	}
	r := t.duration - t.elapsed
	if r < 0 {
		return 0
	}
	return r
}

// Duration returns the total frame count.
func (t *Timer) Duration() int {
	if t == nil {
		return 0
	}
	return t.duration
}

// Reset restarts the timer with the same duration.
func (t *Timer) Reset() {
	if t == nil {
		return
	}
	t.elapsed = 0
}

// ResetWith restarts the timer with a new duration.
func (t *Timer) ResetWith(durationFrames int) {
	if t == nil {
		return
	}
	t.duration = durationFrames
	t.elapsed = 0
}

// RepeatingTimer fires at a fixed interval and counts how many times
// it has triggered. Use it for blinking cursors, spawn waves, etc.
type RepeatingTimer struct {
	interval int
	elapsed  int
	count    int
}

// NewRepeatingTimer creates a timer that triggers every intervalFrames.
func NewRepeatingTimer(intervalFrames int) *RepeatingTimer {
	if intervalFrames <= 0 {
		intervalFrames = 1
	}
	return &RepeatingTimer{interval: intervalFrames}
}

// Tick advances by one frame. Returns true on trigger frames.
func (rt *RepeatingTimer) Tick() bool {
	if rt == nil {
		return false
	}
	rt.elapsed++
	if rt.elapsed >= rt.interval {
		rt.elapsed = 0
		rt.count++
		return true
	}
	return false
}

// Count returns how many times the timer has triggered.
func (rt *RepeatingTimer) Count() int {
	if rt == nil {
		return 0
	}
	return rt.count
}

// Reset clears the elapsed time and trigger count.
func (rt *RepeatingTimer) Reset() {
	if rt == nil {
		return
	}
	rt.elapsed = 0
	rt.count = 0
}
