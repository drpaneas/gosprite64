# Load a Sprite

Use this when you already have a compiled `.sheet` file and want to draw one frame from it.

Load it once during setup:

```go
type Game struct {
	hero *gosprite64.SpriteSheet
}

func (g *Game) Init() {
	sheet, err := gosprite64.LoadSpriteSheet("assets/hero.sheet")
	if err != nil {
		panic(err)
	}
	g.hero = sheet
}
```

Then draw it each frame:

```go
func (g *Game) Draw() {
	gosprite64.ClearScreenWith(gosprite64.Black)
	gosprite64.DrawSprite(g.hero, 0, 32, 48)
}
```

If you still need to create `assets/hero.sheet`, follow the beginner journey sprite step first, then come back here for the quick reminder.
