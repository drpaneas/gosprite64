package main

import (
	"fmt"
	"image/color"

	"github.com/drpaneas/gosprite64"
	"github.com/drpaneas/gosprite64/math2d"
)

// Demonstrates Timer, DrawRegion, and Menu working together.
// A title menu lets you pick 1P or 2P mode. In gameplay, each player
// gets their own viewport (DrawRegion). A countdown timer starts each
// round, and a blinking RepeatingTimer spawns targets to collect.

const (
	screenW = 288
	screenH = 216
)

type GameMode int

const (
	ModeMenu GameMode = iota
	ModePlaying
	ModeGameOver
)

type Player struct {
	x, y  float32
	score int
	port  int
}

type Game struct {
	sm   *gosprite64.StateMachine
	mode GameMode
}

func (g *Game) Init() {
	g.sm = gosprite64.NewStateMachine(&MenuState{sm: nil})
	menuState := g.sm.Current().(*MenuState)
	menuState.sm = g.sm
	g.sm.Init()
}

func (g *Game) Update() { g.sm.Update() }
func (g *Game) Draw()   { g.sm.Draw() }

// --- Menu State ---

type MenuState struct {
	sm   *gosprite64.StateMachine
	menu *gosprite64.Menu
}

func (s *MenuState) Enter() {
	s.menu = gosprite64.NewMenu([]gosprite64.MenuItem{
		{Label: "1 Player", OnConfirm: func() {
			s.sm.Switch(&PlayState{sm: s.sm, playerCount: 1})
		}},
		{Label: "2 Players", OnConfirm: func() {
			s.sm.Switch(&PlayState{sm: s.sm, playerCount: 2})
		}},
	})
	s.menu.X = 100
	s.menu.Y = 100
	s.menu.Wrap = true
}

func (s *MenuState) Update() {
	s.menu.HandleInput()
}

func (s *MenuState) Draw() {
	gosprite64.ClearScreenWith(gosprite64.DarkBlue)
	gosprite64.DrawText("SPLIT SCREEN DEMO", 72, 40, gosprite64.White)
	gosprite64.DrawText("Timer + DrawRegion + Menu", 52, 56, gosprite64.LightGray)
	s.menu.Draw()
}

func (s *MenuState) Exit() {}

// --- Play State ---

type PlayState struct {
	sm          *gosprite64.StateMachine
	playerCount int
	players     []Player
	rng         *math2d.Rand
	targetX     float32
	targetY     float32
	countdown   *gosprite64.Timer
	roundTimer  *gosprite64.Timer
	spawnTimer  *gosprite64.RepeatingTimer
	blinker     *gosprite64.RepeatingTimer
	started     bool
}

func (s *PlayState) Enter() {
	s.rng = math2d.NewRand(42)
	s.countdown = gosprite64.NewTimer(180) // 3 second countdown
	s.roundTimer = gosprite64.NewTimer(600) // 10 second round
	s.spawnTimer = gosprite64.NewRepeatingTimer(90) // new target every 1.5s
	s.blinker = gosprite64.NewRepeatingTimer(30)

	s.players = make([]Player, s.playerCount)
	s.players[0] = Player{x: 60, y: 50, port: 0}
	if s.playerCount > 1 {
		s.players[1] = Player{x: 60, y: 50, port: 1}
	}
	s.spawnTarget()
}

func (s *PlayState) spawnTarget() {
	s.targetX = s.rng.RangeFloat32(10, 120)
	s.targetY = s.rng.RangeFloat32(20, 80)
}

