# Step 10: Score Display

Draw a HUD overlay showing the player's score using DrawText.

## What you will learn

- Drawing HUD text that stays fixed on screen
- Using `fmt.Sprintf` to format dynamic values
- Detecting single button presses with `IsButtonJustPressed`
- The difference between screen-space HUD and world-space gameplay

## What changed from Step 9

Two additions to `PlayState`:

1. A `score int` field that tracks points
2. Score increment on the A button and a `DrawText` call to display it

### Score tracking

```go
type PlayState struct {
	// ... existing fields ...
	score int
}
```

In `Update`, after the movement code:

```go
if gosprite64.IsButtonJustPressed(gosprite64.ButtonA) {
	s.score++
}
```

`IsButtonJustPressed` returns true only on the frame the button goes from released to pressed. Unlike `IsButtonDown` (which is true every frame the button is held), this fires once per press. This prevents the score from racing up while the button is held.

### Drawing the score

At the end of `Draw`, before the transition overlay:

```go
gosprite64.DrawText(fmt.Sprintf("SCORE:%d", s.score), 2, 2, gosprite64.White)
```

This requires adding `"fmt"` to the imports.

`DrawText` always uses screen coordinates, not world coordinates. The text stays pinned to the top-left corner of the display regardless of where the camera is. This is what makes it a HUD element rather than a world element.

## The complete Draw function

```go
func (s *PlayState) Draw() {
	gosprite64.ClearScreen()

	s.scene.Draw(s.camera)

	frame := s.player.Frame()
	gosprite64.DrawWorldSpriteWithOptions(s.charSS, frame, s.playerX, s.playerY, s.camera, gosprite64.DrawSpriteOptions{
		FlipH: s.flipH,
	})

	gosprite64.DrawText(fmt.Sprintf("SCORE:%d", s.score), 2, 2, gosprite64.White)

	if s.fade != nil {
		s.fade.Draw()
	}
}
```

The draw order matters:

1. Clear the screen
2. Draw the tile world (scrolls with camera)
3. Draw the player sprite (world coordinates, scrolls with camera)
4. Draw the score text (screen coordinates, stays fixed)
5. Draw the transition overlay (on top of everything)

### Screen coordinates vs world coordinates

| Function | Coordinate space | Scrolls with camera? |
|----------|-----------------|---------------------|
| `DrawWorldSprite` | World | Yes |
| `DrawWorldSpriteWithOptions` | World | Yes |
| `scene.Draw(camera)` | World | Yes |
| `DrawText` | Screen | No |
| `DrawSprite` | Screen | No |

HUD elements like scores, health bars, and menu text use screen coordinates so they stay in place. Game objects use world coordinates so they move with the camera.

### Formatting tips

The built-in font is monospaced, so columns align naturally. Some useful patterns:

```go
gosprite64.DrawText(fmt.Sprintf("SCORE:%04d", s.score), 2, 2, gosprite64.White)
gosprite64.DrawText(fmt.Sprintf("HP:%d/%d", hp, maxHP), 2, 12, gosprite64.Red)
gosprite64.DrawText(fmt.Sprintf("TIME:%02d:%02d", min, sec), 200, 2, gosprite64.Yellow)
```

`%04d` gives zero-padded display like `SCORE:0001`. Each character is 8 pixels wide in the default font, so you can calculate positions precisely.

## Build and run

```bash
go generate ./examples/platformer
GOENV=n64.env go1.24.5-embedded build -o examples/platformer/game.elf ./examples/platformer
```

The gameplay screen now shows "SCORE:0" in the top-left corner. Press A to increment the score. The text stays fixed on screen while the world scrolls underneath. In a real game you would increment the score on meaningful events like collecting items or defeating enemies rather than on a button press.
