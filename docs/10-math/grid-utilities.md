# Grid Utilities

The `math2d` package provides a generic 2D grid for tile-based collision, puzzle logic, and spatial queries.

```go
import "github.com/drpaneas/gosprite64/math2d"
```

## Grid[T]

`Grid[T]` is a generic two-dimensional array of cells. The type parameter `T` must be `comparable`. Addressing is column-first: `Get(col, row)` where column is the X axis and row is the Y axis.

### Creating a grid

```go
func NewGrid[T comparable](cols, rows int) *Grid[T]
```

```go
// A 20x15 grid of integers (all cells start at zero value)
grid := math2d.NewGrid[int](20, 15)

// A 10x10 grid of booleans
walls := math2d.NewGrid[bool](10, 10)
```

### Dimensions

```go
grid.Cols()  // number of columns (X)
grid.Rows()  // number of rows (Y)
```

## Reading and writing cells

### Get

Returns the cell value at (col, row). Returns the zero value of `T` if the position is out of bounds.

```go
value := grid.Get(5, 3)
```

### Set

Writes a value to (col, row). No-op if out of bounds.

```go
grid.Set(5, 3, 42)
```

### InBounds

Checks whether a position is inside the grid:

```go
if grid.InBounds(col, row) {
    // safe to read or write
}
```

You do not need to call `InBounds` before `Get` or `Set` - they handle out-of-bounds access gracefully. Use `InBounds` when you need the check result for game logic.

## Bulk operations

### Fill

Sets every cell to the given value:

```go
grid.Fill(1)  // all cells become 1
```

### Clear

Resets every cell to the zero value of `T`:

```go
grid.Clear()  // all cells become 0 (for int), false (for bool), etc.
```

## Searching

### FindAll

Returns the positions of all cells matching a predicate:

```go
// Find all cells with value 3
cells := grid.FindAll(func(v int) bool { return v == 3 })
for _, cell := range cells {
    fmt.Printf("Found at col=%d, row=%d\n", cell.Col, cell.Row)
}
```

Each result is a `GridCell`:

```go
type GridCell struct {
    Col, Row int
}
```

### CountValue

Counts how many cells equal a specific value:

```go
solidCount := grid.CountValue(1)
emptyCount := grid.CountValue(0)
```

## Neighbors

### Neighbors4

Returns the values of the four orthogonal neighbors (up, down, left, right) that are within bounds:

```go
neighbors := grid.Neighbors4(5, 5)
// Up to 4 values; fewer at edges/corners
```

Corner cells return 2 neighbors, edge cells return 3, and interior cells return 4.

## Row and column scanning

`ScanRow` and `ScanCol` find consecutive runs of cells with the same non-zero group value. You provide a grouping function that maps cell values to group identifiers; cells returning 0 are treated as empty and break runs.

### ScanRow

Scans a row left-to-right:

```go
type Run struct {
    Start  int  // starting column (for rows) or row (for columns)
    Length int  // number of consecutive cells
    Value  int  // the group identifier
}
```

```go
// Find horizontal runs of matching colors
runs := grid.ScanRow(3, func(v int) int { return v })
for _, r := range runs {
    if r.Length >= 3 {
        // Three or more in a row - clear them
    }
}
```

### ScanCol

Same as `ScanRow` but scans a column top-to-bottom:

```go
runs := grid.ScanCol(5, func(v int) int { return v })
```

## Tile collision map example

A common use case is building a collision grid from a tile map, where solid tiles are marked `true` and empty tiles are `false`:

```go
const (
    tileSize = 16
    mapCols  = 20
    mapRows  = 15
)

func buildCollisionGrid(tileData []int) *math2d.Grid[bool] {
    grid := math2d.NewGrid[bool](mapCols, mapRows)
    for row := 0; row < mapRows; row++ {
        for col := 0; col < mapCols; col++ {
            tileID := tileData[row*mapCols+col]
            if tileID > 0 {
                grid.Set(col, row, true)
            }
        }
    }
    return grid
}

func isSolidAt(grid *math2d.Grid[bool], worldX, worldY float32) bool {
    col := int(worldX) / tileSize
    row := int(worldY) / tileSize
    return grid.Get(col, row)
}

func (g *Game) Update() {
    // Check the tile the player is about to move into
    nextX := g.player.X + g.player.VelX
    nextY := g.player.Y + g.player.VelY

    if isSolidAt(g.collisionGrid, nextX, g.player.Y) {
        g.player.VelX = 0
    }
    if isSolidAt(g.collisionGrid, g.player.X, nextY) {
        g.player.VelY = 0
    }

    g.player.X += g.player.VelX
    g.player.Y += g.player.VelY
}
```

## Puzzle game example

Grids are also useful for match-three or Dr. Mario-style puzzle games. Use `ScanRow` and `ScanCol` to find matches:

```go
const (
    colorRed   = 1
    colorBlue  = 2
    colorGreen = 3
)

func findMatches(grid *math2d.Grid[int], minRun int) []math2d.GridCell {
    matched := math2d.NewGrid[bool](grid.Cols(), grid.Rows())

    identity := func(v int) int { return v }

    for row := 0; row < grid.Rows(); row++ {
        for _, run := range grid.ScanRow(row, identity) {
            if run.Length >= minRun {
                for c := run.Start; c < run.Start+run.Length; c++ {
                    matched.Set(c, row, true)
                }
            }
        }
    }
    for col := 0; col < grid.Cols(); col++ {
        for _, run := range grid.ScanCol(col, identity) {
            if run.Length >= minRun {
                for r := run.Start; r < run.Start+run.Length; r++ {
                    matched.Set(col, r, true)
                }
            }
        }
    }

    return matched.FindAll(func(v bool) bool { return v })
}
```
