package gosprite64

import "testing"

func TestInputRecorderCapture(t *testing.T) {
	rec := NewInputRecorder(1)
	rec.CaptureFrame(0, FrameInput{Buttons: ButtonA | ButtonB, StickX: 50, StickY: -30})
	rec.CaptureFrame(0, FrameInput{Buttons: ButtonA, StickX: 50, StickY: -30})
	rec.CaptureFrame(0, FrameInput{Buttons: 0, StickX: 0, StickY: 0})

	data := rec.Finish()
	if data == nil {
		t.Fatal("Finish should return non-nil data")
	}
	if data.FrameCount != 3 {
		t.Fatalf("expected 3 frames, got %d", data.FrameCount)
	}
	if data.PlayerCount != 1 {
		t.Fatalf("expected 1 player, got %d", data.PlayerCount)
	}
}

func TestInputPlayerPlayback(t *testing.T) {
	rec := NewInputRecorder(1)
	rec.CaptureFrame(0, FrameInput{Buttons: ButtonA, StickX: 10, StickY: 20})
	rec.CaptureFrame(0, FrameInput{Buttons: ButtonB, StickX: -10, StickY: -20})
	rec.CaptureFrame(0, FrameInput{Buttons: 0, StickX: 0, StickY: 0})
	data := rec.Finish()

	player := NewInputPlayer(data)

	input, ok := player.NextFrame(0)
	if !ok {
		t.Fatal("should have frame 0")
	}
	if input.Buttons != ButtonA {
		t.Fatalf("frame 0: expected ButtonA, got %d", input.Buttons)
	}
	if input.StickX != 10 || input.StickY != 20 {
		t.Fatalf("frame 0: expected stick (10,20), got (%d,%d)", input.StickX, input.StickY)
	}

	input, ok = player.NextFrame(0)
	if !ok {
		t.Fatal("should have frame 1")
	}
	if input.Buttons != ButtonB {
		t.Fatalf("frame 1: expected ButtonB, got %d", input.Buttons)
	}

	input, ok = player.NextFrame(0)
	if !ok {
		t.Fatal("should have frame 2")
	}
	if input.Buttons != 0 {
		t.Fatalf("frame 2: expected 0, got %d", input.Buttons)
	}

	_, ok = player.NextFrame(0)
	if ok {
		t.Fatal("should return false after all frames consumed")
	}
}

func TestInputPlayerDone(t *testing.T) {
	rec := NewInputRecorder(1)
	rec.CaptureFrame(0, FrameInput{Buttons: ButtonA})
	data := rec.Finish()

	player := NewInputPlayer(data)
	if player.Done() {
		t.Fatal("should not be done before consuming frames")
	}
	player.NextFrame(0)
	if !player.Done() {
		t.Fatal("should be done after consuming all frames")
	}
}

func TestInputRecorderMultiPlayer(t *testing.T) {
	rec := NewInputRecorder(2)
	rec.CaptureFrame(0, FrameInput{Buttons: ButtonA})
	rec.CaptureFrame(1, FrameInput{Buttons: ButtonB})
	rec.CaptureFrame(0, FrameInput{Buttons: ButtonStart})
	rec.CaptureFrame(1, FrameInput{Buttons: ButtonZ})
	data := rec.Finish()

	if data.PlayerCount != 2 {
		t.Fatalf("expected 2 players, got %d", data.PlayerCount)
	}

	player := NewInputPlayer(data)

	p0f0, _ := player.NextFrame(0)
	p1f0, _ := player.NextFrame(1)
	if p0f0.Buttons != ButtonA {
		t.Fatalf("p0 frame 0: expected ButtonA, got %d", p0f0.Buttons)
	}
	if p1f0.Buttons != ButtonB {
		t.Fatalf("p1 frame 0: expected ButtonB, got %d", p1f0.Buttons)
	}

	p0f1, _ := player.NextFrame(0)
	p1f1, _ := player.NextFrame(1)
	if p0f1.Buttons != ButtonStart {
		t.Fatalf("p0 frame 1: expected ButtonStart, got %d", p0f1.Buttons)
	}
	if p1f1.Buttons != ButtonZ {
		t.Fatalf("p1 frame 1: expected ButtonZ, got %d", p1f1.Buttons)
	}
}

func TestInputRecorderEmpty(t *testing.T) {
	rec := NewInputRecorder(1)
	data := rec.Finish()
	if data.FrameCount != 0 {
		t.Fatal("empty recorder should have 0 frames")
	}
	player := NewInputPlayer(data)
	if !player.Done() {
		t.Fatal("player of empty recording should be done immediately")
	}
}

func TestInputPlayerReset(t *testing.T) {
	rec := NewInputRecorder(1)
	rec.CaptureFrame(0, FrameInput{Buttons: ButtonA})
	rec.CaptureFrame(0, FrameInput{Buttons: ButtonB})
	data := rec.Finish()

	player := NewInputPlayer(data)
	player.NextFrame(0)
	player.NextFrame(0)
	if !player.Done() {
		t.Fatal("should be done")
	}

	player.Reset()
	if player.Done() {
		t.Fatal("should not be done after reset")
	}

	input, ok := player.NextFrame(0)
	if !ok || input.Buttons != ButtonA {
		t.Fatal("after reset, should replay from beginning")
	}
}

func TestReplayDataDeterministic(t *testing.T) {
	rec1 := NewInputRecorder(1)
	rec2 := NewInputRecorder(1)

	inputs := []FrameInput{
		{Buttons: ButtonA, StickX: 10, StickY: 20},
		{Buttons: ButtonA, StickX: 10, StickY: 20},
		{Buttons: ButtonB, StickX: -5, StickY: 0},
		{Buttons: 0, StickX: 0, StickY: 0},
	}
	for _, inp := range inputs {
		rec1.CaptureFrame(0, inp)
		rec2.CaptureFrame(0, inp)
	}

	data1 := rec1.Finish()
	data2 := rec2.Finish()

	if data1.FrameCount != data2.FrameCount {
		t.Fatal("same inputs should produce same frame count")
	}

	player1 := NewInputPlayer(data1)
	player2 := NewInputPlayer(data2)

	for !player1.Done() {
		f1, _ := player1.NextFrame(0)
		f2, _ := player2.NextFrame(0)
		if f1.Buttons != f2.Buttons || f1.StickX != f2.StickX || f1.StickY != f2.StickY {
			t.Fatal("replayed inputs should be identical")
		}
	}
}
