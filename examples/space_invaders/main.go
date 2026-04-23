//go:generate go run github.com/drpaneas/gosprite64/cmd/audiogen -dir .
package main

import (
	"fmt"
	"image/color"

	pigo8 "github.com/drpaneas/gosprite64"
)

const (
	// The library canvas is 288x216, but this example keeps a smaller centered
	// 128x128 gameplay field to show that games can choose a tighter play area.
	logicalCanvasWidth  = 288
	logicalCanvasHeight = 216
	screenW             = 128
	screenH             = 128
	playfieldOffsetX    = (logicalCanvasWidth - screenW) / 2
	playfieldOffsetY    = (logicalCanvasHeight - screenH) / 2

	playerStartX     = 64
	playerStartY     = 120
	playerSpeed      = 2
	bulletSpeed      = 4
	alienSpeed       = 1
	alienBulletSpeed = 2
	initialLives     = 7

	aliensRows  = 5
	aliensCols  = 11
	alienW      = 8
	alienH      = 8
	alienPadX   = 10
	alienPadY   = 10
	alienStartX = 16
	alienStartY = 16
)

// Game state
type (
	bullet struct {
		x, y  int
		speed int
	}
	alien struct {
		x, y   int
		alive  bool
		sprite int
	}
	Game struct {
		// Player
		playerX, playerY int
		lives            int

		// Projectiles
		bullets      []bullet
		alienBullets []bullet

		// Aliens
		aliens []alien

		// Game state
		score    int
		gameOver bool
		paused   bool
		menuItem int // 0 = resume, 1 = quit
	}
)

// Init initializes the game state
func (g *Game) Init() {
	g.resetGame()
}

func (g *Game) resetGame() {
	g.playerX = playerStartX
	g.playerY = playerStartY
	g.lives = initialLives
	g.score = 0
	g.gameOver, g.paused = false, false
	g.menuItem = 0
	g.bullets = g.bullets[:0]
	g.alienBullets = g.alienBullets[:0]
	g.initAliens()
}

// ---- Aliens ----
func (g *Game) initAliens() {
	g.aliens = g.aliens[:0] // Clear previous wave
	for row := 0; row < aliensRows; row++ {
		for col := 0; col < aliensCols; col++ {
			g.aliens = append(g.aliens, alien{
				x:      alienStartX + col*alienPadX,
				y:      alienStartY + row*alienPadY,
				alive:  true,
				sprite: (row % 3) + 1, // Cycle through 3 alien sprites
			})
		}
	}
}

// ---- Input Processing ----
func (g *Game) processInputs() bool {
	if g.gameOver {
		return g.handleGameOverInput()
	}
	// Normal controls
	g.handlePlayerMovement()
	g.handlePlayerShooting()
	return true
}

func (g *Game) handleGameOverInput() bool {
	if pigo8.Btnp(pigo8.O) {
		g.resetGame()
	}
	return false
}

func (g *Game) handlePlayerMovement() {
	if pigo8.Btn(pigo8.LEFT) && g.playerX > 8 {
		g.playerX -= playerSpeed
	}
	if pigo8.Btn(pigo8.RIGHT) && g.playerX < screenW-8 {
		g.playerX += playerSpeed
	}
}

func (g *Game) handlePlayerShooting() {
	if pigo8.Btnp(pigo8.O) {
		g.bullets = append(g.bullets, bullet{
			x:     g.playerX,
			y:     g.playerY - 8,
			speed: bulletSpeed,
		})
	}
}

// Update handles game logic each frame
func (g *Game) Update() {
	if !g.processInputs() {
		return // Skip update if paused or game over
	}
	g.updatePlayerBullets()
	g.updateAliensAndBullets()
	if !g.gameOver {
		g.handleCollisions()
	}
}

func (g *Game) updatePlayerBullets() {
	dst := g.bullets[:0]
	for _, b := range g.bullets {
		b.y -= b.speed
		if b.y >= 0 {
			dst = append(dst, b)
		}
	}
	g.bullets = dst
}

