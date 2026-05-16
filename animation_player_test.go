package gosprite64

import "testing"

func TestAnimationPlayerFPS12(t *testing.T) {
	clip := AnimationClip{Name: "walk", FPS: 12, Frames: []uint16{0, 1, 2, 3}}
	p := NewAnimationPlayer()
	p.Play(clip)
	if p.Frame() != 0 {
		t.Fatalf("Frame() = %d, want 0", p.Frame())
	}
	for i := 0; i < 5; i++ {
		p.Advance(1)
	}
	if p.Frame() != 1 {
		t.Fatalf("after 5 ticks at FPS 12, Frame() = %d, want 1", p.Frame())
	}
}

func TestAnimationPlayerFPS24(t *testing.T) {
	clip := AnimationClip{Name: "run", FPS: 24, Frames: []uint16{0, 1, 2, 3, 4, 5}}
	p := NewAnimationPlayer()
	p.Play(clip)
	frames := make([]int, 0)
	for tick := 0; tick < 15; tick++ {
		p.Advance(1)
		frames = append(frames, p.Frame())
	}
	if frames[1] != 0 {
		t.Fatalf("FPS 24: after 2 ticks expected frame 0 (accumulator 48 < 60), got %d", frames[1])
	}
	if frames[2] != 1 {
		t.Fatalf("FPS 24: after 3 ticks expected frame 1 (accumulator 72 >= 60), got %d", frames[2])
	}
	if frames[4] != 2 {
		t.Fatalf("FPS 24: after 5 ticks expected frame 2, got %d (sequence: %v)", frames[4], frames[:6])
	}
}

func TestAnimationPlayerFPS7(t *testing.T) {
	clip := AnimationClip{Name: "slow", FPS: 7, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.SetLoop(true)
	p.Play(clip)
	for tick := 0; tick < 60; tick++ {
		p.Advance(1)
	}
	if !p.Playing() {
		t.Fatal("looping player should still be playing after 60 ticks")
	}
}

func TestAnimationPlayerFPSAbove60(t *testing.T) {
	clip := AnimationClip{Name: "flash", FPS: 120, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.Play(clip)
	p.Advance(1)
	if p.Frame() != 2 {
		t.Fatalf("FPS 120 non-looping: after 1 tick expected frame 2, got %d", p.Frame())
	}
	if p.Done() {
		t.Fatal("should not be done when landing exactly on the last frame")
	}
	p.Advance(1)
	if p.Frame() != 2 {
		t.Fatalf("after advancing past the end, Frame() = %d, want 2", p.Frame())
	}
	if !p.Done() {
		t.Fatal("should be done after advancing past the last frame")
	}
}

func TestAnimationPlayerLoops(t *testing.T) {
	clip := AnimationClip{Name: "walk", FPS: 60, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.SetLoop(true)
	p.Play(clip)
	for i := 0; i < 4; i++ {
		p.Advance(1)
	}
	if p.Frame() != 1 {
		t.Fatalf("after looping, Frame() = %d, want 1", p.Frame())
	}
	if p.Done() {
		t.Fatal("looping player should not be done")
	}
}

func TestAnimationPlayerStopsAtEnd(t *testing.T) {
	clip := AnimationClip{Name: "die", FPS: 60, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.Play(clip)
	for i := 0; i < 10; i++ {
		p.Advance(1)
	}
	if p.Frame() != 2 {
		t.Fatalf("Frame() = %d, want 2", p.Frame())
	}
	if !p.Done() {
		t.Fatal("expected done")
	}
}

func TestAnimationPlayerPauseResume(t *testing.T) {
	clip := AnimationClip{Name: "walk", FPS: 60, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.Play(clip)
	p.Advance(1)
	p.Pause()
	p.Advance(5)
	if p.Frame() != 1 {
		t.Fatalf("after pause, Frame() = %d, want 1", p.Frame())
	}
	p.Resume()
	p.Advance(1)
	if p.Frame() != 2 {
		t.Fatalf("after resume, Frame() = %d, want 2", p.Frame())
	}
}

func TestAnimationPlayerStop(t *testing.T) {
	clip := AnimationClip{Name: "walk", FPS: 60, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.Play(clip)
	p.Advance(2)
	p.Stop()
	if p.Frame() != 0 {
		t.Fatalf("after stop, Frame() = %d, want 0", p.Frame())
	}
	if !p.Done() {
		t.Fatal("stopped player should be done")
	}
}

func TestAnimationPlayerRestart(t *testing.T) {
	clip := AnimationClip{Name: "walk", FPS: 60, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.Play(clip)
	p.Advance(2)
	p.Restart()
	if p.Frame() != 0 {
		t.Fatalf("after restart, Frame() = %d, want 0", p.Frame())
	}
	if !p.Playing() {
		t.Fatal("restarted player should be playing")
	}
}

func TestAnimationPlayerAdvanceZero(t *testing.T) {
	clip := AnimationClip{Name: "walk", FPS: 60, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.Play(clip)
	p.Advance(0)
	if p.Frame() != 0 {
		t.Fatalf("Frame() = %d, want 0", p.Frame())
	}
}

func TestAnimationPlayerNilSafe(t *testing.T) {
	var p *AnimationPlayer
	p.Advance(1)
	p.Pause()
	p.Resume()
	p.Stop()
	p.Restart()
	if p.Frame() != 0 {
		t.Fatal("nil Frame() should return 0")
	}
	if p.Playing() {
		t.Fatal("nil Playing() should return false")
	}
	if !p.Done() {
		t.Fatal("nil Done() should return true")
	}
}

func TestAnimationPlayerLargeAdvance(t *testing.T) {
	clip := AnimationClip{Name: "walk", FPS: 10, Frames: []uint16{0, 1, 2, 3}}
	p := NewAnimationPlayer()
	p.SetLoop(true)
	p.Play(clip)
	p.Advance(600)
	if !p.Playing() {
		t.Fatal("should still be playing after large advance with loop")
	}
}

func TestAnimationPlayerPlayRejectsEmptyClip(t *testing.T) {
	clip := AnimationClip{Name: "empty", FPS: 12, Frames: []uint16{}}
	p := NewAnimationPlayer()
	p.Play(clip)
	if p.Playing() {
		t.Fatal("playing an empty clip should not enter playing state")
	}
}

func TestAnimationPlayerPlayEmptyClipStopsExistingPlayback(t *testing.T) {
	valid := AnimationClip{Name: "walk", FPS: 12, Frames: []uint16{0, 1, 2}}
	p := NewAnimationPlayer()
	p.Play(valid)
	p.Advance(1)
	if !p.Playing() {
		t.Fatal("should be playing after valid Play")
	}
	empty := AnimationClip{Name: "empty", FPS: 12, Frames: []uint16{}}
	p.Play(empty)
	if p.Playing() {
		t.Fatal("Play(emptyClip) on active player should stop playback")
	}
	if !p.Done() {
		t.Fatal("Play(emptyClip) on active player should set done")
	}
	if p.Frame() != 0 {
		t.Fatalf("Play(emptyClip) should reset frame to 0, got %d", p.Frame())
	}
}

func TestAnimationPlayerRestartWithNoClipIsNoop(t *testing.T) {
	p := NewAnimationPlayer()
	p.Restart()
	if p.Playing() {
		t.Fatal("restart with no clip should not enter playing state")
	}
}
