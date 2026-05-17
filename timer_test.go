package gosprite64

import "testing"

func TestTimerBasic(t *testing.T) {
	tm := NewTimer(5)
	if tm.Done() {
		t.Fatal("new timer should not be done")
	}
	for i := 0; i < 5; i++ {
		tm.Tick()
	}
	if !tm.Done() {
		t.Fatal("timer should be done after 5 ticks")
	}
}

func TestTimerProgress(t *testing.T) {
	tm := NewTimer(10)
	if tm.Progress() != 0 {
		t.Fatalf("expected 0, got %f", tm.Progress())
	}
	for i := 0; i < 5; i++ {
		tm.Tick()
	}
	p := tm.Progress()
	if p < 0.49 || p > 0.51 {
		t.Fatalf("expected ~0.5, got %f", p)
	}
	for i := 0; i < 5; i++ {
		tm.Tick()
	}
	if tm.Progress() != 1 {
		t.Fatalf("expected 1, got %f", tm.Progress())
	}
}

func TestTimerReset(t *testing.T) {
	tm := NewTimer(3)
	tm.Tick()
	tm.Tick()
	tm.Tick()
	if !tm.Done() {
		t.Fatal("should be done")
	}
	tm.Reset()
	if tm.Done() {
		t.Fatal("should not be done after reset")
	}
	if tm.Elapsed() != 0 {
		t.Fatalf("elapsed should be 0 after reset, got %d", tm.Elapsed())
	}
}

func TestTimerResetWithNewDuration(t *testing.T) {
	tm := NewTimer(3)
	tm.Tick()
	tm.Tick()
	tm.Tick()
	tm.ResetWith(10)
	if tm.Done() {
		t.Fatal("should not be done after ResetWith(10)")
	}
	if tm.Duration() != 10 {
		t.Fatalf("duration should be 10, got %d", tm.Duration())
	}
}

func TestTimerElapsed(t *testing.T) {
	tm := NewTimer(10)
	tm.Tick()
	tm.Tick()
	tm.Tick()
	if tm.Elapsed() != 3 {
		t.Fatalf("expected 3, got %d", tm.Elapsed())
	}
}

func TestTimerRemaining(t *testing.T) {
	tm := NewTimer(10)
	tm.Tick()
	tm.Tick()
	if tm.Remaining() != 8 {
		t.Fatalf("expected 8, got %d", tm.Remaining())
	}
}

func TestTimerTickPastDone(t *testing.T) {
	tm := NewTimer(2)
	tm.Tick()
	tm.Tick()
	tm.Tick()
	tm.Tick()
	if tm.Elapsed() != 2 {
		t.Fatalf("elapsed should cap at duration, got %d", tm.Elapsed())
	}
}

func TestTimerZeroDuration(t *testing.T) {
	tm := NewTimer(0)
	if !tm.Done() {
		t.Fatal("zero-duration timer should be immediately done")
	}
}

func TestRepeatingTimerBasic(t *testing.T) {
	rt := NewRepeatingTimer(3)
	triggers := 0
	for i := 0; i < 10; i++ {
		if rt.Tick() {
			triggers++
		}
	}
	if triggers != 3 {
		t.Fatalf("expected 3 triggers in 10 frames at interval 3, got %d", triggers)
	}
}

func TestRepeatingTimerCount(t *testing.T) {
	rt := NewRepeatingTimer(5)
	for i := 0; i < 15; i++ {
		rt.Tick()
	}
	if rt.Count() != 3 {
		t.Fatalf("expected 3 triggers, got %d", rt.Count())
	}
}

func TestRepeatingTimerReset(t *testing.T) {
	rt := NewRepeatingTimer(3)
	for i := 0; i < 6; i++ {
		rt.Tick()
	}
	rt.Reset()
	if rt.Count() != 0 {
		t.Fatalf("count should be 0 after reset, got %d", rt.Count())
	}
}
