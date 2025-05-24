//go:generate go run github.com/drpaneas/gosprite64/cmd/audiogen -dir .

package main

import (
	"fmt"

	. "github.com/drpaneas/gosprite64"
)

// Court boundaries
const (
	courtLeft   = 15
	courtRight  = 296
	courtTop    = 10
	courtBottom = 229
	centerX     = (courtRight + courtLeft) / 2
	centerY     = (courtBottom + courtTop) / 2
	lineLen     = 4
)

// Paddle represents a player or computer paddle
type Paddle struct {
	x, y, width, height, speed float64
	color                      int
}

// Ball holds position, velocity, and rendering info
type Ball struct {
	x, y, size           float64
	dx, dy, speed, boost float64
	color                int
}

// Game encapsulates all game state
type Game struct {
	player        Paddle
	computer      Paddle
	ball          Ball
	playerScore   int
	computerScore int
	Scored        string
}

// Init initializes the game state with default paddle and ball positions
func (g *Game) Init() {
	difficulty := 1.5
	paddleHeight := 20.0
	paddleWidth := 4.0
	g.player = Paddle{courtLeft + paddleWidth*2, centerY + paddleHeight/2, paddleWidth, paddleHeight, 1.0 * difficulty, Blue}
	g.computer = Paddle{courtRight - paddleWidth*3, centerY + paddleHeight/2, paddleWidth, paddleHeight, 0.75 * difficulty, Red}
	ballDy := float64(Flr(Rnd(2))) - 0.5
	g.ball = Ball{x: centerX, y: centerY, size: 2, color: White, dx: 1.0 * difficulty, dy: ballDy, speed: 1.0 * difficulty, boost: 0.05 * difficulty}

	// sound
	switch g.Scored {
	case "Player":
		Music(3, false)
	case "Computer":
		Music(4, false)
	default:
		Music(5, false)
	}
}

// Update handles game logic each frame including input, AI, collisions and scoring
func (g *Game) Update() {
	// Player input
	if Btn(UP) && g.player.y > courtTop+1 {
		g.player.y -= g.player.speed
	}
	if Btn(DOWN) && g.player.y+g.player.height < courtBottom-1 {
		g.player.y += g.player.speed
	}

	// Simple AI: track ball when it's moving toward computer
	mid := g.computer.y + g.computer.height/2
	if g.ball.dx > 0 {
		if mid > g.ball.y && g.computer.y > courtTop+1 {
			g.computer.y -= g.computer.speed
		}
		if mid < g.ball.y && g.computer.y+g.computer.height < courtBottom-1 {
			g.computer.y += g.computer.speed
		}
	} else {
		// return to center
		if mid > ((centerY + g.player.height/2) + g.player.height) {
			g.computer.y -= g.computer.speed
		}
		if mid < ((centerY + g.player.height/2) - g.player.height) {
			g.computer.y += g.computer.speed
		}
	}

	// Collisions
	// 1. Ball vs paddles
	if collide(g.ball, g.computer) {
		g.ball.dx = -(g.ball.dx + g.ball.boost)
		Music(0, false)
	}
	if collide(g.ball, g.player) {
		// adjust dy if player changes paddle angle
		if Btn(UP) || Btn(DOWN) {
			g.ball.dy += sign(g.ball.dy) * g.ball.boost * 2
		}
		g.ball.dx = -(g.ball.dx - g.ball.boost)
		Music(1, false)
	}

	// 2. Ball vs top/bottom
	if g.ball.y <= courtTop+1 || g.ball.y+g.ball.size >= courtBottom-1 {
		g.ball.dy = -g.ball.dy
		Music(2, false)
	}

	// 3. Ball vs Walls (aka scoring)
	if g.ball.x > courtRight {
		g.playerScore++
		g.Scored = "Player"
		g.Init()
	}
	if g.ball.x < courtLeft {
		g.computerScore++
		g.Scored = "Computer"
		g.Init()
	}

	// Move ball
	g.ball.x += g.ball.dx
	g.ball.y += g.ball.dy
}

// Draw renders the game elements to the screen each frame
func (g *Game) Draw() {
	ClearScreen(0)

	// Court outline
	DrawRect(courtLeft, courtTop, courtRight, courtBottom, Pico8Palette[White])

	// Center dashed line
	for y := courtTop; y < courtBottom; y += lineLen * 2 {
		// p8.Line(centerX, float64(y), centerX, float64(y+lineLen), 5)
		Line(centerX, y, centerX, y+lineLen, Pico8Palette[White])
		// Rectfill(centerX, y, centerX, y+lineLen, Pico8Palette[White])
	}

	// Ball and paddles
	DrawRectFill(int(g.ball.x), int(g.ball.y), int(g.ball.x+g.ball.size), int(g.ball.y+g.ball.size), Pico8Palette[g.ball.color])
	DrawRectFill(int(g.player.x), int(g.player.y), int(g.player.x+g.player.width), int(g.player.y+g.player.height), Pico8Palette[g.player.color])
	DrawRectFill(int(g.computer.x), int(g.computer.y), int(g.computer.x+g.computer.width), int(g.computer.y+g.computer.height), Pico8Palette[g.computer.color])

	// Scores
	Print(fmt.Sprint(g.playerScore), centerX/2, 2, Yellow)
	Print(fmt.Sprint(g.computerScore), centerX+centerX/2, 2, Yellow)
}

// collide checks axis-aligned collision between ball and paddle
func collide(b Ball, p Paddle) bool {
	return b.x+b.size >= p.x && b.x <= p.x+p.width &&
		b.y+b.size >= p.y && b.y <= p.y+p.height
}

// sign returns the sign of a float, with 0 treated as +1
func sign(v float64) float64 {
	if v < 0 {
		return -1
	}
	return 1
}

func main() {
	Run(&Game{})
}
