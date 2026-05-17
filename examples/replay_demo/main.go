package main

import (
	"fmt"

	"github.com/drpaneas/gosprite64"
	"github.com/drpaneas/gosprite64/math2d"
)

// This demo records your input while you move a square around,
// then plays it back as a ghost trail. Press A to toggle between
// recording and playback modes.

const (
	screenW = 288
	screenH = 216
)

type Mode int

const (
	ModeRecord Mode = iota
	ModePlayback
)

type Game struct {
	mode     Mode
	recorder *gosprite64.InputRecorder
	player   *gosprite64.InputPlayer
	replay   *gosprite64.ReplayData

	liveX, liveY     float32
	ghostX, ghostY   float32
	trail            []math2d.Vec2
	ghostTrail       []math2d.Vec2
	frameCount       int
	playbackFrame    int
}

func (g *Game) Init() {
	g.liveX = float32(screenW) / 2
	g.liveY = float32(screenH) / 2
	g.mode = ModeRecord
	g.recorder = gosprite64.NewInputRecorder(1)
}

func (g *Game) Update() {
	switch g.mode {
	case ModeRecord:
		g.updateRecord()
	case ModePlayback:
		g.updatePlayback()
	}
}

func (g *Game) updateRecord() {
	speed := float32(2)
	var sx, sy int8

	if gosprite64.IsButtonDown(gosprite64.ButtonDPadUp) {
		sy = -1
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadDown) {
		sy = 1
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft) {
		sx = -1
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadRight) {
		sx = 1
	}

	g.liveX += float32(sx) * speed
	g.liveY += float32(sy) * speed
	g.liveX = math2d.Clamp(g.liveX, 8, screenW-8)
	g.liveY = math2d.Clamp(g.liveY, 18, screenH-8)

	var buttons gosprite64.ButtonMask = 0
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadUp) {
		buttons |= gosprite64.ButtonDPadUp
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadDown) {
		buttons |= gosprite64.ButtonDPadDown
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft) {
		buttons |= gosprite64.ButtonDPadLeft
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadRight) {
		buttons |= gosprite64.ButtonDPadRight
	}

	g.recorder.CaptureFrame(0, gosprite64.FrameInput{
		Buttons: buttons,
		StickX:  sx,
		StickY:  sy,
	})
	g.frameCount++

	if g.frameCount%3 == 0 {
		g.trail = append(g.trail, math2d.Vec2{X: g.liveX, Y: g.liveY})
		if len(g.trail) > 60 {
			g.trail = g.trail[1:]
		}
	}

	if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) && g.frameCount > 30 {
		g.replay = g.recorder.Finish()
		g.player = gosprite64.NewInputPlayer(g.replay)
		g.ghostX = float32(screenW) / 2
		g.ghostY = float32(screenH) / 2
		g.ghostTrail = nil
		g.playbackFrame = 0
		g.mode = ModePlayback
	}
}

func (g *Game) updatePlayback() {
	input, ok := g.player.NextFrame(0)
	if !ok {
		g.player.Reset()
		g.ghostX = float32(screenW) / 2
		g.ghostY = float32(screenH) / 2
		g.ghostTrail = nil
		g.playbackFrame = 0
		return
	}

	speed := float32(2)
	g.ghostX += float32(input.StickX) * speed
	g.ghostY += float32(input.StickY) * speed
	g.ghostX = math2d.Clamp(g.ghostX, 8, screenW-8)
	g.ghostY = math2d.Clamp(g.ghostY, 18, screenH-8)
	g.playbackFrame++

	if g.playbackFrame%3 == 0 {
		g.ghostTrail = append(g.ghostTrail, math2d.Vec2{X: g.ghostX, Y: g.ghostY})
		if len(g.ghostTrail) > 60 {
			g.ghostTrail = g.ghostTrail[1:]
		}
	}

	if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) {
		g.recorder = gosprite64.NewInputRecorder(1)
		g.frameCount = 0
		g.trail = nil
		g.liveX = float32(screenW) / 2
		g.liveY = float32(screenH) / 2
		g.mode = ModeRecord
	}
}

func (g *Game) Draw() {
	gosprite64.ClearScreen()

	gosprite64.DrawRect(4, 14, screenW-4, screenH-4, gosprite64.DarkGray)

	switch g.mode {
	case ModeRecord:
		for i, p := range g.trail {
			alpha := i * 4
			if alpha > 255 {
				alpha = 255
			}
			gosprite64.FillRect(int(p.X)-1, int(p.Y)-1, int(p.X)+1, int(p.Y)+1, gosprite64.DarkGreen)
		}

		gosprite64.FillRect(
			int(g.liveX)-4, int(g.liveY)-4,
			int(g.liveX)+4, int(g.liveY)+4,
			gosprite64.Green,
		)

		gosprite64.DrawText(fmt.Sprintf("RECORDING  frames:%d", g.frameCount), 8, 2, gosprite64.Red)
		if g.frameCount > 30 {
			gosprite64.DrawText("A=playback", 200, 2, gosprite64.Yellow)
		}

	case ModePlayback:
		for _, p := range g.ghostTrail {
			gosprite64.FillRect(int(p.X)-1, int(p.Y)-1, int(p.X)+1, int(p.Y)+1, gosprite64.DarkBlue)
		}

		gosprite64.FillRect(
			int(g.ghostX)-4, int(g.ghostY)-4,
			int(g.ghostX)+4, int(g.ghostY)+4,
			gosprite64.Blue,
		)

		progress := 0
		if g.replay != nil && g.replay.FrameCount > 0 {
			progress = g.playbackFrame * 100 / g.replay.FrameCount
		}
		gosprite64.DrawText(fmt.Sprintf("PLAYBACK  %d%%  frame:%d/%d", progress, g.playbackFrame, g.replay.FrameCount), 8, 2, gosprite64.Blue)
		gosprite64.DrawText("A=record", 216, 2, gosprite64.Yellow)
	}
}

func main() {
	gosprite64.Run(&Game{})
}