func (g *Game) updateAliensAndBullets() {
	// Simple AI: Aliens shoot randomly
	for i := range g.aliens {
		a := &g.aliens[i]
		if !a.alive {
			continue
		}
		if pigo8.Rnd(100) == 0 {
			g.alienBullets = append(g.alienBullets, bullet{
				x:     a.x,
				y:     a.y + alienH,
				speed: alienBulletSpeed,
			})
		}
		if a.y > playerStartY-8 {
			g.gameOver = true
			return
		}
	}
	// Update alien bullets & collisions
	dst := g.alienBullets[:0]
	for _, b := range g.alienBullets {
		b.y += b.speed
		if b.y > screenH {
			continue // Off screen
		}
		if b.x > g.playerX-4 && b.x < g.playerX+8 &&
			b.y > g.playerY-8 && b.y < g.playerY+8 {
			g.lives--
			pigo8.Music(1, false)
			if g.lives <= 0 {
				g.gameOver = true
			}
			continue
		}
		dst = append(dst, b)
	}
	g.alienBullets = dst
}

// ---- Collisions and Win Condition ----
func (g *Game) handleCollisions() {
	dst := g.bullets[:0]
	for _, b := range g.bullets {
		hit := false
		for j := range g.aliens {
			a := &g.aliens[j]
			if a.alive &&
				b.x > a.x-4 && b.x < a.x+alienW &&
				b.y > a.y-8 && b.y < a.y+alienH {
				a.alive = false
				g.score += 10
				pigo8.Music(0, false)
				hit = true
				break
			}
		}
		if !hit {
			dst = append(dst, b)
		}
	}
	g.bullets = dst
	// New wave if all aliens dead
	allDead := true
	for _, a := range g.aliens {
		if a.alive {
			allDead = false
			break
		}
	}
	if allDead {
		g.initAliens()
	}
}

// Draw renders the game elements to the screen each frame
func (g *Game) Draw() {
	pigo8.ClearScreen()
	pigo8.DrawRect(
		playfieldOffsetX,
		playfieldOffsetY,
		playfieldOffsetX+screenW-1,
		playfieldOffsetY+screenH-1,
		pigo8.DarkGray,
	)
	g.drawPlayer()
	g.drawBullets()
	g.drawAliens()
	g.drawUI()
	if g.gameOver {
		g.drawGameOver()
	}
}

func (g *Game) drawPlayer() {
	drawX := playfieldOffsetX + g.playerX
	drawY := playfieldOffsetY + g.playerY

	// Triangle shape
	pigo8.Line(drawX+4, drawY-8, drawX, drawY, pigo8.White)
	pigo8.Line(drawX+4, drawY-8, drawX+8, drawY, pigo8.White)
	pigo8.Line(drawX, drawY, drawX+8, drawY, pigo8.White)
}

func (g *Game) drawBullets() {
	for _, b := range g.bullets {
		pigo8.Rectfill(
			playfieldOffsetX+b.x,
			playfieldOffsetY+b.y,
			playfieldOffsetX+b.x+2,
			playfieldOffsetY+b.y+4,
			pigo8.White,
		)
	}
	for _, b := range g.alienBullets {
		pigo8.Rectfill(
			playfieldOffsetX+b.x,
			playfieldOffsetY+b.y,
			playfieldOffsetX+b.x+2,
			playfieldOffsetY+b.y+4,
			pigo8.Red,
		)
	}
}

var alienColors = []color.Color{
	pigo8.DarkPurple, // sprite 1
	pigo8.DarkGreen,  // sprite 2
	pigo8.Brown,      // sprite 3
}

func (g *Game) drawAliens() {
	for _, a := range g.aliens {
		if !a.alive {
			continue
		}
		c := alienColors[(a.sprite-1)%len(alienColors)]
		drawX := playfieldOffsetX + a.x
		drawY := playfieldOffsetY + a.y
		pigo8.Rectfill(drawX, drawY, drawX+alienW, drawY+alienH, c)
		pigo8.Rectfill(drawX+2, drawY+2, drawX+6, drawY+6, c)
	}
}

func (g *Game) drawUI() {
	pigo8.Print(fmt.Sprintf("score: %d", g.score), playfieldOffsetX, 12, pigo8.White)
	pigo8.Print(fmt.Sprintf("lives: %d", g.lives), playfieldOffsetX+72, 12, pigo8.White)
}

func (g *Game) drawGameOver() {
	pigo8.Print("GAME OVER", playfieldOffsetX+24, playfieldOffsetY+44, pigo8.White)
	pigo8.Print("PRESS O TO RESTART", playfieldOffsetX+8, playfieldOffsetY+64, pigo8.White)
}

func main() {
	pigo8.Run(&Game{})
}
