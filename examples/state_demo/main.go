package main

import (
	"fmt"

	"github.com/drpaneas/gosprite64"
	"github.com/drpaneas/gosprite64/math2d"
)

// Three screens: Title -> Gameplay -> GameOver, plus a Pause overlay.
// A button advances through screens, Start toggles pause during gameplay.

// --- Title Screen ---

type TitleState struct {
	sm      *gosprite64.StateMachine
	blink   int
	visible bool
}

func (s *TitleState) Enter() {
	s.blink = 0
	s.visible = true
}

func (s *TitleState) Update() {
	s.blink++
	if s.blink%30 == 0 {
		s.visible = !s.visible
	}
	if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) {
		s.sm.Switch(&GameplayState{sm: s.sm})
	}
}

func (s *TitleState) Draw() {
	gosprite64.ClearScreenWith(gosprite64.DarkBlue)
	gosprite64.DrawText("STATE MACHINE DEMO", 72, 60, gosprite64.White)
	gosprite64.DrawText("A simple game with 3 screens", 44, 80, gosprite64.LightGray)
	gosprite64.DrawText("and a pause overlay", 68, 92, gosprite64.LightGray)
	if s.visible {
		gosprite64.DrawText("PRESS A TO START", 80, 140, gosprite64.Yellow)
	}
}

func (s *TitleState) Exit() {}

// --- Gameplay Screen ---

type GameplayState struct {
	sm      *gosprite64.StateMachine
	playerX float32
	playerY float32
	score   int
	targetX float32
	targetY float32
	rng     *math2d.Rand
}

func (s *GameplayState) Enter() {
	s.playerX = 144
	s.playerY = 108
	s.rng = math2d.NewRand(99)
	s.spawnTarget()
}

func (s *GameplayState) spawnTarget() {
	s.targetX = s.rng.RangeFloat32(20, 260)
	s.targetY = s.rng.RangeFloat32(30, 190)
}

func (s *GameplayState) Update() {
	if gosprite64.IsButtonJustPressed(gosprite64.ButtonStart) {
		s.sm.Push(&PauseState{sm: s.sm})
		return
	}

	speed := float32(2)
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadUp) {
		s.playerY -= speed
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadDown) {
		s.playerY += speed
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadLeft) {
		s.playerX -= speed
	}
	if gosprite64.IsButtonDown(gosprite64.ButtonDPadRight) {
		s.playerX += speed
	}

	s.playerX = math2d.Clamp(s.playerX, 8, 280)
	s.playerY = math2d.Clamp(s.playerY, 8, 208)

	player := math2d.Rect{X: s.playerX - 4, Y: s.playerY - 4, W: 8, H: 8}
	target := math2d.Rect{X: s.targetX - 3, Y: s.targetY - 3, W: 6, H: 6}
	if player.Overlaps(target) {
		s.score++
		if s.score >= 5 {
			s.sm.Switch(&GameOverState{sm: s.sm, finalScore: s.score})
			return
		}
		s.spawnTarget()
	}
}

func (s *GameplayState) Draw() {
	gosprite64.ClearScreen()

	gosprite64.DrawRect(4, 14, 284, 212, gosprite64.DarkGray)

	gosprite64.FillRect(
		int(s.targetX)-3, int(s.targetY)-3,
		int(s.targetX)+3, int(s.targetY)+3,
		gosprite64.Yellow,
	)

	gosprite64.FillRect(
		int(s.playerX)-4, int(s.playerY)-4,
		int(s.playerX)+4, int(s.playerY)+4,
		gosprite64.Green,
	)

	dist := math2d.Vec2{X: s.playerX, Y: s.playerY}.Distance(
		math2d.Vec2{X: s.targetX, Y: s.targetY},
	)

	gosprite64.DrawText(fmt.Sprintf("score:%d/5  dist:%.0f", s.score, dist), 8, 2, gosprite64.White)
	gosprite64.DrawText("START=pause", 200, 2, gosprite64.DarkGray)
}

func (s *GameplayState) Exit() {}

// --- Pause Overlay ---

type PauseState struct {
	sm    *gosprite64.StateMachine
	frame int
}

func (s *PauseState) Enter() {
	s.frame = 0
}

func (s *PauseState) Update() {
	s.frame++
	if gosprite64.IsButtonJustPressed(gosprite64.ButtonStart) {
		s.sm.Pop()
	}
}

func (s *PauseState) Draw() {
	gosprite64.FillRect(60, 80, 228, 136, gosprite64.DarkPurple)
	gosprite64.DrawRect(60, 80, 228, 136, gosprite64.White)
	gosprite64.DrawText("PAUSED", 120, 96, gosprite64.White)
	if s.frame%40 < 30 {
		gosprite64.DrawText("PRESS START TO RESUME", 72, 116, gosprite64.LightGray)
	}
}

func (s *PauseState) Exit() {}

// --- Game Over Screen ---

type GameOverState struct {
	sm         *gosprite64.StateMachine
	finalScore int
}

func (s *GameOverState) Enter() {}

func (s *GameOverState) Update() {
	if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) {
		s.sm.Switch(&TitleState{sm: s.sm})
	}
}

func (s *GameOverState) Draw() {
	gosprite64.ClearScreenWith(gosprite64.DarkPurple)
	gosprite64.DrawText("GAME OVER", 108, 70, gosprite64.White)
	gosprite64.DrawText(fmt.Sprintf("FINAL SCORE: %d", s.finalScore), 92, 100, gosprite64.Yellow)
	gosprite64.DrawText("PRESS A FOR TITLE", 80, 140, gosprite64.LightGray)
}

func (s *GameOverState) Exit() {}

// --- Main Game Wrapper ---

type Game struct {
	sm *gosprite64.StateMachine
}

func (g *Game) Init() {
	title := &TitleState{}
	g.sm = gosprite64.NewStateMachine(title)
	title.sm = g.sm
	g.sm.Init()
}

func (g *Game) Update() { g.sm.Update() }
func (g *Game) Draw()   { g.sm.Draw() }

func main() {
	gosprite64.Run(&Game{})
}