func (s *PlayState) Update() {
	s.blinker.Tick()

	if !s.countdown.Done() {
		s.countdown.Tick()
		return
	}

	if !s.started {
		s.started = true
	}

	if s.roundTimer.Tick() {
		s.sm.Switch(&ResultState{sm: s.sm, players: s.players})
		return
	}

	if s.spawnTimer.Tick() {
		s.spawnTarget()
	}

	for i := range s.players {
		p := &s.players[i]
		speed := float32(2)

		if gosprite64.PlayerButtonDown(p.port, gosprite64.ButtonDPadUp) {
			p.y -= speed
		}
		if gosprite64.PlayerButtonDown(p.port, gosprite64.ButtonDPadDown) {
			p.y += speed
		}
		if gosprite64.PlayerButtonDown(p.port, gosprite64.ButtonDPadLeft) {
			p.x -= speed
		}
		if gosprite64.PlayerButtonDown(p.port, gosprite64.ButtonDPadRight) {
			p.x += speed
		}
		p.x = math2d.Clamp(p.x, 4, 130)
		p.y = math2d.Clamp(p.y, 14, 96)

		playerRect := math2d.Rect{X: p.x - 3, Y: p.y - 3, W: 6, H: 6}
		targetRect := math2d.Rect{X: s.targetX - 3, Y: s.targetY - 3, W: 6, H: 6}
		if playerRect.Overlaps(targetRect) {
			p.score++
			s.spawnTarget()
		}
	}
}

func (s *PlayState) Draw() {
	gosprite64.ClearScreen()

	panelW := screenW / s.playerCount
	for i := range s.players {
		gosprite64.SetDrawRegion(i*panelW, 0, panelW, screenH)
		s.drawPanel(&s.players[i], i)
		gosprite64.ResetDrawRegion()
	}

	if s.playerCount == 2 {
		gosprite64.DrawLine(screenW/2, 0, screenW/2, screenH, gosprite64.White)
	}

	if !s.countdown.Done() {
		remaining := 3 - s.countdown.Elapsed()/60
		gosprite64.DrawText(fmt.Sprintf("%d", remaining), 140, 100, gosprite64.Yellow)
	}
}

func (s *PlayState) drawPanel(p *Player, index int) {
	colors := []color.Color{gosprite64.Green, gosprite64.Blue}
	playerColor := colors[index%2]

	gosprite64.DrawRect(2, 12, 134, 104, gosprite64.DarkGray)

	if s.started {
		gosprite64.FillRect(
			int(s.targetX)-3, int(s.targetY)-3,
			int(s.targetX)+3, int(s.targetY)+3,
			gosprite64.Yellow,
		)
	}

	gosprite64.FillRect(
		int(p.x)-3, int(p.y)-3,
		int(p.x)+3, int(p.y)+3,
		playerColor,
	)

	gosprite64.DrawText(fmt.Sprintf("P%d: %d", index+1, p.score), 4, 2, gosprite64.White)

	timeLeft := s.roundTimer.Remaining() / 60
	gosprite64.DrawText(fmt.Sprintf("T:%d", timeLeft), 100, 2, gosprite64.LightGray)
}

func (s *PlayState) Exit() {}

// --- Result State ---

type ResultState struct {
	sm      *gosprite64.StateMachine
	players []Player
	menu    *gosprite64.Menu
}

func (s *ResultState) Enter() {
	s.menu = gosprite64.NewMenu([]gosprite64.MenuItem{
		{Label: "Play Again", OnConfirm: func() {
			s.sm.Switch(&PlayState{sm: s.sm, playerCount: len(s.players)})
		}},
		{Label: "Title", OnConfirm: func() {
			s.sm.Switch(&MenuState{sm: s.sm})
		}},
	})
	s.menu.X = 104
	s.menu.Y = 140
	s.menu.Wrap = true
}

func (s *ResultState) Update() {
	s.menu.HandleInput()
}

func (s *ResultState) Draw() {
	gosprite64.ClearScreenWith(gosprite64.DarkPurple)
	gosprite64.DrawText("TIME UP!", 112, 40, gosprite64.White)

	for i, p := range s.players {
		gosprite64.DrawText(
			fmt.Sprintf("P%d SCORE: %d", i+1, p.score),
			88, 70+i*16, gosprite64.Yellow,
		)
	}

	if len(s.players) == 2 {
		if s.players[0].score > s.players[1].score {
			gosprite64.DrawText("P1 WINS!", 112, 110, gosprite64.Green)
		} else if s.players[1].score > s.players[0].score {
			gosprite64.DrawText("P2 WINS!", 112, 110, gosprite64.Blue)
		} else {
			gosprite64.DrawText("TIE!", 128, 110, gosprite64.LightGray)
		}
	}

	s.menu.Draw()
}

func (s *ResultState) Exit() {}

func main() {
	gosprite64.Run(&Game{})
}
