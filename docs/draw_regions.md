# Draw Regions (Split-Screen)

Draw regions let you restrict all drawing to a sub-rectangle of the screen. This is how you implement split-screen multiplayer - each player gets their own viewport where coordinates start at (0, 0).

## Basic usage

```go
// Draw player 1's view in the left half
gosprite64.SetDrawRegion(0, 0, 144, 216)
drawGameForPlayer(0)
gosprite64.ResetDrawRegion()

// Draw player 2's view in the right half
gosprite64.SetDrawRegion(144, 0, 144, 216)
drawGameForPlayer(1)
gosprite64.ResetDrawRegion()
```

Inside `SetDrawRegion`, all coordinates are local to the region. Drawing at (0, 0) puts pixels at the region's top-left corner, not the screen's.

On N64 hardware, `SetDrawRegion` also sets the RDP scissor rectangle, so any drawing that falls outside the region is clipped at the hardware level.

## Two-player layout

The 288x216 canvas divides naturally for split-screen:

```
2-player side-by-side:  SetDrawRegion(0, 0, 144, 216)  // left
                        SetDrawRegion(144, 0, 144, 216) // right

2-player top-bottom:    SetDrawRegion(0, 0, 288, 108)   // top
                        SetDrawRegion(0, 108, 288, 108) // bottom

4-player quadrants:     SetDrawRegion(0, 0, 144, 108)     // top-left
                        SetDrawRegion(144, 0, 144, 108)   // top-right
                        SetDrawRegion(0, 108, 144, 108)   // bottom-left
                        SetDrawRegion(144, 108, 144, 108) // bottom-right
```

## Drawing a divider

Draw region doesn't prevent drawing outside it after `ResetDrawRegion`. Draw divider lines after resetting:

```go
gosprite64.SetDrawRegion(0, 0, 144, 216)
drawPlayer1()
gosprite64.ResetDrawRegion()

gosprite64.SetDrawRegion(144, 0, 144, 216)
drawPlayer2()
gosprite64.ResetDrawRegion()

// Divider line in full-screen space
gosprite64.DrawLine(144, 0, 144, 216, gosprite64.White)
```

## Nesting

Calls can be nested. Each `SetDrawRegion` pushes onto a stack, and `ResetDrawRegion` pops back to the previous region:

```go
gosprite64.SetDrawRegion(0, 0, 144, 216)    // player 1 half
gosprite64.SetDrawRegion(10, 10, 124, 90)   // HUD sub-area within player 1
drawHUD()
gosprite64.ResetDrawRegion()                 // back to player 1 full half
drawGameplay()
gosprite64.ResetDrawRegion()                 // back to full screen
```

## DrawRegion type

You can also work with `DrawRegion` values directly for coordinate math:

```go
region := gosprite64.DrawRegion{X: 50, Y: 30, W: 100, H: 80}

// Convert local coords to screen coords
screenX, screenY := region.Offset(10, 20)  // (60, 50)

// Check if a local point is within the region
region.ContainsPoint(10, 20)  // true
region.ContainsPoint(200, 0)  // false

// Clip and offset a rectangle
x1, y1, x2, y2, ok := region.Clip(0, 0, 150, 150)
// ok=true, coords clamped to region bounds

// Check if region is active (non-zero size)
region.Active()  // true
gosprite64.DrawRegion{}.Active()  // false
```

## Dr. Mario example

Dr. Mario 64 shows 2-4 game boards. With draw regions, each board draws at local coordinates:

```go
func drawBoard(playerIndex int, board *GameBoard) {
    // Each board is 64x136 pixels (8 cols x 17 rows of 8x8 cells)
    boardW := 8 * 8
    boardH := 17 * 8

    // Position boards across the screen
    positions := []struct{ x, y int }{
        {20, 20}, {152, 20},           // 2 players
    }

    pos := positions[playerIndex]
    gosprite64.SetDrawRegion(pos.x, pos.y, boardW, boardH)

    // All drawing is now local to the board
    for row := 0; row < 17; row++ {
        for col := 0; col < 8; col++ {
            cell := board.Get(col, row)
            if cell != Empty {
                drawCell(col*8, row*8, cell)
            }
        }
    }

    gosprite64.ResetDrawRegion()
}
```

## Complete example

See `examples/splitscreen_demo` for a working 1P/2P split-screen game using draw regions, timers, and menus together.
