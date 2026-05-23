# Draw Text

Use this when you want to put quick built-in text on screen for a label, debug line, or title.

```go
func (g *Game) Draw() {
	gosprite64.ClearScreenWith(gosprite64.Black)
	gosprite64.DrawText("HELLO", 16, 24, gosprite64.White)
}
```

`DrawText` uses the built-in 8x8 font, so this is the fastest way to get readable text on screen.
