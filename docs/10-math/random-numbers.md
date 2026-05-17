# Random Numbers

The `math2d` package provides a deterministic random number generator. Import it with:

```go
import "github.com/drpaneas/gosprite64/math2d"
```

## Rand

`Rand` is a seedable pseudo-random number generator using the xoshiro128** algorithm. It is deterministic: the same seed always produces the same sequence. This is important for gameplay that needs to be reproducible - replays, netplay, or consistent procedural generation.

Unlike Go's `math/rand`, each `Rand` instance has its own state. There is no global generator to worry about.

```go
rng := math2d.NewRand(42)  // seed with 42
```

### Integers

```go
rng.Uint32()       // raw 32-bit value
rng.Intn(6)        // [0, 6) - like a die roll 0..5
rng.RangeInt(1, 7) // [1, 7) - die roll 1..6
```

### Floats

```go
rng.Float32()                  // [0.0, 1.0)
rng.RangeFloat32(0.5, 2.0)    // [0.5, 2.0)
```

### Bool

```go
if rng.Bool() {
    // roughly 50/50
}
```

### Re-seeding

Call `Seed` to reset the generator to a known state:

```go
rng := math2d.NewRand(42)
first := rng.Uint32()

rng.Seed(42)
second := rng.Uint32()
// first == second
```

### Typical game usage

```go
func (g *Game) Init() {
    g.rng = math2d.NewRand(12345)

    // Spawn enemies at random positions
    for i := 0; i < 10; i++ {
        x := g.rng.RangeFloat32(20, 268)
        y := g.rng.RangeFloat32(20, 196)
        spawnEnemy(x, y)
    }

    // Pick a random color
    colors := []color.Color{gosprite64.Red, gosprite64.Blue, gosprite64.Green}
    c := colors[g.rng.Intn(len(colors))]
}
```
